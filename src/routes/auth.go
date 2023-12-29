package routes

import (
	"context"
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
}
