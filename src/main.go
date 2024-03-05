// main.go
package main

import (
	"blog-server/database"
	"blog-server/routes"
	"blog-server/types"
	"blog-server/utils"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

func init() {
	// load .env file
	if err := godotenv.Load(); err != nil {
		log.Warn("No .env file found")
	}

	routes.LoadAuthConfig()
	utils.CreateStore()
}

func main() {
	log.SetLevel(log.DebugLevel)
	// set up database
	database.Connect()
	// set up router with auth middleware
	r := mux.NewRouter()
	r.Use(AuthMiddleware)
	// posts
	r.HandleFunc("/posts", routes.FetchPosts).Methods("GET")
	r.HandleFunc("/posts/{category}", routes.FetchPosts).Methods("GET")
	// individual post functions
	r.HandleFunc("/post/{slug}", routes.GetPostBySlug).Methods("GET")
	r.HandleFunc("/post/new", routes.CreatePost).Methods("POST")
	r.HandleFunc("/post/edit", routes.EditPost).Methods("PUT")
	r.HandleFunc("/post/delete/{id}", routes.DeletePost).Methods("DELETE")
	// category
	r.HandleFunc("/categories", routes.GetCategories).Methods("GET")
    r.HandleFunc("/category/new", routes.CreateCategory).Methods("POST")
    r.HandleFunc("/category/delete/{name}", routes.DeleteCategory).Methods("DELETE")
	// tags
	r.HandleFunc("/post/tag", routes.AddPostTag).Methods("PUT")
	r.HandleFunc("/post/tag", routes.DeletePostTag).Methods("DELETE")
	r.HandleFunc("/tags", routes.GetTags).Methods("GET")
	// auth
	r.HandleFunc("/auth/user", routes.GetUserInfo).Methods("GET")
	r.HandleFunc("/auth/github/login", routes.GithubLogin).Methods("GET")
	r.HandleFunc("/auth/github/callback", routes.GithubCallback).Methods("GET")
	r.HandleFunc("/auth/logout", routes.Logout).Methods("GET")
    // api tokens
    r.HandleFunc("/tokens", routes.GetUserTokens).Methods("GET")
    r.HandleFunc("/token/new", routes.CreateToken).Methods("POST")
    r.HandleFunc("/token/edit", routes.EditToken).Methods("PUT")
    r.HandleFunc("/token/delete/{id}", routes.DeleteToken).Methods("DELETE")
    // integrations
    r.HandleFunc("/links", routes.GetUserIntegrations).Methods("GET")
    r.HandleFunc("/links/upsert", routes.UpsertIntegrations).Methods("PUT")

	// modify cors
	c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:5173", "https://f0rbit.github.io", "https://blog.forbit.dev", "http://blog.forbit.dev", "blog.forbit.dev"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	})

	// Start the server
	port := os.Getenv("PORT")
	log.Info("Server started", "port", port)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: c.Handler(r),
	}

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Info("Stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Info("Graceful shutdown complete.")
}

var EXEMPT_URL = []string{"/auth/github/login", "/auth/logout", "/auth/test", "/auth/user", "/auth/github/callback"}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *types.User

        // first check for an "API_TOKEN" in the headers
        auth_token := r.Header.Get(routes.AUTH_HEADER)
        if auth_token != "" {
            user, err := database.GetUserByToken(auth_token);
            if err != nil {
                utils.LogError("Error fetching user by token", err, http.StatusNotFound, w);
                return;
            }
            if user != nil {
                ctx := context.WithValue(r.Context(), "user", user);
                next.ServeHTTP(w, r.WithContext(ctx));
                return;
            } else {
                http.Error(w, "Invalid token", http.StatusUnauthorized);
                return;
            }
        } 

		// Retrieve the user session
		session, err := utils.GetStore().Get(r, "user-session")

		if err != nil {
			utils.LogError("Error obtaining session", err, http.StatusInternalServerError, w)
			return
		}

		if userID, ok := session.Values["user_id"].(int); ok {
			user, err = database.GetUserByID(userID)
			if err != nil {
				utils.LogError("User not found", err, http.StatusNotFound, w)
				return
			}
		} else {
			// check if the path is exempt from auth check
			for _, url := range EXEMPT_URL {
				if url == r.URL.Path {
					ctx := context.WithValue(r.Context(), "user", user)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}
			// otherwise return 401
            utils.Unauthorized(w);
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
