// post.go
package routes

import (
	"blog-server/types"
    "blog-server/database"
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
    post, err := database.FetchPost(database.Slug, slug); 
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
    createdPost, err := database.FetchPost(database.ID, newPost.Id)
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

func insertPost(newPost *types.Post) error {
    db := database.Connection();
    var err error;

	// Insert the new post into the database
    _, err = db.Exec("INSERT INTO posts (slug, title, content, category, publish_at) VALUES (?, ?, ?, ?, ?)",
		newPost.Slug, newPost.Title, newPost.Content, newPost.Category, newPost.PublishAt)

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

    insert, err := db.Prepare("INSERT INTO tags (post_id, tag) VALUES (?, ?)")
    // insert any tags
    for _, s := range newPost.Tags {
        insert.Exec(newPost.Id, s)
    }
    log.Info("Inserted tags", "tags", newPost.Tags);

	return err
}

func updatePost(updatedPost *types.Post) error {
    db := database.Connection();
    var err error;

    log.Info("updating", "publish_at", updatedPost.PublishAt)
	// Update the post in the database
	_, err = db.Exec("UPDATE posts SET slug = ?, title = ?, content = ?, category = ?, archived = ?, publish_at = ? WHERE id = ?",
		updatedPost.Slug, updatedPost.Title, updatedPost.Content, updatedPost.Category, updatedPost.Archived, updatedPost.PublishAt, updatedPost.Id)

    // Update the tags
    // first we drop all the previous tags and then re-add
    _, err = db.Exec("DELETE FROM tags WHERE post_id = ?", updatedPost.Id)
    insert, err := db.Prepare("INSERT INTO tags (post_id, tag) VALUES (?, ?)");

    for _, s := range updatedPost.Tags {
        insert.Exec(updatedPost.Id, s)
    }

	log.Info("Updated Post", "id", updatedPost.Id)

	return err
}

func deletePost(postID int) error {
    var db = database.Connection()
    var err error;

	// Delete the post by ID
	_, err = db.Exec("DELETE FROM posts WHERE id = ?", postID)

	log.Info("Deleted Post", "id", postID)

	return err
}
