package adapter

import (
	"context"
	"testing"

	"github.com/base14/otel-glossary/internal/domain"
)

type mockAdapter struct {
	name string
}

func (m *mockAdapter) Name() string {
	return m.name
}

func (m *mockAdapter) Fetch(_ context.Context, _ FetchOptions) (*FetchResult, error) {
	return nil, nil
}

func (m *mockAdapter) Extract(_ context.Context, _ *FetchResult) ([]*RawMetric, error) {
	return nil, nil
}

func (m *mockAdapter) SourceCategory() domain.SourceCategory {
	return domain.SourceOTEL
}

func (m *mockAdapter) Confidence() domain.ConfidenceLevel {
	return domain.ConfidenceAuthoritative
}

func (m *mockAdapter) ExtractionMethod() domain.ExtractionMethod {
	return domain.ExtractionMetadata
}

func (m *mockAdapter) RepoURL() string {
	return "https://github.com/test/repo"
}

func TestAdapterRegistry_Register(t *testing.T) {
	registry := NewRegistry()
	adapter := &mockAdapter{name: "test-adapter"}

	registry.Register(adapter)

	got, ok := registry.Get("test-adapter")
	if !ok {
		t.Fatal("Register() adapter not found after registration")
	}
	if got.Name() != "test-adapter" {
		t.Errorf("Get() returned adapter with name %q, want %q", got.Name(), "test-adapter")
	}
}

func TestAdapterRegistry_Get_NotFound(t *testing.T) {
	registry := NewRegistry()

	_, ok := registry.Get("nonexistent")
	if ok {
		t.Error("Get() should return false for nonexistent adapter")
	}
}

func TestAdapterRegistry_All(t *testing.T) {
	registry := NewRegistry()
	registry.Register(&mockAdapter{name: "adapter1"})
	registry.Register(&mockAdapter{name: "adapter2"})

	all := registry.All()
	if len(all) != 2 {
		t.Errorf("All() returned %d adapters, want 2", len(all))
	}
}

func TestAdapterRegistry_Names(t *testing.T) {
	registry := NewRegistry()
	registry.Register(&mockAdapter{name: "adapter1"})
	registry.Register(&mockAdapter{name: "adapter2"})

	names := registry.Names()
	if len(names) != 2 {
		t.Errorf("Names() returned %d names, want 2", len(names))
	}

	nameMap := make(map[string]bool)
	for _, name := range names {
		nameMap[name] = true
	}
	if !nameMap["adapter1"] || !nameMap["adapter2"] {
		t.Error("Names() missing expected adapter names")
	}
}
