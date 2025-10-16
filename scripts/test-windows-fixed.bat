@echo off
REM Windows test script with fixed configuration
REM This script sets proper timeout and retry configuration

setlocal enabledelayedexpansion

echo [INFO] Windows Test with Fixed Configuration
echo [INFO] ======================================
echo.

REM Set proper environment variables
set LOG_LEVEL=debug
set DIGIFLAZZ_TIMEOUT=30s
set DIGIFLAZZ_RETRY_ATTEMPTS=3
set CACHE_DB_PATH=.\data\cache.db

REM Create data directory
if not exist "data" mkdir data

echo [INFO] Environment Configuration:
echo   LOG_LEVEL=%LOG_LEVEL%
echo   DIGIFLAZZ_TIMEOUT=%DIGIFLAZZ_TIMEOUT%
echo   DIGIFLAZZ_RETRY_ATTEMPTS=%DIGIFLAZZ_RETRY_ATTEMPTS%
echo   CACHE_DB_PATH=%CACHE_DB_PATH%
echo.

REM Check if binary exists
if not exist "build\gateway-digiflazz-windows-amd64.exe" (
    echo [ERROR] Binary not found: build\gateway-digiflazz-windows-amd64.exe
    echo [INFO] Please build first: make build-windows-cgo
    exit /b 1
)

echo [INFO] Starting Windows binary with fixed configuration...
echo [INFO] Server will start on http://localhost:8080
echo [INFO] Press Ctrl+C to stop
echo.

REM Start the binary
build\gateway-digiflazz-windows-amd64.exe
