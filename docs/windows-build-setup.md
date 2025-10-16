# Windows Build Setup for Digiflazz Gateway

This document explains how to set up the build environment on Windows for the Digiflazz Gateway application, which requires CGO for SQLite support.

## Prerequisites

### 1. Go Installation
- Download and install Go from https://golang.org/dl/
- Make sure Go is added to your PATH
- Verify installation: `go version`

### 2. C Compiler Installation (Required for CGO)

Since the application uses SQLite which requires CGO, you need a C compiler. Choose one of the following options:

#### Option A: TDM-GCC (Recommended)
1. Download TDM-GCC from https://jmeubank.github.io/tdm-gcc/
2. Install with default settings
3. Restart your command prompt
4. Verify installation: `gcc --version`

#### Option B: MinGW-w64
1. Download MinGW-w64 from https://www.mingw-w64.org/
2. Install and add to PATH
3. Restart your command prompt
4. Verify installation: `gcc --version`

#### Option C: Microsoft Visual Studio Build Tools
1. Download Visual Studio Build Tools from https://visualstudio.microsoft.com/downloads/
2. Install with C++ build tools
3. Use Developer Command Prompt or set up environment variables
4. Verify installation: `cl`

## Building the Application

### Quick Build (Windows)
```bash
# Use the Windows-specific build script
scripts\build-windows.bat

# Or build for specific architecture
scripts\build-windows.bat windows amd64
```

### Cross-Platform Build
```bash
# Build for all platforms (requires CGO for SQLite)
scripts\build.bat

# Build for specific platform
scripts\build.bat windows amd64
```

### Using Make
```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Build for Windows specifically
make build-windows
```

## Troubleshooting

### CGO Compilation Errors

If you get errors like:
```
cgo: C compiler "gcc" not found: exec: "gcc": executable file not found in %PATH%
```

**Solution:**
1. Install TDM-GCC (easiest option)
2. Restart your command prompt
3. Verify GCC is in PATH: `gcc --version`

### SQLite Build Errors

If you get SQLite-related errors:
```
# github.com/mattn/go-sqlite3
cgo: C compiler "gcc" not found
```

**Solution:**
1. Ensure CGO is enabled: `set CGO_ENABLED=1`
2. Install a C compiler (TDM-GCC recommended)
3. Restart command prompt after installation

### Alternative: Docker Build

If you continue having CGO issues, you can use Docker:

```bash
# Build with Docker (no local C compiler needed)
docker build -t gateway-digiflazz .

# Run with Docker
docker run -p 8080:8080 gateway-digiflazz
```

## Environment Variables

The application can run without a `.env` file as it has default values:

```bash
# Default configuration (no .env file needed)
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
LOG_LEVEL=info
DIGIFLAZZ_USERNAME=your_username
DIGIFLAZZ_API_KEY=your_api_key
```

## Build Output

After successful build, you'll find:
- `build/gateway-digiflazz-windows-amd64.exe` - Windows executable
- `build/gateway-digiflazz-windows-amd64.zip` - Compressed archive

## Running the Application

```bash
# Run the built executable
build\gateway-digiflazz-windows-amd64.exe

# Or run directly with Go
go run ./cmd/server
```

## Notes

- The application uses SQLite for caching, which requires CGO
- CGO builds are platform-specific (can't cross-compile easily)
- For production deployment, consider using Docker
- The application includes default configuration values
- No `.env` file is required for basic operation


