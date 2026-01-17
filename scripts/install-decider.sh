#!/usr/bin/env bash
# Install decider CLI from GitHub Releases with checksum verification.
# Installs to ./tools/decider/ (repo-local).

set -euo pipefail

REPO="sventorben/decider"
INSTALL_DIR="${INSTALL_DIR:-./tools/decider}"
VERSION_FILE="${VERSION_FILE:-./tools/decider.version}"

# Read pinned version
if [[ ! -f "$VERSION_FILE" ]]; then
    echo "Error: Version file not found: $VERSION_FILE" >&2
    exit 1
fi
VERSION=$(cat "$VERSION_FILE" | tr -d '[:space:]')
if [[ -z "$VERSION" ]]; then
    echo "Error: Version file is empty" >&2
    exit 1
fi

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
    linux)  OS="linux" ;;
    darwin) OS="darwin" ;;
    *)
        echo "Error: Unsupported OS: $OS" >&2
        exit 1
        ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64|amd64)  ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *)
        echo "Error: Unsupported architecture: $ARCH" >&2
        exit 1
        ;;
esac

# Build artifact names
VERSION_NUM="${VERSION#v}"
ARCHIVE_NAME="decider_${VERSION_NUM}_${OS}_${ARCH}.tar.gz"
CHECKSUMS_NAME="checksums.txt"
BASE_URL="https://github.com/${REPO}/releases/download/${VERSION}"
ARCHIVE_URL="${BASE_URL}/${ARCHIVE_NAME}"
CHECKSUMS_URL="${BASE_URL}/${CHECKSUMS_NAME}"

echo "Installing decider ${VERSION} for ${OS}/${ARCH}..."
echo "  Archive: ${ARCHIVE_URL}"

# Create temp directory
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

# Download checksums
echo "Downloading checksums..."
if ! curl -fsSL "$CHECKSUMS_URL" -o "$TMPDIR/$CHECKSUMS_NAME"; then
    echo "Error: Failed to download checksums from $CHECKSUMS_URL" >&2
    exit 1
fi

# Download archive
echo "Downloading archive..."
if ! curl -fsSL "$ARCHIVE_URL" -o "$TMPDIR/$ARCHIVE_NAME"; then
    echo "Error: Failed to download archive from $ARCHIVE_URL" >&2
    exit 1
fi

# Verify checksum
echo "Verifying checksum..."
cd "$TMPDIR"
EXPECTED_CHECKSUM=$(grep "$ARCHIVE_NAME" "$CHECKSUMS_NAME" | awk '{print $1}')
if [[ -z "$EXPECTED_CHECKSUM" ]]; then
    echo "Error: Archive not found in checksums file" >&2
    exit 1
fi

if command -v sha256sum &> /dev/null; then
    ACTUAL_CHECKSUM=$(sha256sum "$ARCHIVE_NAME" | awk '{print $1}')
elif command -v shasum &> /dev/null; then
    ACTUAL_CHECKSUM=$(shasum -a 256 "$ARCHIVE_NAME" | awk '{print $1}')
else
    echo "Error: No sha256sum or shasum command found" >&2
    exit 1
fi

if [[ "$EXPECTED_CHECKSUM" != "$ACTUAL_CHECKSUM" ]]; then
    echo "Error: Checksum verification failed!" >&2
    echo "  Expected: $EXPECTED_CHECKSUM" >&2
    echo "  Actual:   $ACTUAL_CHECKSUM" >&2
    exit 1
fi
echo "Checksum verified."

# Extract and install
echo "Extracting..."
tar -xzf "$ARCHIVE_NAME"

cd - > /dev/null
mkdir -p "$INSTALL_DIR"
mv "$TMPDIR/decider" "$INSTALL_DIR/decider"
chmod +x "$INSTALL_DIR/decider"

echo "Installed to: $INSTALL_DIR/decider"
echo ""
echo "Add to PATH with:"
echo "  export PATH=\"\$PWD/$INSTALL_DIR:\$PATH\""
echo ""
"$INSTALL_DIR/decider" version
