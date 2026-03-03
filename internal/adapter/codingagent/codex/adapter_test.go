package codex

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
	if a.Name() != "codingagent-codex" {
		t.Errorf("expected name 'codingagent-codex', got %q", a.Name())
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
	if a.RepoURL() != "https://github.com/openai/codex" {
		t.Errorf("unexpected repo URL: %q", a.RepoURL())
	}
}

func TestAdapter_ImplementsAdapter(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	var _ adapter.Adapter = a
}

func TestAdapter_Extract(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "codex-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	metricsDir := filepath.Join(tmpDir, "codex-rs", "otel", "src", "metrics")
	if err := os.MkdirAll(metricsDir, 0750); err != nil {
		t.Fatalf("failed to create metrics dir: %v", err)
	}

	rustSource := `//! Metric name constants for OpenTelemetry metrics.

pub(crate) const TOOL_CALL_COUNT_METRIC: &str = "codex.tool.call";
pub(crate) const TOOL_CALL_DURATION_METRIC: &str = "codex.tool.call.duration_ms";
pub(crate) const API_CALL_COUNT_METRIC: &str = "codex.api_request";
pub(crate) const API_CALL_DURATION_METRIC: &str = "codex.api_request.duration_ms";
pub(crate) const SSE_EVENT_COUNT_METRIC: &str = "codex.sse_event";
pub(crate) const SSE_EVENT_DURATION_METRIC: &str = "codex.sse_event.duration_ms";
pub const WEBSOCKET_REQUEST_COUNT_METRIC: &str = "codex.websocket.request";
pub const WEBSOCKET_REQUEST_DURATION_METRIC: &str = "codex.websocket.request.duration_ms";
`
	if err := os.WriteFile(filepath.Join(metricsDir, "names.rs"), []byte(rustSource), 0600); err != nil {
		t.Fatalf("failed to write rust file: %v", err)
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

	if len(metrics) != 8 {
		t.Fatalf("expected 8 metrics, got %d", len(metrics))
	}

	names := make(map[string]*adapter.RawMetric)
	for _, m := range metrics {
		names[m.Name] = m
	}

	// Check counters (both pub(crate) and pub forms)
	for _, name := range []string{"codex.tool.call", "codex.api_request", "codex.sse_event", "codex.websocket.request"} {
		m, ok := names[name]
		if !ok {
			t.Errorf("missing metric %q", name)
			continue
		}
		if m.InstrumentType != "counter" {
			t.Errorf("metric %q: expected counter, got %q", name, m.InstrumentType)
		}
	}

	// Check histograms (duration_ms suffix)
	for _, name := range []string{"codex.tool.call.duration_ms", "codex.api_request.duration_ms", "codex.sse_event.duration_ms", "codex.websocket.request.duration_ms"} {
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

	// Check component metadata
	for _, m := range metrics {
		if m.ComponentName != "codex" {
			t.Errorf("metric %q: expected component name 'codex', got %q", m.Name, m.ComponentName)
		}
		if m.ComponentType != string(domain.ComponentPlatform) {
			t.Errorf("metric %q: expected component type 'platform', got %q", m.Name, m.ComponentType)
		}
	}
}

func TestAdapter_Extract_EmptyFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "codex-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	metricsDir := filepath.Join(tmpDir, "codex-rs", "otel", "src", "metrics")
	if err := os.MkdirAll(metricsDir, 0750); err != nil {
		t.Fatalf("failed to create metrics dir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(metricsDir, "names.rs"), []byte("// empty\n"), 0600); err != nil {
		t.Fatalf("failed to write rust file: %v", err)
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
		t.Errorf("expected 0 metrics for empty file, got %d", len(metrics))
	}
}

func TestAdapter_Extract_MissingFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "codex-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	a := NewAdapter("/tmp/cache")
	result := &adapter.FetchResult{
		RepoPath: tmpDir,
		Commit:   "abc123",
	}

	_, err = a.Extract(context.Background(), result)
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestInferType(t *testing.T) {
	tests := []struct {
		name         string
		expectedType string
		expectedUnit string
	}{
		{"codex.session.duration_ms", "histogram", "ms"},
		{"codex.tool.call", "counter", "count"},
		{"codex.api.request.duration_ms", "histogram", "ms"},
		{"codex.token.input", "counter", "count"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotUnit := inferType(tt.name)
			if gotType != tt.expectedType {
				t.Errorf("inferType(%q) type = %q, want %q", tt.name, gotType, tt.expectedType)
			}
			if gotUnit != tt.expectedUnit {
				t.Errorf("inferType(%q) unit = %q, want %q", tt.name, gotUnit, tt.expectedUnit)
			}
		})
	}
}
