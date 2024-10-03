package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sea-stuff/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

// SetClient sets the Mongo client for the handlers.
func SetClient(mongoClient *mongo.Client) {
	client = mongoClient
}

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

	if data.ExtractID == "" {
		http.Error(w, "Missing extract ID", http.StatusBadRequest)
		return
	}

	// Log the extractId before inserting to confirm it's there
	log.Printf("Received extractId: %s", data.ExtractID)

	// Insert the data into the collection with the correct field name
	collection := client.Database("brandAdherence").Collection("analysis")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, data)
	if err != nil {
		http.Error(w, "Error saving to the database", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Data saved successfully")
}
