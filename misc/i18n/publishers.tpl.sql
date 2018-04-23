\timing off
COPY (
SELECT
  p.id,
  p.uid,
  (SELECT name
   FROM publisher_i18n
   WHERE publisher_id = p.id AND language = 'src')  AS "src.name",
  (SELECT description
   FROM publisher_i18n
   WHERE publisher_id = p.id AND language = 'src')  AS "src.full_name",
  (SELECT name
   FROM publisher_i18n
   WHERE publisher_id = p.id AND language = 'dest') AS "dest.name",
  (SELECT description
   FROM publisher_i18n
   WHERE publisher_id = p.id AND language = 'dest') AS "dest.full_name"
FROM publishers p
ORDER BY p.id
) TO STDOUT WITH CSV HEADER;
