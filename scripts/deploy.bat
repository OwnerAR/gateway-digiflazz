@echo off
REM Cross-platform deployment script for Digiflazz Gateway (Windows)
REM Usage: scripts\deploy.bat [platform] [environment] [version]

setlocal enabledelayedexpansion

REM Default values
set PLATFORM=%1
if "%PLATFORM%"=="" set PLATFORM=auto
set ENVIRONMENT=%2
if "%ENVIRONMENT%"=="" set ENVIRONMENT=production
set VERSION=%3
if "%VERSION%"=="" set VERSION=latest
set SERVICE_NAME=digiflazz-gateway

echo [INFO] Digiflazz Gateway Cross-Platform Deployment
echo [INFO] Platform: %PLATFORM%
echo [INFO] Environment: %ENVIRONMENT%
echo [INFO] Version: %VERSION%
echo.

REM Detect platform
if "%PLATFORM%"=="auto" (
    set PLATFORM=windows
    set ARCH=amd64
)

echo [INFO] Detected platform: %PLATFORM%/%ARCH%

REM Check prerequisites
echo [INFO] Checking prerequisites...

set BINARY_PATH=build\gateway-digiflazz-%PLATFORM%-%ARCH%.exe
if not exist "%BINARY_PATH%" (
    echo [ERROR] Binary not found: %BINARY_PATH%
    echo [INFO] Please run 'make build-all' first
    exit /b 1
)

echo [SUCCESS] Binary found: %BINARY_PATH%

REM Deploy to Windows
if "%PLATFORM%"=="windows" (
    echo [INFO] Deploying to Windows...
    
    REM Create Windows service using NSSM
    if not exist "C:\Program Files\nssm" (
        echo [ERROR] NSSM is not installed. Please install NSSM first.
        echo [INFO] Download from: https://nssm.cc/download
        exit /b 1
    )
    
    REM Stop existing service
    "C:\Program Files\nssm\nssm.exe" stop %SERVICE_NAME% 2>nul
    "C:\Program Files\nssm\nssm.exe" remove %SERVICE_NAME% confirm 2>nul
    
    REM Create service directory
    if not exist "C:\%SERVICE_NAME%" mkdir "C:\%SERVICE_NAME%"
    if not exist "C:\%SERVICE_NAME%\logs" mkdir "C:\%SERVICE_NAME%\logs"
    if not exist "C:\%SERVICE_NAME%\cache" mkdir "C:\%SERVICE_NAME%\cache"
    if not exist "C:\%SERVICE_NAME%\config" mkdir "C:\%SERVICE_NAME%\config"
    
    REM Copy binary and config
    copy "%BINARY_PATH%" "C:\%SERVICE_NAME%\"
    copy "configs\config.yaml" "C:\%SERVICE_NAME%\config\"
    copy "configs\.env.example" "C:\%SERVICE_NAME%\config\.env"
    
    REM Install service
    "C:\Program Files\nssm\nssm.exe" install %SERVICE_NAME% "C:\%SERVICE_NAME%\gateway-digiflazz-%PLATFORM%-%ARCH%.exe"
    "C:\Program Files\nssm\nssm.exe" set %SERVICE_NAME% AppDirectory "C:\%SERVICE_NAME%"
    "C:\Program Files\nssm\nssm.exe" set %SERVICE_NAME% AppStdout "C:\%SERVICE_NAME%\logs\output.log"
    "C:\Program Files\nssm\nssm.exe" set %SERVICE_NAME% AppStderr "C:\%SERVICE_NAME%\logs\error.log"
    "C:\Program Files\nssm\nssm.exe" set %SERVICE_NAME% AppEnvironmentExtra "LOG_LEVEL=info" "SERVER_PORT=8080"
    "C:\Program Files\nssm\nssm.exe" set %SERVICE_NAME% Start SERVICE_AUTO_START
    
    REM Start service
    "C:\Program Files\nssm\nssm.exe" start %SERVICE_NAME%
    
    echo [SUCCESS] Service deployed and started
    echo [INFO] Check status: sc query %SERVICE_NAME%
    echo [INFO] View logs: type "C:\%SERVICE_NAME%\logs\output.log"
)

REM Deploy using Docker
if "%PLATFORM%"=="docker" (
    echo [INFO] Deploying using Docker...
    
    REM Check if Docker is installed
    docker --version >nul 2>&1
    if errorlevel 1 (
        echo [ERROR] Docker is not installed
        exit /b 1
    )
    
    REM Build Docker image
    docker build -t %SERVICE_NAME%:%VERSION% .
    
    REM Stop existing container
    docker stop %SERVICE_NAME% 2>nul
    docker rm %SERVICE_NAME% 2>nul
    
    REM Run new container
    docker run -d ^
        --name %SERVICE_NAME% ^
        --restart unless-stopped ^
        -p 8080:8080 ^
        -v "%cd%\logs:/app/logs" ^
        -v "%cd%\cache:/app/cache" ^
        -e DIGIFLAZZ_USERNAME=%DIGIFLAZZ_USERNAME% ^
        -e DIGIFLAZZ_API_KEY=%DIGIFLAZZ_API_KEY% ^
        -e SERVER_PORT=8080 ^
        -e LOG_LEVEL=info ^
        %SERVICE_NAME%:%VERSION%
    
    echo [SUCCESS] Docker container deployed and started
    echo [INFO] Check status: docker ps ^| findstr %SERVICE_NAME%
    echo [INFO] View logs: docker logs -f %SERVICE_NAME%
)

REM Deploy using Docker Compose
if "%PLATFORM%"=="compose" (
    echo [INFO] Deploying using Docker Compose...
    
    REM Check if Docker Compose is installed
    docker-compose --version >nul 2>&1
    if errorlevel 1 (
        echo [ERROR] Docker Compose is not installed
        exit /b 1
    )
    
    REM Stop existing services
    docker-compose down 2>nul
    
    REM Start services
    docker-compose up -d
    
    echo [SUCCESS] Docker Compose services deployed and started
    echo [INFO] Check status: docker-compose ps
    echo [INFO] View logs: docker-compose logs -f
)

echo [SUCCESS] Deployment completed successfully!



