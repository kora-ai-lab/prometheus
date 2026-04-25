#!/bin/bash
set -e

VERSION="${VERSION:-$(git describe --tags --always --dirty)}"
OUTPUT_DIR="${OUTPUT_DIR:-./release}"

mkdir -p "$OUTPUT_DIR"

PLATFORMS=(
  "darwin/amd64"
  "darwin/arm64"
  "linux/amd64"
  "linux/arm64"
  "windows/amd64"
)

echo "Building prometheus $VERSION..."

for plat in "${PLATFORMS[@]}"; do
  OS="${plat%/*}"
  ARCH="${plat#*/}"

  OUTPUT="$OUTPUT_DIR/prometheus-$OS-$ARCH"
  if [ "$OS" = "windows" ]; then
    OUTPUT="$OUTPUT_DIR/prometheus-$OS-$ARCH.exe"
  fi

  echo "Building $OS/$ARCH..."
  GOOS="$OS" GOARCH="$ARCH" go build -ldflags="-s -w" -o "$OUTPUT" ./cmd/prometheus
done

echo "Build complete:"
ls -lh "$OUTPUT_DIR"