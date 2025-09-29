# KineticaFS Makefile
# Provides comprehensive build and development workflow

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet

# Application parameters
BINARY_NAME=kineticafs
BINARY_UNIX=$(BINARY_NAME)_unix
DOCKER_IMAGE=ghcr.io/argon-chat/kineticafs
DOCS_DIR=docs

# Build flags
LDFLAGS=-ldflags "-s -w"
BUILDFLAGS=
PRODFLAGS=-tags prod

# Colors for terminal output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: help all build clean test deps lint format docs run docker-build docker-run security coverage install dev-setup check

# Default target
all: clean deps format lint test build docs

## Show available targets
help:
	@echo "$(BLUE)KineticaFS Build System$(NC)"
	@echo ""
	@echo "$(GREEN)Available targets:$(NC)"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /' | sed 's/:.*//g'
	@echo ""
	@echo "$(YELLOW)Examples:$(NC)"
	@echo "  make all          # Run full build pipeline"
	@echo "  make dev-setup    # Setup development environment"
	@echo "  make run          # Run the application locally"
	@echo "  make docker-build # Build Docker image"

## Install project dependencies
deps:
	@echo "$(BLUE)Installing dependencies...$(NC)"
	$(GOMOD) download
	$(GOMOD) verify

## Setup development environment
dev-setup: deps
	@echo "$(BLUE)Setting up development environment...$(NC)"
	@if ! command -v golangci-lint > /dev/null; then \
		echo "$(YELLOW)Installing golangci-lint...$(NC)"; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2; \
	fi
	@if ! command -v swag > /dev/null; then \
		echo "$(YELLOW)Installing swag...$(NC)"; \
		$(GOGET) -u github.com/swaggo/swag/cmd/swag; \
	fi
	@echo "$(GREEN)Development environment ready!$(NC)"

## Build the application
build: deps docs
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	$(GOBUILD) $(BUILDFLAGS) $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "$(GREEN)Build completed: $(BINARY_NAME)$(NC)"

## Build production binary
build-prod: deps docs
	@echo "$(BLUE)Building production $(BINARY_NAME)...$(NC)"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILDFLAGS) $(PRODFLAGS) $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "$(GREEN)Production build completed: $(BINARY_NAME)$(NC)"

## Build for Unix/Linux
build-linux: deps docs
	@echo "$(BLUE)Building for Linux...$(NC)"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILDFLAGS) $(LDFLAGS) -o $(BINARY_UNIX) .
	@echo "$(GREEN)Linux build completed: $(BINARY_UNIX)$(NC)"

## Run tests with coverage

## Run tests with coverage (Go 1.22+ covdata format)
test: docs
	@echo "$(BLUE)Running tests...$(NC)"
	$(GOTEST) -v -race -cover ./...
	@echo "$(GREEN)Tests completed$(NC)"

## Generate test coverage report (Go 1.22+ covdata format)
coverage: test
	@echo "$(BLUE)Generating coverage report...$(NC)"
	$(GOCMD) tool covdata textfmt -i=default -o coverage.txt
	$(GOCMD) tool cover -html=coverage.txt -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

## Run linters
lint:
	@echo "$(BLUE)Running linters...$(NC)"
	@export PATH=$$PATH:$(shell go env GOPATH)/bin && \
	if command -v golangci-lint > /dev/null; then \
		golangci-lint run --timeout=5m; \
		echo "$(GREEN)Linting completed$(NC)"; \
	else \
		echo "$(YELLOW)golangci-lint not found. Run 'make dev-setup' first.$(NC)"; \
		$(GOVET) ./...; \
	fi

## Format code
format:
	@echo "$(BLUE)Formatting code...$(NC)"
	$(GOFMT) -s -w .
	@if command -v goimports > /dev/null; then \
		goimports -w .; \
	fi
	@echo "$(GREEN)Code formatted$(NC)"

## Generate swagger documentation
docs: deps
	@echo "$(BLUE)Generating API documentation...$(NC)"
	@export PATH=$$PATH:$(shell go env GOPATH)/bin && \
	if command -v swag > /dev/null; then \
		swag init --generalInfo main.go --output $(DOCS_DIR); \
		echo "$(GREEN)Documentation generated in $(DOCS_DIR)/$(NC)"; \
	else \
		echo "$(YELLOW)swag not found. Run 'make dev-setup' first.$(NC)"; \
	fi

## Run security scan
security:
	@echo "$(BLUE)Running security scan...$(NC)"
	@export PATH=$$PATH:$(shell go env GOPATH)/bin && \
	if command -v gosec > /dev/null; then \
		gosec ./...; \
		echo "$(GREEN)Security scan completed$(NC)"; \
	else \
		echo "$(YELLOW)gosec not found. Run 'make dev-setup' first.$(NC)"; \
	fi

## Clean build artifacts
clean:
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f coverage.out
	rm -f coverage.html
	rm -rf $(DOCS_DIR)
	@echo "$(GREEN)Cleanup completed$(NC)"

## Run the application locally (server mode)
run: build
	@echo "$(BLUE)Starting KineticaFS server on port 3000...$(NC)"
	./$(BINARY_NAME) --server --port 3000

## Run the application with migrations
run-migrate: build
	@echo "$(BLUE)Starting KineticaFS server with migrations...$(NC)"
	./$(BINARY_NAME) --server --migrate --port 3000

## Install the binary to GOPATH/bin
install: build
	@echo "$(BLUE)Installing $(BINARY_NAME) to $(shell go env GOPATH)/bin...$(NC)"
	cp $(BINARY_NAME) $(shell go env GOPATH)/bin/
	@echo "$(GREEN)Installation completed$(NC)"

## Build Docker image
docker-build:
	@echo "$(BLUE)Building Docker image...$(NC)"
	docker build -t $(DOCKER_IMAGE):latest .
	@echo "$(GREEN)Docker image built: $(DOCKER_IMAGE):latest$(NC)"

## Run Docker container
docker-run: docker-build
	@echo "$(BLUE)Running Docker container...$(NC)"
	docker run -p 3000:3000 $(DOCKER_IMAGE):latest

## Run dev with docker compose
docker-dev: docs
	@echo "$(BLUE)Starting development environment with Docker Compose...$(NC)"
	docker compose down; docker volume prune -af; docker compose up -d; sleep 10; go run . -m -s

## Run comprehensive checks (CI pipeline)
check: deps format lint test build docs
	@echo "$(GREEN)All checks passed!$(NC)"

## Watch for changes and rebuild (requires 'entr' tool)
watch:
	@echo "$(BLUE)Watching for changes...$(NC)"
	@if command -v find > /dev/null && command -v entr > /dev/null; then \
		find . -name '*.go' | entr -r make build run; \
	else \
		echo "$(YELLOW)This target requires 'entr' tool. Install it with your package manager.$(NC)"; \
	fi