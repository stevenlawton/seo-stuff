package improvementchain

import "sea-stuff/models"

// MetaDescriptionHandler checks the meta description length
type MetaDescriptionHandler struct {
	BaseHandler
}

// Handle processes the meta description and suggests improvements if necessary
func (h *MetaDescriptionHandler) Handle(page *models.AnalysisData, improvements *[]models.Improvement) {
	if page.MetaDescriptionLength < 50 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Meta Description Too Short",
			Field:    "MetaDescription",
			OldValue: page.MetaDescription,
			NewValue: "Meta description should be between 50 and 160 characters to provide enough context for search engines.",
			Status:   "pending",
		})
	} else if page.MetaDescriptionLength > 160 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Meta Description Too Long",
			Field:    "MetaDescription",
			OldValue: page.MetaDescription,
			NewValue: "Meta description should be between 50 and 160 characters to avoid truncation in search results.",
			Status:   "pending",
		})
	}

	// Call the next handler in the chain, if any
	h.CallNext(page, improvements)
}
