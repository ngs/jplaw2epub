package jplaw2epub

import (
	"fmt"
	"html"

	"go.ngs.io/jplaw-xml"
)

// processItems processes a list of items
func processItems(items []jplaw.Item) string {
	return processItemsWithImages(items, nil)
}

// processItemsWithImages processes a list of items with image support
func processItemsWithImages(items []jplaw.Item, imgProc *ImageProcessor) string {
	if len(items) == 0 {
		return ""
	}

	// Collect titles for list style
	titles := collectItemTitles(items)
	body := openListWithStyle(titles)

	for i := range items {
		body += processItemWithImages(&items[i], imgProc)
	}

	body += htmlOLEnd
	return body
}

// processItem processes a single item
func processItem(item *jplaw.Item) string {
	return processItemWithImages(item, nil)
}

// processItemWithImages processes a single item with image support
func processItemWithImages(item *jplaw.Item, imgProc *ImageProcessor) string {
	body := htmlLI

	// Add item title if not a list number
	if item.ItemTitle != nil && item.ItemTitle.Content != "" && !isListNumber(item.ItemTitle.Content) {
		body += fmt.Sprintf("<strong>%s</strong> ", html.EscapeString(item.ItemTitle.Content))
	}

	// Add item sentences
	for i := range item.ItemSentence.Sentence {
		body += item.ItemSentence.Sentence[i].HTML()
	}

	// Process FigStruct if present
	if len(item.FigStruct) > 0 {
		for _, fig := range item.FigStruct {
			if imgProc != nil {
				if html, err := imgProc.ProcessFigStruct(&fig); err == nil {
					body += html
				}
			}
		}
	}

	// Process StyleStruct if present
	if len(item.StyleStruct) > 0 {
		body += ProcessStyleStructs(item.StyleStruct, imgProc)
	}

	// Process subitems
	if len(item.Subitem1) > 0 {
		body += processSubitem1ListWithImages(item.Subitem1, imgProc)
	}

	body += htmlLIEnd
	return body
}

// processSubitem1List processes a list of Subitem1
func processSubitem1List(subitems []jplaw.Subitem1) string {
	return processSubitem1ListWithImages(subitems, nil)
}

// processSubitem1ListWithImages processes a list of Subitem1 with image support
func processSubitem1ListWithImages(subitems []jplaw.Subitem1, imgProc *ImageProcessor) string {
	titles := collectSubitem1Titles(subitems)
	body := openListWithStyle(titles)

	for i := range subitems {
		body += processSubitem1WithImages(&subitems[i], imgProc)
	}

	body += htmlOLEnd
	return body
}

// processSubitem1 processes a single Subitem1
func processSubitem1(subitem *jplaw.Subitem1) string {
	return processSubitem1WithImages(subitem, nil)
}

// processSubitem1WithImages processes a single Subitem1 with image support
func processSubitem1WithImages(subitem *jplaw.Subitem1, imgProc *ImageProcessor) string {
	body := htmlLI

	// Add title if not a list number
	if subitem.Subitem1Title != nil && subitem.Subitem1Title.Content != "" && !isListNumber(subitem.Subitem1Title.Content) {
		body += fmt.Sprintf("<strong>%s</strong> ", html.EscapeString(subitem.Subitem1Title.Content))
	}

	// Add sentences
	for i := range subitem.Subitem1Sentence.Sentence {
		body += subitem.Subitem1Sentence.Sentence[i].HTML()
	}

	// Process FigStruct if present
	if len(subitem.FigStruct) > 0 {
		for _, fig := range subitem.FigStruct {
			if imgProc != nil {
				if html, err := imgProc.ProcessFigStruct(&fig); err == nil {
					body += html
				}
			}
		}
	}

	// Process StyleStruct if present
	if len(subitem.StyleStruct) > 0 {
		body += ProcessStyleStructs(subitem.StyleStruct, imgProc)
	}

	// Process Subitem2
	if len(subitem.Subitem2) > 0 {
		body += processSubitem2ListWithImages(subitem.Subitem2, imgProc)
	}

	body += htmlLIEnd
	return body
}

// processSubitem2List processes a list of Subitem2
func processSubitem2List(subitems []jplaw.Subitem2) string {
	return processSubitem2ListWithImages(subitems, nil)
}

// processSubitem2ListWithImages processes a list of Subitem2 with image support
func processSubitem2ListWithImages(subitems []jplaw.Subitem2, imgProc *ImageProcessor) string {
	titles := collectSubitem2Titles(subitems)
	body := openListWithStyle(titles)

	for i := range subitems {
		body += processSubitem2WithImages(&subitems[i], imgProc)
	}

	body += htmlOLEnd
	return body
}

// processSubitem2 processes a single Subitem2
func processSubitem2(subitem *jplaw.Subitem2) string {
	return processSubitem2WithImages(subitem, nil)
}

// processSubitem2WithImages processes a single Subitem2 with image support
func processSubitem2WithImages(subitem *jplaw.Subitem2, imgProc *ImageProcessor) string {
	body := htmlLI

	// Add title if not a list number
	if subitem.Subitem2Title != nil && subitem.Subitem2Title.Content != "" && !isListNumber(subitem.Subitem2Title.Content) {
		body += fmt.Sprintf("<strong>%s</strong> ", html.EscapeString(subitem.Subitem2Title.Content))
	}

	// Add sentences
	for i := range subitem.Subitem2Sentence.Sentence {
		body += subitem.Subitem2Sentence.Sentence[i].HTML()
	}

	// Process FigStruct if present
	if len(subitem.FigStruct) > 0 {
		for _, fig := range subitem.FigStruct {
			if imgProc != nil {
				if html, err := imgProc.ProcessFigStruct(&fig); err == nil {
					body += html
				}
			}
		}
	}

	// Process StyleStruct if present
	if len(subitem.StyleStruct) > 0 {
		body += ProcessStyleStructs(subitem.StyleStruct, imgProc)
	}

	body += htmlLIEnd
	return body
}

// collectItemTitles collects titles from items
func collectItemTitles(items []jplaw.Item) []string {
	var titles []string
	for i := range items {
		if items[i].ItemTitle != nil {
			titles = append(titles, items[i].ItemTitle.Content)
		}
	}
	return titles
}

// collectSubitem1Titles collects titles from Subitem1
func collectSubitem1Titles(subitems []jplaw.Subitem1) []string {
	var titles []string
	for i := range subitems {
		if subitems[i].Subitem1Title != nil {
			titles = append(titles, subitems[i].Subitem1Title.Content)
		}
	}
	return titles
}

// collectSubitem2Titles collects titles from Subitem2
func collectSubitem2Titles(subitems []jplaw.Subitem2) []string {
	var titles []string
	for i := range subitems {
		if subitems[i].Subitem2Title != nil {
			titles = append(titles, subitems[i].Subitem2Title.Content)
		}
	}
	return titles
}
