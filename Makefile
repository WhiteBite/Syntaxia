.PHONY: help build dev clean test

help:
	@echo "Syntaxia - Build Commands"
	@echo ""
	@echo "Development:"
	@echo "  make dev          - Run development server with hot reload"
	@echo ""
	@echo "Building:"
	@echo "  make build        - Build for current platform"
	@echo ""
	@echo "Maintenance:"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make test         - Run tests"

dev:
	@echo "Starting development server..."
	@cd backend && wails dev

build:
	@echo "Building for current platform..."
	@mkdir -p build/bin
	@cd backend && wails build -clean

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf build
	@rm -rf backend/build
	@cd backend && go clean

test:
	@echo "Running tests..."
	@cd backend && go test ./...
