package jplaw2epub

import (
	"strings"
	"testing"

	"go.ngs.io/jplaw-xml"
)

// Helper function to create test sentence with MixedContent
func createTestSentenceWithContent(text string) jplaw.Sentence {
	return jplaw.Sentence{
		MixedContent: jplaw.MixedContent{
			Nodes: []jplaw.ContentNode{
				jplaw.TextNode{Text: text},
			},
		},
	}
}

func TestProcessLists(t *testing.T) {
	tests := []struct {
		name     string
		lists    []jplaw.List
		want     string
		contains []string
	}{
		{
			name:  "empty lists",
			lists: []jplaw.List{},
			want:  "",
		},
		{
			name: "single list with sentence",
			lists: []jplaw.List{
				{
					ListSentence: jplaw.ListSentence{
						Sentence: []jplaw.Sentence{
							createTestSentenceWithContent("リスト項目1"),
						},
					},
				},
			},
			contains: []string{
				`<ul class="law-list">`,
				`<li>`,
				"リスト項目1",
				`</li>`,
				`</ul>`,
			},
		},
		{
			name: "multiple lists",
			lists: []jplaw.List{
				{
					ListSentence: jplaw.ListSentence{
						Sentence: []jplaw.Sentence{
							createTestSentenceWithContent("項目A"),
						},
					},
				},
				{
					ListSentence: jplaw.ListSentence{
						Sentence: []jplaw.Sentence{
							createTestSentenceWithContent("項目B"),
						},
					},
				},
			},
			contains: []string{
				`<ul class="law-list">`,
				"項目A",
				"項目B",
				`</ul>`,
			},
		},
		{
			name: "list with sublist1",
			lists: []jplaw.List{
				{
					ListSentence: jplaw.ListSentence{
						Sentence: []jplaw.Sentence{
							createTestSentenceWithContent("親リスト"),
						},
					},
					Sublist1: []jplaw.Sublist1{
						{
							Sublist1Sentence: jplaw.Sublist1Sentence{
								Sentence: []jplaw.Sentence{
									createTestSentenceWithContent("サブリスト1"),
								},
							},
						},
					},
				},
			},
			contains: []string{
				"親リスト",
				`<ul class="law-sublist1">`,
				"サブリスト1",
			},
		},
		{
			name: "list with columns",
			lists: []jplaw.List{
				{
					ListSentence: jplaw.ListSentence{
						Sentence: []jplaw.Sentence{
							createTestSentenceWithContent("リスト内容"),
						},
						Column: []jplaw.Column{
							{
								Num:       1,
								LineBreak: false,
								Sentence: []jplaw.Sentence{
									createTestSentenceWithContent("欄1"),
								},
							},
						},
					},
				},
			},
			contains: []string{
				"リスト内容",
				"欄1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processLists(tt.lists)

			if tt.want != "" && got != tt.want {
				t.Errorf("processLists() = %v, want %v", got, tt.want)
			}

			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("processLists() missing expected content: %s", want)
				}
			}
		})
	}
}

func TestProcessSublist1(t *testing.T) {
	tests := []struct {
		name     string
		sublists []jplaw.Sublist1
		want     string
		contains []string
	}{
		{
			name:     "empty sublist1",
			sublists: []jplaw.Sublist1{},
			want:     "",
		},
		{
			name: "single sublist1",
			sublists: []jplaw.Sublist1{
				{
					Sublist1Sentence: jplaw.Sublist1Sentence{
						Sentence: []jplaw.Sentence{
							createTestSentenceWithContent("サブリスト1項目"),
						},
					},
				},
			},
			contains: []string{
				`<ul class="law-sublist1">`,
				`<li>`,
				"サブリスト1項目",
				`</li>`,
				`</ul>`,
			},
		},
		{
			name: "sublist1 with sublist2",
			sublists: []jplaw.Sublist1{
				{
					Sublist1Sentence: jplaw.Sublist1Sentence{
						Sentence: []jplaw.Sentence{
							createTestSentenceWithContent("レベル1"),
						},
					},
					Sublist2: []jplaw.Sublist2{
						{
							Sublist2Sentence: jplaw.Sublist2Sentence{
								Sentence: []jplaw.Sentence{
									createTestSentenceWithContent("レベル2"),
								},
							},
						},
					},
				},
			},
			contains: []string{
				`<ul class="law-sublist1">`,
				"レベル1",
				`<ul class="law-sublist2">`,
				"レベル2",
			},
		},
		{
			name: "sublist1 with columns",
			sublists: []jplaw.Sublist1{
				{
					Sublist1Sentence: jplaw.Sublist1Sentence{
						Sentence: []jplaw.Sentence{
							createTestSentenceWithContent("サブリスト内容"),
						},
						Column: []jplaw.Column{
							{
								Num: 1,
								Sentence: []jplaw.Sentence{
									createTestSentenceWithContent("欄内容"),
								},
							},
						},
					},
				},
			},
			contains: []string{
				"サブリスト内容",
				"欄内容",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processSublist1(tt.sublists)

			if tt.want != "" && got != tt.want {
				t.Errorf("processSublist1() = %v, want %v", got, tt.want)
			}

			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("processSublist1() missing expected content: %s", want)
				}
			}
		})
	}
}

func TestProcessSublist2(t *testing.T) {
	tests := []struct {
		name     string
		sublists []jplaw.Sublist2
		want     string
		contains []string
	}{
		{
			name:     "empty sublist2",
			sublists: []jplaw.Sublist2{},
			want:     "",
		},
		{
			name: "single sublist2",
			sublists: []jplaw.Sublist2{
				{
					Sublist2Sentence: jplaw.Sublist2Sentence{
						Sentence: []jplaw.Sentence{
							createTestSentenceWithContent("サブリスト2項目"),
						},
					},
				},
			},
			contains: []string{
				`<ul class="law-sublist2">`,
				`<li>`,
				"サブリスト2項目",
				`</li>`,
				`</ul>`,
			},
		},
		{
			name: "sublist2 with sublist3",
			sublists: []jplaw.Sublist2{
				{
					Sublist2Sentence: jplaw.Sublist2Sentence{
						Sentence: []jplaw.Sentence{
							createTestSentenceWithContent("レベル2"),
						},
					},
					Sublist3: []jplaw.Sublist3{
						{
							Sublist3Sentence: jplaw.Sublist3Sentence{
								Sentence: []jplaw.Sentence{
									createTestSentenceWithContent("レベル3"),
								},
							},
						},
					},
				},
			},
			contains: []string{
				`<ul class="law-sublist2">`,
				"レベル2",
				`<ul class="law-sublist3">`,
				"レベル3",
			},
		},
		{
			name: "multiple sublist2 items",
			sublists: []jplaw.Sublist2{
				{
					Sublist2Sentence: jplaw.Sublist2Sentence{
						Sentence: []jplaw.Sentence{
							createTestSentenceWithContent("項目2-1"),
						},
					},
				},
				{
					Sublist2Sentence: jplaw.Sublist2Sentence{
						Sentence: []jplaw.Sentence{
							createTestSentenceWithContent("項目2-2"),
						},
					},
				},
			},
			contains: []string{
				"項目2-1",
				"項目2-2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processSublist2(tt.sublists)

			if tt.want != "" && got != tt.want {
				t.Errorf("processSublist2() = %v, want %v", got, tt.want)
			}

			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("processSublist2() missing expected content: %s\nGot: %s", want, got)
				}
			}
		})
	}
}

func TestProcessSublist3(t *testing.T) {
	tests := []struct {
		name     string
		sublists []jplaw.Sublist3
		want     string
		contains []string
	}{
		{
			name:     "empty sublist3",
			sublists: []jplaw.Sublist3{},
			want:     "",
		},
		{
			name: "single sublist3",
			sublists: []jplaw.Sublist3{
				{
					Sublist3Sentence: jplaw.Sublist3Sentence{
						Sentence: []jplaw.Sentence{
							createTestSentenceWithContent("サブリスト3項目"),
						},
					},
				},
			},
			contains: []string{
				`<ul class="law-sublist3">`,
				`<li>`,
				"サブリスト3項目",
				`</li>`,
				`</ul>`,
			},
		},
		{
			name: "multiple sublist3 items",
			sublists: []jplaw.Sublist3{
				{
					Sublist3Sentence: jplaw.Sublist3Sentence{
						Sentence: []jplaw.Sentence{
							createTestSentenceWithContent("項目3-1"),
						},
					},
				},
				{
					Sublist3Sentence: jplaw.Sublist3Sentence{
						Sentence: []jplaw.Sentence{
							createTestSentenceWithContent("項目3-2"),
						},
					},
				},
			},
			contains: []string{
				`<ul class="law-sublist3">`,
				"項目3-1",
				"項目3-2",
				`</ul>`,
			},
		},
		{
			name: "sublist3 with columns",
			sublists: []jplaw.Sublist3{
				{
					Sublist3Sentence: jplaw.Sublist3Sentence{
						Sentence: []jplaw.Sentence{
							createTestSentenceWithContent("リスト3内容"),
						},
						Column: []jplaw.Column{
							{
								Num: 1,
								Sentence: []jplaw.Sentence{
									createTestSentenceWithContent("欄3内容"),
								},
							},
						},
					},
				},
			},
			contains: []string{
				"リスト3内容",
				"欄3内容",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processSublist3(tt.sublists)

			if tt.want != "" && got != tt.want {
				t.Errorf("processSublist3() = %v, want %v", got, tt.want)
			}

			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("processSublist3() missing expected content: %s", want)
				}
			}
		})
	}
}

// Test for nested list structure
func TestProcessListsNested(t *testing.T) {
	// Create a deeply nested list structure
	lists := []jplaw.List{
		{
			ListSentence: jplaw.ListSentence{
				Sentence: []jplaw.Sentence{
					createTestSentenceWithContent("レベル0"),
				},
			},
			Sublist1: []jplaw.Sublist1{
				{
					Sublist1Sentence: jplaw.Sublist1Sentence{
						Sentence: []jplaw.Sentence{
							createTestSentenceWithContent("レベル1"),
						},
					},
					Sublist2: []jplaw.Sublist2{
						{
							Sublist2Sentence: jplaw.Sublist2Sentence{
								Sentence: []jplaw.Sentence{
									createTestSentenceWithContent("レベル2"),
								},
							},
							Sublist3: []jplaw.Sublist3{
								{
									Sublist3Sentence: jplaw.Sublist3Sentence{
										Sentence: []jplaw.Sentence{
											createTestSentenceWithContent("レベル3"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	got := processLists(lists)

	expectedContent := []string{
		`<ul class="law-list">`,
		"レベル0",
		`<ul class="law-sublist1">`,
		"レベル1",
		`<ul class="law-sublist2">`,
		"レベル2",
		`<ul class="law-sublist3">`,
		"レベル3",
	}

	for _, want := range expectedContent {
		if !strings.Contains(got, want) {
			t.Errorf("processLists() missing expected nested content: %s", want)
		}
	}
}

// Benchmark tests
func BenchmarkProcessLists(b *testing.B) {
	lists := []jplaw.List{
		{
			ListSentence: jplaw.ListSentence{
				Sentence: []jplaw.Sentence{
					createTestSentenceWithContent("ベンチマークリスト"),
				},
			},
			Sublist1: []jplaw.Sublist1{
				{
					Sublist1Sentence: jplaw.Sublist1Sentence{
						Sentence: []jplaw.Sentence{
							createTestSentenceWithContent("サブリスト"),
						},
					},
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processLists(lists)
	}
}

func BenchmarkProcessSublist1(b *testing.B) {
	sublists := []jplaw.Sublist1{
		{
			Sublist1Sentence: jplaw.Sublist1Sentence{
				Sentence: []jplaw.Sentence{
					createTestSentenceWithContent("ベンチマーク項目1"),
					createTestSentenceWithContent("ベンチマーク項目2"),
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processSublist1(sublists)
	}
}

// Test for edge cases in processNumberedParagraph
func TestProcessNumberedParagraphEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		para       *jplaw.Paragraph
		nextPara   *jplaw.Paragraph
		wantInList bool
		contains   []string
	}{
		{
			name: "numbered paragraph with list",
			para: &jplaw.Paragraph{
				Num:          1,
				ParagraphNum: jplaw.ParagraphNum{Content: "1"},
				ParagraphSentence: jplaw.ParagraphSentence{
					Sentence: []jplaw.Sentence{
						createTestSentenceWithContent("段落内容"),
					},
				},
				List: []jplaw.List{
					{
						ListSentence: jplaw.ListSentence{
							Sentence: []jplaw.Sentence{
								createTestSentenceWithContent("リスト項目"),
							},
						},
					},
				},
			},
			nextPara: &jplaw.Paragraph{
				Num: 2,
			},
			wantInList: true,
			contains: []string{
				"段落内容",
				"リスト項目",
				`<ul class="law-list">`,
			},
		},
		{
			name: "numbered paragraph with table",
			para: &jplaw.Paragraph{
				Num:          1,
				ParagraphNum: jplaw.ParagraphNum{Content: "1"},
				ParagraphSentence: jplaw.ParagraphSentence{
					Sentence: []jplaw.Sentence{
						createTestSentenceWithContent("表を含む段落"),
					},
				},
				TableStruct: []jplaw.TableStruct{
					{
						TableStructTitle: &jplaw.TableStructTitle{
							Content: "表1",
						},
					},
				},
			},
			nextPara:   nil,
			wantInList: true, // processNumberedParagraph always sets inList to true
			contains: []string{
				"表を含む段落",
				"表1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &paragraphProcessor{
				inList: false,
			}

			// Create a slice with the test paragraph
			paragraphs := []jplaw.Paragraph{*tt.para}
			if tt.nextPara != nil {
				paragraphs = append(paragraphs, *tt.nextPara)
			}

			p.processNumberedParagraph(tt.para, 0, paragraphs)

			if p.inList != tt.wantInList {
				t.Errorf("processNumberedParagraph() inList = %v, want %v", p.inList, tt.wantInList)
			}

			for _, want := range tt.contains {
				if !strings.Contains(p.body, want) {
					t.Errorf("processNumberedParagraph() missing expected content: %s", want)
				}
			}
		})
	}
}
