#!/bin/bash

# Safe CGO build script that handles cross-compilation issues
# This script builds with CGO only for compatible platforms

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

# Check if platform/arch combination supports CGO
supports_cgo() {
    local platform=$1
    local arch=$2
    
    case "$platform" in
        "linux")
            # Linux ARM64 has CGO cross-compilation issues
            if [ "$arch" = "arm64" ] || [ "$arch" = "arm" ]; then
                return 1
            fi
            ;;
        "windows")
            # Windows 386 has compatibility issues
            if [ "$arch" = "386" ]; then
                return 1
            fi
            ;;
        "darwin")
            # macOS supports both architectures
            return 0
            ;;
        *)
            return 1
            ;;
    esac
    
    return 0
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
    
    # Check if CGO is supported for this combination
    if supports_cgo "$platform" "$arch"; then
        print_status "Building with CGO enabled for ${platform}/${arch}..."
        CGO_ENABLED=1 GOOS=$platform GOARCH=$arch go build \
            -ldflags "-X main.version=${VERSION} -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
            -o "${BUILD_DIR}/${output_name}" \
            ./cmd/server
    else
        print_warning "CGO not supported for ${platform}/${arch}, building without CGO..."
        print_warning "Note: SQLite cache will not work without CGO"
        CGO_ENABLED=0 GOOS=$platform GOARCH=$arch go build \
            -ldflags "-X main.version=${VERSION} -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
            -o "${BUILD_DIR}/${output_name}" \
            ./cmd/server
    fi
    
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

# Build all platforms (safe combinations only)
build_all() {
    print_status "Building for all platforms (CGO-safe combinations)..."
    
    # Linux
    build_platform "linux" "amd64"
    
    # Windows
    build_platform "windows" "amd64"
    
    # macOS
    build_platform "darwin" "amd64"
    build_platform "darwin" "arm64"
    
    print_success "All builds completed!"
    print_warning "Note: ARM64 builds for Linux/Windows require native runners due to CGO limitations"
}

# Build specific platform
build_specific() {
    if [ "$ARCH" = "all" ]; then
        case $PLATFORM in
            "linux")
                build_platform "linux" "amd64"
                print_warning "Linux ARM64/ARM builds disabled due to CGO cross-compilation issues"
                print_warning "Use GitHub Actions with ARM64 runners for ARM builds"
                ;;
            "windows")
                build_platform "windows" "amd64"
                print_warning "Windows ARM64/386 builds may have compatibility issues"
                ;;
            "darwin")
                build_platform "darwin" "amd64"
                build_platform "darwin" "arm64"
                ;;
            "freebsd")
                build_platform "freebsd" "amd64"
                ;;
            *)
                print_error "Unsupported platform: $PLATFORM"
                exit 1
                ;;
        esac
    else
        if supports_cgo "$PLATFORM" "$ARCH"; then
            build_platform "$PLATFORM" "$ARCH"
        else
            print_warning "CGO not supported for ${PLATFORM}/${ARCH}"
            print_warning "Building without CGO (SQLite will not work)"
            build_platform "$PLATFORM" "$ARCH"
        fi
    fi
}

# Show help
show_help() {
    echo "Safe CGO build script for Digiflazz Gateway"
    echo ""
    echo "Usage: $0 [platform] [arch] [version]"
    echo ""
    echo "Arguments:"
    echo "  platform    Target platform (linux, windows, darwin, freebsd, all)"
    echo "  arch        Target architecture (amd64, arm64, arm, all)"
    echo "  version     Version tag (default: latest)"
    echo ""
    echo "CGO Support:"
    echo "  - linux/amd64: ✅ CGO supported"
    echo "  - linux/arm64: ❌ CGO not supported (cross-compilation issues)"
    echo "  - windows/amd64: ✅ CGO supported"
    echo "  - windows/arm64: ⚠️  CGO may have issues"
    echo "  - darwin/amd64: ✅ CGO supported"
    echo "  - darwin/arm64: ✅ CGO supported"
    echo ""
    echo "Examples:"
    echo "  $0                          # Build for all CGO-safe platforms"
    echo "  $0 linux                    # Build for Linux amd64"
    echo "  $0 windows amd64           # Build for Windows x64"
    echo "  $0 darwin arm64            # Build for macOS Apple Silicon"
    echo "  $0 linux amd64 v1.0.0      # Build for Linux x64 with version"
    echo ""
    echo "For ARM64 builds, use GitHub Actions with native runners"
}

# Main execution
main() {
    print_status "Digiflazz Gateway Safe CGO Build"
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
