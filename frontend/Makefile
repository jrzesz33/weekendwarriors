# Golf Gamez Frontend Makefile

# Variables
BINARY_NAME=golf-gamez-frontend
BUILD_DIR=build
WEB_DIR=web
STATIC_DIR=$(WEB_DIR)/static

# Go build flags
GOOS=js
GOARCH=wasm
GO_BUILD_FLAGS=-ldflags="-w -s"

# Development server settings
DEV_PORT=8000
API_PORT=8080

.PHONY: help build build-dev run dev clean test lint fmt tidy deps check assets serve-assets

# Default target
help: ## Show this help message
	@echo "Golf Gamez Frontend Build System"
	@echo "================================"
	@echo ""
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: clean deps build-wasm copy-assets ## Build production version
	@echo "✅ Production build complete"

build-dev: clean deps build-wasm-dev copy-assets ## Build development version
	@echo "✅ Development build complete"

build-wasm: ## Build WebAssembly binary (production)
	@echo "🔨 Building WebAssembly binary (production)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/app.wasm ./cmd

build-wasm-dev: ## Build WebAssembly binary (development)
	@echo "🔨 Building WebAssembly binary (development)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/app.wasm ./cmd

copy-assets: ensure-wasm-exec ## Copy static assets to build directory
	@echo "📦 Copying static assets..."
	@mkdir -p $(BUILD_DIR)
	@cp -r $(WEB_DIR)/* $(BUILD_DIR)/
	@cp wasm_exec.js $(BUILD_DIR)/
	@echo "✅ Assets copied"

ensure-wasm-exec: ## Ensure wasm_exec.js is available
	@if [ ! -f wasm_exec.js ]; then \
		echo "📥 Downloading wasm_exec.js..."; \
		go env GOROOT > /tmp/goroot && \
		cp $$(cat /tmp/goroot)/lib/wasm/wasm_exec.js .; \
	fi

generate-icons: ## Generate PWA icons from source
	@echo "🎨 Generating PWA icons..."
	@mkdir -p $(STATIC_DIR)
	@# Note: You would typically use ImageMagick or similar here
	@echo "⚠️  Icon generation requires manual setup with your preferred tool"
	@echo "   Required sizes: 72, 96, 128, 144, 152, 180, 192, 384, 512"

# Development targets
dev: build-dev serve ## Build and serve development version

run: build serve ## Build and serve production version

serve: ## Serve the built application
	@echo "🚀 Starting development server on http://localhost:$(DEV_PORT)"
	@echo "📡 Make sure the API server is running on http://localhost:$(API_PORT)"
	@go run ./cmd/server/main.go -port=$(DEV_PORT) -dir=$(BUILD_DIR) &
	@../bin/golf-gamez

serve-production: ## Serve production build
	@echo "🚀 Starting production server on http://localhost:$(DEV_PORT)"
	@cd $(BUILD_DIR) && python3 -m http.server $(DEV_PORT) 2>/dev/null || \
	cd $(BUILD_DIR) && python -m SimpleHTTPServer $(DEV_PORT) 2>/dev/null || \
	cd $(BUILD_DIR) && php -S localhost:$(DEV_PORT) 2>/dev/null || \
	echo "❌ No web server available. Install Python or PHP to serve files."

watch: ## Watch for changes and rebuild (requires entr or similar)
	@echo "👀 Watching for changes..."
	@if command -v entr >/dev/null 2>&1; then \
		find . -name "*.go" -not -path "./build/*" | entr -r make build-dev; \
	else \
		echo "❌ 'entr' not found. Install it with your package manager for auto-rebuild."; \
		echo "   Falling back to manual build mode."; \
		make build-dev; \
	fi

# Quality and maintenance targets
test: ## Run tests
	@echo "🧪 Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

lint: ## Run linter
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "⚠️  golangci-lint not found. Install it for better linting."; \
		go vet ./...; \
	fi

fmt: ## Format Go code
	@echo "🎨 Formatting Go code..."
	@go fmt ./...

tidy: ## Tidy Go modules
	@echo "🧹 Tidying Go modules..."
	@go mod tidy

deps: tidy ## Install dependencies
	@echo "📦 Installing dependencies..."
	@go mod download

# Utility targets
clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	@rm -f $(BUILD_DIR)/app.wasm
	@rm -f $(BUILD_DIR)/wasm_exec.js
	@rm -f coverage.out coverage.html
	@rm -f wasm_exec.js

clean-all: clean ## Clean everything including dependencies
	@echo "🧹 Cleaning everything..."
	@go clean -modcache

check: fmt lint test ## Run all quality checks

# Docker targets
build-container: build docker-build docker-run ## Build WASM, Go App and Docker

docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	@docker build -t golf-gamez-frontend .

docker-run: docker-build ## Run Docker container
	@echo "🐳 Running Docker container... on $(DEV_PORT)"
	@docker run -p $(DEV_PORT):80 golf-gamez-frontend

# Deployment targets
deploy-prepare: build ## Prepare for deployment
	@echo "🚀 Preparing deployment package..."
	@tar -czf golf-gamez-frontend.tar.gz -C $(BUILD_DIR) .
	@echo "📦 Deployment package created: golf-gamez-frontend.tar.gz"

deploy-netlify: build ## Deploy to Netlify (requires netlify-cli)
	@echo "🌐 Deploying to Netlify..."
	@if command -v netlify >/dev/null 2>&1; then \
		netlify deploy --prod --dir=$(BUILD_DIR); \
	else \
		echo "❌ Netlify CLI not found. Install it with: npm install -g netlify-cli"; \
	fi

deploy-vercel: build ## Deploy to Vercel (requires vercel-cli)
	@echo "🌐 Deploying to Vercel..."
	@if command -v vercel >/dev/null 2>&1; then \
		vercel --prod $(BUILD_DIR); \
	else \
		echo "❌ Vercel CLI not found. Install it with: npm install -g vercel"; \
	fi

# Information targets
info: ## Show build information
	@echo "Golf Gamez Frontend Build Information"
	@echo "===================================="
	@echo "Go version:     $$(go version)"
	@echo "GOOS:           $(GOOS)"
	@echo "GOARCH:         $(GOARCH)"
	@echo "Build dir:      $(BUILD_DIR)"
	@echo "Web dir:        $(WEB_DIR)"
	@echo "Dev port:       $(DEV_PORT)"
	@echo "API port:       $(API_PORT)"
	@echo ""
	@echo "Dependencies:"
	@go list -m all

size: build ## Show build size information
	@echo "📏 Build size information:"
	@echo "========================="
	@du -h $(BUILD_DIR)/app.wasm
	@echo ""
	@echo "Detailed breakdown:"
	@ls -lh $(BUILD_DIR)/

# Development helpers
install-tools: ## Install development tools
	@echo "🔧 Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Development tools installed"

install-entr: ## Install entr for file watching (macOS)
	@if command -v brew >/dev/null 2>&1; then \
		brew install entr; \
	else \
		echo "❌ Homebrew not found. Install entr manually for auto-rebuild."; \
	fi

# Performance and optimization
optimize: ## Optimize build for production
	@echo "⚡ Optimizing build..."
	@make build
	@if command -v wasm-opt >/dev/null 2>&1; then \
		echo "🔧 Running wasm-opt..."; \
		wasm-opt -Oz $(BUILD_DIR)/app.wasm -o $(BUILD_DIR)/app.wasm; \
	else \
		echo "⚠️  wasm-opt not found. Install binaryen for better optimization."; \
	fi
	@echo "✅ Optimization complete"

analyze: build ## Analyze WebAssembly binary
	@echo "🔍 Analyzing WebAssembly binary..."
	@if command -v wasm-objdump >/dev/null 2>&1; then \
		wasm-objdump -h $(BUILD_DIR)/app.wasm; \
	else \
		echo "⚠️  wasm-objdump not found. Install wabt for analysis tools."; \
	fi

# Quick commands
quick: build-dev ## Quick development build (alias for build-dev)

prod: optimize ## Production build with optimizations

