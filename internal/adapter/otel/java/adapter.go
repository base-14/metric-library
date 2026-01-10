package java

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
	"github.com/base-14/metric-library/internal/fetcher"
)

const repoURL = "https://github.com/open-telemetry/opentelemetry-java-instrumentation"

type Adapter struct {
	fetcher *fetcher.GitFetcher
}

func NewAdapter(cacheDir string) *Adapter {
	return &Adapter{
		fetcher: fetcher.NewGitFetcher(cacheDir),
	}
}

func (a *Adapter) Name() string {
	return "otel-java"
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

	// Walk the instrumentation directory
	instrumentationDir := filepath.Join(result.RepoPath, "instrumentation")
	err := filepath.Walk(instrumentationDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".java") {
			return nil
		}

		// Skip test files
		if strings.Contains(path, "/test/") || strings.Contains(path, "/jmh/") {
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

	// Also scan instrumentation-api-incubator for semantic convention helpers
	apiDir := filepath.Join(result.RepoPath, "instrumentation-api-incubator")
	_ = filepath.Walk(apiDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".java") {
			return nil
		}

		if strings.Contains(path, "/test/") {
			return nil
		}

		componentName := extractApiComponentName(path)

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

	return deduplicateMetrics(metrics), nil
}

func extractComponentName(path, baseDir string) string {
	relPath, _ := filepath.Rel(baseDir, path)
	parts := strings.Split(relPath, string(filepath.Separator))

	if len(parts) > 0 {
		// First part is typically the instrumentation name
		// e.g., "kafka/kafka-clients/..." -> "kafka"
		// or "runtime-telemetry/runtime-telemetry-java8/..." -> "runtime-telemetry"
		return parts[0]
	}

	return filepath.Base(filepath.Dir(path))
}

func extractApiComponentName(path string) string {
	// Extract meaningful name from API path
	// e.g., ".../semconv/db/DbConnectionPoolMetrics.java" -> "db-pool"
	dir := filepath.Dir(path)
	parts := strings.Split(dir, string(filepath.Separator))

	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == "semconv" && i+1 < len(parts) {
			return parts[i+1] + "-semconv"
		}
	}

	return "api"
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
