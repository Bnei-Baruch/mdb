\timing off
COPY (
SELECT
  a.id,
  a.code,
  a.name,
  a.full_name,
  (SELECT name
   FROM author_i18n
   WHERE author_id = a.id AND language = 'src')  AS "src.name",
  (SELECT full_name
   FROM author_i18n
   WHERE author_id = a.id AND language = 'src')  AS "src.full_name",
  (SELECT name
   FROM author_i18n
   WHERE author_id = a.id AND language = 'dest') AS "dest.name",
  (SELECT full_name
   FROM author_i18n
   WHERE author_id = a.id AND language = 'dest') AS "dest.full_name"
FROM authors a
ORDER BY a.id
) TO STDOUT WITH CSV HEADER;
