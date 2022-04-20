#!/bin/bash

export ROOT="$PWD"

export GO_MAIN="$ROOT/cmd/fdns"
export TOOLS_BIN="$ROOT/.tools"

export COVERALLS_TOKEN="${COVERALLS_TOKEN:-dev}"


dependencies() {
    if [ ! -d $TOOLS_BIN ]; then
        mkdir -p $TOOLS_BIN
    fi

    if [ ! -f $TOOLS_BIN/golangci-lint ]; then
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $TOOLS_BIN v1.38.0
    fi

    if [ ! -f $TOOLS_BIN/goveralls ]; then
        GOBIN=$TOOLS_BIN go install github.com/mattn/goveralls@latest
    fi

    if [ ! -f $TOOLS_BIN/static ]; then
        GOBIN=$TOOLS_BIN go install github.com/deweppro/go-static/cmd/static@latest
    fi

    if [ ! -f $TOOLS_BIN/easyjson ]; then
        GOBIN=$TOOLS_BIN go install github.com/mailru/easyjson/...@latest
    fi
}

lints() {
    go mod download
    go build -race -v -o /tmp/bin.test $GO_MAIN
    GOBIN=$TOOLS_BIN go generate ./...
    $TOOLS_BIN/golangci-lint -v run ./...
}

tests() {
    go clean -testcache
    go test -v -race -run Integration ./...

    if [ "$COVERALLS_TOKEN" == "dev" ]; then
        go test -v -race -run Unit ./...
    else
        go test -v -race -run Unit -covermode=atomic -coverprofile=coverage.out ./...
        $TOOLS_BIN/goveralls -coverprofile=coverage.out -repotoken $COVERALLS_TOKEN
    fi
}