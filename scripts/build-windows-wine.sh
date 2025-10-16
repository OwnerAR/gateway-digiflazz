#!/bin/bash

# Windows build script using Wine for CGO cross-compilation
# This script uses Wine to provide Windows headers for cross-compilation

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

# Check if Wine is installed
check_wine() {
    if ! command -v wine &> /dev/null; then
        print_error "Wine is not installed. Please install Wine first:"
        echo "  macOS: brew install wine-stable"
        echo "  Ubuntu/Debian: sudo apt install wine"
        echo "  CentOS/RHEL: sudo yum install wine"
        exit 1
    fi
    
    print_status "Found Wine: $(wine --version)"
}

# Check if MinGW-w64 is installed
check_mingw() {
    if ! command -v x86_64-w64-mingw32-gcc &> /dev/null; then
        print_error "MinGW-w64 is not installed. Please install it first:"
        echo "  macOS: brew install mingw-w64"
        echo "  Ubuntu/Debian: sudo apt install mingw-w64"
        echo "  CentOS/RHEL: sudo yum install mingw64-gcc"
        exit 1
    fi
    
    print_status "Found MinGW-w64: $(x86_64-w64-mingw32-gcc --version | head -n1)"
}

# Clean build directory
clean_build() {
    print_status "Cleaning build directory..."
    rm -rf $BUILD_DIR
    mkdir -p $BUILD_DIR
}

# Build for Windows using Wine
build_windows() {
    local arch=$1
    local output_name="${BINARY_NAME}-windows-${arch}.exe"
    
    print_status "Building for windows/${arch} with CGO enabled..."
    
    # Set environment variables for cross-compilation
    export CGO_ENABLED=1
    export GOOS=windows
    export GOARCH=$arch
    
    if [ "$arch" = "amd64" ]; then
        export CC=x86_64-w64-mingw32-gcc
        export CXX=x86_64-w64-mingw32-g++
    elif [ "$arch" = "386" ]; then
        export CC=i686-w64-mingw32-gcc
        export CXX=i686-w64-mingw32-g++
    else
        print_error "Unsupported architecture: $arch"
        exit 1
    fi
    
    # Set Wine environment
    export WINEPREFIX=/tmp/wine-go-build
    export WINEDEBUG=-all
    
    # Create Wine prefix if it doesn't exist
    if [ ! -d "$WINEPREFIX" ]; then
        print_status "Creating Wine prefix..."
        wineboot --init
    fi
    
    # Build with CGO
    go build \
        -ldflags "-X main.version=${VERSION} -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
        -o "${BUILD_DIR}/${output_name}" \
        ./cmd/server
    
    if [ $? -eq 0 ]; then
        print_success "Built ${output_name}"
        
        # Create archive
        cd $BUILD_DIR
        zip "${output_name}.zip" "${output_name}"
        cd ..
        
        print_success "Created archive for ${output_name}"
    else
        print_error "Failed to build for windows/${arch}"
        exit 1
    fi
}

# Main execution
main() {
    print_status "Digiflazz Gateway Windows Build with CGO (using Wine)"
    print_status "Version: $VERSION"
    print_status "Build Directory: $BUILD_DIR"
    echo ""
    
    # Check prerequisites
    check_wine
    check_mingw
    
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
    
    # Build for Windows architectures
    build_windows "amd64"
    build_windows "386"
    
    # Show build summary
    print_status "Build Summary:"
    echo "Build directory: $BUILD_DIR"
    echo "Files created:"
    ls -la $BUILD_DIR/
    
    print_success "Windows build completed successfully!"
    print_warning "Note: The binaries require Windows C runtime libraries to run on Windows."
}

# Handle help flag
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    echo "Windows build script using Wine for CGO cross-compilation"
    echo ""
    echo "Usage: $0 [version]"
    echo ""
    echo "Arguments:"
    echo "  version     Version tag (default: latest)"
    echo ""
    echo "Requirements:"
    echo "  - Wine (for Windows headers)"
    echo "  - MinGW-w64 (for Windows cross-compiler)"
    echo "  - Go 1.21+"
    echo ""
    echo "Examples:"
    echo "  $0                    # Build with default version"
    echo "  $0 v1.0.0            # Build with specific version"
    exit 0
fi

# Run main function
main
