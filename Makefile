.PHONY: help db db-wait db-stop migrate backend frontend dev

help:
	@echo "make db         start Postgres (docker)"
	@echo "make migrate    start Postgres, then apply module migrations"
	@echo "make backend    run the Go API on :8080 (needs the database)"
	@echo "make frontend   run the Vite dev server on :5173"
	@echo "make dev        whole stack: database, migrations, backend and frontend (Ctrl-C stops both servers)"
	@echo "make db-stop    stop Postgres (data is kept)"

db:
	$(MAKE) -C services db

db-wait: db
	@until docker compose -f services/docker-compose.yml exec -T postgres pg_isready -q; do sleep 1; done

db-stop:
	docker compose -f services/docker-compose.yml stop postgres

migrate: db-wait
	$(MAKE) -C services migrate

backend:
	$(MAKE) -C services run

frontend:
	cd frontend && npm run dev

dev: migrate
	$(MAKE) -j2 backend frontend
