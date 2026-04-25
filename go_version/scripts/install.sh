#!/bin/bash
set -e

INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
RELEASE_URL="${RELEASE_URL:-https://github.com/kora-ai-lab/prometheus/releases/latest/download}"

echo "Detecting platform..."
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  i386|i686) ARCH="386" ;;
esac

case "$OS" in
  darwin) OS="darwin" ;;
  linux) OS="linux" ;;
  windows*) OS="windows" ;;
esac

BINARY_NAME="prometheus"
if [ "$OS" = "windows" ]; then
  BINARY_NAME="prometheus.exe"
fi

URL="$RELEASE_URL/prometheus-$OS-$ARCH"
echo "Downloading from $URL..."

curl -fsSL "$URL" -o "$INSTALL_DIR/prometheus"
chmod +x "$INSTALL_DIR/prometheus"

echo "Installed prometheus to $INSTALL_DIR/prometheus"
echo "Run 'prometheus --help' to get started."