package jplaw2epub

import (
	"fmt"
	"html"
	"strconv"
	"strings"

	jplaw "go.ngs.io/jplaw-xml"
)

// processTableStructs processes multiple table structures
func processTableStructs(tables []jplaw.TableStruct, imgProc *ImageProcessor) string {
	if len(tables) == 0 {
		return ""
	}

	var body strings.Builder
	for i := range tables {
		body.WriteString(processTableStructWithImages(&tables[i], imgProc))
	}
	return body.String()
}

// processTableStructWithImages processes a single table structure with image support
func processTableStructWithImages(tableStruct *jplaw.TableStruct, imgProc *ImageProcessor) string {
	var body strings.Builder

	// Add table title if present
	if tableStruct.TableStructTitle != nil {
		titleHTML := processTextWithRuby(
			tableStruct.TableStructTitle.Content,
			tableStruct.TableStructTitle.Ruby,
		)
		body.WriteString(fmt.Sprintf(`<div class="table-title">%s</div>`, titleHTML))
	}

	// Process the table
	body.WriteString(processTable(&tableStruct.Table))

	// Process remarks if present
	for i := range tableStruct.Remarks {
		body.WriteString(processRemarks(&tableStruct.Remarks[i]))
	}

	return body.String()
}

// processTable processes a table element
func processTable(table *jplaw.Table) string {
	var body strings.Builder

	// Determine table class based on writing mode
	tableClass := "law-table"
	if table.WritingMode == jplaw.WritingModeVertical {
		tableClass += " vertical-writing"
	}

	body.WriteString(fmt.Sprintf(`<div class="table-container"><table class="%s">`, tableClass))

	// Process header rows
	if len(table.TableHeaderRow) > 0 {
		body.WriteString("<thead>")
		for i := range table.TableHeaderRow {
			body.WriteString(processTableHeaderRow(&table.TableHeaderRow[i]))
		}
		body.WriteString("</thead>")
	}

	// Process regular rows
	body.WriteString("<tbody>")
	for i := range table.TableRow {
		body.WriteString(processTableRow(&table.TableRow[i]))
	}
	body.WriteString("</tbody>")

	body.WriteString("</table></div>")
	return body.String()
}

// processTableHeaderRow processes a table header row
func processTableHeaderRow(row *jplaw.TableHeaderRow) string {
	var body strings.Builder
	body.WriteString("<tr>")

	for i := range row.TableHeaderColumn {
		body.WriteString(processTableHeaderColumn(&row.TableHeaderColumn[i]))
	}

	body.WriteString("</tr>")
	return body.String()
}

// processTableRow processes a table row
func processTableRow(row *jplaw.TableRow) string {
	var body strings.Builder
	body.WriteString("<tr>")

	for i := range row.TableColumn {
		body.WriteString(processTableColumn(&row.TableColumn[i]))
	}

	body.WriteString("</tr>")
	return body.String()
}

// processTableHeaderColumn processes a table header column
func processTableHeaderColumn(col *jplaw.TableHeaderColumn) string {
	// TableHeaderColumn has different structure - it has Content directly
	content := processTextWithRuby(col.Content, col.Ruby)
	return fmt.Sprintf("<th>%s</th>", content)
}

// processTableColumn processes a table column
func processTableColumn(col *jplaw.TableColumn) string {
	// Build style from border attributes (all strings)
	style := buildCellStyle(col.BorderTop, col.BorderBottom, col.BorderLeft, col.BorderRight)
	// Convert string span values to int for the helper function
	rowspan := parseSpan(col.Rowspan)
	colspan := parseSpan(col.Colspan)
	attrs := buildCellAttributes(rowspan, colspan, col.Align, col.Valign, style)

	var content strings.Builder
	
	// Process sentences
	for i := range col.Sentence {
		content.WriteString(col.Sentence[i].HTML())
	}
	
	// Process column elements (nested content)
	for i := range col.Column {
		content.WriteString(processColumnElement(&col.Column[i]))
	}
	
	// Process parts
	for i := range col.Part {
		content.WriteString(processPartElement(&col.Part[i]))
	}

	return fmt.Sprintf("<td%s>%s</td>", attrs, content.String())
}

// processColumnElement processes a column element within a table cell
func processColumnElement(col *jplaw.Column) string {
	var content strings.Builder
	
	for i := range col.Sentence {
		content.WriteString(col.Sentence[i].HTML())
	}
	
	if col.LineBreak {
		content.WriteString("<br/>")
	}
	
	return content.String()
}

// processPartElement processes a part element within a table cell
func processPartElement(part *jplaw.Part) string {
	var content strings.Builder
	
	// Process part title if present
	if part.PartTitle.Content != "" {
		titleHTML := processTextWithRuby(part.PartTitle.Content, part.PartTitle.Ruby)
		content.WriteString(fmt.Sprintf(`<div class="part-title">%s</div>`, titleHTML))
	}
	
	// Process articles
	for i := range part.Article {
		// For simplicity, just show article titles in tables
		if part.Article[i].ArticleTitle != nil {
			content.WriteString(fmt.Sprintf(`<div class="article-ref">%s</div>`, 
				html.EscapeString(part.Article[i].ArticleTitle.Content)))
		}
	}
	
	return content.String()
}

// parseSpan converts a span string to an int pointer
func parseSpan(span string) *int {
	if span == "" {
		return nil
	}
	if val, err := strconv.Atoi(span); err == nil && val > 0 {
		return &val
	}
	return nil
}

// buildCellStyle builds the style attribute for a table cell based on borders
func buildCellStyle(top, bottom, left, right string) string {
	var styles []string

	// Apply default border style if not "none"
	if top != "" && top != "none" {
		// Use the specified style, or default to "solid #ccc"
		if top == "solid" || top == "dashed" || top == "dotted" || top == "double" {
			styles = append(styles, fmt.Sprintf("border-top: 1px %s #ccc", top))
		} else {
			styles = append(styles, "border-top: 1px solid #ccc")
		}
	}
	if bottom != "" && bottom != "none" {
		// Use the specified style, or default to "solid #ccc"
		if bottom == "solid" || bottom == "dashed" || bottom == "dotted" || bottom == "double" {
			styles = append(styles, fmt.Sprintf("border-bottom: 1px %s #ccc", bottom))
		} else {
			styles = append(styles, "border-bottom: 1px solid #ccc")
		}
	}
	if left != "" && left != "none" {
		// Use the specified style, or default to "solid #ccc"
		if left == "solid" || left == "dashed" || left == "dotted" || left == "double" {
			styles = append(styles, fmt.Sprintf("border-left: 1px %s #ccc", left))
		} else {
			styles = append(styles, "border-left: 1px solid #ccc")
		}
	}
	if right != "" && right != "none" {
		// Use the specified style, or default to "solid #ccc"
		if right == "solid" || right == "dashed" || right == "dotted" || right == "double" {
			styles = append(styles, fmt.Sprintf("border-right: 1px %s #ccc", right))
		} else {
			styles = append(styles, "border-right: 1px solid #ccc")
		}
	}

	if len(styles) == 0 {
		return ""
	}
	return strings.Join(styles, "; ")
}

// buildCellAttributes builds the attributes for a table cell
func buildCellAttributes(rowspan, colspan *int, align, valign, style string) string {
	var attrs []string

	if rowspan != nil && *rowspan > 1 {
		attrs = append(attrs, fmt.Sprintf(`rowspan="%d"`, *rowspan))
	}
	if colspan != nil && *colspan > 1 {
		attrs = append(attrs, fmt.Sprintf(`colspan="%d"`, *colspan))
	}
	if align != "" {
		attrs = append(attrs, fmt.Sprintf(`align="%s"`, string(align)))
	}
	if valign != "" {
		attrs = append(attrs, fmt.Sprintf(`valign="%s"`, string(valign)))
	}
	if style != "" {
		attrs = append(attrs, fmt.Sprintf(`style="%s"`, html.EscapeString(style)))
	}

	if len(attrs) == 0 {
		return ""
	}
	return " " + strings.Join(attrs, " ")
}