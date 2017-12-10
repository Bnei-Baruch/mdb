#!/usr/bin/env bash

set +x
set -e

MYDIR=$(dirname $(readlink -f "$0"))
ROOTDIR=$(dirname $(dirname ${MYDIR}))

cd ${ROOTDIR}

./mdb roza-match
./mdb roza-match-mdb
./mdb roza-master
