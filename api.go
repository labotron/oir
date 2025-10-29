package odtimagereplacer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// ImageSource represents an image from either URL or base64
type ImageSource struct {
	URL    string `json:"url"`
	Base64 string `json:"base64"`
}

// TemplateSource represents the ODT template source
type TemplateSource struct {
	URL    string `json:"url"`
	Base64 string `json:"base64"`
}

// ReplaceRequest represents the JSON request structure for replacing images
type ReplaceRequest struct {
	Template TemplateSource         `json:"template"`
	Data     map[string]ImageSource `json:"data"`
}

// ReplaceResponse represents the JSON response structure
type ReplaceResponse struct {
	Success      bool     `json:"success"`
	Message      string   `json:"message,omitempty"`
	OutputBase64 string   `json:"output_base64,omitempty"`
	ReplacedTags []string `json:"replaced_tags,omitempty"`
	Error        string   `json:"error,omitempty"`
}

// HTTPClient interface for testing
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// DefaultHTTPClient is the default HTTP client with timeout
var DefaultHTTPClient HTTPClient = &http.Client{
	Timeout: 30 * time.Second,
}

// fetchImageFromURL downloads an image from a URL
func fetchImageFromURL(url string, client HTTPClient) ([]byte, error) {
	if url == "" || url == "null" {
		return nil, fmt.Errorf("invalid URL")
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	// Limit response size to prevent memory exhaustion
	limitReader := io.LimitReader(resp.Body, MaxIndividualFileSize+1)
	data, err := io.ReadAll(limitReader)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if len(data) > MaxIndividualFileSize {
		return nil, fmt.Errorf("%w: image from URL exceeds limit", ErrFileTooLarge)
	}

	return data, nil
}

// decodeBase64Image decodes a base64-encoded image
func decodeBase64Image(b64 string) ([]byte, error) {
	if b64 == "" || b64 == "null" {
		return nil, fmt.Errorf("invalid base64 data")
	}

	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("decode base64: %w", err)
	}

	if len(data) > MaxIndividualFileSize {
		return nil, fmt.Errorf("%w: decoded image exceeds limit", ErrFileTooLarge)
	}

	return data, nil
}

// getImageData retrieves image data from either URL or base64
func getImageData(source ImageSource, client HTTPClient) ([]byte, error) {
	// Try URL first if provided
	if source.URL != "" && source.URL != "null" {
		return fetchImageFromURL(source.URL, client)
	}

	// Try base64 if URL not provided
	if source.Base64 != "" && source.Base64 != "null" {
		return decodeBase64Image(source.Base64)
	}

	return nil, fmt.Errorf("no valid image source provided (URL or base64)")
}

// getTemplateData retrieves ODT template data from either URL or base64
func getTemplateData(source TemplateSource, client HTTPClient) ([]byte, error) {
	// Try URL first if provided
	if source.URL != "" && source.URL != "null" {
		return fetchImageFromURL(source.URL, client)
	}

	// Try base64 if URL not provided
	if source.Base64 != "" && source.Base64 != "null" {
		return decodeBase64Image(source.Base64)
	}

	return nil, fmt.Errorf("no valid template source provided (URL or base64)")
}

// ProcessReplaceRequest processes a replace request and returns the modified ODT
func ProcessReplaceRequest(req ReplaceRequest) (*ReplaceResponse, []byte, error) {
	return ProcessReplaceRequestWithClient(req, DefaultHTTPClient)
}

// ProcessReplaceRequestWithClient processes a replace request with a custom HTTP client
func ProcessReplaceRequestWithClient(req ReplaceRequest, client HTTPClient) (*ReplaceResponse, []byte, error) {
	// Validate request
	if len(req.Data) == 0 {
		return &ReplaceResponse{
			Success: false,
			Error:   "no images to replace",
		}, nil, fmt.Errorf("no images to replace")
	}

	// Get template data
	templateData, err := getTemplateData(req.Template, client)
	if err != nil {
		return &ReplaceResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to get template: %v", err),
		}, nil, fmt.Errorf("get template: %w", err)
	}

	// Create temporary ODT document from template data
	doc, err := NewODTDocumentFromBytes(templateData)
	if err != nil {
		return &ReplaceResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to parse template: %v", err),
		}, nil, fmt.Errorf("parse template: %w", err)
	}

	// Track successfully replaced tags
	replacedTags := make([]string, 0, len(req.Data))
	var lastErr error

	// Process each image replacement
	for tag, imageSource := range req.Data {
		// Get image data
		imageData, err := getImageData(imageSource, client)
		if err != nil {
			lastErr = fmt.Errorf("get image for tag '%s': %w", tag, err)
			continue
		}

		// Determine image path in ODT
		imagePath := fmt.Sprintf("Pictures/%s.png", tag)

		// Try to detect extension from data
		if len(imageData) > 0 {
			ext := detectImageExtension(imageData)
			if ext != "" {
				imagePath = fmt.Sprintf("Pictures/%s%s", tag, ext)
			}
		}

		// Replace the image
		err = doc.ReplaceImageByTag(tag, imagePath, imageData)
		if err != nil {
			lastErr = fmt.Errorf("replace image for tag '%s': %w", tag, err)
			continue
		}

		replacedTags = append(replacedTags, tag)
	}

	// Check if any images were replaced
	if len(replacedTags) == 0 {
		return &ReplaceResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to replace any images: %v", lastErr),
		}, nil, fmt.Errorf("no images replaced: %w", lastErr)
	}

	// Save to bytes
	outputData, err := doc.SaveToBytes()
	if err != nil {
		return &ReplaceResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to save output: %v", err),
		}, nil, fmt.Errorf("save output: %w", err)
	}

	// Create response
	response := &ReplaceResponse{
		Success:      true,
		Message:      fmt.Sprintf("Successfully replaced %d image(s)", len(replacedTags)),
		ReplacedTags: replacedTags,
	}

	return response, outputData, nil
}

// detectImageExtension detects image file extension from magic bytes
func detectImageExtension(data []byte) string {
	if len(data) < 12 {
		return ""
	}

	// Check PNG signature
	if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return ".png"
	}

	// Check JPEG signature
	if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return ".jpg"
	}

	// Check GIF signature
	if data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 {
		return ".gif"
	}

	// Check WebP signature
	if data[0] == 0x52 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x46 &&
		data[8] == 0x57 && data[9] == 0x45 && data[10] == 0x42 && data[11] == 0x50 {
		return ".webp"
	}

	return ""
}

// NewODTDocumentFromBytes creates an ODT document from byte data
func NewODTDocumentFromBytes(data []byte) (*ODTDocument, error) {
	// Validate size
	if len(data) > MaxFileSize {
		return nil, fmt.Errorf("%w: %d bytes (max: %d)", ErrFileTooLarge, len(data), MaxFileSize)
	}

	// Create temporary file
	tmpFile, err := createTempODTFile(data)
	if err != nil {
		return nil, fmt.Errorf("create temp file: %w", err)
	}

	// Open as ODT document
	doc, err := NewODTDocument(tmpFile)
	if err != nil {
		return nil, fmt.Errorf("parse ODT: %w", err)
	}

	return doc, nil
}

// createTempODTFile creates a temporary ODT file from byte data
func createTempODTFile(data []byte) (string, error) {
	tmpFile, err := os.CreateTemp("", "odt-*.odt")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(data); err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}

// ParseReplaceRequest parses a JSON request body
func ParseReplaceRequest(body []byte) (*ReplaceRequest, error) {
	var req ReplaceRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}
	return &req, nil
}
