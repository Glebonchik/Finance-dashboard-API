.PHONY: build run test clean docker-up docker-down migrate-up migrate-down swagger help

# Variables
BINARY_NAME=finance-dashboard
CMD_PATH=./cmd/api
MAIN_GO=$(CMD_PATH)/main.go
DOCKER_COMPOSE=docker-compose.yml

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o bin/$(BINARY_NAME) $(MAIN_GO)
	@echo "Build complete!"

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	go run $(MAIN_GO)

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "Clean complete!"

# Start Docker containers
docker-up:
	@echo "Starting Docker containers..."
	docker-compose up -d
	@echo "Containers started!"

# Stop Docker containers
docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down
	@echo "Containers stopped!"

# Run database migrations up
migrate-up:
	@echo "Running migrations up..."
	go run cmd/migrate/main.go -action up

# Run database migrations down
migrate-down:
	@echo "Running migrations down..."
	go run cmd/migrate/main.go -action down

# Force migration state
migrate-force:
	@echo "Forcing migration state..."
	go run cmd/migrate/main.go -action force

# Check migration version
migrate-version:
	@echo "Checking migration version..."
	go run cmd/migrate/main.go -action version

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	swag init --dir . --generalInfo cmd/api/main.go --output docs
	@echo "Swagger documentation generated!"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Lint the code (requires golangci-lint)
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Help
help:
	@echo "Available targets:"
	@echo "  build           - Build the application"
	@echo "  run             - Run the application"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  clean           - Clean build artifacts"
	@echo "  docker-up       - Start Docker containers"
	@echo "  docker-down     - Stop Docker containers"
	@echo "  migrate-up      - Run database migrations up"
	@echo "  migrate-down    - Run database migrations down"
	@echo "  migrate-force   - Force migration state"
	@echo "  migrate-version - Check migration version"
	@echo "  swagger         - Generate Swagger documentation"
	@echo "  deps            - Download dependencies"
	@echo "  lint            - Run linter (requires golangci-lint)"
	@echo "  fmt             - Format code"
	@echo "  help            - Show this help message"
