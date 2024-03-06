package database

import (
	"blog-server/types"
	"database/sql"
	"errors"
)

// Example function to create a new user in the database
func CreateUser(githubUser *types.GitHubUser) (*types.User, error) {
	// Insert user details into 'users' table
	result, err := db.Exec(`
		INSERT INTO users (github_id, username, email, avatar_url)
		VALUES (?, ?, ?, ?)`,
		githubUser.ID, githubUser.Login, githubUser.Email, githubUser.AvatarURL,
	)
	if err != nil {
		return nil, err
	}

	// Retrieve the inserted user's ID
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	/** @todo - create root category */

	// Return the created user
	return &types.User{
		ID:        int(userID),
		GitHubID:  githubUser.ID,
		Username:  githubUser.Login,
		Email:     githubUser.Email,
		AvatarURL: githubUser.AvatarURL,
	}, nil
}

// getUserByGitHubID retrieves a user from the database by their GitHub ID.
func GetUserByGitHubID(githubID int) (*types.User, error) {
	// Execute a SQL query to fetch the user by GitHub ID
	row := db.QueryRow("SELECT * FROM users WHERE github_id = ?", githubID)

	// Create a User struct to hold the fetched data
	var user types.User

	// Scan the row's values into the user struct
	err := row.Scan(&user.ID, &user.GitHubID, &user.Username, &user.Email, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		// Handle the case when no user with the given GitHub ID is found
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err // Other database error
	}

	return &user, nil
}

func GetUserByID(userID int) (*types.User, error) {
	row := db.QueryRow("SELECT * FROM users WHERE user_id = ?", userID)
	var user types.User

	err := row.Scan(&user.ID, &user.GitHubID, &user.Username, &user.Email, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err // Other database error
	}

	return &user, nil
}

func GetUserByToken(token string) (*types.User, error) {
	row := db.QueryRow("SELECT access_keys.user_id FROM access_keys WHERE key_value = ?", token)

	var userID int

	err := row.Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err // Other database error
	}

	return GetUserByID(userID)
}

func GetTokens(userID int) ([]types.AccessKey, error) {
	var tokens []types.AccessKey
	rows, err := db.Query("SELECT * FROM access_keys WHERE user_id = ?", userID)

	if err != nil {
		return tokens, err
	}

	for rows.Next() {
		var token types.AccessKey
		err := rows.Scan(&token.ID, &token.Value, &token.UserID, &token.Name, &token.Note, &token.Enabled, &token.CreatedAt, &token.UpdatedAt)
		if err != nil {
			return tokens, err
		}
		tokens = append(tokens, token)
	}
    if tokens == nil {
        tokens = make([]types.AccessKey, 0)
    }
	return tokens, err
}

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
