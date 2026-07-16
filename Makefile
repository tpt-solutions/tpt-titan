.PHONY: help setup dev build backend frontend desktop test lint ci check install tidy seed backup

# Default target
help:
	@echo "TPT Titan — task runner"
	@echo ""
	@echo "Setup / dev:"
	@echo "  make setup        Install backend + frontend + desktop dependencies"
	@echo "  make dev          Run backend and frontend dev servers concurrently"
	@echo "  make build        Build backend binary and frontend (into backend/static or dist)"
	@echo ""
	@echo "Backend tooling:"
	@echo "  make seed         Insert demo/seed data (default SQLite DB)"
	@echo "  make backup NAME=manual   Create a full database backup"
	@echo "  make admin ARGS=...       Run management CLI, e.g. make admin ARGS='create-user --username admin --password ...'"
	@echo ""
	@echo "Quality:"
	@echo "  make test         Run all Go tests"
	@echo "  make lint         gofmt + go vet on the backend"
	@echo "  make ci           Run the same checks CI runs (test + lint)"
	@echo ""
	@echo "Misc:"
	@echo "  make tidy         go mod tidy"
	@echo "  make install      Build release binary (requires CGO/gcc for SQLite)"

# Load environment from .env if present so DB_* / JWT_* are picked up.
include .env
export

setup:
	cd backend && go mod download
	cd frontend && npm install
	cd desktop && npm install

dev:
	# Runs backend (go run) and frontend (vite) side by side.
	@echo "Starting backend on :8080 and frontend on :5173 (press Ctrl-C to stop)"
	cd backend && go run . & cd frontend && npm run dev & wait

build:
	cd backend && go build -o ../bin/tpt-titan .
	cd frontend && npm run build

tidy:
	cd backend && go mod tidy

seed: build
	./bin/tpt-titan seed

backup: build
	./bin/tpt-titan backup $(NAME)

admin: build
	./bin/tpt-titan admin $(ARGS)

test:
	cd backend && go test ./...

lint:
	cd backend && gofmt -l . && go vet ./...

ci: test lint

install:
	cd backend && CGO_ENABLED=1 go build -o ../bin/tpt-titan .
