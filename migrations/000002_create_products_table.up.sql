CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    price_cents BIGINT NOT NULL CHECK (price_cents >= 0),
    stock BIGINT NOT NULL CHECK (stock >= 0),
    active BOOLEAN NOT NULL DEFAULT TRUE
);