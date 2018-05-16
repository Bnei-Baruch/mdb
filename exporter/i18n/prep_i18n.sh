#!/usr/bin/env bash

set -e
set +x

SRC=$1
DEST=$2

echo "prepare i18n pack from $SRC to $DEST"
rm -rf archive_i18n/*

# make sql queries with src and dest languages
sed "s/src/${SRC}/g; s/dest/${DEST}/g;" sources.tpl.sql > sources.sql
sed "s/src/${SRC}/g; s/dest/${DEST}/g;" tags.tpl.sql > tags.sql
sed "s/src/${SRC}/g; s/dest/${DEST}/g;" authors.tpl.sql > authors.sql
sed "s/src/${SRC}/g; s/dest/${DEST}/g;" persons.tpl.sql > persons.sql
sed "s/src/${SRC}/g; s/dest/${DEST}/g;" publishers.tpl.sql > publishers.sql

# export data from mdb to csv files
psql -q -h pgsql.mdb.bbdomain.org -U readonly -d mdb -f sources.sql > archive_i18n/sources.csv
psql -q -h pgsql.mdb.bbdomain.org -U readonly -d mdb -f tags.sql > archive_i18n/tags.csv
psql -q -h pgsql.mdb.bbdomain.org -U readonly -d mdb -f authors.sql > archive_i18n/authors.csv
psql -q -h pgsql.mdb.bbdomain.org -U readonly -d mdb -f persons.sql > archive_i18n/persons.csv
psql -q -h pgsql.mdb.bbdomain.org -U readonly -d mdb -f publishers.sql > archive_i18n/publishers.csv


# copy site i18n files in
cp ~/projects/kmedia-mdb/public/locales/${SRC}/common.json archive_i18n/ui_${SRC}.json

if [ -f ~/projects/kmedia-mdb/public/locales/${DEST}/common.json ]; then
    cp ~/projects/kmedia-mdb/public/locales/${DEST}/common.json archive_i18n/ui_${DEST}.json
else
    cp ~/projects/kmedia-mdb/public/locales/${SRC}/common.json archive_i18n/ui_${DEST}.json
fi

# copy mdb auto naming conventions i18n data file
cp ../../data/i18n.json archive_i18n/mdb_autonames.json

# zip it all
zip archive_i18n_${SRC}_${DEST}.zip archive_i18n/*
