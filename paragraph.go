package jplaw2epub

import (
	"fmt"
	"html"

	"go.ngs.io/jplaw-xml"
)

// paragraphProcessor handles paragraph processing state
type paragraphProcessor struct {
	inList         bool
	body           string
	imageProcessor *ImageProcessor
}

// processParagraphs processes all paragraphs in an article
func processParagraphs(paragraphs []jplaw.Paragraph) string {
	return processParagraphsWithImages(paragraphs, nil)
}

// processParagraphsWithImages processes all paragraphs with image support
func processParagraphsWithImages(paragraphs []jplaw.Paragraph, imgProc *ImageProcessor) string {
	p := &paragraphProcessor{imageProcessor: imgProc}

	for i := range paragraphs {
		para := &paragraphs[i]
		if para.Num > 0 {
			p.processNumberedParagraph(para, i, paragraphs)
		} else {
			p.processRegularParagraph(para)
		}
	}

	// Close list if still open
	if p.inList {
		p.body += htmlOLEnd
	}

	return p.body
}

// processNumberedParagraph handles paragraphs with Num attribute
func (p *paragraphProcessor) processNumberedParagraph(para *jplaw.Paragraph, idx int, allParagraphs []jplaw.Paragraph) {
	// Start a new list if not in one
	if !p.inList {
		p.startNumberedList(idx, allParagraphs)
		p.inList = true
	}

	p.body += htmlLI
	p.addParagraphNumber(para)
	p.addParagraphSentences(para)

	if len(para.Item) > 0 {
		p.body += processItemsWithImages(para.Item, p.imageProcessor)
	}

	// Process FigStruct if present
	if len(para.FigStruct) > 0 {
		for _, fig := range para.FigStruct {
			if p.imageProcessor != nil {
				if html, err := p.imageProcessor.ProcessFigStruct(&fig); err == nil {
					p.body += html
				}
			}
		}
	}

	// Process StyleStruct if present
	if len(para.StyleStruct) > 0 {
		p.body += ProcessStyleStructs(para.StyleStruct, p.imageProcessor)
	}

	p.body += htmlLIEnd
}

// processRegularParagraph handles regular paragraphs without Num attribute
func (p *paragraphProcessor) processRegularParagraph(para *jplaw.Paragraph) {
	// Close list if we were in one
	if p.inList {
		p.body += htmlOLEnd
		p.inList = false
	}

	if para.ParagraphNum.Content != "" {
		p.body += fmt.Sprintf("<h4>%s</h4>", html.EscapeString(para.ParagraphNum.Content))
	}

	if len(para.ParagraphSentence.Sentence) > 0 {
		p.body += "<p>"
		for i := range para.ParagraphSentence.Sentence {
			p.body += para.ParagraphSentence.Sentence[i].HTML()
		}
		p.body += "</p>"
	}

	if len(para.Item) > 0 {
		p.body += processItemsWithImages(para.Item, p.imageProcessor)
	}

	// Process FigStruct if present
	if len(para.FigStruct) > 0 {
		for _, fig := range para.FigStruct {
			if p.imageProcessor != nil {
				if html, err := p.imageProcessor.ProcessFigStruct(&fig); err == nil {
					p.body += html
				}
			}
		}
	}

	// Process StyleStruct if present
	if len(para.StyleStruct) > 0 {
		p.body += ProcessStyleStructs(para.StyleStruct, p.imageProcessor)
	}
}

// startNumberedList starts a new numbered list with appropriate style
func (p *paragraphProcessor) startNumberedList(idx int, paragraphs []jplaw.Paragraph) {
	var titles []string
	for i := idx; i < len(paragraphs) && paragraphs[i].Num > 0; i++ {
		titles = append(titles, paragraphs[i].ParagraphNum.Content)
	}
	p.body += openListWithStyle(titles)
}

// addParagraphNumber adds paragraph number if not a list number
func (p *paragraphProcessor) addParagraphNumber(para *jplaw.Paragraph) {
	if para.ParagraphNum.Content != "" && !isListNumber(para.ParagraphNum.Content) {
		p.body += fmt.Sprintf("<strong>%s</strong> ", html.EscapeString(para.ParagraphNum.Content))
	}
}

// addParagraphSentences adds paragraph sentences
func (p *paragraphProcessor) addParagraphSentences(para *jplaw.Paragraph) {
	for i := range para.ParagraphSentence.Sentence {
		p.body += para.ParagraphSentence.Sentence[i].HTML()
	}
}
