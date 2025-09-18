#!/bin/bash

# Build script for Golf Gamez WebAssembly frontend

set -e

echo "Building Golf Gamez WebAssembly frontend..."

# Set GOOS and GOARCH for WebAssembly
export GOOS=js
export GOARCH=wasm

# Build the WebAssembly binary
echo "Compiling Go to WebAssembly..."
go build -o web/static/js/main.wasm main.go

# Copy the WebAssembly support file
echo "Copying WebAssembly support files..."
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" web/static/js/

# Optimize for production if requested
if [ "$1" = "production" ]; then
    echo "Optimizing for production..."
    # Add any additional optimization steps here
    # For now, we'll compress the wasm file if brotli is available
    if command -v brotli &> /dev/null; then
        echo "Compressing WebAssembly binary..."
        brotli -f web/static/js/main.wasm
    fi
fi

echo "Build complete!"
echo "WebAssembly binary: web/static/js/main.wasm"
echo "Support file: web/static/js/wasm_exec.js"