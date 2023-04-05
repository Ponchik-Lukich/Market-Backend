CREATE TABLE IF NOT EXISTS orders
(
    id SERIAL PRIMARY KEY NOT NULL UNIQUE ,
    weight FLOAT,
    delivery_district INT,
    delivery_hours TEXT[],
    cost INT,
    courier_id INT REFERENCES couriers(id) DEFAULT NULL,
    complete_time TIMESTAMP DEFAULT NULL
);