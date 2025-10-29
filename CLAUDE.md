# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go library for manipulating ODT (Open Document Text) files, specifically for replacing images within ODT documents. ODT files are ZIP archives containing XML files that define document structure and content.

**Status**: Recently refactored (2024) with focus on security, performance, and production-ready code.

## Architecture

### Core Functionality (Refactored)

The library now provides a clean, secure API through the `ODTDocument` type:

```go
// Open document
doc, err := NewODTDocument("report.odt")

// Replace image by tag
err = doc.ReplaceImageByTag("image1", "Pictures/photo.png", imageData)

// Save changes
err = doc.Save("output.odt")
```

The library works by:
1. Reading ODT files as ZIP archives with security validation
2. Lazy-loading files only when needed (performance optimization)
3. Parsing and modifying XML content using pre-compiled regex patterns
4. Managing manifest.xml automatically with proper MIME type detection
5. Re-zipping the modified content securely

### Key Files

**New Architecture (Use These):**
- `odt.go`: Main API with `ODTDocument` type and secure methods:
  - `NewODTDocument()`: Opens and validates ODT files
  - `ReplaceImageByTag()`: Replaces images by draw:name tag
  - `AddImage()`: Adds new images to ODT
  - `FindImageTags()`: Lists all image tags in document
  - `Save()`: Writes modified ODT to disk
- `errors.go`: Custom error types for better error handling
- `odt_test.go`: Comprehensive unit tests and benchmarks

**Legacy Files (Deprecated but functional):**
- `unzip_odt.go`: Legacy `AddImage()` and `UnzipMem()` - now wrappers around new API
- `reader.go`: Legacy `Test()` function - now uses new API internally
- `draw_frame.go`: Documentation example of ODT XML structure

**CLI:**
- `cmd/cli/main.go`: Feature-rich CLI with flags for listing tags and replacing images

### ODT Structure

ODT files contain:
- `content.xml`: Document content with `<draw:frame>` elements for images
- `META-INF/manifest.xml`: File manifest listing all resources with MIME types
- `Pictures/`: Directory containing image files
- Other files: styles.xml, settings.xml, meta.xml, Fonts/

Images are referenced in content.xml:
```xml
<draw:frame draw:name="image1" ...>
  <draw:image xlink:href="Pictures/filename.png" />
  <svg:title>{img1}</svg:title>
</draw:frame>
```

The `draw:name` attribute is the tag used for targeting specific images.

## Development Commands

### Run tests
```bash
go test ./...                    # All tests
go test -v ./...                 # Verbose output
go test -bench=. ./...           # Run benchmarks
go test -cover ./...             # With coverage
```

### Build the CLI
```bash
go build -o odt-replacer cmd/cli/main.go
```

### CLI Usage
```bash
# List all image tags in an ODT
./odt-replacer -odt=report.odt -list

# Replace an image
./odt-replacer -odt=report.odt -tag=image1 -image=photo.png -name=Pictures/photo.png -output=result.odt

# Legacy test mode (no flags)
./odt-replacer
```

### Inspect ODT Files
```bash
./unzip_odt.sh <odt-file> [output-dir]
```

Displays file listing, content.xml, placeholders, and image count.

### Module Path
`github.com/suttapak/odtimagereplacer` (Go 1.25.2)

## Security Features (New)

The refactored code includes comprehensive security measures:

1. **Path Traversal Protection**: All file paths validated to prevent `../` attacks
2. **Zip Bomb Protection**:
   - Maximum ODT size: 100MB
   - Maximum individual file size: 50MB
   - Maximum files in archive: 10,000
3. **Input Validation**: Image names, tags, and data validated before processing
4. **Secure File Permissions**: Output files created with 0644 permissions
5. **No Silent Failures**: All errors properly propagated (no ignored errors with `_`)

## Performance Optimizations (New)

1. **Lazy Loading**: Files only loaded from ZIP when accessed, not all upfront
2. **Pre-compiled Regex**: Patterns compiled once and cached with `sync.Once`
3. **Efficient Memory Usage**: Streams data where possible, bounded buffers
4. **Minimal Copying**: Direct byte slice operations, reuse of buffers

## API Design Patterns

1. **Constructor Pattern**: `NewODTDocument()` validates and returns ready-to-use document
2. **Error Wrapping**: All errors wrapped with context using `fmt.Errorf("%w", err)`
3. **Custom Error Types**: Use `errors.Is()` to check for specific error conditions
4. **Method Chaining Ready**: All methods return errors, allowing for future fluent API
5. **Backward Compatibility**: Legacy functions still work via internal delegation

## Testing Strategy

- Unit tests cover all public functions and edge cases
- Security tests for path traversal, oversized files, invalid inputs
- Benchmarks for performance-critical operations
- Helper functions to create test ODT files in memory
- Table-driven tests for comprehensive coverage

## Migration Guide (for existing code)

**Old code:**
```go
// Old way (still works but deprecated)
AddImage("report.odt", "photo.png", imageData)
```

**New code:**
```go
// New way (recommended)
doc, err := NewODTDocument("report.odt")
if err != nil { /* handle error */ }
err = doc.AddImage("Pictures/photo.png", imageData)
if err != nil { /* handle error */ }
err = doc.Save("report.odt")
```

See `EXAMPLES.md` for comprehensive usage examples including HTTP handlers and batch processing.

## Common Pitfalls

1. **Image paths**: Must include `Pictures/` prefix (e.g., `Pictures/photo.png`)
2. **Tag names**: Must exactly match the `draw:name` attribute in content.xml
3. **MIME types**: Automatically detected from file extension, ensure correct extensions
4. **File size limits**: Images over 50MB will be rejected
5. **Concurrent access**: `ODTDocument` is not thread-safe, use separate instances per goroutine

## Future Enhancements

Potential areas for improvement:
- Stream-based processing for very large ODT files
- Support for other ODF formats (ODS, ODP)
- Template variable replacement beyond images
- Concurrent safe document access with locking
