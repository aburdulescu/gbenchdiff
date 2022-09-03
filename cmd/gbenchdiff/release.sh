#!/bin/bash

set -e

out_dir=release-out
exe=gbenchdiff

function release() {
    set -x
    dst=$out_dir/$exe"_"$1"_"$2
    GOOS=$1 GOARC=$2 go build -ldflags "-s -w" -o $dst
    gzip $dst
    sha256sum $dst.gz > $dst.gz.sha256
}

rm -rf $out_dir
mkdir -p $out_dir

release linux amd64
release linux 386
release linux arm
release linux arm64

release darwin amd64
release darwin 386
release darwin arm
release darwin arm64

release windows amd64
release windows 386
