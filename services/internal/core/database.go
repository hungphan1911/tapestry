package core

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// maxOpenConnsPerModule caps each module's pool so that the sum across all
// modules stays well below Postgres's connection limit.
const maxOpenConnsPerModule = 10

// OpenDB opens a pooled connection to a database and verifies it is reachable.
func OpenDB(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	db.SetMaxOpenConns(maxOpenConnsPerModule)
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return db, nil
}
