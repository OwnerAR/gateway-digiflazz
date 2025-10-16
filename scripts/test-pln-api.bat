@echo off
REM Test PLN API directly
REM This script tests the PLN inquiry API endpoint

setlocal enabledelayedexpansion

echo [INFO] PLN API Test Script
echo [INFO] ===================
echo.

REM Set test parameters
set CUSTOMER_NO=321058166634
set REF_ID=12343
set BASE_URL=http://localhost:8080

echo [INFO] Test Parameters:
echo   Customer No: %CUSTOMER_NO%
echo   Ref ID: %REF_ID%
echo   Base URL: %BASE_URL%
echo.

echo [INFO] Testing PLN Inquiry API...
echo.

REM Test the API endpoint
curl -v "%BASE_URL%/otomax/pln/inquiry?customer_no=%CUSTOMER_NO%&ref_id=%REF_ID%"

echo.
echo [INFO] Test completed.
echo [INFO] If you see empty response, check the server logs for details.
