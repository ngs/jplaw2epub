package jplaw2epub

import (
	"fmt"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

// processChapterWithImages processes a single chapter with image support
func processChapterWithImages(book *epub.Epub, chapter *jplaw.Chapter, chapterIdx int, imgProc ImageProcessorInterface) error {
	chapterFilename := fmt.Sprintf("chapter-%d.xhtml", chapterIdx)
	body := buildChapterBody(chapter)

	chapterFilename, err := book.AddSection(body, chapter.ChapterTitle.Content, chapterFilename, "")
	if err != nil {
		return fmt.Errorf("adding chapter: %w", err)
	}

	// Process direct articles under chapter
	if len(chapter.Article) > 0 {
		if err := processArticlesWithImages(book, chapter.Article, chapterFilename, chapterIdx, -1, imgProc); err != nil {
			return fmt.Errorf("processing chapter articles: %w", err)
		}
	}

	// Process articles within sections
	for sIdx := range chapter.Section {
		section := &chapter.Section[sIdx]
		if len(section.Article) > 0 {
			if err := processArticlesWithImages(book, section.Article, chapterFilename, chapterIdx, sIdx, imgProc); err != nil {
				return fmt.Errorf("processing section articles: %w", err)
			}
		}
	}

	return nil
}

// buildChapterBody builds the HTML body for a chapter
func buildChapterBody(chapter *jplaw.Chapter) string {
	chapterTitleHTML := processTextWithRuby(chapter.ChapterTitle.Content, chapter.ChapterTitle.Ruby)
	body := fmt.Sprintf(`<div class="chapter-title">%s</div>`, chapterTitleHTML)

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

	body += htmlDivEnd
	return body
}
