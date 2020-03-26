#!/usr/bin/env bash

GOOS=windows GOARCH=amd64 go build -o fdns_windows64.exe ./../cmd/fdns/main.go
GOOS=linux GOARCH=amd64 go build -o fdns_amd64 ./../cmd/fdns/main.go
GOOS=darwin GOARCH=amd64 go build -o fdns_macos ./../cmd/fdns/main.go
