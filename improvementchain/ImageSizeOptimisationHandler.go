package improvementchain

import (
	"sea-stuff/models"
	"strconv"
)

// ImageSizeOptimisationHandler checks if images are too large
type ImageSizeOptimisationHandler struct {
	BaseHandler
}

func (h *ImageSizeOptimisationHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	for _, img := range version.Images {
		if img.Width > 1920 || img.Height > 1080 {
			*improvements = append(*improvements, models.Improvement{
				Name:     "Image Too Large",
				Field:    "Images",
				OldValue: "Image dimensions: " + strconv.Itoa(img.Width) + "x" + strconv.Itoa(img.Height),
				NewValue: "Reduce image dimensions to 1920x1080 or less for optimal performance",
				Status:   "Pending", // Standardized capitalization
			})
		}
	}

	h.CallNext(version, improvements)
}
