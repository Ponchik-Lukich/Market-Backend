CREATE TABLE IF NOT EXISTS orders
(
    id SERIAL PRIMARY KEY NOT NULL UNIQUE ,
    weight FLOAT,
    delivery_district INT,
    delivery_hours TEXT[],
    cost INT,
    assigned BOOLEAN NOT NULL DEFAULT FALSE
);