package openlit

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
	"github.com/base-14/metric-library/internal/fetcher"
)

const repoURL = "https://github.com/openlit/openlit"

type Adapter struct {
	fetcher *fetcher.GitFetcher
}

func NewAdapter(cacheDir string) *Adapter {
	return &Adapter{
		fetcher: fetcher.NewGitFetcher(cacheDir),
	}
}

func (a *Adapter) Name() string {
	return "openlit"
}

func (a *Adapter) SourceCategory() domain.SourceCategory {
	return domain.SourceVendor
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

	sdkDir := filepath.Join(result.RepoPath, "sdk", "python", "src", "openlit")

	// First, load semantic conventions to resolve constants
	semconvPath := filepath.Join(sdkDir, "semcov", "__init__.py")
	constants := loadSemanticConstants(semconvPath)

	// Parse the main metrics file
	metricsPath := filepath.Join(sdkDir, "otel", "metrics.py")
	defs, err := parseFileWithConstants(metricsPath, constants)
	if err == nil {
		relPath, _ := filepath.Rel(result.RepoPath, metricsPath)
		for _, def := range defs {
			rawMetric := &adapter.RawMetric{
				Name:             def.Name,
				Description:      def.Description,
				Unit:             def.Unit,
				InstrumentType:   def.InstrumentType,
				EnabledByDefault: true,
				ComponentType:    string(domain.ComponentInstrumentation),
				ComponentName:    "openlit",
				SourceLocation:   metricsPath,
				Path:             relPath,
			}
			metrics = append(metrics, rawMetric)
		}
	}

	// Also walk instrumentation directory for any additional metrics
	instrumentationDir := filepath.Join(sdkDir, "instrumentation")
	_ = filepath.Walk(instrumentationDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".py") {
			return nil
		}

		if strings.Contains(path, "/tests/") {
			return nil
		}

		componentName := extractComponentName(path, instrumentationDir)

		defs, err := parseFileWithConstants(path, constants)
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

func loadSemanticConstants(path string) map[string]string {
	constants := make(map[string]string)

	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return constants
	}

	// Match patterns like: CONSTANT_NAME = "value"
	pattern := regexp.MustCompile(`(\w+)\s*=\s*["']([^"']+)["']`)
	matches := pattern.FindAllStringSubmatch(string(content), -1)

	for _, match := range matches {
		if len(match) >= 3 {
			constants[match[1]] = match[2]
		}
	}

	return constants
}

func extractComponentName(path, baseDir string) string {
	relPath, _ := filepath.Rel(baseDir, path)
	parts := strings.Split(relPath, string(filepath.Separator))

	if len(parts) > 0 {
		return parts[0]
	}

	return filepath.Base(filepath.Dir(path))
}

func deduplicateMetrics(metrics []*adapter.RawMetric) []*adapter.RawMetric {
	seen := make(map[string]*adapter.RawMetric)

	for _, m := range metrics {
		key := m.Name
		if existing, ok := seen[key]; ok {
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
