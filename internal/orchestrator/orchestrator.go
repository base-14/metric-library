package orchestrator

import (
	"context"
	"fmt"
	"time"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
	"github.com/base-14/metric-library/internal/store"
)

type Adapter interface {
	Name() string
	Fetch(ctx context.Context, opts adapter.FetchOptions) (*adapter.FetchResult, error)
	Extract(ctx context.Context, result *adapter.FetchResult) ([]*adapter.RawMetric, error)
	SourceCategory() domain.SourceCategory
	Confidence() domain.ConfidenceLevel
	ExtractionMethod() domain.ExtractionMethod
	RepoURL() string
}

type Options struct {
	Commit   string
	CacheDir string
	Force    bool
}

type Result struct {
	AdapterName      string
	Commit           string
	MetricsExtracted int
	MetricsStored    int
	Duration         time.Duration
}

type Extractor struct {
	adapter Adapter
	store   store.Store
}

func NewExtractor(adp Adapter, st store.Store) *Extractor {
	return &Extractor{
		adapter: adp,
		store:   st,
	}
}

func (e *Extractor) Run(ctx context.Context, opts Options) (*Result, error) {
	startTime := time.Now()

	run := &store.ExtractionRun{
		ID:          fmt.Sprintf("%s-%d", e.adapter.Name(), startTime.UnixNano()),
		AdapterName: e.adapter.Name(),
		StartedAt:   startTime,
		Status:      "running",
	}
	if err := e.store.CreateExtractionRun(ctx, run); err != nil {
		return nil, fmt.Errorf("failed to create extraction run: %w", err)
	}

	fetchOpts := adapter.FetchOptions{
		Commit:   opts.Commit,
		CacheDir: opts.CacheDir,
		Force:    opts.Force,
	}

	fetchResult, err := e.adapter.Fetch(ctx, fetchOpts)
	if err != nil {
		run.Status = "failed"
		run.ErrorMessage = err.Error()
		completedAt := time.Now()
		run.CompletedAt = &completedAt
		_ = e.store.UpdateExtractionRun(ctx, run)
		return nil, fmt.Errorf("fetch failed: %w", err)
	}

	run.Commit = fetchResult.Commit

	rawMetrics, err := e.adapter.Extract(ctx, fetchResult)
	if err != nil {
		run.Status = "failed"
		run.ErrorMessage = err.Error()
		completedAt := time.Now()
		run.CompletedAt = &completedAt
		_ = e.store.UpdateExtractionRun(ctx, run)
		return nil, fmt.Errorf("extraction failed: %w", err)
	}

	canonicalMetrics := make([]*domain.CanonicalMetric, 0, len(rawMetrics))
	for _, raw := range rawMetrics {
		canonical := e.convertToCanonical(raw, fetchResult)
		if err := canonical.Validate(); err != nil {
			continue
		}
		canonical.EnsureID()
		canonicalMetrics = append(canonicalMetrics, canonical)
	}

	if err := e.store.UpsertMetrics(ctx, canonicalMetrics); err != nil {
		run.Status = "failed"
		run.ErrorMessage = err.Error()
		completedAt := time.Now()
		run.CompletedAt = &completedAt
		_ = e.store.UpdateExtractionRun(ctx, run)
		return nil, fmt.Errorf("failed to store metrics: %w", err)
	}

	completedAt := time.Now()
	run.CompletedAt = &completedAt
	run.MetricsCount = len(canonicalMetrics)
	run.Status = "completed"
	_ = e.store.UpdateExtractionRun(ctx, run)

	return &Result{
		AdapterName:      e.adapter.Name(),
		Commit:           fetchResult.Commit,
		MetricsExtracted: len(rawMetrics),
		MetricsStored:    len(canonicalMetrics),
		Duration:         time.Since(startTime),
	}, nil
}

func (e *Extractor) convertToCanonical(raw *adapter.RawMetric, fetchResult *adapter.FetchResult) *domain.CanonicalMetric {
	return &domain.CanonicalMetric{
		MetricName:       raw.Name,
		InstrumentType:   domain.InstrumentType(raw.InstrumentType),
		Description:      raw.Description,
		Unit:             raw.Unit,
		Attributes:       raw.Attributes,
		EnabledByDefault: raw.EnabledByDefault,
		ComponentType:    domain.ComponentType(raw.ComponentType),
		ComponentName:    raw.ComponentName,
		SourceCategory:   e.adapter.SourceCategory(),
		SourceName:       e.adapter.Name(),
		SourceLocation:   raw.SourceLocation,
		ExtractionMethod: e.adapter.ExtractionMethod(),
		SourceConfidence: e.adapter.Confidence(),
		Repo:             e.adapter.RepoURL(),
		Path:             raw.Path,
		Commit:           fetchResult.Commit,
		ExtractedAt:      fetchResult.Timestamp,
	}
}
