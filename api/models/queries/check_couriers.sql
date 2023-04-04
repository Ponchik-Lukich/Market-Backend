SELECT COUNT(*) = 0
FROM unnest($1::bigint[]) AS id_list(ids)
WHERE NOT EXISTS(
        SELECT 1
        FROM couriers
        WHERE couriers.id = id_list.ids
    );