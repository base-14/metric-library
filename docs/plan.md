markdown
# Metric Discovery Platform — Plan & To-Do List

## 1. Purpose

Build a continuously updated, searchable catalog of metrics across:
- OpenTelemetry contrib ecosystems
- Prometheus exporters
- Kubernetes, databases, messaging systems
- Cloud provider metrics
- (Optionally) vendor observability platforms

The system must extract metrics from source-of-truth artifacts, normalize them into a canonical schema, preserve provenance and confidence, and expose them via search and UI.

This is a **metric intelligence layer**, not just a catalog.

---

## 2. Design Principles

- **Source-first**: derive metrics from repos or official docs, not hand-curated lists
- **Incremental coverage**: high-confidence sources first
- **Pluggable adapters**: each source added independently
- **Provenance-aware**: always know where a metric came from and how trustworthy it is
- **Deterministic & reproducible**: same inputs → same outputs
- **Search-oriented**: optimize for discovery, not perfect modeling

---

## 3. High-Level Architecture

```

[ Source Adapter ]
↓
[ Extractor ]
↓
[ Normalizer ]
↓
[ Canonical Metric Store ]
↓
[ Search Index ]
↓
[ API + UI ]

```

### Key abstraction: Source Adapter
Each metric source implements an adapter that defines:
- how to fetch data
- where metrics are defined
- how to extract them
- confidence level of extracted metrics

---

## 4. Canonical Metric Schema (Core Fields)

```json
{
  "metric_name": "...",
  "instrument_type": "...",
  "description": "...",
  "unit": "...",
  "attributes": ["..."],

  "component_type": "receiver | exporter | instrumentation | platform",
  "component_name": "...",

  "source_category": "otel | prometheus | kubernetes | cloud | vendor",
  "source_name": "...",
  "source_location": "...",

  "extraction_method": "metadata | ast | scrape | hybrid",
  "source_confidence": "authoritative | derived | documented | vendor_claimed",

  "repo": "...",
  "path": "...",
  "commit": "...",
  "timestamp": "..."
}
```

This schema intentionally encodes **trust and provenance**.

---

## 5. Metric Source Backlog

### 5.1 Backlog Table

| Category   | Metric Source             | Where to Find Metrics | Primary Artifact       | Extraction Method       | Confidence     |
| ---------- | ------------------------- | --------------------- | ---------------------- | ----------------------- | -------------- |
| OTEL       | Collector Contrib         | GitHub repo           | metadata.yaml, Go code | Repo clone + YAML + AST | Authoritative  |
| OTEL       | Go Contrib                | GitHub repo           | Go code                | Repo clone + AST        | Derived        |
| OTEL       | Java Instrumentation      | GitHub repo           | Java code, docs        | Repo clone + docs       | Derived        |
| OTEL       | Python Contrib            | GitHub repo           | Python code, docs      | Repo clone + AST-lite   | Derived        |
| Prometheus | node_exporter             | GitHub repo           | Go code, README        | Repo clone + AST        | Derived        |
| Prometheus | blackbox_exporter         | GitHub repo           | Go code                | Repo clone + AST        | Derived        |
| Prometheus | postgres_exporter         | GitHub repo           | SQL + Go               | Repo clone + hybrid     | Derived        |
| Kubernetes | kube-state-metrics        | GitHub repo           | Go code                | Repo clone + AST        | Derived        |
| Kubernetes | cAdvisor                  | GitHub repo           | Go code                | Repo clone + AST        | Derived        |
| Databases  | PostgreSQL                | Docs + exporters      | HTML + SQL             | Scrape + repo clone     | Documented     |
| Messaging  | Kafka                     | Docs + exporters      | HTML + Go              | Hybrid                  | Documented     |
| Cloud      | AWS CloudWatch            | AWS docs              | HTML docs              | Scrape                  | Documented     |
| Cloud      | GCP Monitoring            | GCP docs              | HTML docs              | Scrape                  | Documented     |
| Semantics  | OTEL Semantic Conventions | GitHub repo           | YAML / Markdown        | Repo clone + parse      | Authoritative  |
| Vendors    | Datadog                   | Docs                  | HTML docs              | Scrape                  | Vendor-claimed |

---

## 6. Phased Execution Plan

### Phase 0 — Foundations

Goal: prepare the ground so all later work is additive.

* Define canonical schema
* Define adapter interface
* Repo scaffolding
* CI setup
* Basic logging and error handling

---

### Phase 1 — OTEL First-Class Sources + MVP UI

Goal: establish a high-confidence baseline with a usable interface.

Sources:

* OpenTelemetry Collector Contrib
* OpenTelemetry Go Contrib

Extraction:

* metadata.yaml parsing
* Go AST parsing

UI (Next.js + React + Tailwind):

* Search interface with filters
* Metric cards and detail view
* Faceted navigation

Outcome:

* Large, authoritative metric corpus
* Searchable MVP with web UI

---

### Phase 2 — Prometheus Exporters

Goal: cover infra metrics most users already know.

Sources:

* node_exporter
* blackbox_exporter
* postgres_exporter
* redis_exporter
* kafka_exporter

Outcome:

* Infra-focused discovery
* Prometheus + OTEL coexistence

---

### Phase 3 — Kubernetes Metrics

Goal: enable workload-centric exploration.

Sources:

* kube-state-metrics
* cAdvisor
* control-plane metrics

Outcome:

* Kubernetes object-level metrics

---

### Phase 4 — Language Instrumentations

Goal: broaden application-level coverage.

Sources:

* Java
* Python
* JS
* .NET

Strategy:

* Start with docs
* Incrementally add AST extraction

---

### Phase 5 — Databases & Messaging

Goal: domain-specific depth.

Sources:

* PostgreSQL, MySQL, MongoDB, Redis
* Kafka, RabbitMQ, Pulsar, NATS

Outcome:

* Strong SRE-aligned browsing

---

### Phase 6 — Cloud Providers

Goal: align with existing mental models.

Sources:

* AWS CloudWatch
* GCP Monitoring
* Azure Monitor

Outcome:

* Cloud-native metric taxonomy

---

### Phase 7 — Semantic Conventions

Goal: enrichment and validation.

Sources:

* OTEL semantic conventions
* Stability specs

Outcome:

* Attribute validation
* Future alert/SLO synthesis

---


## 7. To-Do List (Execution-Ready)

### Foundation ✅

* [x] Write canonical schema (Go domain models)
* [x] Document source adapter interface
* [x] Set up repo structure
* [x] Add CI with GitHub Actions (lint, test, build)

### OTEL Collector Contrib ✅

* [x] Repo fetcher (go-git, pinned by commit)
* [x] metadata.yaml discovery
* [x] YAML parser with validation
* [x] Metric extraction
* [x] Normalization
* [x] 1261 metrics extracted

### OTEL Go Contrib

* [ ] Go file discovery
* [ ] AST walker
* [ ] Instrument detection
* [ ] Option parsing (description, unit)
* [ ] Merge with metadata metrics

### Storage & Search ✅

* [x] Canonical metric store (SQLite)
* [x] Deterministic output ordering
* [x] Full-text search index (FTS5)
* [x] Batch indexing

### Search & Filtering ✅

* [x] Free-text search across metric name, description, component
* [x] Filter by instrument type (counter, gauge, histogram, updowncounter, summary)
* [x] Filter by component type (receiver, exporter, processor, instrumentation, platform)
* [x] Filter by component name
* [x] Filter by source category (otel, prometheus, kubernetes, cloud, vendor)
* [x] Filter by source name
* [x] Filter by confidence level (authoritative, derived, documented, vendor_claimed)
* [ ] Filter by unit
* [ ] Filter by attribute name
* [x] Combined filters (AND logic)
* [x] Pagination support
* [x] Facet counts for filter options

### Prometheus Exporters

* [ ] Exporter adapter
* [ ] AST extraction
* [ ] README enrichment
* [ ] Confidence tagging

### Kubernetes Metrics

* [ ] kube-state-metrics adapter
* [ ] cAdvisor adapter
* [ ] Metric categorization by object type

### Language Instrumentations

* [ ] Java doc-based extraction
* [ ] Python doc-based extraction
* [ ] AST extraction (incremental)

### Databases & Messaging

* [ ] Exporter-based extraction
* [ ] Doc scraping
* [ ] Hybrid merge logic

### Cloud Providers

* [ ] Doc scraper framework
* [ ] CloudWatch metrics ingestion
* [ ] GCP metrics ingestion
* [ ] Azure metrics ingestion

### Semantic Conventions

* [ ] Parse conventions
* [ ] Attribute validation logic
* [ ] Metric enrichment hooks

### Automation

* [ ] Nightly refresh job
* [ ] Change detection
* [ ] Failure alerting

### UI (Next.js + React + Tailwind)

* [x] Next.js project setup with TypeScript
* [x] Tailwind CSS configuration
* [x] API client for backend integration
* [x] Layout and navigation
* [ ] Search bar with debounced input (search as you type)
* [x] Filter sidebar (metric type, component type, source category)
* [ ] Reorder filters: Component Name first, Instrument Type second
* [x] Facet counts in filter options
* [ ] Active filters display with clear buttons
* [x] Metric card component with type badges
* [x] Metric list with responsive grid layout
* [x] Metric detail panel (slide-out or dedicated page)
* [x] Attribute display with enum values
* [ ] Copy-to-clipboard for metric names
* [x] Pagination
* [x] Empty state and loading states
* [ ] GitHub source link for each metric
* [ ] Mobile responsive design
* [ ] Dark mode support (optional)

### DevOps ✅

* [x] Docker multi-stage build (Go backend)
* [x] Docker standalone build (Next.js frontend)
* [x] docker-compose for local development
* [x] make extract command for CLI extraction
* [ ] Helm chart for Kubernetes deployment

---

## 8. Guardrails

* One adapter = one milestone
* No cross-ecosystem mapping until OTEL + Prometheus are stable
* Always preserve raw extraction output
* Never drop metrics silently—flag them

---

## 9. End State

When complete, this system should answer:

* What metrics exist for a given system?
* Where do they come from?
* How trustworthy are they?
* What am I missing?
* How does this map across ecosystems?

### Search Capabilities

Users can discover metrics by:

* **Free-text search**: Find metrics by name, description, or component
* **Instrument type**: "Show me all histograms" or "Find counters for HTTP"
* **Component**: "What metrics does the kafka receiver emit?"
* **Source category**: "Show all Prometheus metrics" or "List OTEL metrics"
* **Confidence**: "Show only authoritative metrics"
* **Attributes**: "Find metrics with the `http.method` attribute"
* **Combined**: "Gauges from kubernetes with high confidence"

This document is the source of truth for execution.
