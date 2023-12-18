// main.go
package main

import (
	"blog-server/routes"
	"context"
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var database = "db/sqlite.db"

func main() {
	initLogging()
	initDatabase()

	// Initialize routes
	r := mux.NewRouter()
	r.HandleFunc("/posts", routes.GetPosts).Methods("GET")
	r.HandleFunc("/posts/{category}", routes.GetPostsByCategory).Methods("GET")
	r.HandleFunc("/post/{id}", routes.GetPostByID).Methods("GET")
	r.HandleFunc("/post/new", routes.CreatePost).Methods("POST")
	r.HandleFunc("/post/edit", routes.EditPost).Methods("PUT")
	r.HandleFunc("/post/delete/{id}", routes.DeletePost).Methods("DELETE")
	r.HandleFunc("/categories", routes.GetCategories).Methods("GET")

	// Start the server
	port := ":8080"
	log.Printf("Server started on port %s", port)
	server := &http.Server{
		Addr: ":8080",
	}
	//log.Fatal(http.ListenAndServe(port, r))

	go func() {
		server.Handler = r
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Graceful shutdown complete.")
}

func initLogging() {
	// Create a log file
	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error creating log file: ", err)
	}
	file.Write([]byte("## NEW INSTANCE ##\n"))

	// Set log output to both console and file
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func initDatabase() {
	db, err := sql.Open("sqlite3", database)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	var version string
	err = db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Database Version: %s", version)
	if !tableExists(db, "categories") {
		createCategories(db)
	}
	if !tableExists(db, "posts") {
		createPosts(db)
	}
}

func tableExists(db *sql.DB, tableName string) bool {
	// Query to check if a table exists in SQLite
	query := "SELECT name FROM sqlite_master WHERE type='table' AND name=?"
	var name string
	err := db.QueryRow(query, tableName).Scan(&name)
	return err == nil
}

func createCategories(db *sql.DB) {
	// create "Categories" table
	statement, err := db.Prepare("DROP TABLE IF EXISTS categories")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Dropped categories table")
	}
	statement.Exec()

	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS categories (name VARCHAR(15) PRIMARY KEY, parent VARCHAR(15))")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Created categories table")
	}
	statement.Exec()

	statement, _ = db.Prepare("INSERT INTO categories (name, parent) VALUES (?, ?)")
	statement.Exec("coding", "root")
	statement.Exec("learning", "coding")
	statement.Exec("devlog", "coding")
	statement.Exec("gamedev", "devlog")
	statement.Exec("webdev", "devlog")
	statement.Exec("code-story", "coding")
	statement.Exec("hobbies", "root")
	statement.Exec("photography", "hobbies")
	statement.Exec("painting", "hobbies")
	statement.Exec("hiking", "hobbies")
	statement.Exec("story", "root")
	statement.Exec("advice", "root")
	log.Printf("Inserted categories")
}

func createPosts(db *sql.DB) {
	statement, err := db.Prepare("DROP TABLE IF EXISTS posts")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Dropped posts table")
	}
	statement.Exec()

	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY AUTOINCREMENT, slug TEXT NOT NULL UNIQUE, title TEXT NOT NULL, content TEXT NOT NULL, category TEXT NOT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		log.Fatal(err)
	} else {
		statement.Exec()
		log.Printf("Created posts table")
	}

	statement, err = db.Prepare("INSERT INTO posts (slug, title, content, category) VALUES (?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	} else {
		statement.Exec("test-post", "test", "this is a test post, first post.", "coding")
		log.Printf("Inserted 'test-post'")
	}

}
