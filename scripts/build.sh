#!/usr/bin/env bash

set -euo pipefail

if [[ -d $PWD/go-cache ]]; then
  export GOPATH=$PWD/go-cache
fi

GOOS="linux" go build -ldflags='-s -w' -o bin/build cmd/build/main.go
GOOS="linux" go build -ldflags='-s -w' -o bin/detect cmd/detect/main.go
