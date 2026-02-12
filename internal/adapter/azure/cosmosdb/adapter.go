package cosmosdb

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
	return "azure-cosmosdb"
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
	return "https://learn.microsoft.com/en-us/azure/azure-monitor/reference/supported-metrics/microsoft-documentdb-databaseaccounts-metrics"
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
	metrics = append(metrics, storageMetrics()...)
	metrics = append(metrics, availabilityMetrics()...)
	metrics = append(metrics, latencyMetrics()...)
	metrics = append(metrics, throughputMetrics()...)
	return metrics, nil
}

func requestMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.cosmosdb.total_requests", Description: "Number of requests made", Unit: "1", InstrumentType: "counter", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
		{Name: "azure.cosmosdb.total_request_units", Description: "Request Units consumed", Unit: "1", InstrumentType: "counter", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
		{Name: "azure.cosmosdb.provisioned_throughput", Description: "Provisioned throughput", Unit: "1", InstrumentType: "gauge", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
		{Name: "azure.cosmosdb.autoscale_max_throughput", Description: "Autoscale max throughput", Unit: "1", InstrumentType: "gauge", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
		{Name: "azure.cosmosdb.metadata_requests", Description: "Count of metadata requests", Unit: "1", InstrumentType: "counter", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
		{Name: "azure.cosmosdb.mongo_requests", Description: "Number of Mongo requests made", Unit: "1", InstrumentType: "counter", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
		{Name: "azure.cosmosdb.mongo_request_charge", Description: "Mongo request units consumed", Unit: "1", InstrumentType: "counter", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
		{Name: "azure.cosmosdb.cassandra_requests", Description: "Number of Cassandra requests made", Unit: "1", InstrumentType: "counter", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
		{Name: "azure.cosmosdb.cassandra_request_charges", Description: "Request Units consumed for Cassandra requests made", Unit: "1", InstrumentType: "counter", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
		{Name: "azure.cosmosdb.gremlin_requests", Description: "Number of Gremlin requests made", Unit: "1", InstrumentType: "counter", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
		{Name: "azure.cosmosdb.gremlin_request_charges", Description: "Request Units consumed for Gremlin requests made", Unit: "1", InstrumentType: "counter", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
	}
}

func storageMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.cosmosdb.data_usage", Description: "Total data usage reported at 5 minutes granularity", Unit: "By", InstrumentType: "gauge", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
		{Name: "azure.cosmosdb.index_usage", Description: "Total index usage reported at 5 minutes granularity", Unit: "By", InstrumentType: "gauge", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
		{Name: "azure.cosmosdb.document_quota", Description: "Total storage quota reported at 5 minutes granularity", Unit: "By", InstrumentType: "gauge", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
		{Name: "azure.cosmosdb.document_count", Description: "Total document count reported at 5 minutes granularity", Unit: "1", InstrumentType: "gauge", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
	}
}

func availabilityMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.cosmosdb.service_availability", Description: "Account requests availability at one hour, day or month granularity", Unit: "%", InstrumentType: "gauge", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
		{Name: "azure.cosmosdb.http_2xx", Description: "Count of requests resulting in HTTP 2xx status codes", Unit: "1", InstrumentType: "counter", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
		{Name: "azure.cosmosdb.http_3xx", Description: "Count of requests resulting in HTTP 3xx status codes", Unit: "1", InstrumentType: "counter", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
		{Name: "azure.cosmosdb.http_4xx", Description: "Count of requests resulting in HTTP 4xx status codes", Unit: "1", InstrumentType: "counter", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
		{Name: "azure.cosmosdb.http_5xx", Description: "Count of requests resulting in HTTP 5xx status codes", Unit: "1", InstrumentType: "counter", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
	}
}

func latencyMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.cosmosdb.server_side_latency", Description: "Server side latency for the account", Unit: "ms", InstrumentType: "gauge", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
		{Name: "azure.cosmosdb.replication_latency", Description: "P99 replication latency across source and target regions for geo-enabled account", Unit: "ms", InstrumentType: "gauge", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
		{Name: "azure.cosmosdb.cassandra_connector_average_replicationlatency", Description: "Average replication latency for Cassandra Connector", Unit: "ms", InstrumentType: "gauge", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
	}
}

func throughputMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.cosmosdb.normalized_ru_consumption", Description: "Max RU consumption percentage per minute", Unit: "%", InstrumentType: "gauge", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
		{Name: "azure.cosmosdb.physical_partition_throughput_info", Description: "Provisioned throughput in RU/s for each physical partition", Unit: "1", InstrumentType: "gauge", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
		{Name: "azure.cosmosdb.physical_partition_size_info", Description: "Data size in KB for each physical partition", Unit: "KBy", InstrumentType: "gauge", ComponentName: "Cosmos DB", ComponentType: "platform", SourceLocation: "Microsoft.DocumentDB/databaseAccounts"},
	}
}
