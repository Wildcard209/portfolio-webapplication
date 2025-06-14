// Makefile for the portfolio web application backend
.PHONY: docs build run dev clean

# Generate Swagger documentation
docs:
	swag init -g main.go -o ./docs

# Build the application
build: docs
	go build -o main .

# Run the application
run: build
	./main

# Run in development mode with hot reload (requires air)
dev: docs
	air

# Clean build artifacts
clean:
	rm -f main
	rm -rf docs/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run tests
test:
	go test ./...

# Install swag if not present
install-swag:
	go install github.com/swaggo/swag/cmd/swag@latest
