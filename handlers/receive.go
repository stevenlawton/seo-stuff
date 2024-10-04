package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sea-stuff/models"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AnalysisData struct {
	ExtractID             string              `json:"extractId" bson:"extractId"`
	URL                   string              `json:"url" bson:"url"`
	Title                 string              `json:"title" bson:"title"`
	TitleLength           int                 `json:"titleLength" bson:"titleLength"`
	MetaDescription       string              `json:"metaDescription" bson:"metaDescription"`
	MetaDescriptionLength int                 `json:"metaDescriptionLength" bson:"metaDescriptionLength"`
	MetaTags              map[string]string   `json:"metaTags" bson:"metaTags"`
	CanonicalURL          string              `json:"canonicalUrl" bson:"canonicalUrl"`
	IsCanonicalCorrect    bool                `json:"isCanonicalCorrect" bson:"isCanonicalCorrect"`
	HTags                 map[string][]string `json:"hTags" bson:"hTags"`
	H1TagCount            int                 `json:"h1TagCount" bson:"h1TagCount"`
	WordCount             int                 `json:"wordCount" bson:"wordCount"`
	PageDepth             int                 `json:"pageDepth" bson:"pageDepth"`
	PageLoadTimeSeconds   float64             `json:"pageLoadTimeSeconds" bson:"pageLoadTimeSeconds"`
	PageSizeBytes         int                 `json:"pageSizeBytes" bson:"pageSizeBytes"`
	Images                []ImageData         `json:"images" bson:"images"`
	InternalLinks         []string            `json:"internalLinks" bson:"internalLinks"`
	InternalLinksWithText []LinkWithText      `json:"internalLinksWithAnchorText" bson:"internalLinksWithAnchorText"`
	ExternalLinks         []string            `json:"externalLinks" bson:"externalLinks"`
	BrokenLinks           []string            `json:"brokenLinks" bson:"brokenLinks"`
	StructuredData        []string            `json:"structuredData" bson:"structuredData"`
	StructuredDataTypes   []string            `json:"structuredDataTypes" bson:"structuredDataTypes"`
	RobotsMetaTag         string              `json:"robotsMetaTag" bson:"robotsMetaTag"`
	Content               string              `json:"content" bson:"content"`
	CommonWords           [][]interface{}     `json:"commonWords" bson:"commonWords"`
	SocialTags            map[string]string   `json:"socialTags" bson:"socialTags"`
	Language              string              `json:"language" bson:"language"`
	Hreflangs             []string            `json:"hreflangs" bson:"hreflangs"`
	Breadcrumbs           []string            `json:"breadcrumbs" bson:"breadcrumbs"`
	IsMobileFriendly      bool                `json:"isMobileFriendly" bson:"isMobileFriendly"`
	ExternalScripts       []string            `json:"externalScripts" bson:"externalScripts"`
	ExternalStylesheets   []string            `json:"externalStylesheets" bson:"externalStylesheets"`
}

type ImageData struct {
	Src    string `json:"src" bson:"src"`
	Alt    string `json:"alt" bson:"alt"`
	Width  string `json:"width" bson:"width"`
	Height string `json:"height" bson:"height"`
}

type LinkWithText struct {
	Href       string `json:"href" bson:"href"`
	AnchorText string `json:"anchorText" bson:"anchorText"`
}

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

	var data AnalysisData
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

	// Convert AnalysisData to ExtractVersion
	extractVersion := models.ExtractVersion{
		ExtractID:             data.ExtractID,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		Title:                 data.Title,
		MetaTags:              data.MetaTags,
		RobotsMetaTag:         data.RobotsMetaTag,
		CommonWords:           flattenCommonWords(data.CommonWords),
		SocialTags:            convertMapToString(data.SocialTags),
		Hreflangs:             data.Hreflangs,
		HTags:                 data.HTags,
		InternalLinks:         data.InternalLinks,
		InternalLinksWithText: convertLinkWithText(data.InternalLinksWithText),
		ExternalLinks:         data.ExternalLinks,
		Images:                convertImages(data.Images),
		StructuredDataTypes:   data.StructuredDataTypes,
		StructuredData:        data.StructuredData,
		ExternalScripts:       data.ExternalScripts,
		ExternalStylesheets:   data.ExternalStylesheets,
		PageLoadTimeSeconds:   data.PageLoadTimeSeconds,
		PageSizeBytes:         data.PageSizeBytes,
		Language:              data.Language,
		Improvements:          []models.Improvement{}, // Assuming improvements are added later in the process
	}

	// Insert or update the data into the collection with upsert
	collection := client.Database("brandAdherence").Collection("pages")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find the page by URL
	filter := bson.M{"url": data.URL}
	var page models.Page
	err = collection.FindOne(ctx, filter).Decode(&page)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Printf("Error fetching page from database: %v", err)
		http.Error(w, "Error fetching page from database", http.StatusInternalServerError)
		return
	}

	// If no existing page is found, create a new one
	if err == mongo.ErrNoDocuments {
		page = models.Page{
			URL:      data.URL,
			Versions: []models.ExtractVersion{extractVersion},
		}
	} else {
		// Check if the version already exists and update it if necessary
		versionUpdated := false
		for i, version := range page.Versions {
			if version.ExtractID == data.ExtractID {
				page.Versions[i] = extractVersion
				versionUpdated = true
				break
			}
		}
		// If the version doesn't exist, add it as a new version
		if !versionUpdated {
			page.Versions = append(page.Versions, extractVersion)
		}
	}

	// Update the page document in the database
	update := bson.M{"$set": bson.M{"versions": page.Versions}}
	opts := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Printf("Error saving to the database: %v", err)
		http.Error(w, "Error saving to the database", http.StatusInternalServerError)
		return
	}

	// Log result of upsert
	if err == mongo.ErrNoDocuments {
		log.Printf("Inserted new page document for URL: %s", data.URL)
	} else {
		log.Printf("Updated page document for URL: %s", data.URL)
	}

	// Respond with a success message
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Data saved successfully")
}

// Helper function to convert []handlers.LinkWithText to []models.LinkWithText
func convertLinkWithText(links []LinkWithText) []models.LinkWithText {
	var convertedLinks []models.LinkWithText
	for _, link := range links {
		convertedLinks = append(convertedLinks, models.LinkWithText{
			Href:       link.Href,
			AnchorText: link.AnchorText,
		})
	}
	return convertedLinks
}

// Helper function to convert map to a JSON string
func convertMapToString(dataMap map[string]string) string {
	jsonString, err := json.Marshal(dataMap)
	if err != nil {
		log.Printf("Error converting map to string: %v", err)
		return ""
	}
	return string(jsonString)
}
func flattenCommonWords(commonWords [][]interface{}) []string {
	// Helper function to convert commonWords from [][]interface{} to []string
	var words []string
	for _, pair := range commonWords {
		if len(pair) > 0 {
			if word, ok := pair[0].(string); ok {
				words = append(words, word)
			}
		}
	}
	return words
}

func convertHTags(hTags map[string][]string) models.HeaderTags {
	// Helper function to convert HTags to HeaderTags struct
	return models.HeaderTags{
		H1: hTags["h1"],
		H2: hTags["h2"],
		H3: hTags["h3"],
	}
}

func convertImages(images []ImageData) []models.Image {
	// Helper function to convert ImageData to Image
	var convertedImages []models.Image
	for _, img := range images {
		width, err := strconv.Atoi(img.Width)
		if err != nil {
			log.Printf("Error converting width to int for image %s: %v", img.Src, err)
			width = 0 // Default value in case of error
		}

		height, err := strconv.Atoi(img.Height)
		if err != nil {
			log.Printf("Error converting height to int for image %s: %v", img.Src, err)
			height = 0 // Default value in case of error
		}

		convertedImages = append(convertedImages, models.Image{
			Src:    img.Src,
			Alt:    img.Alt,
			Width:  width,
			Height: height,
		})
	}
	return convertedImages
}
