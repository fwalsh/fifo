package main

import (
  "fmt"
  "net/http"
)

func main() {
  http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"status":"ok"}`)
  })

  fmt.Println("server running on :8080")
  http.ListenAndServe(":8080", nil)
}
