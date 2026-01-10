package python

import (
	"testing"

	adpt "github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
)

func TestAdapter_Name(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Name() != "otel-python" {
		t.Errorf("expected 'otel-python', got '%s'", a.Name())
	}
}

func TestAdapter_SourceCategory(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.SourceCategory() != domain.SourceOTEL {
		t.Errorf("expected SourceOTEL, got '%s'", a.SourceCategory())
	}
}

func TestAdapter_Confidence(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Confidence() != domain.ConfidenceDerived {
		t.Errorf("expected ConfidenceDerived, got '%s'", a.Confidence())
	}
}

func TestAdapter_ExtractionMethod(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.ExtractionMethod() != domain.ExtractionAST {
		t.Errorf("expected ExtractionAST, got '%s'", a.ExtractionMethod())
	}
}

func TestAdapter_RepoURL(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	expected := "https://github.com/open-telemetry/opentelemetry-python-contrib"
	if a.RepoURL() != expected {
		t.Errorf("expected '%s', got '%s'", expected, a.RepoURL())
	}
}

func TestExtractComponentName(t *testing.T) {
	tests := []struct {
		path     string
		baseDir  string
		expected string
	}{
		{
			path:     "/repo/instrumentation/opentelemetry-instrumentation-flask/src/opentelemetry/instrumentation/flask/__init__.py",
			baseDir:  "/repo/instrumentation",
			expected: "flask",
		},
		{
			path:     "/repo/instrumentation/opentelemetry-instrumentation-django/src/opentelemetry/instrumentation/django/__init__.py",
			baseDir:  "/repo/instrumentation",
			expected: "django",
		},
		{
			path:     "/repo/instrumentation/opentelemetry-instrumentation-system-metrics/src/system_metrics.py",
			baseDir:  "/repo/instrumentation",
			expected: "system-metrics",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := extractComponentName(tt.path, tt.baseDir)
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestDeduplicateMetrics(t *testing.T) {
	metrics := []*adpt.RawMetric{
		{Name: "metric.a", ComponentName: "comp1", Description: ""},
		{Name: "metric.a", ComponentName: "comp1", Description: "Has description"},
		{Name: "metric.b", ComponentName: "comp1", Description: "B metric"},
	}

	result := deduplicateMetrics(metrics)

	if len(result) != 2 {
		t.Fatalf("expected 2 metrics, got %d", len(result))
	}

	metricAFound := false
	for _, m := range result {
		if m.Name == "metric.a" {
			metricAFound = true
			if m.Description != "Has description" {
				t.Errorf("expected metric.a to have description 'Has description', got '%s'", m.Description)
			}
		}
	}

	if !metricAFound {
		t.Error("metric.a not found in results")
	}
}
