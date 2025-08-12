package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"html"
	"io"
	"os"
	"strings"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

// processRubyElements converts Ruby elements to HTML ruby tags
func processRubyElements(rubies []jplaw.Ruby) string {
	var result strings.Builder
	for _, ruby := range rubies {
		if len(ruby.Rt) > 0 {
			result.WriteString("<ruby>")
			result.WriteString(html.EscapeString(ruby.Content))
			for _, rt := range ruby.Rt {
				result.WriteString("<rt>")
				result.WriteString(html.EscapeString(rt.Content))
				result.WriteString("</rt>")
			}
			result.WriteString("</ruby>")
		} else {
			result.WriteString(html.EscapeString(ruby.Content))
		}
	}
	return result.String()
}

// processTextWithRuby processes mixed content (text + Ruby elements)  
// Note: Due to XML parsing limitations, Ruby elements that were inline in the original
// XML are extracted separately, losing their position. As a workaround, we append them.
func processTextWithRuby(content string, rubies []jplaw.Ruby) string {
	if len(rubies) == 0 {
		return html.EscapeString(content)
	}

	// For now, we just append Ruby elements at the end
	// This is not ideal but the jplaw-xml library doesn't preserve position
	var result strings.Builder
	if content != "" {
		result.WriteString(html.EscapeString(content))
	}
	
	// Add Ruby elements (they will appear at the end of the text)
	// In the case of "較(こう)正", this will show the Ruby annotation
	result.WriteString(processRubyElements(rubies))
	return result.String()
}

// getEraString converts Era enum to Japanese string
func getEraString(era jplaw.Era) string {
	switch era {
	case jplaw.EraMeiji:
		return "明治"
	case jplaw.EraTaisho:
		return "大正"
	case jplaw.EraShowa:
		return "昭和"
	case jplaw.EraHeisei:
		return "平成"
	case jplaw.EraReiwa:
		return "令和"
	default:
		return ""
	}
}

// getRubyCSS returns CSS for Ruby text rendering as a style tag
func getRubyCSS() string {
	return `<style type="text/css">
/* Ruby text styling for Japanese phonetic guides */
ruby {
	ruby-align: center;
	ruby-position: over;
}

rt {
	font-size: 0.6em;
	line-height: 1;
	text-align: center;
	color: #666;
}

/* Fallback for older EPUB readers */
ruby > rt {
	display: inline-block;
	font-size: 0.6em;
	line-height: 1;
	text-align: center;
	vertical-align: top;
	color: #666;
}

/* Modern browsers and readers */
@supports (ruby-position: over) {
	ruby {
		ruby-position: over;
	}
}
</style>`
}

// processArticles processes a slice of articles and adds them to the EPUB
func processArticles(book *epub.Epub, articles []jplaw.Article, parentFilename string, chapterIdx, sectionIdx int) error {
	for j := range articles {
		article := &articles[j] // Use pointer to avoid copying
		var subFilename string
		if sectionIdx >= 0 {
			subFilename = fmt.Sprintf("article-%d-%d-%d.xhtml", chapterIdx, sectionIdx, j)
		} else {
			subFilename = fmt.Sprintf("article-%d-%d.xhtml", chapterIdx, j)
		}

		articleTitleHTML := processTextWithRuby(article.ArticleTitle.Content, article.ArticleTitle.Ruby)
		articleCaptionHTML := ""
		if article.ArticleCaption != nil {
			articleCaptionHTML = processTextWithRuby(article.ArticleCaption.Content, article.ArticleCaption.Ruby)
		}
		var articleTitleFull string
		if articleCaptionHTML != "" {
			articleTitleFull = fmt.Sprintf("%s %s", articleTitleHTML, articleCaptionHTML)
		} else {
			articleTitleFull = articleTitleHTML
		}
		body := fmt.Sprintf("%s<h3>%s</h3>", getRubyCSS(), articleTitleFull)
		
		// Process paragraphs
		for paraIdx := range article.Paragraph {
			para := &article.Paragraph[paraIdx]
			
			// Add paragraph number if present
			if para.ParagraphNum.Content != "" {
				body += fmt.Sprintf("<h4>%s</h4>", html.EscapeString(para.ParagraphNum.Content))
			}
			
			// Process paragraph sentences
			if len(para.ParagraphSentence.Sentence) > 0 {
				body += "<p>"
				for sentenceIdx := range para.ParagraphSentence.Sentence {
					sentence := &para.ParagraphSentence.Sentence[sentenceIdx]
					// Use the new HTML() method which properly handles inline Ruby
					body += sentence.HTML()
				}
				body += "</p>"
			}
			
			// Process items within paragraph
			if len(para.Item) > 0 {
				body += "<ol>"
				for itemIdx := range para.Item {
					item := &para.Item[itemIdx]
					body += "<li>"
					
					// Add item title if present
					if item.ItemTitle != nil && item.ItemTitle.Content != "" {
						body += fmt.Sprintf("<strong>%s</strong> ", html.EscapeString(item.ItemTitle.Content))
					}
					
					// Process item sentences
					for sentIdx := range item.ItemSentence.Sentence {
						sent := &item.ItemSentence.Sentence[sentIdx]
						// Use the new HTML() method which properly handles inline Ruby
						body += sent.HTML()
					}
					
					// Process subitems if any
					if len(item.Subitem1) > 0 {
						body += "<ol type='i'>"
						for subIdx := range item.Subitem1 {
							subitem := &item.Subitem1[subIdx]
							body += "<li>"
							if subitem.Subitem1Title != nil {
								body += fmt.Sprintf("<strong>%s</strong> ", html.EscapeString(subitem.Subitem1Title.Content))
							}
							for subSentIdx := range subitem.Subitem1Sentence.Sentence {
								subSent := &subitem.Subitem1Sentence.Sentence[subSentIdx]
								// Use the new HTML() method which properly handles inline Ruby
								body += subSent.HTML()
							}
							body += "</li>"
						}
						body += "</ol>"
					}
					
					body += "</li>"
				}
				body += "</ol>"
			}
		}
		// Use plain text for table of contents
		articleTitlePlain := article.ArticleTitle.Content
		if article.ArticleCaption != nil {
			articleTitlePlain = fmt.Sprintf("%s %s", article.ArticleTitle.Content, article.ArticleCaption.Content)
		}
		_, subSectionErr := book.AddSubSection(parentFilename, body, articleTitlePlain, subFilename, "")
		if subSectionErr != nil {
			return fmt.Errorf("error adding article section: %w", subSectionErr)
		}
	}
	return nil
}

// setupEPUBMetadata sets up the basic EPUB metadata
func setupEPUBMetadata(book *epub.Epub, data *jplaw.Law) error {
	book.SetAuthor(data.LawNum)
	book.SetLang(string(data.Lang))

	// Set description
	eraStr := getEraString(data.Era)
	description := fmt.Sprintf("公布日: %s %d年%d月%d日", eraStr, data.Year, data.PromulgateMonth, data.PromulgateDay)
	description += fmt.Sprintf("\n法令番号: %s", data.LawNum)
	lawTitleWithRuby := processTextWithRuby(data.LawBody.LawTitle.Content, data.LawBody.LawTitle.Ruby)
	description += fmt.Sprintf("\n現行法令名: %s %s", lawTitleWithRuby, data.LawBody.LawTitle.Kana)
	book.SetDescription(description)

	return nil
}

func main() {
	destPath := flag.String("d", "", "Destination file path")
	flag.Parse()

	if *destPath == "" {
		fmt.Println("Destination file path is required")
		os.Exit(1)
	}

	if len(flag.Args()) < 1 {
		fmt.Println("Source file path is required as the first argument")
		os.Exit(1)
	}

	sourcePath := flag.Arg(0)

	xmlFile, err := os.Open(sourcePath)
	if err != nil {
		fmt.Printf("Error opening source file: %v\n", err)
		os.Exit(1)
	}
	defer xmlFile.Close()

	byteValue, err := io.ReadAll(xmlFile)
	if err != nil {
		fmt.Printf("Error reading source file: %v\n", err)
		os.Exit(1)
	}

	var data jplaw.Law
	if err := xml.Unmarshal(byteValue, &data); err != nil {
		fmt.Printf("Error unmarshalling XML: %v\n", err)
		os.Exit(1)
	}

	// Use plain text for EPUB title (no HTML markup allowed)
	book, err := epub.NewEpub(data.LawBody.LawTitle.Content)
	if err != nil {
		fmt.Printf("Error creating epub: %v\n", err)
		os.Exit(1)
	}

	// Setup EPUB metadata and CSS
	if setupErr := setupEPUBMetadata(book, &data); setupErr != nil {
		fmt.Printf("Error setting up EPUB metadata: %v\n", setupErr)
		os.Exit(1)
	}

	for i, chapter := range data.LawBody.MainProvision.Chapter {
		chapterFilename := fmt.Sprintf("chapter-%d.xhtml", i)
		chapterTitleHTML := processTextWithRuby(chapter.ChapterTitle.Content, chapter.ChapterTitle.Ruby)
		body := fmt.Sprintf("%s<h2>%s</h2>", getRubyCSS(), chapterTitleHTML)

		// Process Sections if any
		if len(chapter.Section) > 0 {
			body += "<div class='sections'>"
			for sIdx, section := range chapter.Section {
				sectionTitleHTML := processTextWithRuby(section.SectionTitle.Content, section.SectionTitle.Ruby)
				body += fmt.Sprintf("<h3>%s</h3>", sectionTitleHTML)
				// Add a note about articles in this section
				if len(section.Article) > 0 {
					body += fmt.Sprintf("<p>（%s から %s まで）</p>",
						section.Article[0].ArticleTitle.Content,
						section.Article[len(section.Article)-1].ArticleTitle.Content)
				}
				_ = sIdx // Mark as used
			}
			body += "</div>"
		}

		chapterFilename, addErr := book.AddSection(body, chapter.ChapterTitle.Content, chapterFilename, "")
		if addErr != nil {
			fmt.Printf("Error adding chapter: %v\n", addErr)
			os.Exit(1)
		}

		// Process direct articles under chapter
		if len(chapter.Article) > 0 {
			if err := processArticles(book, chapter.Article, chapterFilename, i, -1); err != nil {
				fmt.Printf("Error processing chapter articles: %v\n", err)
				os.Exit(1)
			}
		}

		// Process articles within sections
		for sIdx := range chapter.Section {
			section := &chapter.Section[sIdx]
			if len(section.Article) > 0 {
				if err := processArticles(book, section.Article, chapterFilename, i, sIdx); err != nil {
					fmt.Printf("Error processing section articles: %v\n", err)
					os.Exit(1)
				}
			}
		}
	}

	writeErr := book.Write(*destPath)
	if writeErr != nil {
		fmt.Printf("Error writing epub: %v\n", writeErr)
		os.Exit(1)
	}
}
