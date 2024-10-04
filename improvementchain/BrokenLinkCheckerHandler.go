package improvementchain

import (
	"fmt"
	"net/http"
	"sea-stuff/models"
)

// BrokenLinkCheckerHandler checks for broken links and provides detailed information
type BrokenLinkCheckerHandler struct {
	BaseHandler
}

func (h *BrokenLinkCheckerHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	if len(version.BrokenLinks) > 0 {
		details := ""
		for _, link := range version.BrokenLinks {
			statusCode := h.getHTTPStatusCode(link)
			details += fmt.Sprintf("- %s (Status: %d)\n", link, statusCode)
		}
		*improvements = append(*improvements, models.Improvement{
			Name:     "Broken Links Found",
			Field:    "Links",
			OldValue: details,
			NewValue: "Replace or remove the broken links listed above",
			Status:   "Pending",
		})
	}
	h.CallNext(version, improvements)
}

func (h *BrokenLinkCheckerHandler) getHTTPStatusCode(link string) int {
	resp, err := http.Head(link)
	if err != nil {
		return 0 // Unable to retrieve status code
	}
	defer resp.Body.Close()
	return resp.StatusCode
}
