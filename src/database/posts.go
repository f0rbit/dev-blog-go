package database

import (
	"blog-server/types"
	"database/sql"
	"strings"
)

type Identifier string

const (
	ID   Identifier = "id"
	Slug Identifier = "slug"
)

func FetchPost(identifier Identifier, needle interface{}) (types.Post, error) {
	const base = `
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
    
    `
	var query string
	if identifier == "id" {
		query = base + " WHERE posts.id = ?"
	} else if identifier == "slug" {
		query = base + " WHERE posts.slug = ?"
	}

	var post types.Post
    var tags sql.NullString

    err := db.QueryRow(query, needle).Scan(
        &post.Id,
        &post.Slug,
        &post.Title,
        &post.Content,
        &post.Category,
        &post.Archived,
        &post.PublishAt,
        &post.CreatedAt,
        &post.UpdatedAt,
        &tags);

    if err != nil {
        return post, err
    }

    if tags.Valid {
        post.Tags = strings.Split(tags.String, ",")
    } else {
        post.Tags = []string{}
    }
    
    return post, nil;
}
