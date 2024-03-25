// post.go
package routes

import (
	"blog-server/database"
	"blog-server/types"
	"blog-server/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetPostBySlug(w http.ResponseWriter, r *http.Request) {
	user := utils.GetUser(r)
	if user == nil {
		utils.Unauthorized(w)
		return
	}

	var slug string = mux.Vars(r)["slug"]

	if len(slug) == 0 {
		utils.LogError("No slug specified", nil, http.StatusBadRequest, w)
		return
	}

	post, err := database.FetchPost(user, database.Slug, slug)
	if err != nil {
		utils.LogError("Error fetching post by slug", err, http.StatusNotFound, w)
		return
	}

	utils.ResponseJSON(post, w)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	user := utils.GetUser(r)
	if user == nil {
		utils.Unauthorized(w)
		return
	}
	var newPost types.Post
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		utils.LogError("Error decoding new post", err, http.StatusBadRequest, w)
		return
	}

    // verify that author id is the logged in user
    if newPost.AuthorID != user.ID {
        utils.LogError("Invalid author ID", errors.New("Create Post authorID doesn't match userID"), http.StatusBadRequest, w);
        return;
    }

	id, err := database.CreatePost(newPost)
	if err != nil {
	utils.LogError("Error creating new post", err, http.StatusInternalServerError, w)
		return
	}

	createdPost, err := database.FetchPost(user, database.ID, id)
	if err != nil {
		utils.LogError("Error fetching created post", err, http.StatusInternalServerError, w)
		return
	}

	utils.ResponseJSON(createdPost, w)
}

func EditPost(w http.ResponseWriter, r *http.Request) {
    user := utils.GetUser(r)
	if user == nil {
		utils.Unauthorized(w)
		return
	}

	// Parse the request body to get the updated post data
	var updatedPost types.Post
	err := json.NewDecoder(r.Body).Decode(&updatedPost)
	if err != nil {
		utils.LogError("Error decoding updated post", err, http.StatusBadRequest, w)
		return
	}

    if updatedPost.AuthorID != user.ID {
        utils.LogError("Invalid author ID", errors.New("Update Post authorID doesn't match userID"), http.StatusBadRequest, w);
        return;
    }

	// Update the post in the database
	err = database.UpdatePost(&updatedPost)
	if err != nil {
		utils.LogError("Error updating post", err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
    user := utils.GetUser(r)
	if user == nil {
		utils.Unauthorized(w)
		return
	}

	// Extract post ID parameter from URL
	postIDStr := mux.Vars(r)["id"]
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		utils.LogError("Error parsing post ID", err, http.StatusBadRequest, w)
		return
	}

    post, err := database.FetchPost(user, database.ID, postID);
    if err != nil {
        utils.LogError("Error fetching post", err, http.StatusInternalServerError, w);
        return;
    }

    if post.AuthorID != user.ID {
        utils.LogError("Invalid author ID", errors.New("Delete post authorID doesn't match userID"), http.StatusBadRequest, w);
        return;
    }

	// Delete the post by ID
	err = database.DeletePost(postID)
	if err != nil {
		utils.LogError("Error deleting post", err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
