CREATE TABLE IF NOT EXISTS orders
(
    id SERIAL PRIMARY KEY,
    weight NUMERIC(10, 2) NOT NULL,
    delivery_district INTEGER NOT NULL,
    delivery_time VARCHAR(255) NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    courier_id INTEGER NULL,
    assigned BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (courier_id) REFERENCES couriers(id)
);