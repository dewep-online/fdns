#!/bin/bash

GO_FILES=$(find . -name '*.go' | grep -vE 'vendor|easyjson|static')


cd $PWD

go generate ./...
goimports -w $GO_FILES
go fmt ./...
golangci-lint -v run ./...