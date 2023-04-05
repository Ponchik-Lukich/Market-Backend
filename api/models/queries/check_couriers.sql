SELECT id_list.ids AS non_existing_courier_id
FROM unnest($1::bigint[]) AS id_list(ids)
WHERE NOT EXISTS(
        SELECT 1
        FROM couriers
        WHERE couriers.id = id_list.ids
    );
