package gemini

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
	"github.com/base-14/metric-library/internal/fetcher"
)

const repoURL = "https://github.com/google-gemini/gemini-cli"

var (
	// Matches: meter.createCounter('name', { or meter.createHistogram("name", {
	metricCallPattern = regexp.MustCompile(`meter\.create(Counter|Histogram)\(\s*['"]([^'"]+)['"]`)
	// Matches: description: 'text' or description: "text"
	descriptionPattern = regexp.MustCompile(`description:\s*['"]([^'"]+)['"]`)
	// Matches: unit: 'text' or unit: "text"
	unitPattern = regexp.MustCompile(`unit:\s*['"]([^'"]+)['"]`)
)

type Adapter struct {
	fetcher *fetcher.GitFetcher
}

func NewAdapter(cacheDir string) *Adapter {
	return &Adapter{
		fetcher: fetcher.NewGitFetcher(cacheDir),
	}
}

func (a *Adapter) Name() string {
	return "codingagent-gemini"
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
	telemetryDir := filepath.Join(result.RepoPath, "packages", "core", "src", "telemetry")

	entries, err := os.ReadDir(telemetryDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read telemetry dir: %w", err)
	}

	var metrics []*adapter.RawMetric

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".ts") {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".test.ts") || strings.HasSuffix(entry.Name(), ".spec.ts") {
			continue
		}

		filePath := filepath.Join(telemetryDir, entry.Name())
		fileMetrics, err := parseTypeScriptFile(filePath, result.RepoPath)
		if err != nil {
			continue
		}
		metrics = append(metrics, fileMetrics...)
	}

	return metrics, nil
}

func parseTypeScriptFile(filePath, repoRoot string) ([]*adapter.RawMetric, error) {
	content, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}

	text := string(content)
	relPath, _ := filepath.Rel(repoRoot, filePath)

	callMatches := metricCallPattern.FindAllStringSubmatchIndex(text, -1)
	var metrics []*adapter.RawMetric

	for _, loc := range callMatches {
		fullMatch := text[loc[0]:loc[1]]
		_ = fullMatch

		instrumentKind := text[loc[2]:loc[3]] // Counter or Histogram
		metricName := text[loc[4]:loc[5]]

		// Look for description and unit in the following ~200 chars
		end := min(loc[1]+200, len(text))
		block := text[loc[1]:end]

		description := ""
		if m := descriptionPattern.FindStringSubmatch(block); len(m) >= 2 {
			description = m[1]
		}

		unit := "count"
		if m := unitPattern.FindStringSubmatch(block); len(m) >= 2 {
			unit = m[1]
		}

		instrumentType := string(domain.InstrumentCounter)
		if instrumentKind == "Histogram" {
			instrumentType = string(domain.InstrumentHistogram)
		}

		metrics = append(metrics, &adapter.RawMetric{
			Name:             metricName,
			Description:      description,
			Unit:             unit,
			InstrumentType:   instrumentType,
			EnabledByDefault: true,
			ComponentType:    string(domain.ComponentPlatform),
			ComponentName:    "gemini-cli",
			SourceLocation:   filePath,
			Path:             relPath,
		})
	}

	return metrics, nil
}
