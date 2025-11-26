.PHONY: help install migrate run build clean test

help: # show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: # install Go dependencies
	@echo "Installing dependencies"
	go mod download
	go mod tidy
	@echo "Dependencies installed!"

migrate: # run database migrations
	@echo "Running database migrations"
	@if [ -z "$(DB_USER)" ]; then \
		echo "Using default DB_USER=postgres"; \
		psql -U postgres -d finance_manager -f migrations/001_create_tables.sql; \
	else \
		psql -U $(DB_USER) -d $(DB_NAME) -f migrations/001_create_tables.sql; \
	fi
	@echo "Migrations completed!"

createdb: # create the database
	@echo "Creating database..."
	@if [ -z "$(DB_USER)" ]; then \
		psql -U postgres -c "CREATE DATABASE finance_manager;"; \
	else \
		psql -U $(DB_USER) -c "CREATE DATABASE $(DB_NAME);"; \
	fi
	@echo "Database created!"

setup: createdb migrate install ## Complete setup (create DB, run migrations, install dependencies)
	@echo "Setup completed! Run 'make run' to start the server."

run: # run the application
	@echo "Starting server"
	go run cmd/api/main.go

build: ## Build the application
	@echo "Building application"
	go build -o bin/finance-manager cmd/api/main.go
	@echo "Binary created at: bin/finance-manager"

dev: # run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
	air

test: # run tests
	go test -v ./...

clean: # clean build artifacts
	@echo "Cleaning"
	rm -rf bin/
	go clean
	@echo "Cleaned!"

docker-up: ## Start PostgreSQL with Docker
	docker run --name finance-postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:14

docker-down: ## Stop PostgreSQL Docker container
	docker stop finance-postgres
	docker rm finance-postgres

.DEFAULT_GOAL := help