package main

import (
	"flag"
	"fmt"
	"os"

	"go.ngs.io/jplaw2epub"
)

func main() {
	destPath, sourcePath := parseFlags()

	xmlFile, err := os.Open(sourcePath)
	if err != nil {
		fmt.Printf("Error opening source file: %v\n", err)
		os.Exit(1)
	}
	defer xmlFile.Close()

	book, err := jplaw2epub.CreateEPUBFromXMLFile(xmlFile)
	if err != nil {
		fmt.Printf("Error creating EPUB file: %v\n", err)
		os.Exit(1)
	}

	if err := jplaw2epub.WriteEPUB(book, destPath); err != nil {
		fmt.Printf("Error writing EPUB file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully created EPUB: %s\n", destPath)
}

func parseFlags() (destPath, sourcePath string) {
	destPathFlag := flag.String("d", "", "Destination file path")
	flag.Parse()

	if *destPathFlag == "" {
		fmt.Println("Destination file path is required")
		os.Exit(1)
	}

	if len(flag.Args()) < 1 {
		fmt.Println("Source file path is required as the first argument")
		os.Exit(1)
	}

	return *destPathFlag, flag.Arg(0)
}
