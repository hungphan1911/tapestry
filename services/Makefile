.PHONY: db migrate run

db:
	docker compose up -d postgres

migrate:
	go run ./cmd/migrate

run:
	go run ./cmd
