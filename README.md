# metric-library

A metric discovery platform for cataloging and searching metrics from OpenTelemetry, Prometheus exporters, and more.

## Purpose

metric-library extracts metric definitions from OpenTelemetry Collector components, Prometheus exporters, Kubernetes tooling, and other sources, storing them in a searchable catalog. It provides a web interface for discovering metrics by name, type, component, and source.

## Architecture

```
┌───────────────────────────────────────────────────────────────────────────────────────┐
│                                     Sources                                           │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌────────────┐           │
│  │otel-contrib│ │  postgres  │ │   redis    │ │    ksm     │ │  cadvisor  │    ...    │
│  │   (yaml)   │ │  (go ast)  │ │  (go ast)  │ │  (go ast)  │ │  (go ast)  │           │
│  └─────┬──────┘ └─────┬──────┘ └─────┬──────┘ └─────┬──────┘ └─────┬──────┘           │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌────────────┐           │
│  │otel-python │ │ otel-java  │ │  otel-js   │ │openllmetry │ │  openlit   │           │
│  │ (py ast)   │ │  (regex)   │ │ (ts parse) │ │ (py ast)   │ │ (py ast)   │           │
│  └─────┬──────┘ └─────┬──────┘ └─────┬──────┘ └─────┬──────┘ └─────┬──────┘           │
│                                                                  15 adapters total    │
└────────┼──────────────┼──────────────┼──────────────┼────────────────────────────────┘
         │              │              │              │              │
         ▼              ▼              ▼              ▼              ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Adapters                                       │
│                                                                             │
│    Each adapter: Fetch (git clone) → Extract (parse) → RawMetric            │
│                                                                             │
└─────────────────────────────────┬───────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                            Orchestrator                                     │
│                                                                             │
│    RawMetric → CanonicalMetric → Store (SQLite + FTS5)                      │
│                                                                             │
└─────────────────────────────────┬───────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                             Enricher                                        │
│                                                                             │
│    Cross-reference with OTel Semantic Conventions                           │
│    Match types: exact, prefix, none                                         │
│                                                                             │
└─────────────────────────────────┬───────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              REST API                                       │
│                                                                             │
│    /api/metrics (search)    /api/facets    /health                          │
│                                                                             │
└─────────────────────────────────┬───────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Next.js Frontend                                  │
│                                                                             │
│    Search bar, filters, metric cards, detail view, semconv badges           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

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
make web-build        # Build frontend
make build-all        # Build both

# Extract metrics
make extract          # Extract metrics from otel-collector-contrib
make enrich           # Enrich metrics with semantic convention data

# Test
make test             # Run Go tests
make web-test         # Run frontend tests
make test-all         # Run all tests

# Lint
make lint             # Run golangci-lint
make web-lint         # Run ESLint
make lint-all         # Run all linters

# CI (runs everything)
make ci               # Full CI: lint, test, build (Go + Web)

# Docker
make docker-build     # Build images
make docker-up        # Start services
make docker-down      # Stop services
make docker-logs      # View logs

# Release (see Releasing section below)
git tag web-v0.2.0 && git push origin web-v0.2.0

# Database
make migrate          # Run migrations (via dbmate)
make migrate-status   # Check migration status
```

### Project Structure

```
metric-library/
├── cmd/glossary/          # Main entry point
├── internal/
│   ├── api/               # REST API handlers
│   ├── store/             # SQLite store + migrations
│   ├── domain/            # Domain models
│   ├── adapter/           # Source adapters
│   ├── enricher/          # Semantic convention enrichment
│   ├── fetcher/           # Git fetcher
│   ├── discovery/         # Metadata discovery
│   ├── parser/            # YAML parser
│   └── orchestrator/      # Extraction orchestration
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
- `semconv_match` - Filter by semantic convention match (exact, prefix, none)
- `limit`, `offset` - Pagination

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | API server port |
| `DATABASE_PATH` | ./data/metric-library.db | SQLite database path |
| `CACHE_DIR` | ./.cache | Git repository cache directory |
| `NEXT_PUBLIC_API_URL` | http://localhost:8080 | API URL for frontend |

## Sources

### Available Sources

| Source | Adapter | Extraction | Metrics | Repository |
|--------|---------|------------|---------|------------|
| OpenTelemetry Collector Contrib | `otel-collector-contrib` | YAML metadata | 1261 | [otel-collector-contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib) |
| OpenTelemetry Semantic Conventions | `otel-semconv` | YAML metadata | 349 | [semantic-conventions](https://github.com/open-telemetry/semantic-conventions) |
| OpenTelemetry Python | `otel-python` | Python AST | 30 | [opentelemetry-python-contrib](https://github.com/open-telemetry/opentelemetry-python-contrib) |
| OpenTelemetry Java | `otel-java` | Regex | 50 | [opentelemetry-java-instrumentation](https://github.com/open-telemetry/opentelemetry-java-instrumentation) |
| OpenTelemetry JS | `otel-js` | TS Parse | 35 | [opentelemetry-js-contrib](https://github.com/open-telemetry/opentelemetry-js-contrib) |
| PostgreSQL Exporter | `prometheus-postgres` | Go AST | 120 | [postgres_exporter](https://github.com/prometheus-community/postgres_exporter) |
| Node Exporter | `prometheus-node` | Go AST | 553 | [node_exporter](https://github.com/prometheus/node_exporter) |
| Redis Exporter | `prometheus-redis` | Go AST | 356 | [redis_exporter](https://github.com/oliver006/redis_exporter) |
| MySQL Exporter | `prometheus-mysql` | Go AST | 222 | [mysqld_exporter](https://github.com/prometheus/mysqld_exporter) |
| MongoDB Exporter | `prometheus-mongodb` | Go AST | 8 | [mongodb_exporter](https://github.com/percona/mongodb_exporter) |
| Kafka Exporter | `prometheus-kafka` | Go AST | 16 | [kafka_exporter](https://github.com/danielqsj/kafka_exporter) |
| kube-state-metrics | `kubernetes-ksm` | Go AST | 261 | [kube-state-metrics](https://github.com/kubernetes/kube-state-metrics) |
| cAdvisor | `kubernetes-cadvisor` | Go AST | 107 | [cadvisor](https://github.com/google/cadvisor) |
| OpenLLMetry | `openllmetry` | Python AST | 30 | [openllmetry](https://github.com/traceloop/openllmetry) |
| OpenLIT | `openlit` | Python AST | 21 | [openlit](https://github.com/openlit/openlit) |

**Total: 3,419 metrics**

### Extract Commands

```bash
make extract-otel         # OpenTelemetry Collector Contrib
make extract-semconv      # OpenTelemetry Semantic Conventions
make extract-otel-python  # OpenTelemetry Python
make extract-otel-java    # OpenTelemetry Java
make extract-otel-js      # OpenTelemetry JS
make extract-postgres     # PostgreSQL Exporter
make extract-node         # Node Exporter
make extract-redis        # Redis Exporter
make extract-mysql        # MySQL Exporter
make extract-mongodb      # MongoDB Exporter
make extract-kafka        # Kafka Exporter
make extract-ksm          # kube-state-metrics
make extract-cadvisor     # cAdvisor
make extract-openllmetry  # OpenLLMetry (LLM observability)
make extract-openlit      # OpenLIT (LLM observability)
make extract-all          # All sources
```

### Semantic Conventions Enrichment

After extracting metrics, you can enrich them with OpenTelemetry Semantic Convention compliance data:

```bash
# First, extract semantic conventions (if not already done)
make extract-semconv

# Then run enrichment against all metrics
make enrich
```

The enricher cross-references each metric name against the 349 semantic convention metrics and assigns one of three match types:

| Match Type | Description | UI Badge |
|------------|-------------|----------|
| `exact` | Metric name exactly matches a semantic convention | SemConv (green) |
| `prefix` | Metric name starts with a semantic convention metric | SemConv~ (amber) |
| `none` | No match found | Custom (red) |

The enricher normalizes metric names by converting underscores to dots before matching, so `http_server_request_duration` matches `http.server.request.duration`.

**Example enrichment results:**
- Exact matches: 410 metrics
- Prefix matches: 29 metrics
- No match: 2798 metrics

### Adding a New Source

1. **Create adapter directory**
   ```
   internal/adapter/prometheus/<name>/
   ```

2. **Implement the Adapter interface** (`adapter.go`)
   ```go
   type Adapter interface {
       Name() string
       SourceCategory() domain.SourceCategory
       Confidence() domain.ConfidenceLevel
       ExtractionMethod() domain.ExtractionMethod
       RepoURL() string
       Fetch(ctx context.Context, opts FetchOptions) (*FetchResult, error)
       Extract(ctx context.Context, result *FetchResult) ([]*RawMetric, error)
   }
   ```

3. **Choose extraction method**
   - **YAML metadata**: For sources with `metadata.yaml` files (like otel-collector-contrib)
   - **Go AST**: For Prometheus exporters using `prometheus.NewDesc()` - use the shared parser at `internal/adapter/prometheus/astparser`
   - **Custom AST**: For sources with unique patterns (like redis_exporter's map-based definitions)

4. **Write tests** (`adapter_test.go`)

5. **Register in main.go**
   ```go
   case "prometheus-<name>":
       adp = <name>.NewAdapter(*cacheDir)
   ```

6. **Add Makefile target**
   ```makefile
   extract-<name>: build
       ./bin/$(BINARY_NAME) extract -adapter prometheus-<name>
   ```

### RawMetric Structure

Each adapter extracts metrics into this intermediate format:

```go
type RawMetric struct {
    Name             string             // Metric name (e.g., "redis_connected_clients")
    Description      string             // Help text
    InstrumentType   string             // counter, gauge, histogram, summary
    Attributes       []domain.Attribute // Labels/dimensions
    EnabledByDefault bool
    ComponentType    string             // receiver, processor, exporter, platform
    ComponentName    string             // e.g., "redis", "pg_stat_database"
    SourceLocation   string             // File path in source repo
    Path             string             // Discovery path
}
```

## Releasing

Docker images are published to GitHub Container Registry (GHCR) via GitHub Actions when tags are pushed.

| Component | Image | Tag Pattern |
|-----------|-------|-------------|
| Web Frontend | `ghcr.io/base-14/metric-library-web` | `web-v*` |
| API Backend | `ghcr.io/base-14/metric-library-api` | `api-v*` |

### Releasing a New Version

```bash
# Web frontend
git tag web-v0.2.0 && git push origin web-v0.2.0

# API backend
git tag api-v0.2.0 && git push origin api-v0.2.0
```

Each tag push triggers a workflow that builds and pushes both versioned and `latest` tags.

### Pulling Images

```bash
docker pull ghcr.io/base-14/metric-library-web:latest
docker pull ghcr.io/base-14/metric-library-api:latest
```
