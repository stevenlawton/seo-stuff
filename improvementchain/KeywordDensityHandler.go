package improvementchain

import (
	"fmt"
	"sea-stuff/models"
	"strings"

	"github.com/reiver/go-porterstemmer"
)

// KeywordDensityHandler checks if the keyword density is within acceptable limits
type KeywordDensityHandler struct {
	BaseHandler
}

func (h *KeywordDensityHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	mainKeyword := extractMainKeyword(version.Title)
	if mainKeyword == "" {
		h.CallNext(version, improvements)
		return
	}

	wordCount := version.WordCount
	if wordCount == 0 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "No Content Found",
			Field:    "WordCount",
			OldValue: "Word count is zero",
			NewValue: "Add more content to the page",
			Status:   "Pending",
		})
		h.CallNext(version, improvements)
		return
	}

	keywordStem := porterstemmer.StemString(strings.ToLower(mainKeyword))
	keywordCount := 0
	for _, word := range version.CommonWords {
		wordStem := porterstemmer.StemString(strings.ToLower(word))
		if keywordStem == wordStem {
			keywordCount++
		}
	}

	keywordDensity := float64(keywordCount) / float64(wordCount) * 100

	if keywordDensity < 1.0 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Low Keyword Density",
			Field:    "Content",
			OldValue: fmt.Sprintf("Keyword density: %.2f%%", keywordDensity),
			NewValue: fmt.Sprintf("Consider adding the keyword '%s' more frequently", mainKeyword),
			Status:   "Pending",
		})
	} else if keywordDensity > 3.0 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "High Keyword Density",
			Field:    "Content",
			OldValue: fmt.Sprintf("Keyword density: %.2f%%", keywordDensity),
			NewValue: "Reduce keyword usage to avoid keyword stuffing",
			Status:   "Pending",
		})
	}
	h.CallNext(version, improvements)
}

// extractMainKeyword extracts a potential keyword from the title by removing stop words
func extractMainKeyword(title string) string {
	// List of common stop words to exclude
	stopWords := []string{"the", "is", "in", "at", "of", "and", "a", "to", "for", "on", "by"}
	words := strings.Fields(strings.ToLower(title))
	for _, word := range words {
		if !contains(stopWords, word) {
			return word // return the first non-stop word as the main keyword
		}
	}
	return ""
}

// contains checks if a slice contains a given string
func contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}
