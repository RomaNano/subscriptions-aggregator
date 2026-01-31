package handlers

import (
	"time"

	"github.com/google/uuid"
)



type CreateSubscriptionRequest struct {
	UserID      uuid.UUID  `json:"user_id" binding:"required"`
	ServiceName string     `json:"service_name" binding:"required"`
	Price       int        `json:"price" binding:"required"`
	StartDate   time.Time  `json:"start_date" binding:"required"` // YYYY-MM-01
	EndDate     *time.Time `json:"end_date"`                       // YYYY-MM-01 or null
}

type UpdateSubscriptionRequest struct {
	ServiceName string     `json:"service_name" binding:"required"`
	Price       int        `json:"price" binding:"required"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date"`
}



type SubscriptionResponse struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TotalResponse struct {
	Total int `json:"total"`
}
