BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Common triggers and functions
CREATE OR REPLACE FUNCTION trigger_sanitize_id()
    RETURNS TRIGGER AS
$$
BEGIN
    IF tg_op = 'INSERT'
    THEN
        new.id = uuid_generate_v4();
    ELSEIF tg_op = 'UPDATE'
    THEN
        new.id = old.id;
    END IF;
    RETURN new;
END;
$$
    LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION sanitize_timestamps()
    RETURNS TRIGGER AS
$$
BEGIN
    IF tg_op = 'INSERT'
    THEN
        new.created_at = now();
        new.updated_at = now();
    ELSEIF tg_op = 'UPDATE'
    THEN
        new.created_at = old.created_at;
        new.updated_at = now();
    END IF;
    RETURN new;
END;
$$
    LANGUAGE plpgsql;

-- User table
CREATE TABLE "user"
(
    id         UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
    username   TEXT UNIQUE CHECK ( length(trim(username)) < 512 AND length(trim(username)) > 0),
    full_name  TEXT CHECK ( length(full_name) < 512),
    password   TEXT        NOT NULL CHECK ( length(password) < 512 AND length(password) > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Session table
CREATE TABLE session
(
    id         UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
    user_id    UUID REFERENCES "user" (id),
    expire_at  TIMESTAMPTZ NOT NULL,
    meta       JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Product table
CREATE TABLE product
(
    id          UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
    name        TEXT        NOT NULL CHECK ( length(name) > 0 AND length(name) < 1024 ),
    price       INTEGER     NOT NULL CHECK ( price > 0 ),
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Cart table
CREATE TABLE cart
(
    id         UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
    user_id    UUID UNIQUE REFERENCES "user" (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Cart items table
CREATE TABLE cart_item
(
    id         UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
    cart_id    UUID REFERENCES cart (id),
    product_id UUID REFERENCES product (id),
    amount     INTEGER CHECK ( amount > 0 ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (cart_id, product_id)
);

-- Cart discount sets table
-- CREATE TABLE cart_discount_set
-- (
--     id         UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
--     cart_id    UUID REFERENCES cart (id),
--     items_set  JSONB,
--     discount   INTEGER,
--     created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
--     updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
-- );

-- Discounts table
CREATE TABLE discount
(
    id         UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
    name       TEXT        NOT NULL CHECK ( length(name) > 0 AND length(name) < 1024 ),
    rule       JSONB       NOT NULL,
    percent    INTEGER     NOT NULL CHECK ( percent > 0 AND percent < 100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Coupon table
CREATE TABLE coupon
(
    id               UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
    code             TEXT UNIQUE NOT NULL CHECK ( length(code) > 5 ),
    name             TEXT        NOT NULL CHECK ( length(name) > 0 AND length(name) < 1024 ),
    cart_id          UUID REFERENCES cart (id),
    expire_at        TIMESTAMPTZ NOT NULL,
    rule             JSONB       NOT NULL,
    discount_percent INTEGER     NOT NULL CHECK ( discount_percent > 0 AND discount_percent < 100),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Order table
CREATE TABLE "order"
(
    id         UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
    user_id    UUID        NOT NULL REFERENCES "user" (id),
    status     TEXT        NOT NULL CHECK ( status = ANY ('{pending, paid}'::TEXT[])),
    total      INTEGER     NOT NULL,
    savings    INTEGER     NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Order table
CREATE TABLE order_item
(
    id         UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
    name       TEXT        NOT NULL,
    order_id   UUID        NOT NULL REFERENCES "order" (id),
    status     TEXT        NOT NULL,
    price      INTEGER     NOT NULL,
    amount     INTEGER     NOT NULL,
    discount   INTEGER     NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Payment table
CREATE TABLE payment
(
    id         UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
    order_id   UUID        NOT NULL REFERENCES "order" (id),
    amount     INTEGER     NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMIT;