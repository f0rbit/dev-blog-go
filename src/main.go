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

func main() {
	log.SetLevel(log.DebugLevel)

	// Initialize routes
	r := mux.NewRouter()
	r.HandleFunc("/posts", routes.GetPosts).Methods("GET")
	r.HandleFunc("/posts/{category}", routes.GetPostsByCategory).Methods("GET")
	r.HandleFunc("/post/{slug}", routes.GetPostBySlug).Methods("GET")
	r.HandleFunc("/post/new", routes.CreatePost).Methods("POST")
	r.HandleFunc("/post/edit", routes.EditPost).Methods("PUT")
	r.HandleFunc("/post/delete/{id}", routes.DeletePost).Methods("DELETE")
	r.HandleFunc("/categories", routes.GetCategories).Methods("GET")

    c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:5173"},
        AllowCredentials: true,
    })

	// Start the server
	port := ":8080"
	server := &http.Server{
		Addr: ":8080",
        Handler: c.Handler(r),
	}
	log.Info("Server started", "port", port)

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
