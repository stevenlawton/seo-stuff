package improvementchain

import (
	"fmt"
	"sea-stuff/models"
)

type SocialMetaTagsHandler struct {
	next Handler
}

func (h *SocialMetaTagsHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *SocialMetaTagsHandler) Handle(page *models.AnalysisData, improvements *[]models.Improvement) {
	requiredTags := []string{"og:title", "og:description", "twitter:title", "twitter:description"}
	for _, tag := range requiredTags {
		if _, exists := page.SocialTags[tag]; !exists {
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
		h.next.Handle(page, improvements)
	}
}
