package java

import (
	"testing"

	adpt "github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
)

func TestAdapter_Name(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Name() != "otel-java" {
		t.Errorf("expected 'otel-java', got '%s'", a.Name())
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
	expected := "https://github.com/open-telemetry/opentelemetry-java-instrumentation"
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
			path:     "/repo/instrumentation/kafka/kafka-clients/library/src/main/java/Metrics.java",
			baseDir:  "/repo/instrumentation",
			expected: "kafka",
		},
		{
			path:     "/repo/instrumentation/runtime-telemetry/runtime-telemetry-java8/library/src/Cpu.java",
			baseDir:  "/repo/instrumentation",
			expected: "runtime-telemetry",
		},
		{
			path:     "/repo/instrumentation/hikaricp-3.0/library/src/HikariMetrics.java",
			baseDir:  "/repo/instrumentation",
			expected: "hikaricp-3.0",
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

func TestExtractApiComponentName(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{
			path:     "/repo/api/src/main/java/semconv/db/DbConnectionPoolMetrics.java",
			expected: "db-semconv",
		},
		{
			path:     "/repo/api/src/main/java/semconv/http/HttpMetrics.java",
			expected: "http-semconv",
		},
		{
			path:     "/repo/api/src/main/java/other/SomeMetrics.java",
			expected: "api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := extractApiComponentName(tt.path)
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
