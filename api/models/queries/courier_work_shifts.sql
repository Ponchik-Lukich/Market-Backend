CREATE TABLE IF NOT EXISTS courier_work_shifts
(
    id         INT PRIMARY KEY,
    courier_id INT,
    shift_time VARCHAR(255),
    FOREIGN KEY (courier_id) REFERENCES couriers (id)
);