package store

import (
	"context"
	"time"

	"github.com/base14/otel-glossary/internal/domain"
)

type SearchQuery struct {
	Text             string
	InstrumentTypes  []domain.InstrumentType
	ComponentTypes   []domain.ComponentType
	ComponentNames   []string
	SourceCategories []domain.SourceCategory
	SourceNames      []string
	ConfidenceLevels []domain.ConfidenceLevel
	SemconvMatches   []domain.SemconvMatch
	Units            []string
	AttributeNames   []string
	Limit            int
	Offset           int
}

type SearchResult struct {
	Metrics []*domain.CanonicalMetric
	Total   int
	Took    time.Duration
}

type FacetCounts struct {
	InstrumentTypes  map[domain.InstrumentType]int
	ComponentTypes   map[domain.ComponentType]int
	ComponentNames   map[string]int
	SourceCategories map[domain.SourceCategory]int
	SourceNames      map[string]int
	ConfidenceLevels map[domain.ConfidenceLevel]int
	SemconvMatches   map[domain.SemconvMatch]int
	Units            map[string]int
}

type FacetQuery struct {
	SourceName string
}

type ExtractionRun struct {
	ID           string
	AdapterName  string
	Commit       string
	StartedAt    time.Time
	CompletedAt  *time.Time
	MetricsCount int
	Status       string
	ErrorMessage string
}

type Store interface {
	// Metrics
	UpsertMetric(ctx context.Context, metric *domain.CanonicalMetric) error
	UpsertMetrics(ctx context.Context, metrics []*domain.CanonicalMetric) error
	GetMetric(ctx context.Context, id string) (*domain.CanonicalMetric, error)
	DeleteMetric(ctx context.Context, id string) error
	DeleteMetricsBySource(ctx context.Context, sourceName string) error

	// Search
	Search(ctx context.Context, query SearchQuery) (*SearchResult, error)
	GetFacetCounts(ctx context.Context) (*FacetCounts, error)
	GetFilteredFacetCounts(ctx context.Context, query FacetQuery) (*FacetCounts, error)

	// Semconv
	GetSemconvMetrics(ctx context.Context) ([]*domain.CanonicalMetric, error)

	// Extraction runs
	CreateExtractionRun(ctx context.Context, run *ExtractionRun) error
	UpdateExtractionRun(ctx context.Context, run *ExtractionRun) error
	GetExtractionRun(ctx context.Context, id string) (*ExtractionRun, error)
	GetLatestExtractionRun(ctx context.Context, adapterName string) (*ExtractionRun, error)

	// Lifecycle
	Close() error
}
