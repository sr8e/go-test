package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var dbPool *sql.DB
var redisClient *redis.Client

func init() {
	dbConnStr := os.Getenv("DB_URL")
	if dbConnStr == "" {
		log.Fatal("environment variable DB_URL not set")
	} else {
		db, err := sql.Open("postgres", dbConnStr)
		if err != nil {
			log.Printf("could not open db pool: %s", err)
		} else {
			dbPool = db
		}
	}
	redisConnStr := os.Getenv("REDIS_URL")
	if redisConnStr == "" {
		log.Fatal("environment variable REDIS_URL not set")
	} else {
		redisOpt, err := redis.ParseURL(redisConnStr)
		if err != nil {
			log.Fatalf("could not parse REDIS_URL: %s", err)
		}
		redisClient = redis.NewClient(redisOpt)
	}
}
