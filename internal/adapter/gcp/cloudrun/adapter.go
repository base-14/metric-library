package cloudrun

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
	return "gcp-cloudrun"
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
	return "https://cloud.google.com/monitoring/api/metrics_gcp#gcp-run"
}

func (a *Adapter) Fetch(_ context.Context, _ adapter.FetchOptions) (*adapter.FetchResult, error) {
	return &adapter.FetchResult{
		Commit:    time.Now().Format("2006-01-02"),
		Timestamp: time.Now(),
	}, nil
}

func (a *Adapter) Extract(_ context.Context, _ *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	var metrics []*adapter.RawMetric
	metrics = append(metrics, requestMetrics()...)
	metrics = append(metrics, containerMetrics()...)
	return metrics, nil
}

func requestMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "run.googleapis.com/request_count", Description: "Number of requests reaching the revision", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Run", ComponentType: "platform", SourceLocation: "run.googleapis.com"},
		{Name: "run.googleapis.com/request_latencies", Description: "Distribution of request latency in milliseconds reaching the revision", Unit: "ms", InstrumentType: "histogram", ComponentName: "Cloud Run", ComponentType: "platform", SourceLocation: "run.googleapis.com"},
		{Name: "run.googleapis.com/response_latencies", Description: "Distribution of response latency in milliseconds for requests to the revision", Unit: "ms", InstrumentType: "histogram", ComponentName: "Cloud Run", ComponentType: "platform", SourceLocation: "run.googleapis.com"},
		{Name: "run.googleapis.com/pending_queue/pending_requests", Description: "Number of requests waiting in the pending queue", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud Run", ComponentType: "platform", SourceLocation: "run.googleapis.com"},
	}
}

func containerMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "run.googleapis.com/container/cpu/utilizations", Description: "CPU utilization of the container instance divided by the container cpu limit", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud Run", ComponentType: "platform", SourceLocation: "run.googleapis.com"},
		{Name: "run.googleapis.com/container/cpu/allocation_time", Description: "CPU allocation of the container instance in seconds", Unit: "s", InstrumentType: "counter", ComponentName: "Cloud Run", ComponentType: "platform", SourceLocation: "run.googleapis.com"},
		{Name: "run.googleapis.com/container/memory/utilizations", Description: "Memory utilization of the container instance divided by the container memory limit", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud Run", ComponentType: "platform", SourceLocation: "run.googleapis.com"},
		{Name: "run.googleapis.com/container/memory/allocation", Description: "Memory allocation of the container instance in MiB", Unit: "MiBy", InstrumentType: "gauge", ComponentName: "Cloud Run", ComponentType: "platform", SourceLocation: "run.googleapis.com"},
		{Name: "run.googleapis.com/container/startup_latencies", Description: "Time from start of container to first request for new instances", Unit: "ms", InstrumentType: "histogram", ComponentName: "Cloud Run", ComponentType: "platform", SourceLocation: "run.googleapis.com"},
		{Name: "run.googleapis.com/container/network/received_bytes_count", Description: "Count of bytes received by the container instance from the network", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud Run", ComponentType: "platform", SourceLocation: "run.googleapis.com"},
		{Name: "run.googleapis.com/container/network/sent_bytes_count", Description: "Count of bytes sent by the container instance over the network", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud Run", ComponentType: "platform", SourceLocation: "run.googleapis.com"},
		{Name: "run.googleapis.com/container/instance_count", Description: "Number of container instances that exist, broken down by state", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud Run", ComponentType: "platform", SourceLocation: "run.googleapis.com"},
		{Name: "run.googleapis.com/container/max_request_concurrencies", Description: "Maximum number of concurrent requests being served by each container instance", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud Run", ComponentType: "platform", SourceLocation: "run.googleapis.com"},
		{Name: "run.googleapis.com/container/billable_instance_time", Description: "Billable time aggregated from all container instances", Unit: "s", InstrumentType: "counter", ComponentName: "Cloud Run", ComponentType: "platform", SourceLocation: "run.googleapis.com"},
	}
}
