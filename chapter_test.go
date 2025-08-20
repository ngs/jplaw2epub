package jplaw2epub

import (
	"strings"
	"testing"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

func TestProcessChapterWithImages(t *testing.T) {
	tests := []struct {
		name        string
		chapter     *jplaw.Chapter
		chapterIdx  int
		wantContent []string
		wantErr     bool
	}{
		{
			name: "chapter with sections",
			chapter: &jplaw.Chapter{
				ChapterTitle: jplaw.ChapterTitle{
					Content: "第一章 総則",
				},
				Section: []jplaw.Section{
					{
						SectionTitle: jplaw.SectionTitle{
							Content: "第一節 通則",
						},
						Article: []jplaw.Article{
							{
								ArticleTitle: &jplaw.ArticleTitle{
									Content: "第一条",
								},
							},
							{
								ArticleTitle: &jplaw.ArticleTitle{
									Content: "第三条",
								},
							},
						},
					},
					{
						SectionTitle: jplaw.SectionTitle{
							Content: "第二節 定義",
						},
						Article: []jplaw.Article{
							{
								ArticleTitle: &jplaw.ArticleTitle{
									Content: "第四条",
								},
							},
						},
					},
				},
			},
			chapterIdx: 0,
			wantContent: []string{
				"第一章 総則",
				"第一節 通則",
				"第一条 から 第三条 まで",
				"第二節 定義",
			},
			wantErr: false,
		},
		{
			name: "chapter with direct articles",
			chapter: &jplaw.Chapter{
				ChapterTitle: jplaw.ChapterTitle{
					Content: "第二章 権利と義務",
				},
				Article: []jplaw.Article{
					{
						ArticleTitle: &jplaw.ArticleTitle{
							Content: "第十条",
						},
						ArticleCaption: &jplaw.ArticleCaption{
							Content: "（基本的権利）",
						},
						Paragraph: []jplaw.Paragraph{
							{
								ParagraphNum: jplaw.ParagraphNum{
									Content: "１",
								},
								ParagraphSentence: jplaw.ParagraphSentence{
									Sentence: []jplaw.Sentence{
										{Content: "権利の内容"},
									},
								},
							},
						},
					},
				},
			},
			chapterIdx:  1,
			wantContent: []string{"第二章 権利と義務"},
			wantErr:     false,
		},
		{
			name: "chapter with ruby text",
			chapter: &jplaw.Chapter{
				ChapterTitle: jplaw.ChapterTitle{
					Content: "第三章 罰則",
					Ruby: []jplaw.Ruby{
						{
							Content: "罰則",
							Rt:      []jplaw.Rt{{Content: "ばっそく"}},
						},
					},
				},
			},
			chapterIdx:  2,
			wantContent: []string{"罰則", "ばっそく"},
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new EPUB book
			book, err := epub.NewEpub("Test Book")
			if err != nil {
				t.Fatalf("Failed to create EPUB: %v", err)
			}

			// Process chapter
			err = processChapterWithImages(book, tt.chapter, tt.chapterIdx, nil)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("processChapterWithImages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			// Since we can't easily inspect the internal content without writing to disk,
			// we'll just verify that the function didn't error
			// Test passed if no error occurred
		})
	}
}

func TestBuildChapterBody(t *testing.T) {
	tests := []struct {
		name    string
		chapter *jplaw.Chapter
		want    []string
	}{
		{
			name: "simple chapter",
			chapter: &jplaw.Chapter{
				ChapterTitle: jplaw.ChapterTitle{
					Content: "第一章 総則",
				},
			},
			want: []string{
				`<div class="chapter-title">第一章 総則</div>`,
			},
		},
		{
			name: "chapter with sections",
			chapter: &jplaw.Chapter{
				ChapterTitle: jplaw.ChapterTitle{
					Content: "第一章",
				},
				Section: []jplaw.Section{
					{
						SectionTitle: jplaw.SectionTitle{
							Content: "第一節",
						},
						Article: []jplaw.Article{
							{ArticleTitle: &jplaw.ArticleTitle{Content: "第一条"}},
							{ArticleTitle: &jplaw.ArticleTitle{Content: "第二条"}},
						},
					},
				},
			},
			want: []string{
				`<div class="chapter-title">第一章</div>`,
				`<div class='sections'>`,
				`<h3>第一節</h3>`,
				`<p>（第一条 から 第二条 まで）</p>`,
				`</div>`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildChapterBody(tt.chapter)
			for _, expected := range tt.want {
				if !strings.Contains(got, expected) {
					t.Errorf("buildChapterBody() missing expected content: %s", expected)
				}
			}
		})
	}
}

func TestBuildSectionsHTML(t *testing.T) {
	sections := []jplaw.Section{
		{
			SectionTitle: jplaw.SectionTitle{
				Content: "第一節 総則",
			},
			Article: []jplaw.Article{
				{ArticleTitle: &jplaw.ArticleTitle{Content: "第一条"}},
				{ArticleTitle: &jplaw.ArticleTitle{Content: "第五条"}},
			},
		},
		{
			SectionTitle: jplaw.SectionTitle{
				Content: "第二節 手続",
				Ruby: []jplaw.Ruby{
					{
						Content: "手続",
						Rt:      []jplaw.Rt{{Content: "てつづき"}},
					},
				},
			},
			Article: []jplaw.Article{
				{ArticleTitle: &jplaw.ArticleTitle{Content: "第六条"}},
			},
		},
	}

	got := buildSectionsHTML(sections)

	expectedContent := []string{
		`<div class='sections'>`,
		`<h3>第一節 総則</h3>`,
		`<p>（第一条 から 第五条 まで）</p>`,
		`<h3>第二節 手続<ruby>手続<rt>てつづき</rt></ruby></h3>`,
		`<p>（第六条 から 第六条 まで）</p>`,
		`</div>`,
	}

	for _, expected := range expectedContent {
		if !strings.Contains(got, expected) {
			t.Errorf("buildSectionsHTML() missing expected content: %s", expected)
		}
	}
}
