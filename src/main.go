// main.go
package main

import (
	"blog-server/routes"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gorilla/mux"
    "github.com/rs/cors"
	_ "github.com/mattn/go-sqlite3"
)

var AUTH_TOKEN = os.Getenv("AUTH_TOKEN")

func main() {
	log.SetLevel(log.DebugLevel)

	r := mux.NewRouter()
    // Middleware for checking 'Auth-Token' header
	authMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				token := r.Header.Get("Auth-Token")
				if token != AUTH_TOKEN {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}

	// Apply the middleware to all routes
	r.Use(authMiddleware)

    // posts
	r.HandleFunc("/posts", routes.GetPosts).Methods("GET")
	r.HandleFunc("/posts/{category}", routes.GetPostsByCategory).Methods("GET")
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

    c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:5173"},
        AllowCredentials: true,
    })

	// Start the server
	port := os.Getenv("PORT") 
	log.Infof("Server started on port %s", port)
	server := &http.Server{
		Addr: ":" + port,
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
