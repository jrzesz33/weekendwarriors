#!/bin/bash

# Golf Gamez Frontend Setup Script
# This script sets up the development environment for the Golf Gamez PWA

set -e  # Exit on any error

echo "🏌️  Golf Gamez Frontend Setup"
echo "=============================="
echo ""

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check Go version
check_go() {
    info "Checking Go installation..."

    if ! command_exists go; then
        error "Go is not installed. Please install Go 1.21+ from https://golang.org/dl/"
        exit 1
    fi

    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    REQUIRED_VERSION="1.21"

    if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
        error "Go version $GO_VERSION is too old. Please install Go $REQUIRED_VERSION or newer."
        exit 1
    fi

    success "Go $GO_VERSION is installed"
}

# Check Make
check_make() {
    info "Checking Make installation..."

    if ! command_exists make; then
        error "Make is not installed. Please install Make for your platform."
        exit 1
    fi

    success "Make is installed"
}

# Setup Go modules
setup_go_modules() {
    info "Setting up Go modules..."

    if [ ! -f "go.mod" ]; then
        error "go.mod not found. Make sure you're in the frontend directory."
        exit 1
    fi

    go mod tidy
    go mod download

    success "Go modules configured"
}

# Download wasm_exec.js
setup_wasm_exec() {
    info "Setting up WebAssembly support..."

    if [ ! -f "wasm_exec.js" ]; then
        GOROOT=$(go env GOROOT)
        if [ -f "$GOROOT/misc/wasm/wasm_exec.js" ]; then
            cp "$GOROOT/misc/wasm/wasm_exec.js" .
            success "wasm_exec.js copied from Go installation"
        else
            warning "Could not find wasm_exec.js in Go installation"
            info "You may need to download it manually or run 'make ensure-wasm-exec'"
        fi
    else
        success "wasm_exec.js already exists"
    fi
}

# Create placeholder icons
setup_icons() {
    info "Setting up PWA icons..."

    mkdir -p web/static

    # Create a simple script to generate icons from the SVG placeholder
    if command_exists convert; then
        info "ImageMagick found - generating PNG icons from SVG..."

        SIZES=(72 96 128 144 152 180 192 384 512)
        for size in "${SIZES[@]}"; do
            if [ ! -f "web/static/icon-$size.png" ]; then
                convert "web/static/icon-placeholder.svg" -resize "${size}x${size}" "web/static/icon-$size.png"
            fi
        done

        success "PWA icons generated"
    else
        warning "ImageMagick not found - using placeholder icons"
        info "Install ImageMagick and run 'make generate-icons' to create proper PWA icons"

        # Copy SVG as fallback
        for size in 72 96 128 144 152 180 192 384 512; do
            if [ ! -f "web/static/icon-$size.png" ]; then
                cp "web/static/icon-placeholder.svg" "web/static/icon-$size.png"
            fi
        done
    fi
}

# Install development tools
install_dev_tools() {
    info "Installing development tools..."

    # Install golangci-lint for linting
    if ! command_exists golangci-lint; then
        info "Installing golangci-lint..."
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
        success "golangci-lint installed"
    else
        success "golangci-lint already installed"
    fi

    # Install entr for file watching (macOS/Linux)
    if [[ "$OSTYPE" == "darwin"* ]]; then
        if command_exists brew && ! command_exists entr; then
            info "Installing entr for file watching..."
            brew install entr
            success "entr installed"
        fi
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if command_exists apt-get && ! command_exists entr; then
            info "You may want to install entr for file watching: sudo apt-get install entr"
        fi
    fi
}

# Verify backend API
check_backend() {
    info "Checking backend API availability..."

    if command_exists curl; then
        if curl -s -f "http://localhost:8080/v1/health" >/dev/null 2>&1; then
            success "Backend API is running on http://localhost:8080"
        else
            warning "Backend API is not running on http://localhost:8080"
            info "Make sure to start the backend server before running the frontend"
        fi
    else
        warning "curl not found - cannot check backend API status"
    fi
}

# Build the application
build_app() {
    info "Building the application..."

    make clean
    make build-dev

    success "Application built successfully"
}

# Create sample environment file
create_env_file() {
    info "Creating environment configuration..."

    if [ ! -f ".env.example" ]; then
        cat > .env.example << 'EOF'
# Golf Gamez Frontend Environment Configuration

# Development server settings
DEV_PORT=8000
API_BASE_URL=http://localhost:8080/v1

# PWA settings
PWA_NAME="Golf Gamez"
PWA_SHORT_NAME="Golf Gamez"
PWA_THEME_COLOR="#2e7d32"

# Debug settings
DEBUG=false
LOG_LEVEL=info

# Feature flags
ENABLE_OFFLINE_MODE=true
ENABLE_PUSH_NOTIFICATIONS=false
ENABLE_ANALYTICS=false
EOF
        success "Created .env.example file"
    fi
}

# Main setup process
main() {
    echo "Starting setup process..."
    echo ""

    # Pre-flight checks
    check_go
    check_make

    # Setup steps
    setup_go_modules
    setup_wasm_exec
    setup_icons
    create_env_file
    install_dev_tools
    build_app
    check_backend

    echo ""
    echo "🎉 Setup completed successfully!"
    echo ""
    echo "Next steps:"
    echo "1. Start the backend API server (if not already running)"
    echo "2. Run 'make dev' to start the development server"
    echo "3. Open http://localhost:8000 in your browser"
    echo "4. For mobile testing, use http://[your-ip]:8000"
    echo ""
    echo "Development commands:"
    echo "  make dev          - Start development server with auto-reload"
    echo "  make build        - Build production version"
    echo "  make test         - Run tests"
    echo "  make lint         - Run linter"
    echo "  make clean        - Clean build artifacts"
    echo ""
    echo "Happy coding! 🏌️‍♀️"
}

# Handle script arguments
case "${1:-setup}" in
    "setup"|"")
        main
        ;;
    "check")
        check_go
        check_make
        check_backend
        ;;
    "tools")
        install_dev_tools
        ;;
    "icons")
        setup_icons
        ;;
    "build")
        build_app
        ;;
    "help")
        echo "Golf Gamez Frontend Setup Script"
        echo ""
        echo "Usage: $0 [command]"
        echo ""
        echo "Commands:"
        echo "  setup (default) - Run full setup process"
        echo "  check          - Check system requirements"
        echo "  tools          - Install development tools"
        echo "  icons          - Generate PWA icons"
        echo "  build          - Build the application"
        echo "  help           - Show this help message"
        ;;
    *)
        error "Unknown command: $1"
        echo "Run '$0 help' for usage information"
        exit 1
        ;;
esac