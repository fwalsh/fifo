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
    fifo "github.com/fwalsh/fifo"   // <-- import our app
)


func TestCreateAndListItems(t *testing.T) {
	// Connect to DB
	connStr := "host=localhost port=5432 user=items_user password=items_pass dbname=items_db sslmode=disable"
	if v := os.Getenv("DATABASE_URL"); v != "" {
		connStr = v
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Start server in-memory
	handler := http.NewServeMux()
	handler.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createItemHandler(db)(w, r)
		case http.MethodGet:
			getItemsHandler(db)(w, r)
		}
	})

	// POST /items
	payload := []byte(`{"name":"pear"}`)
	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// Decode response
	var item map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &item); err != nil {
		t.Fatal(err)
	}
	if item["name"] != "pear" {
		t.Errorf("expected name pear, got %s", item["name"])
	}

	// GET /items
	req = httptest.NewRequest(http.MethodGet, "/items", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("pear")) {
		t.Errorf("expected pear in response, got %s", w.Body.String())
	}
}
