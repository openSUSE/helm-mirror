#!/usr/bin/env bash

set -e
workdir=.cover
profile="coverage.txt"
mode=atomic

rm -rf "$workdir"
mkdir "$workdir"

for d in $(go list ./... | grep -v -E "vendor|fixtures"); do
    f="$workdir/$(echo $d | tr / -).cover"
    go test -covermode="$mode" -coverprofile="$f" $d
done

echo "mode: $mode" >"$profile"
grep -h -v "^mode:" "$workdir"/*.cover >>"$profile"

rm -rf .cover
