.PHONY: all build test lint clean dev testsuite help

# Default target
all: dev

# Development workflow (equivalent to ./dev.sh)
dev:
	@FORCE_COLOR=1 ./dev.sh

# Build only
build:
	@FORCE_COLOR=1 ./dev.sh --skip-tests

# Run tests
test:
	@go test ./...

# Run linter
lint:
	@golangci-lint run

# Run full testsuite
testsuite:
	@FORCE_COLOR=1 ./dev.sh --testsuite

# Clean build artifacts
clean:
	@rm -rf ./bin/

# Show help
help:
	@echo "Available targets:"
	@echo "  all       - Run full development workflow (default)"
	@echo "  dev       - Run development workflow (./dev.sh)"
	@echo "  build     - Build only (skip tests)"
	@echo "  test      - Run unit tests only"
	@echo "  lint      - Run linter only"
	@echo "  testsuite - Run full testsuite"
	@echo "  clean     - Clean build artifacts"
	@echo "  help      - Show this help"
	@echo ""
	@echo "For more options, use: ./dev.sh --help" 