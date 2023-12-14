// main.go
package main

import (
	"blog-server/routes"
    "blog-server/types"
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var database = "../db/sqlite.db"

func main() {
	initLogging()
	initDatabase()

	// Initialize routes
	r := mux.NewRouter()
	// r.HandleFunc("/posts", GetPosts).Methods("GET")
    // r.HandleFunc("/posts/{category}", GetPostsByCategory).Methods("GET")
	// r.HandleFunc("/post/{id}", GetPostByID).Methods("GET")
	// r.HandleFunc("/post/new", CreatePost).Methods("POST")
	// r.HandleFunc("/post/edit", EditPost).Methods("PUT")
	r.HandleFunc("/categories", routes.GetCategories).Methods("GET")

	// Start the server
	port := ":8080"
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}

func initLogging() {
	// Create a log file
	file, err := os.OpenFile("../logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error creating log file: ", err)
	}
	defer file.Close()

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
        seedDatabase()
    }
}

func tableExists(db *sql.DB, tableName string) bool {
	// Query to check if a table exists in SQLite
	query := "SELECT name FROM sqlite_master WHERE type='table' AND name=?"
	var name string
	err := db.QueryRow(query, tableName).Scan(&name)
	return err == nil
}

func seedDatabase() {
    // create "Categories" table
    db, err := sql.Open("sqlite3", database)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

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
    log.Printf("Inserted categories");

    rows, _ := db.Query("SELECT name, parent FROM categories")
	var cat types.Category
	for rows.Next() {
		rows.Scan(&cat.Name, &cat.Parent)
        log.Printf("Name: %s, Parent: %s", cat.Name, cat.Parent)
	}
}
