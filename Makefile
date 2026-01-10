.PHONY: build test lint migrate migrate-down clean run fmt tidy

# Binary name
BINARY_NAME=glossary

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOFMT=$(GOCMD) fmt
GOMOD=$(GOCMD) mod

# Database
DATABASE_URL?=sqlite:./data/glossary.db

# Build the application
build:
	CGO_ENABLED=1 $(GOBUILD) -tags "fts5" -o bin/$(BINARY_NAME) ./cmd/glossary

# Run the application
run: build
	./bin/$(BINARY_NAME)

# Run tests
test:
	CGO_ENABLED=1 $(GOTEST) -v -race -cover -tags "fts5" ./...

# Run tests with coverage report
test-coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	golangci-lint run ./...

# Format code
fmt:
	$(GOFMT) ./...

# Tidy dependencies
tidy:
	$(GOMOD) tidy

# Run database migrations
migrate:
	DATABASE_URL=$(DATABASE_URL) dbmate up

# Rollback last migration
migrate-down:
	DATABASE_URL=$(DATABASE_URL) dbmate down

# Create a new migration
migrate-new:
	@read -p "Migration name: " name; \
	DATABASE_URL=$(DATABASE_URL) dbmate new $$name

# Show migration status
migrate-status:
	DATABASE_URL=$(DATABASE_URL) dbmate status

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install development dependencies
dev-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/amacneil/dbmate/v2@latest

# Run all checks (format, lint, test)
check: fmt lint test

# Default target
all: check build
