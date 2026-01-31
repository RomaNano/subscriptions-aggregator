package service

import (
	"time"

	"github.com/google/uuid"
)

type ListFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	From        *time.Time
	To          *time.Time

	Limit  int
	Offset int
}

type TotalFilter struct {
	UserID      *uuid.UUID
	ServiceName *string

	From time.Time 
	To   time.Time
}
