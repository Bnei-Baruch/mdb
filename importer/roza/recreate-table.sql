DROP TABLE IF EXISTS roza_index;
CREATE TABLE roza_index (
  path          VARCHAR(1024)               NOT NULL,
  sha1          BYTEA                       NOT NULL,
  size          BIGINT                      NOT NULL,
  last_modified TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

INSERT INTO roza_index (path, sha1, size, last_modified)
  SELECT
    path,
    decode(sha1, 'hex'),
    size,
    to_timestamp(last_modified)
  FROM roza_index_tmp;

CREATE INDEX IF NOT EXISTS roza_index_path_idx
ON roza_index USING BTREE (path);

CREATE INDEX IF NOT EXISTS roza_index_sha1_idx
  ON roza_index USING BTREE (sha1);

DROP TABLE IF EXISTS roza_index_tmp;