package kafka

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
)

func TestKafkaAdapter_Name(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Name() != "prometheus-kafka" {
		t.Errorf("expected name 'prometheus-kafka', got %q", a.Name())
	}
}

func TestKafkaAdapter_SourceCategory(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.SourceCategory() != domain.SourcePrometheus {
		t.Errorf("expected source category 'prometheus', got %q", a.SourceCategory())
	}
}

func TestKafkaAdapter_Confidence(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Confidence() != domain.ConfidenceDerived {
		t.Errorf("expected confidence 'derived', got %q", a.Confidence())
	}
}

func TestKafkaAdapter_ExtractionMethod(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.ExtractionMethod() != domain.ExtractionAST {
		t.Errorf("expected extraction method 'ast', got %q", a.ExtractionMethod())
	}
}

func TestKafkaAdapter_RepoURL(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.RepoURL() != "https://github.com/danielqsj/kafka_exporter" {
		t.Errorf("unexpected repo URL: %q", a.RepoURL())
	}
}

func TestKafkaAdapter_ImplementsAdapter(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	var _ adapter.Adapter = a
}

func TestKafkaAdapter_Extract(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kafka-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	goFile := `package main

import "github.com/prometheus/client_golang/prometheus"

var (
	clusterBrokers = prometheus.NewDesc(
		"kafka_brokers",
		"Number of Brokers in the Kafka Cluster.",
		nil, nil,
	)
	topicPartitions = prometheus.NewDesc(
		"kafka_topic_partitions",
		"Number of partitions for this Topic",
		[]string{"topic"}, nil,
	)
	consumergroupLag = prometheus.NewDesc(
		"kafka_consumergroup_lag",
		"Current Approximate Lag of a ConsumerGroup at Topic/Partition",
		[]string{"consumergroup", "topic", "partition"}, nil,
	)
)
`
	if err := os.WriteFile(filepath.Join(tmpDir, "kafka_exporter.go"), []byte(goFile), 0600); err != nil {
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

	if _, ok := names["kafka_brokers"]; !ok {
		t.Error("expected metric 'kafka_brokers'")
	}
	if _, ok := names["kafka_topic_partitions"]; !ok {
		t.Error("expected metric 'kafka_topic_partitions'")
	}
	if _, ok := names["kafka_consumergroup_lag"]; !ok {
		t.Error("expected metric 'kafka_consumergroup_lag'")
	}

	m := names["kafka_consumergroup_lag"]
	if m.Description != "Current Approximate Lag of a ConsumerGroup at Topic/Partition" {
		t.Errorf("unexpected description: %q", m.Description)
	}
	if len(m.Attributes) != 3 {
		t.Errorf("expected 3 attributes, got %d", len(m.Attributes))
	}
	if m.ComponentName != "kafka" {
		t.Errorf("expected component name 'kafka', got %q", m.ComponentName)
	}
}

func TestKafkaAdapter_Extract_SkipsTestFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kafka-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	testFile := `package main

import "github.com/prometheus/client_golang/prometheus"

var testDesc = prometheus.NewDesc(
	"kafka_test_metric",
	"Test metric.",
	nil, nil,
)
`
	if err := os.WriteFile(filepath.Join(tmpDir, "simple_test.go"), []byte(testFile), 0600); err != nil {
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

func TestKafkaAdapter_Extract_EmptyRepo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kafka-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

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
		{"kafka_brokers", "gauge"},
		{"kafka_topic_partitions", "gauge"},
		{"kafka_consumergroup_lag", "gauge"},
		{"kafka_consumergroup_current_offset_sum", "gauge"},
		{"kafka_messages_total", "counter"},
		{"request_duration_seconds_bucket", "histogram"},
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
		{"kafka_exporter.go", "kafka"},
		{"scram_client.go", "scram_client"},
		{"main.go", "main"},
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
