package database

import (
	"blog-server/types"
	"database/sql"
	"errors"
)

func UpsertIntegration(integration types.Integration) error {
    // check if integration exists in fetch_queue
    row := db.QueryRow("SELECT * FROM fetch_queue WHERE user_id = ? AND source = ?", integration.UserID, integration.Source)
    var existing types.Integration
    err := row.Scan(&existing.ID, &existing.UserID, &existing.LastFetch, &existing.Location, &existing.Source, &existing.Data, &existing.CreatedAt, &existing.UpdatedAt)
    if err != nil {
        if !errors.Is(err, sql.ErrNoRows) {
            return err
        } else {
            // insert new integration
            _, err = db.Exec("INSERT INTO fetch_queue (user_id, last_fetch, location, source, data) VALUES (?, ?, ?, ?, ?)", integration.UserID, integration.LastFetch, integration.Location, integration.Source, integration.Data)
            return err
        }
    }
    // update existing integration
    _, err = db.Exec("UPDATE fetch_queue SET location = ?, data = ? WHERE id = ?", integration.Location, integration.Data, existing.ID)
    return err    
}

func GetUserIntegrations(userID int) ([]types.Integration, error) {
    var integrations []types.Integration
    rows, err := db.Query("SELECT * FROM fetch_queue WHERE user_id = ?", userID)

    if err != nil {
        return integrations, err
    }

    for rows.Next() {
        var integration types.Integration
        err := rows.Scan(&integration.ID, &integration.UserID, &integration.LastFetch, &integration.Location, &integration.Source, &integration.Data, &integration.CreatedAt, &integration.UpdatedAt)
        if err != nil {
            return integrations, err
        }
        integrations = append(integrations, integration)
    }
    if integrations == nil {
        integrations = make([]types.Integration, 0)
    }
    return integrations, err
}

func GetIntegration(userID int, source string) (*types.Integration, error) {
    row := db.QueryRow("SELECT * FROM fetch_queue WHERE user_id = ? AND source = ?", userID, source)
    var integration types.Integration
    err := row.Scan(&integration.ID, &integration.UserID, &integration.LastFetch, &integration.Location, &integration.Source, &integration.Data, &integration.CreatedAt, &integration.UpdatedAt)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return &integration, nil
}
