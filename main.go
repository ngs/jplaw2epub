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
func processTextWithRuby(content string, rubies []jplaw.Ruby) string {
	if len(rubies) == 0 {
		return html.EscapeString(content)
	}
	
	var result strings.Builder
	if content != "" {
		result.WriteString(html.EscapeString(content))
	}
	result.WriteString(processRubyElements(rubies))
	return result.String()
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

	book.SetAuthor(data.LawNum)
	book.SetLang(string(data.Lang))
	
	// Add CSS for Ruby text rendering
	rubyCSSContent := `
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
`
	rubyCSS, err := book.AddCSS(rubyCSSContent, "ruby.css")
	if err != nil {
		fmt.Printf("Error adding Ruby CSS: %v\n", err)
		os.Exit(1)
	}
	_ = rubyCSS // Mark as used
	eraStr := ""
	switch data.Era {
	case jplaw.EraMeiji:
		eraStr = "明治"
	case jplaw.EraTaisho:
		eraStr = "大正"
	case jplaw.EraShowa:
		eraStr = "昭和"
	case jplaw.EraHeisei:
		eraStr = "平成"
	case jplaw.EraReiwa:
		eraStr = "令和"
	}
	description := fmt.Sprintf("公布日: %s %d年%d月%d日", eraStr, data.Year, data.PromulgateMonth, data.PromulgateDay)
	description += fmt.Sprintf("\n法令番号: %s", data.LawNum)
	// Include Ruby text in description for better accessibility
	lawTitleWithRuby := processTextWithRuby(data.LawBody.LawTitle.Content, data.LawBody.LawTitle.Ruby)
	description += fmt.Sprintf("\n現行法令名: %s %s", lawTitleWithRuby, data.LawBody.LawTitle.Kana)
	book.SetDescription(description)

	for i, chapter := range data.LawBody.MainProvision.Chapter {
		filename := fmt.Sprintf("chapter-%d.xhtml", i)
		chapterTitleHTML := processTextWithRuby(chapter.ChapterTitle.Content, chapter.ChapterTitle.Ruby)
		body := fmt.Sprintf("<h2>%s</h2>", chapterTitleHTML)
		filename, err = book.AddSection(body, chapter.ChapterTitle.Content, filename, "")
		if err != nil {
			fmt.Printf("Error adding section: %v\n", err)
			os.Exit(1)
		}
		for j, article := range chapter.Article {
			subFilename := fmt.Sprintf("article-%d-%d.xhtml", i, j)
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
			body := fmt.Sprintf("<h3>%s</h3><ol>", articleTitleFull)
			for _, para := range article.Paragraph {
				for _, sentense := range para.ParagraphSentence.Sentence {
					sentenceHTML := processTextWithRuby(sentense.Content, sentense.Ruby)
					body += fmt.Sprintf("<li>%s</li>", sentenceHTML)
				}
			}
			body += "</ol>"
			// Use plain text for table of contents
			articleTitlePlain := article.ArticleTitle.Content
			if article.ArticleCaption != nil {
				articleTitlePlain = fmt.Sprintf("%s %s", article.ArticleTitle.Content, article.ArticleCaption.Content)
			}
			_, err = book.AddSubSection(filename, body, articleTitlePlain, subFilename, "")
			if err != nil {
				fmt.Printf("Error adding section: %v\n", err)
				os.Exit(1)
			}
		}
	}

	err = book.Write(*destPath)
	if err != nil {
		fmt.Printf("Error writing epub: %v\n", err)
		os.Exit(1)
	}

}
