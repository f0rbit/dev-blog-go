// categories.go
package routes

import (
	"blog-server/database"
	"blog-server/types"
	"blog-server/utils"
	"encoding/json"
	"errors"
	"net/http"
)

// GetCategories handles the GET /categories route
func GetCategories(w http.ResponseWriter, r *http.Request) {
    user := utils.GetUser(r);
    if user == nil {
        utils.Unauthorized(w);
        return;
    }
	// Fetch categories from the database
	categories, err := database.GetCategories(user)
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

func CreateCategory(w http.ResponseWriter, r *http.Request) {
    user := utils.GetUser(r);
    if user == nil {
        utils.Unauthorized(w);
        return;
    }

    // create the category from the body
    var newCategory types.Category
    err := json.NewDecoder(r.Body).Decode(&newCategory)
    if err != nil {
        utils.LogError("Error decing new category", err, http.StatusBadRequest, w);
        return;
    }

    if newCategory.Name == "root" {
        utils.LogError("Invalid category name", errors.New("Attempted to create category with name 'root'"), http.StatusBadRequest, w);
        return;
    }

    // verify that the owner_id is the same as the user_id
    if newCategory.OwnerID != user.ID {
        utils.LogError("Invalid userID", errors.New("Create category userID doesn't match userID"), http.StatusBadRequest, w);
        return;
    }

    err = database.CreateCategory(newCategory)
    if err != nil {
        utils.LogError("Error creating category", err, http.StatusInternalServerError, w);
        return;
    }

    // Fetch categories from the database
	categories, err := database.GetCategories(user)
	if err != nil {
        utils.LogError("Error fetching categories after creation", err, http.StatusInternalServerError, w);
		return
	}

    graph := database.ConstructCategoryGraph(categories, "root");

    response :=  map[string]interface{}{
        "categories": categories,
        "graph": graph,
    }

    utils.ResponseJSON(response, w);
    
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
    // TODO: implement delete function
}
