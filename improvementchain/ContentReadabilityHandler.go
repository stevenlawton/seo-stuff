package improvementchain

import (
	"fmt"
	"os"
	"sea-stuff/models"
	"strings"

	"github.com/darkliquid/textstats"
)

// ContentReadabilityHandler checks the readability score of the content
type ContentReadabilityHandler struct {
	BaseHandler
}

func (h *ContentReadabilityHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	readabilityScore := calculateReadabilityScore(version.Content)

	if readabilityScore < 60 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Low Readability Score",
			Field:    "Content",
			OldValue: fmt.Sprintf("Readability score: %.2f", readabilityScore),
			NewValue: "Simplify sentences and use common words to improve readability",
			Status:   "Pending",
		})
	}
	h.CallNext(version, improvements)
}

func calculateReadabilityScore(content string) float64 {
	reader := strings.NewReader(content)
	res, err := textstats.Analyse(reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return res.FleschKincaidReadingEase()
}
