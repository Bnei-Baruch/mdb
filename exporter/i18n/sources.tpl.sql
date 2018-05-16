\timing off
COPY (
WITH RECURSIVE rec_sources AS (
  SELECT
    s.id,
    s.uid,
    s.position,
    (SELECT name
     FROM source_i18n
     WHERE source_id = s.id AND language = 'src') AS "src.name",
    (SELECT name
     FROM source_i18n
     WHERE source_id = s.id AND language = 'dest') AS "dest.name",
    ARRAY [s.id]                                    "path"
  FROM sources s
  WHERE s.parent_id IS NULL
  UNION
  SELECT
    s.id,
    s.uid,
    s.position,
    (SELECT name
     FROM source_i18n
     WHERE source_id = s.id AND language = 'src') AS "src.name",
    (SELECT name
     FROM source_i18n
     WHERE source_id = s.id AND language = 'dest') AS "dest.name",
    rs.path || s.id
  FROM sources s INNER JOIN rec_sources rs ON s.parent_id = rs.id
)
SELECT id, uid, "src.name", "dest.name"
FROM rec_sources
ORDER BY path, position
) TO STDOUT WITH CSV HEADER;
