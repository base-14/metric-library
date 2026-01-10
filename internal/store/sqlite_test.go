package store

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/base-14/metric-library/internal/domain"
)

func setupTestStore(t *testing.T) *SQLiteStore {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	_ = tmpFile.Close()

	t.Cleanup(func() {
		_ = os.Remove(tmpFile.Name())
	})

	store, err := NewSQLiteStore(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	t.Cleanup(func() {
		_ = store.Close()
	})

	if err := store.RunMigrations(context.Background()); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return store
}

func testMetric() *domain.CanonicalMetric {
	return &domain.CanonicalMetric{
		MetricName:       "system.cpu.utilization",
		InstrumentType:   domain.InstrumentGauge,
		Description:      "CPU utilization",
		Unit:             "1",
		EnabledByDefault: true,
		ComponentType:    domain.ComponentReceiver,
		ComponentName:    "hostmetrics",
		SourceCategory:   domain.SourceOTEL,
		SourceName:       "opentelemetry-collector-contrib",
		ExtractionMethod: domain.ExtractionMetadata,
		SourceConfidence: domain.ConfidenceAuthoritative,
		Repo:             "https://github.com/open-telemetry/opentelemetry-collector-contrib",
		Path:             "receiver/hostmetricsreceiver/metadata.yaml",
		Commit:           "abc123",
		ExtractedAt:      time.Now(),
		Attributes: []domain.Attribute{
			{
				Name:        "cpu",
				Type:        "string",
				Description: "CPU number",
				Required:    false,
			},
			{
				Name:        "state",
				Type:        "string",
				Description: "CPU state",
				Required:    true,
				Enum:        []string{"idle", "user", "system"},
			},
		},
	}
}

func TestSQLiteStore_UpsertAndGetMetric(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	metric := testMetric()
	if err := store.UpsertMetric(ctx, metric); err != nil {
		t.Fatalf("UpsertMetric failed: %v", err)
	}

	got, err := store.GetMetric(ctx, metric.ID)
	if err != nil {
		t.Fatalf("GetMetric failed: %v", err)
	}

	if got == nil {
		t.Fatal("GetMetric returned nil")
		return
	}

	if got.MetricName != metric.MetricName {
		t.Errorf("MetricName = %q, want %q", got.MetricName, metric.MetricName)
	}

	if got.InstrumentType != metric.InstrumentType {
		t.Errorf("InstrumentType = %q, want %q", got.InstrumentType, metric.InstrumentType)
	}

	if len(got.Attributes) != len(metric.Attributes) {
		t.Fatalf("Attributes count = %d, want %d", len(got.Attributes), len(metric.Attributes))
	}

	// Check enum values
	for i, attr := range got.Attributes {
		if len(attr.Enum) != len(metric.Attributes[i].Enum) {
			t.Errorf("Attribute[%d].Enum count = %d, want %d", i, len(attr.Enum), len(metric.Attributes[i].Enum))
		}
	}
}

func TestSQLiteStore_UpsertMetric_Update(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	metric := testMetric()
	if err := store.UpsertMetric(ctx, metric); err != nil {
		t.Fatalf("UpsertMetric failed: %v", err)
	}

	// Update the metric
	metric.Description = "Updated description"
	metric.Attributes = []domain.Attribute{
		{Name: "new_attr", Type: "int"},
	}

	if err := store.UpsertMetric(ctx, metric); err != nil {
		t.Fatalf("UpsertMetric update failed: %v", err)
	}

	got, err := store.GetMetric(ctx, metric.ID)
	if err != nil {
		t.Fatalf("GetMetric failed: %v", err)
	}

	if got.Description != "Updated description" {
		t.Errorf("Description = %q, want %q", got.Description, "Updated description")
	}

	if len(got.Attributes) != 1 {
		t.Errorf("Attributes count = %d, want 1", len(got.Attributes))
	}

	if got.Attributes[0].Name != "new_attr" {
		t.Errorf("Attribute name = %q, want %q", got.Attributes[0].Name, "new_attr")
	}
}

func TestSQLiteStore_UpsertMetrics_Batch(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	metrics := []*domain.CanonicalMetric{
		testMetric(),
		func() *domain.CanonicalMetric {
			m := testMetric()
			m.MetricName = "system.memory.usage"
			return m
		}(),
	}

	if err := store.UpsertMetrics(ctx, metrics); err != nil {
		t.Fatalf("UpsertMetrics failed: %v", err)
	}

	for _, m := range metrics {
		got, err := store.GetMetric(ctx, m.ID)
		if err != nil {
			t.Fatalf("GetMetric failed: %v", err)
		}
		if got == nil {
			t.Errorf("Metric %s not found", m.ID)
		}
	}
}

func TestSQLiteStore_GetMetric_NotFound(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	got, err := store.GetMetric(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("GetMetric failed: %v", err)
	}

	if got != nil {
		t.Error("GetMetric should return nil for nonexistent metric")
	}
}

func TestSQLiteStore_DeleteMetric(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	metric := testMetric()
	if err := store.UpsertMetric(ctx, metric); err != nil {
		t.Fatalf("UpsertMetric failed: %v", err)
	}

	if err := store.DeleteMetric(ctx, metric.ID); err != nil {
		t.Fatalf("DeleteMetric failed: %v", err)
	}

	got, err := store.GetMetric(ctx, metric.ID)
	if err != nil {
		t.Fatalf("GetMetric failed: %v", err)
	}

	if got != nil {
		t.Error("Metric should be deleted")
	}
}

func TestSQLiteStore_DeleteMetricsBySource(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	metrics := []*domain.CanonicalMetric{
		testMetric(),
		func() *domain.CanonicalMetric {
			m := testMetric()
			m.MetricName = "system.memory.usage"
			return m
		}(),
		func() *domain.CanonicalMetric {
			m := testMetric()
			m.MetricName = "other.metric"
			m.SourceName = "other-source"
			return m
		}(),
	}

	if err := store.UpsertMetrics(ctx, metrics); err != nil {
		t.Fatalf("UpsertMetrics failed: %v", err)
	}

	if err := store.DeleteMetricsBySource(ctx, "opentelemetry-collector-contrib"); err != nil {
		t.Fatalf("DeleteMetricsBySource failed: %v", err)
	}

	// First two should be deleted
	for i := 0; i < 2; i++ {
		got, _ := store.GetMetric(ctx, metrics[i].ID)
		if got != nil {
			t.Errorf("Metric %d should be deleted", i)
		}
	}

	// Third should remain
	got, _ := store.GetMetric(ctx, metrics[2].ID)
	if got == nil {
		t.Error("Third metric should not be deleted")
	}
}

func TestSQLiteStore_Search_Text(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	metrics := []*domain.CanonicalMetric{
		testMetric(),
		func() *domain.CanonicalMetric {
			m := testMetric()
			m.MetricName = "http.server.request.duration"
			m.Description = "HTTP request duration"
			return m
		}(),
	}

	if err := store.UpsertMetrics(ctx, metrics); err != nil {
		t.Fatalf("UpsertMetrics failed: %v", err)
	}

	result, err := store.Search(ctx, SearchQuery{Text: "cpu"})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if result.Total != 1 {
		t.Errorf("Total = %d, want 1", result.Total)
	}

	if len(result.Metrics) != 1 {
		t.Errorf("Metrics count = %d, want 1", len(result.Metrics))
	}
}

func TestSQLiteStore_Search_TextOrdering(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	metrics := []*domain.CanonicalMetric{
		func() *domain.CanonicalMetric {
			m := testMetric()
			m.MetricName = "http.server.request.duration"
			m.Description = "Measures memory usage"
			return m
		}(),
		func() *domain.CanonicalMetric {
			m := testMetric()
			m.MetricName = "system.memory.utilization"
			m.Description = "System memory utilization"
			return m
		}(),
		func() *domain.CanonicalMetric {
			m := testMetric()
			m.MetricName = "process.runtime.heap"
			m.Description = "Process heap memory usage"
			return m
		}(),
	}

	if err := store.UpsertMetrics(ctx, metrics); err != nil {
		t.Fatalf("UpsertMetrics failed: %v", err)
	}

	result, err := store.Search(ctx, SearchQuery{Text: "memory"})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if result.Total != 3 {
		t.Errorf("Total = %d, want 3", result.Total)
	}

	if len(result.Metrics) != 3 {
		t.Fatalf("Metrics count = %d, want 3", len(result.Metrics))
	}

	// First should be the one with "memory" in metric_name
	if result.Metrics[0].MetricName != "system.memory.utilization" {
		t.Errorf("First result = %s, want system.memory.utilization (metric_name match)", result.Metrics[0].MetricName)
	}

	// The other two have "memory" only in description, ordered alphabetically
	if result.Metrics[1].MetricName != "http.server.request.duration" {
		t.Errorf("Second result = %s, want http.server.request.duration", result.Metrics[1].MetricName)
	}

	if result.Metrics[2].MetricName != "process.runtime.heap" {
		t.Errorf("Third result = %s, want process.runtime.heap", result.Metrics[2].MetricName)
	}
}

func TestSQLiteStore_Search_TextOnlyMatchesMetricNameAndDescription(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	metrics := []*domain.CanonicalMetric{
		func() *domain.CanonicalMetric {
			m := testMetric()
			m.MetricName = "http.request.duration"
			m.Description = "HTTP request duration"
			m.ComponentName = "kafkareceiver"
			m.SourceName = "kafka-exporter"
			return m
		}(),
	}

	if err := store.UpsertMetrics(ctx, metrics); err != nil {
		t.Fatalf("UpsertMetrics failed: %v", err)
	}

	// Should NOT match on component_name
	result, err := store.Search(ctx, SearchQuery{Text: "kafka"})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if result.Total != 0 {
		t.Errorf("Total = %d, want 0 (should not match component_name or source_name)", result.Total)
	}

	// Should match on metric_name
	result, err = store.Search(ctx, SearchQuery{Text: "request"})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if result.Total != 1 {
		t.Errorf("Total = %d, want 1 (should match metric_name)", result.Total)
	}
}

func TestSQLiteStore_Search_Filters(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	metrics := []*domain.CanonicalMetric{
		testMetric(),
		func() *domain.CanonicalMetric {
			m := testMetric()
			m.MetricName = "http.duration"
			m.InstrumentType = domain.InstrumentHistogram
			return m
		}(),
	}

	if err := store.UpsertMetrics(ctx, metrics); err != nil {
		t.Fatalf("UpsertMetrics failed: %v", err)
	}

	tests := []struct {
		name  string
		query SearchQuery
		want  int
	}{
		{
			name:  "filter by instrument type",
			query: SearchQuery{InstrumentTypes: []domain.InstrumentType{domain.InstrumentGauge}},
			want:  1,
		},
		{
			name:  "filter by component type",
			query: SearchQuery{ComponentTypes: []domain.ComponentType{domain.ComponentReceiver}},
			want:  2,
		},
		{
			name:  "filter by source category",
			query: SearchQuery{SourceCategories: []domain.SourceCategory{domain.SourceOTEL}},
			want:  2,
		},
		{
			name:  "filter by confidence",
			query: SearchQuery{ConfidenceLevels: []domain.ConfidenceLevel{domain.ConfidenceAuthoritative}},
			want:  2,
		},
		{
			name:  "combined filters",
			query: SearchQuery{InstrumentTypes: []domain.InstrumentType{domain.InstrumentHistogram}, ComponentTypes: []domain.ComponentType{domain.ComponentReceiver}},
			want:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := store.Search(ctx, tt.query)
			if err != nil {
				t.Fatalf("Search failed: %v", err)
			}

			if result.Total != tt.want {
				t.Errorf("Total = %d, want %d", result.Total, tt.want)
			}
		})
	}
}

func TestSQLiteStore_Search_Pagination(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	for i := 0; i < 25; i++ {
		m := testMetric()
		m.MetricName = "metric." + string(rune('a'+i))
		if err := store.UpsertMetric(ctx, m); err != nil {
			t.Fatalf("UpsertMetric failed: %v", err)
		}
	}

	result, err := store.Search(ctx, SearchQuery{Limit: 10, Offset: 0})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if result.Total != 25 {
		t.Errorf("Total = %d, want 25", result.Total)
	}

	if len(result.Metrics) != 10 {
		t.Errorf("Metrics count = %d, want 10", len(result.Metrics))
	}

	result2, err := store.Search(ctx, SearchQuery{Limit: 10, Offset: 20})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(result2.Metrics) != 5 {
		t.Errorf("Metrics count = %d, want 5", len(result2.Metrics))
	}
}

func TestSQLiteStore_GetFacetCounts(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	metrics := []*domain.CanonicalMetric{
		testMetric(),
		func() *domain.CanonicalMetric {
			m := testMetric()
			m.MetricName = "http.duration"
			m.InstrumentType = domain.InstrumentHistogram
			return m
		}(),
	}

	if err := store.UpsertMetrics(ctx, metrics); err != nil {
		t.Fatalf("UpsertMetrics failed: %v", err)
	}

	facets, err := store.GetFacetCounts(ctx)
	if err != nil {
		t.Fatalf("GetFacetCounts failed: %v", err)
	}

	if facets.InstrumentTypes[domain.InstrumentGauge] != 1 {
		t.Errorf("Gauge count = %d, want 1", facets.InstrumentTypes[domain.InstrumentGauge])
	}

	if facets.InstrumentTypes[domain.InstrumentHistogram] != 1 {
		t.Errorf("Histogram count = %d, want 1", facets.InstrumentTypes[domain.InstrumentHistogram])
	}

	if facets.ComponentTypes[domain.ComponentReceiver] != 2 {
		t.Errorf("Receiver count = %d, want 2", facets.ComponentTypes[domain.ComponentReceiver])
	}
}

func TestSQLiteStore_ExtractionRun(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	run := &ExtractionRun{
		ID:          "run-123",
		AdapterName: "otel-collector-contrib",
		Commit:      "abc123",
		StartedAt:   time.Now(),
		Status:      "running",
	}

	if err := store.CreateExtractionRun(ctx, run); err != nil {
		t.Fatalf("CreateExtractionRun failed: %v", err)
	}

	got, err := store.GetExtractionRun(ctx, "run-123")
	if err != nil {
		t.Fatalf("GetExtractionRun failed: %v", err)
	}

	if got == nil {
		t.Fatal("GetExtractionRun returned nil")
		return
	}

	if got.Status != "running" {
		t.Errorf("Status = %q, want %q", got.Status, "running")
	}

	// Update
	now := time.Now()
	run.CompletedAt = &now
	run.MetricsCount = 100
	run.Status = "completed"

	if err := store.UpdateExtractionRun(ctx, run); err != nil {
		t.Fatalf("UpdateExtractionRun failed: %v", err)
	}

	got, _ = store.GetExtractionRun(ctx, "run-123")
	if got.Status != "completed" {
		t.Errorf("Status = %q, want %q", got.Status, "completed")
	}
	if got.MetricsCount != 100 {
		t.Errorf("MetricsCount = %d, want 100", got.MetricsCount)
	}
}

func TestSQLiteStore_GetLatestExtractionRun(t *testing.T) {
	store := setupTestStore(t)
	ctx := context.Background()

	runs := []*ExtractionRun{
		{ID: "run-1", AdapterName: "adapter1", StartedAt: time.Now().Add(-2 * time.Hour), Status: "completed"},
		{ID: "run-2", AdapterName: "adapter1", StartedAt: time.Now().Add(-1 * time.Hour), Status: "completed"},
		{ID: "run-3", AdapterName: "adapter2", StartedAt: time.Now(), Status: "running"},
	}

	for _, run := range runs {
		if err := store.CreateExtractionRun(ctx, run); err != nil {
			t.Fatalf("CreateExtractionRun failed: %v", err)
		}
	}

	got, err := store.GetLatestExtractionRun(ctx, "adapter1")
	if err != nil {
		t.Fatalf("GetLatestExtractionRun failed: %v", err)
	}

	if got == nil {
		t.Fatal("GetLatestExtractionRun returned nil")
		return
	}

	if got.ID != "run-2" {
		t.Errorf("ID = %q, want %q", got.ID, "run-2")
	}
}
