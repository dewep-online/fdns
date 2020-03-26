#!/usr/bin/env bash

root=`dirname $0`

cd ${root}

mkdir bin

GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o fdns_amd64 cmd/fdns/main.go
#GO111MODULE=on GOOS=windows GOARCH=amd64 go build -o fdns_windows64.exe cmd/fdns/main.go
#GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o fdns_macos cmd/fdns/main.go
