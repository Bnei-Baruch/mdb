#!/usr/bin/env bash

set +x
set -e

FILENAME="20171103-archive"
DATADIR="$(dirname $(readlink -f "$0"))/data"

cd ${DATADIR}

if [ ! -f "${DATADIR}/${FILENAME}" ]; then
    xz -dk "${FILENAME}.xz"
    sed 's/.//;s/.$//' ${FILENAME} > "${FILENAME}.csv"
fi

psql -d mdb -f ../recreate-temp-table.sql
psql -d mdb -c "\copy roza_index_tmp (path, sha1,size,last_modified) from '${DATADIR}/${FILENAME}.csv' csv;"
psql -d mdb -f ../recreate-table.sql

