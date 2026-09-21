package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Levi67/LifeOS-Core/internal/db"
)

func main() {
	// Initialize SQLite database file
	database, err := db.InitDB("./lifeos.db")
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer database.Close()

	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","app":"LifeOS-Core","db":"connected"}`)
	})

	fmt.Println("🚀 CoreHub server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}
