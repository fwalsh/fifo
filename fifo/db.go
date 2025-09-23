package fifo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func InitDB() *sql.DB {
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "items_user")
	pass := getEnv("DB_PASS", "items_pass")
	dbname := getEnv("DB_NAME", "items_db")

	connStr := fmt.Sprintf(
		"host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		host, user, pass, dbname,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("failed to connect to db:", err)
	}

	// Retry settings: configurable by env
	maxRetries, _ := strconv.Atoi(getEnv("DB_MAX_RETRIES", "10"))
	retryDelay, _ := strconv.Atoi(getEnv("DB_RETRY_DELAY", "3")) // seconds

	for i := 0; i < maxRetries; i++ {
		if err := db.Ping(); err == nil {
			log.Println("connected to database")
			return db
		}
		log.Printf("waiting for database... retry %d/%d", i+1, maxRetries)
		time.Sleep(time.Duration(retryDelay) * time.Second)
	}

	log.Fatal("could not connect to database after retries")
	return nil
}
