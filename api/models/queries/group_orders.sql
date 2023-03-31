CREATE TABLE IF NOT EXISTS group_orders
(
    courier_id             int  not null,
    assigned_orders text not null,
    total_cost      int  not null
);