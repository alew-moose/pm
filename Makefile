.PHONY: build test lint

build:
	go build -o pm ./cmd/pm

test:
	go test -v ./...

lint:
	golangci-lint run
