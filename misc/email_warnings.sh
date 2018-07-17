#!/usr/bin/env bash

set +e
set -x

DATE="$(date -d "yesterday 13:00" '+%Y-%m-%d')"
#DATE="2017-04-30"

LINES="$(egrep "${DATE}(.*) level=(warning|error)" /sites/mdb/logs/mdb.log)"

if [ "$LINES" = 0 ];then
	echo "No warnings"
	exit 0
fi

LC=$(echo "$LINES" | grep -c "^")

echo "$LINES" | mail -s "MDB warnings $DATE [$LC]" -r "mdb@bbdomain.org" edoshor@gmail.com

