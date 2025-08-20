package jplaw2epub

import (
	"strings"
	"testing"

	"go.ngs.io/jplaw-xml"
)

func TestProcessRubyElements(t *testing.T) {
	tests := []struct {
		name   string
		rubies []jplaw.Ruby
		want   string
	}{
		{
			name:   "Empty rubies",
			rubies: []jplaw.Ruby{},
			want:   "",
		},
		{
			name: "Single ruby with RT",
			rubies: []jplaw.Ruby{
				{
					Content: "較",
					Rt: []jplaw.Rt{
						{Content: "こう"},
					},
				},
			},
			want: "<ruby>較<rt>こう</rt></ruby>",
		},
		{
			name: "Multiple RT elements",
			rubies: []jplaw.Ruby{
				{
					Content: "振",
					Rt: []jplaw.Rt{
						{Content: "ふ"},
						{Content: "り"},
					},
				},
			},
			want: "<ruby>振<rt>ふ</rt><rt>り</rt></ruby>",
		},
		{
			name: "Ruby without RT",
			rubies: []jplaw.Ruby{
				{
					Content: "漢字",
					Rt:      []jplaw.Rt{},
				},
			},
			want: "漢字",
		},
		{
			name: "Multiple rubies",
			rubies: []jplaw.Ruby{
				{
					Content: "漢",
					Rt:      []jplaw.Rt{{Content: "かん"}},
				},
				{
					Content: "字",
					Rt:      []jplaw.Rt{{Content: "じ"}},
				},
			},
			want: "<ruby>漢<rt>かん</rt></ruby><ruby>字<rt>じ</rt></ruby>",
		},
		{
			name: "HTML escaping",
			rubies: []jplaw.Ruby{
				{
					Content: "<tag>",
					Rt:      []jplaw.Rt{{Content: "たぐ"}},
				},
			},
			want: "<ruby>&lt;tag&gt;<rt>たぐ</rt></ruby>",
		},
		{
			name: "Special characters in ruby content",
			rubies: []jplaw.Ruby{
				{
					Content: "検&査",
					Rt:      []jplaw.Rt{{Content: "けん&さ"}},
				},
			},
			want: "<ruby>検&amp;査<rt>けん&amp;さ</rt></ruby>",
		},
		{
			name: "Unicode characters in ruby",
			rubies: []jplaw.Ruby{
				{
					Content: "日本国憲法",
					Rt:      []jplaw.Rt{{Content: "にっぽんこくけんぽう"}},
				},
			},
			want: "<ruby>日本国憲法<rt>にっぽんこくけんぽう</rt></ruby>",
		},
		{
			name: "Empty RT content",
			rubies: []jplaw.Ruby{
				{
					Content: "漢字",
					Rt:      []jplaw.Rt{{Content: ""}},
				},
			},
			want: "<ruby>漢字<rt></rt></ruby>",
		},
		{
			name: "Nil RT array",
			rubies: []jplaw.Ruby{
				{
					Content: "漢字",
					Rt:      nil,
				},
			},
			want: "漢字",
		},
		{
			name: "Mixed ruby with and without RT",
			rubies: []jplaw.Ruby{
				{
					Content: "有",
					Rt:      []jplaw.Rt{{Content: "ゆう"}},
				},
				{
					Content: "無",
					Rt:      []jplaw.Rt{},
				},
				{
					Content: "効",
					Rt:      []jplaw.Rt{{Content: "こう"}},
				},
			},
			want: "<ruby>有<rt>ゆう</rt></ruby>無<ruby>効<rt>こう</rt></ruby>",
		},
		{
			name: "Very long ruby content",
			rubies: []jplaw.Ruby{
				{
					Content: "超長文字列超長文字列超長文字列超長文字列",
					Rt:      []jplaw.Rt{{Content: "ちょうちょうもじれつ"}},
				},
			},
			want: "<ruby>超長文字列超長文字列超長文字列超長文字列<rt>ちょうちょうもじれつ</rt></ruby>",
		},
		{
			name: "Quotation marks in content",
			rubies: []jplaw.Ruby{
				{
					Content: `"引用"`,
					Rt:      []jplaw.Rt{{Content: "いんよう"}},
				},
			},
			want: "<ruby>&#34;引用&#34;<rt>いんよう</rt></ruby>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processRubyElements(tt.rubies)
			if got != tt.want {
				t.Errorf("processRubyElements() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessTextWithRuby(t *testing.T) {
	tests := []struct {
		name    string
		content string
		rubies  []jplaw.Ruby
		want    string
	}{
		{
			name:    "No rubies",
			content: "普通のテキスト",
			rubies:  []jplaw.Ruby{},
			want:    "普通のテキスト",
		},
		{
			name:    "Empty content with rubies",
			content: "",
			rubies: []jplaw.Ruby{
				{
					Content: "較",
					Rt:      []jplaw.Rt{{Content: "こう"}},
				},
			},
			want: "<ruby>較<rt>こう</rt></ruby>",
		},
		{
			name:    "Content with rubies",
			content: "正又は校正",
			rubies: []jplaw.Ruby{
				{
					Content: "較",
					Rt:      []jplaw.Rt{{Content: "こう"}},
				},
			},
			want: "正又は校正<ruby>較<rt>こう</rt></ruby>",
		},
		{
			name:    "HTML escaping in content",
			content: "<div>テキスト</div>",
			rubies:  []jplaw.Ruby{},
			want:    "&lt;div&gt;テキスト&lt;/div&gt;",
		},
		{
			name:    "Content and rubies with escaping",
			content: "これは<test>",
			rubies: []jplaw.Ruby{
				{
					Content: "&amp;",
					Rt:      []jplaw.Rt{{Content: "あんど"}},
				},
			},
			want: "これは&lt;test&gt;<ruby>&amp;amp;<rt>あんど</rt></ruby>",
		},
		{
			name:    "Multiple rubies with content",
			content: "法令文書",
			rubies: []jplaw.Ruby{
				{
					Content: "法",
					Rt:      []jplaw.Rt{{Content: "ほう"}},
				},
				{
					Content: "令",
					Rt:      []jplaw.Rt{{Content: "れい"}},
				},
			},
			want: "法令文書<ruby>法<rt>ほう</rt></ruby><ruby>令<rt>れい</rt></ruby>",
		},
		{
			name:    "Content with special characters",
			content: "A&B会社の<規約>",
			rubies: []jplaw.Ruby{
				{
					Content: "規約",
					Rt:      []jplaw.Rt{{Content: "きやく"}},
				},
			},
			want: "A&amp;B会社の&lt;規約&gt;<ruby>規約<rt>きやく</rt></ruby>",
		},
		{
			name:    "Empty content and empty rubies",
			content: "",
			rubies:  []jplaw.Ruby{},
			want:    "",
		},
		{
			name:    "Nil rubies",
			content: "テキスト",
			rubies:  nil,
			want:    "テキスト",
		},
		{
			name:    "Content with line breaks",
			content: "第一行\n第二行",
			rubies: []jplaw.Ruby{
				{
					Content: "行",
					Rt:      []jplaw.Rt{{Content: "ぎょう"}},
				},
			},
			want: "第一行\n第二行<ruby>行<rt>ぎょう</rt></ruby>",
		},
		{
			name:    "Unicode content with rubies",
			content: "日本国憲法第九条",
			rubies: []jplaw.Ruby{
				{
					Content: "憲法",
					Rt:      []jplaw.Rt{{Content: "けんぽう"}},
				},
				{
					Content: "九条",
					Rt:      []jplaw.Rt{{Content: "きゅうじょう"}},
				},
			},
			want: "日本国憲法第九条<ruby>憲法<rt>けんぽう</rt></ruby><ruby>九条<rt>きゅうじょう</rt></ruby>",
		},
		{
			name:    "Tabs and spaces in content",
			content: "項目\t内容  説明",
			rubies: []jplaw.Ruby{
				{
					Content: "項目",
					Rt:      []jplaw.Rt{{Content: "こうもく"}},
				},
			},
			want: "項目\t内容  説明<ruby>項目<rt>こうもく</rt></ruby>",
		},
		{
			name:    "Ruby without RT in mixed rubies",
			content: "対象文書",
			rubies: []jplaw.Ruby{
				{
					Content: "対象",
					Rt:      []jplaw.Rt{{Content: "たいしょう"}},
				},
				{
					Content: "文書",
					Rt:      []jplaw.Rt{},
				},
			},
			want: "対象文書<ruby>対象<rt>たいしょう</rt></ruby>文書",
		},
		{
			name:    "Content with quotes",
			content: `「引用文」と"quotation"`,
			rubies: []jplaw.Ruby{
				{
					Content: "引用",
					Rt:      []jplaw.Rt{{Content: "いんよう"}},
				},
			},
			want: `「引用文」と&#34;quotation&#34;<ruby>引用<rt>いんよう</rt></ruby>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processTextWithRuby(tt.content, tt.rubies)
			if got != tt.want {
				t.Errorf("processTextWithRuby(%q, %v) = %v, want %v",
					tt.content, tt.rubies, got, tt.want)
			}
		})
	}
}

func TestProcessTextWithRuby_Integration(t *testing.T) {
	// Test that verifies the note about Ruby elements being appended
	content := "別表第二に掲げる測定器その他の設備であつて、次のいずれかに掲げる"
	rubies := []jplaw.Ruby{
		{
			Content: "較",
			Rt:      []jplaw.Rt{{Content: "こう"}},
		},
	}

	got := processTextWithRuby(content, rubies)

	// Verify content comes first
	if !strings.HasPrefix(got, content) {
		t.Errorf("Content should come first in result")
	}

	// Verify ruby is appended at the end
	if !strings.HasSuffix(got, "<ruby>較<rt>こう</rt></ruby>") {
		t.Errorf("Ruby should be appended at the end")
	}
}

func TestRubyHTML(t *testing.T) {
	tests := []struct {
		name string
		ruby *jplaw.Ruby
		want string
	}{
		{
			name: "Ruby with single RT",
			ruby: &jplaw.Ruby{
				Content: "法",
				Rt:      []jplaw.Rt{{Content: "ほう"}},
			},
			want: "<ruby>法<rt>ほう</rt></ruby>",
		},
		{
			name: "Ruby without RT",
			ruby: &jplaw.Ruby{
				Content: "法",
				Rt:      []jplaw.Rt{},
			},
			want: "法",
		},
		{
			name: "Ruby with nil RT",
			ruby: &jplaw.Ruby{
				Content: "法",
				Rt:      nil,
			},
			want: "法",
		},
		{
			name: "Ruby with multiple RT",
			ruby: &jplaw.Ruby{
				Content: "振仮名",
				Rt: []jplaw.Rt{
					{Content: "ふり"},
					{Content: "がな"},
				},
			},
			want: "<ruby>振仮名<rt>ふり</rt><rt>がな</rt></ruby>",
		},
		{
			name: "Ruby with HTML special characters",
			ruby: &jplaw.Ruby{
				Content: "<特殊>",
				Rt:      []jplaw.Rt{{Content: "とくしゅ"}},
			},
			want: "<ruby>&lt;特殊&gt;<rt>とくしゅ</rt></ruby>",
		},
		{
			name: "Ruby with ampersand",
			ruby: &jplaw.Ruby{
				Content: "A&B",
				Rt:      []jplaw.Rt{{Content: "えーあんどびー"}},
			},
			want: "<ruby>A&amp;B<rt>えーあんどびー</rt></ruby>",
		},
		{
			name: "Ruby with quotes",
			ruby: &jplaw.Ruby{
				Content: `"引用"`,
				Rt:      []jplaw.Rt{{Content: "いんよう"}},
			},
			want: "<ruby>&#34;引用&#34;<rt>いんよう</rt></ruby>",
		},
		{
			name: "Ruby with empty RT content",
			ruby: &jplaw.Ruby{
				Content: "漢字",
				Rt:      []jplaw.Rt{{Content: ""}},
			},
			want: "<ruby>漢字<rt></rt></ruby>",
		},
		{
			name: "Empty ruby content with RT",
			ruby: &jplaw.Ruby{
				Content: "",
				Rt:      []jplaw.Rt{{Content: "よみ"}},
			},
			want: "<ruby><rt>よみ</rt></ruby>",
		},
		{
			name: "Ruby with apostrophe",
			ruby: &jplaw.Ruby{
				Content: "it's",
				Rt:      []jplaw.Rt{{Content: "いっつ"}},
			},
			want: "<ruby>it&#39;s<rt>いっつ</rt></ruby>",
		},
		{
			name: "Ruby with line break in content",
			ruby: &jplaw.Ruby{
				Content: "改\n行",
				Rt:      []jplaw.Rt{{Content: "かいぎょう"}},
			},
			want: "<ruby>改\n行<rt>かいぎょう</rt></ruby>",
		},
		{
			name: "Ruby with tab character",
			ruby: &jplaw.Ruby{
				Content: "タブ\t文字",
				Rt:      []jplaw.Rt{{Content: "たぶもじ"}},
			},
			want: "<ruby>タブ\t文字<rt>たぶもじ</rt></ruby>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rubyHTML(tt.ruby)
			if got != tt.want {
				t.Errorf("rubyHTML() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Benchmark tests
func BenchmarkProcessRubyElements(b *testing.B) {
	rubies := []jplaw.Ruby{
		{Content: "法", Rt: []jplaw.Rt{{Content: "ほう"}}},
		{Content: "令", Rt: []jplaw.Rt{{Content: "れい"}}},
		{Content: "文", Rt: []jplaw.Rt{{Content: "ぶん"}}},
		{Content: "書", Rt: []jplaw.Rt{{Content: "しょ"}}},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processRubyElements(rubies)
	}
}

func BenchmarkProcessTextWithRuby(b *testing.B) {
	content := "これは法令文書のサンプルテキストです"
	rubies := []jplaw.Ruby{
		{Content: "法令", Rt: []jplaw.Rt{{Content: "ほうれい"}}},
		{Content: "文書", Rt: []jplaw.Rt{{Content: "ぶんしょ"}}},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processTextWithRuby(content, rubies)
	}
}

func BenchmarkRubyHTML(b *testing.B) {
	ruby := &jplaw.Ruby{
		Content: "法令",
		Rt:      []jplaw.Rt{{Content: "ほうれい"}},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rubyHTML(ruby)
	}
}

// Edge case tests
func TestProcessTextWithRuby_EdgeCases(t *testing.T) {
	t.Run("Very large content", func(t *testing.T) {
		// Generate a very large string
		largeContent := strings.Repeat("大きなテキスト", 1000)
		rubies := []jplaw.Ruby{
			{Content: "テキスト", Rt: []jplaw.Rt{{Content: "てきすと"}}},
		}

		got := processTextWithRuby(largeContent, rubies)

		// Should handle large content without panic
		if !strings.HasPrefix(got, largeContent) {
			t.Errorf("Large content should be preserved")
		}
		if !strings.HasSuffix(got, "<ruby>テキスト<rt>てきすと</rt></ruby>") {
			t.Errorf("Ruby should be appended at the end")
		}
	})

	t.Run("Many rubies", func(t *testing.T) {
		content := "内容"
		var rubies []jplaw.Ruby
		for i := 0; i < 100; i++ {
			rubies = append(rubies, jplaw.Ruby{
				Content: string(rune('あ' + i)),
				Rt:      []jplaw.Rt{{Content: "よみ"}},
			})
		}

		got := processTextWithRuby(content, rubies)

		// Should handle many rubies without panic
		if !strings.HasPrefix(got, content) {
			t.Errorf("Content should come first")
		}
		// Check that all rubies are appended
		if strings.Count(got, "<ruby>") != 100 {
			t.Errorf("All rubies should be appended")
		}
	})

	t.Run("Deeply nested HTML escaping", func(t *testing.T) {
		content := "<div><span>&lt;nested&gt;</span></div>"
		rubies := []jplaw.Ruby{
			{Content: "&lt;", Rt: []jplaw.Rt{{Content: "レスザン"}}},
		}

		got := processTextWithRuby(content, rubies)
		want := "&lt;div&gt;&lt;span&gt;&amp;lt;nested&amp;gt;&lt;/span&gt;&lt;/div&gt;<ruby>&amp;lt;<rt>レスザン</rt></ruby>"

		if got != want {
			t.Errorf("processTextWithRuby() = %v, want %v", got, want)
		}
	})
}

// Test for concurrent access (if the functions are used in concurrent contexts)
func TestProcessTextWithRuby_Concurrent(t *testing.T) {
	content := "並行処理テスト"
	rubies := []jplaw.Ruby{
		{Content: "並行", Rt: []jplaw.Rt{{Content: "へいこう"}}},
		{Content: "処理", Rt: []jplaw.Rt{{Content: "しょり"}}},
	}

	// Run multiple goroutines
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			result := processTextWithRuby(content, rubies)
			expected := "並行処理テスト<ruby>並行<rt>へいこう</rt></ruby><ruby>処理<rt>しょり</rt></ruby>"
			if result != expected {
				t.Errorf("Concurrent execution failed: got %v, want %v", result, expected)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}
