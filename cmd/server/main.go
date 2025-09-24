package main

import (
	"log"
	"net/http"

	"github.com/fwalsh/fifo/fifo"
)

func main() {
	// Initialize DB connection
	db, err := fifo.InitDB()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	// Routes
	http.HandleFunc("/", fifo.RootHandler(db))        // friendly homepage
	http.HandleFunc("/health", fifo.HealthHandler)    // health check

	// Unified handler for /items
	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fifo.GetItemsHandler(db)(w, r)
		case http.MethodPost:
			fifo.CreateItemHandler(db)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start server
	log.Println("🚀 fifo app running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
