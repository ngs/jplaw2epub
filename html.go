package jplaw2epub

import (
	"fmt"

	"go.ngs.io/jplaw-xml"
)

// HTML constants
const (
	htmlOL    = "<ol>"
	htmlOLEnd = "</ol>"
	htmlLI    = "<li>"
	htmlLIEnd = "</li>"

	// List style types
	listStyleDisc     = "disc"
	listStyleDecimal  = "decimal"
	listStyleCJK      = "cjk-ideographic"
	listStyleKatakana = "katakana-iroha"
	listStyleHiragana = "hiragana-iroha"
)

// openListWithStyle returns an opening list tag with appropriate style
func openListWithStyle(titles []string) string {
	listStyle := getListStyleType(titles)
	if listStyle != "" && listStyle != listStyleDisc {
		return fmt.Sprintf(`<ol style="list-style-type: %s;">`, listStyle)
	}
	return htmlOL
}

// buildArticleTitle builds the full HTML title for an article
func buildArticleTitle(article *jplaw.Article) string {
	articleTitleHTML := processTextWithRuby(article.ArticleTitle.Content, article.ArticleTitle.Ruby)

	if article.ArticleCaption != nil {
		articleCaptionHTML := processTextWithRuby(article.ArticleCaption.Content, article.ArticleCaption.Ruby)
		return fmt.Sprintf("%s %s", articleTitleHTML, articleCaptionHTML)
	}

	return articleTitleHTML
}
