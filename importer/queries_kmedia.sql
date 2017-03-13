
-- number of containers
SELECT count(*)
FROM containers
WHERE virtual_lesson_id != 0 AND content_type_id = 4;

-- number of files (including SHA1 duplicates)
SELECT count(*)
FROM virtual_lessons v
  INNER JOIN containers c ON v.id = c.virtual_lesson_id and c.content_type_id = 4
  INNER JOIN containers_file_assets cfa ON c.id = cfa.container_id
  INNER JOIN file_assets fa ON fa.id = cfa.file_asset_id;

SELECT count(*)
FROM containers c
  INNER JOIN containers_file_assets cfa ON c.id = cfa.container_id
  INNER JOIN file_assets fa ON fa.id = cfa.file_asset_id
where c.content_type_id = 4 and c.virtual_lesson_id = 0;

-- List of files with same SHA1 in different containers
SELECT
  fa.sha1,
  fa.id,
  fa.size,
  fa.name,
  c.id,
  c.name
FROM containers_file_assets cfa
  INNER JOIN file_assets fa ON fa.id = cfa.file_asset_id
  INNER JOIN containers c ON c.id = cfa.container_id
WHERE fa.sha1 IN (
  SELECT fa.sha1
  FROM file_assets fa
    INNER JOIN containers_file_assets cfa ON fa.id = cfa.file_asset_id
  WHERE fa.sha1 IS NOT NULL
  GROUP BY fa.sha1
  HAVING count(DISTINCT cfa.container_id) > 1
  ORDER BY count(DISTINCT cfa.container_id) DESC
)
ORDER BY fa.sha1, fa.id, c.id;

-- dump previous query to csv
-- grep da39a3ee5e6b4b0d3255bfef95601890afd80709 that file to see empty files (physical size = 0)
COPY (
SELECT
  fa.id,
  fa.sha1,
  fa.size,
  fa.name,
--   c.id,
--   c.name,
  s.servername,
--   s.path,
  concat(s.httpurl, '/', fa.name)
FROM containers_file_assets cfa
  INNER JOIN file_assets fa ON fa.id = cfa.file_asset_id
  INNER JOIN containers c ON c.id = cfa.container_id
  inner join servers s on fa.servername_id = s.servername
WHERE fa.sha1 IN (
  SELECT fa.sha1
  FROM file_assets fa
    INNER JOIN containers_file_assets cfa ON fa.id = cfa.file_asset_id
  WHERE fa.sha1 IS NOT NULL
  GROUP BY fa.sha1
  HAVING count(DISTINCT cfa.container_id) > 1
  ORDER BY count(DISTINCT cfa.container_id) DESC
)
ORDER BY fa.sha1, fa.id, c.id)
TO '/var/lib/postgres/data/kmedia_dup_sha1.csv' (
FORMAT CSV );


-- pg_restore:
-- pg_restore --username=postgres --clean --no-owner --jobs=2 ~/projects/kmedia/kabbalahmedia.20170313.backup

COPY
(SELECT
   fa.sha1,
   fa.name
 FROM containers c
   INNER JOIN containers_file_assets cfa ON c.id = cfa.container_id
   INNER JOIN file_assets fa ON fa.id = cfa.file_asset_id AND fa.sha1 IS NOT NULL
 WHERE c.content_type_id = 4 AND c.virtual_lesson_id != 0
 ORDER BY fa.sha1)
TO '/var/lib/postgres/data/kmedia_files_by_sha1.csv' (
FORMAT CSV );