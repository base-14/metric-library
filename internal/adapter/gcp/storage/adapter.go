package storage

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
	return "gcp-storage"
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
	return "https://cloud.google.com/monitoring/api/metrics_gcp#gcp-storage"
}

func (a *Adapter) Fetch(_ context.Context, _ adapter.FetchOptions) (*adapter.FetchResult, error) {
	return &adapter.FetchResult{
		Commit:    time.Now().Format("2006-01-02"),
		Timestamp: time.Now(),
	}, nil
}

func (a *Adapter) Extract(_ context.Context, _ *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	var metrics []*adapter.RawMetric
	metrics = append(metrics, apiMetrics()...)
	metrics = append(metrics, bucketMetrics()...)
	return metrics, nil
}

func apiMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "storage.googleapis.com/api/request_count", Description: "Delta count of API calls grouped by the API method name and response code", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Storage", ComponentType: "platform", SourceLocation: "storage.googleapis.com"},
		{Name: "storage.googleapis.com/network/received_bytes_count", Description: "Delta count of bytes received over the network grouped by the API method name and response code", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud Storage", ComponentType: "platform", SourceLocation: "storage.googleapis.com"},
		{Name: "storage.googleapis.com/network/sent_bytes_count", Description: "Delta count of bytes sent over the network grouped by the API method name and response code", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud Storage", ComponentType: "platform", SourceLocation: "storage.googleapis.com"},
	}
}

func bucketMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "storage.googleapis.com/storage/total_bytes", Description: "Total size of all objects in the bucket in bytes", Unit: "By", InstrumentType: "gauge", ComponentName: "Cloud Storage", ComponentType: "platform", SourceLocation: "storage.googleapis.com"},
		{Name: "storage.googleapis.com/storage/object_count", Description: "Total number of objects per bucket", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud Storage", ComponentType: "platform", SourceLocation: "storage.googleapis.com"},
		{Name: "storage.googleapis.com/storage/total_byte_seconds", Description: "Delta count of bytes received over the network, grouped by the API method name and response code (used for billing)", Unit: "By.s", InstrumentType: "counter", ComponentName: "Cloud Storage", ComponentType: "platform", SourceLocation: "storage.googleapis.com"},
		{Name: "storage.googleapis.com/authn/authentication_count", Description: "Delta count of authentication requests grouped by result and authentication method", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Storage", ComponentType: "platform", SourceLocation: "storage.googleapis.com"},
		{Name: "storage.googleapis.com/authz/acl_based_object_access_count", Description: "Delta count of requests that result in an object being granted access solely due to object ACLs", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Storage", ComponentType: "platform", SourceLocation: "storage.googleapis.com"},
		{Name: "storage.googleapis.com/authz/acl_operations_count", Description: "Usage of ACL operations broken down by type", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Storage", ComponentType: "platform", SourceLocation: "storage.googleapis.com"},
		{Name: "storage.googleapis.com/authz/object_specific_acl_mutation_count", Description: "Delta count of changes made to object specific ACLs", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Storage", ComponentType: "platform", SourceLocation: "storage.googleapis.com"},
		{Name: "storage.googleapis.com/replication/meeting_rpo", Description: "Whether the most recent write to a dual-region or multi-region bucket was replicated to meet the turbo replication RPO", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud Storage", ComponentType: "platform", SourceLocation: "storage.googleapis.com"},
		{Name: "storage.googleapis.com/replication/objects_pending_replication_count", Description: "Number of objects pending replication", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud Storage", ComponentType: "platform", SourceLocation: "storage.googleapis.com"},
	}
}
