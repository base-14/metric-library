package postgres

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/base14/otel-glossary/internal/adapter"
	"github.com/base14/otel-glossary/internal/domain"
)

func TestPostgresAdapter_Name(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Name() != "prometheus-postgres" {
		t.Errorf("expected name 'prometheus-postgres', got %q", a.Name())
	}
}

func TestPostgresAdapter_SourceCategory(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.SourceCategory() != domain.SourcePrometheus {
		t.Errorf("expected source category 'prometheus', got %q", a.SourceCategory())
	}
}

func TestPostgresAdapter_Confidence(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Confidence() != domain.ConfidenceDerived {
		t.Errorf("expected confidence 'derived', got %q", a.Confidence())
	}
}

func TestPostgresAdapter_ExtractionMethod(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.ExtractionMethod() != domain.ExtractionAST {
		t.Errorf("expected extraction method 'ast', got %q", a.ExtractionMethod())
	}
}

func TestPostgresAdapter_RepoURL(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.RepoURL() != "https://github.com/prometheus-community/postgres_exporter" {
		t.Errorf("unexpected repo URL: %q", a.RepoURL())
	}
}

func TestPostgresAdapter_ImplementsAdapter(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	var _ adapter.Adapter = a
}

func TestPostgresAdapter_Extract(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "postgres-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	collectorDir := filepath.Join(tmpDir, "collector")
	if err := os.MkdirAll(collectorDir, 0750); err != nil {
		t.Fatalf("failed to create collector dir: %v", err)
	}

	goFile := `package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	pgUp = prometheus.NewDesc(
		"pg_up",
		"Whether the PostgreSQL server is up",
		nil,
		nil,
	)
	pgStatDatabaseNumbackends = prometheus.NewDesc(
		"pg_stat_database_numbackends",
		"Number of backends currently connected",
		[]string{"datid", "datname"},
		nil,
	)
	pgReplicationLagBytes = prometheus.NewDesc(
		prometheus.BuildFQName("pg", "replication", "lag_bytes"),
		"Replication lag in bytes",
		[]string{"slot_name"},
		nil,
	)
)
`
	if err := os.WriteFile(filepath.Join(collectorDir, "pg_database.go"), []byte(goFile), 0600); err != nil {
		t.Fatalf("failed to write go file: %v", err)
	}

	a := NewAdapter("/tmp/cache")
	result := &adapter.FetchResult{
		RepoPath: tmpDir,
		Commit:   "abc123",
	}

	metrics, err := a.Extract(context.Background(), result)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 3 {
		t.Fatalf("expected 3 metrics, got %d", len(metrics))
	}

	names := make(map[string]*adapter.RawMetric)
	for _, m := range metrics {
		names[m.Name] = m
	}

	if _, ok := names["pg_up"]; !ok {
		t.Error("expected metric 'pg_up'")
	}
	if _, ok := names["pg_stat_database_numbackends"]; !ok {
		t.Error("expected metric 'pg_stat_database_numbackends'")
	}
	if _, ok := names["pg_replication_lag_bytes"]; !ok {
		t.Error("expected metric 'pg_replication_lag_bytes'")
	}

	m := names["pg_stat_database_numbackends"]
	if m.ComponentName != "pg_database" {
		t.Errorf("expected component name 'pg_database', got %q", m.ComponentName)
	}
	if m.ComponentType != "platform" {
		t.Errorf("expected component type 'platform', got %q", m.ComponentType)
	}
	if len(m.Attributes) != 2 {
		t.Errorf("expected 2 attributes, got %d", len(m.Attributes))
	}
}

func TestPostgresAdapter_Extract_SkipsTestFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "postgres-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	collectorDir := filepath.Join(tmpDir, "collector")
	if err := os.MkdirAll(collectorDir, 0750); err != nil {
		t.Fatalf("failed to create collector dir: %v", err)
	}

	testFile := `package collector

import "github.com/prometheus/client_golang/prometheus"

var testMetric = prometheus.NewDesc(
	"test_metric",
	"This should be skipped",
	nil,
	nil,
)
`
	if err := os.WriteFile(filepath.Join(collectorDir, "pg_database_test.go"), []byte(testFile), 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	a := NewAdapter("/tmp/cache")
	result := &adapter.FetchResult{
		RepoPath: tmpDir,
		Commit:   "abc123",
	}

	metrics, err := a.Extract(context.Background(), result)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 0 {
		t.Errorf("expected 0 metrics (test files should be skipped), got %d", len(metrics))
	}
}

func TestPostgresAdapter_Extract_EmptyRepo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "postgres-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	collectorDir := filepath.Join(tmpDir, "collector")
	if err := os.MkdirAll(collectorDir, 0750); err != nil {
		t.Fatalf("failed to create collector dir: %v", err)
	}

	a := NewAdapter("/tmp/cache")
	result := &adapter.FetchResult{
		RepoPath: tmpDir,
		Commit:   "abc123",
	}

	metrics, err := a.Extract(context.Background(), result)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 0 {
		t.Errorf("expected 0 metrics, got %d", len(metrics))
	}
}

func TestInferInstrumentType(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"pg_stat_database_xact_commit_total", "counter"},
		{"pg_stat_database_blks_read_total", "counter"},
		{"http_requests_total", "counter"},
		{"pg_stat_database_numbackends", "gauge"},
		{"pg_up", "gauge"},
		{"pg_replication_lag_bytes", "gauge"},
		{"pg_settings_shared_buffers_bytes", "gauge"},
		{"request_duration_seconds_bucket", "histogram"},
		{"request_duration_seconds_sum", "counter"},
		{"request_duration_seconds_count", "counter"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := inferInstrumentType(tt.name)
			if got != tt.expected {
				t.Errorf("inferInstrumentType(%q) = %q, want %q", tt.name, got, tt.expected)
			}
		})
	}
}

func TestDeriveComponentName(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"pg_stat_database.go", "pg_stat_database"},
		{"pg_replication.go", "pg_replication"},
		{"collector.go", "collector"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got := deriveComponentName(tt.filename)
			if got != tt.expected {
				t.Errorf("deriveComponentName(%q) = %q, want %q", tt.filename, got, tt.expected)
			}
		})
	}
}

func TestLabelsToAttributes(t *testing.T) {
	labels := []string{"datid", "datname", "state"}
	attrs := labelsToAttributes(labels)

	if len(attrs) != 3 {
		t.Fatalf("expected 3 attributes, got %d", len(attrs))
	}

	for i, label := range labels {
		if attrs[i].Name != label {
			t.Errorf("expected attribute name %q, got %q", label, attrs[i].Name)
		}
		if attrs[i].Type != "string" {
			t.Errorf("expected attribute type 'string', got %q", attrs[i].Type)
		}
	}
}
