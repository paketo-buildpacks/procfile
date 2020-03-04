#!/usr/bin/env bash

set -euo pipefail

if [[ -d ../go-cache ]]; then
  GOPATH=$(realpath ../go-cache)
  export GOPATH
fi

GOOS="linux" go build -ldflags='-s -w' -o bin/build cmd/build
GOOS="linux" go build -ldflags='-s -w' -o bin/detect cmd/detect
