#!/bin/sh

docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp -e GOOS=darwin -e GOARCH=amd64 golang:1.8 go build -v
