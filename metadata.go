package jplaw2epub

import (
	"fmt"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

// setupEPUBMetadata sets up the basic EPUB metadata
func setupEPUBMetadata(book *epub.Epub, data *jplaw.Law) {
	book.SetAuthor(data.LawNum)
	book.SetLang(string(data.Lang))

	// Set description
	eraStr := getEraString(data.Era)
	description := fmt.Sprintf("公布日: %s %d年%d月%d日", eraStr, data.Year, data.PromulgateMonth, data.PromulgateDay)
	description += fmt.Sprintf("\n法令番号: %s", data.LawNum)
	lawTitleWithRuby := processTextWithRuby(data.LawBody.LawTitle.Content, data.LawBody.LawTitle.Ruby)
	description += fmt.Sprintf("\n現行法令名: %s %s", lawTitleWithRuby, data.LawBody.LawTitle.Kana)
	book.SetDescription(description)
}
