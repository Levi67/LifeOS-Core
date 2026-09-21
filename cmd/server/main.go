package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "modernc.org/sqlite"

	"lifeos/internal/core"
	"lifeos/internal/modules/tasks"
	"lifeos/internal/modules/notes"
)

func main() {
	// 1. Core Storage Setup
	db, err := sql.Open("sqlite", "lifeos.db")
	if err != nil {
		log.Fatalf("Failed to open SQLite: %v", err)
	}
	defer db.Close()

	// 2. Router Setup
	mux := http.NewServeMux()

	// 3. Module Registry Initialization
	ctx := &core.ModuleCtx{
		DB:  db,
		Mux: mux,
	}

	registry := core.NewRegistry()

	// --- Register Modules Here ---
	registry.Register(tasks.New())
	registry.Register(notes.New())
	// registry.Register(notes.New())  <-- Future modules are added like this!

	if err := registry.InitAll(ctx); err != nil {
		log.Fatalf("Module initialization error: %v", err)
	}

	// 4. Start Core Engine
	log.Println("LifeOS Core Engine listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}