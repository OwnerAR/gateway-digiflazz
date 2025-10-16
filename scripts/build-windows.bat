@echo off
REM Windows build script for Digiflazz Gateway with CGO support
REM This script handles CGO requirements for SQLite

setlocal enabledelayedexpansion

REM Default values
set PLATFORM=%1
if "%PLATFORM%"=="" set PLATFORM=windows
set ARCH=%2
if "%ARCH%"=="" set ARCH=amd64
set VERSION=%3
if "%VERSION%"=="" set VERSION=latest
set BUILD_DIR=build
set BINARY_NAME=gateway-digiflazz

echo [INFO] Digiflazz Gateway Windows Build with CGO
echo [INFO] Platform: %PLATFORM%
echo [INFO] Architecture: %ARCH%
echo [INFO] Version: %VERSION%
echo [INFO] Build Directory: %BUILD_DIR%
echo.

REM Check if Go is installed
go version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Go is not installed or not in PATH
    exit /b 1
)

REM Get Go version
for /f "tokens=3" %%i in ('go version') do set GO_VERSION=%%i
echo [INFO] Using Go version: %GO_VERSION%

REM Check for C compiler
where gcc >nul 2>&1
if errorlevel 1 (
    echo [WARNING] GCC not found. Checking for alternative C compilers...
    where cl >nul 2>&1
    if errorlevel 1 (
        echo [ERROR] No C compiler found. Please install one of the following:
        echo   - TDM-GCC: https://jmeubank.github.io/tdm-gcc/
        echo   - MinGW-w64: https://www.mingw-w64.org/
        echo   - Microsoft Visual Studio Build Tools
        echo.
        echo For quick setup, download and install TDM-GCC.
        exit /b 1
    ) else (
        echo [INFO] Found Microsoft Visual C++ compiler
        set CC=cl
    )
) else (
    echo [INFO] Found GCC compiler
    set CC=gcc
)

REM Clean build directory
echo [INFO] Cleaning build directory...
if exist %BUILD_DIR% rmdir /s /q %BUILD_DIR%
mkdir %BUILD_DIR%

REM Set environment variables for CGO
set CGO_ENABLED=1
set GOOS=%PLATFORM%
set GOARCH=%ARCH%

REM Set output name
set output_name=%BINARY_NAME%-%PLATFORM%-%ARCH%.exe

echo [INFO] Building for %PLATFORM%/%ARCH% with CGO enabled...

REM Build with CGO
go build -ldflags "-X main.version=%VERSION% -X main.buildTime=%date:~0,4%-%date:~5,2%-%date:~8,2%T%time:~0,2%:%time:~3,2%:%time:~6,2%Z" -o "%BUILD_DIR%\%output_name%" ./cmd/server

if errorlevel 1 (
    echo [ERROR] Failed to build for %PLATFORM%/%ARCH%
    echo.
    echo [TROUBLESHOOTING] If you're getting CGO errors:
    echo 1. Install TDM-GCC: https://jmeubank.github.io/tdm-gcc/
    echo 2. Or install MinGW-w64: https://www.mingw-w64.org/
    echo 3. Make sure the C compiler is in your PATH
    echo 4. Restart your command prompt after installation
    exit /b 1
)

echo [SUCCESS] Built %output_name%

REM Create archive
powershell Compress-Archive -Path "%BUILD_DIR%\%output_name%" -DestinationPath "%BUILD_DIR%\%output_name%.zip"

echo [SUCCESS] Created archive for %output_name%

REM Show build summary
echo [INFO] Build Summary:
echo Build directory: %BUILD_DIR%
echo Files created:
dir %BUILD_DIR%

echo [SUCCESS] Windows build completed successfully!
echo.
echo [NOTE] The binary requires CGO runtime libraries.
echo For distribution, consider using Docker or static linking alternatives.


