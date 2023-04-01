CREATE TABLE IF NOT EXISTS orders
(
    id SERIAL PRIMARY KEY NOT NULL UNIQUE ,
    weight FLOAT,
    delivery_district INT,
    delivery_hours TEXT[],
    cost INT,
    courier_id BIGINT NULL,
    assigned BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (courier_id) REFERENCES couriers(id)
);