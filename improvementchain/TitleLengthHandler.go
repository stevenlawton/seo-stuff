package improvementchain

import "sea-stuff/models"

// TitleLengthHandler checks the length of the page title
type TitleLengthHandler struct {
	BaseHandler
}

func (h *TitleLengthHandler) Handle(page *models.AnalysisData, improvements *[]models.Improvement) {
	if page.TitleLength > 60 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Title Length Too Long",
			Field:    "Title",
			OldValue: page.Title,
			NewValue: "Title should be 60 characters or less",
			Status:   "pending",
		})
	}
	h.CallNext(page, improvements)
}
