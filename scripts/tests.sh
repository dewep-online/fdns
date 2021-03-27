#!/bin/bash

cd $PWD

go clean -testcache
go test -v -race -run Unit ./...
go test -v -race -run Integration ./...