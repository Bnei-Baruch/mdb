echo "Dropping DB"
dropdb mdb

echo "Creating DB"
createdb mdb

echo "Migrating schema"
psql -d mdb -f migrations/initial_schema.sql --quiet