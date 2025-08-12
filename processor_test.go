package jplaw2epub

import (
	"strings"
	"testing"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

func TestBuildArticleFilename(t *testing.T) {
	tests := []struct {
		name       string
		chapterIdx int
		sectionIdx int
		articleIdx int
		want       string
	}{
		{
			name:       "Article directly under chapter",
			chapterIdx: 1,
			sectionIdx: -1,
			articleIdx: 2,
			want:       "article-1-2.xhtml",
		},
		{
			name:       "Article under section",
			chapterIdx: 3,
			sectionIdx: 4,
			articleIdx: 5,
			want:       "article-3-4-5.xhtml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildArticleFilename(tt.chapterIdx, tt.sectionIdx, tt.articleIdx)
			if got != tt.want {
				t.Errorf("buildArticleFilename(%d, %d, %d) = %v, want %v",
					tt.chapterIdx, tt.sectionIdx, tt.articleIdx, got, tt.want)
			}
		})
	}
}

func TestBuildArticleBody(t *testing.T) {
	article := &jplaw.Article{
		ArticleTitle: &jplaw.ArticleTitle{
			Content: "第一条",
		},
		Paragraph: []jplaw.Paragraph{
			{
				Num: 0,
				ParagraphSentence: jplaw.ParagraphSentence{
					Sentence: []jplaw.Sentence{
						createTestSentence("条文の内容。"),
					},
				},
			},
		},
	}

	articleTitle := "第一条"
	got := buildArticleBody(article, articleTitle)

	expectedParts := []string{
		"<h3>第一条</h3>",
		"条文の内容。",
	}

	for _, part := range expectedParts {
		if !strings.Contains(got, part) {
			t.Errorf("buildArticleBody() should contain %q\ngot: %v", part, got)
		}
	}
}

func TestGetArticleTitlePlain(t *testing.T) {
	tests := []struct {
		name    string
		article *jplaw.Article
		want    string
	}{
		{
			name: "Article with title only",
			article: &jplaw.Article{
				ArticleTitle: &jplaw.ArticleTitle{
					Content: "第一条",
				},
			},
			want: "第一条",
		},
		{
			name: "Article with title and caption",
			article: &jplaw.Article{
				ArticleTitle: &jplaw.ArticleTitle{
					Content: "第二条",
				},
				ArticleCaption: &jplaw.ArticleCaption{
					Content: "（定義）",
				},
			},
			want: "第二条 （定義）",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getArticleTitlePlain(tt.article)
			if got != tt.want {
				t.Errorf("getArticleTitlePlain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessArticle(t *testing.T) {
	// Create a mock EPUB book
	book, err := epub.NewEpub("Test Book")
	if err != nil {
		t.Fatalf("Failed to create test EPUB: %v", err)
	}

	// Add a parent section first
	parentFilename, err := book.AddSection("<h1>Chapter</h1>", "Chapter 1", "chapter-1.xhtml", "")
	if err != nil {
		t.Fatalf("Failed to add parent section: %v", err)
	}

	article := &jplaw.Article{
		ArticleTitle: &jplaw.ArticleTitle{
			Content: "第一条",
		},
		ArticleCaption: &jplaw.ArticleCaption{
			Content: "（目的）",
		},
		Paragraph: []jplaw.Paragraph{
			{
				ParagraphSentence: jplaw.ParagraphSentence{
					Sentence: []jplaw.Sentence{
						{Content: "この法律の目的。"},
					},
				},
			},
		},
	}

	err = processArticle(book, article, parentFilename, 0, -1, 0)
	if err != nil {
		t.Errorf("processArticle() returned error: %v", err)
	}

	// Verify the article was added successfully
	// Note: go-epub doesn't provide a way to inspect added sections,
	// so we're mainly testing that no error occurred
}

func TestProcessArticles(t *testing.T) {
	// Create a mock EPUB book
	book, err := epub.NewEpub("Test Book")
	if err != nil {
		t.Fatalf("Failed to create test EPUB: %v", err)
	}

	// Add a parent section first
	parentFilename, err := book.AddSection("<h1>Chapter</h1>", "Chapter 1", "chapter-1.xhtml", "")
	if err != nil {
		t.Fatalf("Failed to add parent section: %v", err)
	}

	articles := []jplaw.Article{
		{
			ArticleTitle: &jplaw.ArticleTitle{
				Content: "第一条",
			},
			Paragraph: []jplaw.Paragraph{
				{
					ParagraphSentence: jplaw.ParagraphSentence{
						Sentence: []jplaw.Sentence{
							{Content: "第一条の内容。"},
						},
					},
				},
			},
		},
		{
			ArticleTitle: &jplaw.ArticleTitle{
				Content: "第二条",
			},
			Paragraph: []jplaw.Paragraph{
				{
					ParagraphSentence: jplaw.ParagraphSentence{
						Sentence: []jplaw.Sentence{
							{Content: "第二条の内容。"},
						},
					},
				},
			},
		},
	}

	err = processArticles(book, articles, parentFilename, 0, -1)
	if err != nil {
		t.Errorf("processArticles() returned error: %v", err)
	}
}

func TestProcessArticles_Error(t *testing.T) {
	// Create a mock EPUB book
	book, err := epub.NewEpub("Test Book")
	if err != nil {
		t.Fatalf("Failed to create test EPUB: %v", err)
	}

	// Use an invalid parent filename to trigger an error
	invalidParent := "non-existent-parent.xhtml"

	articles := []jplaw.Article{
		{
			ArticleTitle: &jplaw.ArticleTitle{
				Content: "第一条",
			},
		},
	}

	err = processArticles(book, articles, invalidParent, 0, -1)
	if err == nil {
		t.Error("processArticles() should return error for invalid parent")
	}
}
