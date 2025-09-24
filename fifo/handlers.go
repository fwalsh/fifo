package fifo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Item struct maps to the items table
type Item struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// HealthHandler – simple health check
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// CreateItemHandler – adds an item (JSON API)
func CreateItemHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var item Item
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		err := db.QueryRow(
			"INSERT INTO items (name, created_at) VALUES ($1, NOW()) RETURNING id, created_at",
			item.Name,
		).Scan(&item.ID, &item.CreatedAt)

		if err != nil {
			http.Error(w, "Failed to insert item", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(item)
	}
}

// GetItemsHandler – returns items as JSON
func GetItemsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, created_at FROM items ORDER BY id DESC")
		if err != nil {
			http.Error(w, "Failed to query items", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var items []Item
		for rows.Next() {
			var item Item
			if err := rows.Scan(&item.ID, &item.Name, &item.CreatedAt); err == nil {
				items = append(items, item)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

// RootHandler – friendly homepage (with items + form)
func RootHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handle form submission
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err == nil {
				name := r.FormValue("name")
				if name != "" {
					_, _ = db.Exec("INSERT INTO items (name, created_at) VALUES ($1, NOW())", name)
				}
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Otherwise, render the page
		rows, err := db.Query("SELECT id, name, created_at FROM items ORDER BY id DESC")
		if err != nil {
			http.Error(w, "Failed to query items", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var tableRows string
		for rows.Next() {
			var id int
			var name string
			var createdAt time.Time
			if err := rows.Scan(&id, &name, &createdAt); err == nil {
				tableRows += fmt.Sprintf(
					"<tr><td>%d</td><td>%s</td><td>%s</td></tr>",
					id, name, createdAt.Format("2006-01-02 15:04:05"),
				)
			}
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
		  <!DOCTYPE html>
		  <html>
		  <head>
		    <title>🌀 fifo</title>
		    <style>
		      body { font-family: sans-serif; text-align: center; padding: 2rem; background: #fafafa; }
		      h1 { color: #4B0082; }
		      form { margin: 1rem 0; }
		      input[type=text] { padding: 0.5rem; font-size: 1rem; }
		      button { padding: 0.5rem 1rem; font-size: 1rem; background: #4B0082; color: white; border: none; border-radius: 4px; cursor: pointer; }
		      button:hover { background: #360060; }
		      table { margin: 1rem auto; border-collapse: collapse; width: 80%%; }
		      th, td { border: 1px solid #ddd; padding: 8px; }
		      th { background-color: #eee; }
		    </style>
		  </head>
		  <body>
		    <h1>🌀 Welcome to fifo!</h1>
		    <p>A friendly little Items API, powered by Go + Postgres + Docker + CircleCI.</p>

		    <h2>📦 Add an Item</h2>
		    <form method="POST" action="/">
		      <input type="text" name="name" placeholder="Enter item name" required />
		      <button type="submit">Add</button>
		    </form>

		    <h2>📋 Current Items</h2>
		    <table>
		      <tr><th>ID</th><th>Name</th><th>Created At</th></tr>
		      %s
		    </table>

		    <p>
		      Or check: <a href="/health">Health</a> | <a href="/items">Raw JSON Items</a>
		    </p>
		  </body>
		  </html>
		`, tableRows)
	}
}
