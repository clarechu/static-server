#!/usr/bin/env bash

set -e

mkdir -p pkg/osx pkg/amd-x86_64 pkg/win-x86_64

# ==== osx
go build -o http-server

mv http-server pkg/osx


# ==== linux amd

GOOS=linux go build -o http-server

docker build -t clarechu/http-server:v0.1.0 .


mv http-server pkg/amd-x86_64


# === windows

GOOS=windows go build -o http-server.exe

mv http-server.exe pkg/win-x86_64


cd pkg

tar -cvf http-server-v0.3-linux-amd-x86_64.tar.gz amd-x86_64
tar -cvf http-server-v0.3-macos-darwin.tar.gz osx
tar -cvf http-server-v0.3-win-x86_64.tar.gz win-x86_64
