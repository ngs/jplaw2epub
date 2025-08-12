package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

func TestParseFlags(t *testing.T) {
	// Save original args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Note: This test would normally exit the program on error,
	// so we can only test the happy path in unit tests
	os.Args = []string{"cmd", "-d", "output.epub", "input.xml"}
	
	// We can't fully test parseFlags as it calls os.Exit on error
	// This would require refactoring the function to return errors instead
}

func TestLoadXMLData(t *testing.T) {
	// Create a temporary XML file
	tmpDir := t.TempDir()
	xmlFile := filepath.Join(tmpDir, "test.xml")
	
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
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

	if err := os.WriteFile(xmlFile, []byte(xmlContent), 0644); err != nil {
		t.Fatalf("Failed to create test XML file: %v", err)
	}

	// Test loading
	data := loadXMLData(xmlFile)
	
	if data == nil {
		t.Fatal("loadXMLData returned nil")
	}
	
	if data.LawNum != "令和六年法律第一号" {
		t.Errorf("LawNum = %v, want 令和六年法律第一号", data.LawNum)
	}
	
	if data.LawBody.LawTitle.Content != "テスト法" {
		t.Errorf("LawTitle = %v, want テスト法", data.LawBody.LawTitle.Content)
	}
}

func TestCreateEPUB(t *testing.T) {
	data := &jplaw.Law{
		Era:             jplaw.EraReiwa,
		Year:            2024,
		PromulgateMonth: 1,
		PromulgateDay:   1,
		LawNum:          "令和六年法律第一号",
		Lang:            "ja",
		LawBody: jplaw.LawBody{
			LawTitle: &jplaw.LawTitle{
				Content: "テスト法",
				Kana:    "テストホウ",
			},
		},
	}

	book := createEPUB(data)
	
	if book == nil {
		t.Fatal("createEPUB returned nil")
	}
	
	// The epub package doesn't expose getters for metadata,
	// so we can only verify that the book was created without error
}

func TestSetupEPUBMetadata(t *testing.T) {
	book, err := epub.NewEpub("Test Book")
	if err != nil {
		t.Fatalf("Failed to create test EPUB: %v", err)
	}

	data := &jplaw.Law{
		Era:             jplaw.EraReiwa,
		Year:            2024,
		PromulgateMonth: 3,
		PromulgateDay:   15,
		LawNum:          "令和六年法律第一号",
		Lang:            "ja",
		LawBody: jplaw.LawBody{
			LawTitle: &jplaw.LawTitle{
				Content: "テスト法",
				Kana:    "テストホウ",
				Ruby: []jplaw.Ruby{
					{
						Content: "法",
						Rt:      []jplaw.Rt{{Content: "ほう"}},
					},
				},
			},
		},
	}

	setupEPUBMetadata(book, data)
	
	// The epub package doesn't expose getters for metadata,
	// so we can only verify that the function runs without error
}

func TestBuildChapterBody(t *testing.T) {
	chapter := &jplaw.Chapter{
		ChapterTitle: jplaw.ChapterTitle{
			Content: "第一章　総則",
			Ruby: []jplaw.Ruby{
				{
					Content: "総",
					Rt:      []jplaw.Rt{{Content: "そう"}},
				},
			},
		},
		Section: []jplaw.Section{
			{
				SectionTitle: jplaw.SectionTitle{
					Content: "第一節　通則",
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
		},
	}

	got := buildChapterBody(chapter)

	expectedParts := []string{
		"<h2>第一章　総則<ruby>総<rt>そう</rt></ruby></h2>",
		"<div class='sections'>",
		"<h3>第一節　通則</h3>",
		"（第一条 から 第二条 まで）",
		"</div>",
	}

	for _, part := range expectedParts {
		if !strings.Contains(got, part) {
			t.Errorf("buildChapterBody() should contain %q\ngot: %v", part, got)
		}
	}
}

func TestBuildSectionsHTML(t *testing.T) {
	sections := []jplaw.Section{
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
		{
			SectionTitle: jplaw.SectionTitle{
				Content: "第二節",
				Ruby: []jplaw.Ruby{
					{
						Content: "特",
						Rt:      []jplaw.Rt{{Content: "とく"}},
					},
				},
			},
			Article: []jplaw.Article{},
		},
	}

	got := buildSectionsHTML(sections)

	expectedParts := []string{
		"<div class='sections'>",
		"<h3>第一節</h3>",
		"（第一条 から 第一条 まで）",
		"<h3>第二節<ruby>特<rt>とく</rt></ruby></h3>",
		"</div>",
	}

	for _, part := range expectedParts {
		if !strings.Contains(got, part) {
			t.Errorf("buildSectionsHTML() should contain %q\ngot: %v", part, got)
		}
	}
}

func TestProcessChapter(t *testing.T) {
	book, err := epub.NewEpub("Test Book")
	if err != nil {
		t.Fatalf("Failed to create test EPUB: %v", err)
	}

	chapter := &jplaw.Chapter{
		ChapterTitle: jplaw.ChapterTitle{
			Content: "第一章　総則",
		},
		Article: []jplaw.Article{
			{
				ArticleTitle: &jplaw.ArticleTitle{
					Content: "第一条",
				},
				Paragraph: []jplaw.Paragraph{
					{
						ParagraphSentence: jplaw.ParagraphSentence{
							Sentence: []jplaw.Sentence{
								{Content: "条文内容。"},
							},
						},
					},
				},
			},
		},
		Section: []jplaw.Section{
			{
				SectionTitle: jplaw.SectionTitle{
					Content: "第一節",
				},
				Article: []jplaw.Article{
					{
						ArticleTitle: &jplaw.ArticleTitle{
							Content: "第二条",
						},
					},
				},
			},
		},
	}

	// Process the chapter
	// Note: This will call handleError on failure which exits the program,
	// so we can only test the happy path
	processChapter(book, chapter, 0)
	
	// If we reach here, the test passed (no error occurred)
}

func TestWriteEPUB(t *testing.T) {
	book, err := epub.NewEpub("Test Book")
	if err != nil {
		t.Fatalf("Failed to create test EPUB: %v", err)
	}

	// Add some content
	book.SetAuthor("Test Author")
	book.AddSection("<h1>Test</h1>", "Test Chapter", "", "")

	// Create a temporary file
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test.epub")

	// Write the EPUB
	writeEPUB(book, outputPath)

	// Verify the file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("EPUB file was not created")
	}
}

func TestIntegration_SmallLaw(t *testing.T) {
	// Create a small but complete law structure
	data := &jplaw.Law{
		Era:             jplaw.EraReiwa,
		Year:            2024,
		PromulgateMonth: 1,
		PromulgateDay:   1,
		LawNum:          "令和六年法律第一号",
		Lang:            "ja",
		LawBody: jplaw.LawBody{
			LawTitle: &jplaw.LawTitle{
				Content: "テスト基本法",
				Kana:    "テストキホンホウ",
			},
			MainProvision: jplaw.MainProvision{
				Chapter: []jplaw.Chapter{
					{
						ChapterTitle: jplaw.ChapterTitle{
							Content: "第一章　総則",
						},
						Article: []jplaw.Article{
							{
								ArticleTitle: &jplaw.ArticleTitle{
									Content: "第一条",
								},
								ArticleCaption: &jplaw.ArticleCaption{
									Content: "（目的）",
								},
								Paragraph: []jplaw.Paragraph{
									{
										Num:          1,
										ParagraphNum: jplaw.ParagraphNum{Content: "１"},
										ParagraphSentence: jplaw.ParagraphSentence{
											Sentence: []jplaw.Sentence{
												{Content: "この法律は、テストを目的とする。"},
											},
										},
									},
									{
										Num:          2,
										ParagraphNum: jplaw.ParagraphNum{Content: "２"},
										ParagraphSentence: jplaw.ParagraphSentence{
											Sentence: []jplaw.Sentence{
												{Content: "前項の規定により実施する。"},
											},
										},
										Item: []jplaw.Item{
											{
												ItemTitle: &jplaw.ItemTitle{Content: "一"},
												ItemSentence: jplaw.ItemSentence{
													Sentence: []jplaw.Sentence{
														{Content: "第一号項目"},
													},
												},
												Subitem1: []jplaw.Subitem1{
													{
														Subitem1Title: &jplaw.Subitem1Title{Content: "イ"},
														Subitem1Sentence: jplaw.Subitem1Sentence{
															Sentence: []jplaw.Sentence{
																{Content: "イ項目"},
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
				},
			},
		},
	}

	// Create EPUB
	book := createEPUB(data)
	processChapters(book, data)

	// Write to temporary file
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "integration_test.epub")
	writeEPUB(book, outputPath)

	// Verify the file exists and has content
	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("Failed to stat output file: %v", err)
	}

	if info.Size() == 0 {
		t.Error("Generated EPUB file is empty")
	}
}