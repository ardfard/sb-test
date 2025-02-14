.PHONY: all build run test lint smoke_test clean

# Default target: build the server
all: build

# Build the Go executable from cmd/server/main.go into the bin/ directory.
build:
	@echo "Building the project..."
	mkdir -p bin
	go build -o bin/sb-test ./cmd/server

# Run the server executable.
run: build
	@echo "Running the server..."
	./bin/sb-test

# Run all unit tests.
test:
	@echo "Running unit tests..."
	go test -v ./...

# Run the linter using golangci-lint (ensure it's installed).
lint:
	@echo "Running linters..."
	golangci-lint run --timeout=5m

# Run the smoke test script.
smoke_test:
	@echo "Running smoke test..."
	./smoke_test.sh

# Clean up build artifacts.
clean:
	@echo "Cleaning up..."
	rm -rf bin 