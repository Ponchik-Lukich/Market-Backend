CREATE TABLE IF NOT EXISTS couriers
(
    courier_id   SERIAL PRIMARY KEY,
    courier_type VARCHAR(255) NOT NULL,
    max_weight   REAL         NOT NULL
);