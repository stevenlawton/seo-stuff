package handlers

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// HandleDeleteByExtractID deletes the version and improvement documents associated with a given extractId.
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set up the MongoDB collections
	versionCollection := client.Database("brandAdherence").Collection("versions")
	improvementCollection := client.Database("brandAdherence").Collection("improvements")

	// Delete all versions with the given extract ID
	_, err := versionCollection.DeleteMany(ctx, bson.M{"extractId": extractID})
	if err != nil {
		http.Error(w, "Error deleting versions from the database", http.StatusInternalServerError)
		return
	}

	// Delete all improvements associated with the given extract ID
	_, err = improvementCollection.DeleteMany(ctx, bson.M{"extractId": extractID})
	if err != nil {
		http.Error(w, "Error deleting improvements from the database", http.StatusInternalServerError)
		return
	}

	// Redirect back to the pages list
	http.Redirect(w, r, "/pages", http.StatusSeeOther)
}
