package fifo

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// InitDB initializes and verifies a Postgres connection
func InitDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")

	// Build DSN (disable SSL for local/dev)
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, pass, host, name)

	// Open connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Verify DB is reachable
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
