package improvementchain

import "sea-stuff/models"

// H1TagCountHandler checks if the number of H1 tags is correct (ideally 1)
type H1TagCountHandler struct {
	next Handler
}

func (h *H1TagCountHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *H1TagCountHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	if version.H1TagCount > 1 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Multiple H1 Tags",
			Field:    "HTags.h1",
			OldValue: "Multiple H1 tags found",
			NewValue: "Only one H1 tag should be used per page",
			Status:   "pending",
		})
	} else if version.H1TagCount == 0 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Missing H1 Tag",
			Field:    "HTags.h1",
			OldValue: "No H1 tags found",
			NewValue: "Add a single H1 tag to the page",
			Status:   "pending",
		})
	}

	if h.next != nil {
		h.next.Handle(version, improvements)
	}
}
