# Makefile for Restart Life API

.PHONY: help build run test clean fmt lint deps docker

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

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run --rm -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

docker-push: docker-build ## Push Docker image to registry
	@echo "Pushing Docker image..."
	docker push $(DOCKER_IMAGE)

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