-- rambler up

WITH data(name) AS (VALUES ('PUBLIC_EVENTS'))
INSERT
INTO content_types (name)
SELECT d.name
FROM data AS d
WHERE NOT EXISTS(SELECT ct.name
                 FROM content_types AS ct
                 WHERE ct.name = d.name);

-- rambler down

DELETE
FROM content_types
WHERE name IN ('PUBLIC_EVENTS');
