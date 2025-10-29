package odtimagereplacer

import (
	"fmt"
	"log"
	"os"
)

// Test demonstrates the image replacement functionality.
// This is a legacy function maintained for backward compatibility.
//
// DEPRECATED: Use NewODTDocument() and ReplaceImageByTag() for better error handling.
func Test(path string) {
	// Open ODT document
	doc, err := NewODTDocument(path)
	if err != nil {
		log.Fatal(fmt.Errorf("open document: %w", err))
	}

	// Read the image file
	img, err := os.ReadFile("./exmaple/img1.png")
	if err != nil {
		log.Fatal(fmt.Errorf("read image: %w", err))
	}

	// Replace image by tag
	if err := doc.ReplaceImageByTag("image1", "Pictures/image1.png", img); err != nil {
		log.Fatal(fmt.Errorf("replace image: %w", err))
	}

	// Save modified document
	if err := doc.Save("test.odt"); err != nil {
		log.Fatal(fmt.Errorf("save document: %w", err))
	}

	log.Println("Successfully replaced image in test.odt")
}
