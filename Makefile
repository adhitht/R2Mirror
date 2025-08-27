# Ubuntu Release Downloader - Simplified Makefile

BINARY_NAME=ubuntu-release-downloader
BUILD_DIR=build
MAIN_FILE=cmd/main.go

.PHONY: all build clean test run help setup

# Default target
all: clean build

# Build the application
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
build-all:
	@echo "🔨 Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)
	@echo "✅ Multi-platform build complete"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...
	@echo "✅ Tests complete"

# Run the application in development mode
run:
	@echo "🚀 Running $(BINARY_NAME)..."
	go run $(MAIN_FILE)

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download
	@echo "✅ Dependencies installed"

# Setup development environment
setup: deps
	@echo "🛠️  Setting up development environment..."
	@if [ ! -f .env ]; then \
		echo "Creating .env file..."; \
		go run $(MAIN_FILE) & \
		sleep 2; \
		pkill -f $(MAIN_FILE) || true; \
	fi
	@echo "✅ Setup complete. Edit .env and config.yaml files."

# Show help
help:
	@echo "Available commands:"
	@echo "  build      - Build the application"
	@echo "  build-all  - Build for multiple platforms"
	@echo "  clean      - Clean build artifacts"
	@echo "  test       - Run tests"
	@echo "  run        - Run in development mode"
	@echo "  deps       - Install dependencies"
	@echo "  setup      - Setup development environment"
	@echo "  help       - Show this help"