// categories.go
package routes

import (
	"blog-server/types"
	"database/sql"
	"encoding/json"
	"github.com/charmbracelet/log"
	"net/http"
)

// GetCategories handles the GET /categories route
func GetCategories(w http.ResponseWriter, r *http.Request) {
	// Fetch categories from the database
	categories, err := fetchCategories()
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

func fetchCategories() ([]types.Category, error) {
	var categories []types.Category
	db, err := sql.Open("sqlite3", "db/sqlite.db")
	if err != nil {
		return categories, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT name, parent FROM categories")
	if err != nil {
		return categories, err
	}
	for rows.Next() {
		var category types.Category
		err := rows.Scan(&category.Name, &category.Parent)
		if err != nil {
			return categories, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func getChildrenCateogires(categories []types.Category, parent string) []types.Category {
	var cats []types.Category

	for i := 0; i < len(categories); i++ {
		if categories[i].Parent == parent {
			cats = append(cats, categories[i])
			// add all the children as well
			var children = getChildrenCateogires(categories, categories[i].Name)
			for j := 0; j < len(children); j++ {
				cats = append(cats, children[j])
			}
		}
	}

	return cats
}
