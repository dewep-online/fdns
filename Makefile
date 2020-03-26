RD=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

.PHONY: run
run:
	@go run cmd/fdns/main.go --config=$(RD)/configs/config.yaml

.PHONY: build
build:
	@bash $(RD)/build/build.sh
