package improvementchain

import "sea-stuff/models"

// ImageAltTextHandler checks if images have alt attributes
type ImageAltTextHandler struct {
	next Handler
}

func (h *ImageAltTextHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *ImageAltTextHandler) Handle(page *models.AnalysisData, improvements *[]models.Improvement) {
	for _, img := range page.Images {
		if img.Alt == "" {
			*improvements = append(*improvements, models.Improvement{
				Name:     "Missing Alt Attribute",
				Field:    "Images.alt",
				OldValue: "No alt attribute for image: " + img.Src,
				NewValue: "Add a descriptive alt attribute for accessibility and SEO",
				Status:   "pending",
			})
		}
	}

	if h.next != nil {
		h.next.Handle(page, improvements)
	}
}
