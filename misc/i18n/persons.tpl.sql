\timing off
COPY (
SELECT
  p.id,
  p.uid,
  (SELECT name
   FROM person_i18n
   WHERE person_id = p.id AND language = 'src')  AS "src.name",
  (SELECT description
   FROM person_i18n
   WHERE person_id = p.id AND language = 'src')  AS "src.full_name",
  (SELECT name
   FROM person_i18n
   WHERE person_id = p.id AND language = 'dest') AS "dest.name",
  (SELECT description
   FROM person_i18n
   WHERE person_id = p.id AND language = 'dest') AS "dest.full_name"
FROM persons p
ORDER BY p.id
) TO STDOUT WITH CSV HEADER;
