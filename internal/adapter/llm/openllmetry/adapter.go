package openllmetry

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

const repoURL = "https://github.com/traceloop/openllmetry"

type Adapter struct {
	fetcher *fetcher.GitFetcher
}

func NewAdapter(cacheDir string) *Adapter {
	return &Adapter{
		fetcher: fetcher.NewGitFetcher(cacheDir),
	}
}

func (a *Adapter) Name() string {
	return "openllmetry"
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

	packagesDir := filepath.Join(result.RepoPath, "packages")

	// Load semantic conventions from the semconv-ai package
	semconvPath := filepath.Join(packagesDir, "opentelemetry-semantic-conventions-ai", "opentelemetry", "semconv_ai", "__init__.py")
	constants := loadSemanticConstants(semconvPath)

	// Also add known constants from the official opentelemetry-semconv package
	// These are defined in opentelemetry.semconv._incubating.metrics.gen_ai_metrics
	addOfficialSemconvConstants(constants)

	// Walk the packages directory
	err := filepath.Walk(packagesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".py") {
			return nil
		}

		// Skip test files
		if strings.Contains(path, "/tests/") || strings.HasSuffix(info.Name(), "_test.py") {
			return nil
		}

		componentName := extractComponentName(path, packagesDir)

		defs, parseErr := parseFileWithConstants(path, constants)
		if parseErr != nil {
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

func loadSemanticConstants(path string) map[string]string {
	constants := make(map[string]string)

	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return constants
	}

	// Match patterns like: CONSTANT_NAME = "value" (inside class or at module level)
	pattern := regexp.MustCompile(`(\w+)\s*=\s*["']([^"']+)["']`)
	matches := pattern.FindAllStringSubmatch(string(content), -1)

	for _, match := range matches {
		if len(match) >= 3 {
			constants[match[1]] = match[2]
		}
	}

	return constants
}

func addOfficialSemconvConstants(constants map[string]string) {
	// Constants from opentelemetry.semconv._incubating.metrics.gen_ai_metrics
	// These follow the OpenTelemetry GenAI semantic conventions
	officialConstants := map[string]string{
		"GEN_AI_SERVER_TIME_TO_FIRST_TOKEN": "gen_ai.server.time_to_first_token",
		"GEN_AI_CLIENT_TOKEN_USAGE":         "gen_ai.client.token.usage",
		"GEN_AI_CLIENT_OPERATION_DURATION":  "gen_ai.client.operation.duration",
	}

	for k, v := range officialConstants {
		constants[k] = v
	}
}

func extractComponentName(path, baseDir string) string {
	relPath, _ := filepath.Rel(baseDir, path)
	parts := strings.Split(relPath, string(filepath.Separator))

	if len(parts) > 0 {
		pkg := parts[0]
		// e.g., "opentelemetry-instrumentation-openai" -> "openai"
		if strings.HasPrefix(pkg, "opentelemetry-instrumentation-") {
			return strings.TrimPrefix(pkg, "opentelemetry-instrumentation-")
		}
		if strings.HasPrefix(pkg, "opentelemetry-semantic-conventions-") {
			return "semconv-ai"
		}
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
