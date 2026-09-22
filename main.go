package main

import (
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"net/http"

	_ "modernc.org/sqlite"

	"lifeos/internal/core"
	"lifeos/internal/modules/notes"
	"lifeos/internal/modules/tasks"
	"lifeos/internal/modules/habittracker"
)

// Embed the compiled frontend files from web/dist into the binary.
// Make sure you run `npm run build` in your `web/` directory first!
//go:embed web/dist/*
var webDist embed.FS

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
	registry.Register(habittracker.New())

	if err := registry.InitAll(ctx); err != nil {
		log.Fatalf("Module initialization error: %v", err)
	}

	// 4. Serve Embedded Frontend Static Assets
	distFS, err := fs.Sub(webDist, "web/dist")
	if err != nil {
		log.Fatalf("Failed to create sub filesystem for web dist: %v", err)
	}
	
	// http.ServeMux routes /api/v1/... to modules first.
	// Any path not matching an API endpoint falls back to serving static files from web/dist.
	mux.Handle("/", http.FileServer(http.FS(distFS)))

	// 5. Start Core Engine
	log.Println("LifeOS Core Engine listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}