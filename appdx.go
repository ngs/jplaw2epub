package jplaw2epub

import (
	"fmt"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

const htmlDivEnd = "</div>"

// processAppdxStyles processes AppdxStyle elements (appendix styles)
func processAppdxStyles(book *epub.Epub, styles []jplaw.AppdxStyle, imgProc *ImageProcessor) error {
	if len(styles) == 0 {
		return nil
	}

	// Create a section for appendix styles
	sectionBody := "<h2>様式</h2>"
	sectionFilename := "appdx-styles.xhtml"

	// Add the section to the book
	sectionFilename, err := book.AddSection(sectionBody, "様式", sectionFilename, "")
	if err != nil {
		return fmt.Errorf("adding appendix styles section: %w", err)
	}

	// Process each AppdxStyle
	for i, style := range styles {
		if err := processAppdxStyle(book, &style, sectionFilename, i, imgProc); err != nil {
			return fmt.Errorf("processing appendix style %d: %w", i, err)
		}
	}

	return nil
}

// processAppdxStyle processes a single AppdxStyle
func processAppdxStyle(book *epub.Epub, style *jplaw.AppdxStyle, parentFilename string, idx int, imgProc *ImageProcessor) error {
	// Build the body content
	body := ""

	// Add title if present
	if style.AppdxStyleTitle != nil && style.AppdxStyleTitle.Content != "" {
		titleHTML := processTextWithRuby(style.AppdxStyleTitle.Content, style.AppdxStyleTitle.Ruby)
		body += fmt.Sprintf("<h3>%s</h3>", titleHTML)
	}

	// Add related article reference if present
	if style.RelatedArticleNum != nil && style.RelatedArticleNum.Content != "" {
		body += fmt.Sprintf("<div class='related-articles'><p>関連条文: %s</p></div>",
			style.RelatedArticleNum.Content)
	}

	// Process StyleStruct elements
	if len(style.StyleStruct) > 0 {
		body += ProcessStyleStructs(style.StyleStruct, imgProc)
	}

	// Process remarks if present
	if style.Remarks != nil {
		body += processAppdxRemark(style.Remarks, imgProc)
	}

	// Create a subsection for this style
	subFilename := fmt.Sprintf("appdx-style-%d.xhtml", idx)
	title := "様式"
	if style.AppdxStyleTitle != nil && style.AppdxStyleTitle.Content != "" {
		title = style.AppdxStyleTitle.Content
	}

	_, err := book.AddSubSection(parentFilename, body, title, subFilename, "")
	if err != nil {
		return fmt.Errorf("adding appendix style subsection: %w", err)
	}

	return nil
}

// processAppdxRemark processes a single remark in appendix
func processAppdxRemark(remark *jplaw.Remarks, imgProc *ImageProcessor) string {
	html := "<div class='appdx-remarks'>"
	html += "<div class='remark'>"

	// Add remarks label if present
	if remark.RemarksLabel.Content != "" {
		html += fmt.Sprintf("<p class='remarks-label'>%s</p>",
			processTextWithRuby(remark.RemarksLabel.Content, remark.RemarksLabel.Ruby))
	}

	// Add sentences
	for j := range remark.Sentence {
		html += fmt.Sprintf("<p>%s</p>", remark.Sentence[j].HTML())
	}

	// Add items if present
	if len(remark.Item) > 0 {
		html += processItemsWithImages(remark.Item, imgProc)
	}

	html += htmlDivEnd
	html += htmlDivEnd
	return html
}

// processAppdxFig processes AppdxFig elements (appendix figures)
func processAppdxFig(book *epub.Epub, figures []jplaw.AppdxFig, imgProc *ImageProcessor) error {
	if len(figures) == 0 {
		return nil
	}

	// Create a section for appendix figures
	sectionBody := "<h2>附図</h2>"
	sectionFilename := "appdx-figures.xhtml"

	// Add the section to the book
	sectionFilename, err := book.AddSection(sectionBody, "附図", sectionFilename, "")
	if err != nil {
		return fmt.Errorf("adding appendix figures section: %w", err)
	}

	// Process each AppdxFig
	for i, fig := range figures {
		if err := processAppdxFigItem(book, &fig, sectionFilename, i, imgProc); err != nil {
			return fmt.Errorf("processing appendix figure %d: %w", i, err)
		}
	}

	return nil
}

// processAppdxFigItem processes a single AppdxFig
func processAppdxFigItem(book *epub.Epub, fig *jplaw.AppdxFig, parentFilename string, idx int, imgProc *ImageProcessor) error {
	// Build the body content
	body := ""

	// Add title if present
	if fig.AppdxFigTitle != nil && fig.AppdxFigTitle.Content != "" {
		titleHTML := processTextWithRuby(fig.AppdxFigTitle.Content, fig.AppdxFigTitle.Ruby)
		body += fmt.Sprintf("<h3>%s</h3>", titleHTML)
	}

	// Process FigStruct elements
	for _, figStruct := range fig.FigStruct {
		if imgProc != nil {
			if html, err := imgProc.ProcessFigStruct(&figStruct); err == nil {
				body += html
			}
		}
	}

	// Process TableStruct elements if present
	for _, table := range fig.TableStruct {
		body += processTableStruct(&table)
	}

	// Create a subsection for this figure
	subFilename := fmt.Sprintf("appdx-fig-%d.xhtml", idx)
	title := "附図"
	if fig.AppdxFigTitle != nil && fig.AppdxFigTitle.Content != "" {
		title = fig.AppdxFigTitle.Content
	}

	_, err := book.AddSubSection(parentFilename, body, title, subFilename, "")
	if err != nil {
		return fmt.Errorf("adding appendix figure subsection: %w", err)
	}

	return nil
}

// processTableStruct processes a TableStruct (simplified for now)
func processTableStruct(table *jplaw.TableStruct) string {
	html := "<div class='table-struct'>"

	// Add title if present
	if table.TableStructTitle != nil && table.TableStructTitle.Content != "" {
		titleHTML := processTextWithRuby(table.TableStructTitle.Content, table.TableStructTitle.Ruby)
		html += fmt.Sprintf("<p class='table-title'>%s</p>", titleHTML)
	}

	// For now, just indicate that a table exists
	// Full table processing would require parsing the Table structure
	html += "<p class='table-placeholder'>[表]</p>"

	html += htmlDivEnd
	return html
}
