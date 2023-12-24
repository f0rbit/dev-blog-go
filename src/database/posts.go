package database

import (
	"blog-server/types"
	"database/sql"
	"fmt"
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
		&tags)

	if err != nil {
		return post, err
	}

	if tags.Valid {
		post.Tags = strings.Split(tags.String, ",")
	} else {
		post.Tags = []string{}
	}

	return post, nil
}

func CreatePost(post *types.Post) (int, error) {
	var err error
	// Insert the new post into the database
	_, err = db.Exec(
		`INSERT INTO posts (slug, title, content, category, publish_at) VALUES (?, ?, ?, ?, ?)`,
		post.Slug,
		post.Title,
		post.Content,
		post.Category,
		post.PublishAt)
	// get the id & update data structure
	if err != nil {
		return -1, err
	}
	row := db.QueryRow("SELECT last_insert_rowid()")
	err = row.Scan(&post.Id)
	if err != nil {
		return -1, err
	}
	// insert any tags
	insert, err := db.Prepare("INSERT INTO tags (post_id, tag) VALUES (?, ?)")
	if err != nil {
		return -1, err
	}
	for _, s := range post.Tags {
		_, err = insert.Exec(post.Id, s)
		if err != nil {
			return -1, err
		}
	}
	log.Info("Inserted new post", "slug", post.Slug, "id", post.Id)
	return post.Id, err
}

func DeletePost(id int) error {
	_, err := db.Exec("DELETE FROM posts WHERE id = ?", id)
	if err == nil {
		log.Info("Deleted Post", "id", id)
	}
	return err
}

func UpdatePost(updatedPost *types.Post) error {
	var err error
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
		updatedPost.Id)
	if err != nil {
		return err
	}
	// update tags
	_, err = db.Exec("DELETE FROM tags WHERE post_id = ?", updatedPost.Id)
	if err != nil {
		return err
	}
	insert, err := db.Prepare("INSERT INTO tags (post_id, tag) VALUES (?, ?)")
	if err != nil {
		return err
	}
	for _, s := range updatedPost.Tags {
		_, err = insert.Exec(updatedPost.Id, s)
		if err != nil {
			return err
		}
	}
	log.Info("Updated Post", "id", updatedPost.Id)
	return err
}

func GetPosts(category, tag string, limit, offset int) ([]types.Post, int, error) {
	var posts []types.Post
	var totalPosts int
	log.Info("Searching for posts", "category", category, "tag", tag)

	// given a category, we want to get all the children category and include those in our select post query as an 'IN (<array of categories>)'
	var search_categories []string
	categories, err := GetCategories()
	if err != nil {
		return posts, totalPosts, err
	}

	search_categories = append(search_categories, category)
	if category == "root" || category == "" {
		// if we are searching at root, we can just append all the categories
		/** @todo in this case, remove the WHERE category from the search clause */
		for i := 0; i < len(categories); i++ {
			search_categories = append(search_categories, categories[i].Name)
		}
	} else {
		// go through categories here (they have properties .name and .parent)
		children := getChildrenCategories(categories, category)
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
	for _, s := range search_categories {
		params = append(params, s)
	}

	// get the count of total posts
	if tag == "" {
		err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM posts WHERE category IN (%s)", inClause), params...).Scan(&totalPosts)
	} else {
		params = append(params, tag)
		err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM posts LEFT JOIN tags ON tags.post_id = posts.id WHERE posts.category IN (%s) AND tags.tag = ?", inClause), params...).Scan(&totalPosts)
	}
	if err != nil {
		return posts, totalPosts, err
	}
	log.Infof("Found %d posts", totalPosts)

	where := " posts.category IN (%s) "
	if tag != "" {
		where += " AND tags.tag = ? "
	}

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
        `+where+`
    GROUP BY
        posts.id
    LIMIT ? 
    OFFSET ?`, inClause)

	params = append(params, limit)
	params = append(params, offset)
	log.Info("Searching with limit & offset of", "limit", limit, "offset", offset)

	rows, err := db.Query(query, params...)
	if err != nil {
		return posts, totalPosts, err
	}

	for rows.Next() {
		var post types.Post
		var tags sql.NullString
		err := rows.Scan(&post.Id, &post.Slug, &post.Title, &post.Content, &post.Category, &post.Archived, &post.PublishAt, &post.CreatedAt, &post.UpdatedAt, &tags)

		if err != nil {
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

func getChildrenCategories(categories []types.Category, parent string) []types.Category {
	var cats []types.Category

	for i := 0; i < len(categories); i++ {
		if categories[i].Parent == parent {
			cats = append(cats, categories[i])
			// add all the children as well
			var children = getChildrenCategories(categories, categories[i].Name)
			for j := 0; j < len(children); j++ {
				cats = append(cats, children[j])
			}
		}
	}

	return cats
}