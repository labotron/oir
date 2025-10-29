#!/bin/bash

# Simple script to unzip ODT and show content.xml
# Usage: ./unzip_odt.sh <odt-file> [output-dir]

if [ $# -eq 0 ]; then
    echo "Usage: $0 <odt-file> [output-dir]"
    echo ""
    echo "Examples:"
    echo "  $0 template.odt              # Extract to temp directory"
    echo "  $0 template.odt extracted/   # Extract to 'extracted/' directory"
    exit 1
fi

ODT_FILE="$1"
OUTPUT_DIR="${2:-$(mktemp -d)}"

if [ ! -f "$ODT_FILE" ]; then
    echo "❌ File not found: $ODT_FILE"
    exit 1
fi

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

echo "📦 Unzipping: $ODT_FILE"
echo "📁 Output: $OUTPUT_DIR"
echo ""

# Unzip
unzip -q "$ODT_FILE" -d "$OUTPUT_DIR"

echo "✅ Extracted successfully!"
echo ""
echo "📄 Files:"
ls -la "$OUTPUT_DIR"
echo ""

# Show content.xml with line numbers
if [ -f "$OUTPUT_DIR/content.xml" ]; then
    echo "================================================"
    echo "📝 content.xml (with line numbers):"
    echo "================================================"
    cat -n "$OUTPUT_DIR/content.xml"
    echo ""
fi

# Search for placeholders
echo "================================================"
echo "🔍 Placeholders found:"
echo "================================================"
grep -o '{[^}]*}' "$OUTPUT_DIR/content.xml" 2>/dev/null || echo "None"
echo ""

# Show svg:title elements
echo "================================================"
echo "🏷️  All svg:title elements:"
echo "================================================"
grep -o '<svg:title>[^<]*</svg:title>' "$OUTPUT_DIR/content.xml" 2>/dev/null || echo "None"
echo ""

# Show image count
IMAGE_COUNT=$(grep -c "draw:image" "$OUTPUT_DIR/content.xml" 2>/dev/null || echo "0")
echo "================================================"
echo "🖼️  Images found: $IMAGE_COUNT"
echo "================================================"
echo ""

echo "💡 To view files:"
echo "   cat $OUTPUT_DIR/content.xml"
echo "   cat $OUTPUT_DIR/styles.xml"
echo "   cat $OUTPUT_DIR/META-INF/manifest.xml"
echo ""
echo "💡 To format XML nicely (if xmllint installed):"
echo "   xmllint --format $OUTPUT_DIR/content.xml"
