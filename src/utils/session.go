package utils

import (
	"os"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func CreateStore() {
	var secret = os.Getenv("COOKIE_SECRET")
	store = sessions.NewCookieStore([]byte(secret))

	store.Options = &sessions.Options{
        Domain:   "localhost",
		Path:     "/",
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
	}
}

func GetStore() *sessions.CookieStore {
	return store
}
