.PHONY: build build-go build-windows build-macos build-all build-release clean test lint help

# Binary name
BINARY_NAME=convert
CMD_PATH=./cmd/convert

# Build directory
BIN_DIR=bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-s -w
CGO_ENABLED=1

# Version for release build (same as github.ref_name in release workflow)
# Use --abbrev=0 to get only tag (e.g. v0.0.4) without commit suffix
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build-go: ## Build for current platform
	@echo "Building for current platform..."
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_PATH)

build-release: ## Build with release ldflags (same as CI, for local verification)
	@echo "Building with release ldflags (Version=$(VERSION))..."
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=$(CGO_ENABLED) $(GOBUILD) -ldflags "-s -w -X github.com/sugiyan97/heic-image-converter-cli/internal/cli.Version=$(VERSION)" -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_PATH)

build-windows: ## Build for Windows (amd64)
	@echo "Building for Windows (amd64)..."
	@mkdir -p $(BIN_DIR)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_PATH)

build-macos-arm: ## Build for macOS (Apple Silicon/arm64)
	@echo "Building for macOS (Apple Silicon/arm64)..."
	@mkdir -p $(BIN_DIR)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=$(CGO_ENABLED) $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_PATH)

build-macos: build-macos-arm ## Build for macOS (Apple Silicon only)

build-all: build-windows build-macos ## Build for all platforms (Windows and macOS)

clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BIN_DIR)

test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v -parallel $(shell nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 4) ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -v -parallel $(shell nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 4) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

lint: ## Run linter
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run --config golangci.yml --timeout=5m; \
	else \
		echo "golangci-lint is not installed. Install it from https://golangci-lint.run/"; \
	fi

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

build: build-go ## Alias for build-go

