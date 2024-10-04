package improvementchain

import (
	"fmt"
	"sea-stuff/models"
)

// InternalLinkDepthHandler checks if the page depth is within acceptable limits
type InternalLinkDepthHandler struct {
	BaseHandler
}

func (h *InternalLinkDepthHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	if version.PageDepth > 3 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Page Depth Too High",
			Field:    "PageDepth",
			OldValue: "Page depth: " + fmt.Sprintf("%d", version.PageDepth),
			NewValue: "Reduce page depth to improve crawl efficiency (aim for a depth of 3 or less)",
			Status:   "Pending", // Standardized capitalization
		})
	}

	h.CallNext(version, improvements)
}
