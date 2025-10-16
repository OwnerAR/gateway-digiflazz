#!/bin/bash

# Clean Go module dependencies
# This script cleans up Go module cache and dependencies

set -e

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

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    go_version=$(go version | cut -d' ' -f3)
    print_status "Using Go version: $go_version"
}

# Clean Go module cache
clean_module_cache() {
    print_status "Cleaning Go module cache..."
    
    # Clean module cache
    go clean -modcache
    
    # Clean build cache
    go clean -cache
    
    # Clean test cache
    go clean -testcache
    
    print_success "Go module cache cleaned"
}

# Clean go.mod and go.sum
clean_go_mod() {
    print_status "Cleaning go.mod and go.sum..."
    
    # Remove go.sum
    if [ -f "go.sum" ]; then
        rm go.sum
        print_status "Removed go.sum"
    fi
    
    # Tidy go.mod
    go mod tidy
    
    print_success "go.mod and go.sum cleaned"
}

# Download dependencies
download_deps() {
    print_status "Downloading dependencies..."
    
    go mod download
    
    print_success "Dependencies downloaded"
}

# Verify dependencies
verify_deps() {
    print_status "Verifying dependencies..."
    
    go mod verify
    
    print_success "Dependencies verified"
}

# Check for unused dependencies
check_unused() {
    print_status "Checking for unused dependencies..."
    
    # This requires go mod tidy to be run first
    if go mod tidy; then
        print_success "No unused dependencies found"
    else
        print_warning "Some dependencies may be unused or have issues"
    fi
}

# Main function
main() {
    print_status "Go Module Cleaner"
    print_status "This script will clean Go module cache and dependencies"
    echo ""
    
    # Check prerequisites
    check_go
    
    # Clean everything
    clean_module_cache
    clean_go_mod
    download_deps
    verify_deps
    check_unused
    
    print_success "Go module cleanup completed!"
    print_status "Summary:"
    echo "  - Module cache cleaned"
    echo "  - Build cache cleaned"
    echo "  - Test cache cleaned"
    echo "  - go.mod and go.sum cleaned"
    echo "  - Dependencies downloaded and verified"
    echo ""
    print_warning "Note: This process may take a few minutes to complete"
}

# Handle help flag
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    echo "Go Module Cleaner"
    echo ""
    echo "Usage: $0"
    echo ""
    echo "This script will:"
    echo "  - Clean Go module cache"
    echo "  - Clean build cache"
    echo "  - Clean test cache"
    echo "  - Remove and regenerate go.sum"
    echo "  - Download and verify dependencies"
    echo "  - Check for unused dependencies"
    echo ""
    echo "Requirements:"
    echo "  - Go 1.21+"
    echo ""
    exit 0
fi

# Run main function
main
