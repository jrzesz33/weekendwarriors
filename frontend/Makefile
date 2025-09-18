# Golf Gamez Frontend Makefile

.PHONY: build build-prod clean serve test deps lint help

# Variables
GOOS=js
GOARCH=wasm
OUTPUT_DIR=web/static/js
WASM_FILE=$(OUTPUT_DIR)/main.wasm
WASM_EXEC=$(OUTPUT_DIR)/wasm_exec.js
GO_FILES=$(shell find . -name "*.go" -not -path "./vendor/*")

# Default target
help: ## Show this help message
	@echo "Golf Gamez Frontend Build System"
	@echo "================================"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Dependencies
deps: ## Install Go dependencies
	@echo "Installing Go dependencies..."
	@go mod download
	@go mod tidy

# Development build
build: deps ## Build WebAssembly for development
	@echo "Building WebAssembly for development..."
	@mkdir -p $(OUTPUT_DIR)
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(WASM_FILE) main.go
	@cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" $(WASM_EXEC)
	@echo "Build complete: $(WASM_FILE)"
	@echo "Support file: $(WASM_EXEC)"

# Production build with optimizations
build-prod: deps ## Build WebAssembly for production with optimizations
	@echo "Building WebAssembly for production..."
	@mkdir -p $(OUTPUT_DIR)
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags="-s -w" \
		-gcflags="-l=4" \
		-o $(WASM_FILE) main.go
	@cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" $(WASM_EXEC)
	@echo "Production build complete: $(WASM_FILE)"

	# Compress if tools available
	@if command -v brotli >/dev/null 2>&1; then \
		echo "Compressing with Brotli..."; \
		brotli -f $(WASM_FILE); \
		echo "Compressed: $(WASM_FILE).br"; \
	fi

	@if command -v gzip >/dev/null 2>&1; then \
		echo "Compressing with Gzip..."; \
		gzip -f -k $(WASM_FILE); \
		echo "Compressed: $(WASM_FILE).gz"; \
	fi

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -f $(WASM_FILE)
	@rm -f $(WASM_FILE).br
	@rm -f $(WASM_FILE).gz
	@rm -f $(WASM_EXEC)
	@echo "Clean complete"

# Development server
serve: build ## Build and serve locally for development
	@echo "Starting development server..."
	@echo "Access at: http://localhost:8000"
	@cd web && python3 -m http.server 8000

# Production server
serve-prod: build-prod ## Build for production and serve locally
	@echo "Starting production server..."
	@echo "Access at: http://localhost:8080"
	@cd web && python3 -m http.server 8080

# Run tests
test: ## Run Go tests
	@echo "Running tests..."
	@go test -v ./...

# Linting
lint: ## Run linting tools
	@echo "Running Go linting..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, using go vet"; \
		go vet ./...; \
	fi

# Format code
fmt: ## Format Go code
	@echo "Formatting Go code..."
	@go fmt ./...

# Check for updates
check-updates: ## Check for dependency updates
	@echo "Checking for dependency updates..."
	@go list -u -m all

# Watch and rebuild (requires inotify-tools on Linux)
watch: ## Watch for changes and rebuild automatically
	@echo "Watching for changes..."
	@while true; do \
		inotifywait -r -e modify,create,delete --include=".*\.go$$" .; \
		echo "Changes detected, rebuilding..."; \
		make build; \
		echo "Rebuild complete"; \
	done

# Bundle size analysis
analyze: build-prod ## Analyze bundle size
	@echo "WebAssembly bundle size analysis:"
	@echo "================================="
	@if [ -f $(WASM_FILE) ]; then \
		echo "Uncompressed: $$(du -h $(WASM_FILE) | cut -f1)"; \
	fi
	@if [ -f $(WASM_FILE).br ]; then \
		echo "Brotli:       $$(du -h $(WASM_FILE).br | cut -f1)"; \
	fi
	@if [ -f $(WASM_FILE).gz ]; then \
		echo "Gzip:         $$(du -h $(WASM_FILE).gz | cut -f1)"; \
	fi

# Security check
security: ## Run security checks
	@echo "Running security checks..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not installed, skipping security scan"; \
	fi

# Install development tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Docker build (if Dockerfile exists)
docker-build: ## Build Docker image
	@if [ -f Dockerfile ]; then \
		echo "Building Docker image..."; \
		docker build -t golf-gamez-frontend .; \
	else \
		echo "Dockerfile not found"; \
	fi

# Show file sizes
sizes: ## Show file sizes
	@echo "File sizes:"
	@echo "==========="
	@find web/static -type f \( -name "*.css" -o -name "*.js" -o -name "*.wasm" \) -exec du -h {} \; | sort -h

# Validate HTML/CSS
validate: ## Validate HTML and CSS
	@echo "Validating HTML and CSS..."
	@if command -v html5validator >/dev/null 2>&1; then \
		html5validator web/index.html; \
	else \
		echo "html5validator not installed"; \
	fi

# Create production package
package: build-prod ## Create production package
	@echo "Creating production package..."
	@rm -rf dist
	@mkdir -p dist
	@cp -r web/* dist/
	@echo "Production package created in dist/"

# Quick development setup
dev-setup: deps install-tools build ## Quick development environment setup
	@echo "Development environment ready!"
	@echo "Run 'make serve' to start development server"

# CI/CD build
ci-build: deps test lint build-prod analyze ## CI/CD build pipeline
	@echo "CI/CD build completed successfully"