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

# Replace docs host.
sed -i 's/^HOST: .*$/HOST: app.mdb.bbdomain.org/g' docs/docs.tmpl
sed -i "s/^Release: .*$/Release: ${version}/g" docs/docs.tmpl

echo "Updating docs..."
make docs

echo "Deploying to production"
scp mdb archive@app.mdb.bbdomain.org:/sites/mdb/"mdb-$version"
scp docs/docs.html archive@app.mdb.bbdomain.org:/sites/mdb/docs.html
scp migrations/*.sql archive@app.mdb.bbdomain.org:/sites/mdb/migrations/
ssh archive@app.mdb.bbdomain.org "/sites/mdb/migrations/rambler --configuration=/sites/mdb/migrations/rambler.json apply --all"
ssh archive@app.mdb.bbdomain.org "ln -sf /sites/mdb/mdb-$version /sites/mdb/mdb"
ssh archive@app.mdb.bbdomain.org "supervisorctl restart mdb"

# Revert locally changed docs files.
git reset --hard HEAD
