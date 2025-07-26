package main

import (
	"log"
	"net/http"

	"digitalSign/handlers"
	"digitalSign/middleware"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Apply CORS middleware
	r.Use(middleware.CORS)

	// Serve static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Routes
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")

	// Digital Signature routes
	r.HandleFunc("/digital-signature", handlers.DigitalSignaturePageHandler).Methods("GET")
	r.HandleFunc("/api/generate-keys", handlers.GenerateKeysHandler).Methods("POST")
	r.HandleFunc("/api/sign-document", handlers.SignDocumentHandler).Methods("POST")
	r.HandleFunc("/api/verify-signature", handlers.VerifySignatureHandler).Methods("POST")

	// Steganography routes
	r.HandleFunc("/steganography", handlers.SteganographyPageHandler).Methods("GET")
	r.HandleFunc("/api/hide-message", handlers.HideMessageHandler).Methods("POST")
	r.HandleFunc("/api/extract-message", handlers.ExtractMessageHandler).Methods("POST")

	// Combined operations
	r.HandleFunc("/combined", handlers.CombinedPageHandler).Methods("GET")
	r.HandleFunc("/api/sign-and-hide", handlers.SignAndHideHandler).Methods("POST")
	r.HandleFunc("/api/extract-and-verify", handlers.ExtractAndVerifyHandler).Methods("POST")

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
