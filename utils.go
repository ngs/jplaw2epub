package jplaw2epub

import (
	"strings"

	"go.ngs.io/jplaw-xml"
)

// isListNumber checks if the text is just a Japanese list number
func isListNumber(text string) bool {
	// List numbers that should be skipped
	listNumbers := []string{
		// CJK ideographic numbers
		"一", "二", "三", "四", "五", "六", "七", "八", "九", "十",
		"十一", "十二", "十三", "十四", "十五", "十六", "十七", "十八", "十九", "二十",
		// Katakana iroha
		"イ", "ロ", "ハ", "ニ", "ホ", "ヘ", "ト", "チ", "リ", "ヌ",
		"ル", "ヲ", "ワ", "カ", "ヨ", "タ", "レ", "ソ", "ツ", "ネ",
		// Full-width Arabic numerals
		"１", "２", "３", "４", "５", "６", "７", "８", "９", "１０",
		"１１", "１２", "１３", "１４", "１５", "１６", "１７", "１８", "１９", "２０",
	}

	for _, num := range listNumbers {
		if text == num {
			return true
		}
	}
	return false
}

// getListStyleType determines the CSS list-style-type based on the item titles
func getListStyleType(titles []string) string {
	if len(titles) == 0 {
		return ""
	}

	// Check first title to determine the pattern
	first := titles[0]

	// CJK ideographic (一, 二, 三...)
	cjkNumbers := []string{"一", "二", "三", "四", "五", "六", "七", "八", "九", "十"}
	for _, num := range cjkNumbers {
		if first == num {
			return listStyleCJK
		}
	}

	// Katakana iroha (イ, ロ, ハ...)
	katakanaIroha := []string{"イ", "ロ", "ハ", "ニ", "ホ", "ヘ", "ト", "チ", "リ", "ヌ"}
	for _, kana := range katakanaIroha {
		if first == kana {
			return listStyleKatakana
		}
	}

	// Hiragana iroha (い, ろ, は...)
	hiraganaIroha := []string{"い", "ろ", "は", "に", "ほ", "へ", "と", "ち", "り", "ぬ"}
	for _, kana := range hiraganaIroha {
		if first == kana {
			return listStyleHiragana
		}
	}

	// Full-width Arabic numerals (１, ２, ３...)
	fullWidthNumbers := []string{"１", "２", "３", "４", "５", "６", "７", "８", "９"}
	for _, num := range fullWidthNumbers {
		if first == num {
			return listStyleDecimal
		}
	}

	// Half-width Arabic numerals (1, 2, 3...)
	if strings.HasPrefix(first, "1") {
		return listStyleDecimal
	}

	// Parenthesized numbers (（１）, （２）...)
	if strings.HasPrefix(first, "（") && strings.HasSuffix(first, "）") {
		return listStyleDecimal
	}

	// Default
	return listStyleDisc
}

// getEraString converts Era enum to Japanese string
func getEraString(era jplaw.Era) string {
	switch era {
	case jplaw.EraMeiji:
		return "明治"
	case jplaw.EraTaisho:
		return "大正"
	case jplaw.EraShowa:
		return "昭和"
	case jplaw.EraHeisei:
		return "平成"
	case jplaw.EraReiwa:
		return "令和"
	default:
		return ""
	}
}
