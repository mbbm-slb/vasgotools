#!/bin/bash
# Cross-platform build script for Go projects
# Builds executables for Windows, Linux, and macOS
#
# Usage:
#   chmod +x cross-build.sh
#   ./cross-build.sh
#
# Output:
#   Binaries are created in the bin/ directory

echo "Building Go project for multiple platforms..."
echo ""

# Get module name from go.mod and extract the last part as binary name
if [ ! -f "go.mod" ]; then
    echo "Error: go.mod not found in current directory"
    exit 1
fi

MODULE_NAME=$(grep "^module " go.mod | awk '{print $2}')
BINARY_NAME=$(basename "$MODULE_NAME")

echo "Module: $MODULE_NAME"
echo "Binary name: $BINARY_NAME"
echo ""

# Create output directory
mkdir -p bin

# Build for Windows (amd64)
echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o "bin/${BINARY_NAME}-windows-amd64.exe"
if [ $? -ne 0 ]; then
    echo "Failed to build for Windows amd64"
    exit 1
fi

# Build for Linux (amd64)
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o "bin/${BINARY_NAME}-linux-amd64"
if [ $? -ne 0 ]; then
    echo "Failed to build for Linux amd64"
    exit 1
fi

# Build for macOS (amd64 - Intel)
echo "Building for macOS (amd64 - Intel)..."
GOOS=darwin GOARCH=amd64 go build -o "bin/${BINARY_NAME}-darwin-amd64"
if [ $? -ne 0 ]; then
    echo "Failed to build for macOS amd64"
    exit 1
fi

# Build for macOS (arm64 - Apple Silicon)
echo "Building for macOS (arm64 - Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -o "bin/${BINARY_NAME}-darwin-arm64"
if [ $? -ne 0 ]; then
    echo "Failed to build for macOS arm64"
    exit 1
fi

echo ""
echo "Build completed successfully!"
echo "Binaries are located in the bin/ directory:"
echo "  - ${BINARY_NAME}-windows-amd64.exe (Windows)"
echo "  - ${BINARY_NAME}-linux-amd64 (Linux)"
echo "  - ${BINARY_NAME}-darwin-amd64 (macOS Intel)"
echo "  - ${BINARY_NAME}-darwin-arm64 (macOS Apple Silicon)"
echo ""
