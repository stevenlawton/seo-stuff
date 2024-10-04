package improvementchain

import (
	"fmt"
	"regexp"
	"sea-stuff/models"
	"strings"
)

// KeywordDensityHandler checks if the keyword density is within acceptable limits
type KeywordDensityHandler struct {
	next Handler
}

// SetNext sets the next handler in the chain
func (h *KeywordDensityHandler) SetNext(handler Handler) {
	h.next = handler
}

// Handle calculates keyword density and appends improvement suggestions if necessary
func (h *KeywordDensityHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	// Extract the main keyword from the title by removing common stop words
	mainKeyword := extractMainKeyword(version.Title)
	if mainKeyword == "" {
		if h.next != nil {
			h.next.Handle(version, improvements)
		}
		return
	}

	wordCount := version.WordCount
	if wordCount == 0 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "No Content Found",
			Field:    "WordCount",
			OldValue: "Word count is zero, making keyword density irrelevant",
			NewValue: "Add more content to the page",
			Status:   "pending",
		})
		if h.next != nil {
			h.next.Handle(version, improvements)
		}
		return
	}

	// Calculate how many times the main keyword appears in the content
	keywordCount := 0
	for _, word := range version.CommonWords {
		if isKeywordMatch(mainKeyword, word) {
			keywordCount++
		}
	}

	// Calculate density as a percentage
	keywordDensity := float64(keywordCount) / float64(wordCount) * 100

	// Add improvement suggestions based on the keyword density
	if keywordDensity < 1.0 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Low Keyword Density",
			Field:    "Content",
			OldValue: fmt.Sprintf("Keyword density: %.2f%%", keywordDensity),
			NewValue: "Increase keyword density to between 1% and 3% for effective SEO",
			Status:   "pending",
		})
	} else if keywordDensity > 3.0 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "High Keyword Density",
			Field:    "Content",
			OldValue: fmt.Sprintf("Keyword density: %.2f%%", keywordDensity),
			NewValue: "Reduce keyword density to below 3% to avoid keyword stuffing",
			Status:   "pending",
		})
	}

	if h.next != nil {
		h.next.Handle(version, improvements)
	}
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

// isKeywordMatch checks if the word matches the main keyword
func isKeywordMatch(mainKeyword, word string) bool {
	// Use regex to match the entire word
	regex := fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(mainKeyword))
	matched, _ := regexp.MatchString(regex, strings.ToLower(word))
	return matched
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
