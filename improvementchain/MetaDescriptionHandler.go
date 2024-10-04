package improvementchain

import (
	"fmt"
	"sea-stuff/models"
	"strings"
)

// MetaDescriptionHandler checks the meta description length
type MetaDescriptionHandler struct {
	BaseHandler
}

func (h *MetaDescriptionHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	if version.MetaDescriptionLength < 50 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Meta Description Too Short",
			Field:    "MetaDescription",
			OldValue: version.MetaDescription,
			NewValue: "Meta description should be between 50 and 160 characters",
			Status:   "Pending",
		})
	} else if version.MetaDescriptionLength > 160 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Meta Description Too Long",
			Field:    "MetaDescription",
			OldValue: version.MetaDescription,
			NewValue: "Meta description should not exceed 160 characters",
			Status:   "Pending",
		})
	}

	mainKeyword := extractMainKeyword(version.Title)
	if mainKeyword != "" && !strings.Contains(strings.ToLower(version.MetaDescription), strings.ToLower(mainKeyword)) {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Keyword Missing in Meta Description",
			Field:    "MetaDescription",
			OldValue: version.MetaDescription,
			NewValue: fmt.Sprintf("Include the keyword '%s' in the meta description", mainKeyword),
			Status:   "Pending",
		})
	}
	h.CallNext(version, improvements)
}
