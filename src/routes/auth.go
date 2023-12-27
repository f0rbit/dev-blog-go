package routes

import (
	"net/http"
	"os"
)

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
