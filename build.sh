#!/bin/bash

PACKAGE="pilot-sysmon-backend"
VERSION="v1.0.0-beta"
COMMIT_HASH="$(git rev-parse --short HEAD)"
BUILD_TIME="$(date '+%Y-%m-%d %H:%M:%S')"

LDFLAGS=(
  "-X '${PACKAGE}/endpoints.VERSION=${VERSION}'"
  "-X '${PACKAGE}/endpoints.COMMIT_HASH=${COMMIT_HASH}'"
  "-X '${PACKAGE}/endpoints.BUILD_TIME=${BUILD_TIME}'"
)

go build -ldflags="${LDFLAGS[*]}"
