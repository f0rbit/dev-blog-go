// posts.go
package routes

import (
	"blog-server/database"
	"blog-server/types"
	"blog-server/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func FetchPosts(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parsePaginationParams(r)
	if err != nil {
        utils.LogError("Error parsing params", err, http.StatusBadRequest, w);
		return
	}

	// Extract category parameter from URL
	category := mux.Vars(r)["category"]
	if category == "" {
		category = "root"
	}
	tag := r.URL.Query().Get("tag")

    user := utils.GetUser(r);

    if user == nil {
        utils.Unauthorized(w);
        return;
    }

	posts, totalPosts, err := database.GetPosts(user, category, tag, limit, offset)
	if err != nil {
        utils.LogError("Error fetching posts by category", err, http.StatusInternalServerError, w);
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

    utils.ResponseJSON(response, w);
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
