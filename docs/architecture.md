# Metric Discovery Platform - Architecture Document

## 1. Overview

The Metric Discovery Platform is a continuously updated, searchable catalog of metrics across observability ecosystems. It extracts metrics from source-of-truth artifacts, normalizes them into a canonical schema, preserves provenance and confidence levels, and exposes them via search and API.

### Goals

- **Discovery**: Answer "What metrics exist for system X?"
- **Provenance**: Track where each metric came from
- **Trust**: Encode confidence levels for each metric
- **Search**: Optimize for metric discovery across ecosystems

---

## 2. System Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                        Metric Discovery Platform                     │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐               │
│  │    OTEL      │  │  Prometheus  │  │  Kubernetes  │  ...more      │
│  │   Adapter    │  │   Adapter    │  │   Adapter    │               │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘               │
│         │                 │                 │                        │
│         └─────────────────┼─────────────────┘                        │
│                           ▼                                          │
│                  ┌─────────────────┐                                 │
│                  │    Extractor    │                                 │
│                  │    Pipeline     │                                 │
│                  └────────┬────────┘                                 │
│                           │                                          │
│                           ▼                                          │
│                  ┌─────────────────┐                                 │
│                  │   Normalizer    │                                 │
│                  └────────┬────────┘                                 │
│                           │                                          │
│                           ▼                                          │
│         ┌─────────────────────────────────────┐                      │
│         │       Canonical Metric Store        │                      │
│         │            (SQLite)                 │                      │
│         └─────────────────┬───────────────────┘                      │
│                           │                                          │
│                           ▼                                          │
│                  ┌─────────────────┐                                 │
│                  │   Search Index  │                                 │
│                  │    (FTS5)       │                                 │
│                  └────────┬────────┘                                 │
│                           │                                          │
│                           ▼                                          │
│                  ┌─────────────────┐                                 │
│                  │    REST API     │                                 │
│                  └─────────────────┘                                 │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 3. Component Architecture

### 3.1 Source Adapters

Each metric source implements an adapter that handles source-specific logic.

```go
// Adapter defines the interface for metric source adapters
type Adapter interface {
    // Name returns the unique identifier for this adapter
    Name() string

    // Fetch retrieves raw data from the source
    Fetch(ctx context.Context, opts FetchOptions) (*FetchResult, error)

    // Extract parses raw data and extracts metrics
    Extract(ctx context.Context, data *FetchResult) ([]*RawMetric, error)

    // Confidence returns the confidence level for this source
    Confidence() ConfidenceLevel

    // SourceCategory returns the category (otel, prometheus, kubernetes, etc.)
    SourceCategory() SourceCategory
}

type FetchOptions struct {
    Commit    string    // Pin to specific commit (empty = latest)
    CacheDir  string    // Local cache directory
    Force     bool      // Force re-fetch even if cached
}

type FetchResult struct {
    RepoPath  string
    Commit    string
    Timestamp time.Time
    Files     []string  // Discovered files relevant to extraction
}
```

**Adapter Types:**

| Adapter | Source | Extraction Method |
|---------|--------|-------------------|
| `otel-collector-contrib` | GitHub repo | metadata.yaml + Go AST |
| `otel-go-contrib` | GitHub repo | Go AST |
| `prometheus-exporter` | GitHub repos | Go AST + README |
| `kube-state-metrics` | GitHub repo | Go AST |
| `cloud-provider` | Documentation | HTML scraping |

### 3.2 Extractor Pipeline

The extractor pipeline processes source data through multiple stages:

```
┌─────────────┐    ┌──────────────┐    ┌────────────────┐    ┌────────────┐
│  Fetcher    │───▶│  Discovery   │───▶│   Extraction   │───▶│ Validation │
│             │    │              │    │                │    │            │
│ Git clone   │    │ Find metric  │    │ Parse YAML/    │    │ Schema     │
│ or download │    │ definition   │    │ AST/HTML       │    │ validation │
│             │    │ files        │    │                │    │            │
└─────────────┘    └──────────────┘    └────────────────┘    └────────────┘
```

**Extraction Methods:**

- **YAML Parsing**: For metadata.yaml files (OTEL Collector Contrib)
- **Go AST Parsing**: For Go source code metric definitions
- **HTML Scraping**: For documentation-based sources
- **Hybrid**: Combination of multiple methods

### 3.3 Normalizer

Transforms raw extracted metrics into the canonical schema:

```go
type Normalizer interface {
    Normalize(ctx context.Context, raw *RawMetric) (*CanonicalMetric, error)
}
```

Responsibilities:
- Map source-specific fields to canonical fields
- Standardize unit formats (e.g., "milliseconds" → "ms")
- Validate instrument types
- Deduplicate attributes
- Generate deterministic metric IDs

### 3.4 Canonical Metric Store

SQLite database with the following schema:

```sql
-- Core metrics table
CREATE TABLE metrics (
    id              TEXT PRIMARY KEY,  -- Deterministic hash
    metric_name     TEXT NOT NULL,
    instrument_type TEXT NOT NULL,
    description     TEXT,
    unit            TEXT,

    -- Component info
    component_type  TEXT NOT NULL,
    component_name  TEXT NOT NULL,

    -- Source info
    source_category TEXT NOT NULL,
    source_name     TEXT NOT NULL,
    source_location TEXT,

    -- Provenance
    extraction_method   TEXT NOT NULL,
    source_confidence   TEXT NOT NULL,
    repo                TEXT,
    path                TEXT,
    commit              TEXT,
    extracted_at        TIMESTAMP NOT NULL,

    -- Metadata
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Metric attributes (many-to-many)
CREATE TABLE metric_attributes (
    metric_id       TEXT NOT NULL REFERENCES metrics(id),
    attribute_name  TEXT NOT NULL,
    attribute_type  TEXT,
    description     TEXT,
    required        BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (metric_id, attribute_name)
);

-- Full-text search index
CREATE VIRTUAL TABLE metrics_fts USING fts5(
    metric_name,
    description,
    component_name,
    source_name,
    content=metrics,
    content_rowid=rowid
);

-- Extraction runs for tracking
CREATE TABLE extraction_runs (
    id              TEXT PRIMARY KEY,
    adapter_name    TEXT NOT NULL,
    commit          TEXT,
    started_at      TIMESTAMP NOT NULL,
    completed_at    TIMESTAMP,
    metrics_count   INTEGER,
    status          TEXT NOT NULL,
    error_message   TEXT
);
```

### 3.5 Search Index

Full-text search using SQLite FTS5:

```go
type SearchService interface {
    Search(ctx context.Context, query SearchQuery) (*SearchResult, error)
    Reindex(ctx context.Context) error
}

type SearchQuery struct {
    Text            string            // Free-text search
    SourceCategory  []SourceCategory  // Filter by category
    ComponentType   []ComponentType   // Filter by component type
    Confidence      []ConfidenceLevel // Filter by confidence
    Limit           int
    Offset          int
}

type SearchResult struct {
    Metrics    []*CanonicalMetric
    Total      int
    Took       time.Duration
}
```

### 3.6 REST API

HTTP API for metric discovery:

```
GET  /api/v1/metrics                    # List/search metrics
GET  /api/v1/metrics/{id}               # Get metric by ID
GET  /api/v1/sources                    # List available sources
GET  /api/v1/sources/{name}/metrics     # Metrics by source
GET  /api/v1/components/{name}/metrics  # Metrics by component
GET  /api/v1/health                     # Health check
POST /api/v1/refresh                    # Trigger refresh (admin)
```

---

## 4. Canonical Schema

### 4.1 Domain Types

```go
type CanonicalMetric struct {
    ID                string            `json:"id"`
    MetricName        string            `json:"metric_name"`
    InstrumentType    InstrumentType    `json:"instrument_type"`
    Description       string            `json:"description"`
    Unit              string            `json:"unit"`
    Attributes        []Attribute       `json:"attributes"`

    ComponentType     ComponentType     `json:"component_type"`
    ComponentName     string            `json:"component_name"`

    SourceCategory    SourceCategory    `json:"source_category"`
    SourceName        string            `json:"source_name"`
    SourceLocation    string            `json:"source_location"`

    ExtractionMethod  ExtractionMethod  `json:"extraction_method"`
    SourceConfidence  ConfidenceLevel   `json:"source_confidence"`

    Repo              string            `json:"repo"`
    Path              string            `json:"path"`
    Commit            string            `json:"commit"`
    ExtractedAt       time.Time         `json:"extracted_at"`
}

type Attribute struct {
    Name        string `json:"name"`
    Type        string `json:"type"`
    Description string `json:"description"`
    Required    bool   `json:"required"`
}
```

### 4.2 Enumerations

```go
type InstrumentType string
const (
    InstrumentCounter         InstrumentType = "counter"
    InstrumentUpDownCounter   InstrumentType = "updowncounter"
    InstrumentGauge           InstrumentType = "gauge"
    InstrumentHistogram       InstrumentType = "histogram"
    InstrumentSummary         InstrumentType = "summary"  // Prometheus
)

type ComponentType string
const (
    ComponentReceiver        ComponentType = "receiver"
    ComponentExporter        ComponentType = "exporter"
    ComponentProcessor       ComponentType = "processor"
    ComponentInstrumentation ComponentType = "instrumentation"
    ComponentPlatform        ComponentType = "platform"
)

type SourceCategory string
const (
    SourceOTEL       SourceCategory = "otel"
    SourcePrometheus SourceCategory = "prometheus"
    SourceKubernetes SourceCategory = "kubernetes"
    SourceCloud      SourceCategory = "cloud"
    SourceVendor     SourceCategory = "vendor"
)

type ExtractionMethod string
const (
    ExtractionMetadata ExtractionMethod = "metadata"
    ExtractionAST      ExtractionMethod = "ast"
    ExtractionScrape   ExtractionMethod = "scrape"
    ExtractionHybrid   ExtractionMethod = "hybrid"
)

type ConfidenceLevel string
const (
    ConfidenceAuthoritative ConfidenceLevel = "authoritative"
    ConfidenceDerived       ConfidenceLevel = "derived"
    ConfidenceDocumented    ConfidenceLevel = "documented"
    ConfidenceVendorClaimed ConfidenceLevel = "vendor_claimed"
)
```

---

## 5. Package Structure

```
metric-library/
├── cmd/
│   └── glossary/
│       └── main.go              # Application entry point
├── internal/
│   ├── adapter/                 # Source adapters
│   │   ├── adapter.go           # Interface definitions
│   │   ├── otel/
│   │   │   ├── collector.go     # OTEL Collector Contrib adapter
│   │   │   └── gocontrib.go     # OTEL Go Contrib adapter
│   │   ├── prometheus/
│   │   │   └── exporter.go      # Prometheus exporter adapter
│   │   └── kubernetes/
│   │       └── ksm.go           # kube-state-metrics adapter
│   ├── extractor/               # Extraction logic
│   │   ├── yaml.go              # YAML parser
│   │   ├── goast.go             # Go AST parser
│   │   └── html.go              # HTML scraper
│   ├── normalizer/              # Normalization logic
│   │   └── normalizer.go
│   ├── store/                   # Storage layer
│   │   ├── store.go             # Interface
│   │   └── sqlite.go            # SQLite implementation
│   ├── search/                  # Search service
│   │   └── search.go
│   ├── api/                     # HTTP API
│   │   ├── server.go
│   │   └── handlers.go
│   └── domain/                  # Domain types
│       ├── metric.go
│       └── types.go
├── db/
│   └── migrations/              # dbmate migrations
│       ├── 001_create_metrics.sql
│       └── 002_create_fts.sql
├── web/                         # Next.js frontend
│   ├── src/
│   │   ├── app/                 # App router pages
│   │   │   ├── layout.tsx
│   │   │   ├── page.tsx         # Home/search page
│   │   │   └── metrics/
│   │   │       └── [id]/
│   │   │           └── page.tsx # Metric detail page
│   │   ├── components/          # React components
│   │   │   ├── SearchBar.tsx
│   │   │   ├── FilterSidebar.tsx
│   │   │   ├── MetricCard.tsx
│   │   │   ├── MetricDetail.tsx
│   │   │   └── TypeBadge.tsx
│   │   ├── lib/                 # Utilities
│   │   │   ├── api.ts           # API client
│   │   │   └── types.ts         # TypeScript types
│   │   └── hooks/               # Custom hooks
│   │       └── useMetrics.ts
│   ├── public/
│   ├── tailwind.config.ts
│   ├── next.config.js
│   ├── package.json
│   └── tsconfig.json
├── docs/
│   ├── architecture.md          # This document
│   ├── plan.md
│   └── ui-prototype.jsx         # UI prototype
├── Makefile
├── go.mod
├── go.sum
└── CLAUDE.md
```

---

## 6. Data Flow

### 6.1 Extraction Flow

```
1. Scheduler triggers adapter refresh
   │
   ▼
2. Adapter.Fetch() clones/updates repo
   │
   ▼
3. Adapter discovers metric definition files
   │
   ▼
4. Extractor parses files (YAML/AST/HTML)
   │
   ▼
5. Raw metrics validated against schema
   │
   ▼
6. Normalizer transforms to canonical format
   │
   ▼
7. Store upserts metrics (deduplication by ID)
   │
   ▼
8. Search index updated
```

### 6.2 Query Flow

```
1. API receives search request
   │
   ▼
2. SearchService queries FTS index
   │
   ▼
3. Results filtered by facets
   │
   ▼
4. Metrics loaded from store
   │
   ▼
5. Response returned to client
```

---

## 7. Technology Stack

### Backend

| Component | Technology | Rationale |
|-----------|-----------|-----------|
| Language | Go 1.22+ | Performance, AST parsing, ecosystem fit |
| Database | SQLite + FTS5 | Simple deployment, full-text search |
| DI | uber-go/dig | Clean dependency management |
| HTTP | net/http + chi | Lightweight, standard library |
| Git | go-git | Pure Go git implementation |
| YAML | gopkg.in/yaml.v3 | Standard YAML parsing |
| Migrations | dbmate | SQL migration management |
| Linting | golangci-lint | Code quality |
| Testing | testing + testify | Standard + assertions |

### Frontend

| Component | Technology | Rationale |
|-----------|-----------|-----------|
| Framework | Next.js 14+ | SSR, routing, API routes |
| UI Library | React 18+ | Component-based UI |
| Styling | Tailwind CSS | Utility-first, rapid development |
| Language | TypeScript | Type safety |
| State | React hooks | Simple state management |
| Fetching | SWR or TanStack Query | Caching, revalidation |

---

## 8. Key Design Decisions

### 8.1 SQLite over PostgreSQL

- **Simplicity**: No external database required
- **Portability**: Single file database
- **FTS5**: Built-in full-text search
- **Sufficient Scale**: Expected <1M metrics

### 8.2 Deterministic Metric IDs

Metrics are identified by a hash of:
- `source_category`
- `source_name`
- `component_name`
- `metric_name`

This ensures:
- Same input → same ID (reproducible)
- Idempotent upserts
- Stable references

### 8.3 Confidence Levels

Trust is encoded, not assumed:

| Level | Meaning | Example |
|-------|---------|---------|
| `authoritative` | From official metadata | metadata.yaml |
| `derived` | Extracted from code | Go AST parsing |
| `documented` | From official docs | AWS CloudWatch docs |
| `vendor_claimed` | From vendor docs | Datadog integrations |

### 8.4 Adapter Isolation

Each adapter:
- Runs independently
- Has its own extraction logic
- Can fail without affecting others
- Tracks its own extraction runs

---

## 9. Error Handling

### 9.1 Extraction Errors

- **Log all errors** with context
- **Never drop metrics silently** - flag them
- **Continue on partial failure** - extract what's possible
- **Track error rates** per adapter

### 9.2 Error Categories

```go
type ExtractionError struct {
    Adapter   string
    File      string
    Line      int
    Message   string
    Severity  ErrorSeverity  // warning, error, fatal
}

type ErrorSeverity string
const (
    SeverityWarning ErrorSeverity = "warning"  // Metric extracted with issues
    SeverityError   ErrorSeverity = "error"    // Metric skipped
    SeverityFatal   ErrorSeverity = "fatal"    // Adapter failed
)
```

---

## 10. Testing Strategy

### 10.1 Unit Tests

- Domain type validation
- Normalizer transformations
- Individual extractor functions

### 10.2 Integration Tests

- Adapter extraction from fixture repos
- Store operations
- Search queries

### 10.3 Snapshot Tests

- Extraction output stability
- Ensures deterministic results
- Catches unintended changes

---

## 11. Deployment

### 11.1 Build

```bash
make build    # Build binary
make test     # Run tests
make lint     # Run linter
make migrate  # Run migrations
```

### 11.2 Configuration

Environment variables:
```
DATABASE_PATH=./data/metric-library.db
CACHE_DIR=./cache
PORT=8080
LOG_LEVEL=info
```

### 11.3 Automation

- **Nightly refresh**: Cron job triggers full extraction
- **Change detection**: Only process changed files
- **Failure alerting**: Notify on adapter failures

---

## 12. Future Considerations

### 12.1 Potential Enhancements

- Web UI for browsing
- Metric relationship mapping
- Alert/SLO template generation
- Semantic convention validation

### 12.2 Scaling

If metrics exceed SQLite capacity:
- PostgreSQL with pg_trgm for search
- Redis for caching

---

## 13. References

- [OpenTelemetry Collector Contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib)
- [OTEL Semantic Conventions](https://github.com/open-telemetry/semantic-conventions)
- [Prometheus Exporters](https://prometheus.io/docs/instrumenting/exporters/)
- [kube-state-metrics](https://github.com/kubernetes/kube-state-metrics)
