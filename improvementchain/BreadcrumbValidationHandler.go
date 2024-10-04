package improvementchain

import "sea-stuff/models"

// BreadcrumbValidationHandler checks if breadcrumbs are present
type BreadcrumbValidationHandler struct {
	next Handler
}

func (h *BreadcrumbValidationHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *BreadcrumbValidationHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	if len(version.Breadcrumbs) == 0 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Missing Breadcrumbs",
			Field:    "Breadcrumbs",
			OldValue: "No breadcrumbs found",
			NewValue: "Add breadcrumbs to improve navigation and user experience",
			Status:   "pending",
		})
	}

	if h.next != nil {
		h.next.Handle(version, improvements)
	}
}
