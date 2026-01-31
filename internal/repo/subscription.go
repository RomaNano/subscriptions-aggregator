package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/RomaNano/subscriptions-aggregator/internal/domain"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, s *domain.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	Update(ctx context.Context, s *domain.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter ListFilter) ([]domain.Subscription, error)
}


