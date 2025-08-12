package jplaw2epub

import (
	"fmt"
	"html"

	"go.ngs.io/jplaw-xml"
)

// processItems processes a list of items
func processItems(items []jplaw.Item) string {
	if len(items) == 0 {
		return ""
	}

	// Collect titles for list style
	titles := collectItemTitles(items)
	body := openListWithStyle(titles)

	for i := range items {
		body += processItem(&items[i])
	}

	body += htmlOLEnd
	return body
}

// processItem processes a single item
func processItem(item *jplaw.Item) string {
	body := htmlLI

	// Add item title if not a list number
	if item.ItemTitle != nil && item.ItemTitle.Content != "" && !isListNumber(item.ItemTitle.Content) {
		body += fmt.Sprintf("<strong>%s</strong> ", html.EscapeString(item.ItemTitle.Content))
	}

	// Add item sentences
	for i := range item.ItemSentence.Sentence {
		body += item.ItemSentence.Sentence[i].HTML()
	}

	// Process subitems
	if len(item.Subitem1) > 0 {
		body += processSubitem1List(item.Subitem1)
	}

	body += htmlLIEnd
	return body
}

// processSubitem1List processes a list of Subitem1
func processSubitem1List(subitems []jplaw.Subitem1) string {
	titles := collectSubitem1Titles(subitems)
	body := openListWithStyle(titles)

	for i := range subitems {
		body += processSubitem1(&subitems[i])
	}

	body += htmlOLEnd
	return body
}

// processSubitem1 processes a single Subitem1
func processSubitem1(subitem *jplaw.Subitem1) string {
	body := htmlLI

	// Add title if not a list number
	if subitem.Subitem1Title != nil && subitem.Subitem1Title.Content != "" && !isListNumber(subitem.Subitem1Title.Content) {
		body += fmt.Sprintf("<strong>%s</strong> ", html.EscapeString(subitem.Subitem1Title.Content))
	}

	// Add sentences
	for i := range subitem.Subitem1Sentence.Sentence {
		body += subitem.Subitem1Sentence.Sentence[i].HTML()
	}

	// Process Subitem2
	if len(subitem.Subitem2) > 0 {
		body += processSubitem2List(subitem.Subitem2)
	}

	body += htmlLIEnd
	return body
}

// processSubitem2List processes a list of Subitem2
func processSubitem2List(subitems []jplaw.Subitem2) string {
	titles := collectSubitem2Titles(subitems)
	body := openListWithStyle(titles)

	for i := range subitems {
		body += processSubitem2(&subitems[i])
	}

	body += htmlOLEnd
	return body
}

// processSubitem2 processes a single Subitem2
func processSubitem2(subitem *jplaw.Subitem2) string {
	body := htmlLI

	// Add title if not a list number
	if subitem.Subitem2Title != nil && subitem.Subitem2Title.Content != "" && !isListNumber(subitem.Subitem2Title.Content) {
		body += fmt.Sprintf("<strong>%s</strong> ", html.EscapeString(subitem.Subitem2Title.Content))
	}

	// Add sentences
	for i := range subitem.Subitem2Sentence.Sentence {
		body += subitem.Subitem2Sentence.Sentence[i].HTML()
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
