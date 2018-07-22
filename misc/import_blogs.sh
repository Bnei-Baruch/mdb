#!/usr/bin/env bash

set +e
set -x

BASE_DIR="/sites/mdb"
TIMESTAMP="$(date '+%Y%m%d%H%M%S')"
LOG_FILE="$BASE_DIR/logs/blogs/import_$TIMESTAMP.log"

cd ${BASE_DIR} && ./mdb blog-latest > ${LOG_FILE} 2>&1

WARNINGS="$(egrep -c "level=(warning|error)" ${LOG_FILE})"

if [ "$WARNINGS" = 0 ];then
        echo "No warnings"
        exit 0
fi

echo "Errors in periodic import of blogs to MDB" | mail -s "ERROR: MDB blogs import" -r "mdb@bbdomain.org" -a ${LOG_FILE} edoshor@gmail.com

find "${BASE_DIR}/logs/blogs" -type f -mtime +7 -exec rm -rf {} \;
