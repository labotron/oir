package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/suttapak/odtimagereplacer"
)

func main() {
	// Define command-line flags
	odtPath := flag.String("odt", "", "Path to ODT file")
	imageTag := flag.String("tag", "", "Image tag (draw:name) to replace")
	imagePath := flag.String("image", "", "Path to new image file")
	newImageName := flag.String("name", "", "New image name in ODT (e.g., Pictures/image1.png)")
	output := flag.String("output", "", "Output ODT file path (defaults to overwriting input)")
	listTags := flag.Bool("list", false, "List all image tags in the ODT")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  # List all image tags in an ODT:\n")
		fmt.Fprintf(os.Stderr, "  %s -odt=report.odt -list\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  # Replace an image by tag:\n")
		fmt.Fprintf(os.Stderr, "  %s -odt=report.odt -tag=image1 -image=photo.png -name=Pictures/photo.png -output=result.odt\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  # Run legacy test (for backward compatibility):\n")
		fmt.Fprintf(os.Stderr, "  %s (no flags - runs legacy Test function)\n\n", os.Args[0])
	}

	flag.Parse()

	// If no flags provided, run legacy test for backward compatibility
	if flag.NFlag() == 0 {
		fmt.Println("Running legacy test mode...")
		odtimagereplacer.Test("./report.odt")
		return
	}

	// Validate required flags
	if *odtPath == "" {
		fmt.Fprintf(os.Stderr, "Error: -odt flag is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Open ODT document
	doc, err := odtimagereplacer.NewODTDocument(*odtPath)
	if err != nil {
		log.Fatalf("Error opening ODT: %v", err)
	}

	// List tags mode
	if *listTags {
		tags, err := doc.FindImageTags()
		if err != nil {
			log.Fatalf("Error finding image tags: %v", err)
		}

		fmt.Printf("Found %d image tag(s) in %s:\n", len(tags), *odtPath)
		for i, tag := range tags {
			fmt.Printf("  %d. %s\n", i+1, tag)
		}
		return
	}

	// Replace image mode
	if *imageTag == "" || *imagePath == "" || *newImageName == "" {
		fmt.Fprintf(os.Stderr, "Error: -tag, -image, and -name flags are required for image replacement\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Read new image file
	imageData, err := os.ReadFile(*imagePath)
	if err != nil {
		log.Fatalf("Error reading image file: %v", err)
	}

	// Replace image
	if err := doc.ReplaceImageByTag(*imageTag, *newImageName, imageData); err != nil {
		log.Fatalf("Error replacing image: %v", err)
	}

	// Determine output path
	outputPath := *output
	if outputPath == "" {
		outputPath = *odtPath
	}

	// Save document
	if err := doc.Save(outputPath); err != nil {
		log.Fatalf("Error saving ODT: %v", err)
	}

	fmt.Printf("Successfully replaced image '%s' in %s\n", *imageTag, outputPath)
}
