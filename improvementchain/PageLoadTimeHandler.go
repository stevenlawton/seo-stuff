package improvementchain

import (
	"fmt"
	"sea-stuff/models"
)

// PageLoadTimeHandler checks if the page load time is within acceptable limits
type PageLoadTimeHandler struct {
	BaseHandler
}

func (h *PageLoadTimeHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	if version.PageLoadTimeSeconds > 3.0 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Page Load Time Too High",
			Field:    "PageLoadTimeSeconds",
			OldValue: "Page load time: " + fmt.Sprintf("%.2f", version.PageLoadTimeSeconds) + " seconds",
			NewValue: "Reduce page load time to under 3 seconds",
			Status:   "Pending", // Standardized capitalization
		})
	}

	h.CallNext(version, improvements)
}
