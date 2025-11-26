.PHONY: help install migrate run build clean test dev
help: 
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?
install: 
	@echo "Installing dependencies"
	go mod download
	go mod tidy
	@echo "Dependencies installed!"
migrate: 
	@echo "Running database migrations"
	@if [ -z "$(DB_USER)" ]; then \
		echo "Using default DB_USER=postgres"; \
		psql -U postgres -d finance_manager -f migrations/001_create_tables.sql; \
	else \
		psql -U $(DB_USER) -d $(DB_NAME) -f migrations/001_create_tables.sql; \
	fi
	@echo "Migrations completed!"
createdb: 
	@echo "Creating database"
	@if [ -z "$(DB_USER)" ]; then \
		psql -U postgres -c "CREATE DATABASE finance_manager;"; \
	else \
		psql -U $(DB_USER) -c "CREATE DATABASE $(DB_NAME);"; \
	fi
	@echo "Database created!"
setup: createdb migrate install 
	@echo "Setup completed! Run 'make run' to start the server."
run: 
	@echo "Starting server"
	go run cmd/api/main.go
build: 
	@echo "Building application"
	go build -o bin/finance-manager cmd/api/main.go
	@echo "Binary created at: bin/finance-manager"
dev: 
	air
test: 
	go test -v ./...
clean: 
	@echo "Cleaning"
	rm -rf bin/
	go clean
	@echo "Cleaned!"
docker-up: 
	docker run --name finance-postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:14
docker-down: 
	docker stop finance-postgres
	docker rm finance-postgres
.DEFAULT_GOAL := help