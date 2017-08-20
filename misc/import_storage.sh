#!/usr/bin/env bash

set +e
set -x

BASE_DIR="/sites/mdb"
TIMESTAMP="$(date '+%Y%m%d%H%M%S')"
LOG_FILE="$BASE_DIR/logs/storage_import_$TIMESTAMP.log"

cd ${BASE_DIR} && ./mdb storage > ${LOG_FILE} 2>&1

WARNINGS="$(egrep -c "level=(warning|error)" ${LOG_FILE})"

if [ "$WARNINGS" = "" ];then
	echo "No warnings"
	exit 0
fi

mail -s "ERROR: MDB storage import" -r "mdb@bbdomain.org" -a ${LOG_FILE} edoshor@gmail.com
