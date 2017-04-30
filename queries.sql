-- Nice view of collection -> unit -> file
SELECT
  c.id                         AS "c.id",
  c.properties -> 'film_date'  AS "film_date",
  c.properties -> 'kmedia_id'  AS "kmedia",
  c.type_id,
  ccu.name                     AS "part",
  cu.id                        AS "cu.id",
  cu.type_id,
  cui.name,
  cu.properties -> 'duration'  AS "duration",
  cu.properties -> 'kmedia_id' AS "kmedia",
  f.id                         AS "f.id",
  f.name
FROM collections c INNER JOIN collections_content_units ccu ON c.id = ccu.collection_id
  INNER JOIN content_units cu ON ccu.content_unit_id = cu.id
  INNER JOIN content_unit_i18n cui ON cu.id = cui.content_unit_id AND cui.language = 'en'
  LEFT JOIN files f ON cu.id = f.content_unit_id AND f.language = 'en'
ORDER BY c.properties -> 'film_date' DESC, ccu.name :: INT;

-- Sources: top level collections
SELECT
  a.id,
  a.code,
  s.id,
  s.name
FROM authors a
  INNER JOIN authors_sources AS "as" ON a.id = "as".author_id
  INNER JOIN sources s ON "as".source_id = s.id AND s.parent_id IS NULL
ORDER BY a.id;

-- Sources with their i18n
SELECT
  s.id,
  s.uid,
  s.pattern,
  s.name,
  si.language,
  si.name
FROM sources s INNER JOIN source_i18n si ON s.id = si.source_id
ORDER BY s.id;

-- Delete all sources
DELETE FROM authors_sources;
DELETE FROM author_i18n;
DELETE FROM authors;
DELETE FROM source_i18n;
DELETE FROM sources;


WITH RECURSIVE rec_sources AS (
  SELECT s.*
  FROM sources s
  WHERE s.id = 3715
  UNION
  SELECT s.*
  FROM sources s INNER JOIN rec_sources rs ON s.parent_id = rs.id
)
SELECT string_agg(name, '/')
FROM rec_sources
LIMIT 3;


WITH RECURSIVE rec_sources AS (
  SELECT
    s.id,
    concat(a.code, '/', s.name) path
  FROM sources s INNER JOIN authors_sources x ON s.id = x.source_id
    INNER JOIN authors a ON x.author_id = a.id
  WHERE s.parent_id IS NULL
  UNION
  SELECT
    s.id,
    concat(rs.path, '/', s.name)
  FROM sources s INNER JOIN rec_sources rs ON s.parent_id = rs.id
)
SELECT *
FROM rec_sources;

-- sources with named path
COPY (
WITH RECURSIVE rec_sources AS (
  SELECT
    s.id,
    s.pattern,
    concat(a.code, '/', s.name) path
  FROM sources s INNER JOIN authors_sources x ON s.id = x.source_id
    INNER JOIN authors a ON x.author_id = a.id
  WHERE s.parent_id IS NULL
  UNION
  SELECT
    s.id,
    s.pattern,
    concat(rs.path, '/', s.name)
  FROM sources s INNER JOIN rec_sources rs ON s.parent_id = rs.id
)
SELECT *
FROM rec_sources
WHERE pattern IS NOT NULL
ORDER BY pattern
) TO '/var/lib/postgres/data/mdb_sources_patterns.csv' (
FORMAT CSV );

-- kmedia patterns -> catalogs
COPY (
SELECT
  p.id,
  p.pattern,
  p.lang,
  c.id,
  c.name
FROM catalogs c INNER JOIN catalogs_container_description_patterns cp ON c.id = cp.catalog_id
  INNER JOIN container_description_patterns p ON cp.container_description_pattern_id = p.id
ORDER BY p.pattern
) TO '/var/lib/postgres/data/kmedia_patterns_catalogs.csv' (
FORMAT CSV );

-- kmedia catalogs
WITH RECURSIVE rec_catalogs AS (
  SELECT
    c.id,
    c.name :: TEXT path,
    (SELECT DISTINCT cdp.pattern
     FROM catalogs_container_description_patterns ccdp INNER JOIN container_creation_patterns cdp
         ON ccdp.container_description_pattern_id = cdp.id
     WHERE ccdp.catalog_id = c.id
     LIMIT 1) "pattern",
--     (select count(distinct container_id) from catalogs_containers where catalog_id = c.id) as containers
  1 depth
  FROM catalogs c
  WHERE c.id = 3672
  UNION
  SELECT
    c.id,
    concat(rc.path, '/', c.name),
    (SELECT DISTINCT cdp.pattern
     FROM catalogs_container_description_patterns ccdp INNER JOIN container_creation_patterns cdp
         ON ccdp.container_description_pattern_id = cdp.id
     WHERE ccdp.catalog_id = c.id
     LIMIT 1) "pattern",
--     (select count(distinct container_id) from catalogs_containers where catalog_id = c.id) as containers
    rc.depth + 1
  FROM catalogs c INNER JOIN rec_catalogs rc ON c.parent_id = rc.id
  WHERE rc.depth < 2
)
SELECT *
FROM rec_catalogs
ORDER BY path;

-- kmedia catalogs with i18n
COPY (
SELECT
  c.id,
  c.name,
  (SELECT name
   FROM catalog_descriptions
   WHERE catalog_id = c.id AND lang_id = 'ENG') "en.name",
  (SELECT name
   FROM catalog_descriptions
   WHERE catalog_id = c.id AND lang_id = 'HEB') "he.name",
  (SELECT name
   FROM catalog_descriptions
   WHERE catalog_id = c.id AND lang_id = 'RUS') "ru.name",
  (SELECT name
   FROM catalog_descriptions
   WHERE catalog_id = c.id AND lang_id = 'SPA') "es.name",
  (SELECT name
   FROM catalog_descriptions
   WHERE catalog_id = c.id AND lang_id = 'GER') "de.name",
  (SELECT name
   FROM catalog_descriptions
   WHERE catalog_id = c.id AND lang_id = 'UKR') "ua.name",
  (SELECT name
   FROM catalog_descriptions
   WHERE catalog_id = c.id AND lang_id = 'CHN') "zh.name"
FROM catalogs c
WHERE c.parent_id = 3672
ORDER BY c.name
) TO '/var/lib/postgres/data/kmedia_tvshows.csv' (
FORMAT CSV );




WITH RECURSIVE rec_sources AS (
  SELECT
    s.id,
    s.uid,
    s.pattern,
    s.parent_id,
    s.position,
    s.type_id,
    coalesce((SELECT name
              FROM source_i18n
              WHERE source_id = s.id AND language = 'en'),
             (SELECT name
              FROM source_i18n
              WHERE source_id = s.id AND language = 'he')) "name",
    coalesce((SELECT description
              FROM source_i18n
              WHERE source_id = s.id AND language = 'en'),
             (SELECT description
              FROM source_i18n
              WHERE source_id = s.id AND language = 'he')) "description",
    1                                                      "depth"
  FROM sources s
  WHERE s.parent_id IS NULL
  UNION
  SELECT
    s.id,
    s.uid,
    s.pattern,
    s.parent_id,
    s.position,
    s.type_id,
    coalesce((SELECT name
              FROM source_i18n
              WHERE source_id = s.id AND language = 'en'),
             (SELECT name
              FROM source_i18n
              WHERE source_id = s.id AND language = 'he')) "name",
    coalesce((SELECT description
              FROM source_i18n
              WHERE source_id = s.id AND language = 'en'),
             (SELECT description
              FROM source_i18n
              WHERE source_id = s.id AND language = 'he')) "description",
    depth + 1
  FROM sources s INNER JOIN rec_sources rs ON s.parent_id = rs.id
  WHERE rs.depth < 2
)
SELECT *
FROM rec_sources
ORDER BY depth, parent_id, position;


WITH RECURSIVE rec_tags AS (
  SELECT
    t.id,
    t.uid,
    t.pattern,
    t.parent_id,
    coalesce((SELECT label
              FROM tag_i18n
              WHERE tag_id = t.id AND language = 'en'),
             (SELECT label
              FROM tag_i18n
              WHERE tag_id = t.id AND language = 'he')) "label",
    1                                                   "depth"
  FROM tags t
  WHERE t.parent_id IS NULL
  UNION
  SELECT
    t.id,
    t.uid,
    t.pattern,
    t.parent_id,
    coalesce((SELECT label
              FROM tag_i18n
              WHERE tag_id = t.id AND language = 'en'),
             (SELECT label
              FROM tag_i18n
              WHERE tag_id = t.id AND language = 'he')) "label",
    depth + 1
  FROM tags t INNER JOIN rec_tags rt ON t.parent_id = rt.id
  WHERE rt.depth < 2
)
SELECT *
FROM rec_tags
ORDER BY depth, parent_id, label;

WITH RECURSIVE ffo AS (
    SELECT
      file_id,
      min(operation_id) "o_id"
    FROM files_operations
    GROUP BY file_id
),
    rec_files AS (
    SELECT
      f.created_at,
      f.id,
      f.parent_id,
      f.name,
      f.size,
      --     f.properties "fprops",
      o.id    "o_id",
      ot.name "op"
    --     o.properties "oprops"
    FROM files f INNER JOIN ffo ON f.id = ffo.file_id
      INNER JOIN operations o ON ffo.o_id = o.id
      INNER JOIN operation_types ot ON o.type_id = ot.id
    WHERE f.id = 353547
    UNION
    SELECT
      f.created_at,
      f.id,
      f.parent_id,
      f.name,
      f.size,
      --     f.properties "fprops",
      o.id    "o_id",
      ot.name "op"
    --     o.properties "oprops"
    FROM files f INNER JOIN rec_files rf ON f.parent_id = rf.id
      INNER JOIN ffo ON f.id = ffo.file_id
      INNER JOIN operations o ON ffo.o_id = o.id
      INNER JOIN operation_types ot ON o.type_id = ot.id
  ) SELECT *
    FROM rec_files;


-- find all ancestors of a file
WITH RECURSIVE rf AS (
  SELECT f.*
  FROM files f
  WHERE f.id = 353590
  UNION
  SELECT f.*
  FROM files f INNER JOIN rf ON f.id = rf.parent_id
) SELECT *
  FROM rf;




copy (
WITH RECURSIVE rec_sources AS (
  SELECT
    s.id,
    s.pattern,
    si.name::text path
  FROM sources s
    INNER JOIN source_i18n si on s.id = si.source_id and si.language = 'en'
  WHERE s.parent_id IS NULL
  UNION
  SELECT
    s.id,
    s.pattern,
    concat(rs.path, ' ', si.name)
  FROM sources s
    INNER JOIN source_i18n si on s.id = si.source_id and si.language = 'en'
    INNER JOIN rec_sources rs ON s.parent_id = rs.id

)
SELECT *
FROM rec_sources
WHERE pattern IS NOT NULL
ORDER BY pattern)
to '/var/lib/postgres/data/titles.csv' (format CSV);
