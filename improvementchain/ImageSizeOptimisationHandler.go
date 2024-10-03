package improvementchain

import (
	"sea-stuff/models"
	"strconv"
)

// ImageSizeOptimisationHandler checks if images are too large
type ImageSizeOptimisationHandler struct {
	next Handler
}

func (h *ImageSizeOptimisationHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *ImageSizeOptimisationHandler) Handle(page *models.AnalysisData, improvements *[]models.Improvement) {
	for _, img := range page.Images {
		width, _ := strconv.Atoi(img.Width)
		height, _ := strconv.Atoi(img.Height)

		if width > 1920 || height > 1080 {
			*improvements = append(*improvements, models.Improvement{
				Name:     "Image Too Large",
				Field:    "Images",
				OldValue: "Image dimensions: " + img.Width + "x" + img.Height,
				NewValue: "Reduce image dimensions to 1920x1080 or less for optimal performance",
				Status:   "pending",
			})
		}
	}

	if h.next != nil {
		h.next.Handle(page, improvements)
	}
}
