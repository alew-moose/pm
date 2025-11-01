.PHONY: run build

run:
	make build
	./pm

build:
	go build -o pm ./cmd/pm
