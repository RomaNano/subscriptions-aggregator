package repo

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