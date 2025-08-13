package jplaw2epub

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.ngs.io/jplaw-xml"
)

func TestCreateEPUBFromXMLFile(t *testing.T) {
	// Create a simple test XML
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<Law Era="Reiwa" Year="1" Num="1" LawType="Act" Lang="ja">
  <LawNum>令和元年法律第一号</LawNum>
  <LawBody>
    <LawTitle>テスト法</LawTitle>
    <MainProvision>
      <Chapter Num="1">
        <ChapterTitle>第一章　総則</ChapterTitle>
        <Article Num="1">
          <ArticleTitle>第一条</ArticleTitle>
          <Paragraph Num="1">
            <ParagraphSentence>
              <Sentence>これはテストです。</Sentence>
            </ParagraphSentence>
          </Paragraph>
        </Article>
      </Chapter>
    </MainProvision>
  </LawBody>
</Law>`

	// Create temp file
	tmpfile, err := os.CreateTemp("", "test_*.xml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(xmlContent)); err != nil {
		t.Fatalf("Failed to write test XML: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Open and test
	xmlFile, err := os.Open(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	defer xmlFile.Close()

	book, err := CreateEPUBFromXMLFile(xmlFile)
	if err != nil {
		t.Errorf("CreateEPUBFromXMLFile() error = %v", err)
	}
	if book == nil {
		t.Error("CreateEPUBFromXMLFile() returned nil book")
	}
}

func TestCreateEPUBFromXMLFileWithOptions(t *testing.T) {
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<Law Era="Reiwa" Year="1" Num="1" LawType="Act" Lang="ja">
  <LawNum>令和元年法律第一号</LawNum>
  <LawBody>
    <LawTitle>テスト法</LawTitle>
    <MainProvision>
      <Paragraph Num="1">
        <ParagraphSentence>
          <Sentence>これは直接のParagraphです。</Sentence>
        </ParagraphSentence>
      </Paragraph>
      <Article Num="1">
        <ArticleTitle>第一条</ArticleTitle>
        <Paragraph Num="1">
          <ParagraphSentence>
            <Sentence>これは条文です。</Sentence>
          </ParagraphSentence>
        </Paragraph>
      </Article>
    </MainProvision>
  </LawBody>
</Law>`

	// Create temp file
	tmpfile, err := os.CreateTemp("", "test_*.xml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(xmlContent)); err != nil {
		t.Fatalf("Failed to write test XML: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	xmlFile, err := os.Open(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	defer xmlFile.Close()

	opts := &EPUBOptions{
		MaxImageHeight: "100px",
	}

	book, err := CreateEPUBFromXMLFileWithOptions(xmlFile, opts)
	if err != nil {
		t.Errorf("CreateEPUBFromXMLFileWithOptions() error = %v", err)
	}
	if book == nil {
		t.Error("CreateEPUBFromXMLFileWithOptions() returned nil book")
	}
}

func TestCreateEPUBFromXMLPath(t *testing.T) {
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<Law Era="Reiwa" Year="1" Num="1" LawType="Act" Lang="ja">
  <LawNum>令和元年法律第一号</LawNum>
  <LawBody>
    <LawTitle>テスト法</LawTitle>
    <MainProvision>
      <Chapter Num="1">
        <ChapterTitle>第一章</ChapterTitle>
      </Chapter>
    </MainProvision>
  </LawBody>
</Law>`

	// Create temp file
	tmpfile, err := os.CreateTemp("", "test_*.xml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(xmlContent)); err != nil {
		t.Fatalf("Failed to write test XML: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	book, err := CreateEPUBFromXMLPath(tmpfile.Name())
	if err != nil {
		t.Errorf("CreateEPUBFromXMLPath() error = %v", err)
	}
	if book == nil {
		t.Error("CreateEPUBFromXMLPath() returned nil book")
	}
}

func TestWriteEPUB(t *testing.T) {
	// Create a simple EPUB
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<Law Era="Reiwa" Year="1" Num="1" LawType="Act" Lang="ja">
  <LawNum>令和元年法律第一号</LawNum>
  <LawBody>
    <LawTitle>テスト法</LawTitle>
    <MainProvision>
      <Chapter Num="1">
        <ChapterTitle>第一章</ChapterTitle>
      </Chapter>
    </MainProvision>
  </LawBody>
</Law>`

	// Create temp file for XML
	tmpfile, err := os.CreateTemp("", "test_*.xml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(xmlContent)); err != nil {
		t.Fatalf("Failed to write test XML: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	book, err := CreateEPUBFromXMLPath(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to create EPUB: %v", err)
	}

	// Create temp file for EPUB output
	tmpDir := t.TempDir()
	epubPath := filepath.Join(tmpDir, "test.epub")

	err = WriteEPUB(book, epubPath)
	if err != nil {
		t.Errorf("WriteEPUB() error = %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(epubPath); os.IsNotExist(err) {
		t.Error("WriteEPUB() did not create the file")
	}
}

func TestProcessChaptersWithOptions(t *testing.T) {
	tests := []struct {
		name    string
		data    *jplaw.Law
		opts    *EPUBOptions
		wantErr bool
	}{
		{
			name: "Law with chapters",
			data: &jplaw.Law{
				LawBody: jplaw.LawBody{
					LawTitle: &jplaw.LawTitle{Content: "テスト法"},
					MainProvision: jplaw.MainProvision{
						Chapter: []jplaw.Chapter{
							{
								ChapterTitle: jplaw.ChapterTitle{Content: "第一章"},
								Article: []jplaw.Article{
									{
										ArticleTitle: &jplaw.ArticleTitle{Content: "第一条"},
										Paragraph: []jplaw.Paragraph{
											{
												ParagraphSentence: jplaw.ParagraphSentence{
													Sentence: []jplaw.Sentence{
														createTestSentence("内容"),
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			opts:    nil,
			wantErr: false,
		},
		{
			name: "Law with direct paragraphs",
			data: &jplaw.Law{
				LawBody: jplaw.LawBody{
					LawTitle: &jplaw.LawTitle{Content: "テスト法"},
					MainProvision: jplaw.MainProvision{
						Paragraph: []jplaw.Paragraph{
							{
								ParagraphSentence: jplaw.ParagraphSentence{
									Sentence: []jplaw.Sentence{
										createTestSentence("直接段落"),
									},
								},
							},
						},
					},
				},
			},
			opts:    nil,
			wantErr: false,
		},
		{
			name: "Law with AppdxNote",
			data: &jplaw.Law{
				LawBody: jplaw.LawBody{
					LawTitle: &jplaw.LawTitle{Content: "テスト法"},
					MainProvision: jplaw.MainProvision{
						Paragraph: []jplaw.Paragraph{
							{
								ParagraphSentence: jplaw.ParagraphSentence{
									Sentence: []jplaw.Sentence{
										createTestSentence("本文"),
									},
								},
							},
						},
					},
					AppdxNote: []jplaw.AppdxNote{
						{
							AppdxNoteTitle: &jplaw.AppdxNoteTitle{
								Content: "附則",
							},
						},
					},
				},
			},
			opts:    nil,
			wantErr: false,
		},
		{
			name: "Law with AppdxTable",
			data: &jplaw.Law{
				LawBody: jplaw.LawBody{
					LawTitle: &jplaw.LawTitle{Content: "テスト法"},
					MainProvision: jplaw.MainProvision{
						Paragraph: []jplaw.Paragraph{
							{
								ParagraphSentence: jplaw.ParagraphSentence{
									Sentence: []jplaw.Sentence{
										createTestSentence("本文"),
									},
								},
							},
						},
					},
					AppdxTable: []jplaw.AppdxTable{
						{
							AppdxTableTitle: &jplaw.AppdxTableTitle{
								Content: "附表",
							},
						},
					},
				},
			},
			opts:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book, err := createEPUBFromData(tt.data)
			if err != nil {
				t.Fatalf("Failed to create EPUB: %v", err)
			}

			err = processChaptersWithOptions(book, tt.data, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("processChaptersWithOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBuildChapterBody(t *testing.T) {
	tests := []struct {
		name     string
		chapter  *jplaw.Chapter
		contains []string
	}{
		{
			name: "Chapter with title only",
			chapter: &jplaw.Chapter{
				ChapterTitle: jplaw.ChapterTitle{
					Content: "第一章　総則",
				},
			},
			contains: []string{
				"第一章　総則",
				"chapter-title",
			},
		},
		{
			name: "Chapter with sections",
			chapter: &jplaw.Chapter{
				ChapterTitle: jplaw.ChapterTitle{
					Content: "第一章",
				},
				Section: []jplaw.Section{
					{
						SectionTitle: jplaw.SectionTitle{
							Content: "第一節",
						},
						Article: []jplaw.Article{
							{
								ArticleTitle: &jplaw.ArticleTitle{
									Content: "第一条",
								},
							},
						},
					},
				},
			},
			contains: []string{
				"第一章",
				"第一節",
				"第一条",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildChapterBody(tt.chapter)
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("buildChapterBody() should contain %q\ngot: %v", expected, result)
				}
			}
		})
	}
}

func TestBuildSectionsHTML(t *testing.T) {
	sections := []jplaw.Section{
		{
			SectionTitle: jplaw.SectionTitle{
				Content: "第一節　総則",
			},
			Article: []jplaw.Article{
				{
					ArticleTitle: &jplaw.ArticleTitle{
						Content: "第一条",
					},
				},
				{
					ArticleTitle: &jplaw.ArticleTitle{
						Content: "第二条",
					},
				},
			},
		},
		{
			SectionTitle: jplaw.SectionTitle{
				Content: "第二節　定義",
			},
			Article: []jplaw.Article{
				{
					ArticleTitle: &jplaw.ArticleTitle{
						Content: "第三条",
					},
				},
			},
		},
	}

	result := buildSectionsHTML(sections)

	expected := []string{
		"第一節　総則",
		"第二節　定義",
		"第一条 から 第二条 まで",
		"第三条 から 第三条 まで",
		"sections",
	}

	for _, exp := range expected {
		if !strings.Contains(result, exp) {
			t.Errorf("buildSectionsHTML() should contain %q\ngot: %v", exp, result)
		}
	}
}