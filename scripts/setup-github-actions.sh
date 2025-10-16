#!/bin/bash

# Setup script for GitHub Actions workflows
# This script helps configure GitHub Actions for the project

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

# Check if we're in a git repository
check_git() {
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        print_error "Not in a git repository. Please initialize git first."
        exit 1
    fi
    print_success "Git repository detected"
}

# Check if GitHub remote exists
check_github_remote() {
    if ! git remote get-url origin | grep -q "github.com"; then
        print_warning "GitHub remote not detected. Please add GitHub remote:"
        echo "  git remote add origin https://github.com/username/gateway-digiflazz.git"
        echo "  git push -u origin main"
        return 1
    fi
    print_success "GitHub remote detected"
}

# Create .github directory structure
create_github_structure() {
    print_status "Creating GitHub Actions directory structure..."
    
    mkdir -p .github/workflows
    print_success "Created .github/workflows directory"
}

# Check if workflows exist
check_workflows() {
    local workflows=(
        ".github/workflows/build-windows.yml"
        ".github/workflows/build-all-platforms.yml"
        ".github/workflows/test.yml"
    )
    
    for workflow in "${workflows[@]}"; do
        if [ -f "$workflow" ]; then
            print_success "Found $workflow"
        else
            print_warning "Missing $workflow"
        fi
    done
}

# Create README for GitHub Actions
create_actions_readme() {
    cat > .github/README.md << 'EOF'
# GitHub Actions Workflows

This directory contains GitHub Actions workflows for the Digiflazz Gateway project.

## Workflows

### 1. `build-windows.yml`
- **Trigger**: Push to main/develop, PR to main, manual dispatch
- **Platform**: Windows (amd64, 386)
- **Features**: 
  - CGO enabled for SQLite support
  - Automatic artifact creation
  - Release upload on tags

### 2. `build-all-platforms.yml`
- **Trigger**: Push to main/develop, PR to main, manual dispatch
- **Platforms**: Linux, Windows, macOS
- **Architectures**: 
  - Linux: amd64, arm64, arm
  - Windows: amd64, 386
  - macOS: amd64, arm64
- **Features**:
  - Cross-platform builds
  - CGO enabled for all platforms
  - Automatic release creation

### 3. `test.yml`
- **Trigger**: Push to main/develop, PR to main
- **Features**:
  - Unit tests with coverage
  - Linting with golangci-lint
  - Security scanning with Gosec
  - Code quality checks

## Usage

### Manual Build
```bash
# Trigger manual build with custom version
gh workflow run build-windows.yml -f version=v1.0.0
```

### Release Process
1. Create and push a tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
2. GitHub Actions will automatically create a release with all platform binaries

### Local Testing
```bash
# Test the build process locally
make test
make lint
make build-all
```

## Artifacts

Each workflow creates artifacts:
- **Windows**: `.exe` and `.zip` files
- **Linux**: Binary and `.tar.gz` files  
- **macOS**: Binary and `.tar.gz` files

## Requirements

- Go 1.21+
- CGO enabled for SQLite support
- Windows SDK (automatically available in GitHub Actions Windows runners)

## Troubleshooting

### Build Failures
- Check Go version compatibility
- Verify CGO dependencies
- Review build logs for specific errors

### Missing Artifacts
- Ensure workflows completed successfully
- Check artifact retention settings
- Verify file paths in workflow steps
EOF

    print_success "Created GitHub Actions README"
}

# Main setup function
main() {
    print_status "Setting up GitHub Actions for Digiflazz Gateway"
    echo ""
    
    # Check prerequisites
    check_git
    
    # Create directory structure
    create_github_structure
    
    # Check existing workflows
    check_workflows
    
    # Create documentation
    create_actions_readme
    
    print_success "GitHub Actions setup completed!"
    echo ""
    print_status "Next steps:"
    echo "1. Commit and push the workflow files:"
    echo "   git add .github/"
    echo "   git commit -m 'Add GitHub Actions workflows'"
    echo "   git push origin main"
    echo ""
    echo "2. Check the Actions tab in your GitHub repository"
    echo "3. Trigger a manual build to test the workflows"
    echo ""
    print_warning "Note: Make sure your repository is public or you have GitHub Pro for private repositories"
}

# Handle help flag
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    echo "GitHub Actions setup script for Digiflazz Gateway"
    echo ""
    echo "Usage: $0"
    echo ""
    echo "This script will:"
    echo "  - Check git repository"
    echo "  - Create .github/workflows directory"
    echo "  - Verify workflow files exist"
    echo "  - Create documentation"
    echo ""
    echo "Requirements:"
    echo "  - Git repository"
    echo "  - GitHub remote configured"
    echo ""
    exit 0
fi

# Run main function
main
