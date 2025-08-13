package jplaw2epub

import (
	"encoding/base64"
	"fmt"

	"github.com/go-shiori/go-epub"
)

// CSS styles for EPUB formatting
const epubCSS = `
/* Base styles */
body {
    font-family: "Hiragino Kaku Gothic ProN", "ヒラギノ角ゴ ProN W3", "Meiryo", "メイリオ", sans-serif;
    line-height: 1.6;
    margin: 1em;
}

h1, h2, h3, h4, h5, h6 {
    color: #333;
    margin-top: 1.5em;
    margin-bottom: 0.5em;
}

/* Figure styles */
.figure {
    margin: 1em 0;
    text-align: center;
    page-break-inside: avoid;
}

.figure img {
    max-width: 100%;
    max-height: 70vh; /* Limit height to 70% of viewport height */
    height: auto;
    display: block;
    margin: 0 auto;
}

.figure-title {
    font-weight: bold;
    margin-bottom: 0.5em;
    color: #555;
}

.figure-remark {
    font-size: 0.9em;
    color: #666;
    margin-top: 0.5em;
    font-style: italic;
}

/* Style structure styles */
.style-struct {
    margin: 1em 0;
    border-left: 3px solid #ddd;
    padding-left: 1em;
}

.style-title {
    font-weight: bold;
    margin-bottom: 0.5em;
    color: #444;
}

.style-content {
    margin: 0.5em 0;
}

/* Appendix styles */
.appdx-remarks {
    margin: 1em 0;
    padding: 0.5em;
    background-color: #f9f9f9;
    border-radius: 4px;
}

.remarks-label {
    font-weight: bold;
    color: #666;
}

.remark {
    margin: 0.5em 0;
}

/* Related articles */
.related-articles {
    margin: 0.5em 0;
    padding: 0.5em;
    background-color: #f0f8ff;
    border-radius: 4px;
    font-size: 0.9em;
}

/* Table styles */
.table-struct {
    margin: 1em 0;
}

.table-title {
    font-weight: bold;
    margin-bottom: 0.5em;
    color: #444;
}

.table-placeholder {
    text-align: center;
    padding: 2em;
    background-color: #f5f5f5;
    border: 1px dashed #ccc;
    color: #666;
}

/* List styles */
ol, ul {
    margin: 0.5em 0;
    padding-left: 2em;
}

li {
    margin: 0.25em 0;
}

/* Strong and emphasis */
strong {
    font-weight: bold;
    color: #333;
}

/* Page break hints */
.chapter-break {
    page-break-before: always;
}

.section-break {
    page-break-before: auto;
}

/* Print and e-reader specific styles */
@media print, screen and (max-device-width: 1024px) {
    .figure img {
        max-height: 60vh; /* Slightly smaller on smaller screens */
    }
    
    body {
        margin: 0.5em;
    }
}
`

// AddCSSToEPUB adds CSS stylesheet to the EPUB
func AddCSSToEPUB(book *epub.Epub) error {
	// Create a data URL for the CSS content
	dataURL := fmt.Sprintf("data:text/css;base64,%s", base64.StdEncoding.EncodeToString([]byte(epubCSS)))

	// Add CSS to EPUB using data URL
	cssPath, err := book.AddCSS(dataURL, "styles.css")
	if err != nil {
		return fmt.Errorf("adding CSS to EPUB: %w", err)
	}

	// The CSS file will be automatically linked to all HTML files
	// by the go-epub library
	_ = cssPath // CSS path for reference if needed

	return nil
}
