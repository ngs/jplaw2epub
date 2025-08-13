package jplaw2epub

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/go-shiori/go-epub"
	"go.ngs.io/jplaw-xml"
)

func TestGenerateImageFilename(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
	}{
		{
			name:     "PDF file",
			src:      "./pict/S18F03402004001-001.pdf",
			expected: "S18F03402004001-001.png",
		},
		{
			name:     "JPEG file",
			src:      "./images/test-image.jpg",
			expected: "test-image.png",
		},
		{
			name:     "Already PNG file",
			src:      "./images/test.png",
			expected: "test.png",
		},
		{
			name:     "No extension",
			src:      "image",
			expected: "image.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateImageFilename(tt.src)
			if result != tt.expected {
				t.Errorf("generateImageFilename(%q) = %q, want %q", tt.src, result, tt.expected)
			}
		})
	}
}

func TestGuessContentType(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected string
	}{
		{
			name:     "PNG file",
			filename: "test.png",
			expected: "image/png",
		},
		{
			name:     "JPEG file",
			filename: "test.jpg",
			expected: "image/jpeg",
		},
		{
			name:     "JPEG file with jpeg extension",
			filename: "test.jpeg",
			expected: "image/jpeg",
		},
		{
			name:     "GIF file",
			filename: "test.gif",
			expected: "image/gif",
		},
		{
			name:     "PDF file",
			filename: "test.pdf",
			expected: "application/pdf",
		},
		{
			name:     "Unknown file",
			filename: "test.xyz",
			expected: "application/octet-stream",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := guessContentType(tt.filename)
			if result != tt.expected {
				t.Errorf("guessContentType(%q) = %q, want %q", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestIsPNG(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		expected    bool
	}{
		{
			name:        "PNG content type",
			contentType: "image/png",
			expected:    true,
		},
		{
			name:        "PNG with charset",
			contentType: "image/png; charset=utf-8",
			expected:    true,
		},
		{
			name:        "JPEG content type",
			contentType: "image/jpeg",
			expected:    false,
		},
		{
			name:        "Empty content type",
			contentType: "",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPNG(tt.contentType)
			if result != tt.expected {
				t.Errorf("isPNG(%q) = %v, want %v", tt.contentType, result, tt.expected)
			}
		})
	}
}

func TestBuildImageHTML(t *testing.T) {
	book, _ := epub.NewEpub("Test Book")
	imgProc := &ImageProcessor{
		book:       book,
		imageCache: make(map[string]string),
	}

	tests := []struct {
		name     string
		epubPath string
		fig      *jplaw.FigStruct
		contains []string
	}{
		{
			name:     "Simple figure",
			epubPath: "images/test.png",
			fig: &jplaw.FigStruct{
				Fig: jplaw.Fig{Src: "./pict/test.pdf"},
			},
			contains: []string{
				`<div class="figure">`,
				`<img src="images/test.png" alt="Figure" />`,
				`</div>`,
			},
		},
		{
			name:     "Figure with title",
			epubPath: "images/test.png",
			fig: &jplaw.FigStruct{
				FigStructTitle: &jplaw.FigStructTitle{
					Content: "図1 テスト画像",
				},
				Fig: jplaw.Fig{Src: "./pict/test.pdf"},
			},
			contains: []string{
				`<p class="figure-title">図1 テスト画像</p>`,
				`<img src="images/test.png" alt="Figure" />`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := imgProc.buildImageHTML(tt.epubPath, tt.fig)
			for _, expected := range tt.contains {
				if !bytes.Contains([]byte(result), []byte(expected)) {
					t.Errorf("buildImageHTML() result doesn't contain %q\nGot: %q", expected, result)
				}
			}
		})
	}
}

func TestConvertToPNG(t *testing.T) {
	// Create a simple test image
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	// Encode as PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("Failed to create test PNG: %v", err)
	}

	book, _ := epub.NewEpub("Test Book")
	imgProc := &ImageProcessor{
		book:       book,
		imageCache: make(map[string]string),
	}

	// Test PNG passthrough
	result, err := imgProc.convertToPNG(buf.Bytes(), "image/png")
	if err != nil {
		t.Errorf("convertToPNG failed for PNG: %v", err)
	}

	// Verify it's still a valid PNG
	_, err = png.Decode(bytes.NewReader(result))
	if err != nil {
		t.Errorf("Result is not a valid PNG: %v", err)
	}
}