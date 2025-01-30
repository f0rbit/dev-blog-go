package database

import (
	"blog-server/types"
	"crypto/rand"
	"encoding/hex"

	"github.com/charmbracelet/log"
)

// randToken generates a random hex value.
func randToken(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func CreateToken(token types.AccessKey) (int, error) {
	var err error
	// generate key value
	value, err := randToken(24)
	if err != nil {
		return -1, err
	}

	_, err = db.Exec(
		`INSERT INTO access_keys (key_value, user_id, name, note, enabled) VALUES (?,?,?,?,?)`,
		value,
		token.UserID,
		token.Name,
		token.Note,
		token.Enabled)
	if err != nil {
		return -1, err
	}
	log.Info("Created new token", "user_id", token.UserID, "value", value)
	var id int
	row := db.QueryRow("SELECT last_insert_rowid()")
	err = row.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func GetToken(id int) (types.AccessKey, error) {
	var token types.AccessKey

	row := db.QueryRow("SELECT * FROM access_keys WHERE key_id = ?", id)

	err := row.Scan(&token.ID, &token.Value, &token.UserID, &token.Name, &token.Note, &token.Enabled, &token.CreatedAt, &token.UpdatedAt)
	if err != nil {
		return token, err
	}
	return token, nil
}

func UpdateToken(token types.AccessKey) error {
	_, err := db.Exec(`
    UPDATE
        access_keys
    SET
        name = ?,
        note = ?,
        enabled = ?,
        updated_at = CURRENT_TIMESTAMP
    WHERE
        key_id = ? AND
        user_id = ?;
    `, token.Name, token.Note, token.Enabled, token.ID, token.UserID)

	if err != nil {
		return err
	}
	log.Info("Updated token", "id", token.ID)
	return nil
}

func DeleteToken(tokenID int) error {
	_, err := db.Exec("DELETE FROM access_keys WHERE key_id = ?", tokenID)
	if err == nil {
		log.Info("Delete token", "id", tokenID)
	}
	return err
}

func GetProjectCache(cacheID int) (types.ProjectCache, error) {
	var cache types.ProjectCache
	err := db.QueryRow("SELECT * FROM projects_cache WHERE id = ?", cacheID).Scan(&cache)
	if err != nil {
		return cache, err
	}
	return cache, nil
}

func GetLatestProjectCache(userID int) (types.ProjectCache, error) {
	var cache types.ProjectCache
	err := db.QueryRow("SELECT * FROM projects_cache WHERE user_id = ? ORDER BY fetched_at DESC LIMIT 1", userID).Scan(&cache)
	if err != nil {
		return cache, err
	}
	return cache, nil
}

func InsertProjectCache(cache types.ProjectCache) (int, error) {
	var id int
	err := db.QueryRow("INSERT INTO projects_cache (user_id, status, url, data, fetched_at) VALUES (?, ?, ?, ?, ?)", cache.UserID, cache.Status, cache.URL, cache.Data, cache.FetchedAt).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func FailProjectCache(id int) {
	_, err := db.Exec("UPDATE projects_cache SET status = 'failed' WHERE id = ?", id)
	if err != nil {
		log.Error("Failed to set project cache as failed", "id", id)
	}
}

func UpdateProjectCacheData(id int, data []byte) error {
	_, err := db.Exec("UPDATE projects_cache SET status = 'fetched' data = ? WHERE id = ?", data, id)
	if err != nil {
		return err
	}
	return nil
}

func GetProject(userID int, projectID string) (types.Project, error) {
	var project types.Project
	err := db.QueryRow("SELECT * FROM projects_cache WHERE user_id = ? AND project_id = ?", userID, projectID).Scan(&project)
	if err != nil {
		return project, err
	}
	return project, nil
}

func SetProjectKey(userID int, key string) error {
	_, err := db.Exec("UPDATE devpad_api_tokens SET token = ?, updated_at = CURRENT_TIMESTAMP WHERE user_id = ?", key, userID)
	if err != nil {
		return err
	}
	return nil
}

func GetProjectKey(userID int) (string, error) {
	var key string
	err := db.QueryRow("SELECT token FROM devpad_api_tokens WHERE user_id = ?", userID).Scan(&key)
	if err != nil {
		return "", err
	}
	return key, nil
}
