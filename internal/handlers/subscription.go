package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/RomaNano/subscriptions-aggregator/internal/domain"
	"github.com/RomaNano/subscriptions-aggregator/internal/service"
)

type SubscriptionHandler struct {
	svc service.SubscriptionService
}

func NewSubscriptionHandler(svc service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{svc: svc}
}

// Create creates a new subscription
// @Summary      Create subscription
// @Description  Create a new subscription
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        subscription body CreateSubscriptionRequest true "Subscription data"
// @Success      201 {object} SubscriptionResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /subscriptions [post]
func (h *SubscriptionHandler) Create(c *gin.Context) {
	var req CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	sub := &domain.Subscription{
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	if err := h.svc.Create(c.Request.Context(), sub); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toResponse(sub))
}

// GetByID gets subscription by ID
// @Summary      Get subscription
// @Description  Get subscription by ID
// @Tags         subscriptions
// @Produce      json
// @Param        id path string true "Subscription ID"
// @Success      200 {object} SubscriptionResponse
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	sub, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}
	if sub == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, toResponse(sub))
}

// Update updates subscription by ID
// @Summary      Update subscription
// @Description  Update subscription by ID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id path string true "Subscription ID"
// @Param        subscription body UpdateSubscriptionRequest true "Subscription data"
// @Success      204
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	sub := &domain.Subscription{
		ID:          id,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	if err := h.svc.Update(c.Request.Context(), sub); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Delete deletes subscription by ID
// @Summary      Delete subscription
// @Description  Delete subscription by ID
// @Tags         subscriptions
// @Param        id path string true "Subscription ID"
// @Success      204
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// List lists subscriptions
// @Summary      List subscriptions
// @Description  List subscriptions with filters
// @Tags         subscriptions
// @Produce      json
// @Param        user_id query string false "User ID"
// @Param        service_name query string false "Service name"
// @Param        from query string false "From date (YYYY-MM-01)"
// @Param        to query string false "To date (YYYY-MM-01)"
// @Param        limit query int false "Limit"
// @Param        offset query int false "Offset"
// @Success      200 {array} SubscriptionResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /subscriptions [get]
func (h *SubscriptionHandler) List(c *gin.Context) {
	var (
		userID      *uuid.UUID
		serviceName *string
		from, to    *time.Time
	)

	if v := c.Query("user_id"); v != "" {
		u, err := uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
			return
		}
		userID = &u
	}

	if v := c.Query("service_name"); v != "" {
		serviceName = &v
	}

	if v := c.Query("from"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from"})
			return
		}
		from = &t
	}

	if v := c.Query("to"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to"})
			return
		}
		to = &t
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	subs, err := h.svc.List(c.Request.Context(), service.ListFilter{
		UserID:      userID,
		ServiceName: serviceName,
		From:        from,
		To:          to,
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	resp := make([]SubscriptionResponse, 0, len(subs))
	for i := range subs {
		resp = append(resp, toResponse(&subs[i]))
	}

	c.JSON(http.StatusOK, resp)
}

func toResponse(s *domain.Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:          s.ID,
		UserID:      s.UserID,
		ServiceName: s.ServiceName,
		Price:       s.Price,
		StartDate:   s.StartDate,
		EndDate:     s.EndDate,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
