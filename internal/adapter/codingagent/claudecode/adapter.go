package claudecode

import (
	"context"
	"time"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
)

type Adapter struct{}

func NewAdapter(_ string) *Adapter {
	return &Adapter{}
}

func (a *Adapter) Name() string {
	return "codingagent-claude-code"
}

func (a *Adapter) SourceCategory() domain.SourceCategory {
	return domain.SourceCodingAgent
}

func (a *Adapter) Confidence() domain.ConfidenceLevel {
	return domain.ConfidenceDocumented
}

func (a *Adapter) ExtractionMethod() domain.ExtractionMethod {
	return domain.ExtractionMetadata
}

func (a *Adapter) RepoURL() string {
	return "https://github.com/anthropics/claude-code-monitoring-guide"
}

func (a *Adapter) Fetch(_ context.Context, _ adapter.FetchOptions) (*adapter.FetchResult, error) {
	return &adapter.FetchResult{
		Commit:    time.Now().Format("2006-01-02"),
		Timestamp: time.Now(),
	}, nil
}

func (a *Adapter) Extract(_ context.Context, _ *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	return []*adapter.RawMetric{
		{
			Name:             "claude_code.session.count",
			Description:      "Number of Claude Code sessions started",
			Unit:             "count",
			InstrumentType:   string(domain.InstrumentCounter),
			EnabledByDefault: true,
			ComponentType:    string(domain.ComponentPlatform),
			ComponentName:    "claude-code",
		},
		{
			Name:             "claude_code.lines_of_code.count",
			Description:      "Number of lines of code added or removed",
			Unit:             "count",
			InstrumentType:   string(domain.InstrumentCounter),
			EnabledByDefault: true,
			ComponentType:    string(domain.ComponentPlatform),
			ComponentName:    "claude-code",
			Attributes: []domain.Attribute{
				{Name: "type", Type: "string", Description: "Type of change (added, removed)"},
			},
		},
		{
			Name:             "claude_code.pull_request.count",
			Description:      "Number of pull requests created",
			Unit:             "count",
			InstrumentType:   string(domain.InstrumentCounter),
			EnabledByDefault: true,
			ComponentType:    string(domain.ComponentPlatform),
			ComponentName:    "claude-code",
		},
		{
			Name:             "claude_code.commit.count",
			Description:      "Number of commits created",
			Unit:             "count",
			InstrumentType:   string(domain.InstrumentCounter),
			EnabledByDefault: true,
			ComponentType:    string(domain.ComponentPlatform),
			ComponentName:    "claude-code",
		},
		{
			Name:             "claude_code.cost.usage",
			Description:      "Cost of API usage in USD",
			Unit:             "USD",
			InstrumentType:   string(domain.InstrumentCounter),
			EnabledByDefault: true,
			ComponentType:    string(domain.ComponentPlatform),
			ComponentName:    "claude-code",
			Attributes: []domain.Attribute{
				{Name: "model", Type: "string", Description: "Model used for the request"},
			},
		},
		{
			Name:             "claude_code.token.usage",
			Description:      "Number of tokens consumed",
			Unit:             "tokens",
			InstrumentType:   string(domain.InstrumentCounter),
			EnabledByDefault: true,
			ComponentType:    string(domain.ComponentPlatform),
			ComponentName:    "claude-code",
			Attributes: []domain.Attribute{
				{Name: "type", Type: "string", Description: "Token type (input, output)"},
				{Name: "model", Type: "string", Description: "Model used for the request"},
			},
		},
		{
			Name:             "claude_code.code_edit_tool.decision",
			Description:      "Tool invocation decisions made during coding",
			Unit:             "count",
			InstrumentType:   string(domain.InstrumentCounter),
			EnabledByDefault: true,
			ComponentType:    string(domain.ComponentPlatform),
			ComponentName:    "claude-code",
			Attributes: []domain.Attribute{
				{Name: "tool", Type: "string", Description: "Tool name"},
				{Name: "decision", Type: "string", Description: "Decision outcome (accepted, rejected)"},
				{Name: "language", Type: "string", Description: "Programming language"},
			},
		},
		{
			Name:             "claude_code.active_time.total",
			Description:      "Total active time spent in sessions",
			Unit:             "s",
			InstrumentType:   string(domain.InstrumentCounter),
			EnabledByDefault: true,
			ComponentType:    string(domain.ComponentPlatform),
			ComponentName:    "claude-code",
		},
	}, nil
}
