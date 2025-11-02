.PHONY: run build test lint

run:
	make build
	./pm

build:
	go build -o pm ./cmd/pm

test:
	go test -v ./...

lint:
	golangci-lint run
