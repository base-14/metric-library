package rust

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
	"github.com/base-14/metric-library/internal/fetcher"
)

const repoURL = "https://github.com/open-telemetry/opentelemetry-rust-contrib"

type Adapter struct {
	fetcher *fetcher.GitFetcher
}

func NewAdapter(cacheDir string) *Adapter {
	return &Adapter{
		fetcher: fetcher.NewGitFetcher(cacheDir),
	}
}

func (a *Adapter) Name() string {
	return "otel-rust"
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

	err := filepath.Walk(result.RepoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			// Skip test, target, and example directories
			name := info.Name()
			if name == "target" || name == "tests" || name == "benches" {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".rs") {
			return nil
		}

		// Skip test files
		if strings.Contains(path, "/tests/") || strings.HasSuffix(path, "_test.rs") {
			return nil
		}

		componentName := extractComponentName(path, result.RepoPath)

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
		// First part is typically the crate name
		// e.g., "opentelemetry-instrumentation-actix-web/..." -> "actix-web"
		crateName := parts[0]

		// Extract the meaningful part
		if suffix, found := strings.CutPrefix(crateName, "opentelemetry-instrumentation-"); found {
			return strings.ToLower(suffix)
		}

		if suffix, found := strings.CutPrefix(crateName, "opentelemetry-"); found {
			return strings.ToLower(suffix)
		}

		return strings.ToLower(crateName)
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
