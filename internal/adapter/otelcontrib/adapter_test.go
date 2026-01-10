package otelcontrib

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/base14/otel-glossary/internal/adapter"
	"github.com/base14/otel-glossary/internal/domain"
)

func TestOTelContribAdapter_Name(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Name() != "otel-collector-contrib" {
		t.Errorf("expected name 'otel-collector-contrib', got %q", a.Name())
	}
}

func TestOTelContribAdapter_SourceCategory(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.SourceCategory() != domain.SourceOTEL {
		t.Errorf("expected source category 'otel', got %q", a.SourceCategory())
	}
}

func TestOTelContribAdapter_Confidence(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Confidence() != domain.ConfidenceAuthoritative {
		t.Errorf("expected confidence 'authoritative', got %q", a.Confidence())
	}
}

func TestOTelContribAdapter_ExtractionMethod(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.ExtractionMethod() != domain.ExtractionMetadata {
		t.Errorf("expected extraction method 'metadata', got %q", a.ExtractionMethod())
	}
}

func TestOTelContribAdapter_RepoURL(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.RepoURL() != "https://github.com/open-telemetry/opentelemetry-collector-contrib" {
		t.Errorf("unexpected repo URL: %q", a.RepoURL())
	}
}

func TestOTelContribAdapter_ImplementsAdapter(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	var _ adapter.Adapter = a
}

func TestOTelContribAdapter_Extract(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "otelcontrib-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	structure := map[string]string{
		"receiver/mysqlreceiver/metadata.yaml": `
type: mysqlreceiver

status:
  class: receiver

attributes:
  buffer_pool_data:
    description: The status of buffer pool data.
    type: string
    enum: [dirty, clean]

metrics:
  mysql.buffer_pool.pages:
    enabled: true
    description: The number of pages in the InnoDB buffer pool.
    unit: "{pages}"
    sum:
      value_type: int
      monotonic: false
      aggregation_temporality: cumulative
    attributes: [buffer_pool_data]
`,
		"processor/transformprocessor/metadata.yaml": `
type: transformprocessor

status:
  class: processor

metrics: {}
`,
	}

	for path, content := range structure {
		fullPath := filepath.Join(tmpDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0750); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0600); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
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

	if len(metrics) != 1 {
		t.Errorf("expected 1 metric, got %d", len(metrics))
	}

	if len(metrics) > 0 {
		m := metrics[0]
		if m.Name != "mysql.buffer_pool.pages" {
			t.Errorf("unexpected metric name: %q", m.Name)
		}
		if m.ComponentName != "mysqlreceiver" {
			t.Errorf("unexpected component name: %q", m.ComponentName)
		}
		if m.ComponentType != "receiver" {
			t.Errorf("unexpected component type: %q", m.ComponentType)
		}
	}
}

func TestOTelContribAdapter_Extract_EmptyRepo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "otelcontrib-test-*")
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
		t.Errorf("expected 0 metrics, got %d", len(metrics))
	}
}
