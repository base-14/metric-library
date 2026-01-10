package ksm

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
	"github.com/base-14/metric-library/internal/fetcher"
)

const repoURL = "https://github.com/kubernetes/kube-state-metrics"

type Adapter struct {
	fetcher *fetcher.GitFetcher
}

func NewAdapter(cacheDir string) *Adapter {
	return &Adapter{
		fetcher: fetcher.NewGitFetcher(cacheDir),
	}
}

func (a *Adapter) Name() string {
	return "kubernetes-ksm"
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
	storeDir := filepath.Join(result.RepoPath, "internal", "store")

	entries, err := os.ReadDir(storeDir)
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

		if entry.Name() == "builder.go" || entry.Name() == "utils.go" {
			continue
		}

		filePath := filepath.Join(storeDir, entry.Name())
		componentName := deriveComponentName(entry.Name())

		defs, err := ParseFile(filePath)
		if err != nil {
			continue
		}

		for _, def := range defs {
			rawMetric := &adapter.RawMetric{
				Name:             def.Name,
				Description:      def.Help,
				InstrumentType:   def.MetricType,
				Attributes:       []domain.Attribute{},
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
	return name
}
