#!/usr/bin/sh

# *nix build script

set -xe

go build -ldflags="-s -w" -o ./bin/rpseek
