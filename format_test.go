package jplaw2epub

import (
	"strings"
	"testing"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

func TestProcessAppdxFormats(t *testing.T) {
	tests := []struct {
		name    string
		formats []jplaw.AppdxFormat
		wantErr bool
	}{
		{
			name:    "empty formats",
			formats: []jplaw.AppdxFormat{},
			wantErr: false,
		},
		{
			name: "single format",
			formats: []jplaw.AppdxFormat{
				{
					AppdxFormatTitle: &jplaw.AppdxFormatTitle{
						Content: "様式第一",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple formats",
			formats: []jplaw.AppdxFormat{
				{
					AppdxFormatTitle: &jplaw.AppdxFormatTitle{
						Content: "様式第一",
					},
				},
				{
					AppdxFormatTitle: &jplaw.AppdxFormatTitle{
						Content: "様式第二",
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
				t.Fatalf("Failed to create EPUB: %v", err)
			}

			err = processAppdxFormats(book, tt.formats, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("processAppdxFormats() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProcessAppdxFormat(t *testing.T) {
	tests := []struct {
		name         string
		format       *jplaw.AppdxFormat
		idx          int
		wantContains []string
		wantErr      bool
	}{
		{
			name: "format with title and related article",
			format: &jplaw.AppdxFormat{
				AppdxFormatTitle: &jplaw.AppdxFormatTitle{
					Content: "申請書様式",
				},
				RelatedArticleNum: &jplaw.RelatedArticleNum{
					Content: "（第十条関係）",
				},
				FormatStruct: []jplaw.FormatStruct{
					{
						FormatStructTitle: &jplaw.FormatStructTitle{
							Content: "申請書",
						},
					},
				},
			},
			idx: 0,
			wantContains: []string{
				"申請書様式",
				"（第十条関係）",
			},
			wantErr: false,
		},
		{
			name: "format with ruby text",
			format: &jplaw.AppdxFormat{
				AppdxFormatTitle: &jplaw.AppdxFormatTitle{
					Content: "様式",
					Ruby: []jplaw.Ruby{
						{
							Content: "様式",
							Rt:      []jplaw.Rt{{Content: "ようしき"}},
						},
					},
				},
			},
			idx:     1,
			wantErr: false,
		},
		{
			name: "format without title",
			format: &jplaw.AppdxFormat{
				FormatStruct: []jplaw.FormatStruct{
					{
						FormatStructTitle: &jplaw.FormatStructTitle{
							Content: "記載例",
						},
					},
				},
			},
			idx:     2,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book, err := epub.NewEpub("Test Book")
			if err != nil {
				t.Fatalf("Failed to create EPUB: %v", err)
			}

			err = processAppdxFormat(book, tt.format, tt.idx, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("processAppdxFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProcessFormatStruct(t *testing.T) {
	tests := []struct {
		name         string
		formatStruct *jplaw.FormatStruct
		wantContains []string
	}{
		{
			name: "with title",
			formatStruct: &jplaw.FormatStruct{
				FormatStructTitle: &jplaw.FormatStructTitle{
					Content: "書式タイトル",
				},
			},
			wantContains: []string{
				`<h3>書式タイトル</h3>`,
			},
		},
		{
			name: "with format content",
			formatStruct: &jplaw.FormatStruct{
				Format: jplaw.Format{
					Content: "申請書の内容",
				},
			},
			wantContains: []string{
				"申請書の内容",
			},
		},
		{
			name: "with ruby in title",
			formatStruct: &jplaw.FormatStruct{
				FormatStructTitle: &jplaw.FormatStructTitle{
					Content: "申請",
					Ruby: []jplaw.Ruby{
						{
							Content: "申請",
							Rt:      []jplaw.Rt{{Content: "しんせい"}},
						},
					},
				},
			},
			wantContains: []string{
				"<h3>",
				"申請",
				"しんせい",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processFormatStruct(tt.formatStruct, nil)

			for _, want := range tt.wantContains {
				if !strings.Contains(result, want) {
					t.Errorf("processFormatStruct() missing expected content: %s", want)
				}
			}
		})
	}
}

func TestProcessFormat(t *testing.T) {
	tests := []struct {
		name         string
		format       *jplaw.Format
		wantContains []string
	}{
		{
			name: "simple format",
			format: &jplaw.Format{
				Content: "書式の内容がここに入ります",
			},
			wantContains: []string{
				`<div class="format-content">`,
				`<pre class="format-raw">`,
				"書式の内容がここに入ります",
				`</div>`,
			},
		},
		{
			name: "format with Fig element",
			format: &jplaw.Format{
				Content: "内容 <Fig src=\"image.jpg\"/> テキスト",
			},
			wantContains: []string{
				`<div class="format-content">`,
				`<pre class="format-raw">`,
				"内容 <Fig src=\"image.jpg\"/> テキスト",
				`</div>`,
			},
		},
		{
			name: "format with multiple lines",
			format: &jplaw.Format{
				Content: "第一行\n第二行\n第三行",
			},
			wantContains: []string{
				"第一行",
				"第二行",
				"第三行",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processFormat(tt.format, nil)

			for _, want := range tt.wantContains {
				if !strings.Contains(result, want) {
					t.Errorf("processFormat() missing expected content: %s\nGot: %s", want, result)
				}
			}
		})
	}
}

func TestContainsFigElement(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "starts with Fig element",
			content: "<Fig src=\"image.jpg\"/> more text",
			want:    true,
		},
		{
			name:    "starts with Fig with attributes",
			content: "<Fig src=\"test.png\" alt=\"description\"/>",
			want:    true,
		},
		{
			name:    "no Fig element",
			content: "just plain text without images",
			want:    false,
		},
		{
			name:    "Fig in middle of content",
			content: "text <Fig src=\"test.jpg\"/> more",
			want:    false, // Only checks start
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containsFigElement(tt.content); got != tt.want {
				t.Errorf("containsFigElement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessEmbeddedFigs(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "single Fig element",
			content: "前文 <Fig src=\"image1.jpg\"/> 後文",
			want:    `<div class="embedded-content">前文 <Fig src="image1.jpg"/> 後文</div>`,
		},
		{
			name:    "multiple Fig elements",
			content: "<Fig src=\"img1.jpg\"/> text <Fig src=\"img2.png\"/>",
			want:    `<div class="embedded-content"><Fig src="img1.jpg"/> text <Fig src="img2.png"/></div>`,
		},
		{
			name:    "Fig with various attributes",
			content: "text <Fig src=\"test.jpg\" alt=\"desc\" width=\"100\"/> end",
			want:    `<div class="embedded-content">text <Fig src="test.jpg" alt="desc" width="100"/> end</div>`,
		},
		{
			name:    "no Fig elements",
			content: "plain text without any figures",
			want:    `<div class="embedded-content">plain text without any figures</div>`,
		},
		{
			name:    "nested Fig-like text",
			content: "text about <Figure> but not <Fig src=\"real.jpg\"/> tag",
			want:    `<div class="embedded-content">text about <Figure> but not <Fig src="real.jpg"/> tag</div>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processEmbeddedFigs(tt.content, nil)
			if got != tt.want {
				t.Errorf("processEmbeddedFigs() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestProcessFormatWithMockImageProcessor(t *testing.T) {
	// Create a mock image processor
	mock := &MockImageProcessor{
		ProcessFigStructHTML: `<img src="mocked.jpg" alt="Mocked Image"/>`,
	}

	format := &jplaw.Format{
		Content: "書式内容",
	}

	result := processFormat(format, mock)

	if !strings.Contains(result, "書式内容") {
		t.Errorf("processFormat() should contain format content")
	}

	// Verify mock wasn't called since there are no FigStructs in Format
	if len(mock.ProcessFigStructCalls) != 0 {
		t.Errorf("ProcessFigStruct should not be called for regular format content")
	}
}

func TestProcessAppdxFormatIntegration(t *testing.T) {
	book, err := epub.NewEpub("Test Book")
	if err != nil {
		t.Fatalf("Failed to create EPUB: %v", err)
	}

	format := &jplaw.AppdxFormat{
		AppdxFormatTitle: &jplaw.AppdxFormatTitle{
			Content: "統合テスト様式",
		},
		RelatedArticleNum: &jplaw.RelatedArticleNum{
			Content: "（第一条関係）",
		},
		FormatStruct: []jplaw.FormatStruct{
			{
				FormatStructTitle: &jplaw.FormatStructTitle{
					Content: "記載例",
				},
				Format: jplaw.Format{
					Content: "氏名：＿＿＿＿＿＿＿\n住所：＿＿＿＿＿＿＿",
				},
			},
		},
		Remarks: &jplaw.Remarks{
			RemarksLabel: jplaw.RemarksLabel{
				Content: "備考",
			},
			Sentence: []jplaw.Sentence{
				{
					MixedContent: jplaw.MixedContent{
						Nodes: []jplaw.ContentNode{
							jplaw.TextNode{Text: "記入にあたっては黒インクを使用すること"},
						},
					},
				},
			},
		},
	}

	err = processAppdxFormat(book, format, 0, nil)
	if err != nil {
		t.Errorf("processAppdxFormat() unexpected error: %v", err)
	}
}

// Benchmark tests
func BenchmarkProcessFormat(b *testing.B) {
	format := &jplaw.Format{
		Content: "これは書式のベンチマークテストです。<Fig src=\"test.jpg\"/> 複数行の\n内容を\n含みます。",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processFormat(format, nil)
	}
}

func BenchmarkContainsFigElement(b *testing.B) {
	content := "Long text with <Fig src=\"image.jpg\"/> embedded figure element and more text afterwards"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = containsFigElement(content)
	}
}

func BenchmarkProcessEmbeddedFigs(b *testing.B) {
	content := "<Fig src=\"1.jpg\"/> text <Fig src=\"2.jpg\"/> more <Fig src=\"3.jpg\"/> end"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processEmbeddedFigs(content, nil)
	}
}
