@echo off
REM Test cache behavior for PLN inquiry
REM This script tests both cache miss and cache hit scenarios

setlocal enabledelayedexpansion

echo [INFO] Cache Behavior Test for PLN Inquiry
echo [INFO] ====================================
echo.

REM Set test parameters
set CUSTOMER_NO=543602392932
set REF_ID=12343
set BASE_URL=http://localhost:8080

echo [INFO] Test Parameters:
echo   Customer No: %CUSTOMER_NO%
echo   Ref ID: %REF_ID%
echo   Base URL: %BASE_URL%
echo.

REM Clear cache first
echo [INFO] Clearing cache...
del data\cache.db 2>nul
del cache.db 2>nul

echo [INFO] Test 1: Cache Miss (First Request)
echo [INFO] ===================================
echo.

curl -s "%BASE_URL%/otomax/pln/inquiry?customer_no=%CUSTOMER_NO%&ref_id=%REF_ID%" > response1.json

echo [INFO] Response 1 (Cache Miss):
type response1.json
echo.

echo [INFO] Test 2: Cache Hit (Second Request)
echo [INFO] ===================================
echo.

curl -s "%BASE_URL%/otomax/pln/inquiry?customer_no=%CUSTOMER_NO%&ref_id=%REF_ID%" > response2.json

echo [INFO] Response 2 (Cache Hit):
type response2.json
echo.

echo [INFO] Test 3: Different Ref ID (Cache Hit with Different Ref ID)
echo [INFO] ==========================================================
echo.

set REF_ID=56789
curl -s "%BASE_URL%/otomax/pln/inquiry?customer_no=%CUSTOMER_NO%&ref_id=%REF_ID%" > response3.json

echo [INFO] Response 3 (Cache Hit with Different Ref ID):
type response3.json
echo.

echo [INFO] Analysis:
echo [INFO] ===========
echo.

REM Check if responses are consistent
echo [INFO] Checking response consistency...
fc response1.json response2.json > nul
if errorlevel 1 (
    echo [DIFFERENCE] Response 1 and Response 2 are different
    echo [INFO] This indicates cache behavior issues
) else (
    echo [SAME] Response 1 and Response 2 are identical
    echo [INFO] Cache behavior is working correctly
)

echo.
echo [INFO] Test completed.
echo [INFO] Files created:
echo   - response1.json (Cache Miss)
echo   - response2.json (Cache Hit)
echo   - response3.json (Cache Hit with Different Ref ID)
echo.
echo [INFO] Expected behavior:
echo   - Response 1: Should have ref_id and message (Cache Miss)
echo   - Response 2: Should be identical to Response 1 (Cache Hit - Permanent Cache)
echo   - Response 3: Should have different ref_id but same customer data (Cache Hit)
echo.
echo [INFO] Cache Behavior:
echo   - PLN data is cached permanently (no TTL)
echo   - Customer data (name, meter_no, segment_power) never expires
echo   - Only ref_id changes per request
