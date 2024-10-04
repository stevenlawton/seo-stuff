package improvementchain

import (
	"net/url"
	"sea-stuff/models"
	"strings"
)

// BreadcrumbValidationHandler checks if breadcrumbs are present and valid
type BreadcrumbValidationHandler struct {
	BaseHandler
}

func (h *BreadcrumbValidationHandler) Handle(version *models.ExtractVersion, improvements *[]models.Improvement) {
	if len(version.Breadcrumbs) == 0 {
		*improvements = append(*improvements, models.Improvement{
			Name:     "Missing Breadcrumbs",
			Field:    "Breadcrumbs",
			OldValue: "No breadcrumbs found",
			NewValue: "Add breadcrumbs to improve navigation and user experience",
			Status:   "Pending",
		})
	} else {
		for _, breadcrumb := range version.Breadcrumbs {
			parts := strings.Split(breadcrumb, "|") // Assuming breadcrumb format is "Label|URL"
			if len(parts) != 2 {
				*improvements = append(*improvements, models.Improvement{
					Name:     "Invalid Breadcrumb Format",
					Field:    "Breadcrumbs",
					OldValue: breadcrumb,
					NewValue: "Ensure each breadcrumb has a label and a URL, separated by a '|' character",
					Status:   "Pending",
				})
				continue
			}

			label := strings.TrimSpace(parts[0])
			urlStr := strings.TrimSpace(parts[1])

			if label == "" {
				*improvements = append(*improvements, models.Improvement{
					Name:     "Empty Breadcrumb Label",
					Field:    "Breadcrumbs",
					OldValue: breadcrumb,
					NewValue: "Provide a meaningful label for each breadcrumb",
					Status:   "Pending",
				})
			}

			if urlStr == "" {
				*improvements = append(*improvements, models.Improvement{
					Name:     "Empty Breadcrumb URL",
					Field:    "Breadcrumbs",
					OldValue: breadcrumb,
					NewValue: "Provide a valid URL for each breadcrumb",
					Status:   "Pending",
				})
			} else {
				if _, err := url.ParseRequestURI(urlStr); err != nil {
					*improvements = append(*improvements, models.Improvement{
						Name:     "Invalid Breadcrumb URL",
						Field:    "Breadcrumbs",
						OldValue: urlStr,
						NewValue: "Provide a valid URL for each breadcrumb",
						Status:   "Pending",
					})
				}
			}
		}
	}
	h.CallNext(version, improvements)
}
