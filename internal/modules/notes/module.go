package notes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	//"strconv"

	"lifeos/internal/core"
)

type Note struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Pinned    bool   `json:"pinned"`
	UpdatedAt string `json:"updated_at"`
}

type NotesModule struct {
	db *sql.DB
}

func New() *NotesModule {
	return &NotesModule{}
}

func (m *NotesModule) ID() string   { return "core.notes" }
func (m *NotesModule) Name() string { return "Scratchpad" }

func (m *NotesModule) Init(ctx *core.ModuleCtx) error {
	m.db = ctx.DB

	// 1. Create table
	query := `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL DEFAULT 'Untitled',
		content TEXT NOT NULL DEFAULT '',
		pinned BOOLEAN NOT NULL DEFAULT 0,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := m.db.Exec(query); err != nil {
		return err
	}

	// 2. Register HTTP routes
	ctx.Mux.HandleFunc("GET /api/v1/notes", m.handleGetNotes)
	ctx.Mux.HandleFunc("POST /api/v1/notes", m.handleCreateNote)
	ctx.Mux.HandleFunc("PUT /api/v1/notes", m.handleUpdateNote)
	ctx.Mux.HandleFunc("DELETE /api/v1/notes", m.handleDeleteNote)

	return nil
}

// --- Endpoints ---

func (m *NotesModule) handleGetNotes(w http.ResponseWriter, r *http.Request) {
	rows, err := m.db.Query("SELECT id, title, content, pinned, updated_at FROM notes ORDER BY pinned DESC, updated_at DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	notes := []Note{}
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.Pinned, &n.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		notes = append(notes, n)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func (m *NotesModule) handleCreateNote(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	// Decode body if present, default to empty/untitled
	_ = json.NewDecoder(r.Body).Decode(&input)
	if input.Title == "" {
		input.Title = "New Note"
	}

	res, err := m.db.Exec("INSERT INTO notes (title, content) VALUES (?, ?)", input.Title, input.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Note{
		ID:      int(id),
		Title:   input.Title,
		Content: input.Content,
		Pinned:  false,
	})
}

func (m *NotesModule) handleUpdateNote(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID      int     `json:"id"`
		Title   *string `json:"title"`
		Content *string `json:"content"`
		Pinned  *bool   `json:"pinned"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.ID <= 0 {
		http.Error(w, "Invalid or missing note ID", http.StatusBadRequest)
		return
	}

	// Fetch existing note
	var current Note
	err := m.db.QueryRow("SELECT id, title, content, pinned FROM notes WHERE id = ?", input.ID).
		Scan(&current.ID, &current.Title, &current.Content, &current.Pinned)
	if err == sql.ErrNoRows {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Apply partial updates
	if input.Title != nil {
		current.Title = *input.Title
	}
	if input.Content != nil {
		current.Content = *input.Content
	}
	if input.Pinned != nil {
		current.Pinned = *input.Pinned
	}

	_, err = m.db.Exec(
		"UPDATE notes SET title = ?, content = ?, pinned = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		current.Title, current.Content, current.Pinned, current.ID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(current)
}

func (m *NotesModule) handleDeleteNote(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID int `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.ID <= 0 {
		http.Error(w, "Invalid or missing note ID", http.StatusBadRequest)
		return
	}

	res, err := m.db.Exec("DELETE FROM notes WHERE id = ?", input.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}