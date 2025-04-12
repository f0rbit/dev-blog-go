package database

import (
	"blog-server/types"
	"blog-server/utils"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
)

type Identifier string

const (
	ID   Identifier = "id"
	Slug Identifier = "slug"
)

func FetchPost(user *types.User, identifier Identifier, needle interface{}) (types.Post, error) {
	var post types.Post
	if user == nil {
		return post, errors.New("No user specified")
	}
	var where string
	if identifier == "id" {
		where = "posts.id = ?"
	} else if identifier == "slug" {
		where = "posts.slug = ?"
	}

	var base = `
    SELECT
        posts.id,
        posts.author_id,
        posts.slug,
        posts.title,
        posts.description,
        posts.content,
        posts.format,
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
        posts.author_id = ? AND
        ` + where + `;
    `
	var tags sql.NullString

	err := db.QueryRow(base, user.ID, needle).Scan(
		&post.Id,
		&post.AuthorID,
		&post.Slug,
		&post.Title,
		&post.Description,
		&post.Content,
		&post.Format,
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

	post.Description = utils.GetDescription(post.Content)

	// check for project_id link
	post.ProjectID = GetPostProjectID(post.Id)

	return post, nil
}

// author_id should be inside the post object
func CreatePost(post types.Post) (int, error) {
	var err error
	// Insert the new post into the database
	_, err = db.Exec(
		`INSERT INTO posts (author_id, slug, title, description, content, format, category, publish_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		post.AuthorID,
		post.Slug,
		post.Title,
		post.Description,
		post.Content,
		post.Format,
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
	// insert project_id link
	if post.ProjectID != "" {
		_, err = db.Exec("INSERT INTO posts_projects (post_id, project_id) VALUES (?, ?)", post.Id, post.ProjectID)
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
	// delete project_id link
	_, err = db.Exec("DELETE FROM posts_projects WHERE post_id = ?", id)
	if err != nil {
		return err
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
        description = ?,
        content = ?,
        format = ?,
        category = ?,
        archived = ?,
        publish_at = ? 
    WHERE 
        id = ?`,
		updatedPost.Slug,
		updatedPost.Title,
		updatedPost.Description,
		updatedPost.Content,
		updatedPost.Format,
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

	// update project_id link
	err = UpdatePostProjectID(updatedPost.Id, updatedPost.ProjectID)
	log.Info("Updated Post", "id", updatedPost.Id)
	return err
}

func GetPosts(user *types.User, category, tag string, limit, offset int, project_uuid string) ([]types.Post, int, error) {
	var posts []types.Post
	var totalPosts int

	log.Info("Searching for posts", "category", category, "tag", tag, "project_uuid", project_uuid)

	// Step 1: Get search categories
	search_categories, err := getSearchCategories(user, category)
	if err != nil {
		return posts, totalPosts, errors.Join(errors.New("error getting categories"), err)
	}

	log.Info("Searching through categories", "searchCategories", search_categories)

	// Step 2: Build WHERE clause and parameters
	where_clause, params := buildWhereClause(user.ID, search_categories, tag, project_uuid)

	// Step 3: Get total post count
	totalPosts, err = getTotalPostsCount(where_clause, params, tag)
	if err != nil {
		return posts, totalPosts, errors.Join(errors.New("error getting total posts"), err)
	}

	log.Infof("Found %d posts", totalPosts)

	// Step 4: Fetch paginated posts
	posts, err = fetchPaginatedPosts(where_clause, params, limit, offset)
	if err != nil {
		return posts, totalPosts, errors.Join(errors.New("error fetching posts"), err)
	}

	return posts, totalPosts, nil
}

func getSearchCategories(user *types.User, category string) ([]string, error) {
	var search_categories []string
	categories, err := GetCategories(user)
	if err != nil {
		return nil, err
	}

	search_categories = append(search_categories, category)
	if category == "root" || category == "" {
		for _, cat := range categories {
			search_categories = append(search_categories, cat.Name)
		}
	} else {
		children := utils.GetChildrenCategories(categories, category)
		for _, child := range children {
			search_categories = append(search_categories, child.Name)
		}
	}

	return search_categories, nil
}

func buildWhereClause(authorID int, categories []string, tag, project_uuid string) (string, []any) {
	placeholders := make([]string, len(categories))
	for i := range categories {
		placeholders[i] = "?"
	}
	inClause := strings.Join(placeholders, ",")

	var params []any
	params = append(params, authorID)
	for _, c := range categories {
		params = append(params, c)
	}

	where := "posts.author_id = ? AND posts.category IN (" + inClause + ")"

	if tag != "" {
		where += " AND tags.tag = ?"
		params = append(params, tag)
	}

	if project_uuid != "" {
		where += " AND posts_projects.project_uuid = ?"
		params = append(params, project_uuid)
	}

	return where, params
}

func getTotalPostsCount(where string, params []any, tag string) (int, error) {
	var total int
	var query string

	if tag == "" {
		query = "SELECT COUNT(*) FROM posts WHERE " + where
	} else {
		query = "SELECT COUNT(*) FROM posts LEFT JOIN tags ON tags.post_id = posts.id WHERE " + where
	}

	log.Info("Counting posts", "query", query, "params", params)

	err := db.QueryRow(query, params...).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func fetchPaginatedPosts(where string, params []any, limit, offset int) ([]types.Post, error) {
	var posts []types.Post

	query := `SELECT 
		posts.id, 
		posts.author_id,
		posts.slug, 
		posts.title, 
		posts.description,
		posts.content, 
		posts.format,
		posts.category, 
		posts.archived,
		posts.publish_at,
		posts.created_at, 
		posts.updated_at,
		GROUP_CONCAT(tags.tag) AS tags,
		IFNULL(posts_projects.project_uuid, '') AS project_uuid
	FROM 
		posts 
	LEFT JOIN
		tags ON posts.id = tags.post_id
	LEFT JOIN
		posts_projects ON posts.id = posts_projects.post_id
	WHERE 
		` + where + `
	GROUP BY
		posts.id
	LIMIT ? 
	OFFSET ?`

	params = append(params, limit, offset)

	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post types.Post
		var tags sql.NullString
		var project_uuid string

		err := rows.Scan(&post.Id, &post.AuthorID, &post.Slug, &post.Title, &post.Description, &post.Content, &post.Format, &post.Category, &post.Archived, &post.PublishAt, &post.CreatedAt, &post.UpdatedAt, &tags, &project_uuid)
		if err != nil {
			return nil, err
		}

		if tags.Valid {
			post.Tags = strings.Split(tags.String, ",")
		} else {
			post.Tags = []string{}
		}

		if post.Description == "" {
			post.Description = utils.GetDescription(post.Content)
		}

		post.ProjectID = project_uuid
		posts = append(posts, post)
	}

	if posts == nil {
		posts = make([]types.Post, 0)
	}

	return posts, nil
}

func RemoveCategoryFromPosts(user *types.User, cat_list []string) error {
	params := make([]any, len(cat_list)+1)
	params[0] = user.ID
	for i, c := range cat_list {
		params[i+1] = c
	}
	// Build the IN clause with placeholders
	placeholders := make([]string, len(cat_list))
	for i := range cat_list {
		placeholders[i] = "?"
	}
	inClause := strings.Join(placeholders, ",")

	query := fmt.Sprintf("UPDATE posts SET category = 'root' WHERE author_id = ? AND category IN (%s)", inClause)
	log.Info("Removing category from posts", "query", query, "params", params)
	_, err := db.Exec(query, params...)
	return err
}

func GetPostsByTitle(user *types.User, title string) (*types.Post, error) {
	var post types.Post
	var tags sql.NullString
	err := db.QueryRow(`
    SELECT 
        posts.id, 
        posts.author_id,
        posts.slug, 
        posts.title, 
        posts.description,
        posts.content, 
        posts.format,
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
        posts.author_id = ? AND posts.title = ?
    GROUP BY
        posts.id`, user.ID, title).Scan(
		&post.Id,
		&post.AuthorID,
		&post.Slug,
		&post.Title,
		&post.Description,
		&post.Content,
		&post.Format,
		&post.Category,
		&post.Archived,
		&post.PublishAt,
		&post.CreatedAt,
		&post.UpdatedAt,
		&tags)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if tags.Valid {
		post.Tags = strings.Split(tags.String, ",")
	} else {
		post.Tags = []string{}
	}

	post.Description = utils.GetDescription(post.Content)

	// check for project_id link
	post.ProjectID = GetPostProjectID(post.Id)

	return &post, nil
}

func GetPostProjectID(postID int) string {
	var projectID string
	err := db.QueryRow("SELECT project_uuid FROM posts_projects WHERE post_id = ?", postID).Scan(&projectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ""
		} else {
			log.Error("Error fetching post to project ID", "err", err)
			return ""
		}
	}
	return projectID
}

func UpdatePostProjectID(postID int, projectID string) error {
	_, err := db.Exec("DELETE FROM posts_projects WHERE post_id = ?", postID)
	if err != nil {
		return err
	}
	if projectID != "" {
		_, err = db.Exec("INSERT INTO posts_projects (post_id, project_uuid) VALUES (?, ?)", postID, projectID)
		if err != nil {
			return err
		}
	}
	return nil
}
