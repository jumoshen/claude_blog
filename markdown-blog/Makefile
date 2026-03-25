.PHONY: build run swag-init swag test clean fmt lint docker-build docker-up docker-down

# Build variables
BINARY_NAME=blog
BUILD_DIR=bin
CMD_DIR=cmd/server

# Docker variables
DOCKER_IMAGE=markdown-blog
DOCKER_TAG=latest

# Swagger
SWAG_CMD=swag init -g $(CMD_DIR)/main.go -o ./api

# Default target
all: build

# Build the application
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

# Run the application
run:
	@echo "Running..."
	go run ./$(CMD_DIR)/main.go

# Initialize swagger
swag-init:
	@echo "Initializing Swagger..."
	$(SWAG_CMD)

# Generate swagger docs (alias for swag-init)
swag:
	$(SWAG_CMD)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -rf ./api

# Format code
fmt:
	@echo "Formatting..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting..."
	golangci-lint run

# Run in Docker
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Docker Compose up
docker-up:
	@echo "Starting services..."
	docker-compose up -d

# Docker Compose down
docker-down:
	@echo "Stopping services..."
	docker-compose down

# Docker Compose build and up
docker: docker-build docker-up

# Stop and remove containers, volumes
docker-clean: docker-down
	docker volume prune -f

# Show help
help:
	@echo "Available targets:"
	@echo "  build       - Build the application"
	@echo "  run         - Run the application"
	@echo "  swag-init   - Initialize Swagger documentation"
	@echo "  swag        - Generate Swagger documentation"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-up   - Start Docker Compose services"
	@echo "  docker-down - Stop Docker Compose services"
	@echo "  docker      - Build and start Docker Compose services"
	@echo "  docker-clean - Clean Docker resources"
