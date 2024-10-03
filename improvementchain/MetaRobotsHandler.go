package improvementchain

import "sea-stuff/models"

// MetaRobotsHandler checks if the meta robots tag is correctly set
type MetaRobotsHandler struct {
	next Handler
}

func (h *MetaRobotsHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *MetaRobotsHandler) Handle(page *models.AnalysisData, improvements *[]models.Improvement) {
	if page.RobotsMetaTag == "" {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Missing Robots Meta Tag",
			Field:    "RobotsMetaTag",
			OldValue: "No robots meta tag found",
			NewValue: "Add a robots meta tag to control indexing (e.g., 'index, follow')",
			Status:   "pending",
		})
	} else if page.RobotsMetaTag != "index, follow" {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Improper Robots Meta Tag",
			Field:    "RobotsMetaTag",
			OldValue: page.RobotsMetaTag,
			NewValue: "Consider changing robots meta tag to 'index, follow' to allow indexing",
			Status:   "pending",
		})
	}

	if h.next != nil {
		h.next.Handle(page, improvements)
	}
}
