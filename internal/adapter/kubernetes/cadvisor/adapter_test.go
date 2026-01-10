package cadvisor

import (
	"context"
	"testing"

	"github.com/base14/otel-glossary/internal/adapter"
	"github.com/base14/otel-glossary/internal/domain"
)

func TestAdapterName(t *testing.T) {
	a := NewAdapter(".cache")
	if a.Name() != "kubernetes-cadvisor" {
		t.Errorf("expected name kubernetes-cadvisor, got %s", a.Name())
	}
}

func TestAdapterSourceCategory(t *testing.T) {
	a := NewAdapter(".cache")
	if a.SourceCategory() != domain.SourceKubernetes {
		t.Errorf("expected source category kubernetes, got %s", a.SourceCategory())
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
	if a.ExtractionMethod() != domain.ExtractionAST {
		t.Errorf("expected extraction method ast, got %s", a.ExtractionMethod())
	}
}

func TestAdapterRepoURL(t *testing.T) {
	a := NewAdapter(".cache")
	if a.RepoURL() != "https://github.com/google/cadvisor" {
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

	if len(metrics) < 20 {
		t.Errorf("expected at least 20 metrics, got %d", len(metrics))
	}

	foundCPU := false
	foundMemory := false

	for _, m := range metrics {
		if m.Name == "container_cpu_usage_seconds_total" || m.Name == "container_cpu_user_seconds_total" {
			foundCPU = true
		}
		if m.Name == "container_memory_usage_bytes" || m.Name == "machine_memory_bytes" {
			foundMemory = true
		}
	}

	if !foundCPU {
		t.Error("expected to find CPU metric")
	}
	if !foundMemory {
		t.Error("expected to find memory metric")
	}

	t.Logf("Extracted %d metrics from cAdvisor", len(metrics))
}
