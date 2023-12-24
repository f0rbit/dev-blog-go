// categories.go
package routes

import (
	"blog-server/database"
	"blog-server/types"
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

func getChildrenCategories(categories []types.Category, parent string) []types.Category {
	var cats []types.Category

	for i := 0; i < len(categories); i++ {
		if categories[i].Parent == parent {
			cats = append(cats, categories[i])
			// add all the children as well
			var children = getChildrenCategories(categories, categories[i].Name)
			for j := 0; j < len(children); j++ {
				cats = append(cats, children[j])
			}
		}
	}

	return cats
}
