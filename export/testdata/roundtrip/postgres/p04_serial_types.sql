CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(32) NOT NULL UNIQUE,
    name TEXT NOT NULL,
    price NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    weight_kg REAL,
    in_stock BOOLEAN DEFAULT TRUE,
    released DATE,
    attrs JSON,
    search tsvector,
    created TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE product_images (
    id BIGSERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products ON DELETE CASCADE,
    url TEXT NOT NULL,
    sort_order SMALLINT DEFAULT 0
);
CREATE TABLE categories (
    id smallserial PRIMARY KEY,
    parent_id smallint REFERENCES categories(id),
    name varchar(80) NOT NULL
);
