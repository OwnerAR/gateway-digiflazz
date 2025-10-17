@echo off
echo ========================================
echo Test PLN Cache with Detailed Debug Logging
echo ========================================
echo.

REM Set debug logging
set LOG_LEVEL=debug

echo Starting application with debug logging...
echo.

REM Start the application in background
start /B gateway-digiflazz-windows-amd64.exe

REM Wait for application to start
timeout /t 3 /nobreak > nul

echo.
echo ========================================
echo Test 1: First PLN Inquiry (Should Cache Miss + Store)
echo ========================================
echo.

curl -X GET "http://localhost:8081/otomax/pln/inquiry?customer_no=32105816634&ref_id=1234" ^
  -H "Content-Type: application/json" ^
  -w "\n\nHTTP Status: %%{http_code}\nResponse Time: %%{time_total}s\n" ^
  -s

echo.
echo Waiting 2 seconds...
timeout /t 2 /nobreak > nul

echo.
echo ========================================
echo Test 2: Second PLN Inquiry (Should Cache Hit)
echo ========================================
echo.

curl -X GET "http://localhost:8081/otomax/pln/inquiry?customer_no=32105816634&ref_id=1234" ^
  -H "Content-Type: application/json" ^
  -w "\n\nHTTP Status: %%{http_code}\nResponse Time: %%{time_total}s\n" ^
  -s

echo.
echo ========================================
echo Cache Debug Test Complete
echo Check application logs for detailed cache information
echo ========================================
pause
