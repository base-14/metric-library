package orchestrator

import (
	"context"
	"testing"
	"time"

	"github.com/base14/otel-glossary/internal/adapter"
	"github.com/base14/otel-glossary/internal/domain"
	"github.com/base14/otel-glossary/internal/store"
)

type mockAdapter struct {
	name           string
	sourceCategory domain.SourceCategory
	confidence     domain.ConfidenceLevel
	extraction     domain.ExtractionMethod
	repoURL        string
	fetchResult    *adapter.FetchResult
	rawMetrics     []*adapter.RawMetric
	fetchErr       error
	extractErr     error
}

func (m *mockAdapter) Name() string                              { return m.name }
func (m *mockAdapter) SourceCategory() domain.SourceCategory     { return m.sourceCategory }
func (m *mockAdapter) Confidence() domain.ConfidenceLevel        { return m.confidence }
func (m *mockAdapter) ExtractionMethod() domain.ExtractionMethod { return m.extraction }
func (m *mockAdapter) RepoURL() string                           { return m.repoURL }

func (m *mockAdapter) Fetch(ctx context.Context, opts adapter.FetchOptions) (*adapter.FetchResult, error) {
	if m.fetchErr != nil {
		return nil, m.fetchErr
	}
	return m.fetchResult, nil
}

func (m *mockAdapter) Extract(ctx context.Context, result *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	if m.extractErr != nil {
		return nil, m.extractErr
	}
	return m.rawMetrics, nil
}

type mockStore struct {
	metrics []*domain.CanonicalMetric
	runs    []*store.ExtractionRun
}

func (m *mockStore) UpsertMetric(ctx context.Context, metric *domain.CanonicalMetric) error {
	m.metrics = append(m.metrics, metric)
	return nil
}

func (m *mockStore) UpsertMetrics(ctx context.Context, metrics []*domain.CanonicalMetric) error {
	m.metrics = append(m.metrics, metrics...)
	return nil
}

func (m *mockStore) GetMetric(ctx context.Context, id string) (*domain.CanonicalMetric, error) {
	return nil, nil
}

func (m *mockStore) DeleteMetric(ctx context.Context, id string) error {
	return nil
}

func (m *mockStore) DeleteMetricsBySource(ctx context.Context, sourceName string) error {
	return nil
}

func (m *mockStore) Search(ctx context.Context, query store.SearchQuery) (*store.SearchResult, error) {
	return nil, nil
}

func (m *mockStore) GetFacetCounts(ctx context.Context) (*store.FacetCounts, error) {
	return nil, nil
}

func (m *mockStore) GetFilteredFacetCounts(ctx context.Context, query store.FacetQuery) (*store.FacetCounts, error) {
	return nil, nil
}

func (m *mockStore) CreateExtractionRun(ctx context.Context, run *store.ExtractionRun) error {
	m.runs = append(m.runs, run)
	return nil
}

func (m *mockStore) UpdateExtractionRun(ctx context.Context, run *store.ExtractionRun) error {
	for i, r := range m.runs {
		if r.ID == run.ID {
			m.runs[i] = run
			return nil
		}
	}
	return nil
}

func (m *mockStore) GetExtractionRun(ctx context.Context, id string) (*store.ExtractionRun, error) {
	return nil, nil
}

func (m *mockStore) GetLatestExtractionRun(ctx context.Context, adapterName string) (*store.ExtractionRun, error) {
	return nil, nil
}

func (m *mockStore) GetSemconvMetrics(ctx context.Context) ([]*domain.CanonicalMetric, error) {
	return nil, nil
}

func (m *mockStore) Close() error {
	return nil
}

func TestExtractor_Run(t *testing.T) {
	mockAdp := &mockAdapter{
		name:           "test-adapter",
		sourceCategory: domain.SourceOTEL,
		confidence:     domain.ConfidenceAuthoritative,
		extraction:     domain.ExtractionMetadata,
		repoURL:        "https://github.com/test/repo",
		fetchResult: &adapter.FetchResult{
			RepoPath:  "/tmp/repo",
			Commit:    "abc123",
			Timestamp: time.Now(),
		},
		rawMetrics: []*adapter.RawMetric{
			{
				Name:             "http_requests_total",
				InstrumentType:   "counter",
				Description:      "Total HTTP requests",
				Unit:             "1",
				ComponentType:    "receiver",
				ComponentName:    "httpreceiver",
				EnabledByDefault: true,
			},
			{
				Name:             "http_request_duration",
				InstrumentType:   "histogram",
				Description:      "HTTP request duration",
				Unit:             "ms",
				ComponentType:    "receiver",
				ComponentName:    "httpreceiver",
				EnabledByDefault: true,
			},
		},
	}

	mockSt := &mockStore{}

	ext := NewExtractor(mockAdp, mockSt)
	result, err := ext.Run(context.Background(), Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.MetricsExtracted != 2 {
		t.Errorf("expected 2 metrics extracted, got %d", result.MetricsExtracted)
	}

	if result.MetricsStored != 2 {
		t.Errorf("expected 2 metrics stored, got %d", result.MetricsStored)
	}

	if len(mockSt.metrics) != 2 {
		t.Fatalf("expected 2 metrics in store, got %d", len(mockSt.metrics))
	}

	metric := mockSt.metrics[0]
	if metric.MetricName != "http_requests_total" {
		t.Errorf("expected metric name 'http_requests_total', got '%s'", metric.MetricName)
	}
	if metric.InstrumentType != domain.InstrumentCounter {
		t.Errorf("expected instrument type 'counter', got '%s'", metric.InstrumentType)
	}
	if metric.SourceCategory != domain.SourceOTEL {
		t.Errorf("expected source category 'otel', got '%s'", metric.SourceCategory)
	}
	if metric.Commit != "abc123" {
		t.Errorf("expected commit 'abc123', got '%s'", metric.Commit)
	}
	if metric.ID == "" {
		t.Error("expected metric to have an ID")
	}
}

func TestExtractor_TracksExtractionRun(t *testing.T) {
	mockAdp := &mockAdapter{
		name:           "test-adapter",
		sourceCategory: domain.SourceOTEL,
		confidence:     domain.ConfidenceAuthoritative,
		extraction:     domain.ExtractionMetadata,
		repoURL:        "https://github.com/test/repo",
		fetchResult: &adapter.FetchResult{
			RepoPath:  "/tmp/repo",
			Commit:    "def456",
			Timestamp: time.Now(),
		},
		rawMetrics: []*adapter.RawMetric{
			{
				Name:           "test_metric",
				InstrumentType: "gauge",
				ComponentType:  "processor",
				ComponentName:  "testprocessor",
			},
		},
	}

	mockSt := &mockStore{}

	ext := NewExtractor(mockAdp, mockSt)
	_, err := ext.Run(context.Background(), Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(mockSt.runs) != 1 {
		t.Fatalf("expected 1 extraction run, got %d", len(mockSt.runs))
	}

	run := mockSt.runs[0]
	if run.AdapterName != "test-adapter" {
		t.Errorf("expected adapter name 'test-adapter', got '%s'", run.AdapterName)
	}
	if run.Commit != "def456" {
		t.Errorf("expected commit 'def456', got '%s'", run.Commit)
	}
	if run.Status != "completed" {
		t.Errorf("expected status 'completed', got '%s'", run.Status)
	}
	if run.MetricsCount != 1 {
		t.Errorf("expected metrics count 1, got %d", run.MetricsCount)
	}
}
