package cadvisor

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/base14/otel-glossary/internal/adapter"
	"github.com/base14/otel-glossary/internal/domain"
	"github.com/base14/otel-glossary/internal/fetcher"
)

const repoURL = "https://github.com/google/cadvisor"

type Adapter struct {
	fetcher *fetcher.GitFetcher
}

func NewAdapter(cacheDir string) *Adapter {
	return &Adapter{
		fetcher: fetcher.NewGitFetcher(cacheDir),
	}
}

func (a *Adapter) Name() string {
	return "kubernetes-cadvisor"
}

func (a *Adapter) SourceCategory() domain.SourceCategory {
	return domain.SourceKubernetes
}

func (a *Adapter) Confidence() domain.ConfidenceLevel {
	return domain.ConfidenceAuthoritative
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
	metricsDir := filepath.Join(result.RepoPath, "metrics")

	entries, err := os.ReadDir(metricsDir)
	if err != nil {
		return nil, err
	}

	var metrics []*adapter.RawMetric

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}

		if strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}

		if entry.Name() == "prometheus_fake.go" {
			continue
		}

		filePath := filepath.Join(metricsDir, entry.Name())
		componentName := deriveComponentName(entry.Name())

		defs, err := ParseFile(filePath)
		if err != nil {
			continue
		}

		for _, def := range defs {
			attrs := make([]domain.Attribute, 0, len(def.Labels))
			for _, label := range def.Labels {
				attrs = append(attrs, domain.Attribute{
					Name: label,
					Type: "string",
				})
			}

			rawMetric := &adapter.RawMetric{
				Name:             def.Name,
				Description:      def.Help,
				InstrumentType:   def.MetricType,
				Attributes:       attrs,
				EnabledByDefault: true,
				ComponentType:    string(domain.ComponentPlatform),
				ComponentName:    componentName,
				SourceLocation:   filePath,
				Path:             filePath,
			}

			metrics = append(metrics, rawMetric)
		}
	}

	return metrics, nil
}

func deriveComponentName(filename string) string {
	name := strings.TrimSuffix(filename, ".go")
	name = strings.TrimPrefix(name, "prometheus_")
	name = strings.TrimPrefix(name, "prometheus")
	if name == "" {
		return "container"
	}
	return name
}
