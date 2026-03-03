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
	// Matches single-line: const TOOL_CALL_COUNT = 'gemini_cli.tool.call.count';
	constPattern = regexp.MustCompile(`const\s+(\w+)\s*=\s*'([^']+)';`)
	// Matches multiline: const FOO =\n  'gemini_cli.bar';
	constMultilinePattern = regexp.MustCompile(`const\s+(\w+)\s*=\s*\n\s*'([^']+)';`)
	// Matches description in definition blocks
	descriptionPattern = regexp.MustCompile(`description:\s*['"]([^'"]+)['"]`)
	// Matches unit in definition blocks
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
	metricsPath := filepath.Join(result.RepoPath, "packages", "core", "src", "telemetry", "metrics.ts")

	content, err := os.ReadFile(filepath.Clean(metricsPath))
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", metricsPath, err)
	}

	text := string(content)
	relPath, _ := filepath.Rel(result.RepoPath, metricsPath)

	// Step 1: Parse all const declarations to get constantName -> metricName mapping
	constants := parseConstants(text)

	// Step 2: Parse definition blocks to get metric type, description, and unit
	metrics := parseDefinitionBlocks(text, constants)

	var result2 []*adapter.RawMetric
	for _, m := range metrics {
		result2 = append(result2, &adapter.RawMetric{
			Name:             m.name,
			Description:      m.description,
			Unit:             m.unit,
			InstrumentType:   m.instrumentType,
			EnabledByDefault: true,
			ComponentType:    string(domain.ComponentPlatform),
			ComponentName:    "gemini-cli",
			SourceLocation:   metricsPath,
			Path:             relPath,
		})
	}

	return result2, nil
}

type metricInfo struct {
	name           string
	description    string
	unit           string
	instrumentType string
}

func parseConstants(text string) map[string]string {
	constants := make(map[string]string)

	for _, match := range constPattern.FindAllStringSubmatch(text, -1) {
		if len(match) >= 3 && strings.HasPrefix(match[2], "gemini_cli.") || strings.HasPrefix(match[2], "gen_ai.") {
			constants[match[1]] = match[2]
		}
	}

	for _, match := range constMultilinePattern.FindAllStringSubmatch(text, -1) {
		if len(match) >= 3 && (strings.HasPrefix(match[2], "gemini_cli.") || strings.HasPrefix(match[2], "gen_ai.")) {
			constants[match[1]] = match[2]
		}
	}

	return constants
}

func parseDefinitionBlocks(text string, constants map[string]string) []metricInfo {
	var metrics []metricInfo

	// Find all *COUNTER_DEFINITIONS blocks and *HISTOGRAM_DEFINITIONS blocks
	counterDefs := findDefinitionBlockEntries(text, "COUNTER_DEFINITIONS")
	histogramDefs := findDefinitionBlockEntries(text, "HISTOGRAM_DEFINITIONS")

	for constName, metricName := range constants {
		if entry, ok := counterDefs[constName]; ok {
			desc := ""
			if m := descriptionPattern.FindStringSubmatch(entry); len(m) >= 2 {
				desc = m[1]
			}
			metrics = append(metrics, metricInfo{
				name:           metricName,
				description:    desc,
				unit:           "count",
				instrumentType: string(domain.InstrumentCounter),
			})
		} else if entry, ok := histogramDefs[constName]; ok {
			desc := ""
			if m := descriptionPattern.FindStringSubmatch(entry); len(m) >= 2 {
				desc = m[1]
			}
			unit := "count"
			if m := unitPattern.FindStringSubmatch(entry); len(m) >= 2 {
				unit = m[1]
			}
			metrics = append(metrics, metricInfo{
				name:           metricName,
				description:    desc,
				unit:           unit,
				instrumentType: string(domain.InstrumentHistogram),
			})
		}
	}

	return metrics
}

// findDefinitionBlockEntries finds all entries within blocks whose name ends with suffix.
// It handles COUNTER_DEFINITIONS, PERFORMANCE_COUNTER_DEFINITIONS, HISTOGRAM_DEFINITIONS, etc.
func findDefinitionBlockEntries(text string, suffix string) map[string]string {
	entries := make(map[string]string)

	// Find all blocks matching: const <NAME>_DEFINITIONS = { ... } as const;
	// or const <NAME> = { ... } as const; where name ends with suffix
	blockPattern := regexp.MustCompile(`const\s+\w*` + regexp.QuoteMeta(suffix) + `\s*=\s*\{`)

	for _, loc := range blockPattern.FindAllStringIndex(text, -1) {
		blockStart := loc[1] - 1 // Start at the opening brace
		blockEnd := findMatchingBrace(text, blockStart)
		if blockEnd < 0 {
			continue
		}
		blockContent := text[blockStart:blockEnd]

		// Find entries like: [CONST_NAME]: { ... },
		entryPattern := regexp.MustCompile(`\[(\w+)\]\s*:\s*\{`)
		for _, entryLoc := range entryPattern.FindAllStringSubmatchIndex(blockContent, -1) {
			constName := blockContent[entryLoc[2]:entryLoc[3]]
			entryStart := entryLoc[1] - 1 // Opening brace of entry
			entryEnd := findMatchingBrace(blockContent, entryStart)
			if entryEnd < 0 {
				continue
			}
			entries[constName] = blockContent[entryStart:entryEnd]
		}
	}

	return entries
}

func findMatchingBrace(text string, start int) int {
	if start >= len(text) || text[start] != '{' {
		return -1
	}
	depth := 0
	for i := start; i < len(text); i++ {
		switch text[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i + 1
			}
		}
	}
	return -1
}
