@echo off
echo ========================================
echo Debug Digiflazz Configuration on Windows
echo ========================================
echo.

echo Current Environment Variables:
echo DIGIFLAZZ_USERNAME=%DIGIFLAZZ_USERNAME%
echo DIGIFLAZZ_API_KEY=%DIGIFLAZZ_API_KEY%
echo DIGIFLAZZ_BASE_URL=%DIGIFLAZZ_BASE_URL%
echo DIGIFLAZZ_TIMEOUT=%DIGIFLAZZ_TIMEOUT%
echo DIGIFLAZZ_RETRY_ATTEMPTS=%DIGIFLAZZ_RETRY_ATTEMPTS%
echo DIGIFLAZZ_IP_WHITELIST=%DIGIFLAZZ_IP_WHITELIST%
echo.

echo Setting Default Values if Not Set:
if "%DIGIFLAZZ_TIMEOUT%"=="" set DIGIFLAZZ_TIMEOUT=30s
if "%DIGIFLAZZ_RETRY_ATTEMPTS%"=="" set DIGIFLAZZ_RETRY_ATTEMPTS=3
if "%DIGIFLAZZ_BASE_URL%"=="" set DIGIFLAZZ_BASE_URL=https://api.digiflazz.com
if "%DIGIFLAZZ_IP_WHITELIST%"=="" set DIGIFLAZZ_IP_WHITELIST=52.74.250.133

echo.
echo Updated Environment Variables:
echo DIGIFLAZZ_USERNAME=%DIGIFLAZZ_USERNAME%
echo DIGIFLAZZ_API_KEY=%DIGIFLAZZ_API_KEY%
echo DIGIFLAZZ_BASE_URL=%DIGIFLAZZ_BASE_URL%
echo DIGIFLAZZ_TIMEOUT=%DIGIFLAZZ_TIMEOUT%
echo DIGIFLAZZ_RETRY_ATTEMPTS=%DIGIFLAZZ_RETRY_ATTEMPTS%
echo DIGIFLAZZ_IP_WHITELIST=%DIGIFLAZZ_IP_WHITELIST%
echo.

echo Testing PLN Inquiry API...
echo Customer Number: 321058166
echo Ref ID: 12343
echo.

curl -X GET "http://localhost:8081/otomax/pln/inquiry?customer_no=321058166&ref_id=12343" ^
  -H "Content-Type: application/json" ^
  -v

echo.
echo ========================================
echo Debug Complete
echo ========================================
pause
