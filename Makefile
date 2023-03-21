.PHONY: build test

all: test build

build:
	go build -o build/teredix

test:
	go test -v ./...
