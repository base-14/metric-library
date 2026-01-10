package js

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
	"github.com/base-14/metric-library/internal/fetcher"
)

const repoURL = "https://github.com/open-telemetry/opentelemetry-js-contrib"

type Adapter struct {
	fetcher *fetcher.GitFetcher
}

func NewAdapter(cacheDir string) *Adapter {
	return &Adapter{
		fetcher: fetcher.NewGitFetcher(cacheDir),
	}
}

func (a *Adapter) Name() string {
	return "otel-js"
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
	packagesDir := filepath.Join(result.RepoPath, "packages")

	var metrics []*adapter.RawMetric

	// First pass: collect metrics from semconv.ts files (authoritative source)
	semconvMetrics := make(map[string]*MetricDef)
	err := filepath.Walk(packagesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		if info.Name() != "semconv.ts" {
			return nil
		}

		defs, err := ParseSemconvFile(path)
		if err != nil {
			return nil
		}

		for _, def := range defs {
			semconvMetrics[def.Name] = def
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Second pass: walk through packages and extract metrics
	err = filepath.Walk(packagesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".ts") {
			return nil
		}

		// Skip test files and node_modules
		if strings.Contains(path, "/test/") ||
			strings.Contains(path, ".test.") ||
			strings.Contains(path, ".spec.") ||
			strings.Contains(path, "node_modules") {
			return nil
		}

		componentName := extractComponentName(path, packagesDir)
		if componentName == "" {
			return nil
		}

		// For semconv.ts files, we already processed them
		if info.Name() == "semconv.ts" {
			defs, _ := ParseSemconvFile(path)
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
		}

		// For other files, look for create* calls
		// Read semconv.ts from the same package if it exists
		semconvPath := filepath.Clean(filepath.Join(filepath.Dir(path), "semconv.ts"))
		semconvContent := ""
		if data, err := os.ReadFile(semconvPath); err == nil {
			semconvContent = string(data)
		}

		cleanPath := filepath.Clean(path)
		content, err := os.ReadFile(cleanPath)
		if err != nil {
			return nil
		}

		defs, err := parseInstrumentationContent(string(content), semconvContent)
		if err != nil || len(defs) == 0 {
			return nil
		}

		relPath, _ := filepath.Rel(result.RepoPath, path)

		for _, def := range defs {
			// Enrich with semconv data if available
			if semDef, ok := semconvMetrics[def.Name]; ok {
				if def.Description == "" {
					def.Description = semDef.Description
				}
				if def.Unit == "" {
					def.Unit = semDef.Unit
				}
			}

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
	if len(parts) == 0 {
		return ""
	}

	pkg := parts[0]

	// Extract instrumentation name from package dir
	// e.g., "instrumentation-runtime-node" -> "runtime-node"
	// e.g., "host-metrics" -> "host-metrics"
	if strings.HasPrefix(pkg, "instrumentation-") {
		return strings.TrimPrefix(pkg, "instrumentation-")
	}
	if strings.HasPrefix(pkg, "opentelemetry-") {
		return strings.TrimPrefix(pkg, "opentelemetry-")
	}

	return pkg
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
