package main

import (
	"testing"

	"go.ngs.io/jplaw-xml"
)

func TestIsListNumber(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{"CJK number 一", "一", true},
		{"CJK number 十", "十", true},
		{"Katakana イ", "イ", true},
		{"Katakana ロ", "ロ", true},
		{"Full-width 1", "１", true},
		{"Full-width 10", "１０", true},
		{"Not a list number", "第一条", false},
		{"Empty string", "", false},
		{"Regular text", "これはテキストです", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isListNumber(tt.text); got != tt.want {
				t.Errorf("isListNumber(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}

func TestGetListStyleType(t *testing.T) {
	tests := []struct {
		name   string
		titles []string
		want   string
	}{
		{"Empty titles", []string{}, ""},
		{"CJK ideographic", []string{"一", "二", "三"}, listStyleCJK},
		{"Katakana iroha", []string{"イ", "ロ", "ハ"}, listStyleKatakana},
		{"Hiragana iroha", []string{"い", "ろ", "は"}, listStyleHiragana},
		{"Full-width numbers", []string{"１", "２", "３"}, listStyleDecimal},
		{"Half-width numbers", []string{"1", "2", "3"}, listStyleDecimal},
		{"Parenthesized", []string{"（１）", "（２）"}, listStyleDecimal},
		{"Unknown pattern", []string{"A", "B", "C"}, listStyleDisc},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getListStyleType(tt.titles); got != tt.want {
				t.Errorf("getListStyleType(%v) = %v, want %v", tt.titles, got, tt.want)
			}
		})
	}
}

func TestGetEraString(t *testing.T) {
	tests := []struct {
		name string
		era  jplaw.Era
		want string
	}{
		{"Meiji", jplaw.EraMeiji, "明治"},
		{"Taisho", jplaw.EraTaisho, "大正"},
		{"Showa", jplaw.EraShowa, "昭和"},
		{"Heisei", jplaw.EraHeisei, "平成"},
		{"Reiwa", jplaw.EraReiwa, "令和"},
		{"Unknown", jplaw.Era("Unknown"), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEraString(tt.era); got != tt.want {
				t.Errorf("getEraString(%v) = %v, want %v", tt.era, got, tt.want)
			}
		})
	}
}
