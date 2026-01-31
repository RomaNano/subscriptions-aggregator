package main

// @title           Subscriptions Aggregator API
// @version         1.0
// @description     REST API for managing subscriptions and calculating totals.
// @description     Test assignment for Golang developer position.

// @contact.name    RomaNano
// @contact.url     https://github.com/RomaNano

// @host            localhost:8080
// @BasePath        /api/v1

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/RomaNano/subscriptions-aggregator/docs"

	"github.com/RomaNano/subscriptions-aggregator/internal/config"
	"github.com/RomaNano/subscriptions-aggregator/internal/handlers"
	"github.com/RomaNano/subscriptions-aggregator/internal/httpserver"
	"github.com/RomaNano/subscriptions-aggregator/internal/repo"
	"github.com/RomaNano/subscriptions-aggregator/internal/service"
)

func main() {
	// ---------- config ----------
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// ---------- logger ----------
	logger := httpserver.NewLogger(cfg.Log.Level, cfg.Log.Format)
	slog.SetDefault(logger)

	// ---------- database ----------
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
		cfg.DB.SSLMode,
	)

	ctx := context.Background()

	pg, err := repo.NewPostgres(
		ctx,
		dsn,
		cfg.DB.MaxOpenConns,
		cfg.DB.MaxIdleConns,
		cfg.DB.ConnMaxLifetime,
	)
	if err != nil {
		logger.Error("db init failed", "err", err)
		os.Exit(1)
	}
	defer pg.DB.Close()

	// ---------- repositories ----------
	subRepo := repo.NewSubscriptionPostgres(pg.DB)

	// ---------- services ----------
	subService := service.NewSubscriptionService(subRepo)

	// ---------- handlers ----------
	subHandler := handlers.NewSubscriptionHandler(subService)
	totalHandler := handlers.NewTotalHandler(subService)

	// ---------- gin ----------
	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	// ---------- health ----------
	r.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second)
		defer cancel()

		if err := pg.DB.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "down"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// ---------- swagger (ВАЖНО: НЕ внутри /api/v1) ----------
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ---------- API v1 ----------
	api := r.Group("/api/v1")
	{
		api.POST("/subscriptions", subHandler.Create)
		api.GET("/subscriptions/:id", subHandler.GetByID)
		api.PUT("/subscriptions/:id", subHandler.Update)
		api.DELETE("/subscriptions/:id", subHandler.Delete)
		api.GET("/subscriptions", subHandler.List)

		api.GET("/subscriptions/total", totalHandler.Get)
	}

	// ---------- http server ----------
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:      r,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	// ---------- start ----------
	go func() {
		logger.Info("http server starting",
			"addr", srv.Addr,
			"env", cfg.Env,
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("http server failed", "err", err)
			os.Exit(1)
		}
	}()

	// ---------- graceful shutdown ----------
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", "err", err)
	} else {
		logger.Info("server stopped gracefully")
	}
}
