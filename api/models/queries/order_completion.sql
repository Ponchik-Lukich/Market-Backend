CREATE TABLE IF NOT EXISTS order_completion
(
    id            SERIAL PRIMARY KEY,
    courier_id    INTEGER   NOT NULL,
    order_id      INTEGER   NOT NULL,
    complete_time TIMESTAMP NOT NULL,
    FOREIGN KEY (courier_id) REFERENCES couriers (id),
    FOREIGN KEY (order_id) REFERENCES orders (id)
);