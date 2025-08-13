package jplaw2epub

import (
	"strings"
	"testing"

	jplaw "go.ngs.io/jplaw-xml"
)

func TestProcessTableStructs(t *testing.T) {
	tests := []struct {
		name     string
		tables   []jplaw.TableStruct
		contains []string
	}{
		{
			name:     "Empty tables",
			tables:   []jplaw.TableStruct{},
			contains: []string{},
		},
		{
			name: "Single table",
			tables: []jplaw.TableStruct{
				{
					TableStructTitle: &jplaw.TableStructTitle{Content: "Test Table"},
					Table: jplaw.Table{
						TableRow: []jplaw.TableRow{
							{
								TableColumn: []jplaw.TableColumn{
									{Sentence: []jplaw.Sentence{createTestSentence("Cell 1")}},
								},
							},
						},
					},
				},
			},
			contains: []string{"Test Table", "Cell 1", "table"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imgProc := &ImageProcessor{}
			result := processTableStructs(tt.tables, imgProc)

			for _, expected := range tt.contains {
				if expected != "" && !strings.Contains(result, expected) {
					t.Errorf("Expected result to contain %q, but it didn't", expected)
				}
			}
		})
	}
}

func TestProcessTable(t *testing.T) {
	tests := []struct {
		name     string
		table    *jplaw.Table
		contains []string
	}{
		{
			name: "Simple table",
			table: &jplaw.Table{
				TableRow: []jplaw.TableRow{
					{
						TableColumn: []jplaw.TableColumn{
							{Sentence: []jplaw.Sentence{createTestSentence("Cell")}},
						},
					},
				},
			},
			contains: []string{"<table", "<tbody", "Cell"},
		},
		{
			name: "Table with header row",
			table: &jplaw.Table{
				TableHeaderRow: []jplaw.TableHeaderRow{
					{
						TableHeaderColumn: []jplaw.TableHeaderColumn{
							{Content: "Header"},
						},
					},
				},
				TableRow: []jplaw.TableRow{
					{
						TableColumn: []jplaw.TableColumn{
							{Sentence: []jplaw.Sentence{createTestSentence("Cell")}},
						},
					},
				},
			},
			contains: []string{"<thead", "Header", "Cell"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processTable(tt.table)

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected result to contain %q, but it didn't", expected)
				}
			}
		})
	}
}

func TestProcessTableHeaderRow(t *testing.T) {
	row := &jplaw.TableHeaderRow{
		TableHeaderColumn: []jplaw.TableHeaderColumn{
			{Content: "Header 1"},
			{Content: "Header 2"},
		},
	}

	result := processTableHeaderRow(row)

	expected := []string{"<tr>", "<th", "Header 1", "Header 2"}
	for _, exp := range expected {
		if !strings.Contains(result, exp) {
			t.Errorf("Expected result to contain %q", exp)
		}
	}
}

func TestProcessTableRow(t *testing.T) {
	row := &jplaw.TableRow{
		TableColumn: []jplaw.TableColumn{
			{Sentence: []jplaw.Sentence{createTestSentence("Cell 1")}},
			{Sentence: []jplaw.Sentence{createTestSentence("Cell 2")}},
		},
	}

	result := processTableRow(row)

	expected := []string{"<tr>", "<td", "Cell 1", "Cell 2"}
	for _, exp := range expected {
		if !strings.Contains(result, exp) {
			t.Errorf("Expected result to contain %q", exp)
		}
	}
}

func TestProcessTableHeaderColumn(t *testing.T) {
	col := &jplaw.TableHeaderColumn{Content: "Header"}
	result := processTableHeaderColumn(col)

	if !strings.Contains(result, "<th") || !strings.Contains(result, "Header") {
		t.Errorf("Expected result to contain <th and Header, got %s", result)
	}
}

func TestProcessTableColumn(t *testing.T) {
	tests := []struct {
		name     string
		col      *jplaw.TableColumn
		contains []string
	}{
		{
			name: "Column with sentence",
			col: &jplaw.TableColumn{
				Sentence: []jplaw.Sentence{createTestSentence("Cell content")},
			},
			contains: []string{"<td", "Cell content"},
		},
		{
			name: "Column with borders",
			col: &jplaw.TableColumn{
				Sentence:     []jplaw.Sentence{createTestSentence("Cell")},
				BorderTop:    "solid",
				BorderBottom: "solid",
			},
			contains: []string{"Cell"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processTableColumn(tt.col)

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected result to contain %q, got %s", expected, result)
				}
			}
		})
	}
}

func TestProcessColumnElement(t *testing.T) {
	col := &jplaw.Column{
		Sentence: []jplaw.Sentence{createTestSentence("Column content")},
	}

	result := processColumnElement(col)

	if !strings.Contains(result, "Column content") {
		t.Errorf("Expected result to contain 'Column content', got %s", result)
	}
}

func TestProcessPartElement(t *testing.T) {
	part := &jplaw.Part{
		PartTitle: jplaw.PartTitle{Content: "Part Title"},
	}

	result := processPartElement(part)

	if !strings.Contains(result, "Part Title") {
		t.Errorf("Expected result to contain 'Part Title', got %s", result)
	}
}

func TestParseSpan(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantNil bool
		wantVal int
	}{
		{
			name:    "Valid number",
			input:   "2",
			wantNil: false,
			wantVal: 2,
		},
		{
			name:    "Invalid string",
			input:   "invalid",
			wantNil: true,
		},
		{
			name:    "Empty string",
			input:   "",
			wantNil: true,
		},
		{
			name:    "Zero",
			input:   "0",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseSpan(tt.input)
			if tt.wantNil {
				if got != nil {
					t.Errorf("parseSpan() = %v, want nil", *got)
				}
			} else {
				if got == nil {
					t.Error("parseSpan() = nil, want non-nil")
				} else if *got != tt.wantVal {
					t.Errorf("parseSpan() = %v, want %v", *got, tt.wantVal)
				}
			}
		})
	}
}

func TestBuildCellStyle(t *testing.T) {
	tests := []struct {
		name   string
		top    string
		bottom string
		left   string
		right  string
		want   string
	}{
		{
			name: "No borders",
			want: "",
		},
		{
			name:   "All borders solid",
			top:    "solid",
			bottom: "solid",
			left:   "solid",
			right:  "solid",
			want:   "border",
		},
		{
			name:   "Mixed border styles",
			top:    "dotted",
			bottom: "none",
			want:   "border-top",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildCellStyle(tt.top, tt.bottom, tt.left, tt.right)
			if tt.want != "" && !strings.Contains(got, tt.want) {
				t.Errorf("buildCellStyle() = %v, expected to contain %v", got, tt.want)
			}
			if tt.want == "" && got != "" {
				t.Errorf("buildCellStyle() = %v, expected empty", got)
			}
		})
	}
}

func TestBuildCellAttributes(t *testing.T) {
	tests := []struct {
		name    string
		rowspan *int
		colspan *int
		align   string
		valign  string
		style   string
		want    []string
	}{
		{
			name: "No attributes",
			want: []string{},
		},
		{
			name:    "With rowspan",
			rowspan: intPtr(2),
			want:    []string{"rowspan=\"2\""},
		},
		{
			name:    "With colspan",
			colspan: intPtr(3),
			want:    []string{"colspan=\"3\""},
		},
		{
			name:  "With alignment",
			align: "center",
			want:  []string{"align=\"center\""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildCellAttributes(tt.rowspan, tt.colspan, tt.align, tt.valign, tt.style)
			for _, expected := range tt.want {
				if !strings.Contains(got, expected) {
					t.Errorf("buildCellAttributes() = %v, expected to contain %v", got, expected)
				}
			}
		})
	}
}

// Helper function for creating int pointers
func intPtr(i int) *int {
	return &i
}
