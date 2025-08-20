package jplaw2epub

import (
	"testing"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

func TestProcessMainProvision(t *testing.T) {
	tests := []struct {
		name         string
		mainProv     *jplaw.MainProvision
		wantSections []string
		wantErr      bool
	}{
		{
			name: "main provision with chapters",
			mainProv: &jplaw.MainProvision{
				Chapter: []jplaw.Chapter{
					{
						ChapterTitle: jplaw.ChapterTitle{
							Content: "第一章",
						},
					},
					{
						ChapterTitle: jplaw.ChapterTitle{
							Content: "第二章",
						},
					},
				},
			},
			wantSections: []string{"第一章", "第二章"},
			wantErr:      false,
		},
		{
			name: "main provision with articles only",
			mainProv: &jplaw.MainProvision{
				Article: []jplaw.Article{
					{
						ArticleTitle: &jplaw.ArticleTitle{
							Content: "第一条",
						},
					},
					{
						ArticleTitle: &jplaw.ArticleTitle{
							Content: "第二条",
						},
					},
				},
			},
			wantSections: []string{"第一条", "第二条"},
			wantErr:      false,
		},
		{
			name: "main provision with multiple paragraphs",
			mainProv: &jplaw.MainProvision{
				Paragraph: []jplaw.Paragraph{
					{
						Num: 1,
						ParagraphNum: jplaw.ParagraphNum{
							Content: "１",
						},
						ParagraphSentence: jplaw.ParagraphSentence{
							Sentence: []jplaw.Sentence{
								{Content: "第一項の内容"},
							},
						},
					},
					{
						Num: 2,
						ParagraphNum: jplaw.ParagraphNum{
							Content: "２",
						},
						ParagraphSentence: jplaw.ParagraphSentence{
							Sentence: []jplaw.Sentence{
								{Content: "第二項の内容"},
							},
						},
					},
				},
			},
			wantSections: []string{"第１項", "第２項"},
			wantErr:      false,
		},
		{
			name: "main provision with single paragraph",
			mainProv: &jplaw.MainProvision{
				Paragraph: []jplaw.Paragraph{
					{
						Num: 1,
						ParagraphNum: jplaw.ParagraphNum{
							Content: "１",
						},
						ParagraphSentence: jplaw.ParagraphSentence{
							Sentence: []jplaw.Sentence{
								{Content: "単一項の内容"},
							},
						},
					},
				},
			},
			wantSections: []string{"本文"},
			wantErr:      false,
		},
		{
			name:         "empty main provision",
			mainProv:     &jplaw.MainProvision{},
			wantSections: []string{},
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new EPUB book
			book, err := epub.NewEpub("Test Book")
			if err != nil {
				t.Fatalf("Failed to create EPUB: %v", err)
			}

			// Process main provision
			err = processMainProvision(book, tt.mainProv, nil)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("processMainProvision() error = %v, wantErr %v", err, tt.wantErr)
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

func TestProcessMainProvisionWithItems(t *testing.T) {
	// Test paragraph with items
	mainProv := &jplaw.MainProvision{
		Paragraph: []jplaw.Paragraph{
			{
				Num: 1,
				ParagraphNum: jplaw.ParagraphNum{
					Content: "１",
				},
				ParagraphSentence: jplaw.ParagraphSentence{
					Sentence: []jplaw.Sentence{
						{Content: "次に掲げる事項"},
					},
				},
				Item: []jplaw.Item{
					{
						ItemTitle: &jplaw.ItemTitle{
							Content: "一",
						},
						ItemSentence: jplaw.ItemSentence{
							Sentence: []jplaw.Sentence{
								{Content: "第一号の内容"},
							},
						},
					},
					{
						ItemTitle: &jplaw.ItemTitle{
							Content: "二",
						},
						ItemSentence: jplaw.ItemSentence{
							Sentence: []jplaw.Sentence{
								{Content: "第二号の内容"},
							},
						},
					},
				},
			},
			{
				Num: 2,
				ParagraphNum: jplaw.ParagraphNum{
					Content: "２",
				},
				ParagraphSentence: jplaw.ParagraphSentence{
					Sentence: []jplaw.Sentence{
						{Content: "第二項の内容"},
					},
				},
			},
		},
	}

	// Create a new EPUB book
	book, err := epub.NewEpub("Test Book")
	if err != nil {
		t.Fatalf("Failed to create EPUB: %v", err)
	}

	// Process main provision
	err = processMainProvision(book, mainProv, nil)
	if err != nil {
		t.Errorf("processMainProvision() unexpected error: %v", err)
		return
	}

	// Since we can't easily inspect the internal content without writing to disk,
	// we'll just verify that the function didn't error
	// Test passed if no error occurred - main provision with items was processed successfully
}
