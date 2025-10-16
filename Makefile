# Makefile for Digiflazz Gateway

.PHONY: build run test clean docker-build docker-run help build-all build-linux build-windows build-darwin build-freebsd

# Variables
BINARY_NAME=gateway-digiflazz
DOCKER_IMAGE=gateway-digiflazz
DOCKER_TAG=latest
VERSION?=latest
BUILD_DIR=build

# Build the application for current platform
build:
	@echo "Building $(BINARY_NAME)..."
	@CGO_ENABLED=1 go build -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" -o bin/$(BINARY_NAME) ./cmd/server

# Build for all platforms
build-all:
	@echo "Building for all platforms..."
	@chmod +x scripts/build.sh
	@./scripts/build.sh all all $(VERSION)

# Build for all platforms (CGO-safe)
build-all-safe:
	@echo "Building for all platforms (CGO-safe)..."
	@chmod +x scripts/build-cgo-safe.sh
	@./scripts/build-cgo-safe.sh all all $(VERSION)

# Build for Linux
build-linux:
	@echo "Building for Linux..."
	@chmod +x scripts/build.sh
	@./scripts/build.sh linux all $(VERSION)

# Build for Windows
build-windows:
	@echo "Building for Windows..."
	@chmod +x scripts/build.sh
	@./scripts/build.sh windows all $(VERSION)

# Build for Windows with CGO using Wine
build-windows-cgo:
	@echo "Building for Windows with CGO (using Wine)..."
	@chmod +x scripts/build-windows-wine.sh
	@./scripts/build-windows-wine.sh $(VERSION)

# Setup GitHub Actions
setup-github-actions:
	@echo "Setting up GitHub Actions..."
	@chmod +x scripts/setup-github-actions.sh
	@./scripts/setup-github-actions.sh

# Update GitHub Actions to latest versions
update-github-actions:
	@echo "Updating GitHub Actions to latest versions..."
	@chmod +x scripts/update-github-actions.sh
	@./scripts/update-github-actions.sh

# Test binary functionality
test-binary:
	@echo "Testing binary functionality..."
	@chmod +x scripts/test-binary.sh
	@./scripts/test-binary.sh bin/$(BINARY_NAME)

# Build for macOS
build-darwin:
	@echo "Building for macOS..."
	@chmod +x scripts/build.sh
	@./scripts/build.sh darwin all $(VERSION)

# Build for FreeBSD
build-freebsd:
	@echo "Building for FreeBSD..."
	@chmod +x scripts/build.sh
	@./scripts/build.sh freebsd all $(VERSION)

# Build for specific platform and architecture
build-platform:
	@echo "Building for $(PLATFORM)/$(ARCH)..."
	@chmod +x scripts/build.sh
	@./scripts/build.sh $(PLATFORM) $(ARCH) $(VERSION)

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	@go run ./cmd/server

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/ $(BUILD_DIR)/
	@rm -f coverage.out coverage.html

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run

# Docker build
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Docker run
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

# Docker compose up
docker-up:
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d

# Docker compose down
docker-down:
	@echo "Stopping services with Docker Compose..."
	@docker-compose down

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Generate mocks
mocks:
	@echo "Generating mocks..."
	@go generate ./...

# Cross-compile for specific target
cross-compile:
	@echo "Cross-compiling for $(GOOS)/$(GOARCH)..."
	@CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" -o $(BUILD_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH) ./cmd/server

# Create release packages
release:
	@echo "Creating release packages..."
	@make build-all
	@echo "Release packages created in $(BUILD_DIR)/"

# Help
help:
	@echo "Available commands:"
	@echo "  build              - Build the application for current platform"
	@echo "  build-all          - Build for all platforms"
	@echo "  build-linux        - Build for Linux (all architectures)"
	@echo "  build-windows      - Build for Windows (all architectures)"
	@echo "  build-windows-cgo  - Build for Windows with CGO using Wine"
	@echo "  setup-github-actions - Setup GitHub Actions workflows"
	@echo "  update-github-actions - Update GitHub Actions to latest versions"
	@echo "  test-binary       - Test binary functionality (help, version, config flags)"
	@echo "  build-darwin       - Build for macOS (all architectures)"
	@echo "  build-freebsd      - Build for FreeBSD (all architectures)"
	@echo "  build-platform     - Build for specific platform (PLATFORM=linux ARCH=amd64)"
	@echo "  cross-compile      - Cross-compile for specific target (GOOS=linux GOARCH=amd64)"
	@echo "  release            - Create release packages for all platforms"
	@echo "  run                - Run the application"
	@echo "  test               - Run tests"
	@echo "  test-coverage      - Run tests with coverage"
	@echo "  clean              - Clean build artifacts"
	@echo "  fmt                - Format code"
	@echo "  lint               - Lint code"
	@echo "  docker-build       - Build Docker image"
	@echo "  docker-run         - Run Docker container"
	@echo "  docker-up          - Start services with Docker Compose"
	@echo "  docker-down        - Stop services with Docker Compose"
	@echo "  deps               - Install dependencies"
	@echo "  mocks              - Generate mocks"
	@echo "  help               - Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make build-all                    # Build for all platforms"
	@echo "  make build-platform PLATFORM=linux ARCH=amd64  # Build for Linux x64"
	@echo "  make cross-compile GOOS=windows GOARCH=amd64   # Cross-compile for Windows x64"
	@echo "  make release VERSION=v1.0.0       # Create release with version"
