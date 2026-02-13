package nats

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
)

func TestNATSAdapter_Name(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Name() != "prometheus-nats" {
		t.Errorf("expected name 'prometheus-nats', got %q", a.Name())
	}
}

func TestNATSAdapter_SourceCategory(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.SourceCategory() != domain.SourcePrometheus {
		t.Errorf("expected source category 'prometheus', got %q", a.SourceCategory())
	}
}

func TestNATSAdapter_Confidence(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Confidence() != domain.ConfidenceDerived {
		t.Errorf("expected confidence 'derived', got %q", a.Confidence())
	}
}

func TestNATSAdapter_ExtractionMethod(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.ExtractionMethod() != domain.ExtractionAST {
		t.Errorf("expected extraction method 'ast', got %q", a.ExtractionMethod())
	}
}

func TestNATSAdapter_RepoURL(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.RepoURL() != "https://github.com/nats-io/prometheus-nats-exporter" {
		t.Errorf("unexpected repo URL: %q", a.RepoURL())
	}
}

func TestNATSAdapter_ImplementsAdapter(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	var _ adapter.Adapter = a
}

func TestNATSAdapter_Extract_CoreMetrics(t *testing.T) {
	tmpDir := setupTestRepo(t)

	// Simulates connz.go — uses "system" function param and "connzEndpoint" var
	goFile := `package collector

import "github.com/prometheus/client_golang/prometheus"

var connzEndpoint = "connz"

func createConnzCollector(system string) *connzCollector {
	summaryLabels := []string{"server_id"}
	return &connzCollector{
		numConnections: prometheus.NewDesc(
			prometheus.BuildFQName(system, connzEndpoint, "num_connections"),
			"num_connections",
			summaryLabels, nil,
		),
		total: prometheus.NewDesc(
			prometheus.BuildFQName(system, connzEndpoint, "total"),
			"total",
			summaryLabels, nil,
		),
	}
}
`
	writeFile(t, filepath.Join(tmpDir, "collector", "connz.go"), goFile)

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

	// Core system metrics get "gnatsd_" prefix
	m := names["gnatsd_connz_num_connections"]
	if m == nil {
		t.Fatal("expected metric 'gnatsd_connz_num_connections'")
	}
	if m.Description != "num_connections" {
		t.Errorf("unexpected description: %q", m.Description)
	}
	if len(m.Attributes) != 1 || m.Attributes[0].Name != "server_id" {
		t.Errorf("expected 1 attribute 'server_id', got %v", m.Attributes)
	}

	if names["gnatsd_connz_total"] == nil {
		t.Error("expected metric 'gnatsd_connz_total'")
	}
}

func TestNATSAdapter_Extract_JetStreamMetrics(t *testing.T) {
	tmpDir := setupTestRepo(t)

	// Simulates jsz.go — JetStream metrics
	goFile := `package collector

import "github.com/prometheus/client_golang/prometheus"

func newJszCollector(system, endpoint string) *jszCollector {
	serverLabels := []string{"server_id"}
	return &jszCollector{
		totalStreams: prometheus.NewDesc(
			prometheus.BuildFQName(system, "server", "total_streams"),
			"Total number of streams",
			serverLabels, nil,
		),
		streamMessages: prometheus.NewDesc(
			prometheus.BuildFQName(system, "stream", "total_messages"),
			"Total number of messages from a stream",
			serverLabels, nil,
		),
	}
}
`
	writeFile(t, filepath.Join(tmpDir, "collector", "jsz.go"), goFile)

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

	// JetStream metrics get "jetstream_" prefix
	if names["jetstream_server_total_streams"] == nil {
		t.Error("expected metric 'jetstream_server_total_streams'")
	}
	if names["jetstream_stream_total_messages"] == nil {
		t.Error("expected metric 'jetstream_stream_total_messages'")
	}
}

func TestNATSAdapter_Extract_SkipsTestFiles(t *testing.T) {
	tmpDir := setupTestRepo(t)

	testFile := `package collector

import "github.com/prometheus/client_golang/prometheus"

var testMetric = prometheus.NewDesc(
	prometheus.BuildFQName("gnatsd", "test", "metric"),
	"Test metric.",
	nil, nil,
)
`
	writeFile(t, filepath.Join(tmpDir, "collector", "connz_test.go"), testFile)

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

func TestNATSAdapter_Extract_EmptyRepo(t *testing.T) {
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

func TestSystemPrefix(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"connz.go", "gnatsd"},
		{"healthz.go", "gnatsd"},
		{"accountz.go", "gnatsd"},
		{"accstatz.go", "gnatsd"},
		{"gatewayz.go", "gnatsd"},
		{"leafz.go", "gnatsd"},
		{"jsz.go", "jetstream"},
		{"collector.go", "gnatsd"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got := systemPrefix(tt.filename)
			if got != tt.expected {
				t.Errorf("systemPrefix(%q) = %q, want %q", tt.filename, got, tt.expected)
			}
		})
	}
}

func TestInferInstrumentType(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"gnatsd_connz_in_msgs_total", "counter"},
		{"gnatsd_connz_num_connections", "gauge"},
		{"jetstream_stream_total_messages", "gauge"},
		{"gnatsd_healthz_status", "gauge"},
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
	tmpDir, err := os.MkdirTemp("", "nats-adapter-test-*")
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
