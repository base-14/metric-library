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

	tsSource := `import { ValueType } from '@opentelemetry/api';

const TOOL_CALL_COUNT = 'gemini_cli.tool.call.count';
const TOOL_CALL_LATENCY = 'gemini_cli.tool.call.latency';
const API_REQUEST_COUNT = 'gemini_cli.api.request.count';
const TOKEN_USAGE = 'gemini_cli.token.usage';
const SESSION_COUNT = 'gemini_cli.session.count';

const COUNTER_DEFINITIONS = {
  [TOOL_CALL_COUNT]: {
    description: 'Counts tool calls, tagged by function name and success.',
    valueType: ValueType.INT,
    assign: (c: Counter) => (toolCallCounter = c),
    attributes: {} as { function_name: string; success: boolean; },
  },
  [API_REQUEST_COUNT]: {
    description: 'Counts API requests, tagged by model and status.',
    valueType: ValueType.INT,
    assign: (c: Counter) => (apiRequestCounter = c),
    attributes: {} as { model: string; },
  },
  [TOKEN_USAGE]: {
    description: 'Counts the total number of tokens used.',
    valueType: ValueType.INT,
    assign: (c: Counter) => (tokenUsageCounter = c),
    attributes: {} as { model: string; type: string; },
  },
  [SESSION_COUNT]: {
    description: 'Count of CLI sessions started.',
    valueType: ValueType.INT,
    assign: (c: Counter) => (sessionCounter = c),
    attributes: {} as Record<string, never>,
  },
} as const;

const HISTOGRAM_DEFINITIONS = {
  [TOOL_CALL_LATENCY]: {
    description: 'Latency of tool calls in milliseconds.',
    unit: 'ms',
    valueType: ValueType.INT,
    assign: (h: Histogram) => (toolCallLatencyHistogram = h),
    attributes: {} as { function_name: string; },
  },
} as const;
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
	for _, name := range []string{"gemini_cli.tool.call.count", "gemini_cli.api.request.count", "gemini_cli.token.usage", "gemini_cli.session.count"} {
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
	for _, name := range []string{"gemini_cli.tool.call.latency"} {
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
	m := names["gemini_cli.tool.call.count"]
	if m.Description != "Counts tool calls, tagged by function name and success." {
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

func TestAdapter_Extract_PerformanceMetrics(t *testing.T) {
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
const STARTUP_TIME = 'gemini_cli.startup.duration';
const REGRESSION_DETECTION = 'gemini_cli.performance.regression';

const PERFORMANCE_COUNTER_DEFINITIONS = {
  [REGRESSION_DETECTION]: {
    description: 'Performance regression detection events.',
    valueType: ValueType.INT,
    assign: (c: Counter) => (regressionDetectionCounter = c),
  },
} as const;

const PERFORMANCE_HISTOGRAM_DEFINITIONS = {
  [STARTUP_TIME]: {
    description: 'CLI startup time in milliseconds.',
    unit: 'ms',
    valueType: ValueType.DOUBLE,
    assign: (h: Histogram) => (startupTimeHistogram = h),
  },
} as const;
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

	names := make(map[string]*adapter.RawMetric)
	for _, m := range metrics {
		names[m.Name] = m
	}

	if m, ok := names["gemini_cli.performance.regression"]; !ok {
		t.Error("missing metric gemini_cli.performance.regression")
	} else if m.InstrumentType != "counter" {
		t.Errorf("expected counter, got %q", m.InstrumentType)
	}

	if m, ok := names["gemini_cli.startup.duration"]; !ok {
		t.Error("missing metric gemini_cli.startup.duration")
	} else if m.InstrumentType != "histogram" {
		t.Errorf("expected histogram, got %q", m.InstrumentType)
	}
}

func TestAdapter_Extract_MultilineConst(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gemini-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	metricsDir := filepath.Join(tmpDir, "packages", "core", "src", "telemetry")
	if err := os.MkdirAll(metricsDir, 0750); err != nil {
		t.Fatalf("failed to create metrics dir: %v", err)
	}

	// Test multiline const declarations like:
	// const CONTENT_RETRY_FAILURE_COUNT =
	//   'gemini_cli.chat.content_retry_failure.count';
	tsSource := `
const CONTENT_RETRY_FAILURE_COUNT =
  'gemini_cli.chat.content_retry_failure.count';

const COUNTER_DEFINITIONS = {
  [CONTENT_RETRY_FAILURE_COUNT]: {
    description: 'Counts occurrences of all content retries failing.',
    valueType: ValueType.INT,
    assign: (c: Counter) => (contentRetryFailureCounter = c),
  },
} as const;
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

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	if metrics[0].Name != "gemini_cli.chat.content_retry_failure.count" {
		t.Errorf("unexpected name: %q", metrics[0].Name)
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
