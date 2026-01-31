package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/RomaNano/subscriptions-aggregator/internal/domain"
	"github.com/google/uuid"
)

type SubscriptionPostgres struct {
	db *sql.DB
}

func NewSubscriptionPostgres(db *sql.DB) *SubscriptionPostgres {
	return &SubscriptionPostgres{db: db}
}

func (r *SubscriptionPostgres) Create(ctx context.Context, s *domain.Subscription) error {
	query := `
		INSERT INTO subscriptions
		    (user_id, service_name, price, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		s.UserID,
		s.ServiceName,
		s.Price,
		s.StartDate,
		s.EndDate,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *SubscriptionPostgres) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	query := `
		SELECT id, user_id, service_name, price,
		       start_date, end_date, created_at, updated_at
		FROM subscriptions
		WHERE id = $1
	`

	var s domain.Subscription
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID,
		&s.UserID,
		&s.ServiceName,
		&s.Price,
		&s.StartDate,
		&s.EndDate,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *SubscriptionPostgres) Update(ctx context.Context, s *domain.Subscription) error {
	query := `
		UPDATE subscriptions
		SET service_name = $1,
		    price = $2,
		    start_date = $3,
		    end_date = $4,
		    updated_at = now()
		WHERE id = $5
	`

	res, err := r.db.ExecContext(
		ctx,
		query,
		s.ServiceName,
		s.Price,
		s.StartDate,
		s.EndDate,
		s.ID,
	)
	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *SubscriptionPostgres) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(
		ctx,
		`DELETE FROM subscriptions WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *SubscriptionPostgres) List(ctx context.Context, f ListFilter) ([]domain.Subscription, error) {
	var (
		conds []string
		args  []any
		argN  = 1
	)

	if f.UserID != nil {
		conds = append(conds, fmt.Sprintf("user_id = $%d", argN))
		args = append(args, *f.UserID)
		argN++
	}

	if f.ServiceName != nil {
		conds = append(conds, fmt.Sprintf("service_name = $%d", argN))
		args = append(args, *f.ServiceName)
		argN++
	}

	if f.From != nil {
		conds = append(conds, fmt.Sprintf("start_date <= $%d", argN))
		args = append(args, *f.From)
		argN++
	}

	if f.To != nil {
		conds = append(conds, fmt.Sprintf("(end_date IS NULL OR end_date >= $%d)", argN))
		args = append(args, *f.To)
		argN++
	}

	query := `
		SELECT id, user_id, service_name, price,
		       start_date, end_date, created_at, updated_at
		FROM subscriptions
	`

	if len(conds) > 0 {
		query += " WHERE " + strings.Join(conds, " AND ")
	}

	query += " ORDER BY created_at DESC"

	if f.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argN)
		args = append(args, f.Limit)
		argN++
	}
	if f.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argN)
		args = append(args, f.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.Subscription
	for rows.Next() {
		var s domain.Subscription
		if err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.ServiceName,
			&s.Price,
			&s.StartDate,
			&s.EndDate,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, s)
	}

	return res, rows.Err()
}
