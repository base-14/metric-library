package adapter

import (
	"context"
	"time"

	"github.com/base-14/metric-library/internal/domain"
)

type FetchOptions struct {
	Commit   string
	CacheDir string
	Force    bool
}

type FetchResult struct {
	RepoPath  string
	Commit    string
	Timestamp time.Time
	Files     []string
}

type RawMetric struct {
	Name             string
	InstrumentType   string
	Description      string
	Unit             string
	Attributes       []domain.Attribute
	EnabledByDefault bool
	ComponentType    string
	ComponentName    string
	SourceLocation   string
	Path             string
}

type Adapter interface {
	Name() string
	Fetch(ctx context.Context, opts FetchOptions) (*FetchResult, error)
	Extract(ctx context.Context, result *FetchResult) ([]*RawMetric, error)
	SourceCategory() domain.SourceCategory
	Confidence() domain.ConfidenceLevel
	ExtractionMethod() domain.ExtractionMethod
	RepoURL() string
}

type AdapterRegistry struct {
	adapters map[string]Adapter
}

func NewRegistry() *AdapterRegistry {
	return &AdapterRegistry{
		adapters: make(map[string]Adapter),
	}
}

func (r *AdapterRegistry) Register(adapter Adapter) {
	r.adapters[adapter.Name()] = adapter
}

func (r *AdapterRegistry) Get(name string) (Adapter, bool) {
	adapter, ok := r.adapters[name]
	return adapter, ok
}

func (r *AdapterRegistry) All() []Adapter {
	result := make([]Adapter, 0, len(r.adapters))
	for _, adapter := range r.adapters {
		result = append(result, adapter)
	}
	return result
}

func (r *AdapterRegistry) Names() []string {
	names := make([]string, 0, len(r.adapters))
	for name := range r.adapters {
		names = append(names, name)
	}
	return names
}
