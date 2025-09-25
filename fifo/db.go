package fifo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// InitDB initializes the Postgres connection and retries if DB is not yet ready.
func InitDB() *sql.DB {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbName)

	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			if pingErr := db.Ping(); pingErr == nil {
				log.Println("✅ Connected to database")
				return db
			} else {
				err = pingErr
			}
		}
		log.Printf("waiting for database... retry %d/10\n", i+1)
		time.Sleep(3 * time.Second)
	}

	log.Fatalf("❌ Could not connect to database after retries: %v", err)
	return nil
}
