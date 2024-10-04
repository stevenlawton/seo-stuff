package improvementchain

import "sea-stuff/models"

// MetaRobotsHandler checks if the meta robots tag is correctly set
type MetaRobotsHandler struct {
	BaseHandler
}

func (h *MetaRobotsHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	if version.RobotsMetaTag == "" {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Missing Robots Meta Tag",
			Field:    "RobotsMetaTag",
			OldValue: "No robots meta tag found",
			NewValue: "Add a robots meta tag to control indexing (e.g., 'index, follow')",
			Status:   "Pending",
		})
	} else if version.RobotsMetaTag != "index, follow" {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Improper Robots Meta Tag",
			Field:    "RobotsMetaTag",
			OldValue: version.RobotsMetaTag,
			NewValue: "Consider changing robots meta tag to 'index, follow' to allow indexing",
			Status:   "Pending",
		})
	}

	h.CallNext(version, improvements)
}
