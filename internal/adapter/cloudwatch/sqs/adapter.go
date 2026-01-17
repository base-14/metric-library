package sqs

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
	return "cloudwatch-sqs"
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
	return "https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-available-cloudwatch-metrics.html"
}

func (a *Adapter) Fetch(_ context.Context, _ adapter.FetchOptions) (*adapter.FetchResult, error) {
	return &adapter.FetchResult{
		Commit:    time.Now().Format("2006-01-02"),
		Timestamp: time.Now(),
	}, nil
}

func (a *Adapter) Extract(_ context.Context, _ *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	return []*adapter.RawMetric{
		// Standard Metrics
		{Name: "ApproximateAgeOfOldestMessage", Description: "The age of the oldest unprocessed message in the queue", Unit: "Seconds", InstrumentType: "gauge", ComponentName: "SQS", ComponentType: "platform", SourceLocation: "AWS/SQS"},
		{Name: "ApproximateNumberOfGroupsWithInflightMessages", Description: "For FIFO queues, the number of message groups with one or more in-flight messages", Unit: "Count", InstrumentType: "gauge", ComponentName: "SQS", ComponentType: "platform", SourceLocation: "AWS/SQS"},
		{Name: "ApproximateNumberOfMessagesDelayed", Description: "The number of messages in the queue that are delayed and not immediately available for retrieval", Unit: "Count", InstrumentType: "gauge", ComponentName: "SQS", ComponentType: "platform", SourceLocation: "AWS/SQS"},
		{Name: "ApproximateNumberOfMessagesNotVisible", Description: "The number of in-flight messages that have been received but not yet deleted or expired", Unit: "Count", InstrumentType: "gauge", ComponentName: "SQS", ComponentType: "platform", SourceLocation: "AWS/SQS"},
		{Name: "ApproximateNumberOfMessagesVisible", Description: "The number of messages currently available for retrieval and processing", Unit: "Count", InstrumentType: "gauge", ComponentName: "SQS", ComponentType: "platform", SourceLocation: "AWS/SQS"},
		{Name: "NumberOfEmptyReceives", Description: "The number of ReceiveMessage API calls that returned no messages", Unit: "Count", InstrumentType: "counter", ComponentName: "SQS", ComponentType: "platform", SourceLocation: "AWS/SQS"},
		{Name: "NumberOfDeduplicatedSentMessages", Description: "For FIFO queues, the number of sent messages that were deduplicated and not added to the queue", Unit: "Count", InstrumentType: "counter", ComponentName: "SQS", ComponentType: "platform", SourceLocation: "AWS/SQS"},
		{Name: "NumberOfMessagesDeleted", Description: "The number of messages successfully deleted from the queue", Unit: "Count", InstrumentType: "counter", ComponentName: "SQS", ComponentType: "platform", SourceLocation: "AWS/SQS"},
		{Name: "NumberOfMessagesReceived", Description: "The number of messages returned by the ReceiveMessage API", Unit: "Count", InstrumentType: "counter", ComponentName: "SQS", ComponentType: "platform", SourceLocation: "AWS/SQS"},
		{Name: "NumberOfMessagesSent", Description: "The number of messages successfully added to a queue", Unit: "Count", InstrumentType: "counter", ComponentName: "SQS", ComponentType: "platform", SourceLocation: "AWS/SQS"},
		{Name: "SentMessageSize", Description: "The size of messages successfully sent to the queue", Unit: "Bytes", InstrumentType: "gauge", ComponentName: "SQS", ComponentType: "platform", SourceLocation: "AWS/SQS"},
		// Fair Queue Metrics
		{Name: "ApproximateNumberOfNoisyGroups", Description: "The number of message groups that are considered noisy in a fair queue", Unit: "Count", InstrumentType: "gauge", ComponentName: "SQS", ComponentType: "platform", SourceLocation: "AWS/SQS"},
		{Name: "ApproximateNumberOfMessagesVisibleInQuietGroups", Description: "The number of messages visible excluding messages from noisy message groups", Unit: "Count", InstrumentType: "gauge", ComponentName: "SQS", ComponentType: "platform", SourceLocation: "AWS/SQS"},
		{Name: "ApproximateNumberOfMessagesNotVisibleInQuietGroups", Description: "The number of messages in-flight excluding messages from noisy message groups", Unit: "Count", InstrumentType: "gauge", ComponentName: "SQS", ComponentType: "platform", SourceLocation: "AWS/SQS"},
		{Name: "ApproximateNumberOfMessagesDelayedInQuietGroups", Description: "The number of delayed messages excluding messages from noisy message groups", Unit: "Count", InstrumentType: "gauge", ComponentName: "SQS", ComponentType: "platform", SourceLocation: "AWS/SQS"},
		{Name: "ApproximateAgeOfOldestMessageInQuietGroups", Description: "The age of the oldest non-deleted message excluding messages from noisy message groups", Unit: "Seconds", InstrumentType: "gauge", ComponentName: "SQS", ComponentType: "platform", SourceLocation: "AWS/SQS"},
	}, nil
}
