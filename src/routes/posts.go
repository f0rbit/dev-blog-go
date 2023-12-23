// posts.go
package routes

import (
	"blog-server/database"
	"blog-server/types"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gorilla/mux"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parsePaginationParams(r)
	if err != nil {
		log.Error("Error parsing params", "err", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	tag := r.URL.Query().Get("tag")

	posts, totalPosts, err := fetchPosts("root", limit, offset, tag)
	if err != nil {
		log.Error("Error fetching posts", "err", err)
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

func GetPostsByCategory(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parsePaginationParams(r)
	if err != nil {
		log.Error("Error parsing params", "err", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Extract category parameter from URL
	category := mux.Vars(r)["category"]
	tag := r.URL.Query().Get("tag")

	posts, totalPosts, err := fetchPosts(category, limit, offset, tag)
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

func fetchPosts(category string, limit, offset int, tag string) ([]types.Post, int, error) {
	var posts []types.Post
	var totalPosts int
    var db = database.Connection();

	log.Info("Searching for posts", "category", category, "tag", tag)

	// given a category, we want to get all the children category and include those in our select post query as an 'IN (<array of categories>)'
	var search_categories []any
	categories, err := fetchCategories()
	if err != nil {
		return posts, totalPosts, err
	}
    search_categories = append(search_categories, category)
	if category == "root" {
        // if we are searching at root, we can just append all the categories
        /** @todo in this case, remove the WHERE category from the search clause */
        for i := 0; i < len(categories); i++ {
            search_categories = append(search_categories, categories[i].Name)
        }
	} else {
		// go through categories here (they have properties .name and .parent)
		children := getChildrenCateogires(categories, category)
		for i := 0; i < len(children); i++ {
			search_categories = append(search_categories, children[i].Name)
		}
	}
	log.Info("Searching through categories", "search_categories", search_categories)

	// Build the IN clause with placeholders
	placeholders := make([]string, len(search_categories))
	for i := range search_categories {
		placeholders[i] = "?"
	}
	inClause := strings.Join(placeholders, ",")

	var params []any
	if tag == "" {
		err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM posts WHERE category IN (%s)", inClause), search_categories...).Scan(&totalPosts)
	} else {
		params = append(params, tag)
		params = append(params, search_categories...)
		err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM posts LEFT JOIN tags ON tags.post_id = posts.id WHERE tags.tag = ? AND posts.category IN (%s)", inClause), params...).Scan(&totalPosts)
	}
	if err != nil {
		return posts, totalPosts, err
	}
	log.Infof("Found %d posts", totalPosts)

	// Query to fetch paginated posts for the given category
	query := fmt.Sprintf(`
    SELECT 
        posts.id, 
        posts.slug, 
        posts.title, 
        posts.content, 
        posts.category, 
        posts.archived,
        posts.publish_at,
        posts.created_at, 
        posts.updated_at,
        GROUP_CONCAT(tags.tag) AS tags
    FROM 
        posts 
    LEFT JOIN
        tags ON posts.id = tags.post_id
    WHERE 
        posts.category IN (%s) 
    GROUP BY
        posts.id
    LIMIT ? 
    OFFSET ?`, inClause)

	search_categories = append(search_categories, limit)
	search_categories = append(search_categories, offset)
    log.Info("Searching with limit & offset of", "limit", limit, "offset", offset)

	rows, err := db.Query(query, search_categories...)
	if err != nil {
        log.Error("Error during posts query", "err", err)
		return posts, len(posts), err
	}

	for rows.Next() {
		var post types.Post
        var tags sql.NullString
		err := rows.Scan(&post.Id, &post.Slug, &post.Title, &post.Content, &post.Category, &post.Archived, &post.PublishAt, &post.CreatedAt, &post.UpdatedAt, &tags)

		if err != nil {
            log.Error("Error during scanning posts", "err", err)
			return posts, totalPosts, err
		}

        if tags.Valid {
            post.Tags = strings.Split(tags.String, ",")
        } else {
            post.Tags = []string{}
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
