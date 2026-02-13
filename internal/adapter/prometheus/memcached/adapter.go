package memcached

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/adapter/prometheus/astparser"
	"github.com/base-14/metric-library/internal/domain"
	"github.com/base-14/metric-library/internal/fetcher"
)

const repoURL = "https://github.com/prometheus/memcached_exporter"

type Adapter struct {
	fetcher *fetcher.GitFetcher
}

func NewAdapter(cacheDir string) *Adapter {
	return &Adapter{
		fetcher: fetcher.NewGitFetcher(cacheDir),
	}
}

func (a *Adapter) Name() string {
	return "prometheus-memcached"
}

func (a *Adapter) SourceCategory() domain.SourceCategory {
	return domain.SourcePrometheus
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

func (a *Adapter) Extract(_ context.Context, result *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	exporterDir := filepath.Join(result.RepoPath, "pkg", "exporter")

	entries, err := os.ReadDir(exporterDir)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var metrics []*adapter.RawMetric

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		if strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}

		filePath := filepath.Join(exporterDir, entry.Name())

		defs, err := astparser.ParseFile(filePath)
		if err != nil {
			continue
		}

		for _, def := range defs {
			if def.Name == "" || seen[def.Name] {
				continue
			}
			seen[def.Name] = true

			metrics = append(metrics, &adapter.RawMetric{
				Name:             def.Name,
				Description:      def.Help,
				InstrumentType:   inferInstrumentType(def.Name),
				Attributes:       labelsToAttributes(def.Labels),
				EnabledByDefault: true,
				ComponentType:    string(domain.ComponentPlatform),
				ComponentName:    "memcached",
				SourceLocation:   filePath,
				Path:             filePath,
			})
		}
	}

	return metrics, nil
}

func inferInstrumentType(metricName string) string {
	switch {
	case strings.HasSuffix(metricName, "_total"):
		return string(domain.InstrumentCounter)
	case strings.HasSuffix(metricName, "_bucket"):
		return string(domain.InstrumentHistogram)
	case strings.HasSuffix(metricName, "_sum"):
		return string(domain.InstrumentCounter)
	case strings.HasSuffix(metricName, "_count"):
		return string(domain.InstrumentCounter)
	default:
		return string(domain.InstrumentGauge)
	}
}

func labelsToAttributes(labels []string) []domain.Attribute {
	attrs := make([]domain.Attribute, 0, len(labels))
	for _, label := range labels {
		attrs = append(attrs, domain.Attribute{
			Name: label,
			Type: "string",
		})
	}
	return attrs
}
