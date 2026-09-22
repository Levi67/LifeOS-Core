package habittracker

import (
    "database/sql"
    //"encoding/json"
    //"net/http"

    "lifeos/internal/core"
)

type Habit struct {
    ID      int    `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
    Pinned  bool   `json:"pinned"`
}

type HabitModule struct {
    db *sql.DB
}

func New() *HabitModule {
    return &HabitModule{}
}

func (m *HabitModule) ID() string   { return "core.habit-tracker" }
func (m *HabitModule) Name() string { return "Habit Tracker" }

func (m *HabitModule) Init(ctx *core.ModuleCtx) error {
    m.db = ctx.DB

    query := `
    CREATE TABLE IF NOT EXISTS habits (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL DEFAULT 'Untitled',
        content TEXT NOT NULL DEFAULT '',
        pinned BOOLEAN NOT NULL DEFAULT 0,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`
    if _, err := m.db.Exec(query); err != nil {
        return err
    }

    // Angepasste Handler-Namen für Habits
    /*
    ctx.Mux.HandleFunc("GET /api/v1/habits", m.handleGetHabits)
    ctx.Mux.HandleFunc("POST /api/v1/habits", m.handleCreateHabit)
    ctx.Mux.HandleFunc("PUT /api/v1/habits", m.handleUpdateHabit)
    ctx.Mux.HandleFunc("DELETE /api/v1/habits", m.handleDeleteHabit)
    */

    return nil
}