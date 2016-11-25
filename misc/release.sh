#!/usr/bin/env bash
# Usage: misc/release.sh
# Build package, tag a commit, push it to origin, and then deploy the
# package on production server.

set -e

echo "Building..."
make build

version="$(./mdb version | awk '{print $NF}')"
[ -n "$version" ] || exit 1
echo $version

git commit --allow-empty -a -m "Release $version"
git tag "v$version"
git push origin master
git push origin "v$version"

echo "Deploying to production"
scp mdb root@poc.bbdomain.org:/sites/mdb/"mdb-$version"
ssh root@poc.bbdomain.org "ln -sf /sites/mdb/mdb-$version /sites/mdb/mdb"
ssh root@poc.bbdomain.org "supervisorctl restart mdb"


