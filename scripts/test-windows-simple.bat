@echo off
REM Simple Windows test script
REM This script tests the Windows binary and shows detailed output

setlocal enabledelayedexpansion

echo [INFO] Simple Windows Test for PLN Inquiry
echo [INFO] ====================================
echo.

REM Set debug logging
set LOG_LEVEL=debug

REM Check if binary exists
if not exist "build\gateway-digiflazz-windows-amd64.exe" (
    echo [ERROR] Binary not found: build\gateway-digiflazz-windows-amd64.exe
    echo [INFO] Please build first: make build-windows-cgo
    exit /b 1
)

REM Create data directory
if not exist "data" mkdir data

echo [INFO] Starting Windows binary with debug logging...
echo [INFO] Server will start on http://localhost:8080
echo [INFO] Press Ctrl+C to stop
echo.

REM Start the binary
build\gateway-digiflazz-windows-amd64.exe
