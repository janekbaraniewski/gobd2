.PHONY: help build test clean fmt lint vet

PKG := ./...

help:  ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_\-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-30s %s\n", $$1, $$2}'

build: ## Build the project
	@echo "Building the project..."
	go build $(PKG)

test: ## Run tests
	@echo "Running tests..."
	go test $(PKG) -v

clean: ## Clean up artifacts
	@echo "Cleaning up..."
	go clean

fmt: ## Format the Go code
	@echo "Formatting Go code..."
	go fmt $(PKG)

lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	golangci-lint run ./...

vet: ## Vet examines Go source code and reports suspicious constructs
	@echo "Running Go vet..."
	go vet $(PKG)
