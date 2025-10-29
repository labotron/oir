# ODT Image Replacer

A secure, high-performance Go library for manipulating images in ODT (Open Document Text) files.

## Features

- Replace images by tag name in ODT documents
- Add new images to existing ODT files
- List all image tags in a document
- Security-hardened against path traversal and zip bomb attacks
- Production-ready with comprehensive error handling
- Zero external dependencies
- Backward compatible API

## Installation

```bash
go get github.com/suttapak/odtimagereplacer
```

## Quick Start

```go
package main

import (
    "log"
    "os"
    "github.com/suttapak/odtimagereplacer"
)

func main() {
    // Open ODT document
    doc, err := odtimagereplacer.NewODTDocument("report.odt")
    if err != nil {
        log.Fatal(err)
    }

    // Read new image
    imageData, err := os.ReadFile("photo.png")
    if err != nil {
        log.Fatal(err)
    }

    // Replace image by tag
    err = doc.ReplaceImageByTag("image1", "Pictures/photo.png", imageData)
    if err != nil {
        log.Fatal(err)
    }

    // Save document
    err = doc.Save("report_updated.odt")
    if err != nil {
        log.Fatal(err)
    }
}
```

## CLI Usage

Build the CLI tool:

```bash
go build -o odt-replacer cmd/cli/main.go
```

List all image tags in an ODT:

```bash
./odt-replacer -odt=report.odt -list
```

Replace an image:

```bash
./odt-replacer -odt=report.odt -tag=image1 -image=photo.png -name=Pictures/photo.png -output=result.odt
```

## REST API Usage

Build and run the API server:

```bash
# Build the API
go build -o odt-api cmd/api/main.go

# Run the server
./odt-api

# Or with custom settings
./odt-api -port=3000 -mode=debug
```

Replace images via HTTP API:

```bash
curl -X POST http://localhost:8080/api/replace \
  -H "Content-Type: application/json" \
  -d '{
    "template": {
      "url": "https://example.com/template.odt",
      "base64": null
    },
    "data": {
      "image1": {
        "url": "https://example.com/photo.png",
        "base64": null
      }
    }
  }'
```

Download ODT directly:

```bash
curl -X POST http://localhost:8080/api/replace/download \
  -H "Content-Type: application/json" \
  -d @request.json \
  -o output.odt
```

See [API.md](API.md) for complete API documentation and [API_EXAMPLES.sh](API_EXAMPLES.sh) for more examples.

## API Documentation

### Core Functions

#### `NewODTDocument(path string) (*ODTDocument, error)`
Opens and validates an ODT file.

#### `(*ODTDocument) ReplaceImageByTag(tag, imagePath string, imageData []byte) error`
Replaces an image identified by its `draw:name` tag in the ODT.

#### `(*ODTDocument) AddImage(imagePath string, imageData []byte) error`
Adds a new image to the ODT at the specified path.

#### `(*ODTDocument) FindImageTags() ([]string, error)`
Returns a list of all image tags (draw:name attributes) in the document.

#### `(*ODTDocument) Save(outputPath string) error`
Saves the modified ODT to disk.

## Security Features

- **Path Traversal Protection**: All file paths are validated
- **Zip Bomb Protection**:
  - Maximum ODT size: 100MB
  - Maximum individual file: 50MB
  - Maximum files in archive: 10,000
- **Input Validation**: All user inputs are sanitized
- **No Silent Failures**: All errors are properly propagated

## Performance

- **Lazy Loading**: Files are loaded only when needed
- **Pre-compiled Regex**: Patterns compiled once and cached
- **Minimal Memory Copying**: Efficient byte operations
- **Test Coverage**: 67.9%

## Examples

See [EXAMPLES.md](EXAMPLES.md) for comprehensive usage examples including:
- Basic image replacement
- Finding image tags
- Adding new images
- Batch processing
- HTTP handler integration
- Error handling patterns

## Documentation

- [EXAMPLES.md](EXAMPLES.md) - Comprehensive usage examples
- [CLAUDE.md](CLAUDE.md) - Development guide and architecture
- [REFACTORING.md](REFACTORING.md) - Refactoring summary and improvements

## Testing

Run tests:

```bash
go test ./...                    # All tests
go test -v ./...                 # Verbose
go test -cover ./...             # With coverage
go test -bench=. ./...           # Benchmarks
```

## How It Works

ODT files are ZIP archives containing XML files. This library:

1. Opens the ODT as a ZIP archive with security validation
2. Parses `content.xml` to find image references
3. Modifies the `xlink:href` attribute in `<draw:frame>` elements
4. Updates `META-INF/manifest.xml` with new image entries
5. Re-packages everything back into a valid ODT file

Image tags are defined in the ODT template using the `draw:name` attribute:

```xml
<draw:frame draw:name="image1" ...>
  <draw:image xlink:href="Pictures/photo.png" />
</draw:frame>
```

## Requirements

- Go 1.25.2 or higher
- No external dependencies

## License

[Add your license here]

## Contributing

Contributions are welcome! Please ensure:
- All tests pass
- Code follows Go conventions
- Security best practices are maintained
- Documentation is updated

## Changelog

### Version 2.0.0 (2024)
- Complete refactoring with focus on security and performance
- New `ODTDocument` API with method-based interface
- Comprehensive test suite with 67.9% coverage
- Security hardening against attacks
- Performance optimizations with lazy loading
- Enhanced CLI with flag support
- Extensive documentation

### Version 1.0.0 (Original)
- Basic image replacement functionality
- Simple function-based API

## Credits

Refactored and enhanced with senior-level development practices focusing on security, performance, and maintainability.
