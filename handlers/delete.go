package handlers

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func HandleDeleteByExtractID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the extractId from the form data
	extractID := r.FormValue("extractId")
	if extractID == "" {
		http.Error(w, "Missing extract ID", http.StatusBadRequest)
		return
	}

	// Set up the MongoDB collection
	collection := client.Database("brandAdherence").Collection("analysis")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete all documents with the given extract ID
	_, err := collection.DeleteMany(ctx, bson.M{"extractId": extractID})
	if err != nil {
		http.Error(w, "Error deleting documents from the database", http.StatusInternalServerError)
		return
	}

	// Redirect back to the pages list
	http.Redirect(w, r, "/pages", http.StatusSeeOther)
}
