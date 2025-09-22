package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/lib/pq"
	fifo "github.com/fwalsh/fifo/fifo" // updated import path
)

func TestCreateAndListItems(t *testing.T) {
	// Connect to DB (use env var DATABASE_URL if set, else fallback)
	connStr := "host=localhost port=5432 user=items_user password=items_pass dbname=items_db sslmode=disable"
	if v := os.Getenv("DATABASE_URL"); v != "" {
		connStr = v
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("cannot connect to db: %v", err)
	}

	// Set up router with fifo handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			fifo.CreateItemHandler(db)(w, r)
		case http.MethodGet:
			fifo.GetItemsHandler(db)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// --- Test POST /items ---
	payload := []byte(`{"name":"pear"}`)
	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("POST /items expected 200, got %d", w.Code)
	}

	var created map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if created["name"] != "pear" {
		t.Errorf("expected name pear, got %s", created["name"])
	}

	// --- Test GET /items ---
	req = httptest.NewRequest(http.MethodGet, "/items", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /items expected 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("pear")) {
		t.Errorf("expected pear in response, got %s", w.Body.String())
	}
}
