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

// Get calculates total subscription cost
// @Summary      Get total subscription cost
// @Description  Calculate total cost of subscriptions for a period
// @Tags         subscriptions
// @Produce      json
// @Param        from query string true  "From date (YYYY-MM)"
// @Param        to   query string true  "To date (YYYY-MM)"
// @Param        user_id query string false "User ID"
// @Param        service_name query string false "Service name"
// @Success      200 {object} TotalResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /subscriptions/total [get]
func (h *TotalHandler) Get(c *gin.Context) {
	// --- required params ---
	fromStr := c.Query("from")
	toStr := c.Query("to")

	from, err := time.Parse("2006-01", fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from"})
		return
	}

	to, err := time.Parse("2006-01", toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to"})
		return
	}

	// --- optional user_id ---
	var userID *uuid.UUID
	if v := c.Query("user_id"); v != "" {
		u, err := uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
			return
		}
		userID = &u
	}

	// --- optional service_name ---
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, TotalResponse{
		Total: total,
	})
}
