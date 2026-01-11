#!/bin/bash
set -e  # Exit on error

APP_NAME="tiny-lang"
BUILD_DIR="build"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")

echo "Building $APP_NAME version $VERSION..."

# Clean and recreate build directory
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

# Build function for cleaner code
build() {
    local os=$1
    local arch=$2
    local ext=$3
    local output="$BUILD_DIR/${APP_NAME}-${VERSION}-${os}-${arch}${ext}"
    
    echo "Building $os/$arch..."
    GOOS=$os GOARCH=$arch go build -ldflags "-s -w" -o "$output" .
    
    if [ $? -eq 0 ]; then
        echo "  ✓ Created $output ($(du -h "$output" | cut -f1))"
    else
        echo "  ✗ Failed to build $os/$arch"
        return 1
    fi
}

# macOS
build darwin amd64 ""
build darwin arm64 ""

# Windows
build windows amd64 ".exe"
build windows arm64 ".exe"

# Linux
build linux amd64 ""
build linux arm64 ""

echo ""
echo "Build complete! Binaries in $BUILD_DIR/"
ls -lh "$BUILD_DIR"