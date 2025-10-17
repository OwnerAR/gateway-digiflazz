@echo off
echo ========================================
echo Test PLN Cache with Debug Logging
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
echo Test 1: First PLN Inquiry (Should Cache Miss)
echo ========================================
echo.

curl -X GET "http://localhost:8081/otomax/pln/inquiry?customer_no=32105816634&ref_id=1234" ^
  -H "Content-Type: application/json" ^
  -w "\n\nHTTP Status: %%{http_code}\nResponse Time: %%{time_total}s\n" ^
  -s

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
echo Test 3: Different Customer (Should Cache Miss)
echo ========================================
echo.

curl -X GET "http://localhost:8081/otomax/pln/inquiry?customer_no=543602392932&ref_id=1235" ^
  -H "Content-Type: application/json" ^
  -w "\n\nHTTP Status: %%{http_code}\nResponse Time: %%{time_total}s\n" ^
  -s

echo.
echo ========================================
echo Test 4: Same Customer Again (Should Cache Hit)
echo ========================================
echo.

curl -X GET "http://localhost:8081/otomax/pln/inquiry?customer_no=32105816634&ref_id=1236" ^
  -H "Content-Type: application/json" ^
  -w "\n\nHTTP Status: %%{http_code}\nResponse Time: %%{time_total}s\n" ^
  -s

echo.
echo ========================================
echo Cache Test Complete
echo Check application logs for cache debug information
echo ========================================
pause
