package improvementchain

import "sea-stuff/models"

// Handler interface for the chain of responsibility pattern
type Handler interface {
	Handle(page *models.AnalysisData, improvements *[]models.Improvement)
	SetNext(handler Handler)
}

// BaseHandler struct to handle chaining
type BaseHandler struct {
	next Handler
}

// SetNext sets the next handler in the chain
func (h *BaseHandler) SetNext(handler Handler) {
	h.next = handler
}

// CallNext moves to the next handler, if it exists
func (h *BaseHandler) CallNext(page *models.AnalysisData, improvements *[]models.Improvement) {
	if h.next != nil {
		h.next.Handle(page, improvements)
	}
}
