#!/bin/sh
# Buildtree installer script

set -e

echo "Installing buildtree..."

# Determine OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Check if the OS is supported
if [ "$OS" != "linux" ] && [ "$OS" != "darwin" ]; then
    echo "Unsupported OS: $OS. For Windows, please use install.ps1"
    exit 1
fi

# Create a temporary directory for extraction
TMPDIR=$(mktemp -d)
cleanup() {
    rm -rf "$TMPDIR"
}
trap cleanup EXIT

# Download URL
URL="https://github.com/neomen/buildtree/releases/latest/download/buildtree_${OS}_${ARCH}.tar.gz"
echo "Downloading from $URL..."

# Download and extract to temporary directory
if command -v curl >/dev/null 2>&1; then
    curl -sSL "$URL" | tar -xz -C "$TMPDIR"
elif command -v wget >/dev/null 2>&1; then
    wget -q -O - "$URL" | tar -xz -C "$TMPDIR"
else
    echo "Error: curl or wget is required to download buildtree"
    exit 1
fi

# Check if the binary was extracted
if [ -f "$TMPDIR/buildtree" ]; then
    BINARY_PATH="$TMPDIR/buildtree"
else
    echo "Error: binary not found in the downloaded archive"
    exit 1
fi

# Installation
sudo mv "$BINARY_PATH" /usr/local/bin/buildtree
echo "Installed in /usr/local/bin/buildtree"

# Test the installation
if command -v buildtree >/dev/null 2>&1; then
    echo "Installation successful. You can now use 'buildtree' command"
    buildtree -v
else
    echo "Installation completed, but the binary is not in your PATH."
    echo "Please check if /usr/local/bin is in your PATH environment variable."
fi