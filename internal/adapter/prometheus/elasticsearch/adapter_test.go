package elasticsearch

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
)

func TestElasticsearchAdapter_Name(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Name() != "prometheus-elasticsearch" {
		t.Errorf("expected name 'prometheus-elasticsearch', got %q", a.Name())
	}
}

func TestElasticsearchAdapter_SourceCategory(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.SourceCategory() != domain.SourcePrometheus {
		t.Errorf("expected source category 'prometheus', got %q", a.SourceCategory())
	}
}

func TestElasticsearchAdapter_Confidence(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Confidence() != domain.ConfidenceDerived {
		t.Errorf("expected confidence 'derived', got %q", a.Confidence())
	}
}

func TestElasticsearchAdapter_ExtractionMethod(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.ExtractionMethod() != domain.ExtractionAST {
		t.Errorf("expected extraction method 'ast', got %q", a.ExtractionMethod())
	}
}

func TestElasticsearchAdapter_RepoURL(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.RepoURL() != "https://github.com/prometheus-community/elasticsearch_exporter" {
		t.Errorf("unexpected repo URL: %q", a.RepoURL())
	}
}

func TestElasticsearchAdapter_ImplementsAdapter(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	var _ adapter.Adapter = a
}

func TestElasticsearchAdapter_Extract(t *testing.T) {
	tmpDir := setupTestRepo(t)

	goFile := `package collector

import "github.com/prometheus/client_golang/prometheus"

const namespace = "elasticsearch"

var (
	activePrimaryShards = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "cluster_health", "active_primary_shards"),
		"The number of primary shards in your cluster.",
		[]string{"cluster"}, nil,
	)
	docsCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "indices", "docs_primary"),
		"Count of documents with only primary shards",
		[]string{"index", "cluster"}, nil,
	)
	snapshotCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "snapshot_stats", "number_of_snapshots"),
		"Number of snapshots in a repository.",
		nil, nil,
	)
)
`
	writeFile(t, filepath.Join(tmpDir, "collector", "cluster_health.go"), goFile)

	a := NewAdapter("/tmp/cache")
	metrics, err := a.Extract(context.Background(), &adapter.FetchResult{
		RepoPath: tmpDir,
		Commit:   "abc123",
	})
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 3 {
		t.Fatalf("expected 3 metrics, got %d", len(metrics))
	}

	names := metricsByName(metrics)

	m := names["elasticsearch_cluster_health_active_primary_shards"]
	if m == nil {
		t.Fatal("expected metric 'elasticsearch_cluster_health_active_primary_shards'")
	}
	if m.Description != "The number of primary shards in your cluster." {
		t.Errorf("unexpected description: %q", m.Description)
	}
	if len(m.Attributes) != 1 || m.Attributes[0].Name != "cluster" {
		t.Errorf("expected 1 attribute 'cluster', got %v", m.Attributes)
	}
	if m.InstrumentType != "gauge" {
		t.Errorf("expected gauge, got %q", m.InstrumentType)
	}
	if m.ComponentName != "cluster_health" {
		t.Errorf("expected component 'cluster_health', got %q", m.ComponentName)
	}

	docs := names["elasticsearch_indices_docs_primary"]
	if docs == nil {
		t.Fatal("expected metric 'elasticsearch_indices_docs_primary'")
	}
	if len(docs.Attributes) != 2 {
		t.Errorf("expected 2 attributes, got %d", len(docs.Attributes))
	}
}

func TestElasticsearchAdapter_Extract_StructPattern(t *testing.T) {
	tmpDir := setupTestRepo(t)

	// Metrics defined inside struct arrays (nodes.go pattern)
	goFile := `package collector

import "github.com/prometheus/client_golang/prometheus"

const namespace = "elasticsearch"

func NewNodes() *Nodes {
	return &Nodes{
		nodeMetrics: []*nodeMetric{
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "os", "load1"),
					"Shortterm load average",
					defaultNodeLabels, nil,
				),
			},
			{
				Type: prometheus.CounterValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "os", "cpu_percent"),
					"Percent CPU used by the OS",
					defaultNodeLabels, nil,
				),
			},
		},
	}
}
`
	writeFile(t, filepath.Join(tmpDir, "collector", "nodes.go"), goFile)

	a := NewAdapter("/tmp/cache")
	metrics, err := a.Extract(context.Background(), &adapter.FetchResult{
		RepoPath: tmpDir,
		Commit:   "abc123",
	})
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 2 {
		t.Fatalf("expected 2 metrics, got %d", len(metrics))
	}

	names := metricsByName(metrics)
	if names["elasticsearch_os_load1"] == nil {
		t.Error("expected metric 'elasticsearch_os_load1'")
	}
	if names["elasticsearch_os_cpu_percent"] == nil {
		t.Error("expected metric 'elasticsearch_os_cpu_percent'")
	}
}

func TestElasticsearchAdapter_Extract_SkipsTestFiles(t *testing.T) {
	tmpDir := setupTestRepo(t)

	testFile := `package collector

import "github.com/prometheus/client_golang/prometheus"

const namespace = "elasticsearch"

var testMetric = prometheus.NewDesc(
	prometheus.BuildFQName(namespace, "test", "metric"),
	"Test metric.",
	nil, nil,
)
`
	writeFile(t, filepath.Join(tmpDir, "collector", "nodes_test.go"), testFile)

	a := NewAdapter("/tmp/cache")
	metrics, err := a.Extract(context.Background(), &adapter.FetchResult{
		RepoPath: tmpDir,
		Commit:   "abc123",
	})
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 0 {
		t.Errorf("expected 0 metrics (test files should be skipped), got %d", len(metrics))
	}
}

func TestElasticsearchAdapter_Extract_EmptyRepo(t *testing.T) {
	tmpDir := setupTestRepo(t)

	a := NewAdapter("/tmp/cache")
	metrics, err := a.Extract(context.Background(), &adapter.FetchResult{
		RepoPath: tmpDir,
		Commit:   "abc123",
	})
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
		{"elasticsearch_indices_indexing_index_total", "counter"},
		{"elasticsearch_cluster_health_active_primary_shards", "gauge"},
		{"elasticsearch_os_load1", "gauge"},
		{"elasticsearch_indices_store_size_bytes_primary", "gauge"},
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
		{"cluster_health.go", "cluster_health"},
		{"nodes.go", "nodes"},
		{"indices_mappings.go", "indices_mappings"},
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

// Test helpers

func setupTestRepo(t *testing.T) string {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "es-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tmpDir) })

	if err := os.MkdirAll(filepath.Join(tmpDir, "collector"), 0750); err != nil {
		t.Fatalf("failed to create collector dir: %v", err)
	}
	return tmpDir
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write %s: %v", path, err)
	}
}

func metricsByName(metrics []*adapter.RawMetric) map[string]*adapter.RawMetric {
	m := make(map[string]*adapter.RawMetric)
	for _, metric := range metrics {
		m[metric.Name] = metric
	}
	return m
}
