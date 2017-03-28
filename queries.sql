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

-- kmedia catalogs with no patterns
WITH RECURSIVE rec_catalogs AS (
  SELECT
    c.id,
    c.name :: TEXT path
  FROM catalogs c
  WHERE c.id = 4016
  --   WHERE c.parent_id IS NULL
  --         AND c.id NOT IN (SELECT DISTINCT catalog_id
  --                          FROM catalogs_container_description_patterns)
  --         AND c.id IN (SELECT DISTINCT catalog_id
  --                      FROM catalogs_containers)
  UNION
  SELECT
    c.id,
    concat(rc.path, '/', c.name)
  FROM catalogs c INNER JOIN rec_catalogs rc ON c.parent_id = rc.id
  --   WHERE c.id NOT IN (SELECT catalog_id
  --                      FROM catalogs_container_description_patterns)
  --         AND c.id IN (SELECT DISTINCT catalog_id
  --                      FROM catalogs_containers)
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
   WHERE catalog_id = c.id AND lang_id = 'SPA') "es.name"
--   (SELECT name
--    FROM catalog_descriptions
--    WHERE catalog_id = c.id AND lang_id = 'UKR') "ua.name"
FROM catalogs c
WHERE c.parent_id = 12
ORDER BY c.id
) TO '/var/lib/postgres/data/kmedia_holidays.csv' (
FORMAT CSV );