.PHONY: run
run:
	go run -race main.go -config=./config/config.yaml

.PHONY: build
build:
	bash ./build/build.sh
