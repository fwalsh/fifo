package fifo

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"
)

// Item represents a row in the items table
type Item struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// RootHandler serves the friendly homepage with a form + items list
func RootHandler(db *sql.DB) http.HandlerFunc {
	tmpl := template.Must(template.New("index").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>🍀 Welcome to fifo</title>
			<style>
				body { font-family: sans-serif; text-align: center; margin-top: 50px; }
				h1 { font-size: 2.5em; color: #2E8B57; }
				form { margin: 20px auto; }
				input { padding: 10px; margin: 5px; }
				table { margin: 20px auto; border-collapse: collapse; }
				td, th { padding: 10px; border: 1px solid #ddd; }
			</style>
		</head>
		<body>
			<h1>🌌 fifo API</h1>
			<p>Add an item:</p>
			<form action="/items/create" method="post">
				<input type="text" name="name" placeholder="Item name" required />
				<input type="submit" value="Add" />
			</form>
			<h2>Items</h2>
			<table>
				<tr><th>ID</th><th>Name</th><th>Created At</th></tr>
				{{range .}}
					<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.CreatedAt}}</td></tr>
				{{else}}
					<tr><td colspan="3">No items yet</td></tr>
				{{end}}
			</table>
		</body>
		</html>
	`))

	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, created_at FROM items ORDER BY created_at DESC")
		if err != nil {
			http.Error(w, "failed to query items", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var items []Item
		for rows.Next() {
			var it Item
			if err := rows.Scan(&it.ID, &it.Name, &it.CreatedAt); err == nil {
				items = append(items, it)
			}
		}

		if err := tmpl.Execute(w, items); err != nil {
			log.Printf("template execute error: %v", err)
		}
	}
}

// HealthHandler returns a simple health status
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// ItemsHandler lists items in JSON
func ItemsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, created_at FROM items ORDER BY created_at DESC")
		if err != nil {
			http.Error(w, "failed to query items", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var items []Item
		for rows.Next() {
			var it Item
			if err := rows.Scan(&it.ID, &it.Name, &it.CreatedAt); err == nil {
				items = append(items, it)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

// CreateItemHandler inserts a new item into the database
func CreateItemHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		if name == "" {
			http.Error(w, "missing name", http.StatusBadRequest)
			return
		}

		var id int
		var created time.Time
		err := db.QueryRow(
			"INSERT INTO items (name, created_at) VALUES ($1, NOW()) RETURNING id, created_at",
			name,
		).Scan(&id, &created)

		if err != nil {
			http.Error(w, "failed to insert item", http.StatusInternalServerError)
			return
		}

		item := Item{ID: id, Name: name, CreatedAt: created}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(item)
	}
}
