package lambda

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
	return "cloudwatch-lambda"
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
	return "https://docs.aws.amazon.com/lambda/latest/dg/monitoring-metrics-types.html"
}

func (a *Adapter) Fetch(_ context.Context, _ adapter.FetchOptions) (*adapter.FetchResult, error) {
	return &adapter.FetchResult{
		Commit:    time.Now().Format("2006-01-02"),
		Timestamp: time.Now(),
	}, nil
}

func (a *Adapter) Extract(_ context.Context, _ *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	metrics := []*adapter.RawMetric{}
	metrics = append(metrics, invocationMetrics()...)
	metrics = append(metrics, performanceMetrics()...)
	metrics = append(metrics, concurrencyMetrics()...)
	metrics = append(metrics, asyncMetrics()...)
	metrics = append(metrics, eventSourceMetrics()...)
	return metrics, nil
}

func invocationMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "Invocations", Description: "Number of times function code is invoked, including successful invocations and those resulting in errors", Unit: "Count", InstrumentType: "counter", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "Errors", Description: "Number of invocations that result in a function error", Unit: "Count", InstrumentType: "counter", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "DeadLetterErrors", Description: "Number of failed attempts to send events to a dead-letter queue (DLQ) for async invocations", Unit: "Count", InstrumentType: "counter", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "DestinationDeliveryFailures", Description: "Number of failed attempts to send events to a destination for async invocation and event source mappings", Unit: "Count", InstrumentType: "counter", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "Throttles", Description: "Number of invocation requests that are throttled", Unit: "Count", InstrumentType: "counter", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "OversizedRecordCount", Description: "Number of events over 6 MB from DocumentDB change streams", Unit: "Count", InstrumentType: "counter", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "ProvisionedConcurrencyInvocations", Description: "Number of invocations using provisioned concurrency", Unit: "Count", InstrumentType: "counter", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "ProvisionedConcurrencySpilloverInvocations", Description: "Number of invocations using standard concurrency when provisioned concurrency is exhausted", Unit: "Count", InstrumentType: "counter", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "RecursiveInvocationsDropped", Description: "Number of invocations stopped due to detected infinite recursive loops", Unit: "Count", InstrumentType: "counter", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "SignatureValidationErrors", Description: "Number of code package deployments with signature validation failures", Unit: "Count", InstrumentType: "counter", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
	}
}

func performanceMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "Duration", Description: "Time function code spends processing an event", Unit: "Milliseconds", InstrumentType: "gauge", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "PostRuntimeExtensionsDuration", Description: "Cumulative time runtime spends executing extension code after function completion", Unit: "Milliseconds", InstrumentType: "gauge", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "IteratorAge", Description: "Age of the last record in the event for stream-based sources", Unit: "Milliseconds", InstrumentType: "gauge", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "OffsetLag", Description: "Offset lag for self-managed Kafka and Amazon MSK event sources", Unit: "Count", InstrumentType: "gauge", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
	}
}

func concurrencyMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "ConcurrentExecutions", Description: "Number of function instances processing events", Unit: "Count", InstrumentType: "gauge", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "ProvisionedConcurrentExecutions", Description: "Number of function instances processing events using provisioned concurrency", Unit: "Count", InstrumentType: "gauge", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "ProvisionedConcurrencyUtilization", Description: "Ratio of ProvisionedConcurrentExecutions to total provisioned concurrency", Unit: "Percent", InstrumentType: "gauge", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "UnreservedConcurrentExecutions", Description: "Number of events processed by functions without reserved concurrency", Unit: "Count", InstrumentType: "gauge", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "ClaimedAccountConcurrency", Description: "Concurrency unavailable for on-demand invocations at the Region level", Unit: "Count", InstrumentType: "gauge", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
	}
}

func asyncMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "AsyncEventsReceived", Description: "Number of events successfully queued for processing", Unit: "Count", InstrumentType: "counter", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "AsyncEventAge", Description: "Time between event queuing and function invocation", Unit: "Milliseconds", InstrumentType: "gauge", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "AsyncEventsDropped", Description: "Number of events dropped without executing the function", Unit: "Count", InstrumentType: "counter", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
	}
}

func eventSourceMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "PolledEventCount", Description: "Number of events read from event source", Unit: "Count", InstrumentType: "counter", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "FilteredOutEventCount", Description: "Number of events filtered out by filter criteria", Unit: "Count", InstrumentType: "counter", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "InvokedEventCount", Description: "Number of events that invoked the function", Unit: "Count", InstrumentType: "counter", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "FailedInvokeEventCount", Description: "Number of events that failed to invoke the function", Unit: "Count", InstrumentType: "counter", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "DroppedEventCount", Description: "Number of events dropped due to expiry or retry exhaustion", Unit: "Count", InstrumentType: "counter", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "OnFailureDestinationDeliveredEventCount", Description: "Number of events sent to on-failure destination", Unit: "Count", InstrumentType: "counter", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "DeletedEventCount", Description: "Number of events successfully deleted after processing", Unit: "Count", InstrumentType: "counter", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
		{Name: "ProvisionedPollers", Description: "Number of active event pollers in provisioned mode", Unit: "Count", InstrumentType: "gauge", ComponentName: "Lambda", ComponentType: "platform", SourceLocation: "AWS/Lambda"},
	}
}
