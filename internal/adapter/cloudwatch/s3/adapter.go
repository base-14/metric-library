package s3

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
	return "cloudwatch-s3"
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
	return "https://docs.aws.amazon.com/AmazonS3/latest/userguide/metrics-dimensions.html"
}

func (a *Adapter) Fetch(_ context.Context, _ adapter.FetchOptions) (*adapter.FetchResult, error) {
	return &adapter.FetchResult{
		Commit:    time.Now().Format("2006-01-02"),
		Timestamp: time.Now(),
	}, nil
}

func (a *Adapter) Extract(_ context.Context, _ *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	return []*adapter.RawMetric{
		// Daily Storage Metrics
		{Name: "BucketSizeBytes", Description: "Amount of data stored in a bucket across various storage classes", Unit: "Bytes", InstrumentType: "gauge", ComponentName: "S3", ComponentType: "platform", SourceLocation: "AWS/S3"},
		{Name: "NumberOfObjects", Description: "Total number of objects stored in a bucket", Unit: "Count", InstrumentType: "gauge", ComponentName: "S3", ComponentType: "platform", SourceLocation: "AWS/S3"},
		// Request Metrics
		{Name: "AllRequests", Description: "Total number of HTTP requests made to a bucket", Unit: "Count", InstrumentType: "counter", ComponentName: "S3", ComponentType: "platform", SourceLocation: "AWS/S3"},
		{Name: "GetRequests", Description: "Number of HTTP GET requests", Unit: "Count", InstrumentType: "counter", ComponentName: "S3", ComponentType: "platform", SourceLocation: "AWS/S3"},
		{Name: "PutRequests", Description: "Number of HTTP PUT requests", Unit: "Count", InstrumentType: "counter", ComponentName: "S3", ComponentType: "platform", SourceLocation: "AWS/S3"},
		{Name: "DeleteRequests", Description: "Number of HTTP DELETE requests", Unit: "Count", InstrumentType: "counter", ComponentName: "S3", ComponentType: "platform", SourceLocation: "AWS/S3"},
		{Name: "HeadRequests", Description: "Number of HTTP HEAD requests", Unit: "Count", InstrumentType: "counter", ComponentName: "S3", ComponentType: "platform", SourceLocation: "AWS/S3"},
		{Name: "PostRequests", Description: "Number of HTTP POST requests", Unit: "Count", InstrumentType: "counter", ComponentName: "S3", ComponentType: "platform", SourceLocation: "AWS/S3"},
		{Name: "ListRequests", Description: "Number of HTTP requests that list bucket contents", Unit: "Count", InstrumentType: "counter", ComponentName: "S3", ComponentType: "platform", SourceLocation: "AWS/S3"},
		{Name: "SelectRequests", Description: "Number of S3 SelectObjectContent requests", Unit: "Count", InstrumentType: "counter", ComponentName: "S3", ComponentType: "platform", SourceLocation: "AWS/S3"},
		{Name: "SelectBytesScanned", Description: "Number of bytes scanned by SelectObjectContent requests", Unit: "Bytes", InstrumentType: "counter", ComponentName: "S3", ComponentType: "platform", SourceLocation: "AWS/S3"},
		{Name: "SelectBytesReturned", Description: "Number of bytes returned by SelectObjectContent requests", Unit: "Bytes", InstrumentType: "counter", ComponentName: "S3", ComponentType: "platform", SourceLocation: "AWS/S3"},
		{Name: "BytesDownloaded", Description: "Number of bytes downloaded from the bucket", Unit: "Bytes", InstrumentType: "counter", ComponentName: "S3", ComponentType: "platform", SourceLocation: "AWS/S3"},
		{Name: "BytesUploaded", Description: "Number of bytes uploaded to the bucket", Unit: "Bytes", InstrumentType: "counter", ComponentName: "S3", ComponentType: "platform", SourceLocation: "AWS/S3"},
		{Name: "4xxErrors", Description: "Number of HTTP 4xx client error requests", Unit: "Count", InstrumentType: "counter", ComponentName: "S3", ComponentType: "platform", SourceLocation: "AWS/S3"},
		{Name: "5xxErrors", Description: "Number of HTTP 5xx server error requests", Unit: "Count", InstrumentType: "counter", ComponentName: "S3", ComponentType: "platform", SourceLocation: "AWS/S3"},
		{Name: "FirstByteLatency", Description: "Time from complete request received to response starts", Unit: "Milliseconds", InstrumentType: "gauge", ComponentName: "S3", ComponentType: "platform", SourceLocation: "AWS/S3"},
		{Name: "TotalRequestLatency", Description: "Elapsed time from first byte received to last byte sent", Unit: "Milliseconds", InstrumentType: "gauge", ComponentName: "S3", ComponentType: "platform", SourceLocation: "AWS/S3"},
		// Replication Metrics
		{Name: "ReplicationLatency", Description: "Maximum seconds the destination region lags behind source", Unit: "Seconds", InstrumentType: "gauge", ComponentName: "S3", ComponentType: "platform", SourceLocation: "AWS/S3"},
		{Name: "BytesPendingReplication", Description: "Total bytes of objects pending replication", Unit: "Bytes", InstrumentType: "gauge", ComponentName: "S3", ComponentType: "platform", SourceLocation: "AWS/S3"},
		{Name: "OperationsPendingReplication", Description: "Number of operations pending replication", Unit: "Count", InstrumentType: "gauge", ComponentName: "S3", ComponentType: "platform", SourceLocation: "AWS/S3"},
		{Name: "OperationsFailedReplication", Description: "Number of operations that failed to replicate", Unit: "Count", InstrumentType: "counter", ComponentName: "S3", ComponentType: "platform", SourceLocation: "AWS/S3"},
	}, nil
}
