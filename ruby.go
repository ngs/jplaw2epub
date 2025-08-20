package jplaw2epub

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
// XML are extracted separately, losing their position. As a workaround, we append them
// at the end of the text.
func processTextWithRuby(content string, rubies []jplaw.Ruby) string {
	if len(rubies) == 0 {
		return html.EscapeString(content)
	}

	// Build result with escaped content first
	result := html.EscapeString(content)

	// Append all ruby elements at the end
	// This is the expected behavior based on the jplaw-xml library's structure
	for _, ruby := range rubies {
		result += rubyHTML(&ruby)
	}

	return result
}

// rubyHTML creates HTML ruby element from a Ruby struct
func rubyHTML(ruby *jplaw.Ruby) string {
	if len(ruby.Rt) == 0 {
		return html.EscapeString(ruby.Content)
	}

	var result strings.Builder
	result.WriteString("<ruby>")
	result.WriteString(html.EscapeString(ruby.Content))
	for _, rt := range ruby.Rt {
		result.WriteString("<rt>")
		result.WriteString(html.EscapeString(rt.Content))
		result.WriteString("</rt>")
	}
	result.WriteString("</ruby>")
	return result.String()
}
