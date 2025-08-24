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
if [ "$OS" != "linux" ] && [ "$OS" != "darwin" ] && [ "$OS" != "windows" ]; then
    echo "Unsupported OS: $OS. This script supports Linux, macOS, and Windows (via WSL)."
    exit 1
fi

# For Windows, we need to adjust the binary name
if [ "$OS" = "windows" ]; then
    BINARY_NAME="buildtree.exe"
    ARCHIVE_NAME="buildtree_windows_${ARCH}.tar.gz"
else
    BINARY_NAME="buildtree"
    ARCHIVE_NAME="buildtree_${OS}_${ARCH}.tar.gz"
fi

# Download URL
URL="https://github.com/neomen/buildtree/releases/latest/download/${ARCHIVE_NAME}"
echo "Downloading from $URL..."

# Download and extract
if command -v curl >/dev/null 2>&1; then
    curl -sSL "$URL" | tar xz
elif command -v wget >/dev/null 2>&1; then
    wget -q -O - "$URL" | tar xz
else
    echo "Error: curl or wget is required to download buildtree"
    exit 1
fi

# Installation
if [ "$OS" = "windows" ]; then
    # For Windows, offer to install to a directory in PATH
    echo "Windows installation:"
    echo "1. Move $BINARY_NAME to a directory in your PATH"
    echo "2. Or run directly from current directory: .\\$BINARY_NAME"
    mv "$BINARY_NAME" ./
    echo "Binary downloaded to current directory: $BINARY_NAME"
else
    # For Unix systems, install to /usr/local/bin
    sudo mv "$BINARY_NAME" /usr/local/bin/
    echo "Installed in /usr/local/bin/$BINARY_NAME"

    # Test the installation
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        echo "Installation successful. You can now use '$BINARY_NAME'"
        "$BINARY_NAME" -v
    else
        echo "Installation completed, but the binary is not in your PATH."
    fi
fi