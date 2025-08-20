package jplaw2epub

import (
	"fmt"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

// processAppdxFormats processes appendix formats
func processAppdxFormats(book *epub.Epub, formats []jplaw.AppdxFormat, imgProc ImageProcessorInterface) error {
	if len(formats) == 0 {
		return nil
	}

	for idx, format := range formats {
		if err := processAppdxFormat(book, &format, idx, imgProc); err != nil {
			return fmt.Errorf("processing AppdxFormat %d: %w", idx, err)
		}
	}

	return nil
}

// processAppdxFormat processes a single appendix format
func processAppdxFormat(book *epub.Epub, format *jplaw.AppdxFormat, idx int, imgProc ImageProcessorInterface) error {
	filename := fmt.Sprintf("appdx-format-%d.xhtml", idx)
	body := ""

	// Add title if present
	title := "書式"
	if format.AppdxFormatTitle != nil && format.AppdxFormatTitle.Content != "" {
		title = format.AppdxFormatTitle.Content
		body += fmt.Sprintf(`<div class="chapter-title">%s</div>`, processTextWithRuby(title, format.AppdxFormatTitle.Ruby))
	}

	// Process related article number if present
	if format.RelatedArticleNum != nil && format.RelatedArticleNum.Content != "" {
		body += fmt.Sprintf(`<div class="related-articles">%s</div>`,
			processTextWithRuby(format.RelatedArticleNum.Content, format.RelatedArticleNum.Ruby))
	}

	// Process FormatStructs
	for _, formatStruct := range format.FormatStruct {
		body += processFormatStruct(&formatStruct, imgProc)
	}

	// Add the section to the book
	_, err := book.AddSection(body, title, filename, "")
	if err != nil {
		return fmt.Errorf("adding AppdxFormat section: %w", err)
	}

	return nil
}

// processFormatStruct processes a FormatStruct
func processFormatStruct(formatStruct *jplaw.FormatStruct, imgProc ImageProcessorInterface) string {
	body := `<div class="format-struct">`

	// Add title if present
	if formatStruct.FormatStructTitle != nil && formatStruct.FormatStructTitle.Content != "" {
		body += fmt.Sprintf(`<h3>%s</h3>`,
			processTextWithRuby(formatStruct.FormatStructTitle.Content, formatStruct.FormatStructTitle.Ruby))
	}

	// Process Format content - it's raw XML content
	body += processFormat(&formatStruct.Format, imgProc)

	// Process Remarks
	for i := range formatStruct.Remarks {
		body += processRemarks(&formatStruct.Remarks[i])
	}

	body += htmlDivEnd
	return body
}

// processFormat processes a Format element which contains raw XML
func processFormat(format *jplaw.Format, imgProc ImageProcessorInterface) string {
	body := `<div class="format-content">`

	// Format.Content contains raw XML that may include Fig elements
	// For now, we'll display it as-is or try to extract Fig elements
	if format.Content != "" {
		// Check if it contains Fig elements
		if imgProc != nil && containsFigElement(format.Content) {
			// Try to process embedded Fig elements
			processedContent := processEmbeddedFigs(format.Content, imgProc)
			body += processedContent
		} else {
			// Display as preformatted text
			body += fmt.Sprintf(`<pre class="format-raw">%s</pre>`, format.Content)
		}
	}

	body += htmlDivEnd
	return body
}

// containsFigElement checks if the content contains Fig elements
func containsFigElement(content string) bool {
	return len(content) > 4 && content[:4] == "<Fig"
}

// processEmbeddedFigs processes Fig elements embedded in XML content
func processEmbeddedFigs(content string, _ ImageProcessorInterface) string {
	// For simplicity, just display the content
	// In a real implementation, we'd parse the XML and extract Fig elements
	return fmt.Sprintf(`<div class="embedded-content">%s</div>`, content)
}
