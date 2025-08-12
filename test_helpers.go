package main

import "go.ngs.io/jplaw-xml"

// createTestSentence creates a Sentence with proper MixedContent for testing
func createTestSentence(content string) jplaw.Sentence {
	return jplaw.Sentence{
		Content: content,
		Ruby:    []jplaw.Ruby{},
		MixedContent: jplaw.MixedContent{
			Nodes: []jplaw.ContentNode{
				jplaw.TextNode{Text: content},
			},
		},
	}
}