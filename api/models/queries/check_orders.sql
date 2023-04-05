WITH null_courier_orders AS (
    SELECT id
    FROM orders
    WHERE courier_id IS NULL
),
     mismatched_courier_orders AS (
         SELECT o.id
         FROM orders o
                  JOIN unnest($1::bigint[], $2::bigint[]) t(order_id, courier_id)
                       ON o.id = t.order_id
         WHERE o.courier_id <> t.courier_id
     ),
     nonexistent_orders AS (
         SELECT order_id
         FROM unnest($1::bigint[]) order_id
         WHERE NOT EXISTS (SELECT 1 FROM orders o WHERE o.id = order_id)
     )
SELECT
    CASE
        WHEN EXISTS(SELECT 1 FROM null_courier_orders) THEN 1
        WHEN EXISTS(SELECT 1 FROM mismatched_courier_orders) THEN 2
        WHEN EXISTS(SELECT 1 FROM nonexistent_orders) THEN 3
        ELSE 0
        END AS status
FROM orders
LIMIT 1;
