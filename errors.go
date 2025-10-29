package odtimagereplacer

import "errors"

var (
	// ErrInvalidODT indicates the file is not a valid ODT file
	ErrInvalidODT = errors.New("invalid ODT file format")

	// ErrFileTooLarge indicates the file exceeds maximum allowed size
	ErrFileTooLarge = errors.New("file size exceeds maximum allowed limit")

	// ErrInvalidPath indicates a path traversal attempt or invalid path
	ErrInvalidPath = errors.New("invalid or unsafe file path")

	// ErrContentNotFound indicates content.xml was not found in ODT
	ErrContentNotFound = errors.New("content.xml not found in ODT file")

	// ErrManifestNotFound indicates manifest.xml was not found in ODT
	ErrManifestNotFound = errors.New("META-INF/manifest.xml not found in ODT file")

	// ErrImageNotFound indicates the specified image tag was not found
	ErrImageNotFound = errors.New("image with specified tag not found")

	// ErrInvalidImageName indicates an invalid image filename
	ErrInvalidImageName = errors.New("invalid image filename")

	// ErrTooManyFiles indicates too many files in ZIP archive
	ErrTooManyFiles = errors.New("too many files in archive")
)
