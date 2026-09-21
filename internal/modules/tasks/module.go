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
	ctx.Mux.HandleFunc("DELETE /api/v1/tasks", m.handleDeleteTask)
	ctx.Mux.HandleFunc("PUT /api/v1/tasks", m.handleUpdateTask)

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

func (m *TasksModule) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
    var input struct {
        ID int `json:"id"`
    }

    if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.ID <= 0 {
        http.Error(w, "Invalid or missing task ID", http.StatusBadRequest)
        return
    }

    res, err := m.db.Exec("DELETE FROM tasks WHERE id = ?", input.ID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Check how many rows were actually deleted
    rowsAffected, err := res.RowsAffected()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // If 0 rows were affected, the ID didn't exist in the database
    if rowsAffected == 0 {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (m *TasksModule) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID        int     `json:"id"`
		Title     *string `json:"title"`
		Completed *bool   `json:"completed"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.ID <= 0 {
		http.Error(w, "Invalid or missing task ID", http.StatusBadRequest)
		return
	}

	// 1. Fetch current task to allow partial updates
	var currentTitle string
	var currentCompleted bool
	err := m.db.QueryRow("SELECT title, completed FROM tasks WHERE id = ?", input.ID).Scan(&currentTitle, &currentCompleted)
	if err == sql.ErrNoRows {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. Use provided values or keep existing ones
	newTitle := currentTitle
	if input.Title != nil {
		newTitle = *input.Title
	}

	newCompleted := currentCompleted
	if input.Completed != nil {
		newCompleted = *input.Completed
	}

	// 3. Update SQLite record
	res, err := m.db.Exec("UPDATE tasks SET title = ?, completed = ? WHERE id = ?", newTitle, newCompleted, input.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// 4. Return updated task object
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Task{
		ID:        input.ID,
		Title:     newTitle,
		Completed: newCompleted,
	})
}



func (m *TasksModule) parseBool(b bool) *bool { return &b }