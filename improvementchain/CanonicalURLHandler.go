package improvementchain

import "sea-stuff/models"

// CanonicalURLHandler checks if the canonical URL is correct
type CanonicalURLHandler struct {
	next Handler
}

func (h *CanonicalURLHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *CanonicalURLHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	if !version.IsCanonicalCorrect {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Incorrect Canonical URL",
			Field:    "CanonicalURL",
			OldValue: version.CanonicalURL,
			NewValue: "Update the canonical URL to match the current page URL",
			Status:   "pending",
		})
	}

	if h.next != nil {
		h.next.Handle(version, improvements)
	}
}
