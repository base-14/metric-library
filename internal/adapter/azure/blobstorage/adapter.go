package blobstorage

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
	return "azure-blobstorage"
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
	return "https://learn.microsoft.com/en-us/azure/azure-monitor/reference/supported-metrics/microsoft-storage-storageaccounts-blobservices-metrics"
}

func (a *Adapter) Fetch(_ context.Context, _ adapter.FetchOptions) (*adapter.FetchResult, error) {
	return &adapter.FetchResult{
		Commit:    time.Now().Format("2006-01-02"),
		Timestamp: time.Now(),
	}, nil
}

func (a *Adapter) Extract(_ context.Context, _ *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	var metrics []*adapter.RawMetric
	metrics = append(metrics, transactionMetrics()...)
	metrics = append(metrics, capacityMetrics()...)
	return metrics, nil
}

func transactionMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.blobstorage.availability", Description: "The percentage of availability for the storage service", Unit: "%", InstrumentType: "gauge", ComponentName: "Blob Storage", ComponentType: "platform", SourceLocation: "Microsoft.Storage/storageAccounts/blobServices"},
		{Name: "azure.blobstorage.egress", Description: "The amount of egress data in bytes", Unit: "By", InstrumentType: "counter", ComponentName: "Blob Storage", ComponentType: "platform", SourceLocation: "Microsoft.Storage/storageAccounts/blobServices"},
		{Name: "azure.blobstorage.ingress", Description: "The amount of ingress data in bytes", Unit: "By", InstrumentType: "counter", ComponentName: "Blob Storage", ComponentType: "platform", SourceLocation: "Microsoft.Storage/storageAccounts/blobServices"},
		{Name: "azure.blobstorage.success_server_latency", Description: "The average latency used by Azure Storage to process a successful request in milliseconds", Unit: "ms", InstrumentType: "gauge", ComponentName: "Blob Storage", ComponentType: "platform", SourceLocation: "Microsoft.Storage/storageAccounts/blobServices"},
		{Name: "azure.blobstorage.success_e2e_latency", Description: "The average end-to-end latency of successful requests made to a storage service in milliseconds", Unit: "ms", InstrumentType: "gauge", ComponentName: "Blob Storage", ComponentType: "platform", SourceLocation: "Microsoft.Storage/storageAccounts/blobServices"},
		{Name: "azure.blobstorage.transactions", Description: "The number of requests made to a storage service", Unit: "1", InstrumentType: "counter", ComponentName: "Blob Storage", ComponentType: "platform", SourceLocation: "Microsoft.Storage/storageAccounts/blobServices"},
	}
}

func capacityMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.blobstorage.blob_capacity", Description: "The amount of storage used by the storage account Blob service in bytes", Unit: "By", InstrumentType: "gauge", ComponentName: "Blob Storage", ComponentType: "platform", SourceLocation: "Microsoft.Storage/storageAccounts/blobServices"},
		{Name: "azure.blobstorage.blob_count", Description: "The number of blob objects stored in the storage account", Unit: "1", InstrumentType: "gauge", ComponentName: "Blob Storage", ComponentType: "platform", SourceLocation: "Microsoft.Storage/storageAccounts/blobServices"},
		{Name: "azure.blobstorage.container_count", Description: "The number of containers in the storage account", Unit: "1", InstrumentType: "gauge", ComponentName: "Blob Storage", ComponentType: "platform", SourceLocation: "Microsoft.Storage/storageAccounts/blobServices"},
		{Name: "azure.blobstorage.index_capacity", Description: "The amount of storage used by Azure Data Lake Storage Gen2 hierarchical index", Unit: "By", InstrumentType: "gauge", ComponentName: "Blob Storage", ComponentType: "platform", SourceLocation: "Microsoft.Storage/storageAccounts/blobServices"},
		{Name: "azure.blobstorage.blob_provision_size", Description: "The amount of storage provisioned in the storage account Blob service in bytes", Unit: "By", InstrumentType: "gauge", ComponentName: "Blob Storage", ComponentType: "platform", SourceLocation: "Microsoft.Storage/storageAccounts/blobServices"},
	}
}
