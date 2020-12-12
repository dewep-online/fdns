#!/usr/bin/env bash

root=`dirname $0`
mainfile=${root}/../../..

cd ${root}/bin
rm -rf ./fdns*

GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o fdns_amd64 ${mainfile}
GO111MODULE=on GOOS=windows GOARCH=amd64 go build -o fdns_windows64.exe ${mainfile}
GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o fdns_macos ${mainfile}
