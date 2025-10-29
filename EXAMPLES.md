# ODT Image Replacer - Usage Examples

This document provides practical examples for using the ODT Image Replacer library.

## Basic Usage

### 1. Replace an Image by Tag

```go
package main

import (
    "log"
    "os"

    "github.com/suttapak/odtimagereplacer"
)

func main() {
    // Open the ODT document
    doc, err := odtimagereplacer.NewODTDocument("report.odt")
    if err != nil {
        log.Fatalf("Failed to open ODT: %v", err)
    }

    // Read the new image
    imageData, err := os.ReadFile("new_photo.png")
    if err != nil {
        log.Fatalf("Failed to read image: %v", err)
    }

    // Replace image with tag "image1"
    err = doc.ReplaceImageByTag("image1", "Pictures/new_photo.png", imageData)
    if err != nil {
        log.Fatalf("Failed to replace image: %v", err)
    }

    // Save the modified document
    err = doc.Save("report_updated.odt")
    if err != nil {
        log.Fatalf("Failed to save ODT: %v", err)
    }

    log.Println("Image replaced successfully!")
}
```

### 2. Find All Image Tags

```go
package main

import (
    "fmt"
    "log"

    "github.com/suttapak/odtimagereplacer"
)

func main() {
    // Open the ODT document
    doc, err := odtimagereplacer.NewODTDocument("report.odt")
    if err != nil {
        log.Fatalf("Failed to open ODT: %v", err)
    }

    // Find all image tags
    tags, err := doc.FindImageTags()
    if err != nil {
        log.Fatalf("Failed to find tags: %v", err)
    }

    fmt.Printf("Found %d image tags:\n", len(tags))
    for i, tag := range tags {
        fmt.Printf("%d. %s\n", i+1, tag)
    }
}
```

### 3. Add a New Image

```go
package main

import (
    "log"
    "os"

    "github.com/suttapak/odtimagereplacer"
)

func main() {
    // Open the ODT document
    doc, err := odtimagereplacer.NewODTDocument("report.odt")
    if err != nil {
        log.Fatalf("Failed to open ODT: %v", err)
    }

    // Read the image
    imageData, err := os.ReadFile("logo.png")
    if err != nil {
        log.Fatalf("Failed to read image: %v", err)
    }

    // Add the image to Pictures directory
    err = doc.AddImage("Pictures/logo.png", imageData)
    if err != nil {
        log.Fatalf("Failed to add image: %v", err)
    }

    // Save the document
    err = doc.Save("report_with_logo.odt")
    if err != nil {
        log.Fatalf("Failed to save ODT: %v", err)
    }

    log.Println("Image added successfully!")
}
```

### 4. Replace Multiple Images

```go
package main

import (
    "log"
    "os"

    "github.com/suttapak/odtimagereplacer"
)

func main() {
    // Open the ODT document
    doc, err := odtimagereplacer.NewODTDocument("report.odt")
    if err != nil {
        log.Fatalf("Failed to open ODT: %v", err)
    }

    // Define images to replace
    replacements := map[string]string{
        "image1": "photos/photo1.jpg",
        "image2": "photos/photo2.jpg",
        "logo":   "assets/logo.png",
    }

    // Replace each image
    for tag, imagePath := range replacements {
        imageData, err := os.ReadFile(imagePath)
        if err != nil {
            log.Printf("Warning: Failed to read %s: %v", imagePath, err)
            continue
        }

        // Use the same filename in the ODT
        odtPath := "Pictures/" + filepath.Base(imagePath)
        err = doc.ReplaceImageByTag(tag, odtPath, imageData)
        if err != nil {
            log.Printf("Warning: Failed to replace %s: %v", tag, err)
            continue
        }

        log.Printf("Replaced %s successfully", tag)
    }

    // Save the modified document
    err = doc.Save("report_updated.odt")
    if err != nil {
        log.Fatalf("Failed to save ODT: %v", err)
    }

    log.Println("All images replaced successfully!")
}
```

### 5. Error Handling Example

```go
package main

import (
    "errors"
    "log"
    "os"

    "github.com/suttapak/odtimagereplacer"
)

func main() {
    doc, err := odtimagereplacer.NewODTDocument("report.odt")
    if err != nil {
        // Handle specific errors
        if errors.Is(err, odtimagereplacer.ErrInvalidODT) {
            log.Fatal("File is not a valid ODT document")
        }
        if errors.Is(err, odtimagereplacer.ErrFileTooLarge) {
            log.Fatal("ODT file is too large to process")
        }
        log.Fatalf("Failed to open ODT: %v", err)
    }

    imageData, err := os.ReadFile("photo.png")
    if err != nil {
        log.Fatalf("Failed to read image: %v", err)
    }

    err = doc.ReplaceImageByTag("image1", "Pictures/photo.png", imageData)
    if err != nil {
        if errors.Is(err, odtimagereplacer.ErrImageNotFound) {
            log.Fatal("Image tag 'image1' not found in document")
        }
        log.Fatalf("Failed to replace image: %v", err)
    }

    err = doc.Save("output.odt")
    if err != nil {
        log.Fatalf("Failed to save: %v", err)
    }
}
```

## Command Line Usage

### Build the CLI

```bash
go build -o odt-replacer cmd/cli/main.go
```

### List Image Tags

```bash
./odt-replacer -odt=report.odt -list
```

### Replace an Image

```bash
./odt-replacer -odt=report.odt -tag=image1 -image=photo.png -name=Pictures/photo.png -output=result.odt
```

### Overwrite Original File

```bash
./odt-replacer -odt=report.odt -tag=image1 -image=photo.png -name=Pictures/photo.png
```

## Advanced Usage

### Using with HTTP Handler

```go
package main

import (
    "bytes"
    "io"
    "log"
    "net/http"
    "os"

    "github.com/suttapak/odtimagereplacer"
)

func replaceImageHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse multipart form
    err := r.ParseMultipartForm(10 << 20) // 10 MB max
    if err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    // Get uploaded ODT
    odtFile, _, err := r.FormFile("odt")
    if err != nil {
        http.Error(w, "Missing ODT file", http.StatusBadRequest)
        return
    }
    defer odtFile.Close()

    // Save temporarily
    tmpODT, err := os.CreateTemp("", "upload-*.odt")
    if err != nil {
        http.Error(w, "Internal error", http.StatusInternalServerError)
        return
    }
    defer os.Remove(tmpODT.Name())
    defer tmpODT.Close()

    io.Copy(tmpODT, odtFile)
    tmpODT.Close()

    // Open document
    doc, err := odtimagereplacer.NewODTDocument(tmpODT.Name())
    if err != nil {
        http.Error(w, "Invalid ODT file", http.StatusBadRequest)
        return
    }

    // Get image file
    imageFile, _, err := r.FormFile("image")
    if err != nil {
        http.Error(w, "Missing image file", http.StatusBadRequest)
        return
    }
    defer imageFile.Close()

    imageData, err := io.ReadAll(imageFile)
    if err != nil {
        http.Error(w, "Failed to read image", http.StatusBadRequest)
        return
    }

    // Get tag from form
    tag := r.FormValue("tag")
    if tag == "" {
        http.Error(w, "Missing tag parameter", http.StatusBadRequest)
        return
    }

    // Replace image
    err = doc.ReplaceImageByTag(tag, "Pictures/uploaded.png", imageData)
    if err != nil {
        http.Error(w, "Failed to replace image: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Save to buffer
    outputPath := tmpODT.Name() + ".out"
    err = doc.Save(outputPath)
    if err != nil {
        http.Error(w, "Failed to save document", http.StatusInternalServerError)
        return
    }
    defer os.Remove(outputPath)

    // Send back the modified ODT
    outputData, err := os.ReadFile(outputPath)
    if err != nil {
        http.Error(w, "Failed to read output", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/vnd.oasis.opendocument.text")
    w.Header().Set("Content-Disposition", "attachment; filename=modified.odt")
    w.Write(outputData)
}

func main() {
    http.HandleFunc("/replace", replaceImageHandler)
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Tips and Best Practices

1. **Always check errors**: ODT manipulation can fail for many reasons (corrupted files, missing tags, etc.)

2. **Tag naming in ODT**: When creating ODT templates, give your images meaningful `draw:name` attributes in LibreOffice/OpenOffice

3. **Image paths**: Images should typically be in the `Pictures/` directory within the ODT

4. **File size limits**: The library has built-in protection against zip bombs and oversized files:
   - Maximum ODT size: 100MB
   - Maximum individual file size: 50MB
   - Maximum files in archive: 10,000

5. **MIME type detection**: The library automatically detects MIME types from file extensions (png, jpg, gif, svg, webp)

6. **Backward compatibility**: Legacy functions `AddImage()` and `Test()` are still available but deprecated

7. **Security**: The library includes protection against path traversal attacks and validates all file paths

8. **Performance**: Regex patterns are compiled once and cached for better performance
