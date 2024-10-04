package improvementchain

import (
	"fmt"
	"sea-stuff/models"
	"sea-stuff/utils"
)

type TitleLengthHandler struct {
	BaseHandler
}

func (h *TitleLengthHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	if version.TitleLength > 60 || version.TitleLength < 10 {
		// Prepare the prompt for GPT-4
		prompt := fmt.Sprintf(`Our brand guidelines are as follows: %s

The page content is: %s

The current title is: "%s"

The title is %d characters long, which is outside the recommended length.

Please suggest a new title that aligns with our brand guidelines and is between 50 to 60 characters.`, version.BrandJSON, version.Content, version.Title, version.TitleLength)

		// Call GPT-4 to get a suggestion
		suggestion, err := utils.CallGPT4(prompt)
		if err != nil {
			suggestion = "Unable to generate suggestion at this time."
		}

		*improvements = append(*improvements, models.Improvement{
			Name:     "Title Length Improvement",
			Field:    "Title",
			OldValue: version.Title,
			NewValue: suggestion,
			Status:   "Pending",
		})
	}
	h.CallNext(version, improvements)
}
