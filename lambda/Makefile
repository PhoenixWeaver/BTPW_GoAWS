# Makefile for AWS Lambda Go Testing
# This Makefile provides commands for running different types of tests

.PHONY: test test-unit test-integration test-e2e test-all clean build

# Default target
all: test-all

# Run all tests
test-all: test-unit test-integration test-e2e

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	go test -v ./types/... ./database/... ./api/... ./middleware/...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	go test -v -run TestIntegration ./...

# Run end-to-end tests
test-e2e:
	@echo "Running end-to-end tests..."
	go test -v -run TestE2E ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	go test -v -race ./...

# Run tests with benchmarks
test-bench:
	@echo "Running benchmarks..."
	go test -v -bench=. ./...

# Clean test artifacts
clean:
	@echo "Cleaning test artifacts..."
	rm -f coverage.out coverage.html
	go clean -testcache

# Build the Lambda function
build:
	@echo "Building Lambda function..."
	GOOS=linux GOARCH=amd64 go build -o bootstrap
	zip function.zip bootstrap

# Run tests in verbose mode
test-verbose:
	@echo "Running tests in verbose mode..."
	go test -v ./...

# Run specific test package
test-package:
	@echo "Running tests for package: $(PACKAGE)"
	go test -v ./$(PACKAGE)/...

# Run tests with timeout
test-timeout:
	@echo "Running tests with timeout..."
	go test -v -timeout=30s ./...

# Run tests with short flag (skips integration tests)
test-short:
	@echo "Running short tests..."
	go test -v -short ./...

# Run tests with long flag (includes integration tests)
test-long:
	@echo "Running long tests..."
	go test -v -long ./...

# Help target
help:
	@echo "Available targets:"
	@echo "  test-all        - Run all tests"
	@echo "  test-unit       - Run unit tests only"
	@echo "  test-integration- Run integration tests"
	@echo "  test-e2e        - Run end-to-end tests"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  test-race       - Run tests with race detection"
	@echo "  test-bench      - Run benchmarks"
	@echo "  test-verbose    - Run tests in verbose mode"
	@echo "  test-short      - Run short tests (skip integration)"
	@echo "  test-long       - Run long tests (include integration)"
	@echo "  clean          - Clean test artifacts"
	@echo "  build          - Build Lambda function"
	@echo "  help           - Show this help message"
