package jplaw2epub

import (
	"fmt"
	"testing"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

// MockImageProcessor is a mock implementation of ImageProcessorInterface for testing
type MockImageProcessor struct {
	// Control behavior
	ProcessFigStructFunc func(*jplaw.FigStruct) (string, error)
	ProcessFigStructErr  error
	ProcessFigStructHTML string

	// Track calls
	ProcessFigStructCalls  []*jplaw.FigStruct
	SetMaxImageHeightCalls []string
}

// Ensure MockImageProcessor implements ImageProcessorInterface
var _ ImageProcessorInterface = (*MockImageProcessor)(nil)

// ProcessFigStruct mocks the ProcessFigStruct method
func (m *MockImageProcessor) ProcessFigStruct(fig *jplaw.FigStruct) (string, error) {
	m.ProcessFigStructCalls = append(m.ProcessFigStructCalls, fig)

	if m.ProcessFigStructFunc != nil {
		return m.ProcessFigStructFunc(fig)
	}

	if m.ProcessFigStructErr != nil {
		return "", m.ProcessFigStructErr
	}

	if m.ProcessFigStructHTML != "" {
		return m.ProcessFigStructHTML, nil
	}

	// Default behavior - return simple HTML
	return fmt.Sprintf(`<img src="mock-%s.png" alt="Figure"/>`, fig.Fig.Src), nil
}

// SetMaxImageHeight mocks the SetMaxImageHeight method
func (m *MockImageProcessor) SetMaxImageHeight(height string) {
	m.SetMaxImageHeightCalls = append(m.SetMaxImageHeightCalls, height)
}

// Test using MockImageProcessor
func TestProcessParagraphWithImagesMocked(t *testing.T) {
	tests := []struct {
		name         string
		para         *jplaw.Paragraph
		mockBehavior func(*MockImageProcessor)
		wantContains []string
		wantCalls    int
	}{
		{
			name: "paragraph with FigStruct - success",
			para: &jplaw.Paragraph{
				ParagraphNum: jplaw.ParagraphNum{Content: "1"},
				ParagraphSentence: jplaw.ParagraphSentence{
					Sentence: []jplaw.Sentence{{
						MixedContent: jplaw.MixedContent{
							Nodes: []jplaw.ContentNode{
								jplaw.TextNode{Text: "テスト段落"},
							},
						},
					}},
				},
				FigStruct: []jplaw.FigStruct{
					{
						Fig: jplaw.Fig{
							Src: "image1.jpg",
						},
					},
				},
			},
			mockBehavior: func(m *MockImageProcessor) {
				m.ProcessFigStructHTML = `<div class="figure"><img src="images/image1.jpg" alt="Figure 1"/></div>`
			},
			wantContains: []string{
				"テスト段落",
				`<div class="figure"><img src="images/image1.jpg" alt="Figure 1"/></div>`,
			},
			wantCalls: 1,
		},
		{
			name: "paragraph with multiple FigStructs",
			para: &jplaw.Paragraph{
				ParagraphNum: jplaw.ParagraphNum{Content: "2"},
				ParagraphSentence: jplaw.ParagraphSentence{
					Sentence: []jplaw.Sentence{{
						MixedContent: jplaw.MixedContent{
							Nodes: []jplaw.ContentNode{
								jplaw.TextNode{Text: "複数の図"},
							},
						},
					}},
				},
				FigStruct: []jplaw.FigStruct{
					{Fig: jplaw.Fig{Src: "image1.jpg"}},
					{Fig: jplaw.Fig{Src: "image2.jpg"}},
				},
			},
			mockBehavior: func(m *MockImageProcessor) {
				m.ProcessFigStructFunc = func(fig *jplaw.FigStruct) (string, error) {
					return fmt.Sprintf(`<img src=%q alt="Fig"/>`, fig.Fig.Src), nil
				}
			},
			wantContains: []string{
				"複数の図",
				`<img src="image1.jpg" alt="Fig"/>`,
				`<img src="image2.jpg" alt="Fig"/>`,
			},
			wantCalls: 2,
		},
		{
			name: "paragraph with FigStruct - error handling",
			para: &jplaw.Paragraph{
				ParagraphNum: jplaw.ParagraphNum{Content: "3"},
				ParagraphSentence: jplaw.ParagraphSentence{
					Sentence: []jplaw.Sentence{{
						MixedContent: jplaw.MixedContent{
							Nodes: []jplaw.ContentNode{
								jplaw.TextNode{Text: "エラーテスト"},
							},
						},
					}},
				},
				FigStruct: []jplaw.FigStruct{
					{Fig: jplaw.Fig{Src: "error.jpg"}},
				},
			},
			mockBehavior: func(m *MockImageProcessor) {
				m.ProcessFigStructErr = fmt.Errorf("image processing failed")
			},
			wantContains: []string{
				"エラーテスト",
				// Error case - image HTML should not be included
			},
			wantCalls: 1,
		},
		{
			name: "paragraph without FigStruct",
			para: &jplaw.Paragraph{
				ParagraphNum: jplaw.ParagraphNum{Content: "4"},
				ParagraphSentence: jplaw.ParagraphSentence{
					Sentence: []jplaw.Sentence{{
						MixedContent: jplaw.MixedContent{
							Nodes: []jplaw.ContentNode{
								jplaw.TextNode{Text: "図なし"},
							},
						},
					}},
				},
				FigStruct: []jplaw.FigStruct{},
			},
			mockBehavior: func(m *MockImageProcessor) {
				// No special behavior needed
			},
			wantContains: []string{
				"図なし",
			},
			wantCalls: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock
			mock := &MockImageProcessor{}
			if tt.mockBehavior != nil {
				tt.mockBehavior(mock)
			}

			// Call the function with mock
			result := processParagraphWithImages(tt.para, mock)

			// Verify the HTML contains expected content
			for _, want := range tt.wantContains {
				if !contains(result, want) {
					t.Errorf("Result should contain %q, got: %v", want, result)
				}
			}

			// Verify the number of calls
			if len(mock.ProcessFigStructCalls) != tt.wantCalls {
				t.Errorf("Expected %d calls to ProcessFigStruct, got %d",
					tt.wantCalls, len(mock.ProcessFigStructCalls))
			}

			// Verify the correct FigStructs were passed
			if tt.wantCalls > 0 && len(tt.para.FigStruct) > 0 {
				for i, call := range mock.ProcessFigStructCalls {
					if call.Fig.Src != tt.para.FigStruct[i].Fig.Src {
						t.Errorf("Call %d: expected Src %q, got %q",
							i, tt.para.FigStruct[i].Fig.Src, call.Fig.Src)
					}
				}
			}
		})
	}
}

// Test for processChapterWithImages using mock
func TestProcessChapterWithImagesMocked(t *testing.T) {
	book, err := epub.NewEpub("Test Book")
	if err != nil {
		t.Fatalf("Failed to create EPUB: %v", err)
	}

	chapter := &jplaw.Chapter{
		ChapterTitle: jplaw.ChapterTitle{
			Content: "第一章",
		},
		Article: []jplaw.Article{
			{
				ArticleTitle: &jplaw.ArticleTitle{
					Content: "第一条",
				},
				Paragraph: []jplaw.Paragraph{
					{
						ParagraphSentence: jplaw.ParagraphSentence{
							Sentence: []jplaw.Sentence{{
								MixedContent: jplaw.MixedContent{
									Nodes: []jplaw.ContentNode{
										jplaw.TextNode{Text: "条文内容"},
									},
								},
							}},
						},
						FigStruct: []jplaw.FigStruct{
							{Fig: jplaw.Fig{Src: "article_image.jpg"}},
						},
					},
				},
			},
		},
	}

	t.Run("with mock image processor", func(t *testing.T) {
		mock := &MockImageProcessor{
			ProcessFigStructHTML: `<div class="article-figure"><img src="processed.jpg"/></div>`,
		}

		err := processChapterWithImages(book, chapter, 0, mock)
		if err != nil {
			t.Errorf("processChapterWithImages() error = %v", err)
		}

		// Verify ProcessFigStruct was called
		if len(mock.ProcessFigStructCalls) != 1 {
			t.Errorf("Expected 1 call to ProcessFigStruct, got %d", len(mock.ProcessFigStructCalls))
		}
	})

	t.Run("with nil image processor", func(t *testing.T) {
		// Should not panic with nil processor
		err := processChapterWithImages(book, chapter, 1, nil)
		if err != nil {
			t.Errorf("processChapterWithImages() with nil processor error = %v", err)
		}
	})
}

// Test for processMainProvision with mock
func TestProcessMainProvisionWithMock(t *testing.T) {
	book, err := epub.NewEpub("Test Book")
	if err != nil {
		t.Fatalf("Failed to create EPUB: %v", err)
	}

	mainProv := &jplaw.MainProvision{
		Paragraph: []jplaw.Paragraph{
			{
				Num:          1,
				ParagraphNum: jplaw.ParagraphNum{Content: "1"},
				ParagraphSentence: jplaw.ParagraphSentence{
					Sentence: []jplaw.Sentence{{
						MixedContent: jplaw.MixedContent{
							Nodes: []jplaw.ContentNode{
								jplaw.TextNode{Text: "第一項"},
							},
						},
					}},
				},
				FigStruct: []jplaw.FigStruct{
					{Fig: jplaw.Fig{Src: "para1.jpg"}},
				},
			},
			{
				Num:          2,
				ParagraphNum: jplaw.ParagraphNum{Content: "2"},
				ParagraphSentence: jplaw.ParagraphSentence{
					Sentence: []jplaw.Sentence{{
						MixedContent: jplaw.MixedContent{
							Nodes: []jplaw.ContentNode{
								jplaw.TextNode{Text: "第二項"},
							},
						},
					}},
				},
				FigStruct: []jplaw.FigStruct{
					{Fig: jplaw.Fig{Src: "para2.jpg"}},
				},
			},
		},
	}

	t.Run("multiple paragraphs with images", func(t *testing.T) {
		callCount := 0
		mock := &MockImageProcessor{
			ProcessFigStructFunc: func(fig *jplaw.FigStruct) (string, error) {
				callCount++
				return fmt.Sprintf(`<img src="mock-%d.jpg" alt="Mock %d"/>`, callCount, callCount), nil
			},
		}

		err := processMainProvision(book, mainProv, mock)
		if err != nil {
			t.Errorf("processMainProvision() error = %v", err)
		}

		// Verify both FigStructs were processed
		if len(mock.ProcessFigStructCalls) != 2 {
			t.Errorf("Expected 2 calls to ProcessFigStruct, got %d", len(mock.ProcessFigStructCalls))
		}

		// Verify the correct images were processed
		if mock.ProcessFigStructCalls[0].Fig.Src != "para1.jpg" {
			t.Errorf("First call: expected Src 'para1.jpg', got %q", mock.ProcessFigStructCalls[0].Fig.Src)
		}
		if mock.ProcessFigStructCalls[1].Fig.Src != "para2.jpg" {
			t.Errorf("Second call: expected Src 'para2.jpg', got %q", mock.ProcessFigStructCalls[1].Fig.Src)
		}
	})
}

// Helper function for testing
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || substr == "" ||
		(s != "" && substr != "" &&
			(s[0:len(substr)] == substr || contains(s[1:], substr))))
}

// Test for createImageProcessor with options
func TestCreateImageProcessorWithOptions(t *testing.T) {
	book, err := epub.NewEpub("Test Book")
	if err != nil {
		t.Fatalf("Failed to create EPUB: %v", err)
	}

	t.Run("with nil options", func(t *testing.T) {
		imgProc := createImageProcessor(book, nil)
		if imgProc != nil {
			t.Errorf("Expected nil processor for nil options, got %v", imgProc)
		}
	})

	t.Run("with empty options", func(t *testing.T) {
		opts := &EPUBOptions{}
		imgProc := createImageProcessor(book, opts)
		if imgProc != nil {
			t.Errorf("Expected nil processor for empty options, got %v", imgProc)
		}
	})

	// Note: Full ImageProcessor creation requires API client,
	// which we can't easily test without external dependencies
}

// Benchmark with mock
func BenchmarkProcessParagraphWithMockImages(b *testing.B) {
	para := &jplaw.Paragraph{
		ParagraphNum: jplaw.ParagraphNum{Content: "1"},
		ParagraphSentence: jplaw.ParagraphSentence{
			Sentence: []jplaw.Sentence{{
				MixedContent: jplaw.MixedContent{
					Nodes: []jplaw.ContentNode{
						jplaw.TextNode{Text: "ベンチマークテスト"},
					},
				},
			}},
		},
		FigStruct: []jplaw.FigStruct{
			{Fig: jplaw.Fig{Src: "bench1.jpg"}},
			{Fig: jplaw.Fig{Src: "bench2.jpg"}},
			{Fig: jplaw.Fig{Src: "bench3.jpg"}},
		},
	}

	mock := &MockImageProcessor{
		ProcessFigStructFunc: func(fig *jplaw.FigStruct) (string, error) {
			return fmt.Sprintf(`<img src=%q/>`, fig.Fig.Src), nil
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processParagraphWithImages(para, mock)
	}
}
