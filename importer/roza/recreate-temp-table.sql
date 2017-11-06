DROP TABLE IF EXISTS roza_index_tmp;
CREATE TABLE roza_index_tmp (
  path          VARCHAR(1024) NOT NULL,
  sha1          CHAR(40)      NOT NULL,
  size          BIGINT        NOT NULL,
  last_modified INTEGER       NOT NULL
);
