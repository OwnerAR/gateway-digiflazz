# Windows Debug Script for PLN Inquiry Issues
# This PowerShell script provides detailed debugging information

param(
    [string]$CustomerNo = "32105816634",
    [string]$RefId = "12343",
    [switch]$Verbose
)

Write-Host "Windows Debug Script for PLN Inquiry" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""

# Function to log with timestamp
function Write-Log {
    param([string]$Message, [string]$Level = "INFO")
    $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    $color = switch ($Level) {
        "ERROR" { "Red" }
        "WARNING" { "Yellow" }
        "SUCCESS" { "Green" }
        default { "White" }
    }
    Write-Host "[$timestamp] [$Level] $Message" -ForegroundColor $color
}

# Function to test API endpoint
function Test-PLNAPI {
    param(
        [string]$BaseUrl,
        [string]$CustomerNo,
        [string]$RefId,
        [string]$Description
    )
    
    Write-Log "Testing $Description" "INFO"
    $url = "$BaseUrl/otomax/pln/inquiry?customer_no=$CustomerNo&ref_id=$RefId"
    Write-Log "URL: $url" "INFO"
    
    try {
        $response = Invoke-RestMethod -Uri $url -Method GET -TimeoutSec 30
        Write-Log "Response received" "SUCCESS"
        Write-Host "Response: $($response | ConvertTo-Json -Depth 3)" -ForegroundColor Green
        return $response
    }
    catch {
        Write-Log "Error: $($_.Exception.Message)" "ERROR"
        return $null
    }
}

# Function to check environment
function Test-Environment {
    Write-Log "Checking Environment" "INFO"
    
    # Check if binary exists
    if (Test-Path "build\gateway-digiflazz-windows-amd64.exe") {
        Write-Log "Windows binary found" "SUCCESS"
        $binaryInfo = Get-Item "build\gateway-digiflazz-windows-amd64.exe"
        Write-Log "Binary size: $($binaryInfo.Length) bytes" "INFO"
        Write-Log "Binary created: $($binaryInfo.CreationTime)" "INFO"
    } else {
        Write-Log "Windows binary not found" "ERROR"
        return $false
    }
    
    # Check cache database
    if (Test-Path "cache.db") {
        Write-Log "Cache database found" "SUCCESS"
        $cacheInfo = Get-Item "cache.db"
        Write-Log "Cache size: $($cacheInfo.Length) bytes" "INFO"
        Write-Log "Cache modified: $($cacheInfo.LastWriteTime)" "INFO"
    } else {
        Write-Log "Cache database not found" "WARNING"
    }
    
    # Check data directory
    if (Test-Path "data\cache.db") {
        Write-Log "Data cache database found" "SUCCESS"
        $dataCacheInfo = Get-Item "data\cache.db"
        Write-Log "Data cache size: $($dataCacheInfo.Length) bytes" "INFO"
    } else {
        Write-Log "Data cache database not found" "INFO"
    }
    
    # Check environment variables
    Write-Log "Environment Variables:" "INFO"
    $envVars = @("DIGIFLAZZ_USERNAME", "DIGIFLAZZ_API_KEY", "LOG_LEVEL", "CACHE_DB_PATH")
    foreach ($var in $envVars) {
        $value = [Environment]::GetEnvironmentVariable($var)
        if ($value) {
            if ($var -like "*API_KEY*") {
                Write-Log "$var = [HIDDEN]" "INFO"
            } else {
                Write-Log "$var = $value" "INFO"
            }
        } else {
            Write-Log "$var = [NOT SET]" "WARNING"
        }
    }
    
    return $true
}

# Function to start server and monitor
function Start-ServerMonitor {
    param([string]$Description, [string]$Command)
    
    Write-Log "Starting $Description" "INFO"
    Write-Log "Command: $Command" "INFO"
    
    # Start process
    $process = Start-Process -FilePath "cmd.exe" -ArgumentList "/c", $Command -PassThru -WindowStyle Hidden
    
    # Wait for server to start
    Write-Log "Waiting for server to start..." "INFO"
    Start-Sleep -Seconds 5
    
    # Check if process is running
    if (-not $process.HasExited) {
        Write-Log "Server started successfully" "SUCCESS"
        return $process
    } else {
        Write-Log "Server failed to start" "ERROR"
        return $null
    }
}

# Function to stop server
function Stop-Server {
    param([System.Diagnostics.Process]$Process)
    
    if ($Process -and -not $Process.HasExited) {
        Write-Log "Stopping server..." "INFO"
        $Process.Kill()
        $Process.WaitForExit(5000)
        Write-Log "Server stopped" "SUCCESS"
    }
}

# Main execution
Write-Log "Starting Windows Debug Analysis" "INFO"
Write-Host ""

# Check environment
if (-not (Test-Environment)) {
    Write-Log "Environment check failed" "ERROR"
    exit 1
}

Write-Host ""

# Test 1: Windows Binary
Write-Log "Test 1: Windows Binary" "INFO"
$winProcess = Start-ServerMonitor "Windows Binary" "build\gateway-digiflazz-windows-amd64.exe"

if ($winProcess) {
    $winResponse = Test-PLNAPI "http://localhost:8080" $CustomerNo $RefId "Windows Binary"
    Stop-Server $winProcess
    
    Start-Sleep -Seconds 2
    
    # Test 2: Local Development (if Go is available)
    Write-Host ""
    Write-Log "Test 2: Local Development" "INFO"
    
    $goProcess = Start-ServerMonitor "Local Development" "go run cmd/server/main.go"
    
    if ($goProcess) {
        $localResponse = Test-PLNAPI "http://localhost:8080" $CustomerNo $RefId "Local Development"
        Stop-Server $goProcess
    } else {
        Write-Log "Local development test skipped (Go not available)" "WARNING"
        $localResponse = $null
    }
    
    # Compare results
    Write-Host ""
    Write-Log "Comparison Results" "INFO"
    Write-Host "=================" -ForegroundColor Cyan
    
    if ($winResponse -and $localResponse) {
        $winJson = $winResponse | ConvertTo-Json -Depth 3
        $localJson = $localResponse | ConvertTo-Json -Depth 3
        
        if ($winJson -eq $localJson) {
            Write-Log "Responses are IDENTICAL" "SUCCESS"
        } else {
            Write-Log "Responses are DIFFERENT" "WARNING"
            Write-Host "Windows Response:" -ForegroundColor Yellow
            Write-Host $winJson -ForegroundColor Yellow
            Write-Host "Local Response:" -ForegroundColor Green
            Write-Host $localJson -ForegroundColor Green
        }
    } else {
        Write-Log "Cannot compare - one or both tests failed" "ERROR"
    }
}

Write-Host ""
Write-Log "Debug analysis completed" "SUCCESS"
