package jplaw2epub

import (
	"strings"
	"testing"

	"github.com/go-shiori/go-epub"
	jplaw "go.ngs.io/jplaw-xml"
)

func TestProcessSupplProvisions(t *testing.T) {
	tests := []struct {
		name    string
		supplPr []jplaw.SupplProvision
		wantErr bool
	}{
		{
			name:    "Empty provisions",
			supplPr: []jplaw.SupplProvision{},
			wantErr: false,
		},
		{
			name: "Single provision with title",
			supplPr: []jplaw.SupplProvision{
				{
					AmendLawNum: "1",
					SupplProvisionLabel: jplaw.SupplProvisionLabel{
						Content: "附則",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Multiple provisions",
			supplPr: []jplaw.SupplProvision{
				{
					AmendLawNum: "1",
					SupplProvisionLabel: jplaw.SupplProvisionLabel{
						Content: "附則",
					},
				},
				{
					AmendLawNum: "2",
					SupplProvisionLabel: jplaw.SupplProvisionLabel{
						Content: "附則第二",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book, err := epub.NewEpub("Test Book")
			if err != nil {
				t.Fatalf("Failed to create epub: %v", err)
			}

			imgProc := &ImageProcessor{}
			err = processSupplProvisions(book, tt.supplPr, imgProc)
			if (err != nil) != tt.wantErr {
				t.Errorf("processSupplProvisions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProcessSupplProvision(t *testing.T) {
	tests := []struct {
		name     string
		provision *jplaw.SupplProvision
		idx      int
		wantErr  bool
	}{
		{
			name: "Provision with articles",
			provision: &jplaw.SupplProvision{
				AmendLawNum: "1",
				SupplProvisionLabel: jplaw.SupplProvisionLabel{
					Content: "附則",
				},
				Article: []jplaw.Article{
					{
						Num: "1",
						ArticleTitle: &jplaw.ArticleTitle{
							Content: "Test Article",
						},
						Paragraph: []jplaw.Paragraph{
							{
								Num: 1,
								ParagraphSentence: jplaw.ParagraphSentence{
									Sentence: []jplaw.Sentence{createTestSentence("Test content")},
								},
							},
						},
					},
				},
			},
			idx:     0,
			wantErr: false,
		},
		{
			name: "Provision with paragraphs",
			provision: &jplaw.SupplProvision{
				AmendLawNum: "2",
				SupplProvisionLabel: jplaw.SupplProvisionLabel{
					Content: "附則",
				},
				Paragraph: []jplaw.Paragraph{
					{
						Num: 1,
						ParagraphSentence: jplaw.ParagraphSentence{
							Sentence: []jplaw.Sentence{createTestSentence("Test paragraph")},
						},
					},
				},
			},
			idx:     1,
			wantErr: false,
		},
		{
			name: "Provision with appendix table",
			provision: &jplaw.SupplProvision{
				AmendLawNum: "3",
				SupplProvisionLabel: jplaw.SupplProvisionLabel{
					Content: "附則",
				},
				SupplProvisionAppdxTable: []jplaw.SupplProvisionAppdxTable{
					{
						SupplProvisionAppdxTableTitle: jplaw.SupplProvisionAppdxTableTitle{
							Content: "Test Table",
						},
					},
				},
			},
			idx:     2,
			wantErr: false,
		},
		{
			name: "Provision with appendix style",
			provision: &jplaw.SupplProvision{
				AmendLawNum: "4",
				SupplProvisionLabel: jplaw.SupplProvisionLabel{
					Content: "附則",
				},
				SupplProvisionAppdxStyle: []jplaw.SupplProvisionAppdxStyle{
					{
						SupplProvisionAppdxStyleTitle: jplaw.SupplProvisionAppdxStyleTitle{
							Content: "Test Style",
						},
					},
				},
			},
			idx:     3,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book, err := epub.NewEpub("Test Book")
			if err != nil {
				t.Fatalf("Failed to create epub: %v", err)
			}

			imgProc := &ImageProcessor{}
			err = processSupplProvision(book, tt.provision, tt.idx, imgProc)
			if (err != nil) != tt.wantErr {
				t.Errorf("processSupplProvision() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProcessSupplProvisionAppdxTable(t *testing.T) {
	tests := []struct {
		name     string
		table    *jplaw.SupplProvisionAppdxTable
		contains []string
	}{
		{
			name: "Table with title",
			table: &jplaw.SupplProvisionAppdxTable{
				SupplProvisionAppdxTableTitle: jplaw.SupplProvisionAppdxTableTitle{
					Content: "Test Table",
				},
			},
			contains: []string{"Test Table", "suppl-appdx-table"},
		},
		{
			name: "Table with related article",
			table: &jplaw.SupplProvisionAppdxTable{
				RelatedArticleNum: &jplaw.RelatedArticleNum{
					Content: "第1条",
				},
			},
			contains: []string{"第1条", "related-articles"},
		},
		{
			name: "Table with TableStruct",
			table: &jplaw.SupplProvisionAppdxTable{
				TableStruct: []jplaw.TableStruct{
					{
						TableStructTitle: &jplaw.TableStructTitle{
							Content: "Inner Table",
						},
					},
				},
			},
			contains: []string{"suppl-appdx-table"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imgProc := &ImageProcessor{}
			result := processSupplProvisionAppdxTable(tt.table, imgProc)
			
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected result to contain %q, but it didn't", expected)
				}
			}
		})
	}
}

func TestProcessSupplProvisionAppdxStyle(t *testing.T) {
	tests := []struct {
		name     string
		style    *jplaw.SupplProvisionAppdxStyle
		contains []string
	}{
		{
			name: "Style with title",
			style: &jplaw.SupplProvisionAppdxStyle{
				SupplProvisionAppdxStyleTitle: jplaw.SupplProvisionAppdxStyleTitle{
					Content: "Test Style",
				},
			},
			contains: []string{"Test Style", "suppl-appdx-style"},
		},
		{
			name: "Style with related article",
			style: &jplaw.SupplProvisionAppdxStyle{
				RelatedArticleNum: &jplaw.RelatedArticleNum{
					Content: "第2条",
				},
			},
			contains: []string{"第2条", "related-articles"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imgProc := &ImageProcessor{}
			result := processSupplProvisionAppdxStyle(tt.style, imgProc)
			
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected result to contain %q, but it didn't", expected)
				}
			}
		})
	}
}

func TestProcessSupplProvisionAppdx(t *testing.T) {
	tests := []struct {
		name     string
		appdx    *jplaw.SupplProvisionAppdx
		contains []string
	}{
		{
			name: "Appendix with ArithFormulaNum",
			appdx: &jplaw.SupplProvisionAppdx{
				ArithFormulaNum: &jplaw.ArithFormulaNum{
					Content: "Formula 1",
				},
			},
			contains: []string{"Formula 1", "arith-formula-num"},
		},
		{
			name: "Appendix with RelatedArticleNum",
			appdx: &jplaw.SupplProvisionAppdx{
				RelatedArticleNum: &jplaw.RelatedArticleNum{
					Content: "第3条",
				},
			},
			contains: []string{"第3条", "related-articles"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imgProc := &ImageProcessor{}
			result := processSupplProvisionAppdx(tt.appdx, imgProc)
			
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected result to contain %q, but it didn't", expected)
				}
			}
		})
	}
}