#!/usr/bin/env bash

set +x
set -e

for ((i = 1995 ; i < 2030 ; i++ )); do
   curl --create-dirs -o "data/hebcal_${i}.json" "http://www.hebcal.com/hebcal/?v=1&cfg=json&maj=on&min=on&mod=on&mf=on&c=on&i=on&lg=h&year=${i}"
done
