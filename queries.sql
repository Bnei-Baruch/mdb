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

