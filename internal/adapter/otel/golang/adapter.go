package golang

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
	"github.com/base-14/metric-library/internal/fetcher"
)

const repoURL = "https://github.com/open-telemetry/opentelemetry-go-contrib"

type Adapter struct {
	fetcher *fetcher.GitFetcher
}

func NewAdapter(cacheDir string) *Adapter {
	return &Adapter{
		fetcher: fetcher.NewGitFetcher(cacheDir),
	}
}

func (a *Adapter) Name() string {
	return "otel-go"
}

func (a *Adapter) SourceCategory() domain.SourceCategory {
	return domain.SourceOTEL
}

func (a *Adapter) Confidence() domain.ConfidenceLevel {
	return domain.ConfidenceDerived
}

func (a *Adapter) ExtractionMethod() domain.ExtractionMethod {
	return domain.ExtractionAST
}

func (a *Adapter) RepoURL() string {
	return repoURL
}

func (a *Adapter) Fetch(ctx context.Context, opts adapter.FetchOptions) (*adapter.FetchResult, error) {
	fetchOpts := fetcher.FetchOptions{
		RepoURL: repoURL,
		Commit:  opts.Commit,
		Shallow: true,
		Depth:   1,
		Force:   opts.Force,
	}

	result, err := a.fetcher.Fetch(ctx, fetchOpts)
	if err != nil {
		return nil, err
	}

	return &adapter.FetchResult{
		RepoPath:  result.RepoPath,
		Commit:    result.Commit,
		Timestamp: result.Timestamp,
	}, nil
}

func (a *Adapter) Extract(ctx context.Context, result *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	var metrics []*adapter.RawMetric

	instrumentationDir := filepath.Join(result.RepoPath, "instrumentation")
	err := filepath.Walk(instrumentationDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			// Skip test and vendor directories
			name := info.Name()
			if name == "test" || name == "testdata" || name == "vendor" || name == "example" || name == "examples" {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Skip test files
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		componentName := extractComponentName(path, instrumentationDir)

		defs, err := ParseFile(path)
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(result.RepoPath, path)

		for _, def := range defs {
			rawMetric := &adapter.RawMetric{
				Name:             def.Name,
				Description:      def.Description,
				Unit:             def.Unit,
				InstrumentType:   def.InstrumentType,
				EnabledByDefault: true,
				ComponentType:    string(domain.ComponentInstrumentation),
				ComponentName:    componentName,
				SourceLocation:   path,
				Path:             relPath,
			}

			metrics = append(metrics, rawMetric)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return deduplicateMetrics(metrics), nil
}

func extractComponentName(path, baseDir string) string {
	relPath, _ := filepath.Rel(baseDir, path)
	parts := strings.Split(relPath, string(filepath.Separator))

	if len(parts) > 0 {
		// First part is typically the instrumentation name
		// e.g., "runtime/..." -> "runtime"
		// or "google.golang.org/grpc/..." -> "grpc"
		component := parts[0]

		// Handle nested paths like google.golang.org/grpc
		if strings.Contains(component, ".") {
			// Take last segment of domain-style path
			segments := strings.Split(component, ".")
			component = segments[len(segments)-1]
		}

		return strings.ToLower(component)
	}

	return filepath.Base(filepath.Dir(path))
}

func deduplicateMetrics(metrics []*adapter.RawMetric) []*adapter.RawMetric {
	seen := make(map[string]*adapter.RawMetric)

	for _, m := range metrics {
		key := m.Name + "|" + m.ComponentName
		if existing, ok := seen[key]; ok {
			// Keep the one with more complete info
			if m.Description != "" && existing.Description == "" {
				seen[key] = m
			}
		} else {
			seen[key] = m
		}
	}

	result := make([]*adapter.RawMetric, 0, len(seen))
	for _, m := range seen {
		result = append(result, m)
	}

	return result
}
