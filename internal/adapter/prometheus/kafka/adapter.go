package kafka

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/base14/otel-glossary/internal/adapter"
	"github.com/base14/otel-glossary/internal/adapter/prometheus/astparser"
	"github.com/base14/otel-glossary/internal/domain"
	"github.com/base14/otel-glossary/internal/fetcher"
)

const repoURL = "https://github.com/danielqsj/kafka_exporter"

type Adapter struct {
	fetcher *fetcher.GitFetcher
}

func NewAdapter(cacheDir string) *Adapter {
	return &Adapter{
		fetcher: fetcher.NewGitFetcher(cacheDir),
	}
}

func (a *Adapter) Name() string {
	return "prometheus-kafka"
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

func (a *Adapter) Extract(ctx context.Context, result *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	entries, err := os.ReadDir(result.RepoPath)
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

		filePath := filepath.Join(result.RepoPath, entry.Name())
		componentName := deriveComponentName(entry.Name())

		defs, err := astparser.ParseFile(filePath)
		if err != nil {
			continue
		}

		for _, def := range defs {
			rawMetric := &adapter.RawMetric{
				Name:             def.Name,
				Description:      def.Help,
				InstrumentType:   inferInstrumentType(def.Name),
				Attributes:       labelsToAttributes(def.Labels),
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
	name = strings.TrimSuffix(name, "_exporter")
	return name
}

func inferInstrumentType(metricName string) string {
	switch {
	case strings.HasSuffix(metricName, "_total"):
		return string(domain.InstrumentCounter)
	case strings.HasSuffix(metricName, "_bucket"):
		return string(domain.InstrumentHistogram)
	case strings.HasSuffix(metricName, "_sum") && !strings.Contains(metricName, "offset_sum"):
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
