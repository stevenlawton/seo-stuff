package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sea-stuff/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// SetClient sets the Mongo client for the handlers.
func SetClient(mongoClient *mongo.Client) {
	client = mongoClient
}

// HandlePost handles POST requests to insert or update analysis data.
func HandlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var data models.AnalysisData
	// Decode the JSON payload
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate essential fields
	if data.ExtractID == "" {
		http.Error(w, "Missing extract ID", http.StatusBadRequest)
		return
	}
	if data.URL == "" {
		http.Error(w, "Missing URL", http.StatusBadRequest)
		return
	}
	if _, err := url.ParseRequestURI(data.URL); err != nil {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	// Log the extractId before inserting to confirm it's there
	log.Printf("Received ExtractID: %s, URL: %s", data.ExtractID, data.URL)

	// Insert or update the data into the collection with upsert
	collection := client.Database("brandAdherence").Collection("analysis")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"extractId": data.ExtractID, "url": data.URL}
	update := bson.M{"$set": data}
	opts := options.Update().SetUpsert(true)

	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Printf("Error saving to the database: %v", err)
		http.Error(w, "Error saving to the database", http.StatusInternalServerError)
		return
	}

	// Log result of upsert
	if result.MatchedCount > 0 {
		log.Printf("Updated existing document for ExtractID: %s, URL: %s", data.ExtractID, data.URL)
	} else if result.UpsertedCount > 0 {
		log.Printf("Inserted new document for ExtractID: %s, URL: %s", data.ExtractID, data.URL)
	}

	// Respond with a success message
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Data saved successfully")
}
