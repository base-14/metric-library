package clickhouse

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

const repoURL = "https://github.com/ClickHouse/ClickHouse"

type Adapter struct {
	fetcher *fetcher.GitFetcher
}

func NewAdapter(cacheDir string) *Adapter {
	return &Adapter{
		fetcher: fetcher.NewGitFetcher(cacheDir),
	}
}

func (a *Adapter) Name() string {
	return "prometheus-clickhouse"
}

func (a *Adapter) SourceCategory() domain.SourceCategory {
	return domain.SourcePrometheus
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
	var metrics []*adapter.RawMetric

	// 1. CurrentMetrics.cpp → gauges (system.metrics)
	currentMetricsPath := filepath.Join(result.RepoPath, "src", "Common", "CurrentMetrics.cpp")
	if src, err := os.ReadFile(currentMetricsPath); err == nil { //nolint:gosec // trusted repo path
		for _, def := range parseCurrentMetrics(src) {
			metrics = append(metrics, &adapter.RawMetric{
				Name:             "ClickHouseMetrics_" + def.Name,
				Description:      def.Description,
				InstrumentType:   string(domain.InstrumentGauge),
				EnabledByDefault: true,
				ComponentType:    string(domain.ComponentPlatform),
				ComponentName:    "current_metrics",
				SourceLocation:   currentMetricsPath,
				Path:             currentMetricsPath,
			})
		}
	}

	// 2. ProfileEvents.cpp → counters (system.events)
	profileEventsPath := filepath.Join(result.RepoPath, "src", "Common", "ProfileEvents.cpp")
	if src, err := os.ReadFile(profileEventsPath); err == nil { //nolint:gosec // trusted repo path
		for _, def := range parseProfileEvents(src) {
			metrics = append(metrics, &adapter.RawMetric{
				Name:             "ClickHouseProfileEvents_" + def.Name,
				Description:      def.Description,
				Unit:             def.Unit,
				InstrumentType:   string(domain.InstrumentCounter),
				EnabledByDefault: true,
				ComponentType:    string(domain.ComponentPlatform),
				ComponentName:    "profile_events",
				SourceLocation:   profileEventsPath,
				Path:             profileEventsPath,
			})
		}
	}

	// 3. Async metrics (two files) → gauges (system.asynchronous_metrics)
	asyncFiles := []string{
		filepath.Join(result.RepoPath, "src", "Interpreters", "ServerAsynchronousMetrics.cpp"),
		filepath.Join(result.RepoPath, "src", "Common", "AsynchronousMetrics.cpp"),
	}
	seen := make(map[string]bool)
	for _, asyncPath := range asyncFiles {
		src, err := os.ReadFile(asyncPath) //nolint:gosec // trusted repo path
		if err != nil {
			continue
		}
		for _, def := range parseAsyncMetrics(src) {
			if seen[def.Name] {
				continue
			}
			seen[def.Name] = true
			metrics = append(metrics, &adapter.RawMetric{
				Name:             "ClickHouseAsyncMetrics_" + def.Name,
				Description:      def.Description,
				InstrumentType:   string(domain.InstrumentGauge),
				EnabledByDefault: true,
				ComponentType:    string(domain.ComponentPlatform),
				ComponentName:    "async_metrics",
				SourceLocation:   asyncPath,
				Path:             asyncPath,
			})
		}
	}

	return metrics, nil
}

type metricDef struct {
	Name        string
	Description string
	Unit        string
}

// currentMetricRe matches: M(Name, "description")
var currentMetricRe = regexp.MustCompile(`M\((\w+),\s*"([^"]+)"\)`)

func parseCurrentMetrics(src []byte) []metricDef {
	matches := currentMetricRe.FindAllSubmatch(src, -1)
	defs := make([]metricDef, 0, len(matches))
	for _, m := range matches {
		defs = append(defs, metricDef{
			Name:        string(m[1]),
			Description: string(m[2]),
		})
	}
	return defs
}

// profileEventRe matches: M(Name, "description", ValueType::Type)
var profileEventRe = regexp.MustCompile(`M\((\w+),\s*"([^"]+)",\s*ValueType::(\w+)\)`)

func parseProfileEvents(src []byte) []metricDef {
	matches := profileEventRe.FindAllSubmatch(src, -1)
	defs := make([]metricDef, 0, len(matches))
	for _, m := range matches {
		defs = append(defs, metricDef{
			Name:        string(m[1]),
			Description: string(m[2]),
			Unit:        valueTypeToUnit(string(m[3])),
		})
	}
	return defs
}

// asyncMetricRe matches: new_values["Name"] with static string key (not fmt::format)
var asyncMetricRe = regexp.MustCompile(`new_values\["(\w+)"\]`)

// asyncDescRe matches the description in: new_values["Name"] = { ..., "description" };
// Handles multiline cases where the description may be on the next line.
var asyncDescRe = regexp.MustCompile(`new_values\["(\w+)"\]\s*=\s*\{[^}]*?"([^"]+)"\s*\}`)

func parseAsyncMetrics(src []byte) []metricDef {
	// Build description map from full pattern matches
	descMap := make(map[string]string)
	descMatches := asyncDescRe.FindAllSubmatch(src, -1)
	for _, m := range descMatches {
		descMap[string(m[1])] = string(m[2])
	}

	// Extract all static metric names
	seen := make(map[string]bool)
	var defs []metricDef
	nameMatches := asyncMetricRe.FindAllSubmatch(src, -1)
	for _, m := range nameMatches {
		name := string(m[1])
		if seen[name] {
			continue
		}
		seen[name] = true
		defs = append(defs, metricDef{
			Name:        name,
			Description: descMap[name],
		})
	}
	return defs
}

func valueTypeToUnit(vt string) string {
	switch strings.TrimSpace(vt) {
	case "Bytes":
		return "bytes"
	case "Microseconds":
		return "microseconds"
	case "Milliseconds":
		return "milliseconds"
	case "Nanoseconds":
		return "nanoseconds"
	default:
		return ""
	}
}
