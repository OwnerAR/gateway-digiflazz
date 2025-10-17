@echo off
echo ========================================
echo Test PLN Inquiry with Proper Configuration
echo ========================================
echo.

REM Set proper environment variables
set DIGIFLAZZ_TIMEOUT=30s
set DIGIFLAZZ_RETRY_ATTEMPTS=3
set DIGIFLAZZ_BASE_URL=https://api.digiflazz.com
set DIGIFLAZZ_IP_WHITELIST=52.74.250.133

echo Environment Variables Set:
echo DIGIFLAZZ_TIMEOUT=%DIGIFLAZZ_TIMEOUT%
echo DIGIFLAZZ_RETRY_ATTEMPTS=%DIGIFLAZZ_RETRY_ATTEMPTS%
echo DIGIFLAZZ_BASE_URL=%DIGIFLAZZ_BASE_URL%
echo DIGIFLAZZ_IP_WHITELIST=%DIGIFLAZZ_IP_WHITELIST%
echo.

echo Starting Gateway Application...
echo.

REM Start the application with proper environment
start /B gateway-digiflazz-windows-amd64.exe

REM Wait for application to start
timeout /t 5 /nobreak > nul

echo Testing PLN Inquiry API...
echo.

REM Test PLN Inquiry
curl -X GET "http://localhost:8081/otomax/pln/inquiry?customer_no=321058166&ref_id=12343" ^
  -H "Content-Type: application/json" ^
  -w "\n\nHTTP Status: %%{http_code}\nResponse Time: %%{time_total}s\n" ^
  -s

echo.
echo Testing with different customer number...
curl -X GET "http://localhost:8081/otomax/pln/inquiry?customer_no=543602392932&ref_id=12344" ^
  -H "Content-Type: application/json" ^
  -w "\n\nHTTP Status: %%{http_code}\nResponse Time: %%{time_total}s\n" ^
  -s

echo.
echo ========================================
echo Test Complete
echo ========================================
pause
