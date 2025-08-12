package jplaw2epub

import (
	"testing"

	"go.ngs.io/jplaw-xml"
)

func TestOpenListWithStyle(t *testing.T) {
	tests := []struct {
		name   string
		titles []string
		want   string
	}{
		{
			name:   "Empty titles",
			titles: []string{},
			want:   htmlOL,
		},
		{
			name:   "CJK style",
			titles: []string{"一", "二", "三"},
			want:   `<ol style="list-style-type: cjk-ideographic;">`,
		},
		{
			name:   "Katakana style",
			titles: []string{"イ", "ロ", "ハ"},
			want:   `<ol style="list-style-type: katakana-iroha;">`,
		},
		{
			name:   "Decimal style",
			titles: []string{"１", "２", "３"},
			want:   `<ol style="list-style-type: decimal;">`,
		},
		{
			name:   "Default disc style",
			titles: []string{"A", "B", "C"},
			want:   htmlOL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := openListWithStyle(tt.titles)
			if got != tt.want {
				t.Errorf("openListWithStyle(%v) = %v, want %v", tt.titles, got, tt.want)
			}
		})
	}
}

func TestBuildArticleTitle(t *testing.T) {
	tests := []struct {
		name    string
		article *jplaw.Article
		want    string
	}{
		{
			name: "Article with title only",
			article: &jplaw.Article{
				ArticleTitle: &jplaw.ArticleTitle{
					Content: "第一条",
					Ruby:    []jplaw.Ruby{},
				},
				ArticleCaption: nil,
			},
			want: "第一条",
		},
		{
			name: "Article with title and caption",
			article: &jplaw.Article{
				ArticleTitle: &jplaw.ArticleTitle{
					Content: "第二条",
					Ruby:    []jplaw.Ruby{},
				},
				ArticleCaption: &jplaw.ArticleCaption{
					Content: "（定義）",
					Ruby:    []jplaw.Ruby{},
				},
			},
			want: "第二条 （定義）",
		},
		{
			name: "Article with ruby in title",
			article: &jplaw.Article{
				ArticleTitle: &jplaw.ArticleTitle{
					Content: "第三条",
					Ruby: []jplaw.Ruby{
						{
							Content: "較",
							Rt:      []jplaw.Rt{{Content: "こう"}},
						},
					},
				},
				ArticleCaption: nil,
			},
			want: "第三条<ruby>較<rt>こう</rt></ruby>",
		},
		{
			name: "Article with ruby in both title and caption",
			article: &jplaw.Article{
				ArticleTitle: &jplaw.ArticleTitle{
					Content: "第四条",
					Ruby: []jplaw.Ruby{
						{
							Content: "較",
							Rt:      []jplaw.Rt{{Content: "こう"}},
						},
					},
				},
				ArticleCaption: &jplaw.ArticleCaption{
					Content: "（校正）",
					Ruby: []jplaw.Ruby{
						{
							Content: "正",
							Rt:      []jplaw.Rt{{Content: "せい"}},
						},
					},
				},
			},
			want: "第四条<ruby>較<rt>こう</rt></ruby> （校正）<ruby>正<rt>せい</rt></ruby>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildArticleTitle(tt.article)
			if got != tt.want {
				t.Errorf("buildArticleTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTMLConstants(t *testing.T) {
	// Test that HTML constants are defined correctly
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"htmlOL", htmlOL, "<ol>"},
		{"htmlOLEnd", htmlOLEnd, "</ol>"},
		{"htmlLI", htmlLI, "<li>"},
		{"htmlLIEnd", htmlLIEnd, "</li>"},
		{"listStyleDisc", listStyleDisc, "disc"},
		{"listStyleDecimal", listStyleDecimal, "decimal"},
		{"listStyleCJK", listStyleCJK, "cjk-ideographic"},
		{"listStyleKatakana", listStyleKatakana, "katakana-iroha"},
		{"listStyleHiragana", listStyleHiragana, "hiragana-iroha"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.constant, tt.expected)
			}
		})
	}
}
