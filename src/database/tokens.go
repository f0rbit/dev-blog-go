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
    value, err := randToken(24);
    if err != nil {
        return -1, err
    }

    _, err = db.Exec(
        `INSERT INTO access_keys (key_value, user_id, name, note, enabled) VALUES (?,?,?,?,?)`,
        value,
        token.UserID,
        token.Name,
        token.Note, 
        token.Enabled);
    if err != nil {
        return -1, err;
    }
    log.Info("Created new token", "user_id", token.UserID, "value", value);
    var id int;
    row := db.QueryRow("SELECT last_insert_rowid()")
    err = row.Scan(&id);
    if err != nil {
        return -1, err
    }
    return id, nil;
}

func GetToken(id int) (types.AccessKey, error) {
    var token types.AccessKey;

    row := db.QueryRow("SELECT * FROM access_keys WHERE key_id = ?", id);

    err := row.Scan(&token.ID, &token.Value, &token.UserID, &token.Name, &token.Note, &token.Enabled, &token.CreatedAt, &token.UpdatedAt);
    if err != nil {
        return token, err
    }
    return token, nil;
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
    return nil;
}
