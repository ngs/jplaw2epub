package main

import (
	"html"
	"strings"

	"go.ngs.io/jplaw-xml"
)

// processRubyElements converts Ruby elements to HTML ruby tags
func processRubyElements(rubies []jplaw.Ruby) string {
	var result strings.Builder
	for _, ruby := range rubies {
		if len(ruby.Rt) > 0 {
			result.WriteString("<ruby>")
			result.WriteString(html.EscapeString(ruby.Content))
			for _, rt := range ruby.Rt {
				result.WriteString("<rt>")
				result.WriteString(html.EscapeString(rt.Content))
				result.WriteString("</rt>")
			}
			result.WriteString("</ruby>")
		} else {
			result.WriteString(html.EscapeString(ruby.Content))
		}
	}
	return result.String()
}

// processTextWithRuby processes mixed content (text + Ruby elements)
// Note: Due to XML parsing limitations, Ruby elements that were inline in the original
// XML are extracted separately, losing their position. As a workaround, we append them.
func processTextWithRuby(content string, rubies []jplaw.Ruby) string {
	if len(rubies) == 0 {
		return html.EscapeString(content)
	}

	// For now, we just append Ruby elements at the end
	// This is not ideal but the jplaw-xml library doesn't preserve position
	var result strings.Builder
	if content != "" {
		result.WriteString(html.EscapeString(content))
	}

	// Add Ruby elements (they will appear at the end of the text)
	// In the case of "較(こう)正", this will show the Ruby annotation
	result.WriteString(processRubyElements(rubies))
	return result.String()
}
