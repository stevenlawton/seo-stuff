package improvementchain

import "sea-stuff/models"

// TitleLengthHandler checks the length of the page title
type TitleLengthHandler struct {
	BaseHandler
}

func (h *TitleLengthHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	if version.TitleLength > 60 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Title Length Too Long",
			Field:    "Title",
			OldValue: version.Title,
			NewValue: "Title should be 60 characters or less",
			Status:   "pending",
		})
	}
	h.CallNext(version, improvements)
}
