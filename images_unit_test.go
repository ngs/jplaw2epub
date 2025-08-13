package jplaw2epub

import (
	"strings"
	"testing"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

func TestNewImageProcessor(t *testing.T) {
	book, err := epub.NewEpub("Test Book")
	if err != nil {
		t.Fatalf("Failed to create epub: %v", err)
	}

	// Note: We can't create a real API client in tests, so we use nil
	proc := NewImageProcessor(nil, "test-revision", book)

	if proc == nil {
		t.Error("NewImageProcessor returned nil")
		return
	}
	if proc.revisionID != "test-revision" {
		t.Errorf("revisionID = %v, want %v", proc.revisionID, "test-revision")
	}
	if proc.book != book {
		t.Error("book reference not set correctly")
	}
	if proc.maxImageHeight != "80vh" {
		t.Errorf("Default maxImageHeight = %v, want %v", proc.maxImageHeight, "80vh")
	}
}

func TestSetMaxImageHeight(t *testing.T) {
	book, _ := epub.NewEpub("Test Book")
	proc := NewImageProcessor(nil, "test", book)

	tests := []struct {
		name   string
		height string
	}{
		{"Pixels", "300px"},
		{"Viewport height", "50vh"},
		{"Percentage", "75%"},
		{"Empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proc.SetMaxImageHeight(tt.height)
			if proc.maxImageHeight != tt.height {
				t.Errorf("SetMaxImageHeight(%v) = %v, want %v", tt.height, proc.maxImageHeight, tt.height)
			}
		})
	}
}

func TestProcessFigStruct_NoClient(t *testing.T) {
	book, _ := epub.NewEpub("Test Book")
	proc := NewImageProcessor(nil, "test", book)

	fig := &jplaw.FigStruct{
		Fig: jplaw.Fig{
			Src: "test.pdf",
		},
	}

	// Without API client, should return empty string
	result, err := proc.ProcessFigStruct(fig)
	if err == nil {
		t.Error("ProcessFigStruct should return error when client is nil")
	}
	if result != "" {
		t.Errorf("ProcessFigStruct with nil client should return empty string, got %v", result)
	}
}

func TestProcessFigStruct_EmptySrc(t *testing.T) {
	book, _ := epub.NewEpub("Test Book")
	proc := NewImageProcessor(nil, "test", book)

	fig := &jplaw.FigStruct{
		Fig: jplaw.Fig{
			Src: "", // Empty source
		},
	}

	result, err := proc.ProcessFigStruct(fig)
	if err != nil {
		t.Errorf("ProcessFigStruct with empty src should not return error, got %v", err)
	}
	if result != "" {
		t.Errorf("ProcessFigStruct with empty src should return empty string, got %v", result)
	}
}

func TestConvertToPNG_InvalidData(t *testing.T) {
	book, _ := epub.NewEpub("Test Book")
	proc := NewImageProcessor(nil, "test", book)

	// Test with invalid image data
	invalidData := []byte("not an image")

	_, err := proc.convertToPNG(invalidData, "image/jpeg")
	if err == nil {
		t.Error("convertToPNG should return error for invalid data")
	}
}

func TestBuildImageHTML_WithTitle(t *testing.T) {
	book, _ := epub.NewEpub("Test Book")
	proc := NewImageProcessor(nil, "test", book)
	proc.SetMaxImageHeight("100px")

	fig := &jplaw.FigStruct{
		FigStructTitle: &jplaw.FigStructTitle{
			Content: "図1 テスト画像",
		},
		Fig: jplaw.Fig{
			Src: "test.png",
		},
	}

	html := proc.buildImageHTML("internal/path.png", fig)

	// Check for essential elements
	expectedElements := []string{
		`class="figure"`,
		`class="figure-title"`,
		"図1 テスト画像",
		`src="internal/path.png"`,
		"max-height: 100px",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(html, expected) {
			t.Errorf("buildImageHTML should contain %q\ngot: %v", expected, html)
		}
	}
}

func TestBuildImageHTML_WithRemarks(t *testing.T) {
	book, _ := epub.NewEpub("Test Book")
	proc := NewImageProcessor(nil, "test", book)

	fig := &jplaw.FigStruct{
		Fig: jplaw.Fig{
			Src: "test.png",
		},
		Remarks: []jplaw.Remarks{
			{
				RemarksLabel: jplaw.RemarksLabel{
					Content: "注記",
				},
				Sentence: []jplaw.Sentence{
					createTestSentence("これは注記です"),
				},
			},
		},
	}

	html := proc.buildImageHTML("path.png", fig)

	expectedElements := []string{
		`class="figure-remark"`,
		`class="remarks-label"`,
		"注記",
		"これは注記です",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(html, expected) {
			t.Errorf("buildImageHTML with remarks should contain %q", expected)
		}
	}
}
