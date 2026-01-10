package mongodb

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/base14/otel-glossary/internal/adapter"
	"github.com/base14/otel-glossary/internal/domain"
)

func TestMongoDBAdapter_Name(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Name() != "prometheus-mongodb" {
		t.Errorf("expected name 'prometheus-mongodb', got %q", a.Name())
	}
}

func TestMongoDBAdapter_SourceCategory(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.SourceCategory() != domain.SourcePrometheus {
		t.Errorf("expected source category 'prometheus', got %q", a.SourceCategory())
	}
}

func TestMongoDBAdapter_Confidence(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Confidence() != domain.ConfidenceDerived {
		t.Errorf("expected confidence 'derived', got %q", a.Confidence())
	}
}

func TestMongoDBAdapter_ExtractionMethod(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.ExtractionMethod() != domain.ExtractionAST {
		t.Errorf("expected extraction method 'ast', got %q", a.ExtractionMethod())
	}
}

func TestMongoDBAdapter_RepoURL(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.RepoURL() != "https://github.com/percona/mongodb_exporter" {
		t.Errorf("unexpected repo URL: %q", a.RepoURL())
	}
}

func TestMongoDBAdapter_ImplementsAdapter(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	var _ adapter.Adapter = a
}

func TestMongoDBAdapter_Extract(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mongodb-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	exporterDir := filepath.Join(tmpDir, "exporter")
	if err := os.MkdirAll(exporterDir, 0750); err != nil {
		t.Fatalf("failed to create exporter dir: %v", err)
	}

	goFile := `package exporter

import "github.com/prometheus/client_golang/prometheus"

var (
	mongodbUp = prometheus.NewDesc(
		"mongodb_up",
		"Whether MongoDB is up.",
		nil, nil,
	)
	connectionsDesc = prometheus.NewDesc(
		"mongodb_connections_total",
		"Total number of MongoDB connections.",
		[]string{"state"}, nil,
	)
)
`
	if err := os.WriteFile(filepath.Join(exporterDir, "general_collector.go"), []byte(goFile), 0600); err != nil {
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

	if _, ok := names["mongodb_up"]; !ok {
		t.Error("expected metric 'mongodb_up'")
	}
	if _, ok := names["mongodb_connections_total"]; !ok {
		t.Error("expected metric 'mongodb_connections_total'")
	}

	m := names["mongodb_connections_total"]
	if m.Description != "Total number of MongoDB connections." {
		t.Errorf("unexpected description: %q", m.Description)
	}
	if len(m.Attributes) != 1 || m.Attributes[0].Name != "state" {
		t.Errorf("expected 1 attribute 'state', got %v", m.Attributes)
	}
	if m.InstrumentType != "counter" {
		t.Errorf("expected instrument type 'counter', got %q", m.InstrumentType)
	}
	if m.ComponentName != "general" {
		t.Errorf("expected component name 'general', got %q", m.ComponentName)
	}
}

func TestMongoDBAdapter_Extract_SkipsTestFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mongodb-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	exporterDir := filepath.Join(tmpDir, "exporter")
	if err := os.MkdirAll(exporterDir, 0750); err != nil {
		t.Fatalf("failed to create exporter dir: %v", err)
	}

	testFile := `package exporter

import "github.com/prometheus/client_golang/prometheus"

var testDesc = prometheus.NewDesc(
	"mongodb_test_metric",
	"Test metric.",
	nil, nil,
)
`
	if err := os.WriteFile(filepath.Join(exporterDir, "test_collector_test.go"), []byte(testFile), 0600); err != nil {
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

func TestMongoDBAdapter_Extract_EmptyRepo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mongodb-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	exporterDir := filepath.Join(tmpDir, "exporter")
	if err := os.MkdirAll(exporterDir, 0750); err != nil {
		t.Fatalf("failed to create exporter dir: %v", err)
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
		{"mongodb_connections_total", "counter"},
		{"mongodb_up", "gauge"},
		{"mongodb_memory_bytes", "gauge"},
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
		{"general_collector.go", "general"},
		{"collstats_collector.go", "collstats"},
		{"exporter.go", "exporter"},
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
