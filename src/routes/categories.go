// categories.go
package routes

import (
	"blog-server/database"
	"blog-server/utils"
	"net/http"
)

// GetCategories handles the GET /categories route
func GetCategories(w http.ResponseWriter, r *http.Request) {
	// Fetch categories from the database
	categories, err := database.GetCategories()
	if err != nil {
        utils.LogError("Error fetching categories", err, http.StatusInternalServerError, w);
		return
	}

    graph := database.ConstructCategoryGraph(categories, "root");

    response :=  map[string]interface{}{
        "categories": categories,
        "graph": graph,
    }

    utils.ResponseJSON(response, w);
}
