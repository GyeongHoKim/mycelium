.PHONY: build run test lint fmt clean all release setup

BINARY_NAME := mycelium
BUILD_DIR := ./bin
CMD_DIR := ./cmd/mycelium

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

# Default target
all: lint test build

# Build the daemon binary
build:
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

# Run the daemon
run: build
	$(BUILD_DIR)/$(BINARY_NAME)

# Run tests
test:
	go test -v -race ./...

# Run tests with coverage
test-cover:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	goimports -w .
	golines -w .

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Cross-compile for all platforms
release:
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD_DIR)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)

# Install dependencies
deps:
	go mod download
	go mod tidy

# Install dev tools
tools:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/segmentio/golines@latest

# Setup development environment (install hooks)
setup: tools
	npm install
	npx lefthook install

# Help
help:
	@echo "Available targets:"
	@echo "  all        - Run lint, test, and build (default)"
	@echo "  build      - Build the daemon binary"
	@echo "  run        - Build and run the daemon"
	@echo "  test       - Run tests"
	@echo "  test-cover - Run tests with coverage report"
	@echo "  lint       - Run golangci-lint"
	@echo "  fmt        - Format code with goimports and golines"
	@echo "  clean      - Remove build artifacts"
	@echo "  release    - Cross-compile for all platforms"
	@echo "  deps       - Download and tidy dependencies"
	@echo "  tools      - Install development tools"
	@echo "  setup      - Setup development environment (tools + git hooks)"
