UPDATE orders
SET complete_time = ct.complete_time
FROM unnest($1::bigint[], $2::timestamp[]) ct(order_id, complete_time)
WHERE orders.id = ct.order_id
RETURNING orders.id, orders.weight, orders.delivery_district, orders.delivery_hours, orders.cost, orders.complete_time;
