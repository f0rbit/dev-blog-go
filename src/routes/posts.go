// posts.go
package routes

import (
	"blog-server/database"
	"blog-server/types"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gorilla/mux"
)

func FetchPosts(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parsePaginationParams(r)
	if err != nil {
		log.Error("Error parsing params", "err", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Extract category parameter from URL
	category := mux.Vars(r)["category"]
	if category == "" {
		category = "root"
	}
	tag := r.URL.Query().Get("tag")

	posts, totalPosts, err := database.GetPosts(category, tag, limit, offset)
	if err != nil {
		log.Error("Error fetching posts by category", "err", err)
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
		log.Error("Error converting posts to JSON", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(encoded)
}

func buildPreparedParams(length int) string {
	var placeholders []string
	for i := 0; i < length; i++ {
		placeholders = append(placeholders, "?")
	}
	return "(" + strings.Join(placeholders, ",") + ")"
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
