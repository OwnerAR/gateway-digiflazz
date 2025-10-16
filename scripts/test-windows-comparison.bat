@echo off
REM Windows comparison test script
REM This script compares local development vs Windows binary behavior

setlocal enabledelayedexpansion

echo [INFO] Windows Comparison Test
echo [INFO] ========================
echo.

REM Set debug environment
set LOG_LEVEL=debug
set CACHE_DB_PATH=.\data\cache.db

REM Create data directory
if not exist "data" mkdir data

echo [INFO] Testing Windows Binary vs Local Development
echo [INFO] Environment:
echo   LOG_LEVEL=%LOG_LEVEL%
echo   CACHE_DB_PATH=%CACHE_DB_PATH%
echo   GOOS=%GOOS%
echo   GOARCH=%GOARCH%
echo.

echo [INFO] Step 1: Testing Windows Binary
echo [INFO] ================================
echo.

REM Start Windows binary in background
echo Starting Windows binary...
start /B build\gateway-digiflazz-windows-amd64.exe > windows-binary.log 2>&1

REM Wait for server to start
echo Waiting for server to start...
timeout /t 5 /nobreak > nul

REM Test PLN inquiry
echo Testing PLN inquiry with Windows binary...
curl -s "http://localhost:8080/otomax/pln/inquiry?customer_no=32105816634&ref_id=12343" > windows-response.json

REM Stop Windows binary
echo Stopping Windows binary...
taskkill /F /IM gateway-digiflazz-windows-amd64.exe > nul 2>&1

echo.
echo [INFO] Step 2: Testing Local Development
echo [INFO] ===================================
echo.

REM Test with local development (if Go is available)
where go >nul 2>&1
if errorlevel 1 (
    echo [WARNING] Go not found, skipping local development test
    goto :compare
)

echo Starting local development server...
start /B go run cmd/server/main.go > local-dev.log 2>&1

REM Wait for server to start
echo Waiting for server to start...
timeout /t 5 /nobreak > nul

REM Test PLN inquiry
echo Testing PLN inquiry with local development...
curl -s "http://localhost:8080/otomax/pln/inquiry?customer_no=32105816634&ref_id=12343" > local-response.json

REM Stop local development
echo Stopping local development...
taskkill /F /IM go.exe > nul 2>&1

:compare
echo.
echo [INFO] Step 3: Comparing Results
echo [INFO] ==========================
echo.

echo [INFO] Windows Binary Response:
type windows-response.json
echo.

echo [INFO] Local Development Response:
type local-response.json
echo.

echo [INFO] Log Files:
echo [INFO] Windows Binary Log:
type windows-binary.log
echo.

echo [INFO] Local Development Log:
type local-dev.log
echo.

echo [INFO] Analysis:
echo [INFO] ===========
echo.

REM Check if responses are different
fc windows-response.json local-response.json > nul
if errorlevel 1 (
    echo [DIFFERENCE] Responses are different!
    echo [INFO] Check the log files above for details
) else (
    echo [SAME] Responses are identical
)

echo.
echo [INFO] Test completed. Check the log files for detailed analysis.
echo [INFO] Files created:
echo   - windows-binary.log
echo   - local-dev.log
echo   - windows-response.json
echo   - local-response.json
