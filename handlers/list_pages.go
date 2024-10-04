package handlers

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"html/template"
	"log"
	"net/http"
	"sea-stuff/models"
	"sea-stuff/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type PageWithKey struct {
	models.AnalysisData
	SyntheticKey string
}

func HandleListPages(w http.ResponseWriter, r *http.Request) {
	collection := client.Database("brandAdherence").Collection("analysis")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch distinct extractId values for the dropdown
	log.Println("Attempting to fetch distinct extractId...")
	extractIDs, err := collection.Distinct(ctx, "extractId", bson.M{})
	if err != nil {
		log.Printf("Error fetching distinct extractIds: %v", err)
		http.Error(w, "Error fetching distinct extractIds", http.StatusInternalServerError)
		return
	}

	var extractIDStrings []string
	for _, id := range extractIDs {
		if idStr, ok := id.(string); ok {
			extractIDStrings = append(extractIDStrings, idStr)
		} else {
			log.Printf("Warning: Non-string extractId found: %v", id)
		}
	}

	if len(extractIDStrings) == 0 {
		log.Println("No extractIds found in the collection. Possible reasons could be no documents or incorrect field name.")
	} else {
		log.Printf("Extracted IDs: %v", extractIDStrings)
	}

	// Prepare the filter based on query parameter
	filter := bson.M{}
	selectedExtractID := r.URL.Query().Get("extractId")
	if selectedExtractID != "" {
		filter = bson.M{"extractId": selectedExtractID}
	}

	// Fetch filtered pages from MongoDB with a projection
	projection := bson.D{{"url", 1}, {"title", 1}, {"extractId", 1}}
	opts := options.Find().SetProjection(projection)

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("Error fetching data from the database: %v", err)
		http.Error(w, "Error fetching data from the database", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err = cursor.Close(ctx); err != nil {
			log.Printf("Error closing cursor: %v", err)
		}
	}()

	// Parse documents into a slice
	var pages []PageWithKey
	for cursor.Next(ctx) {
		var page models.AnalysisData
		if err = cursor.Decode(&page); err != nil {
			log.Printf("Error decoding data from the database: %v", err)
			http.Error(w, "Error decoding data from the database", http.StatusInternalServerError)
			return
		}

		// Ensure that extractId is correctly passed to GenerateKey
		syntheticKey := utils.GenerateKey(page.ExtractID, page.URL)
		if page.ExtractID == "" {
			log.Printf("Warning: extractId is empty for page URL: %s", page.URL)
		}
		pages = append(pages, PageWithKey{
			AnalysisData: page,
			SyntheticKey: syntheticKey,
		})
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Error iterating through the data: %v", err)
		http.Error(w, "Error iterating through the data", http.StatusInternalServerError)
		return
	}

	// Load and execute the HTML template
	tmpl, err := template.ParseFiles("templates/pages.html")
	if err != nil {
		log.Printf("Error loading HTML template: %v", err)
		http.Error(w, "Error loading HTML template", http.StatusInternalServerError)
		return
	}

	// Define a struct to pass the pages and extract IDs to the template
	data := struct {
		Pages             []PageWithKey
		ExtractIDs        []string
		SelectedExtractID string
	}{
		Pages:             pages,
		ExtractIDs:        extractIDStrings,
		SelectedExtractID: selectedExtractID,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error rendering HTML: %v", err)
		http.Error(w, "Error rendering HTML", http.StatusInternalServerError)
	}
}
