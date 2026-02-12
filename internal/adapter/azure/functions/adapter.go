package functions

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
	return "azure-functions"
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
	return "https://learn.microsoft.com/en-us/azure/azure-monitor/reference/supported-metrics/microsoft-web-sites-metrics"
}

func (a *Adapter) Fetch(_ context.Context, _ adapter.FetchOptions) (*adapter.FetchResult, error) {
	return &adapter.FetchResult{
		Commit:    time.Now().Format("2006-01-02"),
		Timestamp: time.Now(),
	}, nil
}

func (a *Adapter) Extract(_ context.Context, _ *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	var metrics []*adapter.RawMetric
	metrics = append(metrics, executionMetrics()...)
	metrics = append(metrics, httpMetrics()...)
	return metrics, nil
}

func executionMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.functions.function_execution_count", Description: "Function execution count", Unit: "1", InstrumentType: "counter", ComponentName: "Azure Functions", ComponentType: "platform", SourceLocation: "Microsoft.Web/sites"},
		{Name: "azure.functions.function_execution_units", Description: "Function execution units in MB-milliseconds", Unit: "1", InstrumentType: "counter", ComponentName: "Azure Functions", ComponentType: "platform", SourceLocation: "Microsoft.Web/sites"},
		{Name: "azure.functions.app_connections", Description: "Number of bound sockets existing in the sandbox", Unit: "1", InstrumentType: "gauge", ComponentName: "Azure Functions", ComponentType: "platform", SourceLocation: "Microsoft.Web/sites"},
		{Name: "azure.functions.handles", Description: "Total number of handles currently open by the app process", Unit: "1", InstrumentType: "gauge", ComponentName: "Azure Functions", ComponentType: "platform", SourceLocation: "Microsoft.Web/sites"},
		{Name: "azure.functions.threads", Description: "Number of threads currently active in the app process", Unit: "1", InstrumentType: "gauge", ComponentName: "Azure Functions", ComponentType: "platform", SourceLocation: "Microsoft.Web/sites"},
		{Name: "azure.functions.private_bytes", Description: "The current size in bytes of memory that the app process has allocated that can't be shared with other processes", Unit: "By", InstrumentType: "gauge", ComponentName: "Azure Functions", ComponentType: "platform", SourceLocation: "Microsoft.Web/sites"},
		{Name: "azure.functions.io_read_bytes_per_second", Description: "The rate at which the app process is reading bytes from I/O operations", Unit: "By/s", InstrumentType: "gauge", ComponentName: "Azure Functions", ComponentType: "platform", SourceLocation: "Microsoft.Web/sites"},
		{Name: "azure.functions.io_write_bytes_per_second", Description: "The rate at which the app process is writing bytes to I/O operations", Unit: "By/s", InstrumentType: "gauge", ComponentName: "Azure Functions", ComponentType: "platform", SourceLocation: "Microsoft.Web/sites"},
		{Name: "azure.functions.io_other_bytes_per_second", Description: "The rate at which the app process is issuing bytes to I/O operations that don't involve data, such as control operations", Unit: "By/s", InstrumentType: "gauge", ComponentName: "Azure Functions", ComponentType: "platform", SourceLocation: "Microsoft.Web/sites"},
	}
}

func httpMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.functions.requests", Description: "Total number of requests regardless of their resulting HTTP status code", Unit: "1", InstrumentType: "counter", ComponentName: "Azure Functions", ComponentType: "platform", SourceLocation: "Microsoft.Web/sites"},
		{Name: "azure.functions.http2xx", Description: "Count of requests resulting in an HTTP status code >= 200 but < 300", Unit: "1", InstrumentType: "counter", ComponentName: "Azure Functions", ComponentType: "platform", SourceLocation: "Microsoft.Web/sites"},
		{Name: "azure.functions.http3xx", Description: "Count of requests resulting in an HTTP status code >= 300 but < 400", Unit: "1", InstrumentType: "counter", ComponentName: "Azure Functions", ComponentType: "platform", SourceLocation: "Microsoft.Web/sites"},
		{Name: "azure.functions.http4xx", Description: "Count of requests resulting in an HTTP status code >= 400 but < 500", Unit: "1", InstrumentType: "counter", ComponentName: "Azure Functions", ComponentType: "platform", SourceLocation: "Microsoft.Web/sites"},
		{Name: "azure.functions.http5xx", Description: "Count of requests resulting in an HTTP status code >= 500", Unit: "1", InstrumentType: "counter", ComponentName: "Azure Functions", ComponentType: "platform", SourceLocation: "Microsoft.Web/sites"},
		{Name: "azure.functions.http_response_time", Description: "Time taken for the app to serve requests in seconds", Unit: "s", InstrumentType: "gauge", ComponentName: "Azure Functions", ComponentType: "platform", SourceLocation: "Microsoft.Web/sites"},
		{Name: "azure.functions.average_response_time", Description: "Average time taken for the app to serve requests in seconds", Unit: "s", InstrumentType: "gauge", ComponentName: "Azure Functions", ComponentType: "platform", SourceLocation: "Microsoft.Web/sites"},
	}
}
