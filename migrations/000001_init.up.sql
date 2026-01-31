CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    user_id UUID NOT NULL,
    service_name TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price > 0),

    start_date DATE NOT NULL,
    end_date DATE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CHECK (end_date IS NULL OR end_date >= start_date)
);

CREATE INDEX idx_subscriptions_user_id
    ON subscriptions(user_id);

CREATE INDEX idx_subscriptions_service_name
    ON subscriptions(service_name);

CREATE INDEX idx_subscriptions_period
    ON subscriptions(start_date, end_date);
