package database

import (
	"blog-server/types"
	"database/sql"
	"strings"

    "github.com/charmbracelet/log"
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

func CreatePost(post *types.Post) (int, error) {
    var err error;
	// Insert the new post into the database
    _, err = db.Exec(
        `INSERT INTO posts (slug, title, content, category, publish_at) VALUES (?, ?, ?, ?, ?)`,
		post.Slug, 
        post.Title, 
        post.Content, 
        post.Category, 
        post.PublishAt)
    // get the id & update data structure
	if err != nil { return -1, err } 
	row := db.QueryRow("SELECT last_insert_rowid()")
	err = row.Scan(&post.Id)
	if err != nil { return -1, err }
    // insert any tags
    insert, err := db.Prepare("INSERT INTO tags (post_id, tag) VALUES (?, ?)")
    if err != nil { return -1, err }
    for _, s := range post.Tags {
        _, err = insert.Exec(post.Id, s);
        if err != nil { return -1, err }
    }
    log.Info("Inserted new post", "slug", post.Slug, "id", post.Id);
    return post.Id, err
}

func DeletePost(id int) error {
    _, err := db.Exec("DELETE FROM posts WHERE id = ?", id);
    if err == nil {
        log.Info("Deleted Post", "id", id);
    }
    return err;
}

func UpdatePost(updatedPost *types.Post) error {
    var err error;
    // update post
	_, err = db.Exec(`
    UPDATE 
        posts 
    SET 
        slug = ?,
        title = ?,
        content = ?,
        category = ?,
        archived = ?,
        publish_at = ? 
    WHERE 
        id = ?`,
	updatedPost.Slug, 
    updatedPost.Title, 
    updatedPost.Content, 
    updatedPost.Category, 
    updatedPost.Archived, 
    updatedPost.PublishAt, 
    updatedPost.Id);
    if err != nil { return err }
    // update tags
    _, err = db.Exec("DELETE FROM tags WHERE post_id = ?", updatedPost.Id)
    if err != nil { return err }
    insert, err := db.Prepare("INSERT INTO tags (post_id, tag) VALUES (?, ?)");
    if err != nil { return err }
    for _, s := range updatedPost.Tags {
        _, err = insert.Exec(updatedPost.Id, s)
        if err != nil { return err }
    }
	log.Info("Updated Post", "id", updatedPost.Id)
	return err
}
