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
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.SetLevel(log.DebugLevel)

	// Initialize routes
	r := mux.NewRouter()
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

	// Start the server
	port := ":8080"
	log.Infof("Server started on port %s", port)
	server := &http.Server{
		Addr: ":8080",
	}
	//log.Fatal(http.ListenAndServe(port, r))

	go func() {
		server.Handler = r
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
