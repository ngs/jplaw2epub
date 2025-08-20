package jplaw2epub

import (
	"fmt"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

// processMainProvision processes the main provision content
func processMainProvision(book *epub.Epub, mainProv *jplaw.MainProvision, imgProc *ImageProcessor) error {
	if len(mainProv.Chapter) > 0 {
		// Process chapters
		for i := range mainProv.Chapter {
			if err := processChapterWithImages(book, &mainProv.Chapter[i], i, imgProc); err != nil {
				return err
			}
		}
		return nil
	}

	// No chapters, check for articles
	if len(mainProv.Article) > 0 {
		// Process articles as separate sections for TOC
		for i := range mainProv.Article {
			article := &mainProv.Article[i]
			articleFilename := fmt.Sprintf("article-%d.xhtml", i)
			articleTitle := buildArticleTitle(article)
			body := buildArticleBodyWithImages(article, articleTitle, imgProc)

			articleTitlePlain := getArticleTitlePlain(article)
			_, err := book.AddSection(body, articleTitlePlain, articleFilename, "")
			if err != nil {
				return fmt.Errorf("adding article section: %w", err)
			}
		}
		return nil
	}

	// No chapters or articles, process direct paragraphs
	if len(mainProv.Paragraph) == 0 {
		return nil
	}

	// If there are multiple paragraphs, create separate sections for better TOC
	if len(mainProv.Paragraph) > 1 {
		for i := range mainProv.Paragraph {
			paragraph := &mainProv.Paragraph[i]
			paragraphFilename := fmt.Sprintf("paragraph-%d.xhtml", i)

			// Create title for the paragraph
			paragraphTitle := fmt.Sprintf("第%s項", paragraph.ParagraphNum.Content)
			if paragraph.ParagraphNum.Content == "" && paragraph.Num != 0 {
				paragraphTitle = fmt.Sprintf("第%d項", paragraph.Num)
			}

			// Build paragraph body
			body := fmt.Sprintf("<h3>%s</h3>", paragraphTitle)
			body += processParagraphWithImages(paragraph, imgProc)

			_, err := book.AddSection(body, paragraphTitle, paragraphFilename, "")
			if err != nil {
				return fmt.Errorf("adding paragraph section: %w", err)
			}
		}
		return nil
	}

	// Single paragraph - add as main content
	mainFilename := "main-content.xhtml"
	body := processParagraphsWithImages(mainProv.Paragraph, imgProc)

	if body != "" {
		_, err := book.AddSection(body, "本文", mainFilename, "")
		if err != nil {
			return fmt.Errorf("adding main content: %w", err)
		}
	}

	return nil
}
