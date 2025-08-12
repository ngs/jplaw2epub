package main

import (
	"strings"
	"testing"

	"go.ngs.io/jplaw-xml"
)

func TestProcessParagraphs(t *testing.T) {
	tests := []struct {
		name        string
		paragraphs  []jplaw.Paragraph
		contains    []string
		notContains []string
	}{
		{
			name:       "Empty paragraphs",
			paragraphs: []jplaw.Paragraph{},
			contains:   []string{""},
		},
		{
			name: "Single numbered paragraph",
			paragraphs: []jplaw.Paragraph{
				{
					Num:          1,
					ParagraphNum: jplaw.ParagraphNum{Content: "１"},
					ParagraphSentence: jplaw.ParagraphSentence{
						Sentence: []jplaw.Sentence{
							createTestSentence("第一項の内容。"),
						},
					},
				},
			},
			contains: []string{
				"<ol",
				"<li>",
				"第一項の内容。",
				"</li>",
				"</ol>",
			},
			notContains: []string{
				"<strong>１</strong>", // Should skip list number
			},
		},
		{
			name: "Multiple numbered paragraphs",
			paragraphs: []jplaw.Paragraph{
				{
					Num:          1,
					ParagraphNum: jplaw.ParagraphNum{Content: "１"},
					ParagraphSentence: jplaw.ParagraphSentence{
						Sentence: []jplaw.Sentence{
							createTestSentence("第一項。"),
						},
					},
				},
				{
					Num:          2,
					ParagraphNum: jplaw.ParagraphNum{Content: "２"},
					ParagraphSentence: jplaw.ParagraphSentence{
						Sentence: []jplaw.Sentence{
							createTestSentence("第二項。"),
						},
					},
				},
			},
			contains: []string{
				`<ol style="list-style-type: decimal;">`,
				"第一項。",
				"第二項。",
			},
		},
		{
			name: "Regular paragraph (Num=0)",
			paragraphs: []jplaw.Paragraph{
				{
					Num:          0,
					ParagraphNum: jplaw.ParagraphNum{Content: "前文"},
					ParagraphSentence: jplaw.ParagraphSentence{
						Sentence: []jplaw.Sentence{
							createTestSentence("これは前文です。"),
						},
					},
				},
			},
			contains: []string{
				"<h4>前文</h4>",
				"<p>",
				"これは前文です。",
				"</p>",
			},
			notContains: []string{
				"<ol>",
				"<li>",
			},
		},
		{
			name: "Mixed numbered and regular paragraphs",
			paragraphs: []jplaw.Paragraph{
				{
					Num:          1,
					ParagraphNum: jplaw.ParagraphNum{Content: "１"},
					ParagraphSentence: jplaw.ParagraphSentence{
						Sentence: []jplaw.Sentence{
							createTestSentence("番号付き。"),
						},
					},
				},
				{
					Num:          0,
					ParagraphNum: jplaw.ParagraphNum{Content: "補足"},
					ParagraphSentence: jplaw.ParagraphSentence{
						Sentence: []jplaw.Sentence{
							createTestSentence("通常の段落。"),
						},
					},
				},
			},
			contains: []string{
				"<ol",
				"番号付き。",
				"</ol>",
				"<h4>補足</h4>",
				"通常の段落。",
			},
		},
		{
			name: "Paragraph with items",
			paragraphs: []jplaw.Paragraph{
				{
					Num: 0,
					Item: []jplaw.Item{
						{
							ItemTitle: &jplaw.ItemTitle{Content: "一"},
							ItemSentence: jplaw.ItemSentence{
								Sentence: []jplaw.Sentence{
									createTestSentence("項目内容。"),
								},
							},
						},
					},
				},
			},
			contains: []string{
				"<ol",
				"項目内容。",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processParagraphs(tt.paragraphs)

			for _, contain := range tt.contains {
				if contain != "" && !strings.Contains(got, contain) {
					t.Errorf("processParagraphs() should contain %q\ngot: %v", contain, got)
				}
			}

			for _, notContain := range tt.notContains {
				if strings.Contains(got, notContain) {
					t.Errorf("processParagraphs() should not contain %q\ngot: %v", notContain, got)
				}
			}
		})
	}
}

func TestParagraphProcessor_processNumberedParagraph(t *testing.T) {
	p := &paragraphProcessor{}
	para := &jplaw.Paragraph{
		Num:          1,
		ParagraphNum: jplaw.ParagraphNum{Content: "特別条項"},
		ParagraphSentence: jplaw.ParagraphSentence{
			Sentence: []jplaw.Sentence{
				createTestSentence("特別な内容。"),
			},
		},
	}

	allParagraphs := []jplaw.Paragraph{*para}
	p.processNumberedParagraph(para, 0, allParagraphs)

	if !p.inList {
		t.Error("inList should be true after processing numbered paragraph")
	}

	if !strings.Contains(p.body, "<ol>") {
		t.Error("body should contain opening list tag")
	}

	if !strings.Contains(p.body, "<strong>特別条項</strong>") {
		t.Error("body should contain paragraph number as it's not a list number")
	}

	if !strings.Contains(p.body, "特別な内容。") {
		t.Error("body should contain paragraph content")
	}
}

func TestParagraphProcessor_processRegularParagraph(t *testing.T) {
	tests := []struct {
		name        string
		initialList bool
		para        *jplaw.Paragraph
		wantInList  bool
		contains    []string
	}{
		{
			name:        "Regular paragraph with list open",
			initialList: true,
			para: &jplaw.Paragraph{
				ParagraphNum: jplaw.ParagraphNum{Content: "補足"},
				ParagraphSentence: jplaw.ParagraphSentence{
					Sentence: []jplaw.Sentence{
						createTestSentence("補足内容。"),
					},
				},
			},
			wantInList: false,
			contains:   []string{"</ol>", "<h4>補足</h4>", "<p>", "補足内容。", "</p>"},
		},
		{
			name:        "Regular paragraph without list open",
			initialList: false,
			para: &jplaw.Paragraph{
				ParagraphSentence: jplaw.ParagraphSentence{
					Sentence: []jplaw.Sentence{
						createTestSentence("通常の文章。"),
					},
				},
			},
			wantInList: false,
			contains:   []string{"<p>", "通常の文章。", "</p>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &paragraphProcessor{
				inList: tt.initialList,
			}

			p.processRegularParagraph(tt.para)

			if p.inList != tt.wantInList {
				t.Errorf("inList = %v, want %v", p.inList, tt.wantInList)
			}

			for _, contain := range tt.contains {
				if !strings.Contains(p.body, contain) {
					t.Errorf("body should contain %q\ngot: %v", contain, p.body)
				}
			}
		})
	}
}

func TestParagraphProcessor_startNumberedList(t *testing.T) {
	tests := []struct {
		name       string
		idx        int
		paragraphs []jplaw.Paragraph
		want       string
	}{
		{
			name: "CJK numbered list",
			idx:  0,
			paragraphs: []jplaw.Paragraph{
				{Num: 1, ParagraphNum: jplaw.ParagraphNum{Content: "一"}},
				{Num: 2, ParagraphNum: jplaw.ParagraphNum{Content: "二"}},
				{Num: 0, ParagraphNum: jplaw.ParagraphNum{Content: "補足"}},
			},
			want: `<ol style="list-style-type: cjk-ideographic;">`,
		},
		{
			name: "Decimal numbered list",
			idx:  0,
			paragraphs: []jplaw.Paragraph{
				{Num: 1, ParagraphNum: jplaw.ParagraphNum{Content: "１"}},
				{Num: 2, ParagraphNum: jplaw.ParagraphNum{Content: "２"}},
			},
			want: `<ol style="list-style-type: decimal;">`,
		},
		{
			name: "Unknown pattern",
			idx:  0,
			paragraphs: []jplaw.Paragraph{
				{Num: 1, ParagraphNum: jplaw.ParagraphNum{Content: "A"}},
			},
			want: "<ol>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &paragraphProcessor{}
			p.startNumberedList(tt.idx, tt.paragraphs)

			if p.body != tt.want {
				t.Errorf("startNumberedList() = %v, want %v", p.body, tt.want)
			}
		})
	}
}

func TestParagraphProcessor_addParagraphNumber(t *testing.T) {
	tests := []struct {
		name string
		para *jplaw.Paragraph
		want string
	}{
		{
			name: "List number (should be skipped)",
			para: &jplaw.Paragraph{
				ParagraphNum: jplaw.ParagraphNum{Content: "一"},
			},
			want: "",
		},
		{
			name: "Non-list number",
			para: &jplaw.Paragraph{
				ParagraphNum: jplaw.ParagraphNum{Content: "特別項"},
			},
			want: "<strong>特別項</strong> ",
		},
		{
			name: "Empty paragraph number",
			para: &jplaw.Paragraph{
				ParagraphNum: jplaw.ParagraphNum{Content: ""},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &paragraphProcessor{}
			p.addParagraphNumber(tt.para)

			if p.body != tt.want {
				t.Errorf("addParagraphNumber() = %v, want %v", p.body, tt.want)
			}
		})
	}
}
