-- rambler up

ALTER TABLE collections RENAME COLUMN name TO name_id;
ALTER TABLE collections RENAME COLUMN description TO description_id;

-- rambler down

ALTER TABLE collections RENAME COLUMN name_id TO name;
ALTER TABLE collections RENAME COLUMN description_id TO description;
