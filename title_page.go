package jplaw2epub

import (
	"fmt"
	"html"
	"strings"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

// addTitlePage adds a title page as the first page of the EPUB
func addTitlePage(book *epub.Epub, data *jplaw.Law) error {
	// Build title page content
	var body strings.Builder
	body.WriteString(`<div style="text-align: center; margin-top: 20%;">`)

	// Law title with ruby if available
	body.WriteString(`<h1 style="font-size: 1.5em; margin-bottom: 1em;">`)
	if data.LawBody.LawTitle.Ruby != nil {
		body.WriteString(processTextWithRuby(data.LawBody.LawTitle.Content, data.LawBody.LawTitle.Ruby))
	} else {
		body.WriteString(html.EscapeString(data.LawBody.LawTitle.Content))
	}
	body.WriteString(`</h1>`)

	// Law number
	body.WriteString(`<p style="font-size: 1.2em; margin-bottom: 2em;">`)
	body.WriteString(html.EscapeString(data.LawNum))
	body.WriteString(`</p>`)

	// Promulgation date
	eraStr := getEraString(data.Era)
	body.WriteString(`<p style="margin-bottom: 0.5em;">`)
	body.WriteString(fmt.Sprintf("公布日: %s%d年%d月%d日", eraStr, data.Year, data.PromulgateMonth, data.PromulgateDay))
	body.WriteString(`</p>`)

	// Enact statement if present
	if len(data.LawBody.EnactStatement) > 0 && data.LawBody.EnactStatement[0].Content != "" {
		body.WriteString(`<div style="margin-top: 3em; text-align: left; padding: 0 10%;">`)
		body.WriteString(`<p style="text-indent: 1em;">`)
		enactStmt := &data.LawBody.EnactStatement[0]
		if len(enactStmt.Ruby) > 0 {
			body.WriteString(processTextWithRuby(enactStmt.Content, enactStmt.Ruby))
		} else {
			body.WriteString(html.EscapeString(enactStmt.Content))
		}
		body.WriteString(`</p>`)
		body.WriteString(`</div>`)
	}

	body.WriteString(`</div>`)

	// Add the title page as the first section
	_, err := book.AddSection(body.String(), "タイトルページ", "title.xhtml", "")
	if err != nil {
		return fmt.Errorf("adding title page section: %w", err)
	}

	return nil
}
