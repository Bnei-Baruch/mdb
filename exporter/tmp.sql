\timing off
COPY (
       select cu.properties ->> 'film_date',
              cui.name,
              concat('https://kabbalahmedia.info/he/lessons/cu/', cu.uid) as link
       from content_units cu
                   left join content_unit_i18n cui on cu.id = cui.content_unit_id and cui.language = 'he'
       where id in
             (select distinct cu.id
              from content_units cu
                          inner join content_units_sources cus on cu.id = cus.content_unit_id and cus.source_id = 8
                          inner join files f on cu.id = f.content_unit_id and f.type = 'image')
       order by cu.properties ->> 'film_date'
) TO STDOUT WITH CSV HEADER;
