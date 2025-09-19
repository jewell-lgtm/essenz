.PHONY: build test lint fmt vet clean install help check-tools

# Check tool versions
check-tools:
	@echo "Checking tool versions..."
	@go version
	@echo "Expected Go version: $(shell cat .tool-versions | grep golang | cut -d' ' -f2)"
	@if [ "$(shell go version | grep -o 'go[0-9]\+\.[0-9]\+\.[0-9]\+' | sed 's/go//')" != "$(shell cat .tool-versions | grep golang | cut -d' ' -f2)" ]; then \
		echo "Warning: Go version mismatch. Please run 'asdf install' to install the correct version."; \
	fi

# Build the binary
build:
	go build -o sz ./cmd/essenz

# Run tests
test:
	go test -v ./...

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

# Run all checks
check: check-tools fmt vet lint test

# Show help
help:
	@echo "Available targets:"
	@echo "  check-tools - Check tool versions against .tool-versions"
	@echo "  build       - Build the sz binary"
	@echo "  test        - Run tests"
	@echo "  lint        - Run golangci-lint"
	@echo "  fmt         - Format code with gofmt and goimports"
	@echo "  vet         - Run go vet"
	@echo "  clean       - Clean build artifacts"
	@echo "  install     - Install binary locally"
	@echo "  check       - Run all checks (tools, fmt, vet, lint, test)"
	@echo "  help        - Show this help message"