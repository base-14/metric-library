package gemini

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
)

func TestAdapter_Name(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Name() != "codingagent-gemini" {
		t.Errorf("expected name 'codingagent-gemini', got %q", a.Name())
	}
}

func TestAdapter_SourceCategory(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.SourceCategory() != domain.SourceCodingAgent {
		t.Errorf("expected source category 'codingagent', got %q", a.SourceCategory())
	}
}

func TestAdapter_Confidence(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Confidence() != domain.ConfidenceAuthoritative {
		t.Errorf("expected confidence 'authoritative', got %q", a.Confidence())
	}
}

func TestAdapter_ExtractionMethod(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.ExtractionMethod() != domain.ExtractionAST {
		t.Errorf("expected extraction method 'ast', got %q", a.ExtractionMethod())
	}
}

func TestAdapter_RepoURL(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.RepoURL() != "https://github.com/google-gemini/gemini-cli" {
		t.Errorf("unexpected repo URL: %q", a.RepoURL())
	}
}

func TestAdapter_ImplementsAdapter(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	var _ adapter.Adapter = a
}

func TestAdapter_Extract(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gemini-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	metricsDir := filepath.Join(tmpDir, "packages", "core", "src", "telemetry")
	if err := os.MkdirAll(metricsDir, 0750); err != nil {
		t.Fatalf("failed to create metrics dir: %v", err)
	}

	tsSource := `import { Meter } from '@opentelemetry/api';

export function registerMetrics(meter: Meter) {
  meter.createCounter('gemini_cli.session.count', {
    description: 'Number of CLI sessions started',
    unit: 'count',
  });

  meter.createCounter('gemini_cli.token.usage', {
    description: 'Number of tokens consumed',
    unit: 'tokens',
  });

  meter.createHistogram('gemini_cli.tool.call.latency', {
    description: 'Latency of tool calls',
    unit: 'ms',
  });

  meter.createHistogram('gemini_cli.api.request.duration', {
    description: 'Duration of API requests',
    unit: 'ms',
  });

  meter.createCounter('gemini_cli.error.count', {
    description: 'Number of errors encountered',
    unit: 'count',
  });
}
`
	if err := os.WriteFile(filepath.Join(metricsDir, "metrics.ts"), []byte(tsSource), 0600); err != nil {
		t.Fatalf("failed to write ts file: %v", err)
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

	if len(metrics) != 5 {
		t.Fatalf("expected 5 metrics, got %d", len(metrics))
	}

	names := make(map[string]*adapter.RawMetric)
	for _, m := range metrics {
		names[m.Name] = m
	}

	// Check counters
	for _, name := range []string{"gemini_cli.session.count", "gemini_cli.token.usage", "gemini_cli.error.count"} {
		m, ok := names[name]
		if !ok {
			t.Errorf("missing metric %q", name)
			continue
		}
		if m.InstrumentType != "counter" {
			t.Errorf("metric %q: expected counter, got %q", name, m.InstrumentType)
		}
	}

	// Check histograms
	for _, name := range []string{"gemini_cli.tool.call.latency", "gemini_cli.api.request.duration"} {
		m, ok := names[name]
		if !ok {
			t.Errorf("missing metric %q", name)
			continue
		}
		if m.InstrumentType != "histogram" {
			t.Errorf("metric %q: expected histogram, got %q", name, m.InstrumentType)
		}
		if m.Unit != "ms" {
			t.Errorf("metric %q: expected unit 'ms', got %q", name, m.Unit)
		}
	}

	// Check descriptions
	m := names["gemini_cli.session.count"]
	if m.Description != "Number of CLI sessions started" {
		t.Errorf("unexpected description: %q", m.Description)
	}

	// Check component metadata
	for _, m := range metrics {
		if m.ComponentName != "gemini-cli" {
			t.Errorf("metric %q: expected component name 'gemini-cli', got %q", m.Name, m.ComponentName)
		}
		if m.ComponentType != string(domain.ComponentPlatform) {
			t.Errorf("metric %q: expected component type 'platform', got %q", m.Name, m.ComponentType)
		}
	}
}

func TestAdapter_Extract_SingleQuotesAndDoubleQuotes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gemini-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	metricsDir := filepath.Join(tmpDir, "packages", "core", "src", "telemetry")
	if err := os.MkdirAll(metricsDir, 0750); err != nil {
		t.Fatalf("failed to create metrics dir: %v", err)
	}

	tsSource := `
  meter.createCounter('gemini_cli.single_quoted', {
    description: 'Single quoted metric',
    unit: 'count',
  });

  meter.createCounter("gemini_cli.double_quoted", {
    description: "Double quoted metric",
    unit: "count",
  });
`
	if err := os.WriteFile(filepath.Join(metricsDir, "metrics.ts"), []byte(tsSource), 0600); err != nil {
		t.Fatalf("failed to write ts file: %v", err)
	}

	a := NewAdapter("/tmp/cache")
	result := &adapter.FetchResult{RepoPath: tmpDir, Commit: "abc123"}

	metrics, err := a.Extract(context.Background(), result)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 2 {
		t.Fatalf("expected 2 metrics, got %d", len(metrics))
	}
}

func TestAdapter_Extract_EmptyFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gemini-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	metricsDir := filepath.Join(tmpDir, "packages", "core", "src", "telemetry")
	if err := os.MkdirAll(metricsDir, 0750); err != nil {
		t.Fatalf("failed to create metrics dir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(metricsDir, "metrics.ts"), []byte("// empty\n"), 0600); err != nil {
		t.Fatalf("failed to write ts file: %v", err)
	}

	a := NewAdapter("/tmp/cache")
	result := &adapter.FetchResult{RepoPath: tmpDir, Commit: "abc123"}

	metrics, err := a.Extract(context.Background(), result)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 0 {
		t.Errorf("expected 0 metrics for empty file, got %d", len(metrics))
	}
}

func TestAdapter_Extract_MissingFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gemini-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	a := NewAdapter("/tmp/cache")
	result := &adapter.FetchResult{RepoPath: tmpDir, Commit: "abc123"}

	_, err = a.Extract(context.Background(), result)
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestAdapter_Extract_MultipleFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gemini-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	metricsDir := filepath.Join(tmpDir, "packages", "core", "src", "telemetry")
	if err := os.MkdirAll(metricsDir, 0750); err != nil {
		t.Fatalf("failed to create metrics dir: %v", err)
	}

	metricsTS := `
  meter.createCounter('gemini_cli.session.count', {
    description: 'Sessions',
    unit: 'count',
  });
`
	if err := os.WriteFile(filepath.Join(metricsDir, "metrics.ts"), []byte(metricsTS), 0600); err != nil {
		t.Fatalf("failed to write metrics.ts: %v", err)
	}

	otherTS := `
  meter.createCounter('gemini_cli.other.count', {
    description: 'Other metric',
    unit: 'count',
  });
`
	if err := os.WriteFile(filepath.Join(metricsDir, "other_metrics.ts"), []byte(otherTS), 0600); err != nil {
		t.Fatalf("failed to write other_metrics.ts: %v", err)
	}

	a := NewAdapter("/tmp/cache")
	result := &adapter.FetchResult{RepoPath: tmpDir, Commit: "abc123"}

	metrics, err := a.Extract(context.Background(), result)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 2 {
		t.Fatalf("expected 2 metrics from both files, got %d", len(metrics))
	}
}
