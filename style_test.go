package jplaw2epub

import (
	"strings"
	"testing"

	"go.ngs.io/jplaw-xml"
)

func TestNewStyleProcessor(t *testing.T) {
	proc := NewStyleProcessor(nil)
	if proc == nil {
		t.Error("NewStyleProcessor returned nil")
	}
}

func TestProcessStyleStruct(t *testing.T) {
	proc := NewStyleProcessor(nil)

	tests := []struct {
		name     string
		style    *jplaw.StyleStruct
		contains []string
	}{
		{
			name: "Style with title",
			style: &jplaw.StyleStruct{
				StyleStructTitle: &jplaw.StyleStructTitle{
					Content: "スタイルタイトル",
				},
			},
			contains: []string{
				"style-struct",
				"style-title",
				"スタイルタイトル",
			},
		},
		{
			name: "Style with content",
			style: &jplaw.StyleStruct{
				Style: jplaw.Style{
					Content: "文章1 文章2",
				},
			},
			contains: []string{
				"文章1",
				"文章2",
			},
		},
		{
			name: "Style with remarks",
			style: &jplaw.StyleStruct{
				Remarks: []jplaw.Remarks{
					{
						Sentence: []jplaw.Sentence{
							createTestSentence("項目内容"),
						},
					},
				},
			},
			contains: []string{
				"項目内容",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := proc.ProcessStyleStruct(tt.style)
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("ProcessStyleStruct() should contain %q\ngot: %v", expected, result)
				}
			}
		})
	}
}

func TestProcessStyleStructs(t *testing.T) {
	styles := []jplaw.StyleStruct{
		{
			StyleStructTitle: &jplaw.StyleStructTitle{
				Content: "スタイル1",
			},
		},
		{
			StyleStructTitle: &jplaw.StyleStructTitle{
				Content: "スタイル2",
			},
		},
	}

	result := ProcessStyleStructs(styles, nil)

	if !strings.Contains(result, "スタイル1") {
		t.Error("ProcessStyleStructs should contain first style")
	}
	if !strings.Contains(result, "スタイル2") {
		t.Error("ProcessStyleStructs should contain second style")
	}
}
