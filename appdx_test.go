package jplaw2epub

import (
	"strings"
	"testing"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

func TestProcessAppdxStyles(t *testing.T) {
	tests := []struct {
		name    string
		styles  []jplaw.AppdxStyle
		wantErr bool
	}{
		{
			name:    "Empty styles",
			styles:  []jplaw.AppdxStyle{},
			wantErr: false,
		},
		{
			name: "Single style with title",
			styles: []jplaw.AppdxStyle{
				{
					AppdxStyleTitle: &jplaw.AppdxStyleTitle{
						Content: "様式",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Multiple styles",
			styles: []jplaw.AppdxStyle{
				{
					AppdxStyleTitle: &jplaw.AppdxStyleTitle{
						Content: "様式1",
					},
				},
				{
					AppdxStyleTitle: &jplaw.AppdxStyleTitle{
						Content: "様式2",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new book for each test to avoid conflicts
			book, err := epub.NewEpub("Test Book")
			if err != nil {
				t.Fatalf("Failed to create epub: %v", err)
			}

			err = processAppdxStyles(book, tt.styles, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("processAppdxStyles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProcessAppdxStyle(t *testing.T) {
	book, _ := epub.NewEpub("Test Book")

	// Add a parent section that the subsections can reference
	book.AddSection("<p>Parent Section</p>", "Parent", "parent.xhtml", "")

	tests := []struct {
		name    string
		style   *jplaw.AppdxStyle
		wantErr bool
	}{
		{
			name: "Style with title and content",
			style: &jplaw.AppdxStyle{
				AppdxStyleTitle: &jplaw.AppdxStyleTitle{
					Content: "様式タイトル",
				},
				StyleStruct: []jplaw.StyleStruct{
					{
						StyleStructTitle: &jplaw.StyleStructTitle{
							Content: "構造タイトル",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Style with related article",
			style: &jplaw.AppdxStyle{
				AppdxStyleTitle: &jplaw.AppdxStyleTitle{
					Content: "様式",
				},
				RelatedArticleNum: &jplaw.RelatedArticleNum{
					Content: "第一条関係",
				},
			},
			wantErr: false,
		},
		{
			name: "Style with remarks",
			style: &jplaw.AppdxStyle{
				AppdxStyleTitle: &jplaw.AppdxStyleTitle{
					Content: "様式",
				},
				Remarks: &jplaw.Remarks{
					RemarksLabel: jplaw.RemarksLabel{
						Content: "備考",
					},
					Sentence: []jplaw.Sentence{
						createTestSentence("備考内容"),
					},
				},
			},
			wantErr: false,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := processAppdxStyle(book, tt.style, "parent.xhtml", i, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("processAppdxStyle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProcessAppdxRemark(t *testing.T) {
	tests := []struct {
		name     string
		remark   *jplaw.Remarks
		contains []string
	}{
		{
			name: "Remark with label",
			remark: &jplaw.Remarks{
				RemarksLabel: jplaw.RemarksLabel{
					Content: "備考",
				},
				Sentence: []jplaw.Sentence{
					createTestSentence("備考の内容"),
				},
			},
			contains: []string{
				"appdx-remarks",
				"remarks-label",
				"備考",
				"備考の内容",
			},
		},
		{
			name: "Remark with items",
			remark: &jplaw.Remarks{
				Item: []jplaw.Item{
					{
						ItemTitle: &jplaw.ItemTitle{Content: "一"},
						ItemSentence: jplaw.ItemSentence{
							Sentence: []jplaw.Sentence{
								createTestSentence("項目1"),
							},
						},
					},
				},
			},
			contains: []string{
				"項目1",
			},
		},
		{
			name:   "Empty remark",
			remark: &jplaw.Remarks{},
			contains: []string{
				"appdx-remarks",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processAppdxRemark(tt.remark, nil)
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("processAppdxRemark() should contain %q\ngot: %v", expected, result)
				}
			}
		})
	}
}

func TestProcessAppdxFig(t *testing.T) {
	book, _ := epub.NewEpub("Test Book")

	tests := []struct {
		name    string
		fig     []jplaw.AppdxFig
		wantErr bool
	}{
		{
			name:    "Empty figures",
			fig:     []jplaw.AppdxFig{},
			wantErr: false,
		},
		{
			name: "Single figure",
			fig: []jplaw.AppdxFig{
				{
					AppdxFigTitle: &jplaw.AppdxFigTitle{
						Content: "附図",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := processAppdxFig(book, tt.fig, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("processAppdxFig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProcessAppdxFigItem(t *testing.T) {
	book, _ := epub.NewEpub("Test Book")

	// Add a parent section that the figure subsection can reference
	book.AddSection("<p>Parent Section</p>", "Parent", "parent.xhtml", "")

	fig := &jplaw.AppdxFig{
		AppdxFigTitle: &jplaw.AppdxFigTitle{
			Content: "附図タイトル",
		},
		FigStruct: []jplaw.FigStruct{
			{
				Fig: jplaw.Fig{
					Src: "test.png",
				},
			},
		},
	}

	err := processAppdxFigItem(book, fig, "parent.xhtml", 0, nil)
	if err != nil {
		t.Errorf("processAppdxFigItem() unexpected error = %v", err)
	}
}
