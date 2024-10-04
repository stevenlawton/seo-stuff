package handlers

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"sea-stuff/models"
	"sea-stuff/utils"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func HandlePageDetails(w http.ResponseWriter, r *http.Request) {
	// Extract the synthetic key from the URL
	syntheticKey := strings.TrimPrefix(r.URL.Path, "/pages/")
	if syntheticKey == "" {
		http.Error(w, "Page key is required", http.StatusBadRequest)
		return
	}

	// Parse the synthetic key to get extract_id and URL
	extractID, url, err := utils.ParseKey(syntheticKey)
	if err != nil {
		http.Error(w, "Invalid page key", http.StatusBadRequest)
		return
	}

	// Log the parsed values for debugging
	log.Printf("Parsed extractId: %s, URL: %s", extractID, url)

	// Connect to the MongoDB collection
	collection := client.Database("brandAdherence").Collection("analysis")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch the page document by extractId and URL using projection to limit fields
	var page models.AnalysisData
	filter := bson.M{"extractId": extractID, "url": url}
	projection := bson.M{"_id": 0} // Exclude internal MongoDB ID field
	err = collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&page)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Page not found for extractId: %s, URL: %s", extractID, url)
			http.Error(w, "Page not found", http.StatusNotFound)
		} else {
			log.Printf("Error fetching page details from the database: %v", err)
			http.Error(w, "Error fetching page details from the database", http.StatusInternalServerError)
		}
		return
	}

	// Load and execute the HTML template
	tmpl, err := template.ParseFiles("templates/page_detail.html")
	if err != nil {
		log.Printf("Error loading HTML template: %v", err)
		http.Error(w, "Error loading HTML template", http.StatusInternalServerError)
		return
	}

	// Serve the template with the page data
	if err := tmpl.Execute(w, page); err != nil {
		log.Printf("Error rendering HTML: %v", err)
		http.Error(w, "Error rendering HTML", http.StatusInternalServerError)
	}
}
