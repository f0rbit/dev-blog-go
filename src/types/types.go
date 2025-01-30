package types

import (
	"database/sql"
	"time"
)

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
	Format      string    `json:"format"`
	Category    string    `json:"category"`
	Tags        []string  `json:"tags"`
	Archived    bool      `json:"archived"`
	Description string    `json:"description"`
	ProjectID   string    `json:"project_id"`
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

type FetchLink struct {
	PostID      int       `json:"post_id"`
	FetchSource int       `json:"fetch_source"`
	Identifier  string    `json:"identifier"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Project struct {
	ID             string    `json:"id"`
	ProjectID      string    `json:"project_id"`
	OwnerID        string    `json:"owner_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Specification  string    `json:"specification"`
	RepoURL        string    `json:"repo_url"`
	RepoID         string    `json:"repo_id"`
	IconURL        string    `json:"icon_url"`
	Status         string    `json:"status"`
	Deleted        bool      `json:"deleted"`
	LinkURL        string    `json:"link_url"`
	LinkText       string    `json:"link_text"`
	Visibility     string    `json:"visibility"`
	CurrentVersion string    `json:"current_version"`
	ScanBranch     string    `json:"scan_branch"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ProjectCache struct {
	ID        int            `json:"id"`
	UserID    int            `json:"user_id"`
	Status    string         `json:"status"`
	URL       string         `json:"url"`
	Data      sql.NullString `json:"data"`
	FetchedAt time.Time      `json:"fetched_at"`
}

type ProjectLink struct {
	ID          int       `json:"id"`
	PostID      int       `json:"post_id"`
	ProjectUUID string    `json:"project_uuid"`
	ProjectID   string    `json:"project_id"`
	CreatedAt   time.Time `json:"created_at"`
}
