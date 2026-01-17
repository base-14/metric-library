package dynamodb

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
	return "cloudwatch-dynamodb"
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
	return "https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/metrics-dimensions.html"
}

func (a *Adapter) Fetch(_ context.Context, _ adapter.FetchOptions) (*adapter.FetchResult, error) {
	return &adapter.FetchResult{
		Commit:    time.Now().Format("2006-01-02"),
		Timestamp: time.Now(),
	}, nil
}

func (a *Adapter) Extract(_ context.Context, _ *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	return []*adapter.RawMetric{
		// Capacity Metrics
		{Name: "AccountMaxReads", Description: "Maximum read capacity units usable by an account", Unit: "Count", InstrumentType: "gauge", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "AccountMaxTableLevelReads", Description: "Maximum read capacity units usable by a table or GSI", Unit: "Count", InstrumentType: "gauge", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "AccountMaxTableLevelWrites", Description: "Maximum write capacity units usable by a table or GSI", Unit: "Count", InstrumentType: "gauge", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "AccountMaxWrites", Description: "Maximum write capacity units usable by an account", Unit: "Count", InstrumentType: "gauge", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "AccountProvisionedReadCapacityUtilization", Description: "Percentage of provisioned read capacity utilized by account", Unit: "Percent", InstrumentType: "gauge", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "AccountProvisionedWriteCapacityUtilization", Description: "Percentage of provisioned write capacity utilized by account", Unit: "Percent", InstrumentType: "gauge", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "ConsumedReadCapacityUnits", Description: "Read capacity units consumed over time period", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "ConsumedWriteCapacityUnits", Description: "Write capacity units consumed over time period", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "MaxProvisionedTableReadCapacityUtilization", Description: "Percentage of provisioned read capacity for highest table/GSI", Unit: "Percent", InstrumentType: "gauge", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "MaxProvisionedTableWriteCapacityUtilization", Description: "Percentage of provisioned write capacity for highest table/GSI", Unit: "Percent", InstrumentType: "gauge", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "OnDemandMaxReadRequestUnits", Description: "Specified on-demand read request units for table/GSI", Unit: "Count", InstrumentType: "gauge", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "OnDemandMaxWriteRequestUnits", Description: "Specified on-demand write request units for table/GSI", Unit: "Count", InstrumentType: "gauge", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "ProvisionedReadCapacityUnits", Description: "Provisioned read capacity units for table/GSI", Unit: "Count", InstrumentType: "gauge", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "ProvisionedWriteCapacityUnits", Description: "Provisioned write capacity units for table/GSI", Unit: "Count", InstrumentType: "gauge", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		// Request Metrics
		{Name: "ConditionalCheckFailedRequests", Description: "Number of failed conditional write attempts", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "SuccessfulRequestLatency", Description: "Latency of successful requests to DynamoDB", Unit: "Milliseconds", InstrumentType: "gauge", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "SystemErrors", Description: "Requests generating HTTP 500 status code", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "UserErrors", Description: "Requests generating HTTP 400 status code", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "ReturnedBytes", Description: "Bytes returned by GetRecords operations", Unit: "Bytes", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "ReturnedItemCount", Description: "Items returned by Query, Scan, or ExecuteStatement", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "ReturnedRecordsCount", Description: "Stream records returned by GetRecords operations", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		// Throttling Metrics
		{Name: "ReadThrottleEvents", Description: "Requests exceeding provisioned read capacity", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "WriteThrottleEvents", Description: "Requests exceeding provisioned write capacity", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "ThrottledRequests", Description: "Requests exceeding provisioned throughput limits", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "ReadAccountLimitThrottleEvents", Description: "Read requests throttled due to account limits", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "WriteAccountLimitThrottleEvents", Description: "Write requests throttled due to account limits", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "ReadKeyRangeThroughputThrottleEvents", Description: "Read requests throttled due to partition limits", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "WriteKeyRangeThroughputThrottleEvents", Description: "Write requests throttled due to partition limits", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "ReadMaxOnDemandThroughputThrottleEvents", Description: "Read requests throttled due to on-demand max throughput", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "WriteMaxOnDemandThroughputThrottleEvents", Description: "Write requests throttled due to on-demand max throughput", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "ReadProvisionedThroughputThrottleEvents", Description: "Read requests throttled due to provisioned limits", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "WriteProvisionedThroughputThrottleEvents", Description: "Write requests throttled due to provisioned limits", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		// Replication Metrics
		{Name: "AgeOfOldestUnreplicatedRecord", Description: "Elapsed time since unreplicated record first appeared in table", Unit: "Milliseconds", InstrumentType: "gauge", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "PendingReplicationCount", Description: "Item updates not yet written to replica tables", Unit: "Count", InstrumentType: "gauge", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "ReplicationLatency", Description: "Elapsed time between item appearing in stream and replica", Unit: "Milliseconds", InstrumentType: "gauge", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		// Stream Metrics
		{Name: "ConsumedChangeDataCaptureUnits", Description: "Number of consumed change data capture units", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "FailedToReplicateRecordCount", Description: "Records failed to replicate to Kinesis stream", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "ThrottledPutRecordCount", Description: "Records throttled by Kinesis stream", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		// GSI Metrics
		{Name: "OnlineIndexConsumedWriteCapacity", Description: "Write capacity consumed when adding new GSI", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "OnlineIndexPercentageProgress", Description: "Percentage completion of new GSI creation", Unit: "Percent", InstrumentType: "gauge", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "OnlineIndexThrottleEvents", Description: "Write throttle events during GSI creation", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		// Other Metrics
		{Name: "TimeToLiveDeletedItemCount", Description: "Items deleted by TTL", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		{Name: "TransactionConflict", Description: "Rejected item-level requests due to transaction conflicts", Unit: "Count", InstrumentType: "counter", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/DynamoDB"},
		// Usage Metrics
		{Name: "AccountProvisionedWriteCapacityUnits", Description: "Sum of write capacity units provisioned for all tables/GSIs", Unit: "Count", InstrumentType: "gauge", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/Usage"},
		{Name: "AccountProvisionedReadCapacityUnits", Description: "Sum of read capacity units provisioned for all tables/GSIs", Unit: "Count", InstrumentType: "gauge", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/Usage"},
		{Name: "TableCount", Description: "Number of active tables in account", Unit: "Count", InstrumentType: "gauge", ComponentName: "DynamoDB", ComponentType: "platform", SourceLocation: "AWS/Usage"},
	}, nil
}
