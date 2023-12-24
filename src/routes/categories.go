// categories.go
package routes

import (
	"blog-server/database"
	"encoding/json"
	"net/http"

	"github.com/charmbracelet/log"
)

// GetCategories handles the GET /categories route
func GetCategories(w http.ResponseWriter, r *http.Request) {
	// Fetch categories from the database
	categories, err := database.GetCategories()
	if err != nil {
		log.Error("Error fetching categories", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Convert categories to JSON
	encoded, err := json.Marshal(categories)
	if err != nil {
		log.Error("Error marshaling categories to JSON", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(encoded)
}
