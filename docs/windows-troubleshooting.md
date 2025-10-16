# Windows Troubleshooting Guide

## PLN Inquiry Empty Response Issue

### Problem Description
When running the Windows binary, PLN inquiry returns empty response (`message=""`, `status=""`, `rc=""`) even though the customer number is correct and the same request works fine when running the application directly (local development).

### Root Causes

#### 1. **Cache Database Path Issues**
- Windows may have different working directory behavior
- SQLite database file might not be created in the expected location
- File permissions might be different on Windows

#### 2. **Environment Variables**
- Windows handles environment variables differently
- Default values might not be loaded properly
- Configuration file paths might be different

#### 3. **HTTP Client Differences**
- Windows might have different network stack behavior
- User-Agent headers might be different
- SSL/TLS handling might vary

### Debugging Steps

#### 1. **Enable Debug Logging**
```bash
# Set debug level
set LOG_LEVEL=debug

# Run with debug
scripts\test-windows-debug.bat
```

#### 2. **Check Cache Database**
```bash
# Check if cache database exists
dir data\cache.db

# Check cache directory permissions
dir data\
```

#### 3. **Verify Environment Variables**
```bash
# Check current environment
echo %DIGIFLAZZ_USERNAME%
echo %DIGIFLAZZ_API_KEY%
echo %CACHE_DB_PATH%
```

#### 4. **Test with Different Cache Path**
```bash
# Set custom cache path
set CACHE_DB_PATH=.\custom-cache.db

# Run application
build\gateway-digiflazz-windows-amd64.exe
```

### Solutions

#### 1. **Fixed Cache Path Handling**
The application now automatically:
- Creates `data/` directory if it doesn't exist
- Uses proper Windows path separators
- Falls back to current directory if needed

#### 2. **Enhanced Logging**
Added detailed logging for:
- Cache initialization
- HTTP client configuration
- Raw API responses
- Platform information

#### 3. **Environment Variable Support**
- `CACHE_DB_PATH`: Custom cache database path
- `LOG_LEVEL`: Set to `debug` for detailed logging

### Testing Commands

#### 1. **Basic Test**
```bash
# Start server
build\gateway-digiflazz-windows-amd64.exe

# Test PLN inquiry
curl "http://localhost:8080/otomax/pln/inquiry?customer_no=32105816634&ref_id=12343"
```

#### 2. **Debug Test**
```bash
# Run with debug logging
scripts\test-windows-debug.bat

# In another terminal, test the API
curl "http://localhost:8080/otomax/pln/inquiry?customer_no=32105816634&ref_id=12343"
```

#### 3. **Cache Test**
```bash
# Clear cache and test
del data\cache.db
build\gateway-digiflazz-windows-amd64.exe
```

### Expected Log Output (Debug Mode)

When working correctly, you should see:
```
{"level":"info","msg":"Initializing SQLite cache","cache_path":"data/cache.db"}
{"level":"debug","msg":"Raw response from Digiflazz API","endpoint":"/inquiry-pln","raw_response":"{\"data\":{...}}"}
{"level":"info","msg":"PLN inquiry response received from Digiflazz API","rc":"00","status":"Sukses"}
```

When failing, you might see:
```
{"level":"error","msg":"Failed to initialize SQLite cache","cache_path":"data/cache.db"}
{"level":"error","msg":"Failed to unmarshal Digiflazz API response"}
{"level":"warn","msg":"PLN inquiry returned empty response"}
```

### Common Issues and Fixes

#### 1. **Cache Database Permission Error**
```bash
# Run as Administrator or change permissions
icacls data /grant Everyone:F
```

#### 2. **Environment Variables Not Set**
```bash
# Set required variables
set DIGIFLAZZ_USERNAME=your_username
set DIGIFLAZZ_API_KEY=your_api_key
```

#### 3. **Working Directory Issues**
```bash
# Run from application directory
cd C:\path\to\gateway-digiflazz
build\gateway-digiflazz-windows-amd64.exe
```

### File Structure (Windows)
```
gateway-digiflazz/
├── build/
│   └── gateway-digiflazz-windows-amd64.exe
├── data/
│   └── cache.db (created automatically)
├── scripts/
│   └── test-windows-debug.bat
└── .env (optional)
```

### Performance Notes
- SQLite on Windows might be slower than on Unix systems
- Consider using Redis cache for better Windows performance
- Network timeouts might need adjustment for Windows environments

### Support
If issues persist:
1. Run with `LOG_LEVEL=debug`
2. Check Windows Event Viewer for system errors
3. Verify antivirus software isn't blocking the application
4. Test with different customer numbers
5. Compare with local development environment
