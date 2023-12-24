// post.go
package routes

import (
	"blog-server/database"
	"blog-server/types"
	"blog-server/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetPostBySlug(w http.ResponseWriter, r *http.Request) {
	var slug string = mux.Vars(r)["slug"];
	if len(slug) == 0 {
        utils.LogError("No slug specified", nil, http.StatusBadRequest, w);
		return;
	}

	post, err := database.FetchPost(database.Slug, slug)
	if err != nil {
        utils.LogError("Error fetching post by ID", err, http.StatusInternalServerError, w);
		return
	}

    utils.ResponseJSON(post, w);
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var newPost types.Post
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
        utils.LogError("Error decoding new post", err, http.StatusBadRequest, w);
		return
	}

	id, err := database.CreatePost(&newPost)
	if err != nil {
        utils.LogError("Error creating new post", err, http.StatusInternalServerError, w);
		return
	}

	createdPost, err := database.FetchPost(database.ID, id)
	if err != nil {
        utils.LogError("Error fetching created post", err, http.StatusInternalServerError, w);
		return
	}

    utils.ResponseJSON(createdPost, w);
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the updated post data
	var updatedPost types.Post
	err := json.NewDecoder(r.Body).Decode(&updatedPost)
	if err != nil {
        utils.LogError("Error decoding updated post", err, http.StatusBadRequest, w);
		return
	}

	// Update the post in the database
	err = database.UpdatePost(&updatedPost)
	if err != nil {
        utils.LogError("Error updating post", err, http.StatusInternalServerError, w);
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	// Extract post ID parameter from URL
	postIDStr := mux.Vars(r)["id"]
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
        utils.LogError("Error parsing post ID", err, http.StatusBadRequest, w);
		return
	}

	// Delete the post by ID
	err = database.DeletePost(postID)
	if err != nil {
        utils.LogError("Error deleting post", err, http.StatusInternalServerError, w);
		return
	}

	w.WriteHeader(http.StatusOK)
}
