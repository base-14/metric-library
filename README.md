# OTel Glossary

A metric discovery platform for cataloging and searching OpenTelemetry metrics from various sources.

## Purpose

OTel Glossary extracts metric definitions from OpenTelemetry Collector components (receivers, processors, exporters) and other sources, storing them in a searchable catalog. It provides a web interface for discovering metrics by name, type, component, and source.

## Tech Stack

| Component | Version |
|-----------|---------|
| Go | 1.25 |
| SQLite | FTS5 (full-text search) |
| Node.js | 20 |
| Next.js | 16.1.1 |
| golangci-lint | v2.8.0 |

## Quick Start

```bash
# Start with Docker
make docker-up

# API: http://localhost:8080
# Web: http://localhost:3000
```

## Development

### Prerequisites

- Go 1.25+
- Node.js 20+
- Docker (optional)

### Commands

```bash
# Build
make build            # Build Go binary
cd web && make build  # Build frontend

# Extract metrics
make extract          # Extract metrics from otel-collector-contrib

# Test
make test             # Run Go tests
cd web && make test   # Run frontend tests

# Lint
make lint             # Run golangci-lint
cd web && make lint   # Run ESLint

# CI (runs everything)
make ci               # Full CI: lint, test, build (Go + Web)

# Docker
make docker-build     # Build images
make docker-up        # Start services
make docker-down      # Stop services
make docker-logs      # View logs

# Database
make migrate          # Run migrations (via dbmate)
make migrate-status   # Check migration status
```

### Project Structure

```
otel-glossary/
├── cmd/glossary/          # Main entry point
├── internal/
│   ├── api/               # REST API handlers
│   ├── store/             # SQLite store + migrations
│   ├── domain/            # Domain models
│   ├── adapter/           # Source adapters
│   ├── fetcher/           # Git fetcher
│   ├── discovery/         # Metadata discovery
│   ├── parser/            # YAML parser
│   └── extractor/         # Metric extractor
├── web/                   # Next.js frontend
├── Dockerfile             # Go backend
├── docker-compose.yml     # Local development
└── Makefile
```

## API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /health` | Health check |
| `GET /api/metrics` | Search metrics (supports filters) |
| `GET /api/metrics/{id}` | Get single metric |
| `GET /api/facets` | Get facet counts for filtering |

### Query Parameters

- `q` - Full-text search
- `instrument_type` - Filter by type (counter, gauge, histogram, etc.)
- `component_type` - Filter by component (receiver, processor, exporter)
- `component_name` - Filter by component name
- `source_category` - Filter by source
- `limit`, `offset` - Pagination

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | API server port |
| `DATABASE_PATH` | ./data/glossary.db | SQLite database path |
| `NEXT_PUBLIC_API_URL` | http://localhost:8080 | API URL for frontend |
