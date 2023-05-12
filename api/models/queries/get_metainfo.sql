SELECT c.id,
       c.type,
       c.working_areas,
       c.working_hours,
       SUM(o.cost) AS earnings,
       COUNT(o.id) AS completed_orders
FROM couriers c
         JOIN orders o ON o.courier_id = c.id
WHERE o.courier_id = $1
  AND o.complete_time >= $2
  AND o.complete_time < $3
GROUP BY c.id, c.type, c.working_areas, c.working_hours;