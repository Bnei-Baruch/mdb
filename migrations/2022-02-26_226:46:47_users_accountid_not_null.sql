-- MDB generated migration file
-- rambler up
ALTER TABLE users
    ALTER COLUMN accounts_id SET NOT NULL;

CREATE UNIQUE INDEX users_accaunt_id_idx ON users (accounts_id)
    WHERE accounts_id IS NOT NULL;


-- rambler down

ALTER TABLE users
    ALTER COLUMN accounts_id DROP NOT NULL;

DROP INDEX IF EXISTS users_accaunt_id_idx;
