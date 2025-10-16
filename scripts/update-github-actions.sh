#!/bin/bash

# Update GitHub Actions to latest versions
# This script updates all deprecated actions to their latest versions

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

# Update action versions in workflow files
update_workflow_actions() {
    local workflow_file="$1"
    
    if [ ! -f "$workflow_file" ]; then
        print_warning "Workflow file not found: $workflow_file"
        return 1
    fi
    
    print_status "Updating actions in $workflow_file..."
    
    # Create backup
    cp "$workflow_file" "${workflow_file}.backup"
    
    # Update deprecated actions to latest versions
    sed -i '' 's/actions\/checkout@v3/actions\/checkout@v4/g' "$workflow_file"
    sed -i '' 's/actions\/checkout@v2/actions\/checkout@v4/g' "$workflow_file"
    sed -i '' 's/actions\/setup-go@v3/actions\/setup-go@v5/g' "$workflow_file"
    sed -i '' 's/actions\/setup-go@v4/actions\/setup-go@v5/g' "$workflow_file"
    sed -i '' 's/actions\/cache@v2/actions\/cache@v4/g' "$workflow_file"
    sed -i '' 's/actions\/cache@v3/actions\/cache@v4/g' "$workflow_file"
    sed -i '' 's/actions\/upload-artifact@v2/actions\/upload-artifact@v4/g' "$workflow_file"
    sed -i '' 's/actions\/upload-artifact@v3/actions\/upload-artifact@v4/g' "$workflow_file"
    sed -i '' 's/actions\/download-artifact@v2/actions\/download-artifact@v4/g' "$workflow_file"
    sed -i '' 's/actions\/download-artifact@v3/actions\/download-artifact@v4/g' "$workflow_file"
    sed -i '' 's/codecov\/codecov-action@v1/codecov\/codecov-action@v4/g' "$workflow_file"
    sed -i '' 's/codecov\/codecov-action@v2/codecov\/codecov-action@v4/g' "$workflow_file"
    sed -i '' 's/codecov\/codecov-action@v3/codecov\/codecov-action@v4/g' "$workflow_file"
    
    # Check if file was modified
    if ! cmp -s "$workflow_file" "${workflow_file}.backup"; then
        print_success "Updated actions in $workflow_file"
        rm "${workflow_file}.backup"
    else
        print_status "No updates needed for $workflow_file"
        rm "${workflow_file}.backup"
    fi
}

# Check for deprecated actions
check_deprecated_actions() {
    local workflow_file="$1"
    
    if [ ! -f "$workflow_file" ]; then
        return 1
    fi
    
    print_status "Checking for deprecated actions in $workflow_file..."
    
    # List of deprecated patterns
    local deprecated_patterns=(
        "actions/upload-artifact@v3"
        "actions/download-artifact@v3"
        "actions/cache@v3"
        "actions/setup-go@v4"
        "actions/checkout@v3"
        "codecov/codecov-action@v3"
    )
    
    local found_deprecated=false
    
    for pattern in "${deprecated_patterns[@]}"; do
        if grep -q "$pattern" "$workflow_file"; then
            print_warning "Found deprecated action: $pattern in $workflow_file"
            found_deprecated=true
        fi
    done
    
    if [ "$found_deprecated" = false ]; then
        print_success "No deprecated actions found in $workflow_file"
        return 0
    else
        return 1
    fi
}

# Main function
main() {
    print_status "GitHub Actions Version Updater"
    print_status "Checking and updating deprecated actions..."
    echo ""
    
    # Find all workflow files
    local workflow_files=($(find .github/workflows -name "*.yml" -o -name "*.yaml" 2>/dev/null))
    
    if [ ${#workflow_files[@]} -eq 0 ]; then
        print_warning "No workflow files found in .github/workflows/"
        exit 1
    fi
    
    print_status "Found ${#workflow_files[@]} workflow file(s):"
    for file in "${workflow_files[@]}"; do
        echo "  - $file"
    done
    echo ""
    
    # Check for deprecated actions
    local has_deprecated=false
    for workflow_file in "${workflow_files[@]}"; do
        if check_deprecated_actions "$workflow_file"; then
            has_deprecated=true
        fi
    done
    
    if [ "$has_deprecated" = false ]; then
        print_success "All actions are up to date!"
        exit 0
    fi
    
    echo ""
    print_status "Updating deprecated actions..."
    
    # Update workflow files
    for workflow_file in "${workflow_files[@]}"; do
        update_workflow_actions "$workflow_file"
    done
    
    echo ""
    print_success "GitHub Actions update completed!"
    print_status "Summary of updates:"
    echo "  - actions/checkout: v3/v2 → v4"
    echo "  - actions/setup-go: v4/v3 → v5"
    echo "  - actions/cache: v3/v2 → v4"
    echo "  - actions/upload-artifact: v3/v2 → v4"
    echo "  - actions/download-artifact: v3/v2 → v4"
    echo "  - codecov/codecov-action: v3/v2/v1 → v4"
    echo ""
    print_status "Next steps:"
    echo "1. Review the changes:"
    echo "   git diff .github/workflows/"
    echo ""
    echo "2. Commit the updates:"
    echo "   git add .github/workflows/"
    echo "   git commit -m 'Update GitHub Actions to latest versions'"
    echo "   git push origin main"
    echo ""
    print_warning "Note: Test the workflows after updating to ensure compatibility"
}

# Handle help flag
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    echo "GitHub Actions Version Updater"
    echo ""
    echo "Usage: $0"
    echo ""
    echo "This script will:"
    echo "  - Check all workflow files for deprecated actions"
    echo "  - Update deprecated actions to their latest versions"
    echo "  - Show summary of changes made"
    echo ""
    echo "Supported action updates:"
    echo "  - actions/checkout: v3/v2 → v4"
    echo "  - actions/setup-go: v4/v3 → v5"
    echo "  - actions/cache: v3/v2 → v4"
    echo "  - actions/upload-artifact: v3/v2 → v4"
    echo "  - actions/download-artifact: v3/v2 → v4"
    echo "  - codecov/codecov-action: v3/v2/v1 → v4"
    echo ""
    exit 0
fi

# Run main function
main
