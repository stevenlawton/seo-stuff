package handlers

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"os"
	"sea-stuff/improvementchain"
	"sea-stuff/models"
	"sea-stuff/utils"
	"time"
)

// HandleRunScanners runs the improvement chain on the given page version
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
	collection := client.Database("brandAdherence").Collection("pages")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var page models.Page
	filter := bson.M{"url": url, "versions.extractId": extractID}
	err = collection.FindOne(ctx, filter).Decode(&page)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Page not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching page from database", http.StatusInternalServerError)
		}
		return
	}

	// Find the specific version within the page
	var pageVersion *models.ExtractVersion
	for i, version := range page.Versions {
		if version.ExtractID == extractID {
			pageVersion = &page.Versions[i]
			break
		}
	}

	if pageVersion == nil {
		log.Println("Error: Specified version not found")
		http.Error(w, "Specified version not found", http.StatusNotFound)
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
	internalLinkDepthHandler := &improvementchain.InternalLinkDepthHandler{}
	mobileFriendlinessHandler := &improvementchain.MobileFriendlinessHandler{}
	imageSizeOptimisationHandler := &improvementchain.ImageSizeOptimisationHandler{}
	keywordDensityHandler := &improvementchain.KeywordDensityHandler{}
	breadcrumbValidationHandler := &improvementchain.BreadcrumbValidationHandler{}
	externalScriptEvaluationHandler := &improvementchain.ExternalScriptEvaluationHandler{}
	structuredDataValidationHandler := &improvementchain.StructuredDataValidationHandler{}
	contentReadabilityHandler := &improvementchain.ContentReadabilityHandler{}
	externalLinkQualityHandler := improvementchain.NewExternalLinkQualityHandler(os.Getenv("VIRUS_TOTAL_API_KEY"))
	socialMetaTagsHandler := &improvementchain.SocialMetaTagsHandler{}
	brokenLinkCheckerHandler := &improvementchain.BrokenLinkCheckerHandler{}

	// Set up the chain
	titleHandler.SetNext(metaDescriptionHandler)
	metaDescriptionHandler.SetNext(h1TagCountHandler)
	h1TagCountHandler.SetNext(imageAltTextHandler)
	imageAltTextHandler.SetNext(metaRobotsHandler)
	metaRobotsHandler.SetNext(pageLoadTimeHandler)
	pageLoadTimeHandler.SetNext(canonicalURLHandler)
	canonicalURLHandler.SetNext(internalLinkDepthHandler)
	internalLinkDepthHandler.SetNext(mobileFriendlinessHandler)
	mobileFriendlinessHandler.SetNext(imageSizeOptimisationHandler)
	imageSizeOptimisationHandler.SetNext(keywordDensityHandler)
	keywordDensityHandler.SetNext(breadcrumbValidationHandler)
	breadcrumbValidationHandler.SetNext(externalScriptEvaluationHandler)
	externalScriptEvaluationHandler.SetNext(structuredDataValidationHandler)
	structuredDataValidationHandler.SetNext(contentReadabilityHandler)
	contentReadabilityHandler.SetNext(externalLinkQualityHandler)
	externalLinkQualityHandler.SetNext(socialMetaTagsHandler)
	socialMetaTagsHandler.SetNext(brokenLinkCheckerHandler)

	// Execute the chain
	titleHandler.Handle(pageVersion, &improvements)

	// Update the improvements of the specific version
	pageVersion.Improvements = improvements

	// Update the MongoDB document with the generated improvements
	update := bson.M{
		"$set": bson.M{"versions.$.improvements": pageVersion.Improvements},
	}
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating page improvements in database: %v", err)
		http.Error(w, "Error updating page improvements in database", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully updated improvements for ExtractID: %s, URL: %s", extractID, url)

	// Send the improvements back as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(improvements)
}
