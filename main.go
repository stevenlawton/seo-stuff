package main

import (
	"log"
	"net/http"
	"os"
	"sea-stuff/db"
	"sea-stuff/handlers"
	"time"
)

func main() {
	// Get MongoDB URI from environment variable
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is required")
	}

	// Connect to MongoDB
	client, err := db.ConnectMongo(mongoURI, 10*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.DisconnectMongo(client)

	// Set the Mongo client in handlers
	handlers.SetClient(client)

	// Define the endpoints and start the HTTP server
	http.HandleFunc("/api/receive_data", handlers.HandlePost)
	http.HandleFunc("/pages", handlers.HandleListPages)
	http.HandleFunc("/pages/", handlers.HandlePageDetails)
	http.HandleFunc("/delete_by_extract_id", handlers.HandleDeleteByExtractID)
	http.HandleFunc("/run_scanners", handlers.HandleRunScanners)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
