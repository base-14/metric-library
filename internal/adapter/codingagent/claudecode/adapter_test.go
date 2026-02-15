package claudecode

import (
	"context"
	"testing"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
)

func TestAdapter_Name(t *testing.T) {
	a := NewAdapter("")
	if a.Name() != "codingagent-claude-code" {
		t.Errorf("expected name 'codingagent-claude-code', got %q", a.Name())
	}
}

func TestAdapter_SourceCategory(t *testing.T) {
	a := NewAdapter("")
	if a.SourceCategory() != domain.SourceCodingAgent {
		t.Errorf("expected source category 'codingagent', got %q", a.SourceCategory())
	}
}

func TestAdapter_Confidence(t *testing.T) {
	a := NewAdapter("")
	if a.Confidence() != domain.ConfidenceDocumented {
		t.Errorf("expected confidence 'documented', got %q", a.Confidence())
	}
}

func TestAdapter_ExtractionMethod(t *testing.T) {
	a := NewAdapter("")
	if a.ExtractionMethod() != domain.ExtractionMetadata {
		t.Errorf("expected extraction method 'metadata', got %q", a.ExtractionMethod())
	}
}

func TestAdapter_RepoURL(t *testing.T) {
	a := NewAdapter("")
	expected := "https://github.com/anthropics/claude-code-monitoring-guide"
	if a.RepoURL() != expected {
		t.Errorf("expected repo URL %q, got %q", expected, a.RepoURL())
	}
}

func TestAdapter_ImplementsAdapter(t *testing.T) {
	a := NewAdapter("")
	var _ adapter.Adapter = a
}

func TestAdapter_Fetch(t *testing.T) {
	a := NewAdapter("")
	result, err := a.Fetch(context.Background(), adapter.FetchOptions{})
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}
	if result.Commit == "" {
		t.Error("expected non-empty commit")
	}
}

func TestAdapter_Extract(t *testing.T) {
	a := NewAdapter("")
	fetchResult, _ := a.Fetch(context.Background(), adapter.FetchOptions{})
	metrics, err := a.Extract(context.Background(), fetchResult)
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

	expectedMetrics := []struct {
		name           string
		instrumentType string
		unit           string
		attrCount      int
	}{
		{"claude_code.session.count", "counter", "count", 0},
		{"claude_code.lines_of_code.count", "counter", "count", 1},
		{"claude_code.pull_request.count", "counter", "count", 0},
		{"claude_code.commit.count", "counter", "count", 0},
		{"claude_code.cost.usage", "counter", "USD", 1},
		{"claude_code.token.usage", "counter", "tokens", 2},
		{"claude_code.code_edit_tool.decision", "counter", "count", 3},
		{"claude_code.active_time.total", "counter", "s", 0},
	}

	for _, exp := range expectedMetrics {
		m, ok := names[exp.name]
		if !ok {
			t.Errorf("missing metric %q", exp.name)
			continue
		}
		if m.InstrumentType != exp.instrumentType {
			t.Errorf("metric %q: expected instrument type %q, got %q", exp.name, exp.instrumentType, m.InstrumentType)
		}
		if m.Unit != exp.unit {
			t.Errorf("metric %q: expected unit %q, got %q", exp.name, exp.unit, m.Unit)
		}
		if len(m.Attributes) != exp.attrCount {
			t.Errorf("metric %q: expected %d attributes, got %d", exp.name, exp.attrCount, len(m.Attributes))
		}
		if m.ComponentName != "claude-code" {
			t.Errorf("metric %q: expected component name 'claude-code', got %q", exp.name, m.ComponentName)
		}
		if m.ComponentType != string(domain.ComponentPlatform) {
			t.Errorf("metric %q: expected component type 'platform', got %q", exp.name, m.ComponentType)
		}
	}
}

func TestAdapter_Extract_TokenUsageAttributes(t *testing.T) {
	a := NewAdapter("")
	fetchResult, _ := a.Fetch(context.Background(), adapter.FetchOptions{})
	metrics, err := a.Extract(context.Background(), fetchResult)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	var tokenMetric *adapter.RawMetric
	for _, m := range metrics {
		if m.Name == "claude_code.token.usage" {
			tokenMetric = m
			break
		}
	}

	if tokenMetric == nil {
		t.Fatal("token.usage metric not found")
	}

	attrNames := make(map[string]bool)
	for _, attr := range tokenMetric.Attributes {
		attrNames[attr.Name] = true
	}

	if !attrNames["type"] {
		t.Error("expected 'type' attribute on token.usage")
	}
	if !attrNames["model"] {
		t.Error("expected 'model' attribute on token.usage")
	}
}
