package main

import (
	"strings"
	"testing"

	"go.ngs.io/jplaw-xml"
)

func TestProcessItems(t *testing.T) {
	tests := []struct {
		name  string
		items []jplaw.Item
		want  string
	}{
		{
			name:  "Empty items",
			items: []jplaw.Item{},
			want:  "",
		},
		{
			name: "Single item with title",
			items: []jplaw.Item{
				{
					ItemTitle: &jplaw.ItemTitle{Content: "項目1"},
					ItemSentence: jplaw.ItemSentence{
						Sentence: []jplaw.Sentence{
							createTestSentence("これはテスト文です。"),
						},
					},
				},
			},
			want: "<ol><li><strong>項目1</strong> これはテスト文です。</li></ol>",
		},
		{
			name: "Item with list number title (should be skipped)",
			items: []jplaw.Item{
				{
					ItemTitle: &jplaw.ItemTitle{Content: "一"},
					ItemSentence: jplaw.ItemSentence{
						Sentence: []jplaw.Sentence{
							createTestSentence("第一項目です。"),
						},
					},
				},
			},
			want: `<ol style="list-style-type: cjk-ideographic;"><li>第一項目です。</li></ol>`,
		},
		{
			name: "Multiple items",
			items: []jplaw.Item{
				{
					ItemTitle: &jplaw.ItemTitle{Content: "一"},
					ItemSentence: jplaw.ItemSentence{
						Sentence: []jplaw.Sentence{
							createTestSentence("第一項目。"),
						},
					},
				},
				{
					ItemTitle: &jplaw.ItemTitle{Content: "二"},
					ItemSentence: jplaw.ItemSentence{
						Sentence: []jplaw.Sentence{
							createTestSentence("第二項目。"),
						},
					},
				},
			},
			want: `<ol style="list-style-type: cjk-ideographic;"><li>第一項目。</li><li>第二項目。</li></ol>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processItems(tt.items)
			if got != tt.want {
				t.Errorf("processItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessItem(t *testing.T) {
	tests := []struct {
		name string
		item *jplaw.Item
		want string
	}{
		{
			name: "Item without title",
			item: &jplaw.Item{
				ItemSentence: jplaw.ItemSentence{
					Sentence: []jplaw.Sentence{
						createTestSentence("文章のみ。"),
					},
				},
			},
			want: "<li>文章のみ。</li>",
		},
		{
			name: "Item with subitems",
			item: &jplaw.Item{
				ItemTitle: &jplaw.ItemTitle{Content: "主項目"},
				ItemSentence: jplaw.ItemSentence{
					Sentence: []jplaw.Sentence{
						createTestSentence("主項目の文章。"),
					},
				},
				Subitem1: []jplaw.Subitem1{
					{
						Subitem1Title: &jplaw.Subitem1Title{Content: "イ"},
						Subitem1Sentence: jplaw.Subitem1Sentence{
							Sentence: []jplaw.Sentence{
								createTestSentence("サブ項目。"),
							},
						},
					},
				},
			},
			want: `<li><strong>主項目</strong> 主項目の文章。<ol style="list-style-type: katakana-iroha;"><li>サブ項目。</li></ol></li>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processItem(tt.item)
			if got != tt.want {
				t.Errorf("processItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessSubitem1(t *testing.T) {
	tests := []struct {
		name     string
		subitem  *jplaw.Subitem1
		want     string
		contains []string
	}{
		{
			name: "Subitem1 without title",
			subitem: &jplaw.Subitem1{
				Subitem1Sentence: jplaw.Subitem1Sentence{
					Sentence: []jplaw.Sentence{
						createTestSentence("サブ項目の文章。"),
					},
				},
			},
			want: "<li>サブ項目の文章。</li>",
		},
		{
			name: "Subitem1 with Subitem2",
			subitem: &jplaw.Subitem1{
				Subitem1Title: &jplaw.Subitem1Title{Content: "イ"},
				Subitem1Sentence: jplaw.Subitem1Sentence{
					Sentence: []jplaw.Sentence{
						createTestSentence("第一レベル。"),
					},
				},
				Subitem2: []jplaw.Subitem2{
					{
						Subitem2Title: &jplaw.Subitem2Title{Content: "（１）"},
						Subitem2Sentence: jplaw.Subitem2Sentence{
							Sentence: []jplaw.Sentence{
								createTestSentence("第二レベル。"),
							},
						},
					},
				},
			},
			contains: []string{
				"<li>",
				"第一レベル。",
				"<ol",
				"第二レベル。",
				"</ol>",
				"</li>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processSubitem1(tt.subitem)
			if tt.want != "" && got != tt.want {
				t.Errorf("processSubitem1() = %v, want %v", got, tt.want)
			}
			for _, contain := range tt.contains {
				if !strings.Contains(got, contain) {
					t.Errorf("processSubitem1() should contain %v, got %v", contain, got)
				}
			}
		})
	}
}

func TestProcessSubitem2(t *testing.T) {
	subitem := &jplaw.Subitem2{
		Subitem2Title: &jplaw.Subitem2Title{Content: "詳細項目"},
		Subitem2Sentence: jplaw.Subitem2Sentence{
			Sentence: []jplaw.Sentence{
				createTestSentence("詳細な内容。"),
			},
		},
	}

	got := processSubitem2(subitem)
	want := "<li><strong>詳細項目</strong> 詳細な内容。</li>"

	if got != want {
		t.Errorf("processSubitem2() = %v, want %v", got, want)
	}
}

func TestCollectItemTitles(t *testing.T) {
	items := []jplaw.Item{
		{ItemTitle: &jplaw.ItemTitle{Content: "一"}},
		{ItemTitle: &jplaw.ItemTitle{Content: "二"}},
		{ItemTitle: nil},
		{ItemTitle: &jplaw.ItemTitle{Content: "三"}},
	}

	got := collectItemTitles(items)
	want := []string{"一", "二", "三"}

	if len(got) != len(want) {
		t.Errorf("collectItemTitles() returned %d items, want %d", len(got), len(want))
		return
	}

	for i, title := range got {
		if title != want[i] {
			t.Errorf("collectItemTitles()[%d] = %v, want %v", i, title, want[i])
		}
	}
}

func TestCollectSubitem1Titles(t *testing.T) {
	subitems := []jplaw.Subitem1{
		{Subitem1Title: &jplaw.Subitem1Title{Content: "イ"}},
		{Subitem1Title: &jplaw.Subitem1Title{Content: "ロ"}},
		{Subitem1Title: nil},
	}

	got := collectSubitem1Titles(subitems)
	want := []string{"イ", "ロ"}

	if len(got) != len(want) {
		t.Errorf("collectSubitem1Titles() returned %d items, want %d", len(got), len(want))
		return
	}

	for i, title := range got {
		if title != want[i] {
			t.Errorf("collectSubitem1Titles()[%d] = %v, want %v", i, title, want[i])
		}
	}
}

func TestCollectSubitem2Titles(t *testing.T) {
	subitems := []jplaw.Subitem2{
		{Subitem2Title: &jplaw.Subitem2Title{Content: "（１）"}},
		{Subitem2Title: &jplaw.Subitem2Title{Content: "（２）"}},
	}

	got := collectSubitem2Titles(subitems)
	want := []string{"（１）", "（２）"}

	if len(got) != len(want) {
		t.Errorf("collectSubitem2Titles() returned %d items, want %d", len(got), len(want))
		return
	}

	for i, title := range got {
		if title != want[i] {
			t.Errorf("collectSubitem2Titles()[%d] = %v, want %v", i, title, want[i])
		}
	}
}
