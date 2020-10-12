#!/usr/bin/env bash

set -euo pipefail

GOOS="linux" go build -ldflags='-s -w' -o bin/main github.com/paketo-buildpacks/procfile/cmd/main

strip bin/main
upx -q -9 bin/main

ln -fs main bin/build
ln -fs main bin/detect
