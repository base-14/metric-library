package servicebus

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
	return "azure-servicebus"
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
	return "https://learn.microsoft.com/en-us/azure/azure-monitor/reference/supported-metrics/microsoft-servicebus-namespaces-metrics"
}

func (a *Adapter) Fetch(_ context.Context, _ adapter.FetchOptions) (*adapter.FetchResult, error) {
	return &adapter.FetchResult{
		Commit:    time.Now().Format("2006-01-02"),
		Timestamp: time.Now(),
	}, nil
}

func (a *Adapter) Extract(_ context.Context, _ *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	var metrics []*adapter.RawMetric
	metrics = append(metrics, messageMetrics()...)
	metrics = append(metrics, requestMetrics()...)
	metrics = append(metrics, resourceMetrics()...)
	return metrics, nil
}

func messageMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.servicebus.incoming_messages", Description: "Count of incoming messages for a namespace or entity", Unit: "1", InstrumentType: "counter", ComponentName: "Service Bus", ComponentType: "platform", SourceLocation: "Microsoft.ServiceBus/namespaces"},
		{Name: "azure.servicebus.outgoing_messages", Description: "Count of outgoing messages for a namespace or entity", Unit: "1", InstrumentType: "counter", ComponentName: "Service Bus", ComponentType: "platform", SourceLocation: "Microsoft.ServiceBus/namespaces"},
		{Name: "azure.servicebus.active_messages", Description: "Count of active messages in a Queue/Topic", Unit: "1", InstrumentType: "gauge", ComponentName: "Service Bus", ComponentType: "platform", SourceLocation: "Microsoft.ServiceBus/namespaces"},
		{Name: "azure.servicebus.deadlettered_messages", Description: "Count of dead-lettered messages in a Queue/Topic", Unit: "1", InstrumentType: "gauge", ComponentName: "Service Bus", ComponentType: "platform", SourceLocation: "Microsoft.ServiceBus/namespaces"},
		{Name: "azure.servicebus.scheduled_messages", Description: "Count of scheduled messages in a Queue/Topic", Unit: "1", InstrumentType: "gauge", ComponentName: "Service Bus", ComponentType: "platform", SourceLocation: "Microsoft.ServiceBus/namespaces"},
		{Name: "azure.servicebus.messages", Description: "Count of messages in a Queue/Topic", Unit: "1", InstrumentType: "gauge", ComponentName: "Service Bus", ComponentType: "platform", SourceLocation: "Microsoft.ServiceBus/namespaces"},
		{Name: "azure.servicebus.completed_messages", Description: "Count of messages completed on a Queue/Topic", Unit: "1", InstrumentType: "counter", ComponentName: "Service Bus", ComponentType: "platform", SourceLocation: "Microsoft.ServiceBus/namespaces"},
		{Name: "azure.servicebus.abandoned_messages", Description: "Count of messages abandoned on a Queue/Topic", Unit: "1", InstrumentType: "counter", ComponentName: "Service Bus", ComponentType: "platform", SourceLocation: "Microsoft.ServiceBus/namespaces"},
		{Name: "azure.servicebus.size", Description: "Size of a Queue/Topic in bytes", Unit: "By", InstrumentType: "gauge", ComponentName: "Service Bus", ComponentType: "platform", SourceLocation: "Microsoft.ServiceBus/namespaces"},
	}
}

func requestMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.servicebus.incoming_requests", Description: "Count of incoming requests for a namespace", Unit: "1", InstrumentType: "counter", ComponentName: "Service Bus", ComponentType: "platform", SourceLocation: "Microsoft.ServiceBus/namespaces"},
		{Name: "azure.servicebus.successful_requests", Description: "Count of successful requests for a namespace", Unit: "1", InstrumentType: "counter", ComponentName: "Service Bus", ComponentType: "platform", SourceLocation: "Microsoft.ServiceBus/namespaces"},
		{Name: "azure.servicebus.server_errors", Description: "Count of server errors for a namespace", Unit: "1", InstrumentType: "counter", ComponentName: "Service Bus", ComponentType: "platform", SourceLocation: "Microsoft.ServiceBus/namespaces"},
		{Name: "azure.servicebus.user_errors", Description: "Count of user errors for a namespace", Unit: "1", InstrumentType: "counter", ComponentName: "Service Bus", ComponentType: "platform", SourceLocation: "Microsoft.ServiceBus/namespaces"},
		{Name: "azure.servicebus.throttled_requests", Description: "Count of throttled requests for a namespace", Unit: "1", InstrumentType: "counter", ComponentName: "Service Bus", ComponentType: "platform", SourceLocation: "Microsoft.ServiceBus/namespaces"},
		{Name: "azure.servicebus.server_send_latency", Description: "Latency of Send message operations for Service Bus resources", Unit: "ms", InstrumentType: "gauge", ComponentName: "Service Bus", ComponentType: "platform", SourceLocation: "Microsoft.ServiceBus/namespaces"},
	}
}

func resourceMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.servicebus.active_connections", Description: "Total active connections for a namespace", Unit: "1", InstrumentType: "gauge", ComponentName: "Service Bus", ComponentType: "platform", SourceLocation: "Microsoft.ServiceBus/namespaces"},
		{Name: "azure.servicebus.connections_opened", Description: "Count of connections opened for a namespace", Unit: "1", InstrumentType: "counter", ComponentName: "Service Bus", ComponentType: "platform", SourceLocation: "Microsoft.ServiceBus/namespaces"},
		{Name: "azure.servicebus.connections_closed", Description: "Count of connections closed for a namespace", Unit: "1", InstrumentType: "counter", ComponentName: "Service Bus", ComponentType: "platform", SourceLocation: "Microsoft.ServiceBus/namespaces"},
		{Name: "azure.servicebus.namespace_cpu_usage", Description: "Service Bus premium namespace CPU usage metric", Unit: "%", InstrumentType: "gauge", ComponentName: "Service Bus", ComponentType: "platform", SourceLocation: "Microsoft.ServiceBus/namespaces"},
		{Name: "azure.servicebus.namespace_memory_usage", Description: "Service Bus premium namespace memory usage metric", Unit: "%", InstrumentType: "gauge", ComponentName: "Service Bus", ComponentType: "platform", SourceLocation: "Microsoft.ServiceBus/namespaces"},
	}
}
