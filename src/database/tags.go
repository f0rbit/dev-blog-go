package database

func GetTags() ([]string, error) {
	var tags []string
	rows, err := db.Query("SELECT DISTINCT tag FROM tags")
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
