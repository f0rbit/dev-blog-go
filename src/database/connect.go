package database

import (
	"database/sql"
	"os"

    "github.com/charmbracelet/log"
)

var db *sql.DB;

func Connect() {
    var DB_URL = os.Getenv("DATABASE");
    connection, err := sql.Open("sqlite3", DB_URL + "?parseTime=true");
    if err != nil {
        log.Fatal("Failed to open connection to DB.", "err", err);
        return
    }

    connection.SetMaxOpenConns(10);
    connection.SetMaxIdleConns(5);

    err = connection.Ping()
    if err != nil {
        log.Fatal("Failed to ping DB connection.", "err", err);
        return
    }

    db = connection
}

func Connection() (*sql.DB) {
    return db;
}
