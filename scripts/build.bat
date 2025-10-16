@echo off
REM Cross-platform build script for Digiflazz Gateway (Windows)
REM Usage: scripts\build.bat [platform] [arch] [version]

setlocal enabledelayedexpansion

REM Default values
set PLATFORM=%1
if "%PLATFORM%"=="" set PLATFORM=all
set ARCH=%2
if "%ARCH%"=="" set ARCH=all
set VERSION=%3
if "%VERSION%"=="" set VERSION=latest
set BUILD_DIR=build
set BINARY_NAME=gateway-digiflazz

echo [INFO] Digiflazz Gateway Cross-Platform Build
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

REM Check if GCC is installed (required for CGO)
gcc --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] GCC is not installed or not in PATH
    echo [INFO] Please install TDM-GCC or MinGW-w64 for CGO support
    echo [INFO] Download from: https://jmeubank.github.io/tdm-gcc/
    exit /b 1
)

REM Get Go version
for /f "tokens=3" %%i in ('go version') do set GO_VERSION=%%i
echo [INFO] Using Go version: %GO_VERSION%

REM Get GCC version
for /f "tokens=3" %%i in ('gcc --version 2^>^&1 ^| findstr "gcc"') do set GCC_VERSION=%%i
echo [INFO] Using GCC version: %GCC_VERSION%

REM Clean build directory
echo [INFO] Cleaning build directory...
if exist %BUILD_DIR% rmdir /s /q %BUILD_DIR%
mkdir %BUILD_DIR%

REM Build function
:build_platform
set platform=%1
set arch=%2
set output_name=%BINARY_NAME%-%platform%-%arch%

if "%platform%"=="windows" (
    set output_name=%output_name%.exe
)

echo [INFO] Building for %platform%/%arch%...

set GOOS=%platform%
set GOARCH=%arch%
set CGO_ENABLED=1
set CC=gcc
set CXX=g++

REM Build with optimized flags for Windows compatibility
go build -ldflags "-s -w -X main.version=%VERSION% -X main.buildTime=%date:~0,4%-%date:~5,2%-%date:~8,2%T%time:~0,2%:%time:~3,2%:%time:~6,2%Z" -o "%BUILD_DIR%\%output_name%" ./cmd/server

if errorlevel 1 (
    echo [WARNING] CGO build failed for %platform%/%arch%, trying without CGO...
    set CGO_ENABLED=0
    go build -ldflags "-s -w -X main.version=%VERSION% -X main.buildTime=%date:~0,4%-%date:~5,2%-%date:~8,2%T%time:~0,2%:%time:~3,2%:%time:~6,2%Z" -o "%BUILD_DIR%\%output_name%" ./cmd/server
    
    if errorlevel 1 (
        echo [ERROR] Failed to build for %platform%/%arch% (both CGO and non-CGO)
        exit /b 1
    ) else (
        echo [WARNING] Built without CGO - SQLite cache will not work
    )
)

echo [SUCCESS] Built %output_name%

REM Create archive
if "%platform%"=="windows" (
    powershell Compress-Archive -Path "%BUILD_DIR%\%output_name%" -DestinationPath "%BUILD_DIR%\%output_name%.zip"
) else (
    tar -czf "%BUILD_DIR%\%output_name%.tar.gz" -C "%BUILD_DIR%" "%output_name%"
)

echo [SUCCESS] Created archive for %output_name%
goto :eof

REM Build all platforms
:build_all
echo [INFO] Building for all platforms...

REM Linux
call :build_platform linux amd64
call :build_platform linux arm64
call :build_platform linux arm

REM Windows
call :build_platform windows amd64
call :build_platform windows arm64

REM macOS
call :build_platform darwin amd64
call :build_platform darwin arm64

REM FreeBSD
call :build_platform freebsd amd64
call :build_platform freebsd arm64

echo [SUCCESS] All builds completed!
goto :eof

REM Build specific platform
:build_specific
if "%ARCH%"=="all" (
    if "%PLATFORM%"=="linux" (
        call :build_platform linux amd64
        call :build_platform linux arm64
        call :build_platform linux arm
    ) else if "%PLATFORM%"=="windows" (
        call :build_platform windows amd64
        call :build_platform windows arm64
    ) else if "%PLATFORM%"=="darwin" (
        call :build_platform darwin amd64
        call :build_platform darwin arm64
    ) else if "%PLATFORM%"=="freebsd" (
        call :build_platform freebsd amd64
        call :build_platform freebsd arm64
    ) else (
        echo [ERROR] Unsupported platform: %PLATFORM%
        exit /b 1
    )
) else (
    call :build_platform %PLATFORM% %ARCH%
)
goto :eof

REM Show help
:show_help
echo Cross-platform build script for Digiflazz Gateway
echo.
echo Usage: %0 [platform] [arch] [version]
echo.
echo Arguments:
echo   platform    Target platform (linux, windows, darwin, freebsd, all)
echo   arch        Target architecture (amd64, arm64, arm, all)
echo   version     Version tag (default: latest)
echo.
echo Examples:
echo   %0                          # Build for all platforms
echo   %0 linux                    # Build for all Linux architectures
echo   %0 windows amd64           # Build for Windows x64
echo   %0 darwin arm64            # Build for macOS Apple Silicon
echo   %0 linux amd64 v1.0.0      # Build for Linux x64 with version
echo.
echo Supported platforms:
echo   - linux (amd64, arm64, arm)
echo   - windows (amd64, arm64)
echo   - darwin (amd64, arm64)
echo   - freebsd (amd64, arm64)
goto :eof

REM Handle help flag
if "%1"=="-h" goto show_help
if "%1"=="--help" goto show_help

REM Main execution
if "%PLATFORM%"=="all" (
    call :build_all
) else (
    call :build_specific
)

REM Show build summary
echo [INFO] Build Summary:
echo Build directory: %BUILD_DIR%
echo Files created:
dir %BUILD_DIR%

echo [SUCCESS] Build completed successfully!

