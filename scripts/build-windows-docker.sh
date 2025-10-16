#!/bin/bash

# Windows build script using Docker for CGO cross-compilation
# Alternative to Wine approach

set -e

# Default values
VERSION=${1:-"latest"}
BUILD_DIR="build"
BINARY_NAME="gateway-digiflazz"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is installed
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed. Please install Docker first:"
        echo "  macOS: brew install docker"
        echo "  Or download from: https://www.docker.com/products/docker-desktop"
        exit 1
    fi
    
    print_status "Found Docker: $(docker --version)"
}

# Clean build directory
clean_build() {
    print_status "Cleaning build directory..."
    rm -rf $BUILD_DIR
    mkdir -p $BUILD_DIR
}

# Build using Docker
build_with_docker() {
    print_status "Building Windows binary using Docker..."
    
    # Create temporary Dockerfile for Windows build
    cat > Dockerfile.temp << 'EOF'
FROM golang:1.21-windowsservercore-ltsc2022

WORKDIR /app

# Install build tools
RUN powershell -Command "choco install mingw -y"

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Set environment variables
ENV CGO_ENABLED=1
ENV GOOS=windows
ENV GOARCH=amd64
ENV CC=gcc
ENV CXX=g++

# Build the application
RUN go build -ldflags "-X main.version=latest -X main.buildTime=$(Get-Date -Format 'yyyy-MM-ddTHH:mm:ssZ')" -o gateway-digiflazz-windows-amd64.exe ./cmd/server

# Copy binary out
CMD ["powershell", "-Command", "Copy-Item gateway-digiflazz-windows-amd64.exe /output/"]
EOF

    # Build Docker image
    print_status "Building Docker image..."
    docker build -f Dockerfile.temp -t gateway-windows-builder .
    
    # Create output directory
    mkdir -p $BUILD_DIR
    
    # Run container to extract binary
    print_status "Extracting Windows binary..."
    docker run --rm -v "$(pwd)/$BUILD_DIR:/output" gateway-windows-builder
    
    # Clean up
    rm -f Dockerfile.temp
    
    print_success "Windows binary extracted to $BUILD_DIR/"
}

# Main execution
main() {
    print_status "Digiflazz Gateway Windows Build with CGO (using Docker)"
    print_status "Version: $VERSION"
    print_status "Build Directory: $BUILD_DIR"
    echo ""
    
    # Check prerequisites
    check_docker
    
    # Clean build directory
    clean_build
    
    # Build with Docker
    build_with_docker
    
    # Show build summary
    print_status "Build Summary:"
    echo "Build directory: $BUILD_DIR"
    echo "Files created:"
    ls -la $BUILD_DIR/
    
    print_success "Windows build completed successfully!"
}

# Handle help flag
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    echo "Windows build script using Docker for CGO cross-compilation"
    echo ""
    echo "Usage: $0 [version]"
    echo ""
    echo "Requirements:"
    echo "  - Docker Desktop"
    echo "  - Go 1.21+"
    echo ""
    echo "Examples:"
    echo "  $0                    # Build with default version"
    echo "  $0 v1.0.0            # Build with specific version"
    exit 0
fi

# Run main function
main
