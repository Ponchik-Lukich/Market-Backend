CREATE TABLE IF NOT EXISTS couriers
(
    id            INT PRIMARY KEY,
    type          VARCHAR(255),
    working_areas INT[],
    working_hours TEXT[]
);