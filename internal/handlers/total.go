package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/RomaNano/subscriptions-aggregator/internal/service"
)

type TotalHandler struct {
	svc service.SubscriptionService
}

func NewTotalHandler(svc service.SubscriptionService) *TotalHandler {
	return &TotalHandler{svc: svc}
}

// GET /subscriptions/total?from=YYYY-MM-01&to=YYYY-MM-01
func (h *TotalHandler) Get(c *gin.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")

	if fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to are required"})
		return
	}

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from"})
		return
	}
	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to"})
		return
	}

	var userID *uuid.UUID
	if v := c.Query("user_id"); v != "" {
		u, err := uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
			return
		}
		userID = &u
	}

	var serviceName *string
	if v := c.Query("service_name"); v != "" {
		serviceName = &v
	}

	total, err := h.svc.Total(c.Request.Context(), service.TotalFilter{
		UserID:      userID,
		ServiceName: serviceName,
		From:        from,
		To:          to,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, TotalResponse{Total: total})
}
