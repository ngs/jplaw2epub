// Package jplaw2epub converts Japanese Standard Law XML Schema (法令標準XMLスキーマ) into EPUB files.
//
// This package provides functionality to convert Japanese legal documents in XML format
// into EPUB ebooks with proper formatting, Ruby annotations support, and Japanese-specific
// list styling.
//
// Basic usage:
//
//	book, err := jplaw2epub.CreateEPUBFromXMLPath("law.xml")
//	if err != nil {
//		return err
//	}
//	return jplaw2epub.WriteEPUB(book, "output.epub")
package jplaw2epub

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-shiori/go-epub"
	lawapi "go.ngs.io/jplaw-api-v2"
	"go.ngs.io/jplaw-xml"
)

// EPUBOptions contains options for EPUB creation
type EPUBOptions struct {
	// APIClient is the jplaw API client for downloading images
	APIClient *lawapi.Client
	// RevisionID is the revision ID for fetching attachments
	RevisionID string
	// MaxImageHeight is the maximum height for images (e.g., "300px", "80vh", "50%")
	MaxImageHeight string
}

// CreateEPUBFromXMLFile creates an EPUB file from a jplaw XML file reader.
//
// This function reads XML data from the provided io.Reader, parses it as Japanese
// law data according to the Japanese Standard Law XML Schema, and creates an EPUB
// book with proper formatting and metadata.
//
// Example:
//
//	xmlFile, err := os.Open("law.xml")
//	if err != nil {
//		return err
//	}
//	defer xmlFile.Close()
//
//	book, err := jplaw2epub.CreateEPUBFromXMLFile(xmlFile)
//	if err != nil {
//		return err
//	}
func CreateEPUBFromXMLFile(xmlFile io.Reader) (*epub.Epub, error) {
	return CreateEPUBFromXMLFileWithOptions(xmlFile, nil)
}

// CreateEPUBFromXMLFileWithOptions creates an EPUB file with image support
func CreateEPUBFromXMLFileWithOptions(xmlFile io.Reader, opts *EPUBOptions) (*epub.Epub, error) {
	// Load and parse XML data
	data, err := loadXMLDataFromReader(xmlFile)
	if err != nil {
		return nil, fmt.Errorf("loading XML data: %w", err)
	}

	// Create EPUB
	book, err := createEPUBFromData(data)
	if err != nil {
		return nil, fmt.Errorf("creating EPUB: %w", err)
	}

	// Process chapters and content
	if err := processChaptersWithOptions(book, data, opts); err != nil {
		return nil, fmt.Errorf("processing chapters: %w", err)
	}

	return book, nil
}

// CreateEPUBFromXMLPath creates an EPUB file from a jplaw XML file path.
//
// This is a convenience function that opens the file at the given path and
// calls CreateEPUBFromXMLFile to process it. Images will be automatically
// downloaded and embedded if the filename contains a valid revision ID.
//
// Example:
//
//	book, err := jplaw2epub.CreateEPUBFromXMLPath("law.xml")
//	if err != nil {
//		return err
//	}
func CreateEPUBFromXMLPath(xmlPath string) (*epub.Epub, error) {
	return CreateEPUBFromXMLPathWithOptions(xmlPath, nil)
}

// CreateEPUBFromXMLPathWithOptions creates an EPUB file with image support
func CreateEPUBFromXMLPathWithOptions(xmlPath string, opts *EPUBOptions) (*epub.Epub, error) {
	xmlFile, err := os.Open(xmlPath)
	if err != nil {
		return nil, fmt.Errorf("opening XML file: %w", err)
	}
	defer xmlFile.Close()

	return CreateEPUBFromXMLFileWithOptions(xmlFile, opts)
}

// WriteEPUB writes the EPUB book to the specified path.
//
// The function ensures the directory exists before writing and returns an error
// if the write operation fails.
//
// Example:
//
//	err := jplaw2epub.WriteEPUB(book, "output.epub")
//	if err != nil {
//		return err
//	}
func WriteEPUB(book *epub.Epub, destPath string) error {
	// Ensure directory exists
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	// Write EPUB
	if err := book.Write(destPath); err != nil {
		return fmt.Errorf("writing EPUB file: %w", err)
	}

	return nil
}

// loadXMLDataFromReader loads XML data from an io.Reader
func loadXMLDataFromReader(reader io.Reader) (*jplaw.Law, error) {
	byteValue, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("reading XML data: %w", err)
	}

	var data jplaw.Law
	if err := xml.Unmarshal(byteValue, &data); err != nil {
		return nil, fmt.Errorf("unmarshalling XML: %w", err)
	}

	return &data, nil
}

// createEPUBFromData creates and sets up EPUB from law data
func createEPUBFromData(data *jplaw.Law) (*epub.Epub, error) {
	if data.LawBody.LawTitle == nil {
		return nil, fmt.Errorf("law title is required")
	}
	book, err := epub.NewEpub(data.LawBody.LawTitle.Content)
	if err != nil {
		return nil, fmt.Errorf("creating epub: %w", err)
	}

	setupEPUBMetadata(book, data)

	// Add CSS styles for proper formatting
	if err := AddCSSToEPUB(book); err != nil {
		return nil, fmt.Errorf("adding CSS to EPUB: %w", err)
	}

	return book, nil
}

// createImageProcessor creates an image processor from options
func createImageProcessor(book *epub.Epub, opts *EPUBOptions) ImageProcessorInterface {
	if opts == nil || opts.APIClient == nil || opts.RevisionID == "" {
		return nil
	}

	imgProc := NewImageProcessor(opts.APIClient, opts.RevisionID, book)
	if opts.MaxImageHeight != "" {
		imgProc.SetMaxImageHeight(opts.MaxImageHeight)
	}
	return imgProc
}

// processChaptersWithOptions processes all chapters with image support
func processChaptersWithOptions(book *epub.Epub, data *jplaw.Law, opts *EPUBOptions) error {
	// Create image processor if API client is available
	imgProc := createImageProcessor(book, opts)

	// Add title page as the first page
	if err := addTitlePage(book, data); err != nil {
		return fmt.Errorf("adding title page: %w", err)
	}

	// Process main provision content
	if err := processMainProvision(book, &data.LawBody.MainProvision, imgProc); err != nil {
		return err
	}

	// Process AppdxNote (appendix notes)
	if len(data.LawBody.AppdxNote) > 0 {
		if err := processAppdxNotes(book, data.LawBody.AppdxNote, imgProc); err != nil {
			return fmt.Errorf("processing appendix notes: %w", err)
		}
	}

	// Process AppdxTable (appendix tables)
	if len(data.LawBody.AppdxTable) > 0 {
		if err := processAppdxTables(book, data.LawBody.AppdxTable, imgProc); err != nil {
			return fmt.Errorf("processing appendix tables: %w", err)
		}
	}

	// Process AppdxStyle (appendix styles with images)
	if len(data.LawBody.AppdxStyle) > 0 {
		if err := processAppdxStyles(book, data.LawBody.AppdxStyle, imgProc); err != nil {
			return err
		}
	}

	// Process AppdxFormat (appendix formats)
	if len(data.LawBody.AppdxFormat) > 0 {
		if err := processAppdxFormats(book, data.LawBody.AppdxFormat, imgProc); err != nil {
			return fmt.Errorf("processing appendix formats: %w", err)
		}
	}

	// Process SupplProvision (supplementary provisions)
	if len(data.LawBody.SupplProvision) > 0 {
		if err := processSupplProvisions(book, data.LawBody.SupplProvision, imgProc); err != nil {
			return fmt.Errorf("processing supplementary provisions: %w", err)
		}
	}

	return nil
}
