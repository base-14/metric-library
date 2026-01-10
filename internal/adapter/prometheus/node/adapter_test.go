package node

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/base14/otel-glossary/internal/adapter"
	"github.com/base14/otel-glossary/internal/domain"
)

func TestNodeAdapter_Name(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Name() != "prometheus-node" {
		t.Errorf("expected name 'prometheus-node', got %q", a.Name())
	}
}

func TestNodeAdapter_SourceCategory(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.SourceCategory() != domain.SourcePrometheus {
		t.Errorf("expected source category 'prometheus', got %q", a.SourceCategory())
	}
}

func TestNodeAdapter_Confidence(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Confidence() != domain.ConfidenceDerived {
		t.Errorf("expected confidence 'derived', got %q", a.Confidence())
	}
}

func TestNodeAdapter_ExtractionMethod(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.ExtractionMethod() != domain.ExtractionAST {
		t.Errorf("expected extraction method 'ast', got %q", a.ExtractionMethod())
	}
}

func TestNodeAdapter_RepoURL(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.RepoURL() != "https://github.com/prometheus/node_exporter" {
		t.Errorf("unexpected repo URL: %q", a.RepoURL())
	}
}

func TestNodeAdapter_ImplementsAdapter(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	var _ adapter.Adapter = a
}

func TestNodeAdapter_Extract(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "node-adapter-test-*")
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
	cpuSecondsTotal = prometheus.NewDesc(
		prometheus.BuildFQName("node", "cpu", "seconds_total"),
		"Seconds the CPUs spent in each mode.",
		[]string{"cpu", "mode"},
		nil,
	)
	memoryBytes = prometheus.NewDesc(
		prometheus.BuildFQName("node", "memory", "MemTotal_bytes"),
		"Memory information field MemTotal_bytes.",
		nil,
		nil,
	)
	diskReadBytesTotal = prometheus.NewDesc(
		prometheus.BuildFQName("node", "disk", "read_bytes_total"),
		"The total number of bytes read successfully.",
		[]string{"device"},
		nil,
	)
)
`
	if err := os.WriteFile(filepath.Join(collectorDir, "cpu_linux.go"), []byte(goFile), 0600); err != nil {
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

	if _, ok := names["node_cpu_seconds_total"]; !ok {
		t.Error("expected metric 'node_cpu_seconds_total'")
	}
	if _, ok := names["node_memory_MemTotal_bytes"]; !ok {
		t.Error("expected metric 'node_memory_MemTotal_bytes'")
	}
	if _, ok := names["node_disk_read_bytes_total"]; !ok {
		t.Error("expected metric 'node_disk_read_bytes_total'")
	}

	m := names["node_cpu_seconds_total"]
	if m.ComponentName != "cpu" {
		t.Errorf("expected component name 'cpu', got %q", m.ComponentName)
	}
	if m.ComponentType != "platform" {
		t.Errorf("expected component type 'platform', got %q", m.ComponentType)
	}
	if len(m.Attributes) != 2 {
		t.Errorf("expected 2 attributes, got %d", len(m.Attributes))
	}
	if m.InstrumentType != "counter" {
		t.Errorf("expected instrument type 'counter', got %q", m.InstrumentType)
	}
}

func TestNodeAdapter_Extract_SkipsTestFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "node-adapter-test-*")
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
	if err := os.WriteFile(filepath.Join(collectorDir, "cpu_test.go"), []byte(testFile), 0600); err != nil {
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

func TestNodeAdapter_Extract_EmptyRepo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "node-adapter-test-*")
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

func TestDeriveComponentName(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"cpu_linux.go", "cpu"},
		{"cpu_darwin.go", "cpu"},
		{"memory_linux.go", "memory"},
		{"diskstats_linux.go", "diskstats"},
		{"netdev_common.go", "netdev"},
		{"filesystem_freebsd.go", "filesystem"},
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

func TestInferInstrumentType(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"node_cpu_seconds_total", "counter"},
		{"node_disk_read_bytes_total", "counter"},
		{"node_memory_MemTotal_bytes", "gauge"},
		{"node_filesystem_size_bytes", "gauge"},
		{"node_load1", "gauge"},
		{"node_network_receive_bytes_total", "counter"},
		{"http_request_duration_seconds_bucket", "histogram"},
		{"http_request_duration_seconds_sum", "counter"},
		{"http_request_duration_seconds_count", "counter"},
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
