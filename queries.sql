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
-- DELETE FROM authors_sources;
-- DELETE FROM author_i18n;
-- DELETE FROM authors;
-- DELETE FROM source_i18n;
-- DELETE FROM sources;


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
  FROM rf
  WHERE id != 353590;


-- find all descendants of a file
WITH RECURSIVE rf AS (
  SELECT f.*
  FROM files f
  WHERE f.id = 380600
  UNION
  SELECT f.*
  FROM files f INNER JOIN rf ON f.parent_id = rf.id
) SELECT
    id,
    parent_id,
    name,
    size,
    sha1,
--     created_at,
    content_unit_id,
    properties ->> 'duration'
  FROM rf;

-- update files set content_unit_id=26031 where parent_id=380610;
-- update files set content_unit_id=26032 where parent_id=380662;
-- update files set content_unit_id=26033 where parent_id=380667;
-- update files set content_unit_id=26034 where parent_id=380739;
-- update files set content_unit_id=26035 where parent_id=380791;
update content_units set published=true where id in (26031, 26032, 26033, 26034, 26035);

-- WITH RECURSIVE rf AS (
--   SELECT f.*
--   FROM files f
--   WHERE f.id = 360841
--   UNION
--   SELECT f.*
--   FROM files f INNER JOIN rf ON f.parent_id = rf.id
-- ) SELECT id, parent_id, content_unit_id, name
--   FROM rf
--   WHERE id != 360841 order by parent_id, id;

-- UPDATE files SET content_unit_id = 25378 WHERE id IN
--    (WITH RECURSIVE rf AS ( SELECT f.* FROM files f WHERE f.id = 360841 UNION SELECT f.* FROM files f INNER JOIN rf ON f.parent_id = rf.id ) SELECT id FROM rf);
-- UPDATE files SET content_unit_id = 25382 WHERE id IN
--                                                (WITH RECURSIVE rf AS ( SELECT f.* FROM files f WHERE f.id = 360896 UNION SELECT f.* FROM files f INNER JOIN rf ON f.parent_id = rf.id ) SELECT id FROM rf);
-- UPDATE files SET content_unit_id = 25383 WHERE id IN
--                                                (WITH RECURSIVE rf AS ( SELECT f.* FROM files f WHERE f.id = 360961 UNION SELECT f.* FROM files f INNER JOIN rf ON f.parent_id = rf.id ) SELECT id FROM rf);
-- UPDATE files SET content_unit_id = 25384 WHERE id IN
--                                                (WITH RECURSIVE rf AS ( SELECT f.* FROM files f WHERE f.id = 360975 UNION SELECT f.* FROM files f INNER JOIN rf ON f.parent_id = rf.id ) SELECT id FROM rf);
-- UPDATE files SET content_unit_id = 25385 WHERE id IN
--                                                (WITH RECURSIVE rf AS ( SELECT f.* FROM files f WHERE f.id = 361036 UNION SELECT f.* FROM files f INNER JOIN rf ON f.parent_id = rf.id ) SELECT id FROM rf);
-- UPDATE files SET content_unit_id = 25386 WHERE id IN
--                                                (WITH RECURSIVE rf AS ( SELECT f.* FROM files f WHERE f.id = 361100 UNION SELECT f.* FROM files f INNER JOIN rf ON f.parent_id = rf.id ) SELECT id FROM rf);
-- UPDATE files SET content_unit_id = 25387 WHERE id IN
--                                                (WITH RECURSIVE rf AS ( SELECT f.* FROM files f WHERE f.id = 361155 UNION SELECT f.* FROM files f INNER JOIN rf ON f.parent_id = rf.id ) SELECT id FROM rf);
-- UPDATE files SET content_unit_id = 25388 WHERE id IN
--                                                (WITH RECURSIVE rf AS ( SELECT f.* FROM files f WHERE f.id = 361265 UNION SELECT f.* FROM files f INNER JOIN rf ON f.parent_id = rf.id ) SELECT id FROM rf);
-- UPDATE files SET content_unit_id = 25389 WHERE id IN
--                                                (WITH RECURSIVE rf AS ( SELECT f.* FROM files f WHERE f.id = 361282 UNION SELECT f.* FROM files f INNER JOIN rf ON f.parent_id = rf.id ) SELECT id FROM rf);
-- UPDATE files SET content_unit_id = 25390 WHERE id IN
--                                                (WITH RECURSIVE rf AS ( SELECT f.* FROM files f WHERE f.id = 361340 UNION SELECT f.* FROM files f INNER JOIN rf ON f.parent_id = rf.id ) SELECT id FROM rf);
-- UPDATE files SET content_unit_id = 25391 WHERE id IN
--                                                (WITH RECURSIVE rf AS ( SELECT f.* FROM files f WHERE f.id = 361405 UNION SELECT f.* FROM files f INNER JOIN rf ON f.parent_id = rf.id ) SELECT id FROM rf);
-- UPDATE files SET content_unit_id = 25392 WHERE id IN
--                                                (WITH RECURSIVE rf AS ( SELECT f.* FROM files f WHERE f.id = 361483 UNION SELECT f.* FROM files f INNER JOIN rf ON f.parent_id = rf.id ) SELECT id FROM rf);
-- UPDATE files SET content_unit_id = 25393 WHERE id IN
--                                                (WITH RECURSIVE rf AS ( SELECT f.* FROM files f WHERE f.id = 361638 UNION SELECT f.* FROM files f INNER JOIN rf ON f.parent_id = rf.id ) SELECT id FROM rf);
--
-- UPDATE files SET content_unit_id = 25396 WHERE id IN
--                                                 (WITH RECURSIVE rf AS ( SELECT f.* FROM files f WHERE f.id = 361693 UNION SELECT f.* FROM files f INNER JOIN rf ON f.parent_id = rf.id ) SELECT id FROM rf);
-- UPDATE files SET content_unit_id = 25397 WHERE id IN
--                                                (WITH RECURSIVE rf AS ( SELECT f.* FROM files f WHERE f.id = 361695 UNION SELECT f.* FROM files f INNER JOIN rf ON f.parent_id = rf.id ) SELECT id FROM rf);
-- UPDATE files SET content_unit_id = 25398 WHERE id IN
--                                                (WITH RECURSIVE rf AS ( SELECT f.* FROM files f WHERE f.id = 361696 UNION SELECT f.* FROM files f INNER JOIN rf ON f.parent_id = rf.id ) SELECT id FROM rf);


SELECT
  cu.id,
  cu.properties ->> 'artifact_type'
FROM files f
  INNER JOIN content_units cu ON f.content_unit_id = cu.id AND cu.properties ? 'artifact_type'
WHERE f.parent_id = 1;

UPDATE content_units
SET properties = properties - 'artifact_type'
WHERE id = 1;

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

-- insert into collections (uid, type_id, properties) values
--   ('TYbA7WoZ', 4, '{"active": true, "pattern": "italy", "country": "Italy", "city": "Rome", "start_date": "2017-07-28", "end_date": "2017-07-30", "full_address": "SHG Hotel Antonella, Via Pontina, 00040 Pomezia RM, Italy" }');
--
-- insert into collection_i18n (collection_id, language, name) VALUES
--   (10796, 'ru', 'Переход');

DROP TABLE IF EXISTS file_mappings;
CREATE TABLE file_mappings (
  sha1     CHAR(40),
  k_id     INT    NULL,
  k_cid    INT    NULL,
  m_id     BIGINT NULL,
  m_cuid   BIGINT NULL,
  m_exists BOOLEAN
);


-- update content_units duration property
UPDATE content_units
SET properties = properties || jsonb_build_object('duration', b.duration)
FROM
  (SELECT
     content_unit_id,
     round(a) AS duration
   FROM (SELECT
           content_unit_id,
           avg((properties ->> 'duration') :: REAL)    AS a,
           stddev((properties ->> 'duration') :: REAL) AS s
         FROM files
         WHERE type IN ('audio', 'video') AND content_unit_id IN (SELECT id
                                                                  FROM content_units
                                                                  WHERE properties ? 'duration' IS FALSE)
         GROUP BY content_unit_id) AS t) AS b
WHERE id = b.content_unit_id;


WITH RECURSIVE rs AS (
  SELECT s.*
  FROM sources s
  WHERE s.id = 1754
  UNION
  SELECT s.*
  FROM sources s INNER JOIN rs ON s.parent_id = rs.id
) SELECT *
  FROM rs;

-- kmedia congresses
WITH RECURSIVE rc AS (
  SELECT c.*
  FROM catalogs c
  WHERE c.id = 40
  UNION
  SELECT c.*
  FROM catalogs c INNER JOIN rc ON c.parent_id = rc.id
) SELECT count(cn.*)
  FROM rc
    INNER JOIN catalogs_containers cc ON rc.id = cc.catalog_id
    INNER JOIN containers cn ON cc.container_id = cn.id;


-- all sources for translation (Dima Perkin)
COPY (
WITH RECURSIVE rec_sources AS (
  SELECT
    s.id,
    s.uid,
    s.position,
    (SELECT name
     FROM source_i18n
     WHERE source_id = s.id AND language = 'he') AS "he.name",
    (SELECT name
     FROM source_i18n
     WHERE source_id = s.id AND language = 'ru') AS "ru.name",
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
     WHERE source_id = s.id AND language = 'he') AS "he.name",
    (SELECT name
     FROM source_i18n
     WHERE source_id = s.id AND language = 'ru') AS "ru.name",
    rs.path || s.id
  FROM sources s INNER JOIN rec_sources rs ON s.parent_id = rs.id
  --   WHERE rs.depth < 2
)
SELECT *
FROM rec_sources
ORDER BY path, position
) TO '/var/lib/postgres/data/all_sources.csv'  (
FORMAT CSV );

-- all tags for translation
COPY (
WITH RECURSIVE rec_tags AS (
  SELECT
    t.id,
    t.uid,
    (SELECT label
     FROM tag_i18n
     WHERE tag_id = t.id AND language = 'he') AS "he.name",
    (SELECT label
     FROM tag_i18n
     WHERE tag_id = t.id AND language = 'ru') AS "ru.name",
    ARRAY [t.id]                                    "path"
  FROM tags t
  WHERE t.parent_id IS NULL
  UNION
  SELECT
    t.id,
    t.uid,
    (SELECT label
     FROM tag_i18n
     WHERE tag_id = t.id AND language = 'he') AS "he.name",
    (SELECT label
     FROM tag_i18n
     WHERE tag_id = t.id AND language = 'ru') AS "ru.name",
    rt.path || t.id
  FROM tags t INNER JOIN rec_tags rt ON t.parent_id = rt.id
)
SELECT *
FROM rec_tags
ORDER BY path
) TO '/var/lib/postgres/data/all_tags.csv'  (
FORMAT CSV );


-- all sources for roza mappings
COPY (
WITH RECURSIVE rec_sources AS (
  SELECT
    s.id,
    s.position,
    concat(ai.name, ', ', (SELECT name
                           FROM source_i18n
                           WHERE source_id = s.id AND language = 'he')) AS name,
    ARRAY [s.id]                                                           "path"
  FROM sources s INNER JOIN authors_sources aas ON s.id = aas.source_id
    INNER JOIN authors a ON a.id = aas.author_id
    INNER JOIN author_i18n ai ON a.id = ai.author_id AND ai.language = 'he'
  WHERE s.parent_id IS NULL
  UNION
  SELECT
    s.id,
    s.position,
    concat(rs.name, ', ', (SELECT name
                           FROM source_i18n
                           WHERE source_id = s.id AND language = 'he')) AS name,
    rs.path || s.id
  FROM sources s INNER JOIN rec_sources rs ON s.parent_id = rs.id
)
SELECT *
FROM rec_sources
ORDER BY path, position
) TO '/var/lib/postgres/data/all_sources_roza.csv'  (
FORMAT CSV );


-- manual mapping of existing congresses in MDB
update collections set properties = properties || '{"kmedia_id": 8024}' where id = 10641;
update collections set properties = properties || '{"kmedia_id": 8029}' where id = 10642;
update collections set properties = properties || '{"kmedia_id": 8027}' where id = 10643;
update collections set properties = properties || '{"kmedia_id": 8100}' where id = 10644;
update collections set properties = properties || '{"kmedia_id": 8084}' where id = 10713;
update collections set properties = properties || '{"kmedia_id": 8127}' where id = 10813;

-- kmedia programs with their containers
copy (
WITH RECURSIVE rc AS (
  SELECT c.*
  FROM catalogs c
  WHERE c.id = 3672
  UNION
  SELECT c.*
  FROM catalogs c INNER JOIN rc ON c.parent_id = rc.id
) SELECT
    rc.id,
    rc.parent_id,
    count(1),
    array_agg(cc.container_id)
  FROM rc
    INNER JOIN catalogs_containers cc ON rc.id = cc.catalog_id
  GROUP BY rc.id, rc.parent_id order by rc.parent_id, rc.id
) TO '/var/lib/postgres/data/programs_chapters2.csv'  (
FORMAT CSV );



-- flatten parashat shavua

WITH RECURSIVE rc AS (
  SELECT c.*
  FROM catalogs c
  WHERE c.parent_id = 3624
  UNION
  SELECT c.*
  FROM catalogs c INNER JOIN rc ON c.parent_id = rc.id
) SELECT
    array_agg(cc.container_id)
  FROM rc inner join catalogs_containers cc on rc.id = cc.catalog_id;


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

update collections set properties = properties || '{"kmedia_id": 7899}' where id = 10651;
update collections set properties = properties || '{"kmedia_id": 8159}' where id = 10796;


-- unit original_language from historical send operations metadata
UPDATE content_units cu
SET properties = properties - 'original_language'
WHERE properties ->> 'original_language' = '';

UPDATE content_units cu
SET properties = properties || jsonb_build_object('original_language', tmp.lang) FROM
  (SELECT DISTINCT ON (f.content_unit_id)
     f.content_unit_id AS cuid,
     CASE WHEN o.properties ->> 'language' = 'heb'
       THEN 'he'
     WHEN o.properties ->> 'language' = 'eng'
       THEN 'en'
     WHEN o.properties ->> 'language' = 'rus'
       THEN 'ru'
     WHEN o.properties ->> 'language' = 'mlt'
       THEN 'zz'
     ELSE 'xx' END     AS lang
   FROM operations o INNER JOIN files_operations fo ON o.id = fo.operation_id AND o.type_id = 4
     INNER JOIN files f ON fo.file_id = f.id AND f.content_unit_id IS NOT NULL
   GROUP BY f.content_unit_id, o.id
   ORDER BY f.content_unit_id, o.id DESC) AS tmp
WHERE cu.id = tmp.cuid;

-- cu missing name from historical send operations
SELECT
  cu.id,
  tmp.name
FROM content_units cu
  LEFT JOIN content_unit_i18n cui ON cu.id = cui.content_unit_id
  INNER JOIN
  (SELECT DISTINCT ON (f.content_unit_id)
     f.content_unit_id             AS cuid,
     o.properties ->> 'final_name' AS name
   FROM operations o INNER JOIN files_operations fo ON o.id = fo.operation_id AND o.type_id = 4
     INNER JOIN files f ON fo.file_id = f.id AND f.content_unit_id IS NOT NULL
   GROUP BY f.content_unit_id, o.id
   ORDER BY f.content_unit_id, o.id DESC) AS tmp ON cu.id = tmp.cuid
WHERE cui.content_unit_id IS NULL;


SELECT "content_units".* FROM "content_units" INNER JOIN content_units_sources cus ON id = cus.content_unit_id WHERE (secure=0 AND published IS TRUE) AND "type_id" IN (11) AND "cus"."source_id" IN (1754,1755,1756,1757,1758,1759,1760,1761,1762,1763,1764,1765,1766,1767,1768,1769,1770,1771,1772,1773,1774,1775,1776,1777,1778,1779,1780,1781,1782,1783,1784,1785,1786,1787,1788,1789,1790,1791,1792,1793,1794,1795,1796,1797,1798,1799,1800,1801,1802,1803,1804,1805,1806,1807,1808,1809,1810,1811,1812,1813) ORDER BY (properties->>'film_date')::date desc, created_at desc LIMIT 10;

SELECT "content_units".* FROM "content_units" INNER JOIN content_units_sources cus ON id = cus.content_unit_id WHERE (secure=0 AND published IS TRUE) AND "type_id" IN (11) AND "cus"."source_id" IN (1754,1755,1756,1757,1758,1759,1760,1761,1762,1763,1764,1765,1766,1767,1768,1769,1770,1771,1772,1773,1774,1775,1776,1777,1778,1779,1780,1781,1782,1783,1784,1785,1786,1787,1788,1789,1790,1791,1792,1793,1794,1795,1796,1797,1798,1799,1800,1801,1802,1803,1804,1805,1806,1807,1808,1809,1810,1811,1812,1813) group by content_units.id ORDER BY (properties->>'film_date')::date desc, created_at desc LIMIT 10;

SELECT DISTINCT ON (id) *
FROM
  (
    SELECT "content_units".*
    FROM "content_units"
      INNER JOIN content_units_sources cus ON id = cus.content_unit_id
    WHERE (secure = 0 AND published IS TRUE) AND "type_id" IN (11) AND "cus"."source_id" IN
                                                                       (1754, 1755, 1756, 1757, 1758, 1759, 1760, 1761, 1762, 1763, 1764, 1765, 1766, 1767, 1768, 1769, 1770, 1771, 1772, 1773, 1774, 1775, 1776, 1777, 1778, 1779, 1780, 1781, 1782, 1783, 1784, 1785, 1786, 1787, 1788, 1789, 1790, 1791, 1792, 1793, 1794, 1795, 1796, 1797, 1798, 1799, 1800, 1801, 1802, 1803, 1804, 1805, 1806, 1807, 1808, 1809, 1810, 1811, 1812, 1813)
    ORDER BY (properties ->> 'film_date') :: DATE DESC, created_at DESC
  ) AS t
LIMIT 10;


-- fix missing languages for files based on file name
update files set language='he' where content_unit_id=35509 and name like 'heb_%';
update files set language='en' where content_unit_id=35509 and name like 'eng_%';
update files set language='ru' where content_unit_id=35509 and name like 'rus_%';
update files set language='es' where content_unit_id=35509 and name like 'spa_%';
update files set language='it' where content_unit_id=35509 and name like 'ita_%';
update files set language='de' where content_unit_id=35509 and name like 'ger_%';
update files set language='nl' where content_unit_id=35509 and name like 'dut_%';
update files set language='fr' where content_unit_id=35509 and name like 'fre_%';
update files set language='pt' where content_unit_id=35509 and name like 'por_%';
update files set language='tr' where content_unit_id=35509 and name like 'trk_%';
update files set language='pl' where content_unit_id=35509 and name like 'pol_%';
update files set language='ar' where content_unit_id=35509 and name like 'arb_%';
update files set language='hu' where content_unit_id=35509 and name like 'hun_%';
update files set language='fi' where content_unit_id=35509 and name like 'fin_%';
update files set language='lt' where content_unit_id=35509 and name like 'lit_%';
update files set language='ja' where content_unit_id=35509 and name like 'jpn_%';
update files set language='bg' where content_unit_id=35509 and name like 'bul_%';
update files set language='ka' where content_unit_id=35509 and name like 'geo_%';
update files set language='no' where content_unit_id=35509 and name like 'nor_%';
update files set language='sv' where content_unit_id=35509 and name like 'swe_%';
update files set language='hr' where content_unit_id=35509 and name like 'hrv_%';
update files set language='zh' where content_unit_id=35509 and name like 'chn_%';
update files set language='fa' where content_unit_id=35509 and name like 'far_%';
update files set language='ro' where content_unit_id=35509 and name like 'ron_%';
update files set language='hi' where content_unit_id=35509 and name like 'hin_%';
update files set language='ua' where content_unit_id=35509 and name like 'ukr_%';
update files set language='mk' where content_unit_id=35509 and name like 'mkd_%';
update files set language='sl' where content_unit_id=35509 and name like 'slv_%';
update files set language='lv' where content_unit_id=35509 and name like 'lav_%';
update files set language='sk' where content_unit_id=35509 and name like 'slk_%';
update files set language='cs' where content_unit_id=35509 and name like 'cze_%';
update files set type='video', mime_type='video/mp4' where content_unit_id=35509 and name like '%mp4';
update files set type='audio', mime_type='audio/mpeg' where content_unit_id=35509 and name like '%mp3';
update files set type='text', mime_type='application/msword' where content_unit_id=35509 and name ~ '\.docx?$';


WITH RECURSIVE rec_sources AS (
  SELECT
    s.id,
    s.uid,
    s.position,
    ARRAY [a.code, s.uid] "path"
  FROM sources s INNER JOIN authors_sources aas ON s.id = aas.source_id
    INNER JOIN authors a ON a.id = aas.author_id
  UNION
  SELECT
    s.id,
    s.uid,
    s.position,
    rs.path || s.uid
  FROM sources s INNER JOIN rec_sources rs ON s.parent_id = rs.id
)
SELECT
  cus.content_unit_id,
  array_agg(DISTINCT item)
FROM content_units_sources cus INNER JOIN rec_sources AS rs ON cus.source_id = rs.id
  , unnest(rs.path) item
GROUP BY cus.content_unit_id;


WITH RECURSIVE rec_tags AS (
  SELECT
    t.id,
    t.uid,
    ARRAY [t.uid] :: CHAR(8) [] "path"
  FROM tags t
  WHERE parent_id IS NULL
  UNION
  SELECT
    t.id,
    t.uid,
    (rt.path || t.uid) :: CHAR(8) []
  FROM tags t INNER JOIN rec_tags rt ON t.parent_id = rt.id
)
SELECT
  cut.content_unit_id,
  array_agg(DISTINCT item)
FROM content_units_tags cut INNER JOIN rec_tags AS rt ON cut.tag_id = rt.id
  , unnest(rt.path) item
GROUP BY cut.content_unit_id;


SELECT
  cup.content_unit_id,
  array_agg(p.uid)
FROM content_units_persons cup INNER JOIN persons p ON cup.person_id = p.id
GROUP BY cup.content_unit_id;

SELECT
  content_unit_id,
  array_agg(DISTINCT language)
FROM files
WHERE language NOT IN ('zz', 'xx') AND content_unit_id IS NOT NULL
GROUP BY content_unit_id;


SELECT array_agg(DISTINCT id)
FROM collections
WHERE type_id = 5 AND properties -> 'genres' ?| ARRAY ['educational'];

SELECT
  c.id,
  max(cu.properties ->> 'film_date') max_film_date,
  count(cu.id)
FROM collections c INNER JOIN collections_content_units ccu
    ON c.id = ccu.collection_id AND c.type_id = 5 AND c.secure = 0 AND c.published IS TRUE
  INNER JOIN content_units cu
    ON ccu.content_unit_id = cu.id AND cu.secure = 0 AND cu.published IS TRUE AND cu.properties ? 'film_date'
GROUP BY c.id
ORDER BY max_film_date DESC;


-- fix missing properties for doc files
update files set language='he', type='text' where name ~ '^heb.*\.docx?$';
update files set language='en', type='text' where name ~ '^eng.*\.docx?$';
update files set language='ru', type='text' where name ~ '^rus.*\.docx?$';
update files set language='es', type='text' where name ~ '^spa.*\.docx?$';
update files set language='it', type='text' where name ~ '^ita.*\.docx?$';
update files set language='de', type='text' where name ~ '^ger.*\.docx?$';
update files set language='nl', type='text' where name ~ '^dut.*\.docx?$';
update files set language='fr', type='text' where name ~ '^fre.*\.docx?$';
update files set language='pt', type='text' where name ~ '^por.*\.docx?$';
update files set language='tr', type='text' where name ~ '^trk.*\.docx?$';
update files set language='pl', type='text' where name ~ '^pol.*\.docx?$';
update files set language='ar', type='text' where name ~ '^arb.*\.docx?$';
update files set language='hu', type='text' where name ~ '^hun.*\.docx?$';
update files set language='fi', type='text' where name ~ '^fin.*\.docx?$';
update files set language='lt', type='text' where name ~ '^lit.*\.docx?$';
update files set language='ja', type='text' where name ~ '^jpn.*\.docx?$';
update files set language='bg', type='text' where name ~ '^bul.*\.docx?$';
update files set language='ka', type='text' where name ~ '^geo.*\.docx?$';
update files set language='no', type='text' where name ~ '^nor.*\.docx?$';
update files set language='sv', type='text' where name ~ '^swe.*\.docx?$';
update files set language='hr', type='text' where name ~ '^hrv.*\.docx?$';
update files set language='zh', type='text' where name ~ '^chn.*\.docx?$';
update files set language='fa', type='text' where name ~ '^far.*\.docx?$';
update files set language='ro', type='text' where name ~ '^ron.*\.docx?$';
update files set language='hi', type='text' where name ~ '^hin.*\.docx?$';
update files set language='ua', type='text' where name ~ '^ukr.*\.docx?$';
update files set language='mk', type='text' where name ~ '^mkd.*\.docx?$';
update files set language='sl', type='text' where name ~ '^slv.*\.docx?$';
update files set language='lv', type='text' where name ~ '^lav.*\.docx?$';
update files set language='sk', type='text' where name ~ '^slk.*\.docx?$';
update files set language='cs', type='text' where name ~ '^cze.*\.docx?$';

-- kmedia congresses sorted by catalog hierarchy
WITH RECURSIVE rc AS (
  SELECT
    c.*,
    c.id :: TEXT AS path
  FROM catalogs c
  WHERE c.id = 40
  UNION
  SELECT
    c.*,
    concat(rc.path, '/', c.id) AS path
  FROM catalogs c INNER JOIN rc ON c.parent_id = rc.id
) SELECT
    rc.id,
    rc.parent_id,
    rc.path,
    rc.name
  FROM rc
  ORDER BY rc.path;

8128 |      4564 | 40/4564/8128      | congress_virtual_unityday_2010-08
8129 |      4564 | 40/4564/8129      | congress_virtual_unityday_2010-09
8130 |      4564 | 40/4564/8130      | congress_virtual_unityday_2010-10
8131 |      4564 | 40/4564/8131      | congress_virtual_unityday_2010-11
8132 |      4564 | 40/4564/8132      | congress_virtual_unityday_2010-12
8133 |      4564 | 40/4564/8133      | congress_virtual_unityday_2011-02
8134 |      4564 | 40/4564/8134      | congress_virtual_unityday_2011-03
8135 |      4564 | 40/4564/8135      | congress_virtual_unityday_2011-06
8136 |      4564 | 40/4564/8136      | congress_virtual_unityday_2011-08
8137 |      4564 | 40/4564/8137      | congress_virtual_unityday_2012-04
8166 |      4564 | 40/4564/8166      | congress_virtual_unityday_2011-07


SELECT
  ccu.collection_id,
  array_agg(DISTINCT cu.type_id)
FROM collections_content_units ccu
  INNER JOIN content_units cu ON ccu.content_unit_id = cu.id
  INNER JOIN collections c ON ccu.collection_id = c.id and c.type_id in (1,2)
GROUP BY ccu.collection_id
HAVING not (11 = any(array_agg(DISTINCT cu.type_id))) and not (21 = any(array_agg(DISTINCT cu.type_id)));


-- kmedia catalogs with container count
COPY
(WITH RECURSIVE rec_catalogs AS (
  SELECT
    c.id,
    c.name :: TEXT               path,
    (SELECT count(DISTINCT container_id)
     FROM catalogs_containers
     WHERE catalog_id = c.id) AS containers,
    1                            depth
  FROM catalogs c
  WHERE parent_id IS NULL
  UNION
  SELECT
    c.id,
    concat(rc.path, '/', c.name),
    (SELECT count(DISTINCT container_id)
     FROM catalogs_containers
     WHERE catalog_id = c.id) AS containers,
    rc.depth + 1
  FROM catalogs c INNER JOIN rec_catalogs rc ON c.parent_id = rc.id
)
SELECT
  id,
  path,
  containers
FROM rec_catalogs
ORDER BY path)
TO '/var/lib/postgres/data/catalog_tree.csv' (
FORMAT CSV );

-- containers under any lesson's catalogs
WITH RECURSIVE rec_catalogs AS (
  SELECT c.id
  FROM catalogs c
  WHERE id IN (11, 6932, 6933, 4772, 4541, 3629, 4020, 4700, 3630, 3631, 4841, 4862, 4728, 4016, 4761, 3632)
  UNION
  SELECT c.id
  FROM catalogs c INNER JOIN rec_catalogs rc ON c.parent_id = rc.id
)
SELECT
  DISTINCT cc.container_id
FROM rec_catalogs rc INNER JOIN catalogs_containers cc ON rc.id = cc.catalog_id;


WITH RECURSIVE rec_catalogs AS (
  SELECT c.id
  FROM catalogs c
  WHERE id IN (8154)
  UNION
  SELECT c.id
  FROM catalogs c INNER JOIN rec_catalogs rc ON c.parent_id = rc.id
)
SELECT
  count(DISTINCT cc.container_id)
FROM rec_catalogs rc INNER JOIN catalogs_containers cc ON rc.id = cc.catalog_id;

SELECT
  fa.id,
  fa.sha1,
  array_agg(DISTINCT cfa.container_id)
FROM file_assets fa INNER JOIN containers_file_assets cfa ON fa.id = cfa.file_asset_id AND fa.sha1 IS NOT NULL
GROUP BY fa.id;

-- publishers
WITH RECURSIVE rec_catalogs AS (
  SELECT c.id
  FROM catalogs c
  WHERE id = 7957  -- itonut
  UNION
  SELECT c.id
  FROM catalogs c INNER JOIN rec_catalogs rc ON c.parent_id = rc.id
)
SELECT
  DISTINCT regexp_replace(split_part(fa.name, '.', 1), '^.*_', '') AS publicator
FROM rec_catalogs rc INNER JOIN catalogs_containers cc ON rc.id = cc.catalog_id
  INNER JOIN containers_file_assets cfa ON cc.container_id = cfa.container_id
  INNER JOIN file_assets fa ON cfa.file_asset_id = fa.id AND fa.asset_type = 'zip'
ORDER BY publicator;

SELECT *
FROM "files"
WHERE (uid = '72824d8e6dd0103d90957854c786d6c47a0056ce') OR (id :: TEXT = '72824d8e6dd0103d90957854c786d6c47a0056ce') OR
      (sha1 :: TEXT ~ '72824d8e6dd0103d90957854c786d6c47a0056ce') OR (name ~ '72824d8e6dd0103d90957854c786d6c47a0056ce')
ORDER BY id DESC
LIMIT 50;


-- cu with image files in more than one language
SELECT
  cu.id,
  cu.uid,
  cu.type_id,
  array_agg(DISTINCT f.language)
FROM content_units cu INNER JOIN files f ON cu.id = f.content_unit_id AND f.type = 'image'
GROUP BY cu.id
HAVING count(DISTINCT f.language) > 1
ORDER BY cu.created_at DESC;

-- Published CU's without published files
SELECT
  cu.id,
  cu.uid,
  ct.name,
  cu.created_at,
  cu.properties
FROM content_units cu
  INNER JOIN content_types ct ON cu.type_id = ct.id
  LEFT JOIN files f ON cu.id = f.content_unit_id AND f.published IS TRUE
WHERE cu.secure = 0 AND cu.published IS TRUE AND f.id IS NULL
ORDER BY cu.created_at;

-- kmedia containers mapped to more than one CU
SELECT
  properties ->> 'kmedia_id',
  array_agg(DISTINCT id)
FROM content_units
WHERE properties ? 'kmedia_id'
GROUP BY properties ->> 'kmedia_id'
HAVING count(id) > 1
ORDER BY count(id);

