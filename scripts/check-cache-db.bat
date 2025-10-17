@echo off
echo ========================================
echo Check Cache Database Contents
echo ========================================
echo.

REM Check if cache database exists
if not exist "data\cache.db" (
    echo Cache database not found at data\cache.db
    echo.
    goto :check_alternative
)

echo Cache database found at data\cache.db
echo.

REM Use sqlite3 to check cache contents
echo ========================================
echo Cache Database Contents:
echo ========================================
echo.

sqlite3 data\cache.db "SELECT customer_no, substr(data, 1, 100) as data_preview, created_at, expires_at FROM pln_inquiry_cache ORDER BY created_at DESC;"

echo.
echo ========================================
echo Cache Statistics:
echo ========================================
echo.

sqlite3 data\cache.db "SELECT COUNT(*) as total_entries FROM pln_inquiry_cache;"
sqlite3 data\cache.db "SELECT COUNT(*) as permanent_entries FROM pln_inquiry_cache WHERE expires_at = '0001-01-01 00:00:00';"
sqlite3 data\cache.db "SELECT COUNT(*) as expired_entries FROM pln_inquiry_cache WHERE expires_at < datetime('now');"

goto :end

:check_alternative
echo Checking alternative cache locations...
echo.

if exist "cache.db" (
    echo Found cache.db in current directory
    sqlite3 cache.db "SELECT customer_no, substr(data, 1, 100) as data_preview, created_at, expires_at FROM pln_inquiry_cache ORDER BY created_at DESC;"
) else (
    echo No cache database found in current directory or data subdirectory
)

:end
echo.
echo ========================================
echo Cache Check Complete
echo ========================================
pause
