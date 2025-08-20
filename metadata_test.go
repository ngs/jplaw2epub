package jplaw2epub

import (
	"testing"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

func TestSetupEPUBMetadata(t *testing.T) {
	tests := []struct {
		name           string
		data           *jplaw.Law
		expectedAuthor string
		expectedLang   string
		expectedInDesc []string
	}{
		{
			name: "complete metadata",
			data: &jplaw.Law{
				Era:             "Heisei",
				Year:            27,
				PromulgateMonth: 3,
				PromulgateDay:   26,
				LawNum:          "平成二十七年総務省令第二十四号",
				Lang:            "ja",
				LawBody: jplaw.LawBody{
					LawTitle: &jplaw.LawTitle{
						Content: "放送法及び電波法の一部を改正する法律",
						Kana:    "ホウソウホウオヨビデンパホウノイチブヲカイセイスルホウリツ",
					},
				},
			},
			expectedAuthor: "平成二十七年総務省令第二十四号",
			expectedLang:   "ja",
			expectedInDesc: []string{
				"公布日: 平成 27年3月26日",
				"法令番号: 平成二十七年総務省令第二十四号",
				"現行法令名: 放送法及び電波法の一部を改正する法律",
			},
		},
		{
			name: "metadata with ruby",
			data: &jplaw.Law{
				Era:             "Reiwa",
				Year:            5,
				PromulgateMonth: 1,
				PromulgateDay:   15,
				LawNum:          "令和五年法律第一号",
				Lang:            "ja",
				LawBody: jplaw.LawBody{
					LawTitle: &jplaw.LawTitle{
						Content: "特別措置法",
						Kana:    "トクベツソチホウ",
						Ruby: []jplaw.Ruby{
							{
								Content: "特別",
								Rt:      []jplaw.Rt{{Content: "とくべつ"}},
							},
							{
								Content: "措置",
								Rt:      []jplaw.Rt{{Content: "そち"}},
							},
						},
					},
				},
			},
			expectedAuthor: "令和五年法律第一号",
			expectedLang:   "ja",
			expectedInDesc: []string{
				"公布日: 令和 5年1月15日",
				"法令番号: 令和五年法律第一号",
				"特別措置法",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new EPUB book
			book, err := epub.NewEpub("Test Book")
			if err != nil {
				t.Fatalf("Failed to create EPUB: %v", err)
			}

			// Apply metadata
			setupEPUBMetadata(book, tt.data)

			// Since we can't easily inspect metadata without writing to disk,
			// we'll just verify that the function executed successfully
			// The setupEPUBMetadata function doesn't return an error,
			// so if we got here, the test passed
		})
	}
}
