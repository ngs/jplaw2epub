package jplaw2epub

import (
	"fmt"
	"regexp"
	"strings"

	"go.ngs.io/jplaw-xml"
)

// StyleProcessor handles StyleStruct processing
type StyleProcessor struct {
	imageProcessor *ImageProcessor
}

// NewStyleProcessor creates a new style processor
func NewStyleProcessor(imgProc *ImageProcessor) *StyleProcessor {
	return &StyleProcessor{
		imageProcessor: imgProc,
	}
}

// ProcessStyleStruct processes a StyleStruct and returns HTML
func (sp *StyleProcessor) ProcessStyleStruct(style *jplaw.StyleStruct) string {
	html := `<div class="style-struct">`

	// Add title if present
	if style.StyleStructTitle != nil && style.StyleStructTitle.Content != "" {
		titleHTML := processTextWithRuby(style.StyleStructTitle.Content, style.StyleStructTitle.Ruby)
		html += fmt.Sprintf(`<p class="style-title">%s</p>`, titleHTML)
	}

	// Process Style content
	styleHTML := sp.processStyleContent(style.Style.Content)
	html += styleHTML

	// Add remarks if present
	for i := range style.Remarks {
		remark := &style.Remarks[i]
		html += `<div class="style-remark">`

		// Add remarks label if present
		if remark.RemarksLabel.Content != "" {
			html += fmt.Sprintf(`<p class="remarks-label">%s</p>`,
				processTextWithRuby(remark.RemarksLabel.Content, remark.RemarksLabel.Ruby))
		}

		// Add sentences
		for j := range remark.Sentence {
			html += fmt.Sprintf(`<p>%s</p>`, remark.Sentence[j].HTML())
		}

		// Add items if present
		if len(remark.Item) > 0 {
			html += processItems(remark.Item)
		}

		html += htmlDivEnd
	}

	html += htmlDivEnd
	return html
}

// processStyleContent processes the inner XML content of Style element
func (sp *StyleProcessor) processStyleContent(content string) string {
	// Extract Fig elements from the content
	figPattern := regexp.MustCompile(`<Fig\s+src="([^"]+)"\s*/>`)
	matches := figPattern.FindAllStringSubmatch(content, -1)

	if len(matches) == 0 {
		// No Fig elements, return content as-is (might contain other HTML)
		return fmt.Sprintf(`<div class="style-content">%s</div>`, content)
	}

	// Process each Fig element
	html := ""
	for _, match := range matches {
		if len(match) > 1 {
			src := match[1]

			// Create a temporary FigStruct for processing
			fig := &jplaw.FigStruct{
				Fig: jplaw.Fig{Src: src},
			}

			if sp.imageProcessor != nil {
				if imgHTML, err := sp.imageProcessor.ProcessFigStruct(fig); err == nil {
					html += imgHTML
				}
				// If image processing fails, don't add error text
				// The image itself should be embedded, just skip the error message
			}
			// No image processor available, don't add placeholder text
			// This avoids showing "Image: ./pict/..." text in the EPUB
		}
	}

	// Check if there's any other content besides Fig elements
	remainingContent := figPattern.ReplaceAllString(content, "")
	remainingContent = strings.TrimSpace(remainingContent)
	if remainingContent != "" {
		html += fmt.Sprintf(`<div class="style-content">%s</div>`, remainingContent)
	}

	return html
}

// ProcessStyleStructs processes multiple StyleStructs
func ProcessStyleStructs(styles []jplaw.StyleStruct, imgProc *ImageProcessor) string {
	if len(styles) == 0 {
		return ""
	}

	sp := NewStyleProcessor(imgProc)
	html := ""
	for i := range styles {
		html += sp.ProcessStyleStruct(&styles[i])
	}
	return html
}
