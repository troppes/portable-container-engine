.PHONY: all build clean test build-all

# Variables
BINARY_NAME=pce
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

# Build targets for different platforms
PLATFORMS=linux-amd64 windows-amd64 darwin-amd64 darwin-arm64
PLATFORM_BUILDS=$(addprefix build-,$(PLATFORMS))

# Default target
all: clean build test

# Build for current platform
build:
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME) cmd/pce/main.go

# Build for all platforms
build-all: $(PLATFORM_BUILDS)

build-linux-amd64:
	@echo "Building for Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 cmd/pce/main.go

build-windows-amd64:
	@echo "Building for Windows (amd64)..."
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe cmd/pce/main.go

build-darwin-amd64:
	@echo "Building for macOS (amd64)..."
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 cmd/pce/main.go

build-darwin-arm64:
	@echo "Building for macOS (arm64)..."
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 cmd/pce/main.go

# Run all tests
test: test-unit test-integration

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	@go test -v -skip Integration ./...

# Run integration tests only
test-integration:
	@echo "Running integration tests..."
	@go test -v -run Integration ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean

# Install locally
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp bin/$(BINARY_NAME) $(GOPATH)/bin/

# Development helpers
dev: build
	@./bin/$(BINARY_NAME)

# Show help
help:
	@echo "Available targets:"
	@echo "  all          - Clean, build, and test"
	@echo "  build        - Build for current platform"
	@echo "  build-all    - Build for all platforms"
	@echo "  test         - Run tests"
	@echo "  clean        - Remove build artifacts"
	@echo "  install      - Install to GOPATH/bin"
	@echo "  dev          - Build and run locally"
	@echo "  help         - Show this help message"
