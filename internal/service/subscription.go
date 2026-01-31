package service

import (
	"context"

	"github.com/RomaNano/subscriptions-aggregator/internal/domain"
	"github.com/google/uuid"
)

type SubscriptionService interface {
	Create(ctx context.Context, s *domain.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	Update(ctx context.Context, s *domain.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, f ListFilter) ([]domain.Subscription, error)

	// ключевая ручка
	Total(ctx context.Context, f TotalFilter) (int, error)
}
