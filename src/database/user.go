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
    row := db.QueryRow("SELECT access_keys.user_id FROM access_keys WHERE key_value = ?", token);

    var userID int;

    err := row.Scan(&userID) 
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err // Other database error
    }

    return GetUserByID(userID)
}
