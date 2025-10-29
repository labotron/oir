package odtimagereplacer

import (
	"fmt"
)

// AddImage adds an image file into an existing ODT file.
// This is a legacy function maintained for backward compatibility.
// For new code, use ODTDocument.AddImage() instead.
//
// DEPRECATED: Use NewODTDocument() and ODTDocument.AddImage() for better
// error handling and security features.
func AddImage(odtPath string, imageName string, imageData []byte) error {
	doc, err := NewODTDocument(odtPath)
	if err != nil {
		return fmt.Errorf("open odt: %w", err)
	}

	// Add image to Pictures directory
	imagePath := "Pictures/" + imageName
	if err := doc.AddImage(imagePath, imageData); err != nil {
		return fmt.Errorf("add image: %w", err)
	}

	// Save back to original file
	if err := doc.Save(odtPath); err != nil {
		return fmt.Errorf("save odt: %w", err)
	}

	return nil
}

// UnzipMem lists all files in an ODT archive.
// This is a legacy function maintained for backward compatibility.
//
// DEPRECATED: Use NewODTDocument() and inspect the files directly for better control.
func UnzipMem(path string) error {
	doc, err := NewODTDocument(path)
	if err != nil {
		return err
	}

	// Load all files to list them
	if err := doc.loadAllFiles(); err != nil {
		return err
	}

	// Print file names
	for name := range doc.files {
		fmt.Println(name)
	}

	return nil
}
