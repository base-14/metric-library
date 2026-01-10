package redis

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/base14/otel-glossary/internal/adapter"
	"github.com/base14/otel-glossary/internal/domain"
)

func TestRedisAdapter_Name(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Name() != "prometheus-redis" {
		t.Errorf("expected name 'prometheus-redis', got %q", a.Name())
	}
}

func TestRedisAdapter_SourceCategory(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.SourceCategory() != domain.SourcePrometheus {
		t.Errorf("expected source category 'prometheus', got %q", a.SourceCategory())
	}
}

func TestRedisAdapter_Confidence(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Confidence() != domain.ConfidenceDerived {
		t.Errorf("expected confidence 'derived', got %q", a.Confidence())
	}
}

func TestRedisAdapter_ExtractionMethod(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.ExtractionMethod() != domain.ExtractionAST {
		t.Errorf("expected extraction method 'ast', got %q", a.ExtractionMethod())
	}
}

func TestRedisAdapter_RepoURL(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.RepoURL() != "https://github.com/oliver006/redis_exporter" {
		t.Errorf("unexpected repo URL: %q", a.RepoURL())
	}
}

func TestRedisAdapter_ImplementsAdapter(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	var _ adapter.Adapter = a
}

func TestRedisAdapter_Extract_MetricDescriptions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "redis-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	exporterDir := filepath.Join(tmpDir, "exporter")
	if err := os.MkdirAll(exporterDir, 0750); err != nil {
		t.Fatalf("failed to create exporter dir: %v", err)
	}

	goFile := `package exporter

var metricDescriptions = map[string]struct {
	txt  string
	lbls []string
}{
	"commands_duration_seconds_total": {
		txt:  "Total amount of time in seconds spent per command",
		lbls: []string{"cmd"},
	},
	"connected_slave_lag_seconds": {
		txt:  "Lag of connected slave",
		lbls: []string{"slave_ip", "slave_port"},
	},
	"db_keys": {
		txt:  "Total number of keys by DB",
		lbls: []string{"db"},
	},
}
`
	if err := os.WriteFile(filepath.Join(exporterDir, "exporter.go"), []byte(goFile), 0600); err != nil {
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

	if _, ok := names["redis_commands_duration_seconds_total"]; !ok {
		t.Error("expected metric 'redis_commands_duration_seconds_total'")
	}
	if _, ok := names["redis_connected_slave_lag_seconds"]; !ok {
		t.Error("expected metric 'redis_connected_slave_lag_seconds'")
	}
	if _, ok := names["redis_db_keys"]; !ok {
		t.Error("expected metric 'redis_db_keys'")
	}

	m := names["redis_commands_duration_seconds_total"]
	if m.Description != "Total amount of time in seconds spent per command" {
		t.Errorf("unexpected description: %q", m.Description)
	}
	if len(m.Attributes) != 1 || m.Attributes[0].Name != "cmd" {
		t.Errorf("expected 1 attribute 'cmd', got %v", m.Attributes)
	}
	if m.InstrumentType != "counter" {
		t.Errorf("expected instrument type 'counter', got %q", m.InstrumentType)
	}

	m2 := names["redis_connected_slave_lag_seconds"]
	if len(m2.Attributes) != 2 {
		t.Errorf("expected 2 attributes, got %d", len(m2.Attributes))
	}
}

func TestRedisAdapter_Extract_MetricMapGauges(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "redis-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	exporterDir := filepath.Join(tmpDir, "exporter")
	if err := os.MkdirAll(exporterDir, 0750); err != nil {
		t.Fatalf("failed to create exporter dir: %v", err)
	}

	goFile := `package exporter

var metricMapGauges = map[string]string{
	"connected_clients": "connected_clients",
	"blocked_clients":   "blocked_clients",
	"used_memory":       "used_memory",
}
`
	if err := os.WriteFile(filepath.Join(exporterDir, "gauges.go"), []byte(goFile), 0600); err != nil {
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

	for _, m := range metrics {
		if m.InstrumentType != "gauge" {
			t.Errorf("expected gauge, got %q for %s", m.InstrumentType, m.Name)
		}
		if m.ComponentName != "redis" {
			t.Errorf("expected component name 'redis', got %q", m.ComponentName)
		}
	}
}

func TestRedisAdapter_Extract_MetricMapCounters(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "redis-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	exporterDir := filepath.Join(tmpDir, "exporter")
	if err := os.MkdirAll(exporterDir, 0750); err != nil {
		t.Fatalf("failed to create exporter dir: %v", err)
	}

	goFile := `package exporter

var metricMapCounters = map[string]string{
	"total_connections_received_total": "total_connections_received",
	"total_commands_processed_total":   "total_commands_processed",
	"expired_keys_total":               "expired_keys",
}
`
	if err := os.WriteFile(filepath.Join(exporterDir, "counters.go"), []byte(goFile), 0600); err != nil {
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

	for _, m := range metrics {
		if m.InstrumentType != "counter" {
			t.Errorf("expected counter, got %q for %s", m.InstrumentType, m.Name)
		}
	}
}

func TestRedisAdapter_Extract_SkipsTestFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "redis-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	exporterDir := filepath.Join(tmpDir, "exporter")
	if err := os.MkdirAll(exporterDir, 0750); err != nil {
		t.Fatalf("failed to create exporter dir: %v", err)
	}

	testFile := `package exporter

var metricMapGauges = map[string]string{
	"test_metric": "test_metric",
}
`
	if err := os.WriteFile(filepath.Join(exporterDir, "exporter_test.go"), []byte(testFile), 0600); err != nil {
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

func TestInferInstrumentType(t *testing.T) {
	tests := []struct {
		name       string
		metricType string
		expected   string
	}{
		{"commands_total", "", "counter"},
		{"expired_keys_total", "", "counter"},
		{"connected_clients", "", "gauge"},
		{"used_memory", "", "gauge"},
		{"latency_bucket", "", "histogram"},
		{"some_metric", "gauge", "gauge"},
		{"some_metric", "counter", "counter"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := inferInstrumentType(tt.name, tt.metricType)
			if got != tt.expected {
				t.Errorf("inferInstrumentType(%q, %q) = %q, want %q", tt.name, tt.metricType, got, tt.expected)
			}
		})
	}
}
