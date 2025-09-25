package main

import (
	"log"
	"net/http"

	"github.com/fwalsh/fifo/fifo"
)

func main() {
	// Initialize database (retries internally until ready)
	db := fifo.InitDB()

	// Register routes
	http.HandleFunc("/", fifo.RootHandler(db))
	http.HandleFunc("/health", fifo.HealthHandler)
	http.HandleFunc("/items", fifo.ItemsHandler(db))
	http.HandleFunc("/items/create", fifo.CreateItemHandler(db))

	log.Println("🚀 fifo app running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
