#!/usr/bin/env bash

go build -o http-server

mv http-server pkg/osx


GOOS=linux go build -o http-server

mv http-server pkg/amd-x86_64

