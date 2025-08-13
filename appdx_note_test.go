package jplaw2epub

import (
	"strings"
	"testing"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

func TestProcessAppdxNotes(t *testing.T) {
	tests := []struct {
		name    string
		notes   []jplaw.AppdxNote
		wantErr bool
	}{
		{
			name:    "Empty notes",
			notes:   []jplaw.AppdxNote{},
			wantErr: false,
		},
		{
			name: "Single note with title",
			notes: []jplaw.AppdxNote{
				{
					AppdxNoteTitle: &jplaw.AppdxNoteTitle{
						Content: "附則",
					},
					NoteStruct: []jplaw.NoteStruct{
						{
							Note: jplaw.Note{
								Content: "<Paragraph><ParagraphSentence><Sentence>テスト</Sentence></ParagraphSentence></Paragraph>",
							},
						},
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

			err = processAppdxNotes(book, tt.notes, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("processAppdxNotes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProcessAppdxNote(t *testing.T) {
	tests := []struct {
		name    string
		note    *jplaw.AppdxNote
		wantErr bool
	}{
		{
			name: "Note with title and content",
			note: &jplaw.AppdxNote{
				AppdxNoteTitle: &jplaw.AppdxNoteTitle{
					Content: "附則",
				},
				NoteStruct: []jplaw.NoteStruct{
					{
						NoteStructTitle: &jplaw.NoteStructTitle{
							Content: "詔書",
						},
						Note: jplaw.Note{
							Content: "<Paragraph><ParagraphSentence><Sentence>内容</Sentence></ParagraphSentence></Paragraph>",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Note with related article",
			note: &jplaw.AppdxNote{
				AppdxNoteTitle: &jplaw.AppdxNoteTitle{
					Content: "附則",
				},
				RelatedArticleNum: &jplaw.RelatedArticleNum{
					Content: "第一条関係",
				},
				NoteStruct: []jplaw.NoteStruct{
					{
						Note: jplaw.Note{
							Content: "<Paragraph><ParagraphSentence><Sentence>内容</Sentence></ParagraphSentence></Paragraph>",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Note with remarks",
			note: &jplaw.AppdxNote{
				AppdxNoteTitle: &jplaw.AppdxNoteTitle{
					Content: "附則",
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book, err := epub.NewEpub("Test Book")
			if err != nil {
				t.Fatalf("Failed to create epub: %v", err)
			}

			err = processAppdxNote(book, tt.note, 0, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("processAppdxNote() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProcessNoteStruct(t *testing.T) {
	tests := []struct {
		name     string
		note     *jplaw.NoteStruct
		contains []string
	}{
		{
			name: "Note with title and content",
			note: &jplaw.NoteStruct{
				NoteStructTitle: &jplaw.NoteStructTitle{
					Content: "タイトル",
				},
				Note: jplaw.Note{
					Content: "<Paragraph><ParagraphSentence><Sentence>内容</Sentence></ParagraphSentence></Paragraph>",
				},
			},
			contains: []string{
				"<h3>タイトル</h3>",
				"内容",
			},
		},
		{
			name: "Note with remarks",
			note: &jplaw.NoteStruct{
				Note: jplaw.Note{
					Content: "<Paragraph><ParagraphSentence><Sentence>内容</Sentence></ParagraphSentence></Paragraph>",
				},
				Remarks: []jplaw.Remarks{
					{
						RemarksLabel: jplaw.RemarksLabel{
							Content: "備考",
						},
						Sentence: []jplaw.Sentence{
							createTestSentence("備考内容"),
						},
					},
				},
			},
			contains: []string{
				"内容",
				"備考",
				"備考内容",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processNoteStruct(tt.note, nil)
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("processNoteStruct() should contain %q\ngot: %v", expected, result)
				}
			}
		})
	}
}

func TestProcessNoteContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		contains []string
	}{
		{
			name: "Simple paragraph",
			content: `<Paragraph Num="1">
				<ParagraphSentence>
					<Sentence>これはテストです。</Sentence>
				</ParagraphSentence>
			</Paragraph>`,
			contains: []string{
				"これはテストです。",
			},
		},
		{
			name: "Multiple paragraphs",
			content: `<Paragraph Num="1">
				<ParagraphSentence>
					<Sentence>第一段落</Sentence>
				</ParagraphSentence>
			</Paragraph>
			<Paragraph Num="2">
				<ParagraphSentence>
					<Sentence>第二段落</Sentence>
				</ParagraphSentence>
			</Paragraph>`,
			contains: []string{
				"第一段落",
				"第二段落",
			},
		},
		{
			name: "Paragraph with items",
			content: `<Paragraph Num="1">
				<ParagraphSentence>
					<Sentence>前文</Sentence>
				</ParagraphSentence>
				<Item Num="1">
					<ItemTitle>一</ItemTitle>
					<ItemSentence>
						<Sentence>項目一</Sentence>
					</ItemSentence>
				</Item>
			</Paragraph>`,
			contains: []string{
				"前文",
				"項目一",
			},
		},
		{
			name:     "Invalid XML",
			content:  "Not valid XML",
			contains: []string{"Not valid XML"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processNoteContent(tt.content, nil)
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("processNoteContent() should contain %q\ngot: %v", expected, result)
				}
			}
		})
	}
}

func TestProcessRemarks(t *testing.T) {
	tests := []struct {
		name     string
		remarks  *jplaw.Remarks
		contains []string
	}{
		{
			name: "Remarks with label and sentences",
			remarks: &jplaw.Remarks{
				RemarksLabel: jplaw.RemarksLabel{
					Content: "備考",
				},
				Sentence: []jplaw.Sentence{
					createTestSentence("備考1"),
					createTestSentence("備考2"),
				},
			},
			contains: []string{
				"備考",
				"備考1",
				"備考2",
				`class="appdx-remarks"`,
			},
		},
		{
			name: "Remarks with items",
			remarks: &jplaw.Remarks{
				RemarksLabel: jplaw.RemarksLabel{
					Content: "注記",
				},
				Item: []jplaw.Item{
					{
						ItemTitle: &jplaw.ItemTitle{Content: "一"},
						ItemSentence: jplaw.ItemSentence{
							Sentence: []jplaw.Sentence{
								createTestSentence("項目内容"),
							},
						},
					},
				},
			},
			contains: []string{
				"注記",
				"項目内容",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processRemarks(tt.remarks)
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("processRemarks() should contain %q\ngot: %v", expected, result)
				}
			}
		})
	}
}

func TestProcessAppdxTables(t *testing.T) {
	tests := []struct {
		name    string
		tables  []jplaw.AppdxTable
		wantErr bool
	}{
		{
			name:    "Empty tables",
			tables:  []jplaw.AppdxTable{},
			wantErr: false,
		},
		{
			name: "Single table with title",
			tables: []jplaw.AppdxTable{
				{
					AppdxTableTitle: &jplaw.AppdxTableTitle{
						Content: "附表",
					},
					TableStruct: []jplaw.TableStruct{
						{
							TableStructTitle: &jplaw.TableStructTitle{
								Content: "表1",
							},
						},
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

			err = processAppdxTables(book, tt.tables)
			if (err != nil) != tt.wantErr {
				t.Errorf("processAppdxTables() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProcessAppdxTable(t *testing.T) {
	tests := []struct {
		name    string
		table   *jplaw.AppdxTable
		wantErr bool
	}{
		{
			name: "Table with title and structure",
			table: &jplaw.AppdxTable{
				AppdxTableTitle: &jplaw.AppdxTableTitle{
					Content: "附表",
				},
				TableStruct: []jplaw.TableStruct{
					{
						TableStructTitle: &jplaw.TableStructTitle{
							Content: "表1",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Table with related article",
			table: &jplaw.AppdxTable{
				AppdxTableTitle: &jplaw.AppdxTableTitle{
					Content: "附表",
				},
				RelatedArticleNum: &jplaw.RelatedArticleNum{
					Content: "第一条関係",
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

			err = processAppdxTable(book, tt.table, 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("processAppdxTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}