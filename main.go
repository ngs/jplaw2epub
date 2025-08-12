package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

func main() {
	destPath, sourcePath := parseFlags()
	data := loadXMLData(sourcePath)
	book := createEPUB(data)
	processChapters(book, data)
	writeEPUB(book, destPath)
}

// parseFlags parses command line flags and returns destination and source paths
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

// loadXMLData loads and unmarshals XML data from file
func loadXMLData(sourcePath string) *jplaw.Law {
	xmlFile, err := os.Open(sourcePath)
	if err != nil {
		handleError("opening source file", err)
	}
	defer xmlFile.Close()

	byteValue, err := io.ReadAll(xmlFile)
	if err != nil {
		handleError("reading source file", err)
	}

	var data jplaw.Law
	if err := xml.Unmarshal(byteValue, &data); err != nil {
		handleError("unmarshalling XML", err)
	}

	return &data
}

// createEPUB creates and sets up EPUB
func createEPUB(data *jplaw.Law) *epub.Epub {
	book, err := epub.NewEpub(data.LawBody.LawTitle.Content)
	if err != nil {
		handleError("creating epub", err)
	}

	setupEPUBMetadata(book, data)
	return book
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
func processChapters(book *epub.Epub, data *jplaw.Law) {
	for i := range data.LawBody.MainProvision.Chapter {
		processChapter(book, &data.LawBody.MainProvision.Chapter[i], i)
	}
}

// processChapter processes a single chapter
func processChapter(book *epub.Epub, chapter *jplaw.Chapter, chapterIdx int) {
	chapterFilename := fmt.Sprintf("chapter-%d.xhtml", chapterIdx)
	body := buildChapterBody(chapter)

	chapterFilename, err := book.AddSection(body, chapter.ChapterTitle.Content, chapterFilename, "")
	if err != nil {
		handleError("adding chapter", err)
	}

	// Process direct articles under chapter
	if len(chapter.Article) > 0 {
		if err := processArticles(book, chapter.Article, chapterFilename, chapterIdx, -1); err != nil {
			handleError("processing chapter articles", err)
		}
	}

	// Process articles within sections
	for sIdx := range chapter.Section {
		section := &chapter.Section[sIdx]
		if len(section.Article) > 0 {
			if err := processArticles(book, section.Article, chapterFilename, chapterIdx, sIdx); err != nil {
				handleError("processing section articles", err)
			}
		}
	}
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

// writeEPUB writes the EPUB to file
func writeEPUB(book *epub.Epub, destPath string) {
	if err := book.Write(destPath); err != nil {
		handleError("writing epub", err)
	}
}

// handleError handles errors uniformly
func handleError(context string, err error) {
	fmt.Printf("Error %s: %v\n", context, err)
	os.Exit(1)
}
