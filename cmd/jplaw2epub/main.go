package main

import (
	"flag"
	"fmt"
	"os"

	"go.ngs.io/jplaw2epub"
)

func main() {
	os.Exit(run())
}

func run() int {
	destPath, sourcePath, err := parseFlags()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 1
	}

	xmlFile, openErr := os.Open(sourcePath)
	if openErr != nil {
		fmt.Printf("Error opening source file: %v\n", openErr)
		return 1
	}
	defer xmlFile.Close()

	book, createErr := jplaw2epub.CreateEPUBFromXMLFile(xmlFile)
	if createErr != nil {
		fmt.Printf("Error creating EPUB file: %v\n", createErr)
		return 1
	}

	if writeErr := jplaw2epub.WriteEPUB(book, destPath); writeErr != nil {
		fmt.Printf("Error writing EPUB file: %v\n", writeErr)
		return 1
	}

	fmt.Printf("Successfully created EPUB: %s\n", destPath)
	return 0
}

func parseFlags() (destPath, sourcePath string, err error) {
	destPathFlag := flag.String("d", "", "Destination file path")
	flag.Parse()

	if *destPathFlag == "" {
		return "", "", fmt.Errorf("destination file path is required")
	}

	if len(flag.Args()) < 1 {
		return "", "", fmt.Errorf("source file path is required as the first argument")
	}

	return *destPathFlag, flag.Arg(0), nil
}
