package improvementchain

import (
	"encoding/json"
	"fmt"
	"sea-stuff/models"
)

// SocialMetaTagsHandler checks for the presence of social meta tags
type SocialMetaTagsHandler struct {
	next Handler
}

// SetNext sets the next handler in the chain
func (h *SocialMetaTagsHandler) SetNext(handler Handler) {
	h.next = handler
}

// Handle checks if required social meta tags are present
func (h *SocialMetaTagsHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	var socialTags map[string]string
	if err := json.Unmarshal([]byte(version.SocialTags), &socialTags); err != nil {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Invalid Social Tags Format",
			Field:    "SocialTags",
			OldValue: "Social tags could not be parsed",
			NewValue: "Ensure that social tags are provided in a valid JSON format",
			Status:   "pending",
		})
		if h.next != nil {
			h.next.Handle(version, improvements)
		}
		return
	}

	requiredTags := []string{"og:title", "og:description", "twitter:title", "twitter:description"}
	for _, tag := range requiredTags {
		if _, exists := socialTags[tag]; !exists {
			*improvements = append(*improvements, models.Improvement{
				Name:     "Missing Social Meta Tag",
				Field:    "SocialTags",
				OldValue: fmt.Sprintf("Missing %s tag", tag),
				NewValue: fmt.Sprintf("Add %s tag to improve social sharing", tag),
				Status:   "pending",
			})
		}
	}

	if h.next != nil {
		h.next.Handle(version, improvements)
	}
}
