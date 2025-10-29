package odtimagereplacer

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

const (
	// MaxFileSize is the maximum allowed ODT file size (100MB)
	MaxFileSize = 100 * 1024 * 1024

	// MaxFilesInArchive is the maximum number of files allowed in ZIP
	MaxFilesInArchive = 10000

	// MaxIndividualFileSize is the maximum size for individual files in archive (50MB)
	MaxIndividualFileSize = 50 * 1024 * 1024
)

var (
	// Pre-compiled regex patterns for performance
	drawFrameRegex     *regexp.Regexp
	manifestEndRegex   *regexp.Regexp
	regexOnce          sync.Once
	regexCompileErrors []error
)

// initRegex initializes all regex patterns once
func initRegex() {
	regexOnce.Do(func() {
		var err error

		drawFrameRegex, err = regexp.Compile(`<draw:frame[^>]*>[\s\S]*?</draw:frame>`)
		if err != nil {
			regexCompileErrors = append(regexCompileErrors, fmt.Errorf("draw frame regex: %w", err))
		}

		manifestEndRegex, err = regexp.Compile(`(?i)</manifest:manifest>`)
		if err != nil {
			regexCompileErrors = append(regexCompileErrors, fmt.Errorf("manifest end regex: %w", err))
		}
	})
}

// ODTDocument represents an ODT document with methods for manipulation
type ODTDocument struct {
	path   string
	reader *zip.Reader
	files  map[string][]byte
}

// NewODTDocument creates a new ODT document from a file path
func NewODTDocument(path string) (*ODTDocument, error) {
	// Initialize regex patterns
	initRegex()
	if len(regexCompileErrors) > 0 {
		return nil, fmt.Errorf("regex compilation failed: %v", regexCompileErrors)
	}

	// Validate file path
	if err := validatePath(path); err != nil {
		return nil, err
	}

	// Check file size before reading
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat file: %w", err)
	}

	if fileInfo.Size() > MaxFileSize {
		return nil, fmt.Errorf("%w: %d bytes (max: %d)", ErrFileTooLarge, fileInfo.Size(), MaxFileSize)
	}

	// Read file into memory
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	// Create ZIP reader
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidODT, err)
	}

	// Validate ZIP structure
	if len(reader.File) > MaxFilesInArchive {
		return nil, fmt.Errorf("%w: %d files (max: %d)", ErrTooManyFiles, len(reader.File), MaxFilesInArchive)
	}

	doc := &ODTDocument{
		path:   path,
		reader: reader,
		files:  make(map[string][]byte, len(reader.File)),
	}

	return doc, nil
}

// validatePath checks for path traversal attempts and validates the path
func validatePath(path string) error {
	if path == "" {
		return fmt.Errorf("%w: empty path", ErrInvalidPath)
	}

	// Check for path traversal attempts BEFORE cleaning
	if strings.Contains(path, "..") {
		return fmt.Errorf("%w: path traversal detected", ErrInvalidPath)
	}

	return nil
}

// validateImageName validates image filename for security
func validateImageName(name string) error {
	if name == "" {
		return fmt.Errorf("%w: empty name", ErrInvalidImageName)
	}

	// Check for path traversal
	if strings.Contains(name, "..") || strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return fmt.Errorf("%w: invalid characters in name", ErrInvalidImageName)
	}

	// Check filename length
	if len(name) > 255 {
		return fmt.Errorf("%w: name too long", ErrInvalidImageName)
	}

	return nil
}

// loadFile safely loads a single file from the ZIP archive
func (doc *ODTDocument) loadFile(f *zip.File) ([]byte, error) {
	// Validate file path
	if err := validatePath(f.Name); err != nil {
		return nil, err
	}

	// Check uncompressed size
	if f.UncompressedSize64 > MaxIndividualFileSize {
		return nil, fmt.Errorf("%w: file %s is %d bytes (max: %d)",
			ErrFileTooLarge, f.Name, f.UncompressedSize64, MaxIndividualFileSize)
	}

	rc, err := f.Open()
	if err != nil {
		return nil, fmt.Errorf("open file %s: %w", f.Name, err)
	}
	defer rc.Close()

	// Use LimitReader to prevent zip bomb attacks
	limitReader := io.LimitReader(rc, MaxIndividualFileSize+1)
	data, err := io.ReadAll(limitReader)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", f.Name, err)
	}

	// Double-check size after reading
	if len(data) > MaxIndividualFileSize {
		return nil, fmt.Errorf("%w: file %s exceeds limit after decompression", ErrFileTooLarge, f.Name)
	}

	return data, nil
}

// loadAllFiles loads all files from the ZIP archive into memory
func (doc *ODTDocument) loadAllFiles() error {
	for _, f := range doc.reader.File {
		data, err := doc.loadFile(f)
		if err != nil {
			return err
		}
		doc.files[f.Name] = data
	}
	return nil
}

// getFile retrieves a file from the archive, loading it if necessary
func (doc *ODTDocument) getFile(name string) ([]byte, error) {
	// Check if already loaded
	if data, ok := doc.files[name]; ok {
		return data, nil
	}

	// Load the specific file
	for _, f := range doc.reader.File {
		if f.Name == name {
			data, err := doc.loadFile(f)
			if err != nil {
				return nil, err
			}
			doc.files[name] = data
			return data, nil
		}
	}

	return nil, fmt.Errorf("file %s not found in archive", name)
}

// getContentXML retrieves and caches content.xml
func (doc *ODTDocument) getContentXML() (string, error) {
	data, err := doc.getFile("content.xml")
	if err != nil {
		return "", ErrContentNotFound
	}
	return string(data), nil
}

// getManifestXML retrieves and caches manifest.xml
func (doc *ODTDocument) getManifestXML() (string, error) {
	data, err := doc.getFile("META-INF/manifest.xml")
	if err != nil {
		return "", ErrManifestNotFound
	}
	return string(data), nil
}

// ReplaceImageByTag replaces an image in the ODT by its draw:name tag
func (doc *ODTDocument) ReplaceImageByTag(tag, newImagePath string, newImageData []byte) error {
	// Validate inputs
	if tag == "" {
		return fmt.Errorf("tag cannot be empty")
	}
	if len(newImageData) == 0 {
		return fmt.Errorf("image data cannot be empty")
	}
	if len(newImageData) > MaxIndividualFileSize {
		return fmt.Errorf("%w: image size %d exceeds limit", ErrFileTooLarge, len(newImageData))
	}

	// Validate image path
	imageName := filepath.Base(newImagePath)
	if err := validateImageName(imageName); err != nil {
		return err
	}

	// Get content.xml
	content, err := doc.getContentXML()
	if err != nil {
		return err
	}

	// Build regex pattern for this specific tag
	pattern := fmt.Sprintf(`(<draw:frame[^>]*draw:name="%s"[^>]*>[\s\S]*?xlink:href=")[^"]+("[\s\S]*?</draw:frame>)`, regexp.QuoteMeta(tag))
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("compile replacement regex: %w", err)
	}

	// Replace the image reference
	newContent := re.ReplaceAllString(content, "${1}"+newImagePath+"${2}")

	// Check if replacement occurred
	if newContent == content {
		return fmt.Errorf("%w: tag '%s'", ErrImageNotFound, tag)
	}

	// Update content.xml
	doc.files["content.xml"] = []byte(newContent)

	// Update manifest.xml
	if err := doc.addImageToManifest(newImagePath); err != nil {
		return err
	}

	// Add image data
	doc.files[newImagePath] = newImageData

	return nil
}

// addImageToManifest adds or updates an image entry in manifest.xml
func (doc *ODTDocument) addImageToManifest(imagePath string) error {
	manifest, err := doc.getManifestXML()
	if err != nil {
		return err
	}

	// Detect MIME type from extension
	mimeType := detectMIMEType(imagePath)

	// Check if entry already exists
	entryPattern := fmt.Sprintf(`<manifest:file-entry[^>]*manifest:full-path="%s"[^>]*/>`, regexp.QuoteMeta(imagePath))
	if matched, _ := regexp.MatchString(entryPattern, manifest); matched {
		// Entry already exists, no need to add
		return nil
	}

	// Create new manifest entry
	entry := fmt.Sprintf(`    <manifest:file-entry manifest:full-path="%s" manifest:media-type="%s" />`,
		imagePath, mimeType)

	// Insert before closing tag
	newManifest := manifestEndRegex.ReplaceAllString(manifest, entry+"\n</manifest:manifest>")

	doc.files["META-INF/manifest.xml"] = []byte(newManifest)
	return nil
}

// detectMIMEType returns MIME type based on file extension
func detectMIMEType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

// Save writes the modified ODT back to disk
func (doc *ODTDocument) Save(outputPath string) error {
	// Validate output path
	if err := validatePath(outputPath); err != nil {
		return err
	}

	// Create buffer for ZIP
	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)

	// Write all files from original ZIP, but use modified versions if they exist
	for _, f := range doc.reader.File {
		var data []byte
		var err error

		// Check if this file has been modified
		if modifiedData, exists := doc.files[f.Name]; exists {
			// Use the modified version
			data = modifiedData
		} else {
			// Use the original version
			data, err = doc.loadFile(f)
			if err != nil {
				return fmt.Errorf("load original file %s: %w", f.Name, err)
			}
		}

		// Write to new ZIP
		fw, err := writer.Create(f.Name)
		if err != nil {
			return fmt.Errorf("create zip entry %s: %w", f.Name, err)
		}
		if _, err := fw.Write(data); err != nil {
			return fmt.Errorf("write zip entry %s: %w", f.Name, err)
		}
	}

	// Write any new files that weren't in the original (e.g., new images)
	for name, data := range doc.files {
		// Check if this file already exists in original ZIP
		exists := false
		for _, f := range doc.reader.File {
			if f.Name == name {
				exists = true
				break
			}
		}

		// Only write if it's a new file
		if !exists {
			fw, err := writer.Create(name)
			if err != nil {
				return fmt.Errorf("create zip entry %s: %w", name, err)
			}
			if _, err := fw.Write(data); err != nil {
				return fmt.Errorf("write zip entry %s: %w", name, err)
			}
		}
	}

	// Close ZIP writer
	if err := writer.Close(); err != nil {
		return fmt.Errorf("close zip writer: %w", err)
	}

	// Write to file with secure permissions
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}

// SaveToBytes saves the ODT document to a byte slice
func (doc *ODTDocument) SaveToBytes() ([]byte, error) {
	// Create buffer for ZIP
	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)

	// Write all files from original ZIP, but use modified versions if they exist
	for _, f := range doc.reader.File {
		var data []byte
		var err error

		// Check if this file has been modified
		if modifiedData, exists := doc.files[f.Name]; exists {
			// Use the modified version
			data = modifiedData
		} else {
			// Use the original version
			data, err = doc.loadFile(f)
			if err != nil {
				return nil, fmt.Errorf("load original file %s: %w", f.Name, err)
			}
		}

		// Write to new ZIP
		fw, err := writer.Create(f.Name)
		if err != nil {
			return nil, fmt.Errorf("create zip entry %s: %w", f.Name, err)
		}
		if _, err := fw.Write(data); err != nil {
			return nil, fmt.Errorf("write zip entry %s: %w", f.Name, err)
		}
	}

	// Write any new files that weren't in the original (e.g., new images)
	for name, data := range doc.files {
		// Check if this file already exists in original ZIP
		exists := false
		for _, f := range doc.reader.File {
			if f.Name == name {
				exists = true
				break
			}
		}

		// Only write if it's a new file
		if !exists {
			fw, err := writer.Create(name)
			if err != nil {
				return nil, fmt.Errorf("create zip entry %s: %w", name, err)
			}
			if _, err := fw.Write(data); err != nil {
				return nil, fmt.Errorf("write zip entry %s: %w", name, err)
			}
		}
	}

	// Close ZIP writer
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close zip writer: %w", err)
	}

	return buf.Bytes(), nil
}

// AddImage adds an image to the ODT at the specified path
func (doc *ODTDocument) AddImage(imagePath string, imageData []byte) error {
	// Validate inputs
	if len(imageData) == 0 {
		return fmt.Errorf("image data cannot be empty")
	}
	if len(imageData) > MaxIndividualFileSize {
		return fmt.Errorf("%w: image size %d exceeds limit", ErrFileTooLarge, len(imageData))
	}

	// Validate image path
	if err := validatePath(imagePath); err != nil {
		return err
	}

	// Ensure all files are loaded
	if len(doc.files) == 0 {
		if err := doc.loadAllFiles(); err != nil {
			return err
		}
	}

	// Add or replace the image
	doc.files[imagePath] = imageData

	// Update manifest if needed
	if err := doc.addImageToManifest(imagePath); err != nil {
		return err
	}

	return nil
}

// FindImageTags finds all image tags (draw:name attributes) in the document
func (doc *ODTDocument) FindImageTags() ([]string, error) {
	content, err := doc.getContentXML()
	if err != nil {
		return nil, err
	}

	// Find all draw:frame elements
	frames := drawFrameRegex.FindAllString(content, -1)

	// Extract draw:name attributes
	namePattern := regexp.MustCompile(`draw:name="([^"]+)"`)
	tags := make([]string, 0, len(frames))
	seen := make(map[string]bool)

	for _, frame := range frames {
		matches := namePattern.FindStringSubmatch(frame)
		if len(matches) > 1 {
			tag := matches[1]
			if !seen[tag] {
				tags = append(tags, tag)
				seen[tag] = true
			}
		}
	}

	return tags, nil
}
