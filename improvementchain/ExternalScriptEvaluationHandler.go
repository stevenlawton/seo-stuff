package improvementchain

import (
	"fmt"
	"sea-stuff/models"
)

// ExternalScriptEvaluationHandler checks if external scripts are slowing down the page
type ExternalScriptEvaluationHandler struct {
	next Handler
}

func (h *ExternalScriptEvaluationHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *ExternalScriptEvaluationHandler) Handle(page *models.AnalysisData, improvements *[]models.Improvement) {
	if len(page.ExternalScripts) > 0 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "External Scripts Detected",
			Field:    "ExternalScripts",
			OldValue: "External scripts found: " + fmt.Sprintf("%d", len(page.ExternalScripts)),
			NewValue: "Consider reducing the number of external scripts or loading them asynchronously",
			Status:   "pending",
		})
	}

	if h.next != nil {
		h.next.Handle(page, improvements)
	}
}
