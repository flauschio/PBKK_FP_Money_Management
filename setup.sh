#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Go is installed
echo -n "Checking Go installation... "
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗${NC}"
    echo "Go is not installed. Please install Go 1.21 or higher from https://golang.org/dl/"
    exit 1
fi
echo -e "${GREEN}✓${NC}"

# Check if PostgreSQL is installed
echo -n "Checking PostgreSQL installation... "
if ! command -v psql &> /dev/null; then
    echo -e "${RED}✗${NC}"
    echo "PostgreSQL is not installed. Please install PostgreSQL from https://www.postgresql.org/download/"
    exit 1
fi
echo -e "${GREEN}✓${NC}"

# Get database credentials
echo ""
echo "Database Configuration"
echo "======================"
read -p "Database host (default: localhost): " DB_HOST
DB_HOST=${DB_HOST:-localhost}

read -p "Database port (default: 5432): " DB_PORT
DB_PORT=${DB_PORT:-5432}

read -p "Database user (default: postgres): " DB_USER
DB_USER=${DB_USER:-postgres}

read -sp "Database password: " DB_PASSWORD
echo ""

read -p "Database name (default: finance_manager): " DB_NAME
DB_NAME=${DB_NAME:-finance_manager}

read -p "Server port (default: 8080): " PORT
PORT=${PORT:-8080}

# Create .env file
echo ""
echo -n "Creating .env file... "
cat > .env << EOF
DB_HOST=$DB_HOST
DB_PORT=$DB_PORT
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
DB_NAME=$DB_NAME
PORT=$PORT
EOF
echo -e "${GREEN}✓${NC}"

# Export variables for this session
export DB_HOST DB_PORT DB_USER DB_PASSWORD DB_NAME PORT

# Create database
echo -n "Creating database... "
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "CREATE DATABASE $DB_NAME;" 2>/dev/null
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${YELLOW}⚠${NC} (Database might already exist)"
fi

# Run migrations
echo -n "Running database migrations... "
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f migrations/001_create_tables.sql > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    echo "Failed to run migrations. Please check your database connection."
    exit 1
fi

# Install Go dependencies
echo -n "Installing Go dependencies... "
go mod download > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    echo "Failed to install dependencies."
    exit 1
fi

# Success message
echo ""
echo -e "${GREEN}=================================="
echo "Setup completed successfully!"
echo "==================================${NC}"
echo ""
echo "To start the server, run:"
echo "  go run cmd/api/main.go"
echo ""
echo "Or use make:"
echo "  make run"
echo ""
echo "Then open your browser to:"
echo "  http://localhost:$PORT"
echo ""