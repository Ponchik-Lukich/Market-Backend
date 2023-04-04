SELECT c.id,
       c.type,
       c.working_areas,
       c.working_hours,
       SUM(o.cost) AS earnings,
       COUNT(o.id) AS completed_orders
FROM couriers c
         JOIN order_completion oc ON c.id = oc.courier_id
         JOIN orders o ON o.id = oc.order_id
WHERE oc.courier_id = $1
  AND oc.complete_time >= $2
  AND oc.complete_time < $3
GROUP BY c.id, c.type, c.working_areas, c.working_hours;