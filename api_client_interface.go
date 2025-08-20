package jplaw2epub

import lawapi "go.ngs.io/jplaw-api-v2"

// APIClient defines the interface for the law API client
type APIClient interface {
	GetAttachment(lawRevisionID string, params *lawapi.GetAttachmentParams) (*string, error)
}

// Ensure lawapi.Client implements APIClient
var _ APIClient = (*lawapi.Client)(nil)
