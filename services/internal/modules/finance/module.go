package finance

import (
	"database/sql"
	"embed"

	"github.com/go-chi/chi/v5"

	"github.com/hungphan1911/tapestry/services/internal/core"
)

const moduleName = "finance"

//go:embed db/migrations/*.sql
var migrationsFS embed.FS

// Migration returns the module's embedded migrations.
func Migration() core.Migration {
	return core.Migration{Module: moduleName, FS: migrationsFS, Dir: "db/migrations"}
}

type module struct {
	db      *sql.DB
	handler *Handler
}

// New opens the finance database and wires repository, service and handler.
func New(pg core.PostgresConfig) (core.Module, error) {
	db, err := core.OpenDB(pg.URL(core.DatabaseName(moduleName)))
	if err != nil {
		return nil, err
	}
	svc := NewService(NewRepository(db))
	return &module{db: db, handler: NewHandler(svc)}, nil
}

func (m *module) Name() string                { return moduleName }
func (m *module) RegisterRoutes(r chi.Router) { m.handler.RegisterRoutes(r) }
func (m *module) Close() error                { return m.db.Close() }

var _ core.Module = (*module)(nil)
