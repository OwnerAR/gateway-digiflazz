# Windows Configuration Fix

## Problem Identified

From the Windows log analysis, the main issues causing empty PLN inquiry responses are:

### 1. **Timeout Configuration Issue**
```
"retry_attempts=0 timeout=0s"
```
**Problem**: No timeout and no retry attempts configured, causing requests to fail silently.

### 2. **User-Agent Issue**
```
"user_agent="<nil>"
```
**Problem**: User-Agent header not properly set in HTTP requests.

### 3. **Empty Response from Digiflazz**
```
"message= meter_no= name= rc= ref_id= status="
"response="{{        }  0}"
```
**Problem**: Digiflazz API returns empty response due to configuration issues.

## Solutions Applied

### 1. **Fixed Default Configuration**
```go
// Set default timeout and retry attempts if not configured
if cfg.Digiflazz.Timeout == 0 {
    cfg.Digiflazz.Timeout = 30 * time.Second
}
if cfg.Digiflazz.RetryAttempts == 0 {
    cfg.Digiflazz.RetryAttempts = 3
}
```

### 2. **Added Environment Variable Support**
```bash
# New environment variables
DIGIFLAZZ_TIMEOUT=30s
DIGIFLAZZ_RETRY_ATTEMPTS=3
```

### 3. **Fixed User-Agent Header**
```go
userAgent := fmt.Sprintf("Digiflazz-Gateway/1.0 (%s/%s)", runtime.GOOS, runtime.GOARCH)
httpReq.Header.Set("User-Agent", userAgent)
```

### 4. **Enhanced Logging**
- Platform information logging
- User-Agent header logging
- Raw response logging for debugging

## Testing Commands

### 1. **Test with Fixed Configuration**
```bash
# Build and test with fixed config
make build-windows-cgo
make test-windows-fixed

# Or manually
scripts/test-windows-fixed.bat
```

### 2. **Test PLN API Directly**
```bash
# Test the API endpoint
make test-pln-api

# Or manually
scripts/test-pln-api.bat
```

### 3. **Manual Testing**
```bash
# Set environment variables
set DIGIFLAZZ_TIMEOUT=30s
set DIGIFLAZZ_RETRY_ATTEMPTS=3
set LOG_LEVEL=debug

# Start server
build\gateway-digiflazz-windows-amd64.exe

# Test in another terminal
curl "http://localhost:8080/otomax/pln/inquiry?customer_no=321058166634&ref_id=12343"
```

## Expected Log Output (Fixed)

After applying the fixes, you should see:

```
"timeout": 30000000000
"retry_attempts": 3
"user_agent": "Digiflazz-Gateway/1.0 (windows/amd64)"
"platform": "windows/amd64"
```

And successful PLN inquiry response:
```
"rc": "00"
"status": "Sukses"
"message": "Transaksi Sukses"
```

## Troubleshooting Steps

### 1. **Check Configuration**
```bash
# Verify environment variables are set
echo %DIGIFLAZZ_TIMEOUT%
echo %DIGIFLAZZ_RETRY_ATTEMPTS%
echo %DIGIFLAZZ_USERNAME%
echo %DIGIFLAZZ_API_KEY%
```

### 2. **Test with Debug Logging**
```bash
set LOG_LEVEL=debug
build\gateway-digiflazz-windows-amd64.exe
```

### 3. **Compare with Local Development**
```bash
# Test local development
go run cmd/server/main.go

# Test Windows binary
scripts/test-windows-fixed.bat
```

### 4. **Check Network Connectivity**
```bash
# Test Digiflazz API directly
curl -X POST "https://api.digiflazz.com/v1/inquiry-pln" \
  -H "Content-Type: application/json" \
  -d '{"username":"your_username","customer_no":"321058166634","sign":"your_signature"}'
```

## Common Issues and Solutions

### 1. **Still Getting Empty Response**
- Check if IP is whitelisted in Digiflazz
- Verify username and API key are correct
- Check network connectivity

### 2. **Timeout Errors**
- Increase timeout: `set DIGIFLAZZ_TIMEOUT=60s`
- Check network stability

### 3. **Retry Failures**
- Increase retry attempts: `set DIGIFLAZZ_RETRY_ATTEMPTS=5`
- Check Digiflazz API status

### 4. **Cache Issues**
- Clear cache: `del data\cache.db`
- Check cache permissions

## Environment Variables Reference

```bash
# Required
DIGIFLAZZ_USERNAME=your_username
DIGIFLAZZ_API_KEY=your_api_key

# Optional (with defaults)
DIGIFLAZZ_BASE_URL=https://api.digiflazz.com
DIGIFLAZZ_TIMEOUT=30s
DIGIFLAZZ_RETRY_ATTEMPTS=3
LOG_LEVEL=info
CACHE_DB_PATH=.\data\cache.db
```

## File Structure

```
gateway-digiflazz/
├── build/
│   └── gateway-digiflazz-windows-amd64.exe
├── data/
│   └── cache.db
├── scripts/
│   ├── test-windows-fixed.bat
│   ├── test-pln-api.bat
│   └── debug-windows.ps1
└── .env (optional)
```

## Next Steps

1. **Rebuild the Windows binary** with the fixes
2. **Test with the fixed configuration**
3. **Compare results** with local development
4. **Monitor logs** for any remaining issues

The main issue was the missing timeout and retry configuration, which should now be resolved.
