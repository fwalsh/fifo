package tests

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fwalsh/fifo/fifo"
	_ "github.com/lib/pq"
)

func setupDB(t *testing.T) *sql.DB {
	db := fifo.InitDB()

	// Clean up items before each test
	_, err := db.Exec("DELETE FROM items")
	if err != nil {
		t.Fatalf("failed to clean items: %v", err)
	}
	return db
}

func TestCreateAndListItems(t *testing.T) {
	db := setupDB(t)

	mux := http.NewServeMux()
	mux.Handle("/items", fifo.ItemsHandler(db))
	mux.Handle("/items/create", fifo.CreateItemHandler(db))

	// Create an item
	req := httptest.NewRequest("POST", "/items/create", strings.NewReader("name=test-item"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// List items
	req = httptest.NewRequest("GET", "/items", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "test-item") {
		t.Fatalf("expected response to contain item name, got %s", body)
	}
}
