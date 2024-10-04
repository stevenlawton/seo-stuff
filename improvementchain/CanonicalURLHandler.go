package improvementchain

import (
	"net/url"
	"sea-stuff/models"
)

// CanonicalURLHandler checks if the canonical URL is correct
type CanonicalURLHandler struct {
	BaseHandler
}

func (h *CanonicalURLHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	if !version.IsCanonicalCorrect {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Incorrect Canonical URL",
			Field:    "CanonicalURL",
			OldValue: version.CanonicalURL,
			NewValue: "Update the canonical URL to match the current page URL",
			Status:   "Pending",
		})
	}

	//// Check for multiple canonical tags (assuming this info is in version.CanonicalTagCount)
	//if version.CanonicalTagCount > 1 {
	//	*improvements = append(*improvements, models.Improvement{
	//		Name:     "Multiple Canonical Tags",
	//		Field:    "CanonicalURL",
	//		OldValue: fmt.Sprintf("%d canonical tags found", version.CanonicalTagCount),
	//		NewValue: "Ensure only one canonical tag is present on the page",
	//		Status:   "Pending",
	//	})
	//}

	// Compare canonical URL to actual URL
	if parsedCanonical, err := url.Parse(version.CanonicalURL); err == nil {
		if parsedPageURL, err := url.Parse(version.URL); err == nil {
			if parsedCanonical.Path != parsedPageURL.Path {
				*improvements = append(*improvements, models.Improvement{
					Name:     "Canonical URL Mismatch",
					Field:    "CanonicalURL",
					OldValue: version.CanonicalURL,
					NewValue: "Canonical URL should match the page URL unless intentionally different",
					Status:   "Pending",
				})
			}
		}
	}
	h.CallNext(version, improvements)
}
