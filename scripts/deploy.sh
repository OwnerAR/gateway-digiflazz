#!/bin/bash

# Cross-platform deployment script for Digiflazz Gateway
# Usage: ./scripts/deploy.sh [platform] [environment]

set -e

# Default values
PLATFORM=${1:-"auto"}
ENVIRONMENT=${2:-"production"}
VERSION=${3:-"latest"}
SERVICE_NAME="digiflazz-gateway"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect platform
detect_platform() {
    if [ "$PLATFORM" = "auto" ]; then
        case "$(uname -s)" in
            Linux*)     PLATFORM="linux";;
            Darwin*)    PLATFORM="darwin";;
            CYGWIN*)    PLATFORM="windows";;
            MINGW*)     PLATFORM="windows";;
            FreeBSD*)   PLATFORM="freebsd";;
            *)          PLATFORM="unknown";;
        esac
    fi
    
    # Detect architecture
    case "$(uname -m)" in
        x86_64)     ARCH="amd64";;
        arm64)      ARCH="arm64";;
        armv7l)     ARCH="arm";;
        *)          ARCH="unknown";;
    esac
    
    print_status "Detected platform: $PLATFORM/$ARCH"
}

# Check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    # Check if binary exists
    BINARY_PATH="build/gateway-digiflazz-${PLATFORM}-${ARCH}"
    if [ "$PLATFORM" = "windows" ]; then
        BINARY_PATH="${BINARY_PATH}.exe"
    fi
    
    if [ ! -f "$BINARY_PATH" ]; then
        print_error "Binary not found: $BINARY_PATH"
        print_status "Please run 'make build-all' first"
        exit 1
    fi
    
    print_success "Binary found: $BINARY_PATH"
}

# Deploy to Linux
deploy_linux() {
    print_status "Deploying to Linux..."
    
    # Create systemd service
    sudo tee /etc/systemd/system/${SERVICE_NAME}.service > /dev/null <<EOF
[Unit]
Description=Digiflazz Gateway
After=network.target

[Service]
Type=simple
User=digiflazz
Group=digiflazz
WorkingDirectory=/opt/${SERVICE_NAME}
ExecStart=/opt/${SERVICE_NAME}/gateway-digiflazz-${PLATFORM}-${ARCH}
Restart=always
RestartSec=5
Environment=LOG_LEVEL=info
Environment=SERVER_PORT=8080

[Install]
WantedBy=multi-user.target
EOF

    # Create user and directories
    sudo useradd -r -s /bin/false digiflazz 2>/dev/null || true
    sudo mkdir -p /opt/${SERVICE_NAME}
    sudo mkdir -p /var/log/${SERVICE_NAME}
    sudo mkdir -p /etc/${SERVICE_NAME}
    
    # Copy binary and config
    sudo cp "build/gateway-digiflazz-${PLATFORM}-${ARCH}" /opt/${SERVICE_NAME}/
    sudo cp configs/config.yaml /etc/${SERVICE_NAME}/
    sudo cp configs/.env.example /etc/${SERVICE_NAME}/.env
    
    # Set permissions
    sudo chown -R digiflazz:digiflazz /opt/${SERVICE_NAME}
    sudo chown -R digiflazz:digiflazz /var/log/${SERVICE_NAME}
    sudo chown -R digiflazz:digiflazz /etc/${SERVICE_NAME}
    sudo chmod +x /opt/${SERVICE_NAME}/gateway-digiflazz-${PLATFORM}-${ARCH}
    
    # Enable and start service
    sudo systemctl daemon-reload
    sudo systemctl enable ${SERVICE_NAME}
    sudo systemctl start ${SERVICE_NAME}
    
    print_success "Service deployed and started"
    print_status "Check status: sudo systemctl status ${SERVICE_NAME}"
    print_status "View logs: sudo journalctl -u ${SERVICE_NAME} -f"
}

# Deploy to macOS
deploy_darwin() {
    print_status "Deploying to macOS..."
    
    # Create launchd plist
    sudo tee /Library/LaunchDaemons/com.digiflazz.gateway.plist > /dev/null <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.digiflazz.gateway</string>
    <key>ProgramArguments</key>
    <array>
        <string>/opt/${SERVICE_NAME}/gateway-digiflazz-${PLATFORM}-${ARCH}</string>
    </array>
    <key>WorkingDirectory</key>
    <string>/opt/${SERVICE_NAME}</string>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/var/log/${SERVICE_NAME}/output.log</string>
    <key>StandardErrorPath</key>
    <string>/var/log/${SERVICE_NAME}/error.log</string>
</dict>
</plist>
EOF

    # Create directories
    sudo mkdir -p /opt/${SERVICE_NAME}
    sudo mkdir -p /var/log/${SERVICE_NAME}
    sudo mkdir -p /etc/${SERVICE_NAME}
    
    # Copy binary and config
    sudo cp "build/gateway-digiflazz-${PLATFORM}-${ARCH}" /opt/${SERVICE_NAME}/
    sudo cp configs/config.yaml /etc/${SERVICE_NAME}/
    sudo cp configs/.env.example /etc/${SERVICE_NAME}/.env
    
    # Set permissions
    sudo chmod +x /opt/${SERVICE_NAME}/gateway-digiflazz-${PLATFORM}-${ARCH}
    sudo chown root:wheel /Library/LaunchDaemons/com.digiflazz.gateway.plist
    sudo chmod 644 /Library/LaunchDaemons/com.digiflazz.gateway.plist
    
    # Load and start service
    sudo launchctl load /Library/LaunchDaemons/com.digiflazz.gateway.plist
    sudo launchctl start com.digiflazz.gateway
    
    print_success "Service deployed and started"
    print_status "Check status: sudo launchctl list | grep digiflazz"
    print_status "View logs: tail -f /var/log/${SERVICE_NAME}/output.log"
}

# Deploy using Docker
deploy_docker() {
    print_status "Deploying using Docker..."
    
    # Check if Docker is installed
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed"
        exit 1
    fi
    
    # Build Docker image
    docker build -t ${SERVICE_NAME}:${VERSION} .
    
    # Stop existing container
    docker stop ${SERVICE_NAME} 2>/dev/null || true
    docker rm ${SERVICE_NAME} 2>/dev/null || true
    
    # Run new container
    docker run -d \
        --name ${SERVICE_NAME} \
        --restart unless-stopped \
        -p 8080:8080 \
        -v $(pwd)/logs:/app/logs \
        -v $(pwd)/cache:/app/cache \
        -e DIGIFLAZZ_USERNAME=${DIGIFLAZZ_USERNAME} \
        -e DIGIFLAZZ_API_KEY=${DIGIFLAZZ_API_KEY} \
        -e SERVER_PORT=8080 \
        -e LOG_LEVEL=info \
        ${SERVICE_NAME}:${VERSION}
    
    print_success "Docker container deployed and started"
    print_status "Check status: docker ps | grep ${SERVICE_NAME}"
    print_status "View logs: docker logs -f ${SERVICE_NAME}"
}

# Deploy using Docker Compose
deploy_compose() {
    print_status "Deploying using Docker Compose..."
    
    # Check if Docker Compose is installed
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose is not installed"
        exit 1
    fi
    
    # Stop existing services
    docker-compose down 2>/dev/null || true
    
    # Start services
    docker-compose up -d
    
    print_success "Docker Compose services deployed and started"
    print_status "Check status: docker-compose ps"
    print_status "View logs: docker-compose logs -f"
}

# Show help
show_help() {
    echo "Cross-platform deployment script for Digiflazz Gateway"
    echo ""
    echo "Usage: $0 [platform] [environment] [version]"
    echo ""
    echo "Arguments:"
    echo "  platform     Target platform (linux, darwin, windows, docker, compose, auto)"
    echo "  environment  Deployment environment (production, staging, development)"
    echo "  version      Version tag (default: latest)"
    echo ""
    echo "Examples:"
    echo "  $0                          # Auto-detect platform and deploy"
    echo "  $0 linux                    # Deploy to Linux"
    echo "  $0 docker production        # Deploy using Docker in production"
    echo "  $0 compose staging v1.0.0   # Deploy using Docker Compose with version"
    echo ""
    echo "Supported platforms:"
    echo "  - linux: Native Linux deployment with systemd"
    echo "  - darwin: Native macOS deployment with launchd"
    echo "  - docker: Docker container deployment"
    echo "  - compose: Docker Compose deployment"
    echo "  - auto: Auto-detect platform"
}

# Main execution
main() {
    print_status "Digiflazz Gateway Cross-Platform Deployment"
    print_status "Platform: $PLATFORM"
    print_status "Environment: $ENVIRONMENT"
    print_status "Version: $VERSION"
    echo ""
    
    # Detect platform if auto
    detect_platform
    
    # Check prerequisites
    check_prerequisites
    
    # Deploy based on platform
    case $PLATFORM in
        "linux")
            deploy_linux
            ;;
        "darwin")
            deploy_darwin
            ;;
        "docker")
            deploy_docker
            ;;
        "compose")
            deploy_compose
            ;;
        *)
            print_error "Unsupported platform: $PLATFORM"
            print_status "Supported platforms: linux, darwin, docker, compose"
            exit 1
            ;;
    esac
    
    print_success "Deployment completed successfully!"
}

# Handle help flag
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    show_help
    exit 0
fi

# Run main function
main



