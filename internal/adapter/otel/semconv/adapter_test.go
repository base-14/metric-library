package semconv

import (
	"context"
	"testing"

	"github.com/base14/otel-glossary/internal/adapter"
	"github.com/base14/otel-glossary/internal/domain"
)

func TestAdapterName(t *testing.T) {
	a := NewAdapter(".cache")
	if a.Name() != "otel-semconv" {
		t.Errorf("expected name otel-semconv, got %s", a.Name())
	}
}

func TestAdapterSourceCategory(t *testing.T) {
	a := NewAdapter(".cache")
	if a.SourceCategory() != domain.SourceOTEL {
		t.Errorf("expected source category otel, got %s", a.SourceCategory())
	}
}

func TestAdapterConfidence(t *testing.T) {
	a := NewAdapter(".cache")
	if a.Confidence() != domain.ConfidenceAuthoritative {
		t.Errorf("expected confidence authoritative, got %s", a.Confidence())
	}
}

func TestAdapterExtractionMethod(t *testing.T) {
	a := NewAdapter(".cache")
	if a.ExtractionMethod() != domain.ExtractionMetadata {
		t.Errorf("expected extraction method metadata, got %s", a.ExtractionMethod())
	}
}

func TestAdapterRepoURL(t *testing.T) {
	a := NewAdapter(".cache")
	if a.RepoURL() != "https://github.com/open-telemetry/semantic-conventions" {
		t.Errorf("expected repo URL, got %s", a.RepoURL())
	}
}

func TestAdapterImplementsInterface(t *testing.T) {
	a := NewAdapter(".cache")
	var _ adapter.Adapter = a
}

func TestAdapterExtractIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	a := NewAdapter(".cache")

	ctx := context.Background()
	fetchResult, err := a.Fetch(ctx, adapter.FetchOptions{})
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	metrics, err := a.Extract(ctx, fetchResult)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) < 50 {
		t.Errorf("expected at least 50 metrics, got %d", len(metrics))
	}

	foundHTTP := false
	foundSystem := false

	for _, m := range metrics {
		if m.Name == "http.server.request.duration" {
			foundHTTP = true
			if m.InstrumentType != "histogram" {
				t.Errorf("expected http metric to be histogram, got %s", m.InstrumentType)
			}
		}
		if m.Name == "system.cpu.time" {
			foundSystem = true
		}
	}

	if !foundHTTP {
		t.Error("expected to find http.server.request.duration metric")
	}
	if !foundSystem {
		t.Error("expected to find system.cpu.time metric")
	}

	t.Logf("Extracted %d metrics from semantic-conventions", len(metrics))
}
