.PHONY: build test

all: test build

build:
	go build -o build/teredix ./main.go

test:
	go test -v ./...
