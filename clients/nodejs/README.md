# ODT Image Replacer - Node.js Client

A simple Node.js client library for the ODT Image Replacer API.

## Installation

```bash
npm install
```

Or if you're using this as a standalone package:

```bash
npm install axios
```

## Quick Start

### Simple Usage (One Function)

```javascript
const { replaceImages } = require('./index');

// Replace images in one line
await replaceImages(
  './template.odt',           // Template path
  {
    image1: './photo1.png',   // Images by tag name
    image2: './photo2.jpg'
  },
  './output.odt'              // Output path
);
```

## API Reference

### Class: `ODTImageReplacerClient`

Main client class for interacting with the API.

#### Constructor

```javascript
const client = new ODTImageReplacerClient(baseURL, options);
```

**Parameters:**
- `baseURL` (string): API base URL (default: `'http://localhost:8080'`)
- `options` (object):
  - `timeout` (number): Request timeout in milliseconds (default: `30000`)

#### Methods

##### `replaceImages(params, options)`

Replace images in an ODT document.

```javascript
const buffer = await client.replaceImages({
  template: {
    filePath: './template.odt'  // or url, or base64
  },
  data: {
    image1: {
      filePath: './photo.png'   // or url, or base64
    }
  }
}, {
  returnBase64: false  // Set true to get base64 string instead of Buffer
});
```

**Returns:** `Promise<Buffer|string>`

##### `replaceImagesAndSave(params, outputPath)`

Replace images and save directly to a file.

```javascript
await client.replaceImagesAndSave({
  template: { filePath: './template.odt' },
  data: {
    image1: { filePath: './photo.png' }
  }
}, './output.odt');
```

**Returns:** `Promise<void>`

##### `downloadModifiedODT(params, outputPath)`

Use the download endpoint to save directly.

```javascript
await client.downloadModifiedODT({
  template: { filePath: './template.odt' },
  data: {
    image1: { filePath: './photo.png' }
  }
}, './output.odt');
```

**Returns:** `Promise<void>`

##### `healthCheck()`

Check if the API is running.

```javascript
const health = await client.healthCheck();
// { status: 'healthy', service: 'odt-image-replacer' }
```

**Returns:** `Promise<Object>`

##### `getInfo()`

Get API information.

```javascript
const info = await client.getInfo();
// { service: '...', version: '...', endpoints: {...} }
```

**Returns:** `Promise<Object>`

---

### Function: `replaceImages()`

Convenience function for quick usage without creating a client instance.

```javascript
await replaceImages(templatePath, images, outputPath, options);
```

**Parameters:**
- `templatePath` (string): Path to ODT template file
- `images` (object): Map of tag names to image file paths
- `outputPath` (string): Path where to save output
- `options` (object):
  - `apiUrl` (string): API base URL (default: `'http://localhost:8080'`)

**Returns:** `Promise<void>`

---

## Usage Examples

### Example 1: Local Files

```javascript
const { replaceImages } = require('odt-image-replacer-client');

await replaceImages(
  './template.odt',
  {
    logo: './logo.png',
    signature: './signature.jpg'
  },
  './output.odt'
);
```

### Example 2: Using URLs

```javascript
const { ODTImageReplacerClient } = require('odt-image-replacer-client');

const client = new ODTImageReplacerClient('http://localhost:8080');

await client.replaceImagesAndSave({
  template: {
    url: 'https://example.com/template.odt'
  },
  data: {
    logo: {
      url: 'https://example.com/logo.png'
    }
  }
}, './output.odt');
```

### Example 3: Mixed Sources

```javascript
const client = new ODTImageReplacerClient();

await client.replaceImagesAndSave({
  template: {
    filePath: './template.odt'  // Local file
  },
  data: {
    logo: {
      url: 'https://example.com/logo.png'  // From URL
    },
    photo: {
      filePath: './photo.jpg'  // Local file
    },
    signature: {
      base64: 'iVBORw0KGgo...'  // Base64 encoded
    }
  }
}, './output.odt');
```

### Example 4: Get Buffer for Further Processing

```javascript
const client = new ODTImageReplacerClient();

const buffer = await client.replaceImages({
  template: { filePath: './template.odt' },
  data: {
    image1: { filePath: './photo.png' }
  }
});

// Now you can upload to S3, send via HTTP, etc.
console.log(`Generated ${buffer.length} bytes`);
```

### Example 5: Get Base64 String

```javascript
const client = new ODTImageReplacerClient();

const base64 = await client.replaceImages({
  template: { filePath: './template.odt' },
  data: {
    image1: { filePath: './photo.png' }
  }
}, { returnBase64: true });

// Send base64 to frontend, store in database, etc.
console.log(base64.substring(0, 50) + '...');
```

### Example 6: Error Handling

```javascript
const { replaceImages } = require('odt-image-replacer-client');

try {
  await replaceImages(
    './template.odt',
    { image1: './photo.png' },
    './output.odt'
  );
  console.log('Success!');
} catch (error) {
  console.error('Failed:', error.message);
}
```

### Example 7: Custom Configuration

```javascript
const client = new ODTImageReplacerClient('http://production:3000', {
  timeout: 60000  // 60 seconds for large files
});

await client.replaceImagesAndSave({
  template: { filePath: './large-template.odt' },
  data: {
    image1: { filePath: './high-res-photo.png' }
  }
}, './output.odt');
```

### Example 8: Using Buffers

```javascript
const client = new ODTImageReplacerClient();
const fs = require('fs');

// Read files into Buffers
const templateBuffer = fs.readFileSync('./template.odt');
const imageBuffer = fs.readFileSync('./photo.png');

// Use Buffers directly - auto-converted to base64
const outputBuffer = await client.replaceImages({
  template: {
    buffer: templateBuffer  // Buffer auto-converts to base64
  },
  data: {
    image1: {
      buffer: imageBuffer  // Buffer auto-converts to base64
    }
  }
});

// Can also save directly
fs.writeFileSync('./output.odt', outputBuffer);
```

### Example 9: Batch Processing

```javascript
const client = new ODTImageReplacerClient();

const employees = ['john', 'jane', 'bob'];

for (const employee of employees) {
  await client.replaceImagesAndSave({
    template: { filePath: './employee-template.odt' },
    data: {
      photo: { filePath: `./photos/${employee}.jpg` },
      signature: { filePath: `./signatures/${employee}.png` }
    }
  }, `./output/report-${employee}.odt`);

  console.log(`Generated report for ${employee}`);
}
```

## TypeScript Support

The package includes TypeScript definitions. Import types as needed:

```typescript
import {
  ODTImageReplacerClient,
  ReplaceImagesParams,
  ImageSource
} from 'odt-image-replacer-client';

const client = new ODTImageReplacerClient('http://localhost:8080');

const params: ReplaceImagesParams = {
  template: { filePath: './template.odt' },
  data: {
    image1: { filePath: './photo.png' }
  }
};

await client.replaceImagesAndSave(params, './output.odt');
```

## Image Source Options

Each image source can be provided in **four** ways:

### 1. Buffer (Most Flexible)
```javascript
const fs = require('fs');
const imageBuffer = fs.readFileSync('./photo.png');

{
  image1: { buffer: imageBuffer }  // Auto-converts to base64
}
```

### 2. File Path (Local File)
```javascript
{
  image1: { filePath: './photo.png' }
}
```

### 3. URL (Download from Web)
```javascript
{
  image1: { url: 'https://example.com/photo.png' }
}
```

### 4. Base64 (Encoded String)
```javascript
{
  image1: { base64: 'iVBORw0KGgoAAAANSUhEUgAAAAUA...' }
}
```

The same options apply to the template source.

**Note:** When using `buffer`, it will be automatically converted to base64 before sending to the API.

## Requirements

- Node.js 12 or higher
- Running ODT Image Replacer API server
- axios package (installed automatically)

## Starting the API Server

Before using the client, start the API server:

```bash
# Build the API
go build -o odt-api cmd/api/main.go

# Run the API
./odt-api

# Or with custom settings
./odt-api -port=3000 -mode=debug
```

The API will start on `http://localhost:8080` by default.

## Error Handling

The client throws descriptive errors:

```javascript
try {
  await replaceImages(...);
} catch (error) {
  // Possible errors:
  // - "Template source must be provided (url, base64, or filePath)"
  // - "Image source for tag 'image1' must be provided"
  // - "failed to get template: invalid URL"
  // - "no images to replace"
  // - Network errors from axios
  console.error(error.message);
}
```

## Testing

Run the examples:

```bash
node example.js
```

Or run individual examples:

```bash
node -e "require('./example').example6()"  # Health check
```

## License

MIT

## Support

For issues or questions about the API, see the main project documentation.
