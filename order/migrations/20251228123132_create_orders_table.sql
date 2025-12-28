-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
CREATE TYPE order_status AS ENUM (
            'PENDING_PAYMENT',
            'PAID',
            'CANCELLED'
        );
END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_method') THEN
CREATE TYPE payment_method AS ENUM (
            'PAYMENT_METHOD_UNKNOWN_UNSPECIFIED',
            'PAYMENT_METHOD_CARD',
            'PAYMENT_METHOD_SBP',
            'PAYMENT_METHOD_CREDIT_CARD',
            'PAYMENT_METHOD_INVESTOR_MONEY'
        );
END IF;
END $$;
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS orders (
                                      uuid             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_uuid         UUID NOT NULL,
    part_uuids        UUID[] NOT NULL DEFAULT '{}',
    total_price       NUMERIC(14,2) NOT NULL CHECK (total_price >= 0),
    transaction_uuid  UUID NULL,
    order_status      order_status NOT NULL,
    payment_method    payment_method NULL,

    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_orders_user_uuid ON orders(user_uuid);
CREATE INDEX IF NOT EXISTS idx_orders_transaction_uuid ON orders(transaction_uuid);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(order_status);


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_orders_updated_at()
RETURNS trigger AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_set_updated_at_orders
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION set_orders_updated_at();


-- +goose Down
DROP TABLE IF EXISTS orders;

-- +goose StatementBegin
DO $$
BEGIN
    -- Типы удаляем только если больше нигде не используются (иначе будет ошибка).
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
        DROP TYPE order_status;
    END IF;

    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_method') THEN
        DROP TYPE payment_method;
    END IF;
END $$;
-- +goose StatementEnd

DROP TRIGGER IF EXISTS trg_set_updated_at_orders ON orders;

-- +goose StatementBegin
DROP FUNCTION IF EXISTS set_orders_updated_at();
-- +goose StatementEnd