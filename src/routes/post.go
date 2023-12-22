// post.go
package routes

import (
	"blog-server/types"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/charmbracelet/log"

	"github.com/gorilla/mux"
)

var database string = os.Getenv("DATABASE")

func GetPostBySlug(w http.ResponseWriter, r *http.Request) {
	// Extract post ID parameter from URL
	var slug string = mux.Vars(r)["slug"]
	if len(slug) == 0 {
		log.Error("Called GetPostBySlug without slug specified!")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Fetch the post by ID
	post, err := fetchPostBySlug(slug)
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
	err = insertPost(&newPost)
	if err != nil {
		log.Error("Error creating new post", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Fetch the complete post with the new ID
	createdPost, err := fetchPostByID(newPost.Id)
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
	err = updatePost(&updatedPost)
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
	err = deletePost(postID)
	if err != nil {
		log.Error("Error deleting post", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func fetchPostBySlug(slug string) (types.Post, error) {
	var post types.Post

	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return post, err
	}
	defer db.Close()

	err = db.QueryRow("SELECT posts.id, posts.slug, posts.title, posts.content, posts.category, posts.created_at, updated_at, GROUP_CONACT(tags.tag) AS tags FROM posts LEFT JOIN tags ON posts.id = tags.post_id GROUP BY posts.id WHERE posts.slug = ?", slug).Scan(&post.Id, &post.Slug, &post.Title, &post.Content, &post.Category, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return post, err
	}

	return post, nil
}

func fetchPostByID(postID int) (types.Post, error) {
	var post types.Post

	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return post, err
	}
	defer db.Close()

	// Query to fetch the post by ID
	err = db.QueryRow("SELECT posts.id, posts.slug, posts.title, posts.content, posts.category, posts.created_at, updated_at, GROUP_CONACT(tags.tag) AS tags FROM posts LEFT JOIN tags ON posts.id = tags.post_id GROUP BY posts.id WHERE id = ?", postID).Scan(&post.Id, &post.Slug, &post.Title, &post.Content, &post.Category, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return post, err
	}

	return post, nil
}

func insertPost(newPost *types.Post) error {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return err
	}
	defer db.Close()

	// Insert the new post into the database
	_, err = db.Exec("INSERT INTO posts (slug, title, content, category) VALUES (?, ?, ?, ?)",
		newPost.Slug, newPost.Title, newPost.Content, newPost.Category)

	if err != nil {
		return err
	}

	// Assuming your database is configured to auto-increment the ID,
	// retrieve the last inserted ID using the LastInsertId method
	row := db.QueryRow("SELECT last_insert_rowid()")

	err = row.Scan(&newPost.Id)
	if err != nil {
		return err
	}

	log.Infof("Inserted new post '%s' with ID: %d", newPost.Slug, newPost.Id)

	return err
}

func updatePost(updatedPost *types.Post) error {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return err
	}
	defer db.Close()

	// Update the post in the database
	_, err = db.Exec("UPDATE posts SET slug = ?, title = ?, content = ?, category = ? WHERE id = ?",
		updatedPost.Slug, updatedPost.Title, updatedPost.Content, updatedPost.Category, updatedPost.Id)

	log.Info("Updated Post", "id", updatedPost.Id)

	return err
}

func deletePost(postID int) error {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return err
	}
	defer db.Close()

	// Delete the post by ID
	_, err = db.Exec("DELETE FROM posts WHERE id = ?", postID)

	log.Info("Deleted Post", "id", postID)

	return err
}
