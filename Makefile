.PHONY: build test lint migrate migrate-down clean run fmt tidy ci ci-go ci-web docker-build docker-up docker-down docker-logs extract

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
MIGRATIONS_DIR=./internal/store/migrations

# Build the application
build:
	CGO_ENABLED=1 $(GOBUILD) -tags "fts5" -o bin/$(BINARY_NAME) ./cmd/glossary

# Run the application
run: build
	./bin/$(BINARY_NAME)

# Run metric extraction
extract: build
	./bin/$(BINARY_NAME) extract

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
	DATABASE_URL=$(DATABASE_URL) dbmate -d $(MIGRATIONS_DIR) up

# Rollback last migration
migrate-down:
	DATABASE_URL=$(DATABASE_URL) dbmate -d $(MIGRATIONS_DIR) down

# Create a new migration
migrate-new:
	@read -p "Migration name: " name; \
	DATABASE_URL=$(DATABASE_URL) dbmate -d $(MIGRATIONS_DIR) new $$name

# Show migration status
migrate-status:
	DATABASE_URL=$(DATABASE_URL) dbmate -d $(MIGRATIONS_DIR) status

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

# CI targets - mirrors GitHub Actions workflow
ci-go: lint test build

ci-web:
	cd web && npm ci && npm run lint && npm test && npm run build

ci: ci-go ci-web

# Docker targets
docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

# Default target
all: check build
