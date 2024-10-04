package improvementchain

import "sea-stuff/models"

// MobileFriendlinessHandler checks if the page is mobile-friendly
type MobileFriendlinessHandler struct {
	next Handler
}

func (h *MobileFriendlinessHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *MobileFriendlinessHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	if !version.IsMobileFriendly {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Not Mobile-Friendly",
			Field:    "IsMobileFriendly",
			OldValue: "The page is not optimized for mobile devices",
			NewValue: "Ensure that the page is responsive and adjusts correctly on mobile devices",
			Status:   "pending",
		})
	}

	if h.next != nil {
		h.next.Handle(version, improvements)
	}
}
