package cockroachdb

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
)

func TestCockroachDBAdapter_Name(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Name() != "prometheus-cockroachdb" {
		t.Errorf("expected name 'prometheus-cockroachdb', got %q", a.Name())
	}
}

func TestCockroachDBAdapter_SourceCategory(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.SourceCategory() != domain.SourcePrometheus {
		t.Errorf("expected source category 'prometheus', got %q", a.SourceCategory())
	}
}

func TestCockroachDBAdapter_Confidence(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Confidence() != domain.ConfidenceAuthoritative {
		t.Errorf("expected confidence 'authoritative', got %q", a.Confidence())
	}
}

func TestCockroachDBAdapter_ExtractionMethod(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.ExtractionMethod() != domain.ExtractionAST {
		t.Errorf("expected extraction method 'ast', got %q", a.ExtractionMethod())
	}
}

func TestCockroachDBAdapter_RepoURL(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.RepoURL() != "https://github.com/cockroachdb/cockroach" {
		t.Errorf("unexpected repo URL: %q", a.RepoURL())
	}
}

func TestCockroachDBAdapter_ImplementsAdapter(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	var _ adapter.Adapter = a
}

func TestCockroachDBAdapter_Extract_VarMetadata(t *testing.T) {
	tmpDir := setupTestRepo(t)

	goFile := `package kvserver

import "github.com/cockroachdb/cockroach/pkg/util/metric"

var metaReplicaCount = metric.Metadata{
	Name:        "replicas",
	Help:        "Number of replicas",
	Measurement: "Replicas",
	Unit:        metric.Unit_COUNT,
}

var metaRaftLeaderCount = metric.Metadata{
	Name:        "replicas.leaders",
	Help:        "Number of raft leaders",
	Measurement: "Raft Leaders",
	Unit:        metric.Unit_COUNT,
}
`
	writeFile(t, filepath.Join(tmpDir, "pkg", "kv", "kvserver", "metrics.go"), goFile)

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

	m := names["replicas"]
	if m == nil {
		t.Fatal("expected metric 'replicas'")
	}
	if m.Description != "Number of replicas" {
		t.Errorf("unexpected description: %q", m.Description)
	}
	if m.Unit != "count" {
		t.Errorf("unexpected unit: %q", m.Unit)
	}
	if m.InstrumentType != "gauge" {
		t.Errorf("unexpected instrument type: %q", m.InstrumentType)
	}
	if m.ComponentName != "cockroachdb" {
		t.Errorf("unexpected component name: %q", m.ComponentName)
	}

	if names["replicas.leaders"] == nil {
		t.Error("expected metric 'replicas.leaders'")
	}
}

func TestCockroachDBAdapter_Extract_InlineMetadata(t *testing.T) {
	tmpDir := setupTestRepo(t)

	goFile := `package rpc

import "github.com/cockroachdb/cockroach/pkg/util/metric"

type Metrics struct {
	ConnectionHealthy *metric.Gauge
}

func makeMetrics() Metrics {
	meta := metric.Metadata{
		Name: "rpc.connection.healthy",
		Help: "Gauge of current connections in a healthy state",
		Unit: metric.Unit_COUNT,
	}
	return Metrics{
		ConnectionHealthy: metric.NewGauge(meta),
	}
}
`
	writeFile(t, filepath.Join(tmpDir, "pkg", "rpc", "metrics.go"), goFile)

	a := NewAdapter("/tmp/cache")
	metrics, err := a.Extract(context.Background(), &adapter.FetchResult{
		RepoPath: tmpDir,
		Commit:   "abc123",
	})
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	m := metrics[0]
	if m.Name != "rpc.connection.healthy" {
		t.Errorf("expected name 'rpc.connection.healthy', got %q", m.Name)
	}
}

func TestCockroachDBAdapter_Extract_WithUnit(t *testing.T) {
	tmpDir := setupTestRepo(t)

	goFile := `package server

import "github.com/cockroachdb/cockroach/pkg/util/metric"

var metaSysMem = metric.Metadata{
	Name: "sys.rss",
	Help: "Current process RSS",
	Unit: metric.Unit_BYTES,
}

var metaLatency = metric.Metadata{
	Name: "exec.latency",
	Help: "Execution latency",
	Unit: metric.Unit_NANOSECONDS,
}

var metaCpuPercent = metric.Metadata{
	Name: "sys.cpu.combined.percent-normalized",
	Help: "Current combined cpu utilization",
	Unit: metric.Unit_PERCENT,
}
`
	writeFile(t, filepath.Join(tmpDir, "pkg", "server", "status.go"), goFile)

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

	if m := names["sys.rss"]; m == nil || m.Unit != "bytes" {
		t.Errorf("expected sys.rss with unit 'bytes', got %v", m)
	}
	if m := names["exec.latency"]; m == nil || m.Unit != "nanoseconds" {
		t.Errorf("expected exec.latency with unit 'nanoseconds', got %v", m)
	}
	if m := names["sys.cpu.combined.percent-normalized"]; m == nil || m.Unit != "percent" {
		t.Errorf("expected cpu metric with unit 'percent', got %v", m)
	}
}

func TestCockroachDBAdapter_Extract_SkipsTestFiles(t *testing.T) {
	tmpDir := setupTestRepo(t)

	testFile := `package kvserver

import "github.com/cockroachdb/cockroach/pkg/util/metric"

var metaTest = metric.Metadata{
	Name: "test.metric",
	Help: "Should not be extracted",
}
`
	writeFile(t, filepath.Join(tmpDir, "pkg", "kv", "metrics_test.go"), testFile)

	a := NewAdapter("/tmp/cache")
	metrics, err := a.Extract(context.Background(), &adapter.FetchResult{
		RepoPath: tmpDir,
		Commit:   "abc123",
	})
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 0 {
		t.Errorf("expected 0 metrics from test files, got %d", len(metrics))
	}
}

func TestCockroachDBAdapter_Extract_Deduplicates(t *testing.T) {
	tmpDir := setupTestRepo(t)

	// Same metric defined in two files
	goFile1 := `package pkg1

import "github.com/cockroachdb/cockroach/pkg/util/metric"

var meta1 = metric.Metadata{
	Name: "duplicate.metric",
	Help: "First occurrence",
}
`
	goFile2 := `package pkg2

import "github.com/cockroachdb/cockroach/pkg/util/metric"

var meta2 = metric.Metadata{
	Name: "duplicate.metric",
	Help: "Second occurrence",
}
`
	writeFile(t, filepath.Join(tmpDir, "pkg", "a", "metrics.go"), goFile1)
	writeFile(t, filepath.Join(tmpDir, "pkg", "b", "metrics.go"), goFile2)

	a := NewAdapter("/tmp/cache")
	metrics, err := a.Extract(context.Background(), &adapter.FetchResult{
		RepoPath: tmpDir,
		Commit:   "abc123",
	})
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 1 {
		t.Errorf("expected 1 deduplicated metric, got %d", len(metrics))
	}
}

func TestCockroachDBAdapter_Extract_EmptyRepo(t *testing.T) {
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

func TestParseMetricMetadata(t *testing.T) {
	src := `package test

import "github.com/cockroachdb/cockroach/pkg/util/metric"

var metaReplicaCount = metric.Metadata{
	Name:        "replicas",
	Help:        "Number of replicas",
	Measurement: "Replicas",
	Unit:        metric.Unit_COUNT,
}

var metaBytes = metric.Metadata{
	Name: "sys.host.disk.read.bytes",
	Help: "Bytes read from all disks since this process started",
	Unit: metric.Unit_BYTES,
}
`
	defs, err := parseMetricMetadata("test.go", []byte(src))
	if err != nil {
		t.Fatalf("parseMetricMetadata failed: %v", err)
	}

	if len(defs) != 2 {
		t.Fatalf("expected 2 defs, got %d", len(defs))
	}

	if defs[0].Name != "replicas" {
		t.Errorf("expected name 'replicas', got %q", defs[0].Name)
	}
	if defs[0].Help != "Number of replicas" {
		t.Errorf("unexpected help: %q", defs[0].Help)
	}
	if defs[0].Unit != "Unit_COUNT" {
		t.Errorf("unexpected unit: %q", defs[0].Unit)
	}

	if defs[1].Name != "sys.host.disk.read.bytes" {
		t.Errorf("expected name 'sys.host.disk.read.bytes', got %q", defs[1].Name)
	}
	if defs[1].Unit != "Unit_BYTES" {
		t.Errorf("unexpected unit: %q", defs[1].Unit)
	}
}

func TestMapUnit(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Unit_BYTES", "bytes"},
		{"Unit_COUNT", "count"},
		{"Unit_NANOSECONDS", "nanoseconds"},
		{"Unit_SECONDS", "seconds"},
		{"Unit_PERCENT", "percent"},
		{"Unit_TIMESTAMP_NS", "nanoseconds"},
		{"Unit_TIMESTAMP_SEC", "seconds"},
		{"Unit_CONST", ""},
		{"Unit_UNSET", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := mapUnit(tt.input)
			if got != tt.expected {
				t.Errorf("mapUnit(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// Test helpers

func setupTestRepo(t *testing.T) string {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "cockroachdb-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tmpDir) })

	if err := os.MkdirAll(filepath.Join(tmpDir, "pkg"), 0750); err != nil {
		t.Fatalf("failed to create pkg dir: %v", err)
	}
	return tmpDir
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		t.Fatalf("failed to create dir for %s: %v", path, err)
	}
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
