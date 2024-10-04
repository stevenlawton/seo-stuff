package improvementchain

import (
	"fmt"
	"sea-stuff/models"
	"strings"
)

// H1TagCountHandler checks if the number of H1 tags is correct (ideally 1)
type H1TagCountHandler struct {
	BaseHandler
}

func (h *H1TagCountHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	h1Tags := version.HTags["h1"]

	if len(h1Tags) == 0 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Missing H1 Tag",
			Field:    "HTags.h1",
			OldValue: "No H1 tags found",
			NewValue: "Add a single H1 tag to the page",
			Status:   "Pending",
		})
	} else if len(h1Tags) > 1 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Multiple H1 Tags",
			Field:    "HTags.h1",
			OldValue: fmt.Sprintf("%d H1 tags found", len(h1Tags)),
			NewValue: "Use only one H1 tag per page",
			Status:   "Pending",
		})
	} else {
		h1Content := strings.TrimSpace(h1Tags[0])
		if h1Content == "" {
			*improvements = append(*improvements, models.Improvement{
				Name:     "Empty H1 Tag",
				Field:    "HTags.h1",
				OldValue: "H1 tag is empty",
				NewValue: "Provide descriptive content in the H1 tag",
				Status:   "Pending",
			})
		} else if len(h1Content) > 70 {
			*improvements = append(*improvements, models.Improvement{
				Name:     "Long H1 Tag",
				Field:    "HTags.h1",
				OldValue: h1Content,
				NewValue: "Shorten the H1 tag to 70 characters or fewer",
				Status:   "Pending",
			})
		}
	}
	h.CallNext(version, improvements)
}
