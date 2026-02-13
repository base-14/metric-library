.PHONY: build test lint migrate migrate-down clean run fmt tidy ci ci-go ci-web docker-build docker-up docker-down docker-rebuild docker-logs \
	extract extract-otel extract-postgres extract-node extract-redis extract-clickhouse extract-cockroachdb extract-elasticsearch extract-memcached extract-nats extract-ksm extract-cadvisor extract-semconv extract-all enrich \
	extract-otel-python extract-otel-java extract-otel-dotnet extract-otel-go extract-otel-rust extract-otel-js extract-openllmetry extract-openlit \
	extract-gcp-compute extract-gcp-cloudsql extract-gcp-gke extract-gcp-loadbalancing extract-gcp-pubsub extract-gcp-cloudrun extract-gcp-storage extract-gcp-cloudfunctions \
	extract-azure-vm extract-azure-sqldatabase extract-azure-aks extract-azure-appgateway extract-azure-servicebus extract-azure-functions extract-azure-blobstorage extract-azure-cosmosdb \
	web-build web-test web-lint build-all test-all lint-all version version-set release

# Binary name
BINARY_NAME=metric-library

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOFMT=$(GOCMD) fmt
GOMOD=$(GOCMD) mod

# Database
DATABASE_URL?=sqlite:./data/metric-library.db
MIGRATIONS_DIR=./internal/store/migrations

# Build the application
build:
	CGO_ENABLED=1 $(GOBUILD) -tags "fts5" -o bin/$(BINARY_NAME) ./cmd/glossary

# Run the application
run: build
	./bin/$(BINARY_NAME)

# Run metric extraction (default: otel-collector-contrib)
extract: extract-otel

extract-otel: build
	./bin/$(BINARY_NAME) extract -adapter otel-collector-contrib

extract-postgres: build
	./bin/$(BINARY_NAME) extract -adapter prometheus-postgres

extract-node: build
	./bin/$(BINARY_NAME) extract -adapter prometheus-node

extract-redis: build
	./bin/$(BINARY_NAME) extract -adapter prometheus-redis

extract-mysql: build
	./bin/$(BINARY_NAME) extract -adapter prometheus-mysql

extract-mongodb: build
	./bin/$(BINARY_NAME) extract -adapter prometheus-mongodb

extract-kafka: build
	./bin/$(BINARY_NAME) extract -adapter prometheus-kafka

extract-clickhouse: build
	./bin/$(BINARY_NAME) extract -adapter prometheus-clickhouse

extract-cockroachdb: build
	./bin/$(BINARY_NAME) extract -adapter prometheus-cockroachdb

extract-elasticsearch: build
	./bin/$(BINARY_NAME) extract -adapter prometheus-elasticsearch

extract-memcached: build
	./bin/$(BINARY_NAME) extract -adapter prometheus-memcached

extract-nats: build
	./bin/$(BINARY_NAME) extract -adapter prometheus-nats

extract-ksm: build
	./bin/$(BINARY_NAME) extract -adapter kubernetes-ksm

extract-cadvisor: build
	./bin/$(BINARY_NAME) extract -adapter kubernetes-cadvisor

extract-semconv: build
	./bin/$(BINARY_NAME) extract -adapter otel-semconv

extract-otel-python: build
	./bin/$(BINARY_NAME) extract -adapter otel-python

extract-otel-java: build
	./bin/$(BINARY_NAME) extract -adapter otel-java

extract-otel-dotnet: build
	./bin/$(BINARY_NAME) extract -adapter otel-dotnet

extract-otel-go: build
	./bin/$(BINARY_NAME) extract -adapter otel-go

extract-otel-rust: build
	./bin/$(BINARY_NAME) extract -adapter otel-rust

extract-otel-js: build
	./bin/$(BINARY_NAME) extract -adapter otel-js

extract-openllmetry: build
	./bin/$(BINARY_NAME) extract -adapter openllmetry

extract-openlit: build
	./bin/$(BINARY_NAME) extract -adapter openlit

extract-all: build
	./bin/$(BINARY_NAME) extract -adapter otel-collector-contrib
	./bin/$(BINARY_NAME) extract -adapter otel-semconv
	./bin/$(BINARY_NAME) extract -adapter otel-python
	./bin/$(BINARY_NAME) extract -adapter otel-java
	./bin/$(BINARY_NAME) extract -adapter otel-dotnet
	./bin/$(BINARY_NAME) extract -adapter otel-go
	./bin/$(BINARY_NAME) extract -adapter otel-js
	./bin/$(BINARY_NAME) extract -adapter otel-rust
	./bin/$(BINARY_NAME) extract -adapter prometheus-postgres
	./bin/$(BINARY_NAME) extract -adapter prometheus-node
	./bin/$(BINARY_NAME) extract -adapter prometheus-redis
	./bin/$(BINARY_NAME) extract -adapter prometheus-mysql
	./bin/$(BINARY_NAME) extract -adapter prometheus-mongodb
	./bin/$(BINARY_NAME) extract -adapter prometheus-kafka
	./bin/$(BINARY_NAME) extract -adapter prometheus-clickhouse
	./bin/$(BINARY_NAME) extract -adapter prometheus-cockroachdb
	./bin/$(BINARY_NAME) extract -adapter prometheus-elasticsearch
	./bin/$(BINARY_NAME) extract -adapter prometheus-memcached
	./bin/$(BINARY_NAME) extract -adapter prometheus-nats
	./bin/$(BINARY_NAME) extract -adapter kubernetes-ksm
	./bin/$(BINARY_NAME) extract -adapter kubernetes-cadvisor
	./bin/$(BINARY_NAME) extract -adapter openllmetry
	./bin/$(BINARY_NAME) extract -adapter openlit
	./bin/$(BINARY_NAME) extract -adapter cloudwatch-ec2
	./bin/$(BINARY_NAME) extract -adapter cloudwatch-rds
	./bin/$(BINARY_NAME) extract -adapter cloudwatch-lambda
	./bin/$(BINARY_NAME) extract -adapter cloudwatch-s3
	./bin/$(BINARY_NAME) extract -adapter cloudwatch-dynamodb
	./bin/$(BINARY_NAME) extract -adapter cloudwatch-alb
	./bin/$(BINARY_NAME) extract -adapter cloudwatch-sqs
	./bin/$(BINARY_NAME) extract -adapter cloudwatch-apigateway
	./bin/$(BINARY_NAME) extract -adapter gcp-compute
	./bin/$(BINARY_NAME) extract -adapter gcp-cloudsql
	./bin/$(BINARY_NAME) extract -adapter gcp-gke
	./bin/$(BINARY_NAME) extract -adapter gcp-loadbalancing
	./bin/$(BINARY_NAME) extract -adapter gcp-pubsub
	./bin/$(BINARY_NAME) extract -adapter gcp-cloudrun
	./bin/$(BINARY_NAME) extract -adapter gcp-storage
	./bin/$(BINARY_NAME) extract -adapter gcp-cloudfunctions
	./bin/$(BINARY_NAME) extract -adapter azure-vm
	./bin/$(BINARY_NAME) extract -adapter azure-sqldatabase
	./bin/$(BINARY_NAME) extract -adapter azure-aks
	./bin/$(BINARY_NAME) extract -adapter azure-appgateway
	./bin/$(BINARY_NAME) extract -adapter azure-servicebus
	./bin/$(BINARY_NAME) extract -adapter azure-functions
	./bin/$(BINARY_NAME) extract -adapter azure-blobstorage
	./bin/$(BINARY_NAME) extract -adapter azure-cosmosdb

# Individual CloudWatch extractions
extract-cloudwatch-ec2: build
	./bin/$(BINARY_NAME) extract -adapter cloudwatch-ec2

extract-cloudwatch-rds: build
	./bin/$(BINARY_NAME) extract -adapter cloudwatch-rds

extract-cloudwatch-lambda: build
	./bin/$(BINARY_NAME) extract -adapter cloudwatch-lambda

extract-cloudwatch-s3: build
	./bin/$(BINARY_NAME) extract -adapter cloudwatch-s3

extract-cloudwatch-dynamodb: build
	./bin/$(BINARY_NAME) extract -adapter cloudwatch-dynamodb

extract-cloudwatch-alb: build
	./bin/$(BINARY_NAME) extract -adapter cloudwatch-alb

extract-cloudwatch-sqs: build
	./bin/$(BINARY_NAME) extract -adapter cloudwatch-sqs

extract-cloudwatch-apigateway: build
	./bin/$(BINARY_NAME) extract -adapter cloudwatch-apigateway

# Individual GCP extractions
extract-gcp-compute: build
	./bin/$(BINARY_NAME) extract -adapter gcp-compute

extract-gcp-cloudsql: build
	./bin/$(BINARY_NAME) extract -adapter gcp-cloudsql

extract-gcp-gke: build
	./bin/$(BINARY_NAME) extract -adapter gcp-gke

extract-gcp-loadbalancing: build
	./bin/$(BINARY_NAME) extract -adapter gcp-loadbalancing

extract-gcp-pubsub: build
	./bin/$(BINARY_NAME) extract -adapter gcp-pubsub

extract-gcp-cloudrun: build
	./bin/$(BINARY_NAME) extract -adapter gcp-cloudrun

extract-gcp-storage: build
	./bin/$(BINARY_NAME) extract -adapter gcp-storage

extract-gcp-cloudfunctions: build
	./bin/$(BINARY_NAME) extract -adapter gcp-cloudfunctions

# Individual Azure extractions
extract-azure-vm: build
	./bin/$(BINARY_NAME) extract -adapter azure-vm

extract-azure-sqldatabase: build
	./bin/$(BINARY_NAME) extract -adapter azure-sqldatabase

extract-azure-aks: build
	./bin/$(BINARY_NAME) extract -adapter azure-aks

extract-azure-appgateway: build
	./bin/$(BINARY_NAME) extract -adapter azure-appgateway

extract-azure-servicebus: build
	./bin/$(BINARY_NAME) extract -adapter azure-servicebus

extract-azure-functions: build
	./bin/$(BINARY_NAME) extract -adapter azure-functions

extract-azure-blobstorage: build
	./bin/$(BINARY_NAME) extract -adapter azure-blobstorage

extract-azure-cosmosdb: build
	./bin/$(BINARY_NAME) extract -adapter azure-cosmosdb

# Enrich metrics with semconv data
enrich: build
	./bin/$(BINARY_NAME) enrich

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

# Web targets
web-build:
	cd web && npm run build

web-test:
	cd web && npm test

web-lint:
	cd web && npm run lint

# Combined targets
build-all: build web-build

test-all: test web-test

lint-all: lint web-lint

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
	COMPOSE_BAKE=true docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-rebuild: docker-down docker-build docker-up

docker-logs:
	docker compose logs -f

# Version management
VERSION_FILE=VERSION
CURRENT_VERSION=$(shell cat $(VERSION_FILE))

version:
	@echo $(CURRENT_VERSION)

version-set:
	@if [ -z "$(V)" ]; then echo "Usage: make version-set V=x.y.z"; exit 1; fi
	@echo "$(V)" > $(VERSION_FILE)
	@sed -i '' 's/appVersion: ".*"/appVersion: "$(V)"/' deploy/helm/metric-library/Chart.yaml
	@sed -i '' 's/tag: ".*"/tag: "$(V)"/' deploy/helm/metric-library/values.yaml
	@sed -i '' 's/tag: ".*"/tag: "$(V)"/' local/values.yaml
	@echo "Version updated to $(V)"
	@echo "Files updated:"
	@echo "  - VERSION"
	@echo "  - deploy/helm/metric-library/Chart.yaml (appVersion)"
	@echo "  - deploy/helm/metric-library/values.yaml (image tags)"
	@echo "  - local/values.yaml (image tags)"

release: version
	@echo "Creating release tags for version $(CURRENT_VERSION)..."
	git tag -a api-v$(CURRENT_VERSION) -m "API release $(CURRENT_VERSION)"
	git tag -a web-v$(CURRENT_VERSION) -m "Web release $(CURRENT_VERSION)"
	@echo "Tags created. Push with: git push --tags"

# Default target
all: check build
