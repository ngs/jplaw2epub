package jplaw2epub

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"strings"
	"testing"

	"github.com/go-shiori/go-epub"
	lawapi "go.ngs.io/jplaw-api-v2"
	"go.ngs.io/jplaw-xml"
)

// MockAPIClient mocks the lawapi.Client for testing
type MockAPIClient struct {
	// GetAttachment behavior control
	GetAttachmentFunc func(lawRevisionID string, params *lawapi.GetAttachmentParams) (*string, error)
	GetAttachmentErr  error
	GetAttachmentData map[string]string // Map of src to base64 encoded data

	// Track calls
	GetAttachmentCalls []struct {
		RevisionID string
		Params     *lawapi.GetAttachmentParams
	}
}

// GetAttachment mocks the GetAttachment method
func (m *MockAPIClient) GetAttachment(lawRevisionID string, params *lawapi.GetAttachmentParams) (*string, error) {
	// Track the call
	m.GetAttachmentCalls = append(m.GetAttachmentCalls, struct {
		RevisionID string
		Params     *lawapi.GetAttachmentParams
	}{
		RevisionID: lawRevisionID,
		Params:     params,
	})

	// Use custom function if provided
	if m.GetAttachmentFunc != nil {
		return m.GetAttachmentFunc(lawRevisionID, params)
	}

	// Return error if set
	if m.GetAttachmentErr != nil {
		return nil, m.GetAttachmentErr
	}

	// Return data from map if available
	if m.GetAttachmentData != nil && params != nil && params.Src != nil {
		if data, ok := m.GetAttachmentData[*params.Src]; ok {
			return &data, nil
		}
	}

	// Default: return empty string
	empty := ""
	return &empty, nil
}

// Helper function to create a simple PNG image
func createTestPNGData(width, height int, fillColor color.Color) ([]byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with color
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, fillColor)
		}
	}

	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Helper function to create base64 encoded PNG data
func createBase64PNG() string {
	data, _ := createTestPNGData(10, 10, color.RGBA{255, 0, 0, 255})
	return base64.StdEncoding.EncodeToString(data)
}

const (
	testRevisionID     = "test-revision"
	defaultImageHeight = "80vh"
)

func TestImageProcessorWithMockClient(t *testing.T) {
	tests := []struct {
		name             string
		setupMock        func(*MockAPIClient)
		fig              *jplaw.FigStruct
		wantErr          bool
		wantErrContains  string
		wantHTMLContains []string
	}{
		{
			name: "successful image download and processing",
			setupMock: func(m *MockAPIClient) {
				m.GetAttachmentData = map[string]string{
					"test.png": createBase64PNG(),
				}
			},
			fig: &jplaw.FigStruct{
				Fig: jplaw.Fig{
					Src: "test.png",
				},
				FigStructTitle: &jplaw.FigStructTitle{
					Content: "テスト画像",
				},
			},
			wantErr: false,
			wantHTMLContains: []string{
				`<div class="figure"`,
				`<img`,
				`alt="Figure"`,
				`<p class="figure-title">テスト画像</p>`,
			},
		},
		{
			name: "image download error",
			setupMock: func(m *MockAPIClient) {
				m.GetAttachmentErr = fmt.Errorf("network error")
			},
			fig: &jplaw.FigStruct{
				Fig: jplaw.Fig{
					Src: "error.png",
				},
			},
			wantErr:         true,
			wantErrContains: "downloading image",
		},
		{
			name: "cached image",
			setupMock: func(m *MockAPIClient) {
				// No setup needed - we'll test caching behavior
				m.GetAttachmentData = map[string]string{
					"cached.png": createBase64PNG(),
				}
			},
			fig: &jplaw.FigStruct{
				Fig: jplaw.Fig{
					Src: "cached.png",
				},
			},
			wantErr: false,
		},
		{
			name: "empty src",
			setupMock: func(m *MockAPIClient) {
				// No setup needed
			},
			fig: &jplaw.FigStruct{
				Fig: jplaw.Fig{
					Src: "",
				},
			},
			wantErr: false,
		},
		{
			name: "with remarks",
			setupMock: func(m *MockAPIClient) {
				m.GetAttachmentData = map[string]string{
					"with-remarks.png": createBase64PNG(),
				}
			},
			fig: &jplaw.FigStruct{
				Fig: jplaw.Fig{
					Src: "with-remarks.png",
				},
				Remarks: []jplaw.Remarks{
					{
						RemarksLabel: jplaw.RemarksLabel{
							Content: "備考",
						},
						Sentence: []jplaw.Sentence{
							{
								MixedContent: jplaw.MixedContent{
									Nodes: []jplaw.ContentNode{
										jplaw.TextNode{Text: "これは備考です"},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
			wantHTMLContains: []string{
				`<div class="figure"`,
				`<div class="figure-remark">`,
				`<p class="remarks-label">備考</p>`,
				"これは備考です",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			mockClient := &MockAPIClient{}
			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}

			// Create EPUB book
			book, err := epub.NewEpub("Test Book")
			if err != nil {
				t.Fatalf("Failed to create EPUB: %v", err)
			}

			// Create ImageProcessor with mock client
			imgProc := &ImageProcessor{
				client:         mockClient,
				revisionID:     testRevisionID,
				book:           book,
				imageCache:     make(map[string]string),
				maxImageHeight: defaultImageHeight,
			}

			// Process the figure
			html, err := imgProc.ProcessFigStruct(tt.fig)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessFigStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.wantErrContains != "" {
				if !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("ProcessFigStruct() error = %v, want error containing %v", err, tt.wantErrContains)
				}
				return
			}

			// Check HTML output
			for _, want := range tt.wantHTMLContains {
				if !strings.Contains(html, want) {
					t.Errorf("ProcessFigStruct() HTML missing expected content: %s\nGot: %s", want, html)
				}
			}
		})
	}
}

func TestImageProcessorCaching(t *testing.T) {
	// Create mock client with tracking
	mockClient := &MockAPIClient{
		GetAttachmentData: map[string]string{
			"test.png": createBase64PNG(),
		},
	}

	book, err := epub.NewEpub("Test Book")
	if err != nil {
		t.Fatalf("Failed to create EPUB: %v", err)
	}

	imgProc := &ImageProcessor{
		client:         mockClient,
		revisionID:     testRevisionID,
		book:           book,
		imageCache:     make(map[string]string),
		maxImageHeight: defaultImageHeight,
	}

	fig := &jplaw.FigStruct{
		Fig: jplaw.Fig{
			Src: "test.png",
		},
	}

	// First call - should download
	_, err = imgProc.ProcessFigStruct(fig)
	if err != nil {
		t.Errorf("First ProcessFigStruct() unexpected error: %v", err)
	}

	if len(mockClient.GetAttachmentCalls) != 1 {
		t.Errorf("Expected 1 GetAttachment call, got %d", len(mockClient.GetAttachmentCalls))
	}

	// Second call - should use cache
	_, err = imgProc.ProcessFigStruct(fig)
	if err != nil {
		t.Errorf("Second ProcessFigStruct() unexpected error: %v", err)
	}

	// Should still be only 1 call (cached)
	if len(mockClient.GetAttachmentCalls) != 1 {
		t.Errorf("Expected 1 GetAttachment call (cached), got %d", len(mockClient.GetAttachmentCalls))
	}
}

func TestImageProcessorWithNilClient(t *testing.T) {
	book, err := epub.NewEpub("Test Book")
	if err != nil {
		t.Fatalf("Failed to create EPUB: %v", err)
	}

	// Create ImageProcessor without client
	imgProc := &ImageProcessor{
		client:         nil,
		revisionID:     testRevisionID,
		book:           book,
		imageCache:     make(map[string]string),
		maxImageHeight: defaultImageHeight,
	}

	fig := &jplaw.FigStruct{
		Fig: jplaw.Fig{
			Src: "test.png",
		},
	}

	_, err = imgProc.ProcessFigStruct(fig)
	if err == nil {
		t.Error("Expected error with nil client, got nil")
	}

	if !strings.Contains(err.Error(), "API client is not configured") {
		t.Errorf("Expected 'API client is not configured' error, got: %v", err)
	}
}

func TestDownloadImage(t *testing.T) {
	tests := []struct {
		name        string
		src         string
		setupMock   func(*MockAPIClient)
		wantErr     bool
		wantContent string
	}{
		{
			name: "successful download",
			src:  "image.jpg",
			setupMock: func(m *MockAPIClient) {
				data := "test image data"
				m.GetAttachmentData = map[string]string{
					"image.jpg": data,
				}
			},
			wantErr:     false,
			wantContent: "test image data",
		},
		{
			name: "API error",
			src:  "error.jpg",
			setupMock: func(m *MockAPIClient) {
				m.GetAttachmentErr = fmt.Errorf("API error")
			},
			wantErr: true,
		},
		{
			name: "nil attachment",
			src:  "nil.jpg",
			setupMock: func(m *MockAPIClient) {
				m.GetAttachmentFunc = func(string, *lawapi.GetAttachmentParams) (*string, error) {
					return nil, nil
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAPIClient{}
			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}

			imgProc := &ImageProcessor{
				client:     mockClient,
				revisionID: testRevisionID,
			}

			data, contentType, err := imgProc.downloadImage(tt.src)

			if (err != nil) != tt.wantErr {
				t.Errorf("downloadImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if string(data) != tt.wantContent {
					t.Errorf("downloadImage() data = %v, want %v", string(data), tt.wantContent)
				}

				if contentType == "" {
					t.Error("downloadImage() contentType is empty")
				}
			}

			// Verify API was called with correct parameters
			if len(mockClient.GetAttachmentCalls) > 0 {
				call := mockClient.GetAttachmentCalls[0]
				if call.RevisionID != testRevisionID {
					t.Errorf("GetAttachment called with wrong RevisionID: %v", call.RevisionID)
				}
				if call.Params == nil || call.Params.Src == nil || *call.Params.Src != tt.src {
					t.Errorf("GetAttachment called with wrong src parameter")
				}
			}
		})
	}
}

// Test ImageProcessor creation with mock client
func TestNewImageProcessorWithMockClient(t *testing.T) {
	mockClient := &MockAPIClient{}
	book, _ := epub.NewEpub("Test Book")

	imgProc := NewImageProcessor(mockClient, testRevisionID, book)

	if imgProc.client != mockClient {
		t.Error("ImageProcessor client not set correctly")
	}

	if imgProc.revisionID != testRevisionID {
		t.Error("ImageProcessor revisionID not set correctly")
	}

	if imgProc.book != book {
		t.Error("ImageProcessor book not set correctly")
	}

	if imgProc.maxImageHeight != defaultImageHeight {
		t.Errorf("ImageProcessor default maxImageHeight = %v, want %v", imgProc.maxImageHeight, defaultImageHeight)
	}

	if imgProc.imageCache == nil {
		t.Error("ImageProcessor imageCache not initialized")
	}
}

// Benchmark with mock client
func BenchmarkProcessFigStructWithMockClient(b *testing.B) {
	mockClient := &MockAPIClient{
		GetAttachmentData: map[string]string{
			"bench.png": createBase64PNG(),
		},
	}

	book, _ := epub.NewEpub("Test Book")
	imgProc := &ImageProcessor{
		client:         mockClient,
		revisionID:     testRevisionID,
		book:           book,
		imageCache:     make(map[string]string),
		maxImageHeight: defaultImageHeight,
	}

	fig := &jplaw.FigStruct{
		Fig: jplaw.Fig{
			Src: "bench.png",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Clear cache each iteration to test full processing
		imgProc.imageCache = make(map[string]string)
		_, _ = imgProc.ProcessFigStruct(fig)
	}
}

// Test for validating mock client interface compatibility
func TestMockClientInterfaceCompatibility(t *testing.T) {
	// This test ensures our MockAPIClient can be used wherever lawapi.Client is expected
	var _ interface {
		GetAttachment(string, *lawapi.GetAttachmentParams) (*string, error)
	} = (*MockAPIClient)(nil)
}
