
.PHONY: install
install:
	go install github.com/osspkg/devtool@latest

.PHONY: setup
setup:
	devtool setup-lib

.PHONY: lint
lint:
	devtool lint

.PHONY: license
license:
	devtool license

.PHONY: build
build:
	devtool build --arch=amd64

.PHONY: tests
tests:
	devtool test

.PHONY: pre-commite
pre-commite: setup lint build tests

.PHONY: ci
ci: install setup lint build tests

run_back:
	go run cmd/fdns/main.go --config=config/config.dev.yaml

nslookup:
	nslookup -port=8053 google.com 127.0.0.1
	nslookup -port=8053 adstop.org 127.0.0.1
	nslookup -port=8053 yandex.ru 127.0.0.1
	nslookup -port=8053 vk.com 127.0.0.1
	nslookup -port=8053 dewep.pro 127.0.0.1
	nslookup -port=8053 dewep.online 127.0.0.1
	nslookup -port=8053 googleads.github.io 127.0.0.1
