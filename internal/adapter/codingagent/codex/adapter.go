package codex

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
	"github.com/base-14/metric-library/internal/fetcher"
)

const repoURL = "https://github.com/openai/codex"

var metricPattern = regexp.MustCompile(`pub\s+const\s+\w+:\s*&str\s*=\s*"([^"]+)"`)

type Adapter struct {
	fetcher *fetcher.GitFetcher
}

func NewAdapter(cacheDir string) *Adapter {
	return &Adapter{
		fetcher: fetcher.NewGitFetcher(cacheDir),
	}
}

func (a *Adapter) Name() string {
	return "codingagent-codex"
}

func (a *Adapter) SourceCategory() domain.SourceCategory {
	return domain.SourceCodingAgent
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

func (a *Adapter) Extract(_ context.Context, result *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	namesPath := filepath.Join(result.RepoPath, "codex-rs", "otel", "src", "metrics", "names.rs")

	content, err := os.ReadFile(filepath.Clean(namesPath))
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", namesPath, err)
	}

	matches := metricPattern.FindAllStringSubmatch(string(content), -1)

	var metrics []*adapter.RawMetric
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		name := match[1]
		instrumentType, unit := inferType(name)

		relPath, _ := filepath.Rel(result.RepoPath, namesPath)
		metrics = append(metrics, &adapter.RawMetric{
			Name:             name,
			Description:      inferDescription(name),
			Unit:             unit,
			InstrumentType:   instrumentType,
			EnabledByDefault: true,
			ComponentType:    string(domain.ComponentPlatform),
			ComponentName:    "codex",
			SourceLocation:   namesPath,
			Path:             relPath,
		})
	}

	return metrics, nil
}

func inferType(name string) (instrumentType, unit string) {
	if strings.HasSuffix(name, "_ms") || strings.HasSuffix(name, ".duration_ms") {
		return string(domain.InstrumentHistogram), "ms"
	}
	return string(domain.InstrumentCounter), "count"
}

func inferDescription(name string) string {
	parts := strings.Split(name, ".")
	if len(parts) < 2 {
		return name
	}
	// Remove prefix, join remaining parts with spaces
	desc := strings.Join(parts[1:], " ")
	desc = strings.ReplaceAll(desc, "_", " ")
	return cases.Title(language.English).String(desc)
}
