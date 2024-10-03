package improvementchain

import (
	"fmt"
	"sea-stuff/models"
)

type ContentReadabilityHandler struct {
	next Handler
}

func (h *ContentReadabilityHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *ContentReadabilityHandler) Handle(page *models.AnalysisData, improvements *[]models.Improvement) {
	// Assume we use a readability scoring algorithm (simplified for demonstration)
	readabilityScore := calculateReadabilityScore(page.Content)

	if readabilityScore < 60 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Low Readability Score",
			Field:    "Content",
			OldValue: fmt.Sprintf("Readability score: %d", readabilityScore),
			NewValue: "Rewrite content to improve readability (aim for a score above 60)",
			Status:   "pending",
		})
	}

	if h.next != nil {
		h.next.Handle(page, improvements)
	}
}

func calculateReadabilityScore(content string) int {
	// Simplified placeholder for readability score calculation
	return 50 // Placeholder value
}
