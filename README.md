# jplaw2epub

Go library and command line tool to convert [Japanese Standard Law XML Schema (法令標準XMLスキーマ)][xmldoc] into EPUB files.

## Library Usage

```go
package main

import (
    "fmt"
    "os"
    
    "go.ngs.io/jplaw2epub"
)

func main() {
    // Method 1: Create EPUB from XML file path
    book, err := jplaw2epub.CreateEPUBFromXMLPath("law.xml")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    if err := jplaw2epub.WriteEPUB(book, "output.epub"); err != nil {
        fmt.Printf("Error writing EPUB: %v\n", err)
        return
    }

    // Method 2: Create EPUB from io.Reader
    xmlFile, err := os.Open("law.xml")
    if err != nil {
        fmt.Printf("Error opening file: %v\n", err)
        return
    }
    defer xmlFile.Close()

    book2, err := jplaw2epub.CreateEPUBFromXMLFile(xmlFile)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    if err := jplaw2epub.WriteEPUB(book2, "output2.epub"); err != nil {
        fmt.Printf("Error writing EPUB: %v\n", err)
        return
    }

    fmt.Println("EPUB files created successfully!")
}
```

### API Functions

- `CreateEPUBFromXMLPath(xmlPath string) (*epub.Epub, error)` - Creates an EPUB from a file path
- `CreateEPUBFromXMLFile(xmlFile io.Reader) (*epub.Epub, error)` - Creates an EPUB from an io.Reader
- `WriteEPUB(book *epub.Epub, destPath string) error` - Writes an EPUB book to a file

## Command Line Usage

Install the CLI tool:

```sh
go install go.ngs.io/jplaw2epub/cmd/jplaw2epub@latest

jplaw2epub -d mylaw.epub path/to/law.xml
```

Or build from source:

```sh
git clone https://github.com/ngs/jplaw2epub.git
cd jplaw2epub
go build -o jplaw2epub ./cmd/jplaw2epub
./jplaw2epub -d mylaw.epub path/to/law.xml
```

## Installation as Go Library

Add to your Go project:

```sh
go get go.ngs.io/jplaw2epub
```

## Features

- Support for Ruby (ルビ) annotations for Japanese phonetic guides
- Dynamic list styling based on content (CJK ideographic, katakana-iroha, hiragana-iroha)
- Proper paragraph and item hierarchy handling
- Support for Subitem1 and Subitem2 structures
- Full EPUB metadata support with Japanese era information
- Clean Go library API for programmatic usage
- Cross-platform CLI tool with binary releases

## Testing

Run the test suite:

```sh
go test -v ./...
```

Generate coverage report:

```sh
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Author

[Atsushi Nagase]

## License

MIT License

Copyright &copy; 2025 [Atsushi Nagase]. See [LICENSE.md](LICENSE.md) for details.

[Atsushi Nagase]: https://ngs.io/
[xmldoc]: https://laws.e-gov.go.jp/docs/law-data-basic/419a603-xml-schema-for-japanese-law/