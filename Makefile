.PHONY: build test lint fmt vet clean install help

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
check: fmt vet lint test

# Show help
help:
	@echo "Available targets:"
	@echo "  build    - Build the sz binary"
	@echo "  test     - Run tests"
	@echo "  lint     - Run golangci-lint"
	@echo "  fmt      - Format code with gofmt and goimports"
	@echo "  vet      - Run go vet"
	@echo "  clean    - Clean build artifacts"
	@echo "  install  - Install binary locally"
	@echo "  check    - Run fmt, vet, lint, and test"
	@echo "  help     - Show this help message"