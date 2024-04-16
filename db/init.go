package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var dbPool *sql.DB

func init() {
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		log.Fatal("environment variable DB_URL not set")
	} else {
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("could not open db pool: %w", err)
		} else {
			dbPool = db
		}
	}
}
