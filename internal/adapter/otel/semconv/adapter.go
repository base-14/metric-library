package semconv

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/base14/otel-glossary/internal/adapter"
	"github.com/base14/otel-glossary/internal/domain"
	"github.com/base14/otel-glossary/internal/fetcher"
)

const repoURL = "https://github.com/open-telemetry/semantic-conventions"

type Adapter struct {
	fetcher *fetcher.GitFetcher
}

func NewAdapter(cacheDir string) *Adapter {
	return &Adapter{
		fetcher: fetcher.NewGitFetcher(cacheDir),
	}
}

func (a *Adapter) Name() string {
	return "otel-semconv"
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
	modelDir := filepath.Join(result.RepoPath, "model")

	var metrics []*adapter.RawMetric

	err := filepath.Walk(modelDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if info.Name() != "metrics.yaml" {
			return nil
		}

		defs, parseErr := ParseFile(path)
		if parseErr != nil {
			return nil
		}

		componentName := deriveComponentName(path, modelDir)

		for _, def := range defs {
			attrs := make([]domain.Attribute, 0, len(def.Attributes))
			for _, attr := range def.Attributes {
				attrs = append(attrs, domain.Attribute{
					Name:     attr.Ref,
					Type:     "string",
					Required: attr.RequirementLevel == "required",
				})
			}

			rawMetric := &adapter.RawMetric{
				Name:             def.Name,
				Description:      def.Brief,
				InstrumentType:   mapInstrumentType(def.Instrument),
				Unit:             def.Unit,
				Attributes:       attrs,
				EnabledByDefault: true,
				ComponentType:    string(domain.ComponentInstrumentation),
				ComponentName:    componentName,
				SourceLocation:   path,
				Path:             path,
			}

			metrics = append(metrics, rawMetric)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func deriveComponentName(filePath, modelDir string) string {
	rel, err := filepath.Rel(modelDir, filePath)
	if err != nil {
		return "unknown"
	}

	dir := filepath.Dir(rel)
	if dir == "." {
		return "general"
	}

	return strings.ReplaceAll(dir, string(filepath.Separator), ".")
}

func mapInstrumentType(instrument string) string {
	switch instrument {
	case "counter":
		return string(domain.InstrumentCounter)
	case "updowncounter":
		return string(domain.InstrumentUpDownCounter)
	case "gauge":
		return string(domain.InstrumentGauge)
	case "histogram":
		return string(domain.InstrumentHistogram)
	default:
		return string(domain.InstrumentGauge)
	}
}
