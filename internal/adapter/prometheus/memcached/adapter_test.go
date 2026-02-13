package memcached

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
)

func TestMemcachedAdapter_Name(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Name() != "prometheus-memcached" {
		t.Errorf("expected name 'prometheus-memcached', got %q", a.Name())
	}
}

func TestMemcachedAdapter_SourceCategory(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.SourceCategory() != domain.SourcePrometheus {
		t.Errorf("expected source category 'prometheus', got %q", a.SourceCategory())
	}
}

func TestMemcachedAdapter_Confidence(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Confidence() != domain.ConfidenceDerived {
		t.Errorf("expected confidence 'derived', got %q", a.Confidence())
	}
}

func TestMemcachedAdapter_ExtractionMethod(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.ExtractionMethod() != domain.ExtractionAST {
		t.Errorf("expected extraction method 'ast', got %q", a.ExtractionMethod())
	}
}

func TestMemcachedAdapter_RepoURL(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.RepoURL() != "https://github.com/prometheus/memcached_exporter" {
		t.Errorf("unexpected repo URL: %q", a.RepoURL())
	}
}

func TestMemcachedAdapter_ImplementsAdapter(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	var _ adapter.Adapter = a
}

func TestMemcachedAdapter_Extract(t *testing.T) {
	tmpDir := setupTestRepo(t)

	goFile := `package exporter

import "github.com/prometheus/client_golang/prometheus"

const (
	Namespace           = "memcached"
	subsystemLruCrawler = "lru_crawler"
	subsystemSlab       = "slab"
)

func New() *Exporter {
	return &Exporter{
		up: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "up"),
			"Could the memcached server be reached.",
			nil, nil,
		),
		commands: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "commands_total"),
			"Total number of all requests broken down by command and status.",
			[]string{"command", "status"}, nil,
		),
		lruCrawlerEnabled: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystemLruCrawler, "enabled"),
			"Whether the LRU crawler is enabled.",
			nil, nil,
		),
		slabHits: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystemSlab, "lru_hits_total"),
			"Number of get_hits to the LRU.",
			[]string{"slab", "lru"}, nil,
		),
	}
}
`
	writeFile(t, filepath.Join(tmpDir, "pkg", "exporter", "exporter.go"), goFile)

	a := NewAdapter("/tmp/cache")
	metrics, err := a.Extract(context.Background(), &adapter.FetchResult{
		RepoPath: tmpDir,
		Commit:   "abc123",
	})
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 4 {
		t.Fatalf("expected 4 metrics, got %d", len(metrics))
	}

	names := metricsByName(metrics)

	m := names["memcached_up"]
	if m == nil {
		t.Fatal("expected metric 'memcached_up'")
	}
	if m.Description != "Could the memcached server be reached." {
		t.Errorf("unexpected description: %q", m.Description)
	}
	if m.InstrumentType != "gauge" {
		t.Errorf("expected gauge, got %q", m.InstrumentType)
	}

	cmd := names["memcached_commands_total"]
	if cmd == nil {
		t.Fatal("expected metric 'memcached_commands_total'")
	}
	if cmd.InstrumentType != "counter" {
		t.Errorf("expected counter, got %q", cmd.InstrumentType)
	}
	if len(cmd.Attributes) != 2 {
		t.Fatalf("expected 2 attributes, got %d", len(cmd.Attributes))
	}

	crawler := names["memcached_lru_crawler_enabled"]
	if crawler == nil {
		t.Fatal("expected metric 'memcached_lru_crawler_enabled'")
	}

	slab := names["memcached_slab_lru_hits_total"]
	if slab == nil {
		t.Fatal("expected metric 'memcached_slab_lru_hits_total'")
	}
	if len(slab.Attributes) != 2 {
		t.Fatalf("expected 2 attributes, got %d", len(slab.Attributes))
	}
}

func TestMemcachedAdapter_Extract_SkipsTestFiles(t *testing.T) {
	tmpDir := setupTestRepo(t)

	testFile := `package exporter

import "github.com/prometheus/client_golang/prometheus"

const Namespace = "memcached"

var testMetric = prometheus.NewDesc(
	prometheus.BuildFQName(Namespace, "", "test"),
	"Test metric.",
	nil, nil,
)
`
	writeFile(t, filepath.Join(tmpDir, "pkg", "exporter", "exporter_test.go"), testFile)

	a := NewAdapter("/tmp/cache")
	metrics, err := a.Extract(context.Background(), &adapter.FetchResult{
		RepoPath: tmpDir,
		Commit:   "abc123",
	})
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 0 {
		t.Errorf("expected 0 metrics, got %d", len(metrics))
	}
}

func TestMemcachedAdapter_Extract_EmptyRepo(t *testing.T) {
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
		t.Errorf("expected 0 metrics, got %d", len(metrics))
	}
}

func TestInferInstrumentType(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"memcached_commands_total", "counter"},
		{"memcached_up", "gauge"},
		{"memcached_slab_lru_hits_total", "counter"},
		{"memcached_lru_crawler_enabled", "gauge"},
		{"request_duration_seconds_bucket", "histogram"},
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

// Test helpers

func setupTestRepo(t *testing.T) string {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "memcached-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tmpDir) })

	if err := os.MkdirAll(filepath.Join(tmpDir, "pkg", "exporter"), 0750); err != nil {
		t.Fatalf("failed to create pkg/exporter dir: %v", err)
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
