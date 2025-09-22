package main

import (
	"fmt"
	"net/http"

	"github.com/fwalsh/fifo/fifo" // import the library code
)

func main() {
	db := fifo.InitDB()
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message":"welcome to fifo API"}`)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok"}`)
	})

	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fifo.GetItemsHandler(db)(w, r)
		case http.MethodPost:
			fifo.CreateItemHandler(db)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("server running on :8080")
	http.ListenAndServe(":8080", nil)
}
