-- MDB generated migration file
-- rambler up
ALTER TABLE users
    ADD COLUMN accounts_id VARCHAR(36) UNIQUE,
    ADD COLUMN disabled    BOOLEAN DEFAULT FALSE NOT NULL;

-- rambler down
ALTER TABLE users
    DROP COLUMN accounts_id,
    DROP COLUMN disabled;
