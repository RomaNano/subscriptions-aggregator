package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/RomaNano/subscriptions-aggregator/internal/service"
)

func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidData),
		errors.Is(err, service.ErrInvalidPeriod):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}
