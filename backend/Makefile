# Makefile для ARK CLI

.PHONY: build clean test help

# Переменные
BINARY_NAME=ark
BUILD_DIR=build
VERSION=1.0.0
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

# Цвета для вывода
GREEN=\033[0;32m
NC=\033[0m # No Color

# Основные цели
build: clean
	@echo "$(GREEN)Building ARK CLI...$(NC)"
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/ark
	@echo "$(GREEN)Build completed: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

build-all: clean
	@echo "$(GREEN)Building ARK CLI for all platforms...$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/ark
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/ark
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/ark
	@echo "$(GREEN)Cross-platform build completed$(NC)"

clean:
	@echo "$(GREEN)Cleaning build directory...$(NC)"
	@rm -rf $(BUILD_DIR)

test:
	@echo "$(GREEN)Running tests...$(NC)"
	go test ./...

install: build
	@echo "$(GREEN)Installing ARK CLI...$(NC)"
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)ARK CLI installed to /usr/local/bin/$(BINARY_NAME)$(NC)"

# Цели для разработки
dev-build:
	@echo "$(GREEN)Building for development...$(NC)"
	go build -race -o $(BUILD_DIR)/$(BINARY_NAME)-dev ./cmd/ark

lint:
	@echo "$(GREEN)Running linter...$(NC)"
	golangci-lint run

lint-full:
	@echo "$(GREEN)Running full linting with all checks...$(NC)"
	golangci-lint run --config .golangci.yml

lint-deadcode:
	@echo "$(GREEN)Checking for dead code...$(NC)"
	golangci-lint run --disable-all --enable unused,ineffassign

lint-dupl:
	@echo "$(GREEN)Checking for code duplication...$(NC)"
	golangci-lint run --disable-all --enable dupl

analyze-complexity:
	@echo "$(GREEN)Analyzing code complexity...$(NC)"
	golangci-lint run --disable-all --enable gocyclo,gocognit

lint-report:
	@echo "$(GREEN)Generating lint report...$(NC)"
	golangci-lint run --config .golangci.yml --out-format json > golangci-lint-report.json || true
	@echo "$(GREEN)Report saved to golangci-lint-report.json$(NC)"

format:
	@echo "$(GREEN)Formatting code...$(NC)"
	go fmt ./...

# Цели для CLI команд
test-index:
	@echo "$(GREEN)Testing index command...$(NC)"
	./$(BUILD_DIR)/$(BINARY_NAME) index --help

test-solve:
	@echo "$(GREEN)Testing solve command...$(NC)"
	./$(BUILD_DIR)/$(BINARY_NAME) solve --help

test-result:
	@echo "$(GREEN)Testing result command...$(NC)"
	./$(BUILD_DIR)/$(BINARY_NAME) result --help

# Performance targets
perf-test:
	@echo "$(GREEN)Running performance tests...$(NC)"
	go test -bench=. -benchmem ./...

bench:
	@echo "$(GREEN)Running benchmarks...$(NC)"
	go test -bench=. -benchmem -benchtime=5s ./... | tee benchmark.txt

pprof-cpu:
	@echo "$(GREEN)Running CPU profiling...$(NC)"
	go test -cpuprofile=cpu.prof -bench=. ./internal/context
	go tool pprof -http=:8080 cpu.prof

pprof-mem:
	@echo "$(GREEN)Running memory profiling...$(NC)"
	go test -memprofile=mem.prof -bench=. ./internal/context
	go tool pprof -http=:8080 mem.prof

leak-check:
	@echo "$(GREEN)Running leak detection tests...$(NC)"
	go test -race -count=1 ./...

stress-test:
	@echo "$(GREEN)Running stress tests...$(NC)"
	go test -race -count=10 -timeout=10m ./...

coverage:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report: coverage.html$(NC)"

# Справка
help:
	@echo "$(GREEN)ARK CLI Makefile$(NC)"
	@echo ""
	@echo "Available targets:"
	@echo "  build       - Build ARK CLI for current platform"
	@echo "  build-all   - Build ARK CLI for all platforms (Linux, Windows, macOS)"
	@echo "  clean       - Clean build directory"
	@echo "  test        - Run tests"
	@echo "  install     - Install ARK CLI to /usr/local/bin"
	@echo "  dev-build   - Build with race detection for development"
	@echo "  lint        - Run linter"
	@echo "  format      - Format code"
	@echo "  test-index  - Test index command"
	@echo "  test-solve  - Test solve command"
	@echo "  test-result - Test result command"
	@echo "  perf-test   - Run performance tests"
	@echo "  bench       - Run benchmarks"
	@echo "  pprof-cpu   - Run CPU profiling"
	@echo "  pprof-mem   - Run memory profiling"
	@echo "  leak-check  - Run leak detection tests"
	@echo "  stress-test - Run stress tests"
	@echo "  coverage    - Run tests with coverage"
	@echo "  help        - Show this help"
