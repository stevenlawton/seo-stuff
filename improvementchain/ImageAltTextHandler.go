package improvementchain

import (
	"sea-stuff/models"
	"strings"
)

// ImageAltTextHandler checks if all images have meaningful alt text
type ImageAltTextHandler struct {
	BaseHandler
}

// Handle checks if each image has an alt text and suggests improvements if not
func (h *ImageAltTextHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	altTextUsage := make(map[string]int)
	for _, image := range version.Images {
		altText := strings.TrimSpace(image.Alt)
		if altText == "" {
			*improvements = append(*improvements, models.Improvement{
				Name:     "Missing Alt Text",
				Field:    "Image",
				OldValue: image.Src,
				NewValue: "Add descriptive alt text for this image",
				Status:   "Pending",
			})
		} else {
			altTextUsage[altText]++
		}
	}

	// Detect repetitive alt texts
	for altText, count := range altTextUsage {
		if count > 1 {
			*improvements = append(*improvements, models.Improvement{
				Name:     "Repetitive Alt Text",
				Field:    "Image",
				OldValue: altText,
				NewValue: "Use unique alt text for each image to improve accessibility",
				Status:   "Pending",
			})
		}
	}
	h.CallNext(version, improvements)
}
