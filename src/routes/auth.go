package routes

import (
	"blog-server/database"
	"blog-server/types"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/charmbracelet/log"
)

var authConfig *oauth2.Config

func LoadAuthConfig() {
	authConfig = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT"),
		ClientSecret: os.Getenv("GITHUB_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/github/callback",
		Scopes:       []string{"read:user"},
		Endpoint:     github.Endpoint,
	}
}

const AUTH_HEADER = "Auth-Token"

var AUTH_TOKEN = os.Getenv("AUTH_TOKEN")

func TryToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token != AUTH_TOKEN {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func GithubLogin(w http.ResponseWriter, r *http.Request) {
	url := authConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GithubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := authConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Use token to access user's GitHub data
	log.Info("Token: " + token.AccessToken)
    handleGitHubUser(token);
}


// Example logic to handle user registration or retrieval
func handleGitHubUser(token *oauth2.Token) (*types.User, error) {
	// Fetch user details from GitHub
	githubUser, err := fetchGitHubUserDetails(token)
	if err != nil {
		return nil, err
	}

	// Check if the user already exists in the database
	user, err := database.GetUserByGitHubID(githubUser.ID);
	if err != nil {
		return nil, err
	}

	// If the user doesn't exist, create a new user record
	if user == nil {
		user, err = database.CreateUser(githubUser)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func fetchGitHubUserDetails(token *oauth2.Token) (*types.GitHubUser, error) {
	// Create an HTTP client with the OAuth2 token
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))

	// Make a GET request to the GitHub API to retrieve user details
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Invalid status code")
	}

	// Parse the JSON response into a GitHubUser struct
	var githubUser types.GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, err
	}

	return &githubUser, nil
}

