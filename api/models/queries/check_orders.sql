WITH id_list(ids) AS (SELECT unnest($1::bigint[])),
     missing_ids AS (SELECT ids
                     FROM id_list
                     WHERE NOT EXISTS(
                             SELECT 1
                             FROM orders
                             WHERE orders.id = id_list.ids
                         )),
     assigned_ids AS (SELECT ids
                      FROM id_list
                      WHERE EXISTS(
                                    SELECT 1
                                    FROM orders
                                    WHERE orders.id = id_list.ids
                                      AND orders.assigned = true
                                ))
SELECT CASE
           WHEN (SELECT COUNT(*) FROM missing_ids) > 0 THEN 1
           WHEN (SELECT COUNT(*) FROM assigned_ids) > 0 THEN 2
           ELSE 0
           END AS result;