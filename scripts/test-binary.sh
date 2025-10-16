#!/bin/bash

# Test script for binary functionality
# This script tests the help, version, and config flags

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

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Test help flag
test_help() {
    local binary=$1
    print_status "Testing help flag..."
    
    if timeout 5s "$binary" -help > /tmp/help_output.txt 2>&1; then
        if grep -q "Digiflazz Gateway API Server" /tmp/help_output.txt; then
            print_success "Help flag works correctly"
            return 0
        else
            print_error "Help flag output is incorrect"
            cat /tmp/help_output.txt
            return 1
        fi
    else
        print_error "Help flag test failed or timed out"
        cat /tmp/help_output.txt
        return 1
    fi
}

# Test version flag
test_version() {
    local binary=$1
    print_status "Testing version flag..."
    
    if timeout 5s "$binary" -version > /tmp/version_output.txt 2>&1; then
        if grep -q "Digiflazz Gateway API Server" /tmp/version_output.txt; then
            print_success "Version flag works correctly"
            return 0
        else
            print_error "Version flag output is incorrect"
            cat /tmp/version_output.txt
            return 1
        fi
    else
        print_error "Version flag test failed or timed out"
        cat /tmp/version_output.txt
        return 1
    fi
}

# Test config flag
test_config() {
    local binary=$1
    print_status "Testing config flag..."
    
    if timeout 5s "$binary" -config > /tmp/config_output.txt 2>&1; then
        if grep -q "Digiflazz Gateway API Server - Configuration" /tmp/config_output.txt; then
            print_success "Config flag works correctly"
            return 0
        else
            print_error "Config flag output is incorrect"
            cat /tmp/config_output.txt
            return 1
        fi
    else
        print_error "Config flag test failed or timed out"
        cat /tmp/config_output.txt
        return 1
    fi
}

# Test that server doesn't start with help flag
test_no_server_start() {
    local binary=$1
    print_status "Testing that server doesn't start with help flag..."
    
    # Start binary with help flag and check if it exits quickly
    local start_time=$(date +%s)
    if timeout 5s "$binary" -help > /tmp/no_server_output.txt 2>&1; then
        local end_time=$(date +%s)
        local duration=$((end_time - start_time))
        
        if [ $duration -lt 3 ]; then
            print_success "Server doesn't start with help flag (exited in ${duration}s)"
            return 0
        else
            print_error "Server may have started with help flag (took ${duration}s)"
            cat /tmp/no_server_output.txt
            return 1
        fi
    else
        print_error "Help flag test failed"
        cat /tmp/no_server_output.txt
        return 1
    fi
}

# Main test function
main() {
    local binary=${1:-"./gateway-digiflazz"}
    
    print_status "Testing binary: $binary"
    echo ""
    
    # Check if binary exists
    if [ ! -f "$binary" ]; then
        print_error "Binary not found: $binary"
        exit 1
    fi
    
    # Make binary executable if it's not
    chmod +x "$binary"
    
    local tests_passed=0
    local total_tests=4
    
    # Run tests
    if test_help "$binary"; then
        ((tests_passed++))
    fi
    echo ""
    
    if test_version "$binary"; then
        ((tests_passed++))
    fi
    echo ""
    
    if test_config "$binary"; then
        ((tests_passed++))
    fi
    echo ""
    
    if test_no_server_start "$binary"; then
        ((tests_passed++))
    fi
    echo ""
    
    # Summary
    print_status "Test Summary: $tests_passed/$total_tests tests passed"
    
    if [ $tests_passed -eq $total_tests ]; then
        print_success "All tests passed! Binary is working correctly."
        exit 0
    else
        print_error "Some tests failed. Please check the binary implementation."
        exit 1
    fi
}

# Handle help flag
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    echo "Binary test script for Digiflazz Gateway"
    echo ""
    echo "Usage: $0 [binary_path]"
    echo ""
    echo "Arguments:"
    echo "  binary_path    Path to the binary to test (default: ./gateway-digiflazz)"
    echo ""
    echo "Tests:"
    echo "  - Help flag functionality"
    echo "  - Version flag functionality"
    echo "  - Config flag functionality"
    echo "  - Server doesn't start with help flag"
    echo ""
    exit 0
fi

# Run main function
main "$@"
