#!/bin/bash

# Cross-platform build script for Digiflazz Gateway
# Usage: ./scripts/build.sh [platform] [arch]

set -e

# Default values
PLATFORM=${1:-"all"}
ARCH=${2:-"all"}
VERSION=${3:-"latest"}
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

# Clean build directory
clean_build() {
    print_status "Cleaning build directory..."
    rm -rf $BUILD_DIR
    mkdir -p $BUILD_DIR
}

# Build for specific platform and architecture
build_platform() {
    local platform=$1
    local arch=$2
    local output_name="${BINARY_NAME}-${platform}-${arch}"
    
    if [ "$platform" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    print_status "Building for ${platform}/${arch}..."
    
    CGO_ENABLED=1 GOOS=$platform GOARCH=$arch go build \
        -ldflags "-X main.version=${VERSION} -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
        -o "${BUILD_DIR}/${output_name}" \
        ./cmd/server
    
    if [ $? -eq 0 ]; then
        print_success "Built ${output_name}"
        
        # Create archive
        if [ "$platform" = "windows" ]; then
            zip "${BUILD_DIR}/${output_name}.zip" "${BUILD_DIR}/${output_name}"
        else
            tar -czf "${BUILD_DIR}/${output_name}.tar.gz" -C "${BUILD_DIR}" "${output_name}"
        fi
        
        print_success "Created archive for ${output_name}"
    else
        print_error "Failed to build for ${platform}/${arch}"
        exit 1
    fi
}

# Build all platforms
build_all() {
    print_status "Building for all platforms..."
    
    # Linux
    build_platform "linux" "amd64"
    build_platform "linux" "arm64"
    build_platform "linux" "arm"
    
    # Windows
    build_platform "windows" "amd64"
    build_platform "windows" "arm64"
    
    # macOS
    build_platform "darwin" "amd64"
    build_platform "darwin" "arm64"
    
    # FreeBSD
    build_platform "freebsd" "amd64"
    build_platform "freebsd" "arm64"
    
    print_success "All builds completed!"
}

# Build specific platform
build_specific() {
    if [ "$ARCH" = "all" ]; then
        case $PLATFORM in
            "linux")
                build_platform "linux" "amd64"
                build_platform "linux" "arm64"
                build_platform "linux" "arm"
                ;;
            "windows")
                build_platform "windows" "amd64"
                build_platform "windows" "arm64"
                ;;
            "darwin")
                build_platform "darwin" "amd64"
                build_platform "darwin" "arm64"
                ;;
            "freebsd")
                build_platform "freebsd" "amd64"
                build_platform "freebsd" "arm64"
                ;;
            *)
                print_error "Unsupported platform: $PLATFORM"
                exit 1
                ;;
        esac
    else
        build_platform "$PLATFORM" "$ARCH"
    fi
}

# Show help
show_help() {
    echo "Cross-platform build script for Digiflazz Gateway"
    echo ""
    echo "Usage: $0 [platform] [arch] [version]"
    echo ""
    echo "Arguments:"
    echo "  platform    Target platform (linux, windows, darwin, freebsd, all)"
    echo "  arch        Target architecture (amd64, arm64, arm, all)"
    echo "  version     Version tag (default: latest)"
    echo ""
    echo "Examples:"
    echo "  $0                          # Build for all platforms"
    echo "  $0 linux                    # Build for all Linux architectures"
    echo "  $0 windows amd64           # Build for Windows x64"
    echo "  $0 darwin arm64            # Build for macOS Apple Silicon"
    echo "  $0 linux amd64 v1.0.0      # Build for Linux x64 with version"
    echo ""
    echo "Supported platforms:"
    echo "  - linux (amd64, arm64, arm)"
    echo "  - windows (amd64, arm64)"
    echo "  - darwin (amd64, arm64)"
    echo "  - freebsd (amd64, arm64)"
}

# Main execution
main() {
    print_status "Digiflazz Gateway Cross-Platform Build"
    print_status "Platform: $PLATFORM"
    print_status "Architecture: $ARCH"
    print_status "Version: $VERSION"
    print_status "Build Directory: $BUILD_DIR"
    echo ""
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    # Check Go version
    go_version=$(go version | cut -d' ' -f3)
    print_status "Using Go version: $go_version"
    
    # Clean build directory
    clean_build
    
    # Build based on parameters
    if [ "$PLATFORM" = "all" ]; then
        build_all
    else
        build_specific
    fi
    
    # Show build summary
    print_status "Build Summary:"
    echo "Build directory: $BUILD_DIR"
    echo "Files created:"
    ls -la $BUILD_DIR/
    
    print_success "Build completed successfully!"
}

# Handle help flag
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    show_help
    exit 0
fi

# Run main function
main

