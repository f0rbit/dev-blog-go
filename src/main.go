// main.go
package main

import (
	"blog-server/database"
	"blog-server/routes"
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

    routes.LoadAuthConfig();
    utils.CreateStore();
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
	// tags
	r.HandleFunc("/post/tag", routes.AddPostTag).Methods("PUT")
	r.HandleFunc("/post/tag", routes.DeletePostTag).Methods("DELETE")
	r.HandleFunc("/tags", routes.GetTags).Methods("GET")
	// auth
	r.HandleFunc("/auth/test", routes.TryToken).Methods("GET")
    r.HandleFunc("/auth/user", routes.GetLogin).Methods("GET")
	r.HandleFunc("/auth/github/login", routes.GithubLogin).Methods("GET")
	r.HandleFunc("/auth/github/callback", routes.GithubCallback).Methods("GET")

	// modify cors
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://f0rbit.github.io"},
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

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}
		if r.Header.Get(routes.AUTH_HEADER) == routes.AUTH_TOKEN {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
