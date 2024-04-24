.PHONY: help build test clean fmt lint vet doc

PKG := github.com/janekbaraniewski/gobd2/gobd2

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

doc: ## Run godoc server and open the documentation in a web browser
	@echo "Starting godoc server at http://localhost:8080/"
	@godoc -http=:8080
