\timing off
COPY (
WITH RECURSIVE rec_tags AS (
  SELECT
    t.id,
    t.uid,
    (SELECT label
     FROM tag_i18n
     WHERE tag_id = t.id AND language = 'src') AS "src.name",
    (SELECT label
     FROM tag_i18n
     WHERE tag_id = t.id AND language = 'dest') AS "dest.name",
    ARRAY [t.id]                                    "path"
  FROM tags t
  WHERE t.parent_id IS NULL
  UNION
  SELECT
    t.id,
    t.uid,
    (SELECT label
     FROM tag_i18n
     WHERE tag_id = t.id AND language = 'src') AS "src.name",
    (SELECT label
     FROM tag_i18n
     WHERE tag_id = t.id AND language = 'dest') AS "dest.name",
    rt.path || t.id
  FROM tags t INNER JOIN rec_tags rt ON t.parent_id = rt.id
)
SELECT id, uid, "src.name", "dest.name"
FROM rec_tags
ORDER BY path
) TO STDOUT WITH CSV HEADER;
