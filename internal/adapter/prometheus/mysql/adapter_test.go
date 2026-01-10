package mysql

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
)

func TestMySQLAdapter_Name(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Name() != "prometheus-mysql" {
		t.Errorf("expected name 'prometheus-mysql', got %q", a.Name())
	}
}

func TestMySQLAdapter_SourceCategory(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.SourceCategory() != domain.SourcePrometheus {
		t.Errorf("expected source category 'prometheus', got %q", a.SourceCategory())
	}
}

func TestMySQLAdapter_Confidence(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Confidence() != domain.ConfidenceDerived {
		t.Errorf("expected confidence 'derived', got %q", a.Confidence())
	}
}

func TestMySQLAdapter_ExtractionMethod(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.ExtractionMethod() != domain.ExtractionAST {
		t.Errorf("expected extraction method 'ast', got %q", a.ExtractionMethod())
	}
}

func TestMySQLAdapter_RepoURL(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.RepoURL() != "https://github.com/prometheus/mysqld_exporter" {
		t.Errorf("unexpected repo URL: %q", a.RepoURL())
	}
}

func TestMySQLAdapter_ImplementsAdapter(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	var _ adapter.Adapter = a
}

func TestMySQLAdapter_Extract(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mysql-adapter-test-*")
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

const namespace = "mysql"

var (
	globalStatusDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "global_status", "uptime_seconds"),
		"Number of seconds since the server started.",
		nil, nil,
	)
	queriesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "global_status", "queries_total"),
		"Total number of queries executed by the server.",
		[]string{"type"}, nil,
	)
)
`
	if err := os.WriteFile(filepath.Join(collectorDir, "global_status.go"), []byte(goFile), 0600); err != nil {
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

	if len(metrics) != 2 {
		t.Fatalf("expected 2 metrics, got %d", len(metrics))
	}

	names := make(map[string]*adapter.RawMetric)
	for _, m := range metrics {
		names[m.Name] = m
	}

	if _, ok := names["mysql_global_status_uptime_seconds"]; !ok {
		t.Error("expected metric 'mysql_global_status_uptime_seconds'")
	}
	if _, ok := names["mysql_global_status_queries_total"]; !ok {
		t.Error("expected metric 'mysql_global_status_queries_total'")
	}

	m := names["mysql_global_status_queries_total"]
	if m.Description != "Total number of queries executed by the server." {
		t.Errorf("unexpected description: %q", m.Description)
	}
	if len(m.Attributes) != 1 || m.Attributes[0].Name != "type" {
		t.Errorf("expected 1 attribute 'type', got %v", m.Attributes)
	}
	if m.InstrumentType != "counter" {
		t.Errorf("expected instrument type 'counter', got %q", m.InstrumentType)
	}
	if m.ComponentName != "global_status" {
		t.Errorf("expected component name 'global_status', got %q", m.ComponentName)
	}
}

func TestMySQLAdapter_Extract_SkipsTestFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mysql-adapter-test-*")
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

const namespace = "mysql"

var testDesc = prometheus.NewDesc(
	prometheus.BuildFQName(namespace, "test", "metric"),
	"Test metric.",
	nil, nil,
)
`
	if err := os.WriteFile(filepath.Join(collectorDir, "test_collector_test.go"), []byte(testFile), 0600); err != nil {
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

func TestMySQLAdapter_Extract_EmptyRepo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mysql-adapter-test-*")
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
		t.Errorf("expected 0 metrics for empty repo, got %d", len(metrics))
	}
}

func TestInferInstrumentType(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"mysql_global_status_queries_total", "counter"},
		{"mysql_global_status_threads_running", "gauge"},
		{"mysql_global_status_uptime_seconds", "gauge"},
		{"mysql_info_schema_table_rows", "gauge"},
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
		{"global_status.go", "global_status"},
		{"info_schema_innodb.go", "info_schema_innodb"},
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
