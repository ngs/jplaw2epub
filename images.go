package jplaw2epub

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"path"
	"path/filepath"
	"strings"

	"github.com/gen2brain/go-fitz"
	"github.com/go-shiori/go-epub"
	lawapi "go.ngs.io/jplaw-api-v2"
	"go.ngs.io/jplaw-xml"
)

// ImageProcessor handles image processing for EPUB
type ImageProcessor struct {
	client         APIClient
	revisionID     string
	book           *epub.Epub
	imageCache     map[string]string // maps src to EPUB internal path
	maxImageHeight string            // maximum height for images (CSS value)
}

// NewImageProcessor creates a new image processor
func NewImageProcessor(client APIClient, revisionID string, book *epub.Epub) *ImageProcessor {
	return &ImageProcessor{
		client:         client,
		revisionID:     revisionID,
		book:           book,
		imageCache:     make(map[string]string),
		maxImageHeight: "80vh", // default height
	}
}

// SetMaxImageHeight sets the maximum image height
func (ip *ImageProcessor) SetMaxImageHeight(height string) {
	ip.maxImageHeight = height
}

// ProcessFigStruct processes a FigStruct and returns HTML
func (ip *ImageProcessor) ProcessFigStruct(fig *jplaw.FigStruct) (string, error) {
	if fig.Fig.Src == "" {
		return "", nil
	}

	// Check cache first
	if epubPath, exists := ip.imageCache[fig.Fig.Src]; exists {
		return ip.buildImageHTML(epubPath, fig), nil
	}

	if ip.client == nil {
		return "", fmt.Errorf("API client is not configured")
	}

	// Download image
	imageData, contentType, err := ip.downloadImage(fig.Fig.Src)
	if err != nil {
		return "", fmt.Errorf("downloading image %s: %w", fig.Fig.Src, err)
	}

	// Convert to PNG if necessary
	if !isPNG(contentType) {
		imageData, err = ip.convertToPNG(imageData, contentType)
		if err != nil {
			return "", fmt.Errorf("converting image to PNG: %w", err)
		}
	}

	// Add image to EPUB
	epubPath, err := ip.addImageToEPUB(fig.Fig.Src, imageData)
	if err != nil {
		return "", fmt.Errorf("adding image to EPUB: %w", err)
	}

	// Cache the path
	ip.imageCache[fig.Fig.Src] = epubPath

	return ip.buildImageHTML(epubPath, fig), nil
}

// downloadImage downloads an image from the API
func (ip *ImageProcessor) downloadImage(src string) (data []byte, contentType string, err error) {
	params := &lawapi.GetAttachmentParams{
		Src: lawapi.StringPtr(src),
	}

	attachment, err := ip.client.GetAttachment(ip.revisionID, params)
	if err != nil {
		return nil, "", err
	}

	if attachment == nil {
		return nil, "", fmt.Errorf("attachment is nil")
	}

	// The API returns binary data as a string
	data = []byte(*attachment)

	// Guess content type from file extension
	contentType = guessContentType(src)

	return data, contentType, nil
}

// convertToPNG converts various image formats to PNG
func (ip *ImageProcessor) convertToPNG(data []byte, contentType string) ([]byte, error) {
	reader := bytes.NewReader(data)

	// Decode the image based on content type
	var img image.Image
	var err error

	switch {
	case strings.Contains(contentType, "jpeg") || strings.Contains(contentType, "jpg"):
		// Import jpeg decoder
		img, _, err = image.Decode(reader)
	case strings.Contains(contentType, "gif"):
		// Import gif decoder
		img, _, err = image.Decode(reader)
	case strings.Contains(contentType, "pdf"):
		// Convert PDF to PNG using go-fitz
		return ip.convertPDFToPNG(data)
	default:
		// Try generic decode
		img, _, err = image.Decode(reader)
	}

	if err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}

	// Encode as PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("encoding PNG: %w", err)
	}

	return buf.Bytes(), nil
}

// addImageToEPUB adds an image to the EPUB and returns its internal path
func (ip *ImageProcessor) addImageToEPUB(src string, data []byte) (string, error) {
	// Generate a unique filename based on the source
	filename := generateImageFilename(src)

	// Create a data URL from the PNG data
	// This avoids file system issues with temporary files
	dataURL := fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(data))

	// Add image to EPUB using data URL
	epubPath, err := ip.book.AddImage(dataURL, filename)
	if err != nil {
		return "", fmt.Errorf("adding image to EPUB: %w", err)
	}

	return epubPath, nil
}

// buildImageHTML builds the HTML for an image
func (ip *ImageProcessor) buildImageHTML(epubPath string, fig *jplaw.FigStruct) string {
	html := `<div class="figure" style="page-break-inside: avoid; margin: 1em 0; text-align: center;">`

	// Add title if present
	if fig.FigStructTitle != nil && fig.FigStructTitle.Content != "" {
		titleHTML := processTextWithRuby(fig.FigStructTitle.Content, fig.FigStructTitle.Ruby)
		html += fmt.Sprintf(`<p class="figure-title">%s</p>`, titleHTML)
	}

	// Add image with configurable size constraints for EPUB readers
	imgStyle := fmt.Sprintf("max-width: 100%%; max-height: %s; "+
		"height: auto; display: block; margin: 0 auto; page-break-inside: avoid;",
		ip.maxImageHeight)
	html += fmt.Sprintf(`<img src=%q alt="Figure" style=%q />`, epubPath, imgStyle)

	// Add remarks if present
	for i := range fig.Remarks {
		remark := &fig.Remarks[i]
		html += `<div class="figure-remark">`

		// Add remarks label if present
		if remark.RemarksLabel.Content != "" {
			html += fmt.Sprintf(`<p class="remarks-label">%s</p>`,
				processTextWithRuby(remark.RemarksLabel.Content, remark.RemarksLabel.Ruby))
		}

		// Add sentences
		for j := range remark.Sentence {
			html += fmt.Sprintf(`<p>%s</p>`, remark.Sentence[j].HTML())
		}

		// Add items if present
		if len(remark.Item) > 0 {
			html += processItems(remark.Item)
		}

		html += htmlDivEnd
	}

	html += htmlDivEnd
	return html
}

// isPNG checks if the content type is PNG
func isPNG(contentType string) bool {
	return strings.Contains(strings.ToLower(contentType), "png")
}

// guessContentType guesses content type from file extension
func guessContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}

// generateImageFilename generates a unique filename for an image
func generateImageFilename(src string) string {
	// Extract the base name
	base := path.Base(src)

	// Change extension to .png
	ext := filepath.Ext(base)
	if ext != "" {
		base = strings.TrimSuffix(base, ext)
	}

	return base + ".png"
}

// convertPDFToPNG converts a PDF to PNG format using go-fitz
func (ip *ImageProcessor) convertPDFToPNG(pdfData []byte) ([]byte, error) {
	// Create a new document from the PDF data
	doc, err := fitz.NewFromMemory(pdfData)
	if err != nil {
		return nil, fmt.Errorf("opening PDF: %w", err)
	}
	defer doc.Close()

	// Get the first page (index 0)
	if doc.NumPage() == 0 {
		return nil, fmt.Errorf("PDF has no pages")
	}

	// Render the first page as an image
	// Using a scale factor of 2.0 for better quality (144 DPI instead of 72 DPI)
	img, err := doc.Image(0)
	if err != nil {
		return nil, fmt.Errorf("rendering PDF page: %w", err)
	}

	// Encode the image as PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("encoding PNG: %w", err)
	}

	return buf.Bytes(), nil
}
