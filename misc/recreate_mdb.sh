#!/usr/bin/env bash

echo "Dropping DB"
dropdb mdb

echo "Creating DB"
createdb mdb

echo "Migrating schema"
cd migrations && rambler apply -a

echo "Seeding database"
psql -d mdb -f data/seed.sql