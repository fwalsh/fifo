package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func initDB() *sql.DB {
	connStr := fmt.Sprintf(
		"host=localhost port=5432 user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_USER", "items_user"),
		getEnv("DB_PASS", "items_pass"),
		getEnv("DB_NAME", "items_db"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("failed to connect to db:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping db:", err)
	}

	log.Println("connected to database")
	return db
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
