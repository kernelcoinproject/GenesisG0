#!/bin/bash

set -e

# Source file
SOURCE="genesis.go"
OUTPUT_DIR="builds"

# Create output directory
mkdir -p "$OUTPUT_DIR"

echo "Building and compressing for all platforms..."

# Windows x86 (32-bit)
echo "Building Windows x86..."
GOOS=windows GOARCH=386 go build -o "$OUTPUT_DIR/genesis-windows-x86.exe" "$SOURCE"
zip -j "$OUTPUT_DIR/genesis-windows-x86.zip" "$OUTPUT_DIR/genesis-windows-x86.exe"
rm "$OUTPUT_DIR/genesis-windows-x86.exe"
echo "✓ Windows x86 compressed"

# Windows x64 (64-bit)
echo "Building Windows x64..."
GOOS=windows GOARCH=amd64 go build -o "$OUTPUT_DIR/genesis-windows-x64.exe" "$SOURCE"
zip -j "$OUTPUT_DIR/genesis-windows-x64.zip" "$OUTPUT_DIR/genesis-windows-x64.exe"
rm "$OUTPUT_DIR/genesis-windows-x64.exe"
echo "✓ Windows x64 compressed"

# Linux x86 (32-bit)
echo "Building Linux x86..."
GOOS=linux GOARCH=386 go build -o "$OUTPUT_DIR/genesis-linux-x86" "$SOURCE"
tar -czf "$OUTPUT_DIR/genesis-linux-x86.tar.gz" -C "$OUTPUT_DIR" genesis-linux-x86
rm "$OUTPUT_DIR/genesis-linux-x86"
echo "✓ Linux x86 compressed"

# Linux x64 (64-bit)
echo "Building Linux x64..."
GOOS=linux GOARCH=amd64 go build -o "$OUTPUT_DIR/genesis-linux-x64" "$SOURCE"
tar -czf "$OUTPUT_DIR/genesis-linux-x64.tar.gz" -C "$OUTPUT_DIR" genesis-linux-x64
rm "$OUTPUT_DIR/genesis-linux-x64"
echo "✓ Linux x64 compressed"

# Linux ARM (32-bit)
echo "Building Linux ARM..."
GOOS=linux GOARCH=arm go build -o "$OUTPUT_DIR/genesis-linux-arm" "$SOURCE"
tar -czf "$OUTPUT_DIR/genesis-linux-arm.tar.gz" -C "$OUTPUT_DIR" genesis-linux-arm
rm "$OUTPUT_DIR/genesis-linux-arm"
echo "✓ Linux ARM compressed"

# Linux ARM64 (64-bit)
echo "Building Linux ARM64..."
GOOS=linux GOARCH=arm64 go build -o "$OUTPUT_DIR/genesis-linux-arm64" "$SOURCE"
tar -czf "$OUTPUT_DIR/genesis-linux-arm64.tar.gz" -C "$OUTPUT_DIR" genesis-linux-arm64
rm "$OUTPUT_DIR/genesis-linux-arm64"
echo "✓ Linux ARM64 compressed"

# macOS x64 (Intel)
echo "Building macOS x64..."
GOOS=darwin GOARCH=amd64 go build -o "$OUTPUT_DIR/genesis-macos-x64" "$SOURCE"
tar -czf "$OUTPUT_DIR/genesis-macos-x64.tar.gz" -C "$OUTPUT_DIR" genesis-macos-x64
rm "$OUTPUT_DIR/genesis-macos-x64"
echo "✓ macOS x64 compressed"

# macOS ARM (Apple Silicon)
echo "Building macOS ARM..."
GOOS=darwin GOARCH=arm64 go build -o "$OUTPUT_DIR/genesis-macos-arm" "$SOURCE"
tar -czf "$OUTPUT_DIR/genesis-macos-arm.tar.gz" -C "$OUTPUT_DIR" genesis-macos-arm
rm "$OUTPUT_DIR/genesis-macos-arm"
echo "✓ macOS ARM compressed"

echo ""
echo "All builds complete!"
echo "Files saved to: $OUTPUT_DIR"
ls -lh "$OUTPUT_DIR"
