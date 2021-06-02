#!/usr/bin/env bash

set -euo pipefail

GOOS="linux" go build -ldflags='-s -w' -o bin/main github.com/paketo-buildpacks/procfile/cmd/main
GOOS="windows" GOARCH="amd64" go build -ldflags='-s -w' -o bin/main.exe github.com/paketo-buildpacks/procfile/cmd/main

if [ "${STRIP:-false}" != "false" ]; then
  strip bin/main bin/main.exe
fi

if [ "${COMPRESS:-false}" != "false" ]; then
  upx -q -9 bin/main bin/main.exe
fi

ln -fs main bin/build
ln -fs main bin/detect
ln -fs main.exe bin/build.exe
ln -fs main.exe bin/detect.exe
