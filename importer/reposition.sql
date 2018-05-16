UPDATE collections_content_units
SET position = s.row_number
FROM (select
        cu.id,
        cu.properties ->> 'film_date',
        ccu.position,
        ROW_NUMBER()
        OVER (
          ORDER BY cu.properties ->> 'film_date' ) as row_number
      from
        collections_content_units ccu inner join content_units cu
          on ccu.content_unit_id = cu.id and ccu.collection_id = 10683
      order by cu.properties ->> 'film_date') AS s
WHERE content_unit_id = s.id and collection_id = 10683;