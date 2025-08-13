package jplaw2epub

import (
	"fmt"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

// processSupplProvisions processes supplementary provisions
func processSupplProvisions(book *epub.Epub, provisions []jplaw.SupplProvision, imgProc *ImageProcessor) error {
	if len(provisions) == 0 {
		return nil
	}

	for idx, provision := range provisions {
		if err := processSupplProvision(book, &provision, idx, imgProc); err != nil {
			return fmt.Errorf("processing SupplProvision %d: %w", idx, err)
		}
	}

	return nil
}

// processSupplProvision processes a single supplementary provision
func processSupplProvision(book *epub.Epub, provision *jplaw.SupplProvision, idx int, imgProc *ImageProcessor) error {
	filename := fmt.Sprintf("suppl-provision-%d.xhtml", idx)
	body := ""

	// Add title
	title := "附則"
	if provision.SupplProvisionLabel.Content != "" {
		title = provision.SupplProvisionLabel.Content
	}
	body += fmt.Sprintf(`<div class="chapter-title">%s</div>`, processTextWithRuby(title, provision.SupplProvisionLabel.Ruby))

	// Add amendment law number if present
	if provision.AmendLawNum != "" {
		body += fmt.Sprintf(`<div class="amend-law-num">（%s）</div>`, provision.AmendLawNum)
	}

	// Process chapters if present
	if len(provision.Chapter) > 0 {
		for i := range provision.Chapter {
			chapterTitle := processTextWithRuby(provision.Chapter[i].ChapterTitle.Content, provision.Chapter[i].ChapterTitle.Ruby)
			body += fmt.Sprintf(`<h3>%s</h3>`, chapterTitle)
			
			// Process articles in chapter
			for j := range provision.Chapter[i].Article {
				article := &provision.Chapter[i].Article[j]
				articleTitle := buildArticleTitle(article)
				body += buildArticleBodyWithImages(article, articleTitle, imgProc)
			}
		}
	}

	// Process direct articles
	if len(provision.Article) > 0 {
		for i := range provision.Article {
			article := &provision.Article[i]
			articleTitle := buildArticleTitle(article)
			body += buildArticleBodyWithImages(article, articleTitle, imgProc)
		}
	}

	// Process direct paragraphs
	if len(provision.Paragraph) > 0 {
		body += processParagraphsWithImages(provision.Paragraph, imgProc)
	}

	// Process supplementary provision appendix tables
	if len(provision.SupplProvisionAppdxTable) > 0 {
		for _, table := range provision.SupplProvisionAppdxTable {
			body += processSupplProvisionAppdxTable(&table, imgProc)
		}
	}

	// Process supplementary provision appendix styles
	if len(provision.SupplProvisionAppdxStyle) > 0 {
		for _, style := range provision.SupplProvisionAppdxStyle {
			body += processSupplProvisionAppdxStyle(&style, imgProc)
		}
	}

	// Process supplementary provision appendix
	if len(provision.SupplProvisionAppdx) > 0 {
		for _, appdx := range provision.SupplProvisionAppdx {
			body += processSupplProvisionAppdx(&appdx, imgProc)
		}
	}

	// Add the section to the book
	sectionTitle := title
	if provision.AmendLawNum != "" {
		sectionTitle = fmt.Sprintf("%s（%s）", title, provision.AmendLawNum)
	}
	
	_, err := book.AddSection(body, sectionTitle, filename, "")
	if err != nil {
		return fmt.Errorf("adding SupplProvision section: %w", err)
	}

	return nil
}

// processSupplProvisionAppdxTable processes supplementary provision appendix table
func processSupplProvisionAppdxTable(table *jplaw.SupplProvisionAppdxTable, imgProc *ImageProcessor) string {
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
func processSupplProvisionAppdxStyle(style *jplaw.SupplProvisionAppdxStyle, imgProc *ImageProcessor) string {
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
func processSupplProvisionAppdx(appdx *jplaw.SupplProvisionAppdx, imgProc *ImageProcessor) string {
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