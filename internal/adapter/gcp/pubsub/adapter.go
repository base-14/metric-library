package pubsub

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
	return "gcp-pubsub"
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
	return "https://cloud.google.com/monitoring/api/metrics_gcp#gcp-pubsub"
}

func (a *Adapter) Fetch(_ context.Context, _ adapter.FetchOptions) (*adapter.FetchResult, error) {
	return &adapter.FetchResult{
		Commit:    time.Now().Format("2006-01-02"),
		Timestamp: time.Now(),
	}, nil
}

func (a *Adapter) Extract(_ context.Context, _ *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	var metrics []*adapter.RawMetric
	metrics = append(metrics, topicMetrics()...)
	metrics = append(metrics, subscriptionMetrics()...)
	metrics = append(metrics, snapshotMetrics()...)
	return metrics, nil
}

func topicMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "pubsub.googleapis.com/topic/send_message_operation_count", Description: "Cumulative count of publish message operations grouped by result", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/topic/send_request_count", Description: "Cumulative count of publish requests grouped by result", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/topic/byte_cost", Description: "Cost of operations in bytes, used to measure quota utilization", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/topic/message_sizes", Description: "Distribution of publish message sizes in bytes", Unit: "By", InstrumentType: "histogram", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/topic/config_updates_count", Description: "Cumulative count of configuration changes for each topic grouped by operation type and result", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/topic/oldest_unacked_message_age_by_region", Description: "Age in seconds of the oldest unacknowledged message in a topic by region", Unit: "s", InstrumentType: "gauge", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/topic/num_unacked_messages_by_region", Description: "Number of unacknowledged messages in a topic by region", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/topic/num_retained_acked_messages_by_region", Description: "Number of acknowledged messages retained in a topic by region", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
	}
}

func subscriptionMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "pubsub.googleapis.com/subscription/pull_request_count", Description: "Cumulative count of pull requests grouped by result", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/subscription/pull_message_operation_count", Description: "Cumulative count of pull message operations grouped by result", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/subscription/streaming_pull_response_count", Description: "Cumulative count of streaming pull responses grouped by result", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/subscription/push_request_count", Description: "Cumulative count of push attempts grouped by result", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/subscription/push_request_latencies", Description: "Distribution of push request latencies in microseconds", Unit: "us", InstrumentType: "histogram", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/subscription/sent_message_count", Description: "Cumulative count of messages sent by Pub/Sub to subscriber clients grouped by delivery type", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/subscription/byte_cost", Description: "Cumulative cost of operations in bytes, used to measure quota utilization", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/subscription/backlog_bytes", Description: "Total byte size of unacknowledged messages (backlog messages) in a subscription", Unit: "By", InstrumentType: "gauge", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/subscription/num_undelivered_messages", Description: "Number of unacknowledged messages (backlog messages) in a subscription", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/subscription/num_outstanding_messages", Description: "Number of messages delivered to a subscription's push endpoint but not yet acknowledged", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/subscription/oldest_unacked_message_age", Description: "Age in seconds of the oldest unacknowledged message in a subscription", Unit: "s", InstrumentType: "gauge", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/subscription/ack_message_count", Description: "Cumulative count of messages acknowledged by Acknowledge requests grouped by delivery type", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/subscription/modify_ack_deadline_message_operation_count", Description: "Cumulative count of ModifyAckDeadline message operations grouped by result", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/subscription/dead_letter_message_count", Description: "Cumulative count of messages published to dead letter topic, grouped by result", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/subscription/config_updates_count", Description: "Cumulative count of configuration changes for each subscription grouped by operation type and result", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/subscription/ack_latencies", Description: "Distribution of ack latencies in milliseconds, from when Pub/Sub sends a message to a subscriber client until Pub/Sub receives an Acknowledge", Unit: "ms", InstrumentType: "histogram", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
	}
}

func snapshotMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "pubsub.googleapis.com/snapshot/backlog_bytes", Description: "Total byte size of messages retained in a snapshot", Unit: "By", InstrumentType: "gauge", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/snapshot/backlog_bytes_by_region", Description: "Total byte size of messages retained in a snapshot by region", Unit: "By", InstrumentType: "gauge", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/snapshot/num_messages", Description: "Number of messages retained in a snapshot", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/snapshot/oldest_message_age", Description: "Age in seconds of the oldest message retained in a snapshot", Unit: "s", InstrumentType: "gauge", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
		{Name: "pubsub.googleapis.com/snapshot/config_updates_count", Description: "Cumulative count of configuration changes grouped by operation type and result", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Pub/Sub", ComponentType: "platform", SourceLocation: "pubsub.googleapis.com"},
	}
}
