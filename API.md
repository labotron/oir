# ODT Image Replacer API Documentation

A RESTful API built with Gin framework for replacing images in ODT documents using JSON requests.

## Quick Start

### Build and Run

```bash
# Build the API server
go build -o odt-api cmd/api/main.go

# Run with default settings (port 8080)
./odt-api

# Run with custom settings
./odt-api -port=3000 -host=localhost -mode=debug
```

### Command Line Options

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | `8080` | Server port |
| `-host` | `0.0.0.0` | Server host |
| `-mode` | `release` | Gin mode: `debug`, `release`, or `test` |

## API Endpoints

### 1. Health Check

Check if the API is running.

**Endpoint:** `GET /health`

**Response:**
```json
{
  "status": "healthy",
  "service": "odt-image-replacer"
}
```

**Example:**
```bash
curl http://localhost:8080/health
```

---

### 2. Service Information

Get information about the API service.

**Endpoint:** `GET /info`

**Response:**
```json
{
  "service": "ODT Image Replacer API",
  "version": "2.0.0",
  "description": "Replace images in ODT documents via JSON API",
  "endpoints": {
    "POST /api/replace": "Replace images and return JSON with base64 output",
    "POST /api/replace/download": "Replace images and download ODT file directly",
    "GET  /health": "Health check endpoint",
    "GET  /info": "Service information"
  }
}
```

**Example:**
```bash
curl http://localhost:8080/info
```

---

### 3. Replace Images (JSON Response)

Replace images in an ODT document and receive the result as base64-encoded data in JSON.

**Endpoint:** `POST /api/replace`

**Request Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "template": {
    "url": "https://example.com/template.odt",
    "base64": null
  },
  "data": {
    "image1": {
      "url": "https://example.com/photo1.png",
      "base64": null
    },
    "image2": {
      "url": null,
      "base64": "iVBORw0KGgoAAAANSUhEUgAAAAUA..."
    }
  }
}
```

**Request Body Fields:**

- `template.url` (string): URL to download the ODT template
- `template.base64` (string): Base64-encoded ODT template (use if URL is null)
- `data` (object): Map of image tag names to image sources
  - Each key is the `draw:name` tag in the ODT
  - Each value has `url` or `base64` for the image source

**Response (Success):**
```json
{
  "success": true,
  "message": "Successfully replaced 2 image(s)",
  "output_base64": "UEsDBBQAAAAIAOB/...",
  "replaced_tags": ["image1", "image2"]
}
```

**Response (Error):**
```json
{
  "success": false,
  "error": "failed to get template: invalid URL"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/replace \
  -H "Content-Type: application/json" \
  -d @example.json
```

---

### 4. Replace Images (File Download)

Replace images in an ODT document and download the result directly as an ODT file.

**Endpoint:** `POST /api/replace/download`

**Request:** Same as `/api/replace`

**Response:** Binary ODT file with headers:
```
Content-Type: application/vnd.oasis.opendocument.text
Content-Disposition: attachment; filename=output.odt
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/replace/download \
  -H "Content-Type: application/json" \
  -d @example.json \
  -o output.odt
```

---

## Request Format Details

### Template Source

The `template` object specifies where to get the ODT template:

**Option 1: From URL**
```json
{
  "template": {
    "url": "https://example.com/template.odt",
    "base64": null
  }
}
```

**Option 2: From Base64**
```json
{
  "template": {
    "url": null,
    "base64": "UEsDBBQAAAAIAOB/..."
  }
}
```

**Note:** If both are provided, URL takes precedence.

### Image Sources

The `data` object maps ODT image tags to image sources:

**Mixed Sources Example:**
```json
{
  "data": {
    "logo": {
      "url": "https://example.com/logo.png",
      "base64": null
    },
    "signature": {
      "url": null,
      "base64": "iVBORw0KGgoAAAANSUhEUgAAAAUA..."
    }
  }
}
```

**Tag Matching:**
- Tag names (e.g., "logo", "signature") must match the `draw:name` attribute in the ODT
- In LibreOffice, right-click an image → Properties → Options → Name

---

## Complete Example

### 1. Create Example JSON File

Save as `request.json`:

```json
{
  "template": {
    "url": null,
    "base64": "UEsDBBQAAAAIAOB/Y1n5H7IWAAAA..."
  },
  "data": {
    "image1": {
      "url": "https://picsum.photos/200/300",
      "base64": null
    },
    "image2": {
      "url": "https://picsum.photos/300/200",
      "base64": null
    }
  }
}
```

### 2. Send Request

```bash
# Get JSON response with base64 output
curl -X POST http://localhost:8080/api/replace \
  -H "Content-Type: application/json" \
  -d @request.json

# Download ODT file directly
curl -X POST http://localhost:8080/api/replace/download \
  -H "Content-Type: application/json" \
  -d @request.json \
  -o output.odt
```

### 3. Decode Base64 Output (if using /api/replace)

```bash
# Using jq and base64 command
curl -X POST http://localhost:8080/api/replace \
  -H "Content-Type: application/json" \
  -d @request.json | \
  jq -r '.output_base64' | \
  base64 -d > output.odt
```

---

## Image Format Support

The API automatically detects image formats based on magic bytes:

- **PNG** (.png)
- **JPEG** (.jpg, .jpeg)
- **GIF** (.gif)
- **WebP** (.webp)

Images are stored in the ODT's `Pictures/` directory with the appropriate extension.

---

## Security Features

### File Size Limits

- Maximum ODT template size: **100MB**
- Maximum individual image size: **50MB**
- Maximum files in ODT archive: **10,000**

### Protected Against

- Path traversal attacks
- Zip bomb attacks
- Memory exhaustion
- Invalid file formats

### Timeouts

- HTTP request timeout: **30 seconds**

---

## Error Handling

### Common Error Responses

**Invalid JSON:**
```json
{
  "success": false,
  "error": "Invalid JSON: unexpected end of JSON input"
}
```

**Missing Template:**
```json
{
  "success": false,
  "error": "no valid template source provided (URL or base64)"
}
```

**Invalid Image Tag:**
```json
{
  "success": false,
  "message": "Successfully replaced 1 image(s)",
  "replaced_tags": ["image1"]
}
```

**File Too Large:**
```json
{
  "success": false,
  "error": "file size exceeds maximum allowed limit: 52428800 bytes (max: 52428800)"
}
```

---

## Docker Deployment (Optional)

### Dockerfile

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o odt-api cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/odt-api .
EXPOSE 8080
CMD ["./odt-api"]
```

### Build and Run

```bash
docker build -t odt-api .
docker run -p 8080:8080 odt-api
```

---

## Testing the API

### Using curl

```bash
# Health check
curl http://localhost:8080/health

# Service info
curl http://localhost:8080/info

# Replace images
curl -X POST http://localhost:8080/api/replace \
  -H "Content-Type: application/json" \
  -d '{
    "template": {"url": "https://example.com/template.odt", "base64": null},
    "data": {
      "image1": {"url": "https://example.com/photo.png", "base64": null}
    }
  }'
```

### Using Postman

1. Create a new POST request
2. URL: `http://localhost:8080/api/replace`
3. Headers: `Content-Type: application/json`
4. Body: Select "raw" and paste JSON
5. Send

### Using JavaScript/Node.js

```javascript
const axios = require('axios');
const fs = require('fs');

async function replaceImages() {
  const response = await axios.post('http://localhost:8080/api/replace', {
    template: {
      url: 'https://example.com/template.odt',
      base64: null
    },
    data: {
      image1: {
        url: 'https://example.com/photo.png',
        base64: null
      }
    }
  });

  // Save the result
  const buffer = Buffer.from(response.data.output_base64, 'base64');
  fs.writeFileSync('output.odt', buffer);

  console.log('Replaced tags:', response.data.replaced_tags);
}

replaceImages();
```

### Using Python

```python
import requests
import base64

response = requests.post('http://localhost:8080/api/replace', json={
    'template': {
        'url': 'https://example.com/template.odt',
        'base64': None
    },
    'data': {
        'image1': {
            'url': 'https://example.com/photo.png',
            'base64': None
        }
    }
})

result = response.json()
if result['success']:
    # Decode and save the ODT
    odt_data = base64.b64decode(result['output_base64'])
    with open('output.odt', 'wb') as f:
        f.write(odt_data)
    print(f"Replaced tags: {result['replaced_tags']}")
```

---

## Performance Considerations

- The API processes requests synchronously
- Large images or many replacements may take longer
- For high-volume production use, consider:
  - Load balancing with multiple instances
  - Request queuing system
  - Async processing with job queue

---

## Monitoring

### Logging

The API uses Gin's default logger which logs:
- Request method and path
- Response status code
- Response time
- Client IP

**Example log output:**
```
[GIN] 2024/10/29 - 09:21:45 | 200 |  1.234567ms |  127.0.0.1 | POST  "/api/replace"
```

### Custom Logging

For production, consider adding structured logging:

```go
import "github.com/sirupsen/logrus"

// Add custom middleware
router.Use(LoggingMiddleware())
```

---

## Troubleshooting

### Issue: "connection refused"
**Solution:** Ensure the server is running and the port is not blocked by firewall

### Issue: "file too large"
**Solution:** Reduce image sizes or compress images before sending

### Issue: "image tag not found"
**Solution:** Verify the tag name matches exactly with `draw:name` in the ODT template

### Issue: "invalid ODT file"
**Solution:** Ensure the template is a valid ODT file created by LibreOffice/OpenOffice

---

## Security Best Practices

1. **Use HTTPS in production** - Add TLS/SSL certificates
2. **Add authentication** - Use API keys or JWT tokens
3. **Rate limiting** - Prevent abuse with rate limits
4. **Input validation** - Already implemented in the API
5. **CORS configuration** - Configure allowed origins

---

## Support

For issues or questions:
- Check the logs for error messages
- Verify your JSON request format matches the examples
- Ensure image URLs are accessible
- Test with the `/health` endpoint first
