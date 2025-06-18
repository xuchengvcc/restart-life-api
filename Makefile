# Makefile for Restart Life API

.PHONY: help build run test clean fmt lint deps docker docker-build docker-up docker-down docker-logs

# Variables
APP_NAME := restart-life-api
VERSION := v0.1.0
BUILD_DIR := build
DOCKER_IMAGE := $(APP_NAME):$(VERSION)

# Go related variables
GO_VERSION := 1.21
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := gofmt
GOLINT := golangci-lint

# Build related
MAIN_PATH := ./cmd/server
BINARY_NAME := $(APP_NAME)
BINARY_PATH := $(BUILD_DIR)/$(BINARY_NAME)

help: ## Show this help message
	@echo 'Management commands for $(APP_NAME):'
	@echo ''
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) { \
			printf "  %-20s%s\n", $$1, $$2 \
		} \
	}' $(MAKEFILE_LIST)

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

build: deps ## Build the application
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-w -s" -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Build completed: $(BINARY_PATH)"

build-local: deps ## Build for local development
	@echo "Building $(APP_NAME) for local..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Local build completed: $(BINARY_PATH)"

run: ## Run the application
	@echo "Running $(APP_NAME)..."
	$(GOCMD) run $(MAIN_PATH)

dev: ## Run in development mode with hot reload
	@echo "Starting development server..."
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air

test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	@echo "Generating coverage report..."
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

bench: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

fmt: ## Format Go code
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	$(GOCMD) mod tidy

lint: ## Run linter
	@echo "Running linter..."
	@which $(GOLINT) > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	$(GOLINT) run

clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .
	@echo "Docker image built: $(DOCKER_IMAGE)"

docker-up: ## Start all services with docker-compose
	@echo "Starting all services..."
	docker-compose up -d
	@echo "Services started! Visit http://localhost:8080"

docker-down: ## Stop all services
	@echo "Stopping all services..."
	docker-compose down

docker-down-volumes: ## Stop all services and remove volumes
	@echo "Stopping all services and removing volumes..."
	docker-compose down -v

docker-logs: ## Show logs from all services
	@echo "Showing logs..."
	docker-compose logs -f

docker-logs-app: ## Show logs from app service only
	@echo "Showing app logs..."
	docker-compose logs -f app

docker-restart: ## Restart all services
	@echo "Restarting all services..."
	docker-compose restart

docker-rebuild: ## Rebuild and restart the app service
	@echo "Rebuilding and restarting app..."
	docker-compose build app
	docker-compose restart app

docker-shell: ## Get shell access to app container
	@echo "Entering app container..."
	docker exec -it restart-life-api sh

docker-mysql: ## Connect to MySQL in container
	@echo "Connecting to MySQL..."
	docker exec -it restart-mysql mysql -u root -p

docker-redis: ## Connect to Redis in container
	@echo "Connecting to Redis..."
	docker exec -it restart-redis redis-cli

docker-clean: ## Clean up Docker resources
	@echo "Cleaning up Docker resources..."
	docker system prune -f
	docker volume prune -f

docker-setup: docker-build docker-up ## Quick setup: build and start services
	@echo "Docker setup completed!"
	@echo "Services:"
	@echo "  - API Server: http://localhost:8080"
	@echo "  - Health Check: http://localhost:8080/health"
	@echo "  - Database Admin: http://localhost:8081"
	@echo "  - Redis Commander: http://localhost:8082"

migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	@which migrate > /dev/null || (echo "Please install golang-migrate" && exit 1)
	migrate -path migrations -database "mysql://root:password@tcp(localhost:3306)/restart_life_dev" up

migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	migrate -path migrations -database "mysql://root:password@tcp(localhost:3306)/restart_life_dev" down

migrate-create: ## Create new migration file (usage: make migrate-create NAME=migration_name)
	@echo "Creating migration: $(NAME)..."
	@test -n "$(NAME)" || (echo "NAME is required. Usage: make migrate-create NAME=migration_name" && exit 1)
	migrate create -ext sql -dir migrations $(NAME)

setup-dev: ## Setup development environment
	@echo "Setting up development environment..."
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Development environment setup complete!"

install: build ## Install the application
	@echo "Installing $(APP_NAME)..."
	cp $(BINARY_PATH) /usr/local/bin/
	@echo "Installation completed!"

.DEFAULT_GOAL := help 