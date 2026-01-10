package ksm

import (
	"context"
	"testing"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
)

func TestAdapterName(t *testing.T) {
	a := NewAdapter(".cache")
	if a.Name() != "kubernetes-ksm" {
		t.Errorf("expected name kubernetes-ksm, got %s", a.Name())
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
	if a.RepoURL() != "https://github.com/kubernetes/kube-state-metrics" {
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

	foundPodInfo := false
	foundDeploymentCreated := false

	for _, m := range metrics {
		if m.Name == "kube_pod_info" {
			foundPodInfo = true
			if m.ComponentName != "pod" {
				t.Errorf("expected pod component, got %s", m.ComponentName)
			}
			if m.ComponentType != string(domain.ComponentPlatform) {
				t.Errorf("expected platform component type, got %s", m.ComponentType)
			}
		}
		if m.Name == "kube_deployment_created" {
			foundDeploymentCreated = true
		}
	}

	if !foundPodInfo {
		t.Error("expected to find kube_pod_info metric")
	}
	if !foundDeploymentCreated {
		t.Error("expected to find kube_deployment_created metric")
	}

	t.Logf("Extracted %d metrics from kube-state-metrics", len(metrics))
}
