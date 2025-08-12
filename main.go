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

// Removed baseCSS - list styles are now determined dynamically

// HTML constants
const (
	htmlOL    = "<ol>"
	htmlOLEnd = "</ol>"
	htmlLI    = "<li>"
	htmlLIEnd = "</li>"

	// List style types
	listStyleDisc     = "disc"
	listStyleDecimal  = "decimal"
	listStyleCJK      = "cjk-ideographic"
	listStyleKatakana = "katakana-iroha"
	listStyleHiragana = "hiragana-iroha"
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

// isListNumber checks if the text is just a Japanese list number
func isListNumber(text string) bool {
	// List numbers that should be skipped
	listNumbers := []string{
		// CJK ideographic numbers
		"一", "二", "三", "四", "五", "六", "七", "八", "九", "十",
		"十一", "十二", "十三", "十四", "十五", "十六", "十七", "十八", "十九", "二十",
		// Katakana iroha
		"イ", "ロ", "ハ", "ニ", "ホ", "ヘ", "ト", "チ", "リ", "ヌ",
		"ル", "ヲ", "ワ", "カ", "ヨ", "タ", "レ", "ソ", "ツ", "ネ",
		// Full-width Arabic numerals
		"１", "２", "３", "４", "５", "６", "７", "８", "９", "１０",
		"１１", "１２", "１３", "１４", "１５", "１６", "１７", "１８", "１９", "２０",
	}

	for _, num := range listNumbers {
		if text == num {
			return true
		}
	}
	return false
}

// getListStyleType determines the CSS list-style-type based on the item titles
func getListStyleType(titles []string) string {
	if len(titles) == 0 {
		return ""
	}

	// Check first title to determine the pattern
	first := titles[0]

	// CJK ideographic (一, 二, 三...)
	cjkNumbers := []string{"一", "二", "三", "四", "五", "六", "七", "八", "九", "十"}
	for _, num := range cjkNumbers {
		if first == num {
			return listStyleCJK
		}
	}

	// Katakana iroha (イ, ロ, ハ...)
	katakanaIroha := []string{"イ", "ロ", "ハ", "ニ", "ホ", "ヘ", "ト", "チ", "リ", "ヌ"}
	for _, kana := range katakanaIroha {
		if first == kana {
			return listStyleKatakana
		}
	}

	// Hiragana iroha (い, ろ, は...)
	hiraganaIroha := []string{"い", "ろ", "は", "に", "ほ", "へ", "と", "ち", "り", "ぬ"}
	for _, kana := range hiraganaIroha {
		if first == kana {
			return listStyleHiragana
		}
	}

	// Full-width Arabic numerals (１, ２, ３...)
	fullWidthNumbers := []string{"１", "２", "３", "４", "５", "６", "７", "８", "９"}
	for _, num := range fullWidthNumbers {
		if first == num {
			return listStyleDecimal
		}
	}

	// Half-width Arabic numerals (1, 2, 3...)
	if strings.HasPrefix(first, "1") {
		return listStyleDecimal
	}

	// Parenthesized numbers (（１）, （２）...)
	if strings.HasPrefix(first, "（") && strings.HasSuffix(first, "）") {
		return listStyleDecimal
	}

	// Default
	return listStyleDisc
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
		body := fmt.Sprintf("<h3>%s</h3>", articleTitleFull)

		// Group paragraphs with Num attribute into lists
		var inList bool

		for paraIdx := range article.Paragraph {
			para := &article.Paragraph[paraIdx]

			// Check if this paragraph should be in a list (has Num > 0)
			if para.Num > 0 {
				// Start a new list if not in one
				if !inList {
					// Determine list style from paragraph numbers if possible
					var paraNumTitles []string
					for i := paraIdx; i < len(article.Paragraph) && article.Paragraph[i].Num > 0; i++ {
						paraNumTitles = append(paraNumTitles, article.Paragraph[i].ParagraphNum.Content)
					}
					listStyle := getListStyleType(paraNumTitles)
					if listStyle != "" && listStyle != listStyleDisc {
						body += fmt.Sprintf(`<ol style="list-style-type: %s;">`, listStyle)
					} else {
						body += htmlOL
					}
					inList = true
				}

				body += htmlLI

				// Add paragraph number if present (as a heading within the list item)
				if para.ParagraphNum.Content != "" {
					// Skip if it's just a number that matches the list style
					if !isListNumber(para.ParagraphNum.Content) {
						body += fmt.Sprintf("<strong>%s</strong> ", html.EscapeString(para.ParagraphNum.Content))
					}
				}

				// Process paragraph sentences
				if len(para.ParagraphSentence.Sentence) > 0 {
					for sentenceIdx := range para.ParagraphSentence.Sentence {
						sentence := &para.ParagraphSentence.Sentence[sentenceIdx]
						body += sentence.HTML()
					}
				}

				// Process items within paragraph
				if len(para.Item) > 0 {
					// Collect item titles to determine list style
					var itemTitles []string
					for _, item := range para.Item {
						if item.ItemTitle != nil {
							itemTitles = append(itemTitles, item.ItemTitle.Content)
						}
					}
					itemListStyle := getListStyleType(itemTitles)

					if itemListStyle != "" && itemListStyle != listStyleDisc {
						body += fmt.Sprintf(`<ol style="list-style-type: %s;">`, itemListStyle)
					} else {
						body += htmlOL
					}

					for itemIdx := range para.Item {
						item := &para.Item[itemIdx]
						body += htmlLI

						// Add item title if present (skip if it's just a list number)
						if item.ItemTitle != nil && item.ItemTitle.Content != "" {
							if !isListNumber(item.ItemTitle.Content) {
								body += fmt.Sprintf("<strong>%s</strong> ", html.EscapeString(item.ItemTitle.Content))
							}
						}

						// Process item sentences
						for sentIdx := range item.ItemSentence.Sentence {
							sent := &item.ItemSentence.Sentence[sentIdx]
							body += sent.HTML()
						}

						// Process Subitem1 if any
						if len(item.Subitem1) > 0 {
							// Collect subitem titles to determine list style
							var subitemTitles []string
							for _, subitem := range item.Subitem1 {
								if subitem.Subitem1Title != nil {
									subitemTitles = append(subitemTitles, subitem.Subitem1Title.Content)
								}
							}
							subitemListStyle := getListStyleType(subitemTitles)

							if subitemListStyle != "" && subitemListStyle != listStyleDisc {
								body += fmt.Sprintf(`<ol style="list-style-type: %s;">`, subitemListStyle)
							} else {
								body += htmlOL
							}

							for subIdx := range item.Subitem1 {
								subitem := &item.Subitem1[subIdx]
								body += htmlLI
								if subitem.Subitem1Title != nil && subitem.Subitem1Title.Content != "" {
									if !isListNumber(subitem.Subitem1Title.Content) {
										body += fmt.Sprintf("<strong>%s</strong> ", html.EscapeString(subitem.Subitem1Title.Content))
									}
								}
								for subSentIdx := range subitem.Subitem1Sentence.Sentence {
									subSent := &subitem.Subitem1Sentence.Sentence[subSentIdx]
									body += subSent.HTML()
								}

								// Process Subitem2 if any
								if len(subitem.Subitem2) > 0 {
									// Collect subitem2 titles to determine list style
									var subitem2Titles []string
									for _, subitem2 := range subitem.Subitem2 {
										if subitem2.Subitem2Title != nil {
											subitem2Titles = append(subitem2Titles, subitem2.Subitem2Title.Content)
										}
									}
									subitem2ListStyle := getListStyleType(subitem2Titles)

									if subitem2ListStyle != "" && subitem2ListStyle != listStyleDisc {
										body += fmt.Sprintf(`<ol style="list-style-type: %s;">`, subitem2ListStyle)
									} else {
										body += htmlOL
									}

									for sub2Idx := range subitem.Subitem2 {
										subitem2 := &subitem.Subitem2[sub2Idx]
										body += htmlLI
										if subitem2.Subitem2Title != nil && subitem2.Subitem2Title.Content != "" {
											if !isListNumber(subitem2.Subitem2Title.Content) {
												body += fmt.Sprintf("<strong>%s</strong> ", html.EscapeString(subitem2.Subitem2Title.Content))
											}
										}
										for sub2SentIdx := range subitem2.Subitem2Sentence.Sentence {
											sub2Sent := &subitem2.Subitem2Sentence.Sentence[sub2SentIdx]
											body += sub2Sent.HTML()
										}
										body += htmlLIEnd
									}
									body += htmlOLEnd
								}

								body += htmlLIEnd
							}
							body += htmlOLEnd
						}

						body += htmlLIEnd
					}
					body += htmlOLEnd
				}

				body += htmlLIEnd

			} else {
				// Close list if we were in one
				if inList {
					body += htmlOLEnd
					inList = false
				}

				// Process as regular paragraph (not in a list)
				if para.ParagraphNum.Content != "" {
					body += fmt.Sprintf("<h4>%s</h4>", html.EscapeString(para.ParagraphNum.Content))
				}

				if len(para.ParagraphSentence.Sentence) > 0 {
					body += "<p>"
					for sentenceIdx := range para.ParagraphSentence.Sentence {
						sentence := &para.ParagraphSentence.Sentence[sentenceIdx]
						body += sentence.HTML()
					}
					body += "</p>"
				}

				// Process items (if not in numbered paragraph)
				if len(para.Item) > 0 {
					// Similar item processing as above
					var itemTitles []string
					for _, item := range para.Item {
						if item.ItemTitle != nil {
							itemTitles = append(itemTitles, item.ItemTitle.Content)
						}
					}
					itemListStyle := getListStyleType(itemTitles)

					if itemListStyle != "" && itemListStyle != listStyleDisc {
						body += fmt.Sprintf(`<ol style="list-style-type: %s;">`, itemListStyle)
					} else {
						body += htmlOL
					}

					for itemIdx := range para.Item {
						item := &para.Item[itemIdx]
						body += htmlLI

						if item.ItemTitle != nil && item.ItemTitle.Content != "" {
							if !isListNumber(item.ItemTitle.Content) {
								body += fmt.Sprintf("<strong>%s</strong> ", html.EscapeString(item.ItemTitle.Content))
							}
						}

						for sentIdx := range item.ItemSentence.Sentence {
							sent := &item.ItemSentence.Sentence[sentIdx]
							body += sent.HTML()
						}

						body += htmlLIEnd
					}
					body += htmlOLEnd
				}
			}
		}

		// Close list if still open at the end
		if inList {
			body += htmlOLEnd
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
		body := fmt.Sprintf("<h2>%s</h2>", chapterTitleHTML)

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
