#!/bin/bash

#################################################
source $(dirname "$0")/env.sh
cd $ROOT
dependencies
#################################################

front() {
  cd web && npm ci --no-delete --cache=/tmp && npm run build
}

# sudo apt install gcc-aarch64-linux-gnu

case $1 in
arm64)
  rm -rf $ROOT/build/bin/fdns_$1 && \
  GOBIN=$TOOLS_BIN go generate ./... && \
  GO111MODULE=on CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc \
      go build -a -o $ROOT/build/bin/fdns_$1 $GO_MAIN
  ;;
front)
  front
  ;;
*)
  rm -rf $ROOT/build/bin/fdns_$1 && \
  GOBIN=$TOOLS_BIN go generate ./... && \
  GO111MODULE=on CGO_ENABLED=1 GOOS=linux GOARCH=$1 go build -o $ROOT/build/bin/fdns_$1 $GO_MAIN
  ;;
esac