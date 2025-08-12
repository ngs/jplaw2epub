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
	"go.ngs.io/jplaw-xml"
)

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
	if err := processChapters(book, data); err != nil {
		return nil, fmt.Errorf("processing chapters: %w", err)
	}

	return book, nil
}

// CreateEPUBFromXMLPath creates an EPUB file from a jplaw XML file path.
//
// This is a convenience function that opens the file at the given path and
// calls CreateEPUBFromXMLFile to process it.
//
// Example:
//
//	book, err := jplaw2epub.CreateEPUBFromXMLPath("law.xml")
//	if err != nil {
//		return err
//	}
func CreateEPUBFromXMLPath(xmlPath string) (*epub.Epub, error) {
	xmlFile, err := os.Open(xmlPath)
	if err != nil {
		return nil, fmt.Errorf("opening XML file: %w", err)
	}
	defer xmlFile.Close()

	return CreateEPUBFromXMLFile(xmlFile)
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
	book, err := epub.NewEpub(data.LawBody.LawTitle.Content)
	if err != nil {
		return nil, fmt.Errorf("creating epub: %w", err)
	}

	setupEPUBMetadata(book, data)
	return book, nil
}

// setupEPUBMetadata sets up the basic EPUB metadata
func setupEPUBMetadata(book *epub.Epub, data *jplaw.Law) {
	book.SetAuthor(data.LawNum)
	book.SetLang(string(data.Lang))

	// Set description
	eraStr := getEraString(data.Era)
	description := fmt.Sprintf("公布日: %s %d年%d月%d日", eraStr, data.Year, data.PromulgateMonth, data.PromulgateDay)
	description += fmt.Sprintf("\n法令番号: %s", data.LawNum)
	lawTitleWithRuby := processTextWithRuby(data.LawBody.LawTitle.Content, data.LawBody.LawTitle.Ruby)
	description += fmt.Sprintf("\n現行法令名: %s %s", lawTitleWithRuby, data.LawBody.LawTitle.Kana)
	book.SetDescription(description)
}

// processChapters processes all chapters in the law
func processChapters(book *epub.Epub, data *jplaw.Law) error {
	for i := range data.LawBody.MainProvision.Chapter {
		if err := processChapter(book, &data.LawBody.MainProvision.Chapter[i], i); err != nil {
			return err
		}
	}
	return nil
}

// processChapter processes a single chapter
func processChapter(book *epub.Epub, chapter *jplaw.Chapter, chapterIdx int) error {
	chapterFilename := fmt.Sprintf("chapter-%d.xhtml", chapterIdx)
	body := buildChapterBody(chapter)

	chapterFilename, err := book.AddSection(body, chapter.ChapterTitle.Content, chapterFilename, "")
	if err != nil {
		return fmt.Errorf("adding chapter: %w", err)
	}

	// Process direct articles under chapter
	if len(chapter.Article) > 0 {
		if err := processArticles(book, chapter.Article, chapterFilename, chapterIdx, -1); err != nil {
			return fmt.Errorf("processing chapter articles: %w", err)
		}
	}

	// Process articles within sections
	for sIdx := range chapter.Section {
		section := &chapter.Section[sIdx]
		if len(section.Article) > 0 {
			if err := processArticles(book, section.Article, chapterFilename, chapterIdx, sIdx); err != nil {
				return fmt.Errorf("processing section articles: %w", err)
			}
		}
	}

	return nil
}

// buildChapterBody builds the HTML body for a chapter
func buildChapterBody(chapter *jplaw.Chapter) string {
	chapterTitleHTML := processTextWithRuby(chapter.ChapterTitle.Content, chapter.ChapterTitle.Ruby)
	body := fmt.Sprintf("<h2>%s</h2>", chapterTitleHTML)

	// Process Sections if any
	if len(chapter.Section) > 0 {
		body += buildSectionsHTML(chapter.Section)
	}

	return body
}

// buildSectionsHTML builds HTML for sections
func buildSectionsHTML(sections []jplaw.Section) string {
	body := "<div class='sections'>"

	for sIdx := range sections {
		section := &sections[sIdx]
		sectionTitleHTML := processTextWithRuby(section.SectionTitle.Content, section.SectionTitle.Ruby)
		body += fmt.Sprintf("<h3>%s</h3>", sectionTitleHTML)

		// Add a note about articles in this section
		if len(section.Article) > 0 {
			body += fmt.Sprintf("<p>（%s から %s まで）</p>",
				section.Article[0].ArticleTitle.Content,
				section.Article[len(section.Article)-1].ArticleTitle.Content)
		}
	}

	body += "</div>"
	return body
}
