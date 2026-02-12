package servicebus

import (
	"context"
	"testing"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
)

func TestAdapterName(t *testing.T) {
	a := NewAdapter(".cache")
	if a.Name() != "azure-servicebus" {
		t.Errorf("expected name azure-servicebus, got %s", a.Name())
	}
}

func TestAdapterSourceCategory(t *testing.T) {
	a := NewAdapter(".cache")
	if a.SourceCategory() != domain.SourceCloud {
		t.Errorf("expected source category cloud, got %s", a.SourceCategory())
	}
}

func TestAdapterConfidence(t *testing.T) {
	a := NewAdapter(".cache")
	if a.Confidence() != domain.ConfidenceDocumented {
		t.Errorf("expected confidence documented, got %s", a.Confidence())
	}
}

func TestAdapterExtractionMethod(t *testing.T) {
	a := NewAdapter(".cache")
	if a.ExtractionMethod() != domain.ExtractionScrape {
		t.Errorf("expected extraction method scrape, got %s", a.ExtractionMethod())
	}
}

func TestAdapterImplementsInterface(t *testing.T) {
	a := NewAdapter(".cache")
	var _ adapter.Adapter = a
}

func TestAdapterExtract(t *testing.T) {
	a := NewAdapter(".cache")
	ctx := context.Background()

	fetchResult, err := a.Fetch(ctx, adapter.FetchOptions{})
	if err != nil {
		t.Fatalf("unexpected fetch error: %v", err)
	}

	metrics, err := a.Extract(ctx, fetchResult)
	if err != nil {
		t.Fatalf("unexpected extract error: %v", err)
	}

	if len(metrics) != 20 {
		t.Errorf("expected 20 metrics, got %d", len(metrics))
	}

	for _, m := range metrics {
		if m.Name == "" {
			t.Error("metric name should not be empty")
		}
		if m.ComponentName != "Service Bus" {
			t.Errorf("expected component name Service Bus, got %s", m.ComponentName)
		}
		if m.ComponentType != "platform" {
			t.Errorf("expected component type platform, got %s", m.ComponentType)
		}
		if m.SourceLocation != "Microsoft.ServiceBus/namespaces" {
			t.Errorf("expected source location Microsoft.ServiceBus/namespaces, got %s", m.SourceLocation)
		}
	}
}
