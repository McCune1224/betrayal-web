# =============================================================================
# Betrayal - Development Commands
# =============================================================================
#
# Usage: make <target>
#
# Run 'make help' to see all available commands.
#
# =============================================================================

.PHONY: help dev dev-backend dev-frontend test test-backend build build-backend build-frontend clean install db-up db-down db-migrate

# Default target
help:
	@echo "Betrayal Game - Available Commands"
	@echo ""
	@echo "Development:"
	@echo "  make dev            - Run backend and frontend (requires 2 terminals)"
	@echo "  make dev-backend    - Run backend server"
	@echo "  make dev-frontend   - Run frontend dev server"
	@echo ""
	@echo "Testing:"
	@echo "  make test           - Run all tests"
	@echo "  make test-backend   - Run backend tests only"
	@echo "  make test-verbose   - Run backend tests with verbose output"
	@echo ""
	@echo "Building:"
	@echo "  make build          - Build both backend and frontend"
	@echo "  make build-backend  - Build backend binary"
	@echo "  make build-frontend - Build frontend for production"
	@echo ""
	@echo "Database:"
	@echo "  make db-up          - Start local Postgres (Docker)"
	@echo "  make db-down        - Stop local Postgres"
	@echo "  make db-migrate     - Run database migrations"
	@echo "  make db-reset       - Reset database (down + up + migrate)"
	@echo ""
	@echo "Setup:"
	@echo "  make install        - Install all dependencies"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make sqlc           - Generate sqlc code"

# =============================================================================
# Development
# =============================================================================

dev:
	@echo "Run these commands in separate terminals:"
	@echo "  Terminal 1: make dev-backend"
	@echo "  Terminal 2: make dev-frontend"

dev-backend:
	cd backend && go run cmd/server/main.go

dev-frontend:
	cd frontend && bun run dev

# =============================================================================
# Testing
# =============================================================================

test: test-backend

test-backend:
	cd backend && go test ./...

test-verbose:
	cd backend && go test ./... -v

test-coverage:
	cd backend && go test ./... -cover

# =============================================================================
# Building
# =============================================================================

build: build-backend build-frontend

build-backend:
	cd backend && go build -o bin/server cmd/server/main.go

build-frontend:
	cd frontend && bun run build

# =============================================================================
# Database
# =============================================================================

db-up:
	docker run --name social-deduction-postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=social_deduction \
		-p 5432:5432 \
		-d postgres:15
	@echo "Postgres started. Connection: postgres://postgres:postgres@localhost:5432/social_deduction"

db-down:
	docker stop social-deduction-postgres || true
	docker rm social-deduction-postgres || true

db-migrate:
	cd backend && migrate -path internal/db/migrations -database "postgres://postgres:postgres@localhost:5432/social_deduction?sslmode=disable" up

db-migrate-down:
	cd backend && migrate -path internal/db/migrations -database "postgres://postgres:postgres@localhost:5432/social_deduction?sslmode=disable" down

db-reset: db-down db-up
	@echo "Waiting for Postgres to start..."
	@sleep 3
	$(MAKE) db-migrate

# =============================================================================
# Setup & Utilities
# =============================================================================

install:
	cd backend && go mod download
	cd frontend && bun install

clean:
	rm -rf backend/bin
	rm -rf frontend/.svelte-kit
	rm -rf frontend/build

sqlc:
	cd backend && sqlc generate

check:
	cd backend && go build ./...
	cd frontend && bun run check

# =============================================================================
# Quick Start (first time setup)
# =============================================================================

setup: install db-up
	@sleep 3
	$(MAKE) db-migrate
	@echo ""
	@echo "Setup complete! Run 'make dev-backend' and 'make dev-frontend' in separate terminals."
