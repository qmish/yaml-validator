#!/bin/bash
# Кросс-платформенная сборка (Linux/macOS)
set -e

OUT_DIR="bin"
mkdir -p "$OUT_DIR"

echo "Building for linux/amd64..."
GOOS=linux GOARCH=amd64 go build -o "$OUT_DIR/yaml-validator-linux-amd64" .

echo "Building for linux/arm64..."
GOOS=linux GOARCH=arm64 go build -o "$OUT_DIR/yaml-validator-linux-arm64" .

echo "Building for darwin/amd64..."
GOOS=darwin GOARCH=amd64 go build -o "$OUT_DIR/yaml-validator-darwin-amd64" .

echo "Building for darwin/arm64..."
GOOS=darwin GOARCH=arm64 go build -o "$OUT_DIR/yaml-validator-darwin-arm64" .

echo "Building for windows/amd64..."
GOOS=windows GOARCH=amd64 go build -o "$OUT_DIR/yaml-validator-windows-amd64.exe" .

echo "Done. Binaries in $OUT_DIR/"
