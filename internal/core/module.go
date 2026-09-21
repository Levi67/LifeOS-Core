package core

import (
	"database/sql"
	"net/http"
)

// ModuleCtx gives modules access to core engine services
type ModuleCtx struct {
	DB  *sql.DB
	Mux *http.ServeMux
}

// Module is the interface every LifeOS module must implement
type Module interface {
	ID() string          // Unique identifier, e.g., "core.tasks"
	Name() string        // Display name
	Init(ctx *ModuleCtx) error // Setup DB tables, register HTTP routes
}

// Registry manages active modules
type Registry struct {
	modules map[string]Module
}

func NewRegistry() *Registry {
	return &Registry{
		modules: make(map[string]Module),
	}
}

func (r *Registry) Register(m Module) {
	r.modules[m.ID()] = m
}

func (r *Registry) InitAll(ctx *ModuleCtx) error {
	for id, m := range r.modules {
		if err := m.Init(ctx); err != nil {
			return err
		}
		println("[Module Loaded]:", id, "-", m.Name())
	}
	return nil
}