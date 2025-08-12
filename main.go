package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"html"
	"io"
	"os"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

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

	book, err := epub.NewEpub(data.LawBody.LawTitle.Content)
	if err != nil {
		fmt.Printf("Error creating epub: %v\n", err)
		os.Exit(1)
	}

	book.SetAuthor(data.LawNum)
	book.SetLang(string(data.Lang))
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
	description += fmt.Sprintf("\n現行法令名: %s %s", data.LawBody.LawTitle.Content, data.LawBody.LawTitle.Kana)
	book.SetDescription(description)

	for i, chapter := range data.LawBody.MainProvision.Chapter {
		filename := fmt.Sprintf("chapter-%d.xhtml", i)
		body := fmt.Sprintf("<h2>%s</h2>", html.EscapeString(chapter.ChapterTitle.Content))
		filename, err = book.AddSection(body, chapter.ChapterTitle.Content, filename, "")
		if err != nil {
			fmt.Printf("Error adding section: %v\n", err)
			os.Exit(1)
		}
		for j, article := range chapter.Article {
			subFilename := fmt.Sprintf("article-%d-%d.xhtml", i, j)
			articleCaption := ""
			if article.ArticleCaption != nil {
				articleCaption = article.ArticleCaption.Content
			}
			articleTitle := fmt.Sprintf("%s %s", article.ArticleTitle.Content, articleCaption)
			body := fmt.Sprintf("<h3>%s</h3><ol>", html.EscapeString(articleTitle))
			for _, para := range article.Paragraph {
				for _, sentense := range para.ParagraphSentence.Sentence {
					body += fmt.Sprintf("<li>%s</li>", html.EscapeString(sentense.Content))
				}
			}
			body += "</ol>"
			_, err = book.AddSubSection(filename, body, articleTitle, subFilename, "")
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
