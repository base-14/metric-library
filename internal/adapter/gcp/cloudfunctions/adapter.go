package cloudfunctions

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
	return "gcp-cloudfunctions"
}

func (a *Adapter) SourceCategory() domain.SourceCategory {
	return domain.SourceCloud
}

func (a *Adapter) Confidence() domain.ConfidenceLevel {
	return domain.ConfidenceDocumented
}

func (a *Adapter) ExtractionMethod() domain.ExtractionMethod {
	return domain.ExtractionScrape
}

func (a *Adapter) RepoURL() string {
	return "https://cloud.google.com/monitoring/api/metrics_gcp#gcp-cloudfunctions"
}

func (a *Adapter) Fetch(_ context.Context, _ adapter.FetchOptions) (*adapter.FetchResult, error) {
	return &adapter.FetchResult{
		Commit:    time.Now().Format("2006-01-02"),
		Timestamp: time.Now(),
	}, nil
}

func (a *Adapter) Extract(_ context.Context, _ *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	return []*adapter.RawMetric{
		{Name: "cloudfunctions.googleapis.com/function/execution_count", Description: "Count of function executions broken down by status", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Functions", ComponentType: "platform", SourceLocation: "cloudfunctions.googleapis.com"},
		{Name: "cloudfunctions.googleapis.com/function/execution_times", Description: "Distribution of functions execution times in nanoseconds", Unit: "ns", InstrumentType: "histogram", ComponentName: "Cloud Functions", ComponentType: "platform", SourceLocation: "cloudfunctions.googleapis.com"},
		{Name: "cloudfunctions.googleapis.com/function/user_memory_bytes", Description: "Distribution of each function's working set of memory during execution in bytes", Unit: "By", InstrumentType: "histogram", ComponentName: "Cloud Functions", ComponentType: "platform", SourceLocation: "cloudfunctions.googleapis.com"},
		{Name: "cloudfunctions.googleapis.com/function/network_egress", Description: "Delta of outgoing network traffic in bytes", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud Functions", ComponentType: "platform", SourceLocation: "cloudfunctions.googleapis.com"},
		{Name: "cloudfunctions.googleapis.com/function/active_instances", Description: "Number of active function instances", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud Functions", ComponentType: "platform", SourceLocation: "cloudfunctions.googleapis.com"},
		{Name: "cloudfunctions.googleapis.com/function/instance_count", Description: "Number of function instances, broken down by state (active, idle)", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud Functions", ComponentType: "platform", SourceLocation: "cloudfunctions.googleapis.com"},
	}, nil
}
