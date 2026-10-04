package core

import (
	"io/fs"

	"github.com/go-chi/chi/v5"
)

// Module is a self-contained feature area that owns its own database.
// main mounts its routes under /api/v1/<Name()>.
type Module interface {
	Name() string
	RegisterRoutes(r chi.Router)
	Close() error
}

// Migration describes where a module's goose migrations live. FS is rooted so
// that Dir resolves to the folder containing the .sql files.
type Migration struct {
	Module string
	FS     fs.FS
	Dir    string
}
