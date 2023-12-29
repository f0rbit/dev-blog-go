package types

import "time"

type Category struct {
	Name   string `json:"name"`
	Parent string `json:"parent"`
}

type CategoryNode struct {
    Name string `json:"name"`
    Children []CategoryNode `json:"children"`
}

type Post struct {
	Id        int       `json:"id"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	Tags      []string  `json:"tags"`
	Archived  bool      `json:"archived"`
    PublishAt time.Time `json:"publish_at" time_format:"sql_datetime"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostsResponse struct {
	Posts       []Post `json:"posts"`
	TotalPosts  int    `json:"total_posts"`
	TotalPages  int    `json:"total_pages"`
	PerPage     int    `json:"per_page"`
	CurrentPage int    `json:"current_page"`
}


// GitHubUser represents the structure of a GitHub user's JSON response.
type GitHubUser struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	// Add other user-related fields as needed
}

// User represents a user's information.
type User struct {
	ID        int
	GitHubID  int
	Username  string
	Email     string
	AvatarURL string
}
