#!/bin/bash
# Cross-platform build script for Syntaxia
# Builds for Windows, macOS, and Linux

set -e

SKIP_WINDOWS=false
SKIP_MACOS=false
SKIP_LINUX=false
CLEAN=false

# Check if wails is installed
if ! command -v wails &> /dev/null; then
    echo "❌ Wails not found. Install with: go install github.com/wailsapp/wails/v2/cmd/wails@latest"
    exit 1
fi

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-windows)
            SKIP_WINDOWS=true
            shift
            ;;
        --skip-macos)
            SKIP_MACOS=true
            shift
            ;;
        --skip-linux)
            SKIP_LINUX=true
            shift
            ;;
        --clean)
            CLEAN=true
            shift
            ;;
        *)
            echo "Unknown option: $1"
            echo "Usage: $0 [--skip-windows] [--skip-macos] [--skip-linux] [--clean]"
            exit 1
            ;;
    esac
done

echo "=== Syntaxia - Cross-Platform Build ==="
echo ""

# Change to backend directory
cd backend

# Build for Windows
if [ "$SKIP_WINDOWS" = false ]; then
    echo "Building for Windows (amd64)..."
    wails build -clean
    if [ $? -eq 0 ]; then
        echo "✓ Windows build complete"
    else
        echo "✗ Windows build failed"
    fi
    echo ""
fi

# Build for macOS
if [ "$SKIP_MACOS" = false ]; then
    echo "Building for macOS (universal)..."
    wails build -clean
    if [ $? -eq 0 ]; then
        echo "✓ macOS build complete"
    else
        echo "✗ macOS build failed"
    fi
    echo ""
fi

# Build for Linux
if [ "$SKIP_LINUX" = false ]; then
    echo "Building for Linux (amd64)..."
    wails build -clean
    if [ $? -eq 0 ]; then
        echo "✓ Linux build complete"
    else
        echo "✗ Linux build failed"
    fi
    echo ""
fi

cd ..

echo "=== Build Summary ==="
echo ""

# List built files
if [ -d "backend/bin" ]; then
    find backend/bin -type f -exec ls -lh {} \; | awk '{print "  " $9 " - " $5}'
fi

echo ""
echo "Build artifacts are in: backend/bin/"
