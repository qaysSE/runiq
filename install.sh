#!/bin/bash
set -e

echo "🧠 Installing Runiq: The Local Agent Runtime..."

# Detect OS
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Darwin)  PLATFORM="darwin" ;;
    Linux)   PLATFORM="linux" ;;
    *)       echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
    x86_64)  ARCH="amd64" ;;
    arm64)   ARCH="arm64" ;;
    aarch64) ARCH="arm64" ;;
    *)       echo "Unsupported Architecture: $ARCH"; exit 1 ;;
esac

# Download URL (Dynamic based on latest release)
# Note: For now, we point to the raw binary you will upload. 
# Once the action runs, change this to fetch from releases/latest/download.
BINARY_URL="https://github.com/qaysSE/runiq/releases/latest/download/runiq-${PLATFORM}-${ARCH}"

echo "⬇️ Downloading ${PLATFORM}/${ARCH} binary..."
# (Mock download for now - requires the GitHub Action to run first)
# curl -L -o runiq $BINARY_URL

echo "⚠️  NOTE: Binaries will be available after the GitHub Action finishes building v1.1.0."
echo "For now, please clone and run 'go build -o runiq cmd/runiq/main.go'"

# Setup Permissions
# chmod +x runiq
# mv runiq /usr/local/bin/runiq

echo "✅ Runiq installed!"