package odtimagereplacer

import (
	"archive/zip"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidatePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"valid path", "/tmp/test.odt", false},
		{"valid relative path", "test.odt", false},
		{"empty path", "", true},
		{"path traversal", "../../../etc/passwd", true},
		{"path traversal in middle", "/tmp/../../../etc/passwd", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateImageName(t *testing.T) {
	tests := []struct {
		name    string
		imgName string
		wantErr bool
	}{
		{"valid name", "image1.png", false},
		{"valid name with numbers", "img123.jpg", false},
		{"empty name", "", true},
		{"path traversal", "../image.png", true},
		{"with slash", "path/image.png", true},
		{"with backslash", "path\\image.png", true},
		{"too long", strings.Repeat("a", 256) + ".png", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateImageName(tt.imgName)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateImageName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDetectMIMEType(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"image.png", "image/png"},
		{"image.PNG", "image/png"},
		{"photo.jpg", "image/jpeg"},
		{"photo.jpeg", "image/jpeg"},
		{"animation.gif", "image/gif"},
		{"vector.svg", "image/svg+xml"},
		{"modern.webp", "image/webp"},
		{"unknown.xyz", "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := detectMIMEType(tt.path)
			if got != tt.expected {
				t.Errorf("detectMIMEType(%q) = %q, want %q", tt.path, got, tt.expected)
			}
		})
	}
}

func TestNewODTDocument_InvalidFile(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr error
	}{
		{"non-existent file", "/tmp/nonexistent.odt", nil},
		{"empty path", "", ErrInvalidPath},
		{"path traversal", "../../../etc/passwd", ErrInvalidPath},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewODTDocument(tt.path)
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("NewODTDocument() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				t.Errorf("NewODTDocument() expected error but got none")
			}
		})
	}
}

func TestODTDocument_AddImage(t *testing.T) {
	// Create a minimal valid ODT file for testing
	tmpDir := t.TempDir()
	testODT := filepath.Join(tmpDir, "test.odt")

	// Create a minimal ODT structure
	if err := createMinimalODT(testODT); err != nil {
		t.Fatalf("Failed to create test ODT: %v", err)
	}

	doc, err := NewODTDocument(testODT)
	if err != nil {
		t.Fatalf("NewODTDocument() error = %v", err)
	}

	// Test adding an image
	imageData := []byte("fake image data")
	err = doc.AddImage("Pictures/test.png", imageData)
	if err != nil {
		t.Errorf("AddImage() error = %v", err)
	}

	// Verify image was added
	if data, ok := doc.files["Pictures/test.png"]; !ok {
		t.Error("Image not found in files map")
	} else if !bytes.Equal(data, imageData) {
		t.Error("Image data mismatch")
	}

	// Test adding image with invalid data
	err = doc.AddImage("Pictures/empty.png", []byte{})
	if err == nil {
		t.Error("AddImage() should fail with empty data")
	}

	// Test adding image that's too large
	largeData := make([]byte, MaxIndividualFileSize+1)
	err = doc.AddImage("Pictures/large.png", largeData)
	if err == nil {
		t.Error("AddImage() should fail with oversized data")
	}
}

func TestODTDocument_FindImageTags(t *testing.T) {
	tmpDir := t.TempDir()
	testODT := filepath.Join(tmpDir, "test.odt")

	// Create ODT with test content
	if err := createODTWithContent(testODT, testContentXML); err != nil {
		t.Fatalf("Failed to create test ODT: %v", err)
	}

	doc, err := NewODTDocument(testODT)
	if err != nil {
		t.Fatalf("NewODTDocument() error = %v", err)
	}

	tags, err := doc.FindImageTags()
	if err != nil {
		t.Errorf("FindImageTags() error = %v", err)
	}

	expectedTags := []string{"image1", "photo1"}
	if len(tags) != len(expectedTags) {
		t.Errorf("FindImageTags() found %d tags, want %d", len(tags), len(expectedTags))
	}

	// Check that expected tags are present
	tagMap := make(map[string]bool)
	for _, tag := range tags {
		tagMap[tag] = true
	}

	for _, expected := range expectedTags {
		if !tagMap[expected] {
			t.Errorf("Expected tag %q not found", expected)
		}
	}
}

// Test content XML with draw:frame elements
const testContentXML = `<?xml version="1.0" encoding="UTF-8"?>
<office:document-content xmlns:draw="urn:oasis:names:tc:opendocument:xmlns:drawing:1.0"
    xmlns:xlink="http://www.w3.org/1999/xlink"
    xmlns:svg="urn:oasis:names:tc:opendocument:xmlns:svg-compatible:1.0"
    xmlns:office="urn:oasis:names:tc:opendocument:xmlns:office:1.0">
    <office:body>
        <draw:frame draw:name="image1" draw:style-name="fr1">
            <draw:image xlink:href="Pictures/img1.png" />
            <svg:title>{img1}</svg:title>
        </draw:frame>
        <draw:frame draw:name="photo1" draw:style-name="fr2">
            <draw:image xlink:href="Pictures/photo.jpg" />
        </draw:frame>
    </office:body>
</office:document-content>`

const testManifestXML = `<?xml version="1.0" encoding="UTF-8"?>
<manifest:manifest xmlns:manifest="urn:oasis:names:tc:opendocument:xmlns:manifest:1.0">
    <manifest:file-entry manifest:full-path="/" manifest:media-type="application/vnd.oasis.opendocument.text" />
    <manifest:file-entry manifest:full-path="content.xml" manifest:media-type="text/xml" />
</manifest:manifest>`

// Helper function to create a minimal valid ODT file for testing
func createMinimalODT(path string) error {
	return createODTWithContent(path, `<?xml version="1.0" encoding="UTF-8"?>
<office:document-content xmlns:office="urn:oasis:names:tc:opendocument:xmlns:office:1.0">
    <office:body></office:body>
</office:document-content>`)
}

// Helper function to create ODT with custom content.xml
func createODTWithContent(path string, content string) error {
	// Create files map
	files := make(map[string][]byte)
	files["mimetype"] = []byte("application/vnd.oasis.opendocument.text")
	files["content.xml"] = []byte(content)
	files["META-INF/manifest.xml"] = []byte(testManifestXML)

	// Create ZIP file directly
	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)

	for name, data := range files {
		fw, err := writer.Create(name)
		if err != nil {
			return err
		}
		if _, err := fw.Write(data); err != nil {
			return err
		}
	}

	if err := writer.Close(); err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(path, buf.Bytes(), 0644)
}

func TestODTDocument_ReplaceImageByTag(t *testing.T) {
	tmpDir := t.TempDir()
	testODT := filepath.Join(tmpDir, "test.odt")

	// Create ODT with test content
	if err := createODTWithContent(testODT, testContentXML); err != nil {
		t.Fatalf("Failed to create test ODT: %v", err)
	}

	doc, err := NewODTDocument(testODT)
	if err != nil {
		t.Fatalf("NewODTDocument() error = %v", err)
	}

	// Replace image
	newImage := []byte("new image data")
	err = doc.ReplaceImageByTag("image1", "Pictures/newimg.png", newImage)
	if err != nil {
		t.Errorf("ReplaceImageByTag() error = %v", err)
	}

	// Verify content.xml was updated
	content, err := doc.getContentXML()
	if err != nil {
		t.Fatalf("getContentXML() error = %v", err)
	}

	if !strings.Contains(content, "Pictures/newimg.png") {
		t.Error("content.xml does not contain new image path")
	}

	// Verify image was added
	if data, ok := doc.files["Pictures/newimg.png"]; !ok {
		t.Error("New image not found in files")
	} else if !bytes.Equal(data, newImage) {
		t.Error("New image data mismatch")
	}

	// Test replacing non-existent tag
	err = doc.ReplaceImageByTag("nonexistent", "Pictures/fake.png", newImage)
	if err == nil {
		t.Error("ReplaceImageByTag() should fail with non-existent tag")
	}
}

func TestODTDocument_Save(t *testing.T) {
	tmpDir := t.TempDir()
	testODT := filepath.Join(tmpDir, "test.odt")
	outputODT := filepath.Join(tmpDir, "output.odt")

	// Create minimal ODT
	if err := createMinimalODT(testODT); err != nil {
		t.Fatalf("Failed to create test ODT: %v", err)
	}

	doc, err := NewODTDocument(testODT)
	if err != nil {
		t.Fatalf("NewODTDocument() error = %v", err)
	}

	// Save to new location
	err = doc.Save(outputODT)
	if err != nil {
		t.Errorf("Save() error = %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outputODT); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}

	// Test saving with invalid path
	err = doc.Save("../../../etc/passwd")
	if err == nil {
		t.Error("Save() should fail with path traversal")
	}
}

func TestODTDocument_PreservesAllFiles(t *testing.T) {
	tmpDir := t.TempDir()
	testODT := filepath.Join(tmpDir, "test.odt")
	outputODT := filepath.Join(tmpDir, "output.odt")

	// Create ODT with test content
	if err := createODTWithContent(testODT, testContentXML); err != nil {
		t.Fatalf("Failed to create test ODT: %v", err)
	}

	// Open and modify document
	doc, err := NewODTDocument(testODT)
	if err != nil {
		t.Fatalf("NewODTDocument() error = %v", err)
	}

	// Count original files
	originalFileCount := len(doc.reader.File)

	// Replace an image (modifies content.xml and manifest.xml)
	newImage := []byte("new image data")
	err = doc.ReplaceImageByTag("image1", "Pictures/newimg.png", newImage)
	if err != nil {
		t.Fatalf("ReplaceImageByTag() error = %v", err)
	}

	// Save document
	err = doc.Save(outputODT)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Open saved document and verify all files are present
	savedDoc, err := NewODTDocument(outputODT)
	if err != nil {
		t.Fatalf("Failed to open saved ODT: %v", err)
	}

	// Count files in saved document (should have 1 more: the new image)
	savedFileCount := len(savedDoc.reader.File)
	expectedCount := originalFileCount + 1 // original files + new image

	if savedFileCount != expectedCount {
		t.Errorf("File count mismatch: got %d files, expected %d", savedFileCount, expectedCount)
	}

	// Verify all original files are still present
	originalFiles := make(map[string]bool)
	for _, f := range doc.reader.File {
		originalFiles[f.Name] = true
	}

	for _, f := range savedDoc.reader.File {
		if f.Name != "Pictures/newimg.png" {
			if !originalFiles[f.Name] {
				t.Errorf("Saved document has unexpected file: %s", f.Name)
			}
		}
	}

	// Verify the new image exists
	foundNewImage := false
	for _, f := range savedDoc.reader.File {
		if f.Name == "Pictures/newimg.png" {
			foundNewImage = true
			break
		}
	}
	if !foundNewImage {
		t.Error("New image not found in saved document")
	}

	// Verify mimetype is preserved
	foundMimetype := false
	for _, f := range savedDoc.reader.File {
		if f.Name == "mimetype" {
			foundMimetype = true
			break
		}
	}
	if !foundMimetype {
		t.Error("mimetype file not preserved in saved document")
	}
}

func BenchmarkNewODTDocument(b *testing.B) {
	tmpDir := b.TempDir()
	testODT := filepath.Join(tmpDir, "bench.odt")

	if err := createMinimalODT(testODT); err != nil {
		b.Fatalf("Failed to create test ODT: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewODTDocument(testODT)
		if err != nil {
			b.Fatalf("NewODTDocument() error = %v", err)
		}
	}
}

func BenchmarkReplaceImageByTag(b *testing.B) {
	tmpDir := b.TempDir()
	testODT := filepath.Join(tmpDir, "bench.odt")

	if err := createODTWithContent(testODT, testContentXML); err != nil {
		b.Fatalf("Failed to create test ODT: %v", err)
	}

	doc, err := NewODTDocument(testODT)
	if err != nil {
		b.Fatalf("NewODTDocument() error = %v", err)
	}

	imageData := []byte("benchmark image data")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset document for each iteration
		if err := createODTWithContent(testODT, testContentXML); err != nil {
			b.Fatal(err)
		}
		doc, _ = NewODTDocument(testODT)

		err := doc.ReplaceImageByTag("image1", "Pictures/bench.png", imageData)
		if err != nil {
			b.Fatalf("ReplaceImageByTag() error = %v", err)
		}
	}
}
