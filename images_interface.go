package jplaw2epub

import "go.ngs.io/jplaw-xml"

// ImageProcessorInterface defines the interface for image processing
type ImageProcessorInterface interface {
	ProcessFigStruct(fig *jplaw.FigStruct) (string, error)
	SetMaxImageHeight(height string)
}

// Ensure ImageProcessor implements ImageProcessorInterface
var _ ImageProcessorInterface = (*ImageProcessor)(nil)
