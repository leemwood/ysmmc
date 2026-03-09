#!/bin/bash

VERSION=${1:-"1.0.0"}
OUTPUT_DIR="build"

mkdir -p $OUTPUT_DIR

echo "Building version $VERSION..."

# Linux AMD64
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=$VERSION" -o $OUTPUT_DIR/backend-$VERSION-linux-amd64 ./cmd/server
echo "Built backend-$VERSION-linux-amd64"

# Linux ARM64
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=$VERSION" -o $OUTPUT_DIR/backend-$VERSION-linux-arm64 ./cmd/server
echo "Built backend-$VERSION-linux-arm64"

# Windows AMD64
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=$VERSION" -o $OUTPUT_DIR/backend-$VERSION-windows-amd64.exe ./cmd/server
echo "Built backend-$VERSION-windows-amd64.exe"

# macOS AMD64
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=$VERSION" -o $OUTPUT_DIR/backend-$VERSION-darwin-amd64 ./cmd/server
echo "Built backend-$VERSION-darwin-amd64"

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=$VERSION" -o $OUTPUT_DIR/backend-$VERSION-darwin-arm64 ./cmd/server
echo "Built backend-$VERSION-darwin-arm64"

echo ""
echo "Build complete! Binaries are in the $OUTPUT_DIR directory:"
ls -la $OUTPUT_DIR
