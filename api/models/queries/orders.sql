CREATE TABLE IF NOT EXISTS orders
(
    order_id       SERIAL PRIMARY KEY,
    weight         REAL    NOT NULL,
    regions        INTEGER NOT NULL,
    delivery_hours TEXT    NOT NULL,
    cost           INTEGER NOT NULL,
    completed_time TIMESTAMP
);