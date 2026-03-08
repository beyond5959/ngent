#!/bin/bash

set -e

# Ngent Installer Script
# Usage: curl -sSL https://raw.githubusercontent.com/beyond5959/ngent/master/install.sh | bash

REPO="beyond5959/ngent"
BINARY_NAME="ngent"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print functions
info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS
detect_os() {
    local os
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    case "$os" in
        linux)
            echo "linux"
            ;;
        darwin)
            echo "darwin"
            ;;
        mingw*|msys*|cygwin*)
            echo "windows"
            ;;
        *)
            error "Unsupported operating system: $os"
            exit 1
            ;;
    esac
}

# Detect architecture
detect_arch() {
    local arch
    arch=$(uname -m)
    case "$arch" in
        x86_64|amd64)
            echo "x86_64"
            ;;
        arm64|aarch64)
            echo "aarch64"
            ;;
        *)
            error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
}

# Get latest release version
get_latest_version() {
    local version
    version=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$version" ]; then
        error "Failed to get latest version"
        exit 1
    fi
    echo "$version"
}

# Download and install
download_and_install() {
    local version="$1"
    local os="$2"
    local arch="$3"
    
    local suffix="${arch}-${os}"
    local filename
    local extract_cmd
    
    if [ "$os" = "windows" ]; then
        filename="${BINARY_NAME}-${version}-${suffix}.zip"
        extract_cmd="unzip -o"
    else
        filename="${BINARY_NAME}-${version}-${suffix}.tar.gz"
        extract_cmd="tar -xzf"
    fi
    
    local download_url="https://github.com/${REPO}/releases/download/${version}/${filename}"
    local tmp_dir
    tmp_dir=$(mktemp -d)
    
    info "Downloading ${BINARY_NAME} ${version} for ${os}/${arch}..."
    info "URL: ${download_url}"
    
    if ! curl -sSL -o "${tmp_dir}/${filename}" "$download_url"; then
        error "Failed to download ${filename}"
        rm -rf "$tmp_dir"
        exit 1
    fi
    
    info "Extracting..."
    cd "$tmp_dir"
    if ! $extract_cmd "$filename" 2>/dev/null; then
        error "Failed to extract ${filename}"
        rm -rf "$tmp_dir"
        exit 1
    fi
    
    # Check if binary exists
    if [ ! -f "$BINARY_NAME" ]; then
        error "Binary not found in archive"
        rm -rf "$tmp_dir"
        exit 1
    fi
    
    # Make executable
    chmod +x "$BINARY_NAME"
    
    # Install
    info "Installing to ${INSTALL_DIR}..."
    
    # Check if we need sudo
    if [ -w "$INSTALL_DIR" ]; then
        mv "$BINARY_NAME" "${INSTALL_DIR}/"
    else
        warn "Need sudo permission to install to ${INSTALL_DIR}"
        if command -v sudo >/dev/null 2>&1; then
            sudo mv "$BINARY_NAME" "${INSTALL_DIR}/"
        else
            error "sudo not found and cannot write to ${INSTALL_DIR}"
            error "Please run with sudo or set INSTALL_DIR to a writable directory"
            rm -rf "$tmp_dir"
            exit 1
        fi
    fi
    
    # Cleanup
    rm -rf "$tmp_dir"
    
    success "${BINARY_NAME} ${version} installed successfully!"
}

# Verify installation
verify_installation() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local version
        version=$($BINARY_NAME --version 2>/dev/null || echo "unknown")
        success "Installation verified: ${BINARY_NAME} is available in PATH"
        info "Version: ${version}"
    else
        warn "${BINARY_NAME} is installed but not in your PATH"
        info "Add ${INSTALL_DIR} to your PATH or run: export PATH=\"${INSTALL_DIR}:\$PATH\""
    fi
}

# Main
main() {
    echo "========================================"
    echo "  ${BINARY_NAME} Installer"
    echo "========================================"
    echo ""
    
    # Check for required commands
    if ! command -v curl >/dev/null 2>&1; then
        error "curl is required but not installed"
        exit 1
    fi
    
    # Detect system
    local os
    local arch
    os=$(detect_os)
    arch=$(detect_arch)
    
    info "Detected OS: ${os}"
    info "Detected Architecture: ${arch}"
    
    # Check if Windows
    if [ "$os" = "windows" ]; then
        warn "Windows detected. Please download the .zip file manually from:"
        info "https://github.com/${REPO}/releases"
        exit 0
    fi
    
    # Get version
    local version
    version=$(get_latest_version)
    info "Latest version: ${version}"
    
    # Download and install
    download_and_install "$version" "$os" "$arch"
    
    # Verify
    verify_installation
    
    echo ""
    success "Installation complete!"
    info "Run '${BINARY_NAME} --help' to get started"
}

# Run main function
main "$@"
