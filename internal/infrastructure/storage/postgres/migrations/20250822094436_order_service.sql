-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS orders (
    order_uid TEXT PRIMARY KEY,
    track_number TEXT NOT NULL,
    entry TEXT NOT NULL,
    locale VARCHAR(3) NOT NULL,
    internal_signature TEXT,
    customer_id TEXT NOT NULL,
    delivery_service TEXT,
    shardkey TEXT,
    sm_id INT NOT NULL,
    date_created TIMESTAMP default NOW(),
    oof_shard TEXT
);

CREATE TABLE IF NOT EXISTS deliveries (
    order_uid TEXT PRIMARY KEY REFERENCES orders(order_uid),
    del_name TEXT NOT NULL,
    phone TEXT NOT NULL,
    zip TEXT,
    city TEXT,
    address TEXT,
    region TEXT,
    email TEXT
);

CREATE TABLE IF NOT EXISTS payments (
    order_uid TEXT PRIMARY KEY REFERENCES orders(order_uid),
    transaction_id TEXT NOT NULL,
    request_id TEXT,
    currency TEXT,
    provider TEXT,
    amount INT,
    payment_dt BIGINT,
    bank TEXT,
    delivery_cost INT,
    goods_total INT,
    custom_fee INT
);

CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    order_uid TEXT NOT NULL REFERENCES orders(order_uid),
    chrt_id BIGINT,
    track_number TEXT,
    price INT,
    rid TEXT,
    item_name TEXT,
    sale INT,
    item_size TEXT,
    total_price INT,
    nm_id BIGINT,
    brand TEXT,
    status INT
);

CREATE INDEX IF NOT EXISTS idx_items_order_uid ON items(order_uid);
-- CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders(customer_id);
-- CREATE INDEX IF NOT EXISTS idx_orders_date_created ON orders(date_created);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS deliveries;
DROP TABLE IF EXISTS orders;

-- +goose StatementEnd