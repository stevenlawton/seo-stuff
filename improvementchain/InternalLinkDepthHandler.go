package improvementchain

import (
	"fmt"
	"sea-stuff/models"
)

// InternalLinkDepthHandler checks if the page depth is within acceptable limits
type InternalLinkDepthHandler struct {
	next Handler
}

func (h *InternalLinkDepthHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *InternalLinkDepthHandler) Handle(page *models.AnalysisData, improvements *[]models.Improvement) {
	if page.PageDepth > 3 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Page Depth Too High",
			Field:    "PageDepth",
			OldValue: "Page depth: " + fmt.Sprintf("%d", page.PageDepth),
			NewValue: "Reduce page depth to improve crawl efficiency (aim for a depth of 3 or less)",
			Status:   "pending",
		})
	}

	if h.next != nil {
		h.next.Handle(page, improvements)
	}
}
