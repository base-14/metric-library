package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/base14/otel-glossary/internal/domain"
	"github.com/base14/otel-glossary/internal/store"
)

type mockStore struct {
	metrics []*domain.CanonicalMetric
}

func (m *mockStore) UpsertMetric(ctx context.Context, metric *domain.CanonicalMetric) error {
	return nil
}

func (m *mockStore) UpsertMetrics(ctx context.Context, metrics []*domain.CanonicalMetric) error {
	return nil
}

func (m *mockStore) GetMetric(ctx context.Context, id string) (*domain.CanonicalMetric, error) {
	for _, metric := range m.metrics {
		if metric.ID == id {
			return metric, nil
		}
	}
	return nil, nil
}

func (m *mockStore) DeleteMetric(ctx context.Context, id string) error {
	return nil
}

func (m *mockStore) DeleteMetricsBySource(ctx context.Context, sourceName string) error {
	return nil
}

func (m *mockStore) Search(ctx context.Context, query store.SearchQuery) (*store.SearchResult, error) {
	var results []*domain.CanonicalMetric
	for _, metric := range m.metrics {
		if len(query.InstrumentTypes) > 0 {
			found := false
			for _, t := range query.InstrumentTypes {
				if metric.InstrumentType == t {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if len(query.ComponentTypes) > 0 {
			found := false
			for _, t := range query.ComponentTypes {
				if metric.ComponentType == t {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		results = append(results, metric)
	}

	total := len(results)
	offset := query.Offset
	if offset > len(results) {
		offset = len(results)
	}
	limit := query.Limit
	if limit == 0 {
		limit = 20
	}
	end := offset + limit
	if end > len(results) {
		end = len(results)
	}

	return &store.SearchResult{
		Metrics: results[offset:end],
		Total:   total,
	}, nil
}

func (m *mockStore) GetFacetCounts(ctx context.Context) (*store.FacetCounts, error) {
	return &store.FacetCounts{
		InstrumentTypes:  map[domain.InstrumentType]int{domain.InstrumentCounter: 5, domain.InstrumentGauge: 3},
		ComponentTypes:   map[domain.ComponentType]int{domain.ComponentReceiver: 4, domain.ComponentProcessor: 2},
		SourceCategories: map[domain.SourceCategory]int{domain.SourceOTEL: 8},
		ConfidenceLevels: map[domain.ConfidenceLevel]int{domain.ConfidenceAuthoritative: 8},
		SemconvMatches:   map[domain.SemconvMatch]int{domain.SemconvMatchExact: 2, domain.SemconvMatchNone: 6},
	}, nil
}

func (m *mockStore) GetFilteredFacetCounts(ctx context.Context, query store.FacetQuery) (*store.FacetCounts, error) {
	return m.GetFacetCounts(ctx)
}

func (m *mockStore) CreateExtractionRun(ctx context.Context, run *store.ExtractionRun) error {
	return nil
}

func (m *mockStore) UpdateExtractionRun(ctx context.Context, run *store.ExtractionRun) error {
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

func newTestMetrics() []*domain.CanonicalMetric {
	return []*domain.CanonicalMetric{
		{
			ID:               "metric1",
			MetricName:       "mysql.buffer_pool.pages",
			InstrumentType:   domain.InstrumentCounter,
			Description:      "Number of pages in buffer pool",
			Unit:             "{pages}",
			ComponentType:    domain.ComponentReceiver,
			ComponentName:    "mysql",
			SourceCategory:   domain.SourceOTEL,
			SourceName:       "otel-collector-contrib",
			ExtractionMethod: domain.ExtractionMetadata,
			SourceConfidence: domain.ConfidenceAuthoritative,
			ExtractedAt:      time.Now(),
		},
		{
			ID:               "metric2",
			MetricName:       "system.cpu.utilization",
			InstrumentType:   domain.InstrumentGauge,
			Description:      "CPU utilization percentage",
			Unit:             "1",
			ComponentType:    domain.ComponentReceiver,
			ComponentName:    "hostmetrics",
			SourceCategory:   domain.SourceOTEL,
			SourceName:       "otel-collector-contrib",
			ExtractionMethod: domain.ExtractionMetadata,
			SourceConfidence: domain.ConfidenceAuthoritative,
			ExtractedAt:      time.Now(),
		},
	}
}

func TestAPI_SearchMetrics(t *testing.T) {
	ms := &mockStore{metrics: newTestMetrics()}
	handler := NewHandler(ms)

	req := httptest.NewRequest(http.MethodGet, "/api/metrics", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp SearchResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Metrics) != 2 {
		t.Errorf("expected 2 metrics, got %d", len(resp.Metrics))
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Total)
	}
}

func TestAPI_SearchMetrics_WithFilters(t *testing.T) {
	ms := &mockStore{metrics: newTestMetrics()}
	handler := NewHandler(ms)

	req := httptest.NewRequest(http.MethodGet, "/api/metrics?instrument_type=counter", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp SearchResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Metrics) != 1 {
		t.Errorf("expected 1 metric, got %d", len(resp.Metrics))
	}
}

func TestAPI_SearchMetrics_Pagination(t *testing.T) {
	ms := &mockStore{metrics: newTestMetrics()}
	handler := NewHandler(ms)

	req := httptest.NewRequest(http.MethodGet, "/api/metrics?limit=1&offset=0", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	var resp SearchResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Metrics) != 1 {
		t.Errorf("expected 1 metric, got %d", len(resp.Metrics))
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Total)
	}
}

func TestAPI_GetMetric(t *testing.T) {
	ms := &mockStore{metrics: newTestMetrics()}
	handler := NewHandler(ms)

	req := httptest.NewRequest(http.MethodGet, "/api/metrics/metric1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var metric domain.CanonicalMetric
	if err := json.NewDecoder(w.Body).Decode(&metric); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if metric.ID != "metric1" {
		t.Errorf("expected metric ID 'metric1', got %q", metric.ID)
	}
}

func TestAPI_GetMetric_NotFound(t *testing.T) {
	ms := &mockStore{metrics: newTestMetrics()}
	handler := NewHandler(ms)

	req := httptest.NewRequest(http.MethodGet, "/api/metrics/nonexistent", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestAPI_GetFacets(t *testing.T) {
	ms := &mockStore{metrics: newTestMetrics()}
	handler := NewHandler(ms)

	req := httptest.NewRequest(http.MethodGet, "/api/facets", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var facets FacetResponse
	if err := json.NewDecoder(w.Body).Decode(&facets); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if facets.InstrumentTypes["counter"] != 5 {
		t.Errorf("expected counter count 5, got %d", facets.InstrumentTypes["counter"])
	}
}

func TestAPI_HealthCheck(t *testing.T) {
	ms := &mockStore{}
	handler := NewHandler(ms)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
