#!/usr/bin/env sh

set +e
set -x

BASE_DIR="/app"
TIMESTAMP="$(date '+%Y%m%d%H%M%S')"
LOG_FILE="/tmp/import_storage_$TIMESTAMP.log"

cleanup() {
  find "/tmp" -type f -name 'import_storage_*.log' -mtime +7 -exec rm -rf {} \;
}

cd ${BASE_DIR} && ./mdb storage > ${LOG_FILE} 2>&1

WARNINGS="$(grep -Ec "level=(warning|error)" ${LOG_FILE})"

if [ "$WARNINGS" = 0 ];then
        echo "No warnings"
        cleanup
        exit 0
fi

echo "Errors in periodic import of storage catalog to MDB" | mail -s "ERROR: MDB storage import" -r "mdb@bbdomain.org" -a ${LOG_FILE} edoshor@gmail.com

cleanup
exit 1

