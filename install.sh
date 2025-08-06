#!/bin/sh
# install.sh

set -e

# OS and Architecture definition
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
[ "$ARCH" = "x86_64" ] && ARCH="amd64"
[ "$ARCH" = "aarch64" ] && ARCH="arm64"

# Download
URL="https://github.com/neomen/buildtree/releases/latest/download/buildtree_${OS}_${ARCH}.tar.gz"
echo "Download from the $URL..."
curl -L "$URL" | tar xz

# Installation
sudo mv buildtree /usr/local/bin/
echo "Installed in /usr/local/bin/buildtree"
buildtree -h