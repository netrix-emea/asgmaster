#!/bin/bash

BUILD_DIR=bin
NAME=asgmaster

mkdir $BUILD_DIR

for GOOS in darwin linux; do
    for GOARCH in amd64; do
        GOOS=$GOOS GOARCH=$GOARCH go build -v -o $BUILD_DIR/$NAME-$GOOS-$GOARCH
    done
done
