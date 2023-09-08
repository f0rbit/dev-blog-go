// main.go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var db *gorm.DB

func main() {
    // Initialize the database connection
    dsn := "user=postgres password=example dbname=mydb sslmode=disable"
    var err error
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // Create a router using Gorilla Mux
    r := mux.NewRouter()

    // Define your routes and handlers
    r.HandleFunc("/", YourHandler)

    // Start the web server
    fmt.Println("Server is running on port 8080...")
    http.Handle("/", r)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func YourHandler(w http.ResponseWriter, r *http.Request) {
    // Your route handler logic here
}
