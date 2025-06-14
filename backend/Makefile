// Makefile for the portfolio web application backend
.PHONY: docs build run dev clean

docs:
	swag init -g main.go -o ./docs

build: docs
	go build -o main .

run: build
	./main

dev: docs
	air

clean:
	rm -f main
	rm -rf docs/

deps:
	go mod download
	go mod tidy

fmt:
	go fmt ./...

test:
	go test ./...

install-swag:
	go install github.com/swaggo/swag/cmd/swag@latest
