package jplaw2epub

import (
	"fmt"
	"html"
	"strings"

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

// processParagraphWithImages processes a single paragraph with image support
func processParagraphWithImages(para *jplaw.Paragraph, imgProc *ImageProcessor) string {
	p := &paragraphProcessor{imageProcessor: imgProc}

	if para.Num > 0 {
		// For numbered paragraphs, we need to handle them differently
		// Create a single-item list
		p.body += openListWithStyle([]string{para.ParagraphNum.Content})
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

		// Process TableStruct if present
		if len(para.TableStruct) > 0 {
			p.body += processTableStructs(para.TableStruct, p.imageProcessor)
		}

		// Process StyleStruct if present
		if len(para.StyleStruct) > 0 {
			p.body += ProcessStyleStructs(para.StyleStruct, p.imageProcessor)
		}

		// Process List if present
		if len(para.List) > 0 {
			p.body += processLists(para.List)
		}

		p.body += htmlLIEnd
		p.body += htmlOLEnd
	} else {
		p.processRegularParagraph(para)
	}

	return p.body
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

	// Process TableStruct if present
	if len(para.TableStruct) > 0 {
		p.body += processTableStructs(para.TableStruct, p.imageProcessor)
	}

	// Process StyleStruct if present
	if len(para.StyleStruct) > 0 {
		p.body += ProcessStyleStructs(para.StyleStruct, p.imageProcessor)
	}

	// Process List if present
	if len(para.List) > 0 {
		p.body += processLists(para.List)
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

	// Process TableStruct if present
	if len(para.TableStruct) > 0 {
		p.body += processTableStructs(para.TableStruct, p.imageProcessor)
	}

	// Process StyleStruct if present
	if len(para.StyleStruct) > 0 {
		p.body += ProcessStyleStructs(para.StyleStruct, p.imageProcessor)
	}

	// Process List if present
	if len(para.List) > 0 {
		p.body += processLists(para.List)
	}
}

// processLists processes List elements
func processLists(lists []jplaw.List) string {
	if len(lists) == 0 {
		return ""
	}

	var body strings.Builder
	body.WriteString(`<ul class="law-list">`)

	for _, list := range lists {
		body.WriteString("<li>")

		// Process ListSentence
		for i := range list.ListSentence.Sentence {
			body.WriteString(list.ListSentence.Sentence[i].HTML())
		}

		// Process Columns if present
		for i := range list.ListSentence.Column {
			body.WriteString(processColumnElement(&list.ListSentence.Column[i]))
		}

		// Process Sublist1 if present
		if len(list.Sublist1) > 0 {
			body.WriteString(processSublist1(list.Sublist1))
		}

		body.WriteString("</li>")
	}

	body.WriteString("</ul>")
	return body.String()
}

// processSublist1 processes Sublist1 elements
func processSublist1(sublists []jplaw.Sublist1) string {
	if len(sublists) == 0 {
		return ""
	}

	var body strings.Builder
	body.WriteString(`<ul class="law-sublist1">`)

	for _, sublist := range sublists {
		body.WriteString("<li>")

		// Process Sublist1Sentence
		for i := range sublist.Sublist1Sentence.Sentence {
			body.WriteString(sublist.Sublist1Sentence.Sentence[i].HTML())
		}

		// Process Columns if present
		for i := range sublist.Sublist1Sentence.Column {
			body.WriteString(processColumnElement(&sublist.Sublist1Sentence.Column[i]))
		}

		// Process Sublist2 if present (recursive structure)
		if len(sublist.Sublist2) > 0 {
			body.WriteString(processSublist2(sublist.Sublist2))
		}

		body.WriteString("</li>")
	}

	body.WriteString("</ul>")
	return body.String()
}

// processSublist2 processes Sublist2 elements
func processSublist2(sublists []jplaw.Sublist2) string {
	if len(sublists) == 0 {
		return ""
	}

	var body strings.Builder
	body.WriteString(`<ul class="law-sublist2">`)

	for _, sublist := range sublists {
		body.WriteString("<li>")

		// Process Sublist2Sentence
		for i := range sublist.Sublist2Sentence.Sentence {
			body.WriteString(sublist.Sublist2Sentence.Sentence[i].HTML())
		}

		// Process Columns if present
		for i := range sublist.Sublist2Sentence.Column {
			body.WriteString(processColumnElement(&sublist.Sublist2Sentence.Column[i]))
		}

		// Process Sublist3 if present (recursive structure)
		if len(sublist.Sublist3) > 0 {
			body.WriteString(processSublist3(sublist.Sublist3))
		}

		body.WriteString("</li>")
	}

	body.WriteString("</ul>")
	return body.String()
}

// processSublist3 processes Sublist3 elements
func processSublist3(sublists []jplaw.Sublist3) string {
	if len(sublists) == 0 {
		return ""
	}

	var body strings.Builder
	body.WriteString(`<ul class="law-sublist3">`)

	for _, sublist := range sublists {
		body.WriteString("<li>")

		// Process Sublist3Sentence
		for i := range sublist.Sublist3Sentence.Sentence {
			body.WriteString(sublist.Sublist3Sentence.Sentence[i].HTML())
		}

		// Process Columns if present
		for i := range sublist.Sublist3Sentence.Column {
			body.WriteString(processColumnElement(&sublist.Sublist3Sentence.Column[i]))
		}

		body.WriteString("</li>")
	}

	body.WriteString("</ul>")
	return body.String()
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
