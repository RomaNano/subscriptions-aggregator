package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/RomaNano/subscriptions-aggregator/internal/domain"
	"github.com/RomaNano/subscriptions-aggregator/internal/repo"
)

var (
	ErrInvalidPeriod = errors.New("invalid period")
	ErrInvalidData   = errors.New("invalid subscription data")
)

type subscriptionService struct {
	repo repo.SubscriptionRepository
}

func NewSubscriptionService(r repo.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: r}
}

func validateSubscription(s *domain.Subscription) error {
	if s.UserID == uuid.Nil {
		return ErrInvalidData
	}
	if s.ServiceName == "" {
		return ErrInvalidData
	}
	if s.Price <= 0 {
		return ErrInvalidData
	}
	if s.EndDate != nil && s.EndDate.Before(s.StartDate) {
		return ErrInvalidData
	}
	return nil
}

func (s *subscriptionService) Create(ctx context.Context, sub *domain.Subscription) error {
	if err := validateSubscription(sub); err != nil {
		return err
	}
	return s.repo.Create(ctx, sub)
}

func (s *subscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *subscriptionService) Update(ctx context.Context, sub *domain.Subscription) error {
	if err := validateSubscription(sub); err != nil {
		return err
	}
	return s.repo.Update(ctx, sub)
}

func (s *subscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *subscriptionService) List(ctx context.Context, f ListFilter) ([]domain.Subscription, error) {
	rf := repo.ListFilter{
		UserID:      f.UserID,
		ServiceName: f.ServiceName,
		From:        f.From,
		To:          f.To,
		Limit:       f.Limit,
		Offset:     f.Offset,
	}
	return s.repo.List(ctx, rf)
}

func firstOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func monthsBetweenInclusive(a, b time.Time) int {
	ay, am := a.Year(), int(a.Month())
	by, bm := b.Year(), int(b.Month())
	return (by-ay)*12 + (bm-am) + 1
}

func (s *subscriptionService) Total(ctx context.Context, f TotalFilter) (int, error) {
	if f.To.Before(f.From) {
		return 0, ErrInvalidPeriod
	}


	rf := repo.ListFilter{
		UserID:      f.UserID,
		ServiceName: f.ServiceName,
		From:        &f.To,   
		To:          &f.From, 
		Limit:       0,
		Offset:     0,
	}

	subs, err := s.repo.List(ctx, rf)
	if err != nil {
		return 0, err
	}

	from := firstOfMonth(f.From)
	to := firstOfMonth(f.To)

	total := 0

	for _, sub := range subs {
		subStart := firstOfMonth(sub.StartDate)

		subEnd := to
		if sub.EndDate != nil {
			se := firstOfMonth(*sub.EndDate)
			if se.Before(subEnd) {
				subEnd = se
			}
		}

		effStart := subStart
		if effStart.Before(from) {
			effStart = from
		}
		effEnd := subEnd
		if effEnd.After(to) {
			effEnd = to
		}

		if effEnd.Before(effStart) {
			continue
		}

		months := monthsBetweenInclusive(effStart, effEnd)
		total += months * sub.Price
	}

	return total, nil
}
