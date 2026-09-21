package tasks

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"lifeos/internal/core"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type TasksModule struct {
	db *sql.DB
}

func New() *TasksModule {
	return &TasksModule{}
}

func (m *TasksModule) ID() string   { return "core.tasks" }
func (m *TasksModule) Name() string { return "To-Do List" }

func (m *TasksModule) Init(ctx *core.ModuleCtx) error {
	m.db = ctx.DB

	// 1. Create table
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := m.db.Exec(query); err != nil {
		return err
	}

	// 2. Register API endpoints directly on core router
	ctx.Mux.HandleFunc("GET /api/v1/tasks", m.handleGetTasks)
	ctx.Mux.HandleFunc("POST /api/v1/tasks", m.handleCreateTask)

	return nil
}

// --- Endpoints ---

func (m *TasksModule) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := m.db.Query("SELECT id, title, completed FROM tasks ORDER BY id DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		// Scan directly into the struct fields
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (m *TasksModule) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Title == "" {
		http.Error(w, "Invalid task payload", http.StatusBadRequest)
		return
	}

	res, err := m.db.Exec("INSERT INTO tasks (title, completed) VALUES (?, 0)", input.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Task{ID: int(id), Title: input.Title, Completed: false})
}

func (m *TasksModule) parseBool(b bool) *bool { return &b }