package main

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
