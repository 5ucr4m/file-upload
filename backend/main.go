package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"fileupload/database"
	"fileupload/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Initialize database
	if err := database.InitDB("./fileupload.db"); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Create uploads directory
	os.MkdirAll("uploads/files", 0755)
	os.MkdirAll("uploads/temp", 0755)

	// Setup router
	router := mux.NewRouter()

	// Upload endpoints
	router.HandleFunc("/api/upload/init", handlers.InitUpload).Methods("POST")
	router.HandleFunc("/api/upload/chunk", handlers.UploadChunk).Methods("POST")
	router.HandleFunc("/api/upload/complete", handlers.CompleteUpload).Methods("POST")
	router.HandleFunc("/api/upload/cancel", handlers.CancelUpload).Methods("DELETE")
	router.HandleFunc("/api/upload/status", handlers.GetUploadStatus).Methods("GET")

	// File endpoints
	router.HandleFunc("/api/files", handlers.GetFiles).Methods("GET")
	router.HandleFunc("/api/files/file", handlers.GetFile).Methods("GET")
	router.HandleFunc("/api/files/download", handlers.DownloadFile).Methods("GET")
	router.HandleFunc("/api/files/versions", handlers.GetFileVersions).Methods("GET")

	// Share endpoints
	router.HandleFunc("/api/share/create", handlers.CreateShareLink).Methods("POST")
	router.HandleFunc("/api/share/access", handlers.GetSharedFile).Methods("POST")

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
