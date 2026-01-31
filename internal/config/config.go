package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string `yaml:"env" env:"ENV" env-default:"local"`
	HTTP HTTP   `yaml:"http"`
	DB   DB     `yaml:"db"`
	Log  Log    `yaml:"log"`
}

type HTTP struct {
	Host          string        `yaml:"host" env:"HTTP_HOST" env-default:"0.0.0.0"`
	Port          int           `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
	ReadTimeout   time.Duration `yaml:"read_timeout" env:"HTTP_READ_TIMEOUT" env-default:"5s"`
	WriteTimeout  time.Duration `yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT" env-default:"5s"`
	IdleTimeout   time.Duration `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"HTTP_SHUTDOWN_TIMEOUT" env-default:"10s"`
}

type DB struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"db"`
	Port     int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"DB_USER" env-default:"app"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:"app"`
	Name     string `yaml:"name" env:"DB_NAME" env-default:"app"`
	SSLMode  string `yaml:"sslmode" env:"DB_SSLMODE" env-default:"disable"`

	MaxOpenConns    int           `yaml:"max_open_conns" env:"DB_MAX_OPEN_CONNS" env-default:"10"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env:"DB_MAX_IDLE_CONNS" env-default:"5"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env:"DB_CONN_MAX_LIFETIME" env-default:"30m"`
}

type Log struct {
	Level  string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	Format string `yaml:"format" env:"LOG_FORMAT" env-default:"json"`
}

func Load() (*Config, error) {
	path := os.Getenv("APP_CONFIG")
	if path == "" {
		path = "config.yaml"
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}

	return &cfg, nil
}
