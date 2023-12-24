// post.go
package routes

import (
	"blog-server/database"
	"blog-server/types"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/charmbracelet/log"

	"github.com/gorilla/mux"
)

func GetPostBySlug(w http.ResponseWriter, r *http.Request) {
	// Extract post ID parameter from URL
	var slug string = mux.Vars(r)["slug"]
	if len(slug) == 0 {
		log.Error("Called GetPostBySlug without slug specified!")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Fetch the post by ID
	// post, err := fetchPostBySlug(slug)
	post, err := database.FetchPost(database.Slug, slug)
	if err != nil {
		log.Error("Error fetching post by ID", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Encode the post to JSON
	encoded, err := json.Marshal(post)
	if err != nil {
		log.Error("Error converting post to JSON", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(encoded)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the new post data
	var newPost types.Post
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		log.Error("Error decoding new post", "err", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Insert the new post into the database
	id, err := database.CreatePost(&newPost)
	if err != nil {
		log.Error("Error creating new post", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Fetch the complete post with the new ID
	createdPost, err := database.FetchPost(database.ID, id)
	if err != nil {
		log.Error("Error fetching created post", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Encode the complete post to JSON
	encoded, err := json.Marshal(createdPost)
	if err != nil {
		log.Error("Error converting created post to JSON", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(encoded)
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the updated post data
	var updatedPost types.Post
	err := json.NewDecoder(r.Body).Decode(&updatedPost)
	if err != nil {
		log.Error("Error decoding updated post", "err", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Update the post in the database
	err = database.UpdatePost(&updatedPost)
	if err != nil {
		log.Error("Error updating post", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	// Extract post ID parameter from URL
	postIDStr := mux.Vars(r)["id"]
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Error("Error parsing post ID", "err", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Delete the post by ID
	err = database.DeletePost(postID)
	if err != nil {
		log.Error("Error deleting post", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
