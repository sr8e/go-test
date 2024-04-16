package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var dbPool *sql.DB

func init() {
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		log.Fatal("environment variable DB_URL not set")
	} else {
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("could not open db pool: %s", err)
		} else {
			dbPool = db
		}
	}
}
