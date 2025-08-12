package jplaw2epub

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.ngs.io/jplaw-xml"
)

const (
	testLawTitle = "テスト法"
	testXMLContent = `<?xml version="1.0" encoding="UTF-8"?>
<Law Era="Reiwa" Year="2024" LawType="Act" Num="1" Lang="ja">
	<LawNum>令和六年法律第一号</LawNum>
	<LawBody>
		<LawTitle>テスト法</LawTitle>
		<MainProvision>
			<Chapter Num="1">
				<ChapterTitle>第一章　総則</ChapterTitle>
				<Article Num="1">
					<ArticleTitle>第一条</ArticleTitle>
					<Paragraph Num="1">
						<ParagraphNum/>
						<ParagraphSentence>
							<Sentence>これはテストです。</Sentence>
						</ParagraphSentence>
					</Paragraph>
				</Article>
			</Chapter>
		</MainProvision>
	</LawBody>
</Law>`
)

func TestLoadXMLDataFromReader(t *testing.T) {
	reader := strings.NewReader(testXMLContent)
	data, err := loadXMLDataFromReader(reader)
	if err != nil {
		t.Fatalf("loadXMLDataFromReader failed: %v", err)
	}

	if data == nil {
		t.Fatal("Expected non-nil law data")
	}

	if data.LawBody.LawTitle.Content != testLawTitle {
		t.Errorf("Expected law title '%s', got '%s'", testLawTitle, data.LawBody.LawTitle.Content)
	}
}

func TestCreateEPUBFromData(t *testing.T) {
	// Create test law data
	data := &jplaw.Law{
		LawNum: "令和六年法律第一号",
		Year:   2024,
		Era:    "Reiwa",
		Lang:   "ja",
		LawBody: jplaw.LawBody{
			LawTitle: &jplaw.LawTitle{
				Content: testLawTitle,
			},
		},
	}

	book, err := createEPUBFromData(data)
	if err != nil {
		t.Fatalf("createEPUBFromData failed: %v", err)
	}

	if book == nil {
		t.Fatal("Expected non-nil EPUB book")
	}

	if book.Title() != testLawTitle {
		t.Errorf("Expected title '%s', got '%s'", testLawTitle, book.Title())
	}
}

func TestSetupEPUBMetadata(t *testing.T) {
	data := &jplaw.Law{
		LawNum:          "令和六年法律第一号",
		Year:            2024,
		Era:             "Reiwa",
		Lang:            "ja",
		PromulgateMonth: 4,
		PromulgateDay:   1,
		LawBody: jplaw.LawBody{
			LawTitle: &jplaw.LawTitle{
				Content: testLawTitle,
				Kana:    "テストホウ",
			},
		},
	}

	book, err := createEPUBFromData(data)
	if err != nil {
		t.Fatalf("Failed to create EPUB: %v", err)
	}

	// Verify metadata was set
	if book.Author() != data.LawNum {
		t.Errorf("Expected author '%s', got '%s'", data.LawNum, book.Author())
	}

	if book.Lang() != string(data.Lang) {
		t.Errorf("Expected lang '%s', got '%s'", string(data.Lang), book.Lang())
	}
}

func TestIntegration_SmallLaw(t *testing.T) {
	// Create a temporary XML file
	tmpDir := t.TempDir()
	xmlFile := filepath.Join(tmpDir, "test.xml")

	if err := os.WriteFile(xmlFile, []byte(testXMLContent), 0o644); err != nil {
		t.Fatalf("Failed to create test XML file: %v", err)
	}

	// Test the library
	book, err := CreateEPUBFromXMLPath(xmlFile)
	if err != nil {
		t.Fatalf("CreateEPUBFromXMLPath failed: %v", err)
	}

	if book == nil {
		t.Fatal("Expected non-nil EPUB book")
	}

	// Write EPUB to temporary file
	epubPath := filepath.Join(tmpDir, "test.epub")
	if err := WriteEPUB(book, epubPath); err != nil {
		t.Fatalf("WriteEPUB failed: %v", err)
	}

	// Verify the file was created
	if _, err := os.Stat(epubPath); os.IsNotExist(err) {
		t.Fatal("EPUB file was not created")
	}
}
