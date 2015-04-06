#!/usr/bin/env bash

set -e

function lint() {
    go get github.com/golang/lint/golint
    golint ./...
}

function vet() {
    local -r flags="$1"
    go get golang.org/x/tools/cmd/vet
    go tool vet $flags $PWD
}


case "$TRAVIS_GO_VERSION" in
    1.1)
        # vet isn't compatible with 1.1
        lint
        ;;
    1.2)
        lint
        # vet doesn't support -copylocks flag in 1.2
        vet
        ;;
    *)
        lint
        vet -copylocks=false
        ;;
esac
