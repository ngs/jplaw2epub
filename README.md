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

## Web API Server

A web API server is available that provides REST and GraphQL endpoints for converting Japanese law documents to EPUB format.

### Installation

```sh
go install go.ngs.io/jplaw2epub/cmd/jplaw2epub-server@latest
```

Or build from source:

```sh
go build -o jplaw2epub-server ./cmd/jplaw2epub-server
```

### Running the Server

```sh
# Use automatic port selection
jplaw2epub-server

# Specify port via flag
jplaw2epub-server -port 8080

# Specify port via environment variable
PORT=8080 jplaw2epub-server
```

### API Endpoints

#### REST API

- **POST /convert** - Convert XML to EPUB
  ```sh
  curl -X POST -H "Content-Type: application/xml" \
    --data-binary @law.xml \
    http://localhost:8080/convert -o output.epub
  ```

- **GET /epubs/{law_id}** - Get EPUB by law ID (uses jplaw-api-v2)
  ```sh
  curl http://localhost:8080/epubs/325AC0000000131 -o radio_act.epub
  ```

- **GET /health** - Health check endpoint

#### GraphQL API

The server includes a GraphQL API powered by [gqlgen](https://gqlgen.com/) for querying Japanese law data.

- **POST/GET /graphql** - GraphQL endpoint
- **GET /graphiql** - Interactive GraphQL playground

##### GraphQL Schema

The GraphQL implementation is located in `/graphql/` directory with:
- `schema.graphqls` - GraphQL schema definition
- `gqlgen.yml` - gqlgen configuration
- Generated resolvers and models

##### Example Queries

Search laws by category and type:
```graphql
query {
  laws(
    categoryCode: [CONSTITUTION, CRIMINAL]
    lawType: [ACT]
    limit: 5
  ) {
    totalCount
    laws {
      lawInfo {
        lawId
        lawNum
        lawType
        promulgationDate
      }
      revisionInfo {
        lawTitle
        lawTitleKana
      }
    }
  }
}
```

Get law revisions:
```graphql
query {
  revisions(lawId: "325AC0000000131") {
    lawInfo {
      lawNum
      promulgationDate
    }
    revisions {
      amendmentLawTitle
      amendmentEnforcementDate
      currentRevisionStatus
    }
  }
}
```

Keyword search:
```graphql
query {
  keyword(keyword: "無線", limit: 3) {
    totalCount
    items {
      lawInfo {
        lawId
      }
      revisionInfo {
        lawTitle
      }
      sentences {
        text
        position
      }
    }
  }
}
```

### GraphQL Development

To regenerate GraphQL code after schema changes:

```sh
cd graphql
gqlgen generate
```

The GraphQL implementation uses:
- Schema-first approach with type-safe code generation
- Direct binding to jplaw-api-v2 types
- Enum support for all law categories and types
- Custom converters for enum value mapping

### Docker Deployment

See [cmd/jplaw2epub-server/README.md](cmd/jplaw2epub-server/README.md) for Docker and Google Cloud Run deployment instructions.

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