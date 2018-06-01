-- rambler up
INSERT INTO users (email) VALUES
  -- Operator for operations.
  ('operator1@dev.com'),
  ('operator2@dev.com'),
  ('operator3@dev.com'),
  ('operator@dev.com');


DROP TABLE IF EXISTS batch_convert;
CREATE TABLE batch_convert (
  file_id       BIGINT REFERENCES files                          NOT NULL,
  operation_id  BIGINT REFERENCES operations                     NULL,
  request_at    TIMESTAMP WITH TIME ZONE                         NULL,
  request_error TEXT                                             NULL
);

INSERT INTO sources (id, uid, pattern, type_id, position, name) VALUES
  (1, 'L2jMWyce', 'test-source-pattern-1', 1, 0, 'test-source-name-1'),
  (2, '5sLqsXjD', 'test-source-pattern-2', 1, 0, 'test-source-name-2'),
  (3, 'oMq3uU8L', 'test-source-pattern-3', 1, 0, 'test-source-name-3'),
  (4, 'DVSS0xAR', 'test-source-pattern-4', 1, 0, 'test-source-name-4'),
  (5, 'AwGBQX2L', 'test-source-pattern-5', 1, 0, 'test-source-name-5'),
  (6, 'cSyh3vQM', 'test-source-pattern-6', 1, 0, 'test-source-name-6'),
  (7, '43BXTx3C', 'test-source-pattern-7', 1, 0, 'test-source-name-7'),
  (8, 'yUcfylRm', 'test-source-pattern-8', 1, 0, 'test-source-name-8'),
  (9, 'nbXRIizB', 'test-source-pattern-9', 1, 0, 'test-source-name-9'),
  (10, 'dvIBxCLT', 'test-source-pattern-10', 1, 0, 'test-source-name-10'),
  (11, 'Uh0TFSOM', 'test-source-pattern-11', 1, 0, 'test-source-name-11'),
  (12, 'B0H0duMW', 'test-source-pattern-12', 1, 0, 'test-source-name-12'),
  (13, 'L8daQ2n7', 'test-source-pattern-13', 1, 0, 'test-source-name-13'),
  (14, '0jt04pix', 'test-source-pattern-14', 1, 0, 'test-source-name-14'),
  (15, '0QygF8Ib', 'test-source-pattern-15', 1, 0, 'test-source-name-15');

INSERT INTO tags (id, uid, pattern) VALUES
  (1, 'L2jMWyce', 'test-tag-pattern-1'),
  (2, '5sLqsXjD', 'test-tag-pattern-2'),
  (3, 'oMq3uU8L', 'test-tag-pattern-3'),
  (4, 'DVSS0xAR', 'test-tag-pattern-4'),
  (5, 'AwGBQX2L', 'test-tag-pattern-5'),
  (6, 'cSyh3vQM', 'test-tag-pattern-6'),
  (7, '43BXTx3C', 'test-tag-pattern-7'),
  (8, 'yUcfylRm', 'test-tag-pattern-8'),
  (9, 'nbXRIizB', 'test-tag-pattern-9'),
  (10, 'dvIBxCLT', 'test-tag-pattern-10'),
  (11, 'Uh0TFSOM', 'test-tag-pattern-11'),
  (12, 'B0H0duMW', 'test-tag-pattern-12'),
  (13, 'L8daQ2n7', 'test-tag-pattern-13'),
  (14, '0jt04pix', 'test-tag-pattern-14'),
  (15, '0QygF8Ib', 'test-tag-pattern-15');


DO $a$
DECLARE ver integer;
BEGIN
  SELECT current_setting('server_version_num') INTO ver;
  IF (ver >= 90700) THEN
    EXECUTE 'CREATE OR REPLACE FUNCTION pg_current_xlog_insert_location() RETURNS pg_lsn AS $$ SELECT pg_current_wal_insert_lsn();$$ LANGUAGE SQL;';
  END IF;
END
$a$;