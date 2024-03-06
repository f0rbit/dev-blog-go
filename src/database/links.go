package database

import (
	"blog-server/types"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/charmbracelet/log"
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

// need a new type to handle the subquery
type FetchLink struct {
	PostID     int    `json:"post_id"`
	Identifier string `json:"identifier"`
}
type IntegrationWithRawLinks struct {
	types.Integration `json:",inline"`
	FetchLinks        string `json:"fetch_links"`
}
type IntegrationWithLinks struct {
    types.Integration `json:",inline"`
    FetchLinks        []FetchLink `json:"fetch_links"`
}

func GetUserIntegrations(userID int) ([]IntegrationWithLinks, error) {
	var integrations []IntegrationWithLinks
	// old query -- rows, err := db.Query("SELECT * FROM fetch_queue WHERE user_id = ?", userID)
	// in the new query we want to include the fetch_links as well as a property "fetch_links" under each integration as a subquery
    // note sqlite3 is being used
    rows, err := db.Query("SELECT fetch_queue.*, (SELECT json_group_array(json_object('post_id', fetch_links.post_id, 'identifier', fetch_links.identifier)) FROM fetch_links WHERE fetch_links.fetch_source = fetch_queue.id) as fetch_links FROM fetch_queue WHERE user_id = ?", userID)

	if err != nil {
		return integrations, err
	}
	defer rows.Close()


    
	for rows.Next() {
		var integration IntegrationWithRawLinks
		err := rows.Scan(&integration.ID, &integration.UserID, &integration.LastFetch, &integration.Location, &integration.Source, &integration.Data, &integration.CreatedAt, &integration.UpdatedAt, &integration.FetchLinks)
		if err != nil {
			return nil, err // Updated for clarity
		}

		// Additional handling for the FetchLinks field
		var fetchLinks []FetchLink
		if err := json.Unmarshal([]byte(integration.FetchLinks), &fetchLinks); err == nil {
			integration.FetchLinks = "" // Resetting to empty if necessary, or handle appropriately
			// Assume you have a way to associate these back if needed, or change the struct accordingly.
		} else {
			log.Error("Error unmarshalling fetch links", "error", err)
            continue
        }

        // copy integration into integrationWithLinks, getting FetchLinks from fetchLinks
        integrationWithLinks := IntegrationWithLinks{
            Integration: integration.Integration,
            FetchLinks: fetchLinks,
        }

        integrations = append(integrations, integrationWithLinks)
	}

	return integrations, nil
}

func GetIntegrationBySource(userID int, source string) (*types.Integration, error) {
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

func GetIntegrationByID(id int) (*types.Integration, error) {
	row := db.QueryRow("SELECT * FROM fetch_queue WHERE id = ?", id)
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

func GetFetchLinkBySlug(linkID int, slug string) (*types.FetchLink, error) {
	row := db.QueryRow("SELECT * FROM fetch_links WHERE fetch_source = ? AND identifier = ?", linkID, slug)
	var link types.FetchLink
	err := row.Scan(&link.PostID, &link.FetchSource, &link.Identifier, &link.CreatedAt, &link.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &link, nil
}

func CreateFetchLink(link types.FetchLink) error {
	_, err := db.Exec("INSERT INTO fetch_links (post_id, fetch_source, identifier) VALUES (?, ?, ?)", link.PostID, link.FetchSource, link.Identifier)
	log.Info("Created fetch link", "link", link)
	return err
}

func DeleteIntegration(id int) error {
	_, err := db.Exec("DELETE FROM fetch_queue WHERE id = ?", id)
	// delete all fetch_links as well
	_, err = db.Exec("DELETE FROM fetch_links WHERE fetch_source = ?", id)
	return err
}

func SetIntegrationLastFetched(linkID int) error {
	_, err := db.Exec("UPDATE fetch_queue SET last_fetch = CURRENT_TIMESTAMP WHERE id = ?", linkID)
	return err
}
