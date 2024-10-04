package improvementchain

import "sea-stuff/models"

// MobileFriendlinessHandler checks if the page is mobile-friendly
type MobileFriendlinessHandler struct {
	BaseHandler
}

func (h *MobileFriendlinessHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	if !version.IsMobileFriendly {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Not Mobile-Friendly",
			Field:    "IsMobileFriendly",
			OldValue: "The page is not optimized for mobile devices",
			NewValue: "Ensure that the page is responsive and adjusts correctly on mobile devices",
			Status:   "Pending", // Standardized capitalization
		})
	}

	h.CallNext(version, improvements)
}
