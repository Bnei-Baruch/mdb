#!/usr/bin/env bash
# Usage: misc/release.sh
# Build package, tag a commit, push it to origin, and then deploy the
# package on production server.

set -e

# Replace docs host.
sed -i 's/^HOST: .*$/HOST: poc.bbdomain.org:8080/g' docs.tmpl

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
scp mdb archive@poc.bbdomain.org:/sites/mdb/"mdb-$version"
scp docs.html archive@poc.bbdomain.org:/sites/mdb/docs.html
scp migrations/*.sql archive@poc.bbdomain.org:/sites/mdb/migrations/
ssh archive@poc.bbdomain.org "/sites/mdb/migrations/rambler --configuration=/sites/mdb/migrations/rambler.json apply --all"
ssh archive@poc.bbdomain.org "ln -sf /sites/mdb/mdb-$version /sites/mdb/mdb"
ssh archive@poc.bbdomain.org "supervisorctl restart mdb"

