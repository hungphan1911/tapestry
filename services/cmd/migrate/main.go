// Command migrate creates each module's database (if missing) and applies its
// migrations. Run it as a deploy step before starting the server.
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"

	"github.com/hungphan1911/tapestry/services/internal/core"
	"github.com/hungphan1911/tapestry/services/internal/modules"
)

func main() {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Loading .env failed: %v", err)
	}
	cfg, err := core.LoadConfig()
	if err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	for _, m := range modules.Migrations() {
		if err := migrate(cfg.Postgres, m); err != nil {
			log.Fatalf("Migrating %s failed: %v", m.Module, err)
		}
	}
}

func migrate(pg core.PostgresConfig, m core.Migration) error {
	name := core.DatabaseName(m.Module)
	if err := ensureDatabase(pg, name); err != nil {
		return err
	}

	db, err := core.OpenDB(pg.URL(name))
	if err != nil {
		return err
	}
	defer db.Close()

	// Each module has its own database, so the default goose version table
	// does not collide between modules.
	goose.SetBaseFS(m.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	fmt.Printf("== %s (%s)\n", m.Module, name)
	return goose.Up(db, m.Dir)
}

func ensureDatabase(pg core.PostgresConfig, name string) error {
	admin, err := core.OpenDB(pg.URL(pg.AdminDB))
	if err != nil {
		return err
	}
	defer admin.Close()

	return createDatabaseIfMissing(admin, name)
}

func createDatabaseIfMissing(admin *sql.DB, name string) error {
	var exists bool
	err := admin.QueryRow(`select exists (select 1 from pg_database where datname = $1)`, name).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check database %q: %w", name, err)
	}
	if exists {
		return nil
	}
	// identifiers cannot be parameterised; quote it instead
	if _, err := admin.Exec(`create database ` + quoteIdent(name)); err != nil {
		return fmt.Errorf("create database %q: %w", name, err)
	}
	return nil
}

func quoteIdent(s string) string {
	out := `"`
	for _, r := range s {
		if r == '"' {
			out += `"`
		}
		out += string(r)
	}
	return out + `"`
}
