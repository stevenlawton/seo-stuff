package handlers

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"sea-stuff/improvementchain"
	"sea-stuff/models"
	"sea-stuff/utils"
	"time"
)

// HandleRunScanners runs the improvement chain on the given page
func HandleRunScanners(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to /run_scanners")

	if r.Method != http.MethodPost {
		log.Println("Error: Only POST method is allowed")
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body to get the synthetic key
	var requestData struct {
		SyntheticKey string `json:"SyntheticKey"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Decode the synthetic key using utils.ParseKey
	extractID, url, err := utils.ParseKey(requestData.SyntheticKey)
	if err != nil {
		log.Printf("Error parsing synthetic key: %v", err)
		http.Error(w, "Invalid synthetic key format", http.StatusBadRequest)
		return
	}

	log.Printf("ExtractID: %s, URL: %s", extractID, url)

	if extractID == "" || url == "" {
		log.Println("Error: ExtractID and URL are required")
		http.Error(w, "extractId and url are required", http.StatusBadRequest)
		return
	}

	// Fetch the page data from MongoDB
	collection := client.Database("brandAdherence").Collection("analysis")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var page models.AnalysisData
	filter := bson.M{"extractId": extractID, "url": url}
	err = collection.FindOne(ctx, filter).Decode(&page)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Page not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching page from database", http.StatusInternalServerError)
		}
		return
	}

	// Create and execute the improvement chain
	improvements := []models.Improvement{}

	titleHandler := &improvementchain.TitleLengthHandler{}
	metaDescriptionHandler := &improvementchain.MetaDescriptionHandler{}
	h1TagCountHandler := &improvementchain.H1TagCountHandler{}
	imageAltTextHandler := &improvementchain.ImageAltTextHandler{}
	metaRobotsHandler := &improvementchain.MetaRobotsHandler{}
	pageLoadTimeHandler := &improvementchain.PageLoadTimeHandler{}
	canonicalURLHandler := &improvementchain.CanonicalURLHandler{}

	// Set up the chain
	titleHandler.SetNext(metaDescriptionHandler)
	metaDescriptionHandler.SetNext(h1TagCountHandler)
	h1TagCountHandler.SetNext(imageAltTextHandler)
	imageAltTextHandler.SetNext(metaRobotsHandler)
	metaRobotsHandler.SetNext(pageLoadTimeHandler)
	pageLoadTimeHandler.SetNext(canonicalURLHandler)

	// Execute the chain
	titleHandler.Handle(&page, &improvements)

	// Update the MongoDB document with the generated improvements
	update := bson.M{
		"$set": bson.M{
			"improvements": improvements,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Printf("Error updating page improvements in database: %v", err)
		http.Error(w, "Error updating page improvements in database", http.StatusInternalServerError)
		return
	}

	// Send the improvements back as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(improvements)
}
