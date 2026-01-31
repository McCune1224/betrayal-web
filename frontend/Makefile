# Makefile for Frontend Development
# Usage: make <command>

.PHONY: help install dev build preview check check-watch clean typecheck lint format

# Default target
help:
	@echo "Frontend Development Commands:"
	@echo ""
	@echo "  make install      - Install dependencies with bun"
	@echo "  make dev          - Start development server"
	@echo "  make build        - Build for production"
	@echo "  make preview      - Preview production build"
	@echo ""
	@echo "  make check        - Run svelte-check (TypeScript + Svelte validation)"
	@echo "  make check-watch  - Run svelte-check in watch mode"
	@echo "  make typecheck    - Alias for 'make check'"
	@echo ""
	@echo "  make clean        - Clean build artifacts"
	@echo "  make format       - Format code (if prettier is installed)"
	@echo "  make lint         - Lint code (if eslint is installed)"
	@echo ""

# Install dependencies
install:
	bun install

# Start development server
dev:
	bun run dev

# Build for production
build:
	bun run build

# Preview production build
preview:
	bun run preview

# Run svelte-check (validates TypeScript and Svelte)
check:
	bun run check

# Run svelte-check in watch mode
check-watch:
	bun run check:watch

# Alias for check
typecheck: check

# Clean build artifacts
clean:
	rm -rf .svelte-kit
	rm -rf build
	rm -rf dist
	@echo "Cleaned build artifacts"

# Format code (requires prettier to be installed)
format:
	@if command -v prettier >/dev/null 2>&1; then \
		prettier --write "src/**/*.{ts,svelte}" ; \
	else \
		echo "Prettier not installed. Run: bun add -d prettier"; \
	fi

# Lint code (requires eslint to be installed)
lint:
	@if command -v eslint >/dev/null 2>&1; then \
		eslint "src/**/*.{ts,svelte}" ; \
	else \
		echo "ESLint not installed. Run: bun add -d eslint"; \
	fi

# Full validation (check types and build)
validate: check build
	@echo "✅ Validation complete - TypeScript checks and build passed"

# Quick check - just type checking without build
quick: check
	@echo "✅ Quick check complete"
