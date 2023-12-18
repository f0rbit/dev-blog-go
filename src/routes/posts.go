// posts.go
package routes

import (
	"blog-server/types"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
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

func GetPostsByCategory(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parsePaginationParams(r)
	if err != nil {
		log.Println("Error parsing params:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Extract category parameter from URL
	category := mux.Vars(r)["category"]

	posts, totalPosts, err := fetchPostsByCategory(category, limit, offset)
	if err != nil {
		log.Println("Error fetching posts by category:", err)
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

func fetchPostsByCategory(category string, limit, offset int) ([]types.Post, int, error) {
	var posts []types.Post
	var totalPosts int

	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return posts, totalPosts, err
	}
	defer db.Close()

	log.Printf("Searching for posts under category %s", category)

	// given a category, we want to get all the children category and include those in our select post query as an 'IN (<array of categories>)'
	var search_categories []any
	categories, err := fetchCategories()
	if err != nil {
		return posts, totalPosts, err
	}
	search_categories = append(search_categories, category)
	// go through categories here (they have properties .name and .parent)
	children := getChildrenCateogires(categories, category)
	for i := 0; i < len(children); i++ {
		search_categories = append(search_categories, children[i].Name)
	}
	log.Printf("Searching through categories: %v", search_categories)

	// Build the IN clause with placeholders
	placeholders := make([]string, len(search_categories))
	for i := range search_categories {
		placeholders[i] = "?"
	}
	inClause := strings.Join(placeholders, ",")

	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM posts WHERE category IN (%s)", inClause), search_categories...).Scan(&totalPosts)
	if err != nil {
		return posts, totalPosts, err
	}
	log.Printf("Found %d posts", totalPosts)

	// Query to fetch paginated posts for the given category
	query := fmt.Sprintf("SELECT id, slug, title, content, category FROM posts WHERE category IN (%s) LIMIT ? OFFSET ?", inClause)
	log.Printf("DB Query: %s", query)

	search_categories = append(search_categories, limit)
	search_categories = append(search_categories, offset)

	rows, err := db.Query(query, search_categories...)
	if err != nil {
		return posts, len(posts), err
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

func buildPreparedParams(length int) string {
	var placeholders []string
	for i := 0; i < length; i++ {
		placeholders = append(placeholders, "?")
	}
	return "(" + strings.Join(placeholders, ",") + ")"
}

func fetchPosts(limit, offset int) ([]types.Post, int, error) {
	var posts []types.Post
	var totalPosts int

	db, err := sql.Open("sqlite3", database)
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
