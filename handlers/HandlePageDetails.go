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

	// Parse the synthetic key to get extractID and URL
	extractID, urlStr, err := utils.ParseKey(syntheticKey)
	if err != nil {
		http.Error(w, "Invalid page key", http.StatusBadRequest)
		return
	}

	// Log the parsed values for debugging
	log.Printf("Parsed extractId: %s, URL: %s", extractID, urlStr)

	// Connect to the MongoDB collection
	collection := client.Database("brandAdherence").Collection("pages")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch the page document by URL and find the specific version using the extractID
	var page models.Page
	filter := bson.M{"url": urlStr}
	projection := bson.M{
		"_id": 0,
		"url": 1,
		"versions": bson.M{
			"$elemMatch": bson.M{"extractId": extractID},
		},
	}

	err = collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&page)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Page not found for extractId: %s, URL: %s", extractID, urlStr)
			http.Error(w, "Page not found", http.StatusNotFound)
		} else {
			log.Printf("Error fetching page details from the database: %v", err)
			http.Error(w, "Error fetching page details from the database", http.StatusInternalServerError)
		}
		return
	}

	if len(page.Versions) == 0 {
		log.Println("Error: Specified version not found")
		http.Error(w, "Specified version not found", http.StatusNotFound)
		return
	}

	pageVersion := &page.Versions[0]

	// Load and execute the HTML template
	tmpl, err := template.ParseFiles("templates/page_detail.html")
	if err != nil {
		log.Printf("Error loading HTML template: %v", err)
		http.Error(w, "Error loading HTML template", http.StatusInternalServerError)
		return
	}

	data := struct {
		PageURL string
		Version *models.ExtractVersion
	}{
		PageURL: page.URL,
		Version: pageVersion,
	}

	// Serve the template with the page version data
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error rendering HTML: %v", err)
		http.Error(w, "Error rendering HTML", http.StatusInternalServerError)
	}
}
