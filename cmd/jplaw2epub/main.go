package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	lawapi "go.ngs.io/jplaw-api-v2"
	"go.ngs.io/jplaw2epub"
)

func main() {
	os.Exit(run())
}

func run() int {
	opts, err := parseFlags()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 1
	}

	xmlFile, openErr := os.Open(opts.sourcePath)
	if openErr != nil {
		fmt.Printf("Error opening source file: %v\n", openErr)
		return 1
	}
	defer xmlFile.Close()

	// Create EPUB options
	epubOpts := createEPUBOptions(opts)

	book, createErr := jplaw2epub.CreateEPUBFromXMLFileWithOptions(xmlFile, epubOpts)
	if createErr != nil {
		fmt.Printf("Error creating EPUB file: %v\n", createErr)
		return 1
	}

	if writeErr := jplaw2epub.WriteEPUB(book, opts.destPath); writeErr != nil {
		fmt.Printf("Error writing EPUB file: %v\n", writeErr)
		return 1
	}

	fmt.Printf("Successfully created EPUB: %s\n", opts.destPath)
	return 0
}

type options struct {
	destPath       string
	sourcePath     string
	downloadImages bool
	maxImageHeight string
}

func parseFlags() (*options, error) {
	destPathFlag := flag.String("d", "", "Destination file path")
	downloadImagesFlag := flag.Bool("no-images", false, "Skip downloading and embedding images")
	maxImageHeightFlag := flag.String("max-image-height", "80vh", "Maximum image height (e.g., '300px', '80vh', '50%')")
	// For backward compatibility, also accept the old -images flag
	oldImagesFlag := flag.Bool("images", false, "Download and embed images (deprecated, images are embedded by default)")
	flag.Parse()

	if *destPathFlag == "" {
		return nil, fmt.Errorf("destination file path is required")
	}

	if len(flag.Args()) < 1 {
		return nil, fmt.Errorf("source file path is required as the first argument")
	}

	// Default to downloading images unless explicitly disabled
	downloadImages := !*downloadImagesFlag || *oldImagesFlag
	
	opts := &options{
		destPath:       *destPathFlag,
		sourcePath:     flag.Arg(0),
		downloadImages: downloadImages,
		maxImageHeight: *maxImageHeightFlag,
	}

	return opts, nil
}

func extractRevisionIDFromPath(path string) string {
	base := filepath.Base(path)
	base = strings.TrimSuffix(base, filepath.Ext(base))
	// The entire base name (without extension) is the revision ID for the API
	// Format: lawID_date_revisionID.xml -> use the whole base as revision ID
	if strings.Count(base, "_") >= 2 {
		return base
	}
	return ""
}

func createEPUBOptions(opts *options) *jplaw2epub.EPUBOptions {
	if !opts.downloadImages {
		return nil
	}

	// Extract revision ID from source path
	revisionID := extractRevisionIDFromPath(opts.sourcePath)
	if revisionID == "" {
		fmt.Println("Warning: Could not extract revision ID from filename, images will not be downloaded")
		return nil
	}

	// Create API client
	client := lawapi.NewClient()

	return &jplaw2epub.EPUBOptions{
		APIClient:      client,
		RevisionID:     revisionID,
		MaxImageHeight: opts.maxImageHeight,
	}
}
