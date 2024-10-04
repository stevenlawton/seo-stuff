package improvementchain

import "sea-stuff/models"

type StructuredDataValidationHandler struct {
	next Handler
}

func (h *StructuredDataValidationHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *StructuredDataValidationHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	if len(version.StructuredData) == 0 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Missing Structured Data",
			Field:    "StructuredData",
			OldValue: "No structured data found",
			NewValue: "Add structured data (e.g., JSON-LD) to improve search engine understanding",
			Status:   "pending",
		})
	} else {
		// Additional validation logic for structured data could be added here
		for _, data := range version.StructuredDataTypes {
			if data == "" {
				*improvements = append(*improvements, models.Improvement{
					Name:     "Invalid Structured Data Type",
					Field:    "StructuredData",
					OldValue: "Empty or incorrect structured data type",
					NewValue: "Ensure the structured data type is valid (e.g., Article, Product, etc.)",
					Status:   "pending",
				})
			}
		}
	}

	if h.next != nil {
		h.next.Handle(version, improvements)
	}
}
