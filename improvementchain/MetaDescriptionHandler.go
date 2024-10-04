package improvementchain

import "sea-stuff/models"

// MetaDescriptionHandler checks the meta description length
type MetaDescriptionHandler struct {
	BaseHandler
}

// Handle processes the meta description and suggests improvements if necessary
func (h *MetaDescriptionHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	if version.MetaDescriptionLength < 50 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Meta Description Too Short",
			Field:    "MetaDescription",
			OldValue: version.MetaDescription,
			NewValue: "Meta description should be between 50 and 160 characters to provide enough context for search engines.",
			Status:   "pending",
		})
	} else if version.MetaDescriptionLength > 160 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Meta Description Too Long",
			Field:    "MetaDescription",
			OldValue: version.MetaDescription,
			NewValue: "Meta description should be between 50 and 160 characters to avoid truncation in search results.",
			Status:   "pending",
		})
	}

	// Call the next handler in the chain, if any
	h.CallNext(version, improvements)
}
