# Services architecture

A modular monolith: one Go binary, several feature modules, each with its own
database on a shared Postgres server and a shared role.

## Layout

```
cmd/main.go            HTTP server; mounts every module under /api/v1/<name>
cmd/migrate/main.go    creates each module's database and applies its migrations
internal/core/         shared code: config, DB helper, errors, Module interface
internal/modules/      module registry (modules.go) + one folder per module
internal/modules/<m>/  models, repository, service, handler, module.go, db/migrations
```

## Adding a module

1. Create `internal/modules/<name>/` with `models.go`, `repository.go`,
   `service.go`, `handler.go`, and `db/migrations/*.sql`.
2. Add `module.go` (copy `finance/module.go`): an embedded `Migration()` and a
   `New(core.PostgresConfig) (core.Module, error)` that opens its own database
   with `core.OpenDB(pg.URL(core.DatabaseName("<name>")))`.
3. Register it in `internal/modules/modules.go` (both `Migrations()` and `Build()`).
4. Run `make migrate`. The database is created automatically (default name is
   the module name; override with `<NAME>_DB_NAME`).

## Boundary rules

- **No cross-module data access.** No foreign keys, joins, or queries into
  another module's database. Reference other modules' data by ID only.
- **Talk in-process through interfaces.** If module A needs module B, B exports
  a small interface (e.g. `finance.Reader`) and A receives it in its
  constructor. A never imports B's repository or opens B's database.
- **Keep internals unexported.** Only `New`, `Migration`, and any deliberately
  shared interfaces/types are exported from a module package.
- **Shared code goes in `internal/core`** (errors, config, HTTP helpers). Core
  must not import any module.
- **Each module owns its migrations.** One goose version table per database, so
  versions never collide between modules.
- **Cap connections per module** (`maxOpenConnsPerModule` in `core/database.go`)
  so the total stays below Postgres's `max_connections`.

## Configuration (`services/.env`)

| Variable | Required | Default | Purpose |
|---|---|---|---|
| `POSTGRES_USER` / `POSTGRES_PASSWORD` | yes | | shared role |
| `POSTGRES_DB` | no | `postgres` | admin DB used to create module DBs (also read by docker-compose) |
| `POSTGRES_HOST` / `POSTGRES_PORT` | no | `localhost` / `5432` | |
| `POSTGRES_SSLMODE` | no | `disable` | |
| `<MODULE>_DB_NAME` | no | module name | e.g. `FINANCE_DB_NAME` |
| `HTTP_PORT` | no | `8080` | |

## Deploying

Run `go run ./cmd/migrate` (or the built binary) once per deploy, before
starting the server. Keep migrations backward-compatible with the previous
app version.
