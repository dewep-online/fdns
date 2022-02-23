SHELL=/bin/bash

.PHONY: run_back run_front
run_back:
	go generate ./...
	go run -race cmd/fdns/main.go -config=./configs/config.dev.yaml
run_front:
	cd web && npm ci --no-delete --cache=/tmp && npm run start

.PHONY: build_back build_font
build_back:
	bash scripts/build.sh back
build_font:
	bash scripts/build.sh front

.PHONY: linter
linter:
	bash scripts/linter.sh

.PHONY: tests
tests:
	bash scripts/tests.sh

.PHONY: develop_up develop_down
develop_up:
	bash scripts/docker.sh docker_up
develop_down:
	bash scripts/docker.sh docker_down

.PHONY: ci
ci:
	bash scripts/ci.sh

deb: build_font build_back
	deb-builder build
