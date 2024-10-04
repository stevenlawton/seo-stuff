package handlers

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"net/http"
	"net/url"
	"sea-stuff/models"
	"sea-stuff/utils"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type PageWithKey struct {
	URL          string
	SyntheticKey string
	Title        string
	Improvements []models.Improvement
	ExtractID    string
}

func HandleListPages(w http.ResponseWriter, r *http.Request) {
	pages, err := getAllPages()
	if err != nil {
		http.Error(w, "Error retrieving pages", http.StatusInternalServerError)
		return
	}

	// Extract unique domains
	domainSet := make(map[string]struct{})
	for _, page := range pages {
		domain := getDomainFromURL(page.URL)
		if domain != "" {
			domainSet[domain] = struct{}{}
		}
	}
	var domains []string
	for domain := range domainSet {
		domains = append(domains, domain)
	}
	sort.Strings(domains)

	// Get selected domain from query parameters
	selectedDomain := r.URL.Query().Get("domain")

	// Filter pages based on selected domain
	var filteredPages []PageWithKey
	if selectedDomain != "" {
		for _, page := range pages {
			if getDomainFromURL(page.URL) == selectedDomain {
				// For each version in page.Versions, create a PageWithKey
				for _, version := range page.Versions {
					syntheticKey := utils.GenerateKey(version.ExtractID, page.URL)
					filteredPages = append(filteredPages, PageWithKey{
						URL:          page.URL,
						SyntheticKey: syntheticKey,
						Title:        version.Title,
						Improvements: version.Improvements,
						ExtractID:    version.ExtractID,
					})
				}
			}
		}
	} else {
		// If no domain is selected, show all pages
		for _, page := range pages {
			for _, version := range page.Versions {
				syntheticKey := utils.GenerateKey(version.ExtractID, page.URL)
				filteredPages = append(filteredPages, PageWithKey{
					URL:          page.URL,
					SyntheticKey: syntheticKey,
					Title:        version.Title,
					Improvements: version.Improvements,
					ExtractID:    version.ExtractID,
				})
			}
		}
	}

	data := struct {
		Domains           []string
		SelectedDomain    string
		Pages             []PageWithKey
		SelectedExtractID string
	}{
		Domains:           domains,
		SelectedDomain:    selectedDomain,
		Pages:             filteredPages,
		SelectedExtractID: "", // Modify as needed
	}

	tmpl, err := template.ParseFiles("templates/pages.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func getAllPages() ([]models.Page, error) {
	collection, err := getMongoCollection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Fetch all pages with their versions
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pages []models.Page
	if err = cursor.All(ctx, &pages); err != nil {
		return nil, err
	}

	return pages, nil
}

func getDomainFromURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

func getMongoCollection() (*mongo.Collection, error) {
	return client.Database("brandAdherence").Collection("pages"), nil
}
