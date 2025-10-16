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
