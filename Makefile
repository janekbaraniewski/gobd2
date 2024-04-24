.PHONY: help build test clean fmt lint vet doc

PKG := github.com/janekbaraniewski/gobd2/gobd2
BINARY_NAME=gobd2-cli
BUILD_DIR=./bin
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

help:  ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_\-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-30s %s\n", $$1, $$2}'

build: ## Build the project
	@echo "Building the project..."
	go build $(PKG)

.PHONY: build-cli
build-cli:
	@echo "Building the CLI application for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH) ./cmd
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH)"

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
