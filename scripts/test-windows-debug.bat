@echo off
REM Windows debug test script for PLN inquiry
REM This script helps debug issues with Windows binary

setlocal enabledelayedexpansion

echo [INFO] Windows Debug Test for PLN Inquiry
echo [INFO] =====================================
echo.

REM Set environment variables for debugging
set LOG_LEVEL=debug
set CACHE_DB_PATH=.\data\cache.db

REM Create data directory if it doesn't exist
if not exist "data" mkdir data

echo [INFO] Environment Variables:
echo   LOG_LEVEL=%LOG_LEVEL%
echo   CACHE_DB_PATH=%CACHE_DB_PATH%
echo.

REM Check if binary exists
if not exist "build\gateway-digiflazz-windows-amd64.exe" (
    echo [ERROR] Binary not found: build\gateway-digiflazz-windows-amd64.exe
    echo [INFO] Please build the application first using:
    echo   make build-windows-cgo
    echo   OR
    echo   scripts\build-windows.bat
    exit /b 1
)

echo [INFO] Starting application with debug logging...
echo [INFO] Press Ctrl+C to stop the server
echo.

REM Start the application
build\gateway-digiflazz-windows-amd64.exe

echo.
echo [INFO] Application stopped
