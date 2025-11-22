#!/bin/bash
echo "mode: set" > coverage.out
for dir in $(go list ./... | grep -v /vendor/ | grep -v /test/); do
    go test -coverprofile=profile.out -covermode=set $dir
    if [ -f profile.out ]; then
        tail -n +2 profile.out >> coverage.out
        rm profile.out
    fi
done
