package jplaw2epub

import (
	"fmt"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

const defaultSupplProvisionTitle = "附則"

// processSupplProvisions processes supplementary provisions
func processSupplProvisions(book *epub.Epub, provisions []jplaw.SupplProvision, imgProc ImageProcessorInterface) error {
	if len(provisions) == 0 {
		return nil
	}

	for idx := range provisions {
		if err := processSupplProvision(book, &provisions[idx], idx, imgProc); err != nil {
			return fmt.Errorf("processing SupplProvision %d: %w", idx, err)
		}
	}

	return nil
}

// processSupplProvision processes a single supplementary provision
func processSupplProvision(book *epub.Epub, provision *jplaw.SupplProvision, idx int, imgProc ImageProcessorInterface) error {
	filename := fmt.Sprintf("suppl-provision-%d.xhtml", idx)

	// Build the body content
	body := buildSupplProvisionBody(provision, imgProc)

	// Get the section title
	title := getSupplProvisionTitle(provision)
	sectionTitle := title
	if provision.AmendLawNum != "" {
		sectionTitle = fmt.Sprintf("%s（%s）", title, provision.AmendLawNum)
	}

	// Add the section to the book
	_, err := book.AddSection(body, sectionTitle, filename, "")
	if err != nil {
		return fmt.Errorf("adding SupplProvision section: %w", err)
	}

	return nil
}

// buildSupplProvisionBody builds the HTML body for a supplementary provision
func buildSupplProvisionBody(provision *jplaw.SupplProvision, imgProc ImageProcessorInterface) string {
	var body string

	// Add title
	title := getSupplProvisionTitle(provision)
	body += fmt.Sprintf(`<div class="chapter-title">%s</div>`, processTextWithRuby(title, provision.SupplProvisionLabel.Ruby))

	// Add amendment law number if present
	if provision.AmendLawNum != "" {
		body += fmt.Sprintf(`<div class="amend-law-num">（%s）</div>`, provision.AmendLawNum)
	}

	// Process chapters
	body += processSupplProvisionChapters(provision, imgProc)

	// Process direct articles
	body += processSupplProvisionArticles(provision.Article, imgProc)

	// Process direct paragraphs
	if len(provision.Paragraph) > 0 {
		body += processParagraphsWithImages(provision.Paragraph, imgProc)
	}

	// Process appendixes
	body += processSupplProvisionAppendixes(provision, imgProc)

	return body
}

// getSupplProvisionTitle gets the title for a supplementary provision
func getSupplProvisionTitle(provision *jplaw.SupplProvision) string {
	if provision.SupplProvisionLabel.Content != "" {
		return provision.SupplProvisionLabel.Content
	}
	return defaultSupplProvisionTitle
}

// processSupplProvisionChapters processes chapters in a supplementary provision
func processSupplProvisionChapters(provision *jplaw.SupplProvision, imgProc ImageProcessorInterface) string {
	if len(provision.Chapter) == 0 {
		return ""
	}

	var body string
	for i := range provision.Chapter {
		chapterTitle := processTextWithRuby(provision.Chapter[i].ChapterTitle.Content, provision.Chapter[i].ChapterTitle.Ruby)
		body += fmt.Sprintf(`<h3>%s</h3>`, chapterTitle)
		body += processSupplProvisionArticles(provision.Chapter[i].Article, imgProc)
	}
	return body
}

// processSupplProvisionArticles processes articles
func processSupplProvisionArticles(articles []jplaw.Article, imgProc ImageProcessorInterface) string {
	if len(articles) == 0 {
		return ""
	}

	var body string
	for i := range articles {
		article := &articles[i]
		articleTitle := buildArticleTitle(article)
		body += buildArticleBodyWithImages(article, articleTitle, imgProc)
	}
	return body
}

// processSupplProvisionAppendixes processes all appendix types
func processSupplProvisionAppendixes(provision *jplaw.SupplProvision, imgProc ImageProcessorInterface) string {
	var body string

	// Process appendix tables
	for i := range provision.SupplProvisionAppdxTable {
		body += processSupplProvisionAppdxTable(&provision.SupplProvisionAppdxTable[i], imgProc)
	}

	// Process appendix styles
	for i := range provision.SupplProvisionAppdxStyle {
		body += processSupplProvisionAppdxStyle(&provision.SupplProvisionAppdxStyle[i], imgProc)
	}

	// Process supplementary provision appendix
	for i := range provision.SupplProvisionAppdx {
		body += processSupplProvisionAppdx(&provision.SupplProvisionAppdx[i], imgProc)
	}

	return body
}

// processSupplProvisionAppdxTable processes supplementary provision appendix table
func processSupplProvisionAppdxTable(table *jplaw.SupplProvisionAppdxTable, imgProc ImageProcessorInterface) string {
	body := `<div class="suppl-appdx-table">`

	// Add title if present
	if table.SupplProvisionAppdxTableTitle.Content != "" {
		body += fmt.Sprintf(`<h4>%s</h4>`,
			processTextWithRuby(table.SupplProvisionAppdxTableTitle.Content, table.SupplProvisionAppdxTableTitle.Ruby))
	}

	// Add related article number if present
	if table.RelatedArticleNum != nil && table.RelatedArticleNum.Content != "" {
		body += fmt.Sprintf(`<div class="related-articles">%s</div>`,
			processTextWithRuby(table.RelatedArticleNum.Content, table.RelatedArticleNum.Ruby))
	}

	// Process TableStructs
	for _, tableStruct := range table.TableStruct {
		body += processTableStructWithImages(&tableStruct, imgProc)
	}

	body += htmlDivEnd
	return body
}

// processSupplProvisionAppdxStyle processes supplementary provision appendix style
func processSupplProvisionAppdxStyle(style *jplaw.SupplProvisionAppdxStyle, imgProc ImageProcessorInterface) string {
	body := `<div class="suppl-appdx-style">`

	// Add title if present
	if style.SupplProvisionAppdxStyleTitle.Content != "" {
		body += fmt.Sprintf(`<h4>%s</h4>`,
			processTextWithRuby(style.SupplProvisionAppdxStyleTitle.Content, style.SupplProvisionAppdxStyleTitle.Ruby))
	}

	// Add related article number if present
	if style.RelatedArticleNum != nil && style.RelatedArticleNum.Content != "" {
		body += fmt.Sprintf(`<div class="related-articles">%s</div>`,
			processTextWithRuby(style.RelatedArticleNum.Content, style.RelatedArticleNum.Ruby))
	}

	// Process StyleStructs
	if len(style.StyleStruct) > 0 {
		body += ProcessStyleStructs(style.StyleStruct, imgProc)
	}

	body += htmlDivEnd
	return body
}

// processSupplProvisionAppdx processes supplementary provision appendix
func processSupplProvisionAppdx(appdx *jplaw.SupplProvisionAppdx, _ ImageProcessorInterface) string {
	body := `<div class="suppl-appdx">`

	// Add arithmetic formula number if present
	if appdx.ArithFormulaNum != nil && appdx.ArithFormulaNum.Content != "" {
		body += fmt.Sprintf(`<div class="arith-formula-num">%s</div>`,
			processTextWithRuby(appdx.ArithFormulaNum.Content, appdx.ArithFormulaNum.Ruby))
	}

	// Add related article number if present
	if appdx.RelatedArticleNum != nil && appdx.RelatedArticleNum.Content != "" {
		body += fmt.Sprintf(`<div class="related-articles">%s</div>`,
			processTextWithRuby(appdx.RelatedArticleNum.Content, appdx.RelatedArticleNum.Ruby))
	}

	// Process ArithFormula
	for _, formula := range appdx.ArithFormula {
		body += `<div class="arith-formula">`
		if formula.Num != 0 {
			body += fmt.Sprintf(`<span class="formula-num">(%d)</span>`, formula.Num)
		}
		// ArithFormula contains complex content, for now just display as text
		body += `<span class="formula-content">[算式]</span>`
		body += htmlDivEnd
	}

	body += htmlDivEnd
	return body
}
