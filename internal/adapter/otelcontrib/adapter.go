package otelcontrib

import (
	"context"

	"github.com/base14/otel-glossary/internal/adapter"
	"github.com/base14/otel-glossary/internal/discovery"
	"github.com/base14/otel-glossary/internal/domain"
	"github.com/base14/otel-glossary/internal/extractor"
	"github.com/base14/otel-glossary/internal/fetcher"
	"github.com/base14/otel-glossary/internal/parser"
)

const (
	repoURL    = "https://github.com/open-telemetry/opentelemetry-collector-contrib"
	sourceName = "opentelemetry-collector-contrib"
)

type Adapter struct {
	fetcher   *fetcher.GitFetcher
	discovery *discovery.MetadataDiscovery
	parser    *parser.MetadataParser
}

func NewAdapter(cacheDir string) *Adapter {
	return &Adapter{
		fetcher:   fetcher.NewGitFetcher(cacheDir),
		discovery: discovery.NewMetadataDiscovery(),
		parser:    parser.NewMetadataParser(),
	}
}

func (a *Adapter) Name() string {
	return "otel-collector-contrib"
}

func (a *Adapter) SourceCategory() domain.SourceCategory {
	return domain.SourceOTEL
}

func (a *Adapter) Confidence() domain.ConfidenceLevel {
	return domain.ConfidenceAuthoritative
}

func (a *Adapter) ExtractionMethod() domain.ExtractionMethod {
	return domain.ExtractionMetadata
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
	files, err := a.discovery.FindMetadataFiles(result.RepoPath)
	if err != nil {
		return nil, err
	}

	var metrics []*adapter.RawMetric

	for _, file := range files {
		meta, err := a.parser.ParseFile(file.Path)
		if err != nil {
			continue
		}

		ext := extractor.NewMetricExtractor(sourceName, file.ComponentName, file.ComponentType)
		extracted, err := ext.Extract(meta)
		if err != nil {
			continue
		}

		for _, m := range extracted {
			rawMetric := &adapter.RawMetric{
				Name:             m.MetricName,
				Description:      m.Description,
				Unit:             m.Unit,
				InstrumentType:   string(m.InstrumentType),
				Attributes:       m.Attributes,
				EnabledByDefault: m.EnabledByDefault,
				ComponentName:    file.ComponentName,
				ComponentType:    file.ComponentType,
				SourceLocation:   file.Path,
				Path:             file.Path,
			}

			metrics = append(metrics, rawMetric)
		}
	}

	return metrics, nil
}
