.PHONY: build test

build:
	go build -o build/teredix main.go

test:
	go test -v ./...
