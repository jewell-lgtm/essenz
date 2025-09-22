.PHONY: all build test lint fmt vet clean install help check-tools install-tools setup-pre-commit check

# Default target
all: build

# Check tool versions against .tool-versions
check-tools:
	@echo "Checking tool versions against .tool-versions..."
	@if command -v asdf >/dev/null 2>&1; then \
		echo "Using asdf-managed tools:"; \
		asdf current; \
		echo ""; \
		echo "Expected versions from .tool-versions:"; \
		cat .tool-versions; \
		echo ""; \
		if ! asdf current | grep -q "$(shell cat .tool-versions | grep golang | cut -d' ' -f2)"; then \
			echo "⚠️  Go version mismatch detected. Run 'asdf install' to fix."; \
		else \
			echo "✅ All tool versions match .tool-versions"; \
		fi; \
	else \
		echo "asdf not found, checking system tools:"; \
		go version; \
		echo "Expected Go version: $(shell cat .tool-versions | grep golang | cut -d' ' -f2)"; \
		if [ "$(shell go version | grep -o 'go[0-9]\+\.[0-9]\+\.[0-9]\+' | sed 's/go//')" != "$(shell cat .tool-versions | grep golang | cut -d' ' -f2)" ]; then \
			echo "⚠️  Go version mismatch. Please install asdf and run 'asdf install'."; \
		fi; \
	fi

# Build the binary
build:
	go build -o sz ./cmd/essenz

# Run tests
test:
	go test -v ./...

# Run tests in sandboxed environments (local caches and skipping env-dependent specs)
test-sandbox:
	mkdir -p .gocache
	ESSENZ_SANDBOX=1 GOCACHE=$(PWD)/.gocache go test -v ./...

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Run go vet
vet:
	go vet ./...

# Clean build artifacts
clean:
	rm -f sz

# Install binary locally
install:
	go install ./cmd/essenz

# Install tools via asdf
install-tools:
	@if command -v asdf >/dev/null 2>&1; then \
		echo "Installing tools from .tool-versions via asdf..."; \
		asdf install; \
		echo "✅ Tools installed successfully"; \
	else \
		echo "❌ asdf not found. Please install asdf first:"; \
		echo "   https://asdf-vm.com/guide/getting-started.html"; \
		exit 1; \
	fi

# Setup pre-commit hooks
setup-pre-commit:
	pre-commit install
	pre-commit install --hook-type commit-msg

# Run all checks
check: check-tools fmt vet lint test

# Show help
help:
	@echo "Available targets:"
	@echo "  all          - Default target (build)"
	@echo "  build        - Build the sz binary"
	@echo "  check-tools  - Check tool versions"
	@echo "  check        - Run all checks"
