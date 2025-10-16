# Cross-Platform Build and Deployment Guide

## Overview
This guide covers building and deploying the Digiflazz Gateway application across multiple platforms and architectures.

## Supported Platforms

### Operating Systems
- **Linux**: amd64, arm64, arm
- **Windows**: amd64, arm64
- **macOS**: amd64, arm64
- **FreeBSD**: amd64, arm64

### Deployment Methods
- **Native**: Direct binary deployment
- **Docker**: Containerized deployment
- **Docker Compose**: Multi-service deployment

## Build Commands

### Using Makefile
```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Build for specific platform
make build-linux
make build-windows
make build-darwin
make build-freebsd

# Build for specific platform and architecture
make build-platform PLATFORM=linux ARCH=amd64

# Cross-compile for specific target
make cross-compile GOOS=windows GOARCH=amd64

# Create release packages
make release VERSION=v1.0.0
```

### Using Build Scripts
```bash
# Linux/macOS
./scripts/build.sh                    # Build for all platforms
./scripts/build.sh linux              # Build for all Linux architectures
./scripts/build.sh windows amd64     # Build for Windows x64
./scripts/build.sh darwin arm64      # Build for macOS Apple Silicon
./scripts/build.sh linux amd64 v1.0.0 # Build with version

# Windows
scripts\build.bat                     # Build for all platforms
scripts\build.bat linux              # Build for all Linux architectures
scripts\build.bat windows amd64      # Build for Windows x64
scripts\build.bat darwin arm64       # Build for macOS Apple Silicon
```

## Deployment Commands

### Using Deployment Scripts
```bash
# Linux/macOS
./scripts/deploy.sh                   # Auto-detect platform and deploy
./scripts/deploy.sh linux             # Deploy to Linux
./scripts/deploy.sh docker            # Deploy using Docker
./scripts/deploy.sh compose           # Deploy using Docker Compose

# Windows
scripts\deploy.bat                    # Auto-detect platform and deploy
scripts\deploy.bat windows            # Deploy to Windows
scripts\deploy.bat docker             # Deploy using Docker
scripts\deploy.bat compose            # Deploy using Docker Compose
```

## Platform-Specific Deployment

### Linux Deployment
```bash
# Build for Linux
make build-linux

# Deploy natively
./scripts/deploy.sh linux

# Deploy using Docker
./scripts/deploy.sh docker

# Deploy using Docker Compose
./scripts/deploy.sh compose
```

**Native Linux Deployment:**
- Creates systemd service
- Installs to `/opt/digiflazz-gateway/`
- Logs to `/var/log/digiflazz-gateway/`
- Config in `/etc/digiflazz-gateway/`

**Service Management:**
```bash
# Check status
sudo systemctl status digiflazz-gateway

# Start/stop/restart
sudo systemctl start digiflazz-gateway
sudo systemctl stop digiflazz-gateway
sudo systemctl restart digiflazz-gateway

# View logs
sudo journalctl -u digiflazz-gateway -f
```

### Windows Deployment
```bash
# Build for Windows
make build-windows

# Deploy natively
scripts\deploy.bat windows

# Deploy using Docker
scripts\deploy.bat docker

# Deploy using Docker Compose
scripts\deploy.bat compose
```

**Native Windows Deployment:**
- Creates Windows service using NSSM
- Installs to `C:\digiflazz-gateway\`
- Logs to `C:\digiflazz-gateway\logs\`
- Config in `C:\digiflazz-gateway\config\`

**Service Management:**
```cmd
# Check status
sc query digiflazz-gateway

# Start/stop/restart
sc start digiflazz-gateway
sc stop digiflazz-gateway
sc stop digiflazz-gateway && sc start digiflazz-gateway

# View logs
type "C:\digiflazz-gateway\logs\output.log"
```

### macOS Deployment
```bash
# Build for macOS
make build-darwin

# Deploy natively
./scripts/deploy.sh darwin

# Deploy using Docker
./scripts/deploy.sh docker

# Deploy using Docker Compose
./scripts/deploy.sh compose
```

**Native macOS Deployment:**
- Creates launchd service
- Installs to `/opt/digiflazz-gateway/`
- Logs to `/var/log/digiflazz-gateway/`
- Config in `/etc/digiflazz-gateway/`

**Service Management:**
```bash
# Check status
sudo launchctl list | grep digiflazz

# Start/stop/restart
sudo launchctl start com.digiflazz.gateway
sudo launchctl stop com.digiflazz.gateway

# View logs
tail -f /var/log/digiflazz-gateway/output.log
```

## Docker Deployment

### Single Container
```bash
# Build Docker image
docker build -t digiflazz-gateway:latest .

# Run container
docker run -d \
  --name digiflazz-gateway \
  --restart unless-stopped \
  -p 8080:8080 \
  -v $(pwd)/logs:/app/logs \
  -v $(pwd)/cache:/app/cache \
  -e DIGIFLAZZ_USERNAME=your_username \
  -e DIGIFLAZZ_API_KEY=your_api_key \
  digiflazz-gateway:latest
```

### Multi-Architecture Docker
```bash
# Build for multiple architectures
docker buildx build --platform linux/amd64,linux/arm64 -t digiflazz-gateway:latest .

# Deploy using multi-arch image
docker run -d \
  --name digiflazz-gateway \
  --platform linux/amd64 \
  -p 8080:8080 \
  digiflazz-gateway:latest
```

### Docker Compose
```bash
# Start all services
docker-compose up -d

# Start with multi-architecture support
docker-compose -f docker-compose.multiarch.yml up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## CI/CD with GitHub Actions

### Automatic Builds
The GitHub Actions workflow automatically:
- Runs tests on every push
- Builds for all platforms on tags
- Creates release packages
- Builds and pushes Docker images

### Manual Release
```bash
# Create a new release
git tag v1.0.0
git push origin v1.0.0

# GitHub Actions will automatically:
# 1. Build for all platforms
# 2. Create release archives
# 3. Build Docker images
# 4. Create GitHub release
```

## Build Artifacts

### Binary Files
- `gateway-digiflazz-linux-amd64`
- `gateway-digiflazz-linux-arm64`
- `gateway-digiflazz-linux-arm`
- `gateway-digiflazz-windows-amd64.exe`
- `gateway-digiflazz-windows-arm64.exe`
- `gateway-digiflazz-darwin-amd64`
- `gateway-digiflazz-darwin-arm64`
- `gateway-digiflazz-freebsd-amd64`
- `gateway-digiflazz-freebsd-arm64`

### Archive Files
- `gateway-digiflazz-linux-amd64.tar.gz`
- `gateway-digiflazz-windows-amd64.zip`
- `gateway-digiflazz-darwin-amd64.tar.gz`
- `gateway-digiflazz-freebsd-amd64.tar.gz`

## Environment Configuration

### Environment Variables
```bash
# Digiflazz API Configuration
DIGIFLAZZ_USERNAME=your_username
DIGIFLAZZ_API_KEY=your_api_key
DIGIFLAZZ_BASE_URL=https://api.digiflazz.com

# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
LOG_LEVEL=info

# Cache Configuration
CACHE_ENABLED=true
CACHE_TTL=24h

# Otomax Configuration
OTOMAX_SECRET_KEY=your_secret_key
```

### Configuration Files
- `configs/config.yaml` - Main configuration
- `configs/.env.example` - Environment variables template
- `nginx.conf` - Nginx configuration (if using reverse proxy)

## Monitoring and Maintenance

### Health Checks
```bash
# Check application health
curl http://localhost:8080/health

# Check specific endpoints
curl http://localhost:8080/api/v1/balance
curl http://localhost:8080/api/v1/pln/stats
```

### Log Management
```bash
# View application logs
tail -f logs/app.log

# View system logs (Linux)
sudo journalctl -u digiflazz-gateway -f

# View Docker logs
docker logs -f digiflazz-gateway
```

### Performance Monitoring
```bash
# Check resource usage
htop
docker stats digiflazz-gateway

# Check cache statistics
curl http://localhost:8080/api/v1/pln/stats
```

## Troubleshooting

### Common Issues

1. **Build Failures**
   - Check Go version (requires 1.21+)
   - Verify all dependencies are installed
   - Check build environment

2. **Deployment Issues**
   - Verify binary exists in build directory
   - Check service permissions
   - Verify configuration files

3. **Runtime Issues**
   - Check application logs
   - Verify environment variables
   - Check network connectivity

### Debug Commands
```bash
# Check build artifacts
ls -la build/

# Test binary
./build/gateway-digiflazz-linux-amd64 --help

# Check service status
sudo systemctl status digiflazz-gateway

# View detailed logs
sudo journalctl -u digiflazz-gateway --no-pager
```

## Security Considerations

1. **Binary Security**
   - Verify binary signatures
   - Use trusted build environments
   - Regular security updates

2. **Service Security**
   - Run as non-root user
   - Restrict file permissions
   - Use secure configuration

3. **Network Security**
   - Use HTTPS in production
   - Implement rate limiting
   - Monitor access logs

## Performance Optimization

1. **Build Optimization**
   - Use static linking
   - Strip debug symbols
   - Optimize for target platform

2. **Runtime Optimization**
   - Configure appropriate cache settings
   - Monitor memory usage
   - Optimize database connections

3. **Deployment Optimization**
   - Use appropriate platform-specific optimizations
   - Configure resource limits
   - Implement health checks



