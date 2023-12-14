// categories.go
package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// Category struct represents a category in the system
type Category struct {
	Name   string `json:"name"`
	Parent string `json:"parent"`
}

// GetCategories handles the GET /categories route
func GetCategories(w http.ResponseWriter, r *http.Request) {
	// Fetch categories from the database
	categories, err := fetchCategories()
	if err != nil {
		log.Println("Error fetching categories:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Convert categories to JSON
	jsonCategories, err := json.Marshal(categories)
	if err != nil {
		log.Println("Error marshaling categories to JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonCategories)
}

func fetchCategories() ([]Category, error) {
	db, err := sql.Open("sqlite3", "../db/sqlite.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// TODO: Implement database query to fetch categories
	rows, err := db.Query("SELECT name, parent FROM categories")
	// Handle errors appropriately
	if err != nil {
		log.Fatal(err)
	}

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.Name, &category.Parent)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}
