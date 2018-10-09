#!/usr/bin/env bash

set -e
workdir=.cover
profile="cover.out"
mode=count

rm -rf "$workdir"
mkdir "$workdir"

for d in $(go list ./... | grep -v vendor); do
    f="$workdir/$(echo $d | tr / -).cover"
    go test -covermode="$mode" -coverprofile="$f" $d
done

echo "mode: $mode" >"$profile"
grep -h -v "^mode:" "$workdir"/*.cover >>"$profile"

rm -rf .cover
