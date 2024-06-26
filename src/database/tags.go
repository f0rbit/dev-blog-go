package database

import (
	"blog-server/types"
	"errors"
)

func GetTags(user *types.User) ([]string, error) {
	var tags []string
	if user == nil {
		return tags, errors.New("Invalid user reference")
	}
	rows, err := db.Query("SELECT DISTINCT tags.tag FROM posts LEFT JOIN tags ON tags.post_id = posts.id WHERE posts.author_id = ? AND tags.tag IS NOT NULL", user.ID)
	if err != nil {
		return tags, err
	}

	for rows.Next() {
		var tag string
		err := rows.Scan(&tag)
		if err != nil {
			return tags, err
		}
		tags = append(tags, tag)
	}

	return tags, err
}

func CreateTag(postID int, tag string) error {
	_, err := db.Exec("INSERT INTO tags (post_id, tag) VALUES (?, ?)", postID, tag)
	return err
}

func DeleteTag(postID int, tag string) error {
	_, err := db.Exec("DELETE FROM tags WHERE post_id = ? AND tag = ?", postID, tag)
	return err
}
