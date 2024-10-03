package improvementchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url" // Correct package for parsing URLs
	"sea-stuff/models"
)

type ExternalLinkQualityHandler struct {
	next   Handler
	apiKey string
}

func NewExternalLinkQualityHandler(apiKey string) *ExternalLinkQualityHandler {
	return &ExternalLinkQualityHandler{apiKey: apiKey}
}

func (h *ExternalLinkQualityHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *ExternalLinkQualityHandler) Handle(page *models.AnalysisData, improvements *[]models.Improvement) {
	for _, link := range page.ExternalLinks {
		domain := extractDomain(link)
		if domain == "" {
			continue
		}

		// Make VirusTotal API call to get domain reputation
		if isMalicious, err := h.checkDomainWithVirusTotal(domain); err == nil && isMalicious {
			*improvements = append(*improvements, models.Improvement{
				Name:     "Potentially Malicious External Link",
				Field:    "ExternalLinks",
				OldValue: link,
				NewValue: "Replace with a link to a more authoritative or safer source",
				Status:   "pending",
			})
		}
	}

	if h.next != nil {
		h.next.Handle(page, improvements)
	}
}

func (h *ExternalLinkQualityHandler) checkDomainWithVirusTotal(domain string) (bool, error) {
	url := fmt.Sprintf("https://www.virustotal.com/api/v3/domains/%s", domain)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("x-apikey", h.apiKey)
	req.Header.Set("accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("VirusTotal API responded with status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var response struct {
		Data struct {
			Attributes struct {
				LastAnalysisStats struct {
					Malicious int `json:"malicious"`
				} `json:"last_analysis_stats"`
			} `json:"attributes"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, err
	}

	// If the domain has been marked as malicious by any source
	if response.Data.Attributes.LastAnalysisStats.Malicious > 0 {
		return true, nil
	}

	return false, nil
}

// Helper function to extract domain from URL
func extractDomain(link string) string {
	parsedURL, err := url.Parse(link) // Corrected function call
	if err != nil {
		return ""
	}
	return parsedURL.Host
}
