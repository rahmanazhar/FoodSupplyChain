# Makefile for FoodSupplyChain project

.PHONY: all build test clean docker-build docker-compose-up docker-compose-down

# Variables
BINARY_NAME=foodsupplychain
DOCKER_COMPOSE=docker-compose

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean

# Build all services
all: clean build

# Build the project
build:
	@echo "Building services..."
	cd cmd/inventory && $(GOBUILD)
	cd cmd/shipment && $(GOBUILD)
	cd cmd/gateway && $(GOBUILD)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) ./...

# Clean build files
clean:
	@echo "Cleaning build files..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Docker commands
docker-build:
	@echo "Building Docker images..."
	docker build -t foodsupplychain-inventory:latest -f build/inventory/Dockerfile .
	docker build -t foodsupplychain-shipment:latest -f build/shipment/Dockerfile .
	docker build -t foodsupplychain-gateway:latest -f build/gateway/Dockerfile .

# Docker Compose commands
docker-compose-up:
	@echo "Starting services with Docker Compose..."
	$(DOCKER_COMPOSE) up -d

docker-compose-down:
	@echo "Stopping services..."
	$(DOCKER_COMPOSE) down

# Development helpers
dev-deps:
	@echo "Installing development dependencies..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	@echo "Running linter..."
	golangci-lint run

# Database commands
db-migrate:
	@echo "Running database migrations..."
	# Add migration command here

# Generate API documentation
docs:
	@echo "Generating API documentation..."
	# Add swagger/openapi generation command here
