package fifo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
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

	// Retry loop: try up to 10 times, waiting 3s between attempts
	for i := 0; i < 10; i++ {
		if err := db.Ping(); err == nil {
			log.Println("connected to database")
			return db
		}
		log.Println("waiting for database to be ready...")
		time.Sleep(3 * time.Second)
	}

	log.Fatal("could not connect to database after retries")
	return nil
}
