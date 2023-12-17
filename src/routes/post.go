// post.go
package routes

import (
	"blog-server/types"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetPostByID(w http.ResponseWriter, r *http.Request) {
	// Extract post ID parameter from URL
	postIDStr := mux.Vars(r)["id"]
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Println("Error parsing post ID:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Fetch the post by ID
	post, err := fetchPostByID(postID)
	if err != nil {
		log.Println("Error fetching post by ID:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Encode the post to JSON
	encoded, err := json.Marshal(post)
	if err != nil {
		log.Println("Error converting post to JSON:", err)
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
		log.Println("Error decoding new post:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Insert the new post into the database
	err = insertPost(&newPost)
	if err != nil {
		log.Println("Error creating new post:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the updated post data
	var updatedPost types.Post
	err := json.NewDecoder(r.Body).Decode(&updatedPost)
	if err != nil {
		log.Println("Error decoding updated post:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Update the post in the database
	err = updatePost(&updatedPost)
	if err != nil {
		log.Println("Error updating post:", err)
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
		log.Println("Error parsing post ID:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Delete the post by ID
	err = deletePost(postID)
	if err != nil {
		log.Println("Error deleting post:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func fetchPostByID(postID int) (types.Post, error) {
	var post types.Post

	db, err := sql.Open("sqlite3", "../db/sqlite.db")
	if err != nil {
		return post, err
	}
	defer db.Close()

	// Query to fetch the post by ID
	err = db.QueryRow("SELECT id, slug, title, content, category FROM posts WHERE id = ?", postID).
		Scan(&post.Id, &post.Slug, &post.Title, &post.Content, &post.Category)

	if err != nil {
		return post, err
	}

	return post, nil
}

func insertPost(newPost *types.Post) error {
	db, err := sql.Open("sqlite3", "../db/sqlite.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Insert the new post into the database
	_, err = db.Exec("INSERT INTO posts (slug, title, content, category) VALUES (?, ?, ?, ?)",
		newPost.Slug, newPost.Title, newPost.Content, newPost.Category)

	return err
}

func updatePost(updatedPost *types.Post) error {
	db, err := sql.Open("sqlite3", "../db/sqlite.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Update the post in the database
	_, err = db.Exec("UPDATE posts SET slug = ?, title = ?, content = ?, category = ? WHERE id = ?",
		updatedPost.Slug, updatedPost.Title, updatedPost.Content, updatedPost.Category, updatedPost.Id)

	return err
}

func deletePost(postID int) error {
	db, err := sql.Open("sqlite3", "../db/sqlite.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Delete the post by ID
	_, err = db.Exec("DELETE FROM posts WHERE id = ?", postID)

	return err
}
