#!/bin/bash
# Build script for Linux only

set -e

echo "=== Building Syntaxia for Linux ==="
echo ""

# Get version from git tag or use dev
VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "dev")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(date +%Y-%m-%d)

echo "Version: $VERSION"
echo "Commit: $GIT_COMMIT"
echo ""

# Build ldflags for version injection
LDFLAGS="-X syntaxia/infrastructure/version.Version=$VERSION -X syntaxia/infrastructure/version.GitCommit=$GIT_COMMIT -X syntaxia/infrastructure/version.BuildDate=$BUILD_DATE"

# Check if wails is installed
if ! command -v wails &> /dev/null; then
    echo "❌ Wails not found. Install with: go install github.com/wailsapp/wails/v2/cmd/wails@latest"
    exit 1
fi

cd backend

echo "Building Linux executable..."
wails build -clean -ldflags "$LDFLAGS"

if [ $? -eq 0 ]; then
    echo ""
    echo "✓ Build successful!"
    echo ""
    
    # Show output file
    if [ -f "bin/Syntaxia" ]; then
        ls -lh bin/Syntaxia | awk '{print "Output: bin/Syntaxia\nSize: " $5}'
    fi
else
    echo ""
    echo "✗ Build failed!"
    exit 1
fi

cd ..
