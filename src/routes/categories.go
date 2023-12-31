// categories.go
package routes

import (
	"blog-server/database"
	"blog-server/types"
	"blog-server/utils"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

// GetCategories handles the GET /categories route
func GetCategories(w http.ResponseWriter, r *http.Request) {
    user := utils.GetUser(r);
    if user == nil {
        utils.Unauthorized(w);
        return;
    }

    serveCategories(user, w)
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

    serveCategories(user, w)
    
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
    user := utils.GetUser(r);
    if user == nil {
        utils.Unauthorized(w);
        return;
    }
    
    category := mux.Vars(r)["name"]

    if category == "" {
        utils.LogError("No category provided", errors.New("No category provided in delete"), http.StatusBadRequest, w);
        return;
    }

    fetched, err := database.GetCategory(user, category);
    if err != nil {
        utils.LogError("Error fetching category", err, http.StatusInternalServerError, w);
        return;
    }

    if fetched.OwnerID != user.ID {
        utils.LogError("Invalid user ID", errors.New("User does not own the category for deletion"), http.StatusUnauthorized, w);
        return;
    }

    err = database.DeleteCategory(user, category)
    if err != nil {
        utils.LogError("Error deleting category", err, http.StatusInternalServerError, w);
        return;
    }
    
    serveCategories(user, w)
}

func serveCategories(user *types.User, w http.ResponseWriter) {
    // Fetch categories from the database
	categories, err := database.GetCategories(user)
	if err != nil {
        utils.LogError("Error fetching categories after creation", err, http.StatusInternalServerError, w);
		return
	}

    graph := database.ConstructCategoryGraph(categories, "root", user.ID);

    response :=  map[string]interface{}{
        "categories": categories,
        "graph": graph,
    }

    utils.ResponseJSON(response, w);
}
