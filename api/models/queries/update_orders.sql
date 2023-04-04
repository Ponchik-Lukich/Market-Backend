UPDATE orders
SET assigned = true
WHERE id = ANY($1)
RETURNING id, cost, delivery_hours, delivery_district, weight, (
    SELECT complete_time
    FROM order_completion
    WHERE order_id = orders.id
);