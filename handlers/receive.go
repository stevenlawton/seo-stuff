package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
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

// HandlePost handles incoming POST requests.
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

	if data.ExtractID == "" || data.URL == "" {
		http.Error(w, "Missing required fields: extractid or URL", http.StatusBadRequest)
		return
	}

	// Connect to the MongoDB collection
	collection := client.Database("brandAdherence").Collection("analysis")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insert the data into the collection
	_, err = collection.InsertOne(ctx, bson.M{
		"extractid":             data.ExtractID,
		"url":                   data.URL,
		"title":                 data.Title,
		"titlelength":           data.TitleLength,
		"metadescription":       data.MetaDescription,
		"metadescriptionlength": data.MetaDescriptionLength,
		"metatags":              data.MetaTags,
		"canonicalurl":          data.CanonicalURL,
		"htags":                 data.HTags,
		"h1tagcount":            data.H1TagCount,
		"wordcount":             data.WordCount,
		"pageloadtimeseconds":   data.PageLoadTimeSeconds,
		"images":                data.Images,
		"internallinks":         data.InternalLinks,
		"externallinks":         data.ExternalLinks,
		"brokenlinks":           data.BrokenLinks,
		"structureddata":        data.StructuredData,
		"robotsmetatag":         data.RobotsMetaTag,
		"content":               data.Content,
		"improvements":          data.Improvements,
	})
	if err != nil {
		http.Error(w, "Error saving to the database", http.StatusInternalServerError)
		return
	}
	// Respond with a success message
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Data saved successfully")
}
