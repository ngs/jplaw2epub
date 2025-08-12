package jplaw2epub

import (
	"fmt"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

// processArticles processes a slice of articles and adds them to the EPUB
func processArticles(book *epub.Epub, articles []jplaw.Article, parentFilename string, chapterIdx, sectionIdx int) error {
	for j := range articles {
		if err := processArticle(book, &articles[j], parentFilename, chapterIdx, sectionIdx, j); err != nil {
			return err
		}
	}
	return nil
}

// processArticle processes a single article
func processArticle(book *epub.Epub, article *jplaw.Article, parentFilename string, chapterIdx, sectionIdx, articleIdx int) error {
	subFilename := buildArticleFilename(chapterIdx, sectionIdx, articleIdx)
	articleTitle := buildArticleTitle(article)
	body := buildArticleBody(article, articleTitle)

	articleTitlePlain := getArticleTitlePlain(article)
	_, err := book.AddSubSection(parentFilename, body, articleTitlePlain, subFilename, "")
	if err != nil {
		return fmt.Errorf("error adding article section: %w", err)
	}
	return nil
}

// buildArticleFilename generates the filename for an article
func buildArticleFilename(chapterIdx, sectionIdx, articleIdx int) string {
	if sectionIdx >= 0 {
		return fmt.Sprintf("article-%d-%d-%d.xhtml", chapterIdx, sectionIdx, articleIdx)
	}
	return fmt.Sprintf("article-%d-%d.xhtml", chapterIdx, articleIdx)
}

// buildArticleBody builds the HTML body for an article
func buildArticleBody(article *jplaw.Article, articleTitle string) string {
	body := fmt.Sprintf("<h3>%s</h3>", articleTitle)
	body += processParagraphs(article.Paragraph)
	return body
}

// getArticleTitlePlain returns the plain text title for an article
func getArticleTitlePlain(article *jplaw.Article) string {
	if article.ArticleCaption != nil {
		return fmt.Sprintf("%s %s", article.ArticleTitle.Content, article.ArticleCaption.Content)
	}
	return article.ArticleTitle.Content
}
