# Default makefile for rn
# NOTE: this isn't really needed since the container is using Air for hot reloading
.DEFAULT_GOAL := build
fmt: 
	@go fmt ./...

lint: fmt
	@golint ./...

vet: fmt
	@go vet ./...

build: vet
	@go build -o bin/main cmd/api/main.go

run:
	@go run cmd/api/main.go

clean:
	@go clean
	rm bin/main

.PHONY: fmt lint vet build clean
