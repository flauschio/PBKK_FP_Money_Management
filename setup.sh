RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' 
echo -n "Checking Go installation"
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗${NC}"
    echo "Go is not installed. Please install Go 1.21 or higher from https://golang.org/dl/"
    exit 1
fi
echo -e "${GREEN}✓${NC}"
echo -n "Checking PostgreSQL installation"
if ! command -v psql &> /dev/null; then
    echo -e "${RED}✗${NC}"
    echo "PostgreSQL is not installed. Please install PostgreSQL from https://www.postgresql.org/download/"
    exit 1
fi
echo -e "${GREEN}✓${NC}"
echo -e "${BLUE}Database Configuration${NC}"
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
echo ""
echo -n "Creating .env file"
cat > .env << EOF
DB_HOST=$DB_HOST
DB_PORT=$DB_PORT
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
DB_NAME=$DB_NAME
PORT=$PORT
GIN_MODE=release
EOF
echo -e "${GREEN}✓${NC}"
export DB_HOST DB_PORT DB_USER DB_PASSWORD DB_NAME PORT
echo -n "Creating database"
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "CREATE DATABASE $DB_NAME;" 2>/dev/null
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${YELLOW}⚠${NC} (Database might already exist)"
fi
echo -n "Running database migrations"
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f migrations/001_create_tables.sql > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    echo "Failed to run migrations. Please check your database connection."
    exit 1
fi
echo -n "Installing Go dependencies"
go mod download > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    echo "Failed to install dependencies."
    exit 1
fi
echo -e "${GREEN}Setup completed successfully!${NC}"
echo -e "${BLUE}To start the server, run:${NC}"
echo "go run cmd/api/main.go"
echo -e "${BLUE}Or use make:${NC}"
echo "make run"
echo -e "${BLUE}Then open your browser to:${NC}"
echo "http://localhost:$PORT"