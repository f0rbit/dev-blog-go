// posts.go
package routes

import (
	"blog-server/types"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parsePaginationParams(r)
	if err != nil {
		log.Println("Error parsing params:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	posts, totalPosts, err := fetchPosts(limit, offset)
	if err != nil {
		log.Println("Error fetching posts:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Calculate pagination information
	totalPages := (totalPosts + limit - 1) / limit
	currentPage := (offset / limit) + 1

	// Create a response structure including pagination information
	response := types.PostsResponse{
		Posts:       posts,
		TotalPosts:  totalPosts,
		TotalPages:  totalPages,
		PerPage:     limit,
		CurrentPage: currentPage,
	}

	encoded, err := json.Marshal(response)
	if err != nil {
		log.Println("Error converting posts to JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(encoded)
}

func fetchPosts(limit, offset int) ([]types.Post, int, error) {
	var posts []types.Post
	var totalPosts int

	db, err := sql.Open("sqlite3", "../db/sqlite.db")
	if err != nil {
		return posts, totalPosts, err
	}
	defer db.Close()

	// Query to get total number of posts
	err = db.QueryRow("SELECT COUNT(*) FROM posts").Scan(&totalPosts)
	if err != nil {
		return posts, totalPosts, err
	}

	// Query to fetch paginated posts
	rows, err := db.Query("SELECT id, slug, title, content, category FROM posts LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return posts, totalPosts, err
	}

	for rows.Next() {
		var post types.Post
		err := rows.Scan(&post.Id, &post.Slug, &post.Title, &post.Content, &post.Category)
		if err != nil {
			return posts, totalPosts, err
		}
		posts = append(posts, post)
	}

	return posts, totalPosts, nil
}

func parsePaginationParams(r *http.Request) (int, int, error) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	limit := 10
	offset := 0
	var err error
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0, err
		}
	}
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return 0, 0, err
		}
	}
	return limit, offset, nil
}
