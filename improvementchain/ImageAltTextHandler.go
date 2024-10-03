package improvementchain

import (
	"sea-stuff/models"
)

// ImageAltTextHandler checks if all images have meaningful alt text
type ImageAltTextHandler struct {
	next Handler
}

// SetNext sets the next handler in the chain
func (h *ImageAltTextHandler) SetNext(handler Handler) {
	h.next = handler
}

// Handle checks if each image has an alt text and suggests improvements if not
func (h *ImageAltTextHandler) Handle(page *models.AnalysisData, improvements *[]models.Improvement) {
	for _, image := range page.Images {
		if image.Alt == "" {
			*improvements = append(*improvements, models.Improvement{
				Name:     "Missing Alt Text",
				Field:    "Image",
				OldValue: image.Src,
				NewValue: "Add descriptive alt text for this image",
				Status:   "pending",
			})
		}
	}

	if h.next != nil {
		h.next.Handle(page, improvements)
	}
}
