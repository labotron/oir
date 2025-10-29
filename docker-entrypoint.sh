#!/bin/sh
set -e

# Docker entrypoint script for ODT Image Replacer API
# Handles initialization, environment setup, and graceful shutdown

echo "╔══════════════════════════════════════════════════════════╗"
echo "║      ODT Image Replacer API - Docker Container          ║"
echo "╚══════════════════════════════════════════════════════════╝"
echo ""

# Print environment information
echo "Environment Information:"
echo "  Container User: $(whoami)"
echo "  Working Directory: $(pwd)"
echo "  Go Version: $(go version 2>/dev/null || echo 'N/A (production build)')"
echo "  Timezone: ${TZ:-UTC}"
echo ""

# Print configuration
echo "API Configuration:"
echo "  Mode: ${GIN_MODE:-release}"
echo "  Host: ${HOST:-0.0.0.0}"
echo "  Port: ${PORT:-8080}"
echo ""

# Check if binary exists
if [ ! -f "./odt-api" ]; then
    echo "ERROR: odt-api binary not found!"
    exit 1
fi

# Verify binary is executable
if [ ! -x "./odt-api" ]; then
    echo "ERROR: odt-api binary is not executable!"
    exit 1
fi

# Check temp directory
if [ ! -d "/tmp/odt-temp" ]; then
    echo "WARNING: /tmp/odt-temp directory not found, creating..."
    mkdir -p /tmp/odt-temp
fi

# Set temporary directory for Go
export TMPDIR=/tmp/odt-temp

# Trap signals for graceful shutdown
trap 'echo ""; echo "Shutting down gracefully..."; kill -TERM $PID; wait $PID' TERM INT

echo "Starting ODT Image Replacer API..."
echo "════════════════════════════════════════════════════════════"
echo ""

# Start the API server in background
./odt-api "$@" &
PID=$!

# Wait for the process
wait $PID
EXIT_CODE=$?

echo ""
echo "════════════════════════════════════════════════════════════"
echo "API server stopped with exit code: $EXIT_CODE"

exit $EXIT_CODE
