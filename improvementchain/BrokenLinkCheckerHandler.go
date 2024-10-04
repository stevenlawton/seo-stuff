package improvementchain

import "sea-stuff/models"

type BrokenLinkCheckerHandler struct {
	next Handler
}

func (h *BrokenLinkCheckerHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *BrokenLinkCheckerHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	for _, brokenLink := range version.BrokenLinks {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Broken Link",
			Field:    "Links",
			OldValue: brokenLink,
			NewValue: "Replace or remove the broken link",
			Status:   "pending",
		})
	}

	if h.next != nil {
		h.next.Handle(version, improvements)
	}
}
