package types

import "time"

type Category struct {
	Name    string `json:"name"`
	Parent  string `json:"parent"`
	OwnerID int    `json:"owner_id"`
}

type CategoryNode struct {
	Name     string         `json:"name"`
	Children []CategoryNode `json:"children"`
	OwnerID  int            `json:"owner_id"`
}

type Post struct {
	Id          int       `json:"id"`
	Slug        string    `json:"slug"`
	AuthorID    int       `json:"author_id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Category    string    `json:"category"`
	Tags        []string  `json:"tags"`
	Archived    bool      `json:"archived"`
	Description string    `json:"description"`
	PublishAt   time.Time `json:"publish_at" time_format:"sql_datetime"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PostsResponse struct {
	Posts       []Post `json:"posts"`
	TotalPosts  int    `json:"total_posts"`
	TotalPages  int    `json:"total_pages"`
	PerPage     int    `json:"per_page"`
	CurrentPage int    `json:"current_page"`
}

type GitHubUser struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type User struct {
	ID        int    `json:"user_id"`
	GitHubID  int    `json:"github_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AccessKey struct {
	ID        int       `json:"id"`
	Value     string    `json:"value"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Note      string    `json:"note"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Integration struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	LastFetch time.Time `json:"last_fetch"`
	Location  string    `json:"location"`
	Source    string    `json:"source"`
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
