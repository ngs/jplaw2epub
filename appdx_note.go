package jplaw2epub

import (
	"encoding/xml"
	"fmt"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

const defaultAppdxNoteTitle = "附則"

// processAppdxNotes processes appendix notes
func processAppdxNotes(book *epub.Epub, notes []jplaw.AppdxNote, imgProc *ImageProcessor) error {
	if len(notes) == 0 {
		return nil
	}

	for idx, note := range notes {
		if err := processAppdxNote(book, &note, idx, imgProc); err != nil {
			return fmt.Errorf("processing AppdxNote %d: %w", idx, err)
		}
	}

	return nil
}

// processAppdxNote processes a single appendix note
func processAppdxNote(book *epub.Epub, note *jplaw.AppdxNote, idx int, imgProc *ImageProcessor) error {
	filename := fmt.Sprintf("appdx-note-%d.xhtml", idx)
	body := ""

	// Add title if present
	title := defaultAppdxNoteTitle
	if note.AppdxNoteTitle != nil && note.AppdxNoteTitle.Content != "" {
		title = note.AppdxNoteTitle.Content
		body += fmt.Sprintf(`<div class="chapter-title">%s</div>`, processTextWithRuby(title, note.AppdxNoteTitle.Ruby))
	}

	// Process related article number if present
	if note.RelatedArticleNum != nil && note.RelatedArticleNum.Content != "" {
		body += fmt.Sprintf(`<div class="related-articles">%s</div>`,
			processTextWithRuby(note.RelatedArticleNum.Content, note.RelatedArticleNum.Ruby))
	}

	// Process NoteStructs
	for _, noteStruct := range note.NoteStruct {
		body += processNoteStruct(&noteStruct, imgProc)
	}

	// Process FigStructs
	for _, figStruct := range note.FigStruct {
		if imgProc != nil {
			html, err := imgProc.ProcessFigStruct(&figStruct)
			if err != nil {
				// Log error but continue
				fmt.Printf("Warning: failed to process FigStruct: %v\n", err)
			} else {
				body += html
			}
		}
	}

	// Process TableStructs
	for _, tableStruct := range note.TableStruct {
		body += processTableStructWithImages(&tableStruct, imgProc)
	}

	// Process Remarks
	if note.Remarks != nil {
		body += processRemarks(note.Remarks)
	}

	// Add the section to the book
	_, err := book.AddSection(body, title, filename, "")
	if err != nil {
		return fmt.Errorf("adding AppdxNote section: %w", err)
	}

	return nil
}

// processNoteStruct processes a NoteStruct
func processNoteStruct(noteStruct *jplaw.NoteStruct, imgProc *ImageProcessor) string {
	body := `<div class="note-struct">`

	// Add title if present
	if noteStruct.NoteStructTitle != nil && noteStruct.NoteStructTitle.Content != "" {
		body += fmt.Sprintf(`<h3>%s</h3>`,
			processTextWithRuby(noteStruct.NoteStructTitle.Content, noteStruct.NoteStructTitle.Ruby))
	}

	// Process Note content
	// Note contains innerxml which we need to parse manually
	noteContent := noteStruct.Note.Content

	// The Note content may contain Paragraph and Item elements as raw XML
	// We need to process them properly
	body += processNoteContent(noteContent, imgProc)

	// Process Remarks
	for i := range noteStruct.Remarks {
		body += processRemarks(&noteStruct.Remarks[i])
	}

	body += htmlDivEnd
	return body
}

// processNoteContent processes the inner content of a Note
func processNoteContent(content string, imgProc *ImageProcessor) string {
	// Parse the XML content as a fragment
	// We'll wrap it in a root element to make it valid XML
	wrappedContent := "<root>" + content + "</root>"

	// Define a structure to parse the Note content
	type NoteContentRoot struct {
		Paragraphs []jplaw.Paragraph `xml:"Paragraph"`
	}

	var root NoteContentRoot
	if err := xml.Unmarshal([]byte(wrappedContent), &root); err != nil {
		// If parsing fails, return the content as-is wrapped in a div
		return fmt.Sprintf("<div class='note-content'>%s</div>", content)
	}

	// Process paragraphs using the existing paragraph processor
	if len(root.Paragraphs) > 0 {
		return processParagraphsWithImages(root.Paragraphs, imgProc)
	}

	// If no paragraphs were found, return the content wrapped in a div
	return fmt.Sprintf("<div class='note-content'>%s</div>", content)
}

// processRemarks processes remarks
func processRemarks(remarks *jplaw.Remarks) string {
	body := `<div class="appdx-remarks">`

	// Add label if present
	if remarks.RemarksLabel.Content != "" {
		body += fmt.Sprintf(`<p class="remarks-label">%s</p>`,
			processTextWithRuby(remarks.RemarksLabel.Content, remarks.RemarksLabel.Ruby))
	}

	// Process sentences
	for i := range remarks.Sentence {
		body += fmt.Sprintf(`<p class="remark">%s</p>`, remarks.Sentence[i].HTML())
	}

	// Process items
	if len(remarks.Item) > 0 {
		body += processItems(remarks.Item)
	}

	body += htmlDivEnd
	return body
}

// processAppdxTables processes appendix tables
func processAppdxTables(book *epub.Epub, tables []jplaw.AppdxTable, imgProc *ImageProcessor) error {
	if len(tables) == 0 {
		return nil
	}

	for idx, table := range tables {
		if err := processAppdxTable(book, &table, idx, imgProc); err != nil {
			return fmt.Errorf("processing AppdxTable %d: %w", idx, err)
		}
	}

	return nil
}

// processAppdxTable processes a single appendix table
func processAppdxTable(book *epub.Epub, table *jplaw.AppdxTable, idx int, imgProc *ImageProcessor) error {
	filename := fmt.Sprintf("appdx-table-%d.xhtml", idx)
	body := ""

	// Add title if present
	title := "附表"
	if table.AppdxTableTitle != nil && table.AppdxTableTitle.Content != "" {
		title = table.AppdxTableTitle.Content
		body += fmt.Sprintf(`<div class="chapter-title">%s</div>`, processTextWithRuby(title, table.AppdxTableTitle.Ruby))
	}

	// Process related article number if present
	if table.RelatedArticleNum != nil && table.RelatedArticleNum.Content != "" {
		body += fmt.Sprintf(`<div class="related-articles">%s</div>`,
			processTextWithRuby(table.RelatedArticleNum.Content, table.RelatedArticleNum.Ruby))
	}

	// Process TableStructs
	for _, tableStruct := range table.TableStruct {
		body += processTableStructWithImages(&tableStruct, imgProc)
	}

	// Process Remarks
	if table.Remarks != nil {
		body += processRemarks(table.Remarks)
	}

	// Add the section to the book
	_, err := book.AddSection(body, title, filename, "")
	if err != nil {
		return fmt.Errorf("adding AppdxTable section: %w", err)
	}

	return nil
}
