BINARY     := bin/web
PKG        := ./cmd/web
DB         ?= app.db
MIGRATIONS := migrations
PORT       := 8080

GOOSE := go tool goose

.PHONY: help generate build run dev test vet fmt check clean \
        migrate-new migrate-status reset-db

help: ## List targets
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | awk -F':.*## ' '{printf "%-16s %s\n", $$1, $$2}'

# templ and sqlc are pinned by the `tool` block in go.mod, so no install step.
generate: ## Run templ and sqlc codegen
	go tool templ generate
	go tool sqlc generate

build: generate ## Build the binary into bin/
	go build -o $(BINARY) $(PKG)

run: generate ## Generate, then run the server
	go run $(PKG) -db $(DB)

dev: ## Regenerate templates on change, restart the app, reload the browser
	go tool templ generate --watch \
		--proxy="http://localhost:$(PORT)" \
		--open-browser=false \
		--cmd="go run $(PKG) -db $(DB)"

test: generate ## Run tests
	go test ./...

vet: ## Static checks
	go vet ./...

fmt: ## Format Go and templ files
	go tool templ fmt .
	go fmt ./...

check: fmt vet test ## Everything before a commit

clean: ## Remove build output
	rm -rf bin

migrate-new: ## Create a migration: make migrate-new name=add_votes_index
	@test -n "$(name)" || { echo "usage: make migrate-new name=<snake_case_name>"; exit 1; }
	$(GOOSE) -dir $(MIGRATIONS) create $(name) sql

migrate-status: ## Show applied and pending migrations for $(DB)
	$(GOOSE) -dir $(MIGRATIONS) sqlite3 $(DB) status

reset-db: ## Delete the local SQLite DB; the next run recreates it
	rm -f $(DB) $(DB)-wal $(DB)-shm
