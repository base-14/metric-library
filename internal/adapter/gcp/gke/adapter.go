package gke

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
	return "gcp-gke"
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
	return "https://cloud.google.com/monitoring/api/metrics_kubernetes"
}

func (a *Adapter) Fetch(_ context.Context, _ adapter.FetchOptions) (*adapter.FetchResult, error) {
	return &adapter.FetchResult{
		Commit:    time.Now().Format("2006-01-02"),
		Timestamp: time.Now(),
	}, nil
}

func (a *Adapter) Extract(_ context.Context, _ *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	var metrics []*adapter.RawMetric
	metrics = append(metrics, containerMetrics()...)
	metrics = append(metrics, nodeMetrics()...)
	metrics = append(metrics, podMetrics()...)
	metrics = append(metrics, clusterMetrics()...)
	return metrics, nil
}

func containerMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "kubernetes.io/container/cpu/core_usage_time", Description: "Cumulative CPU usage on all cores in seconds", Unit: "s", InstrumentType: "counter", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/container/cpu/limit_cores", Description: "CPU cores limit of the container", Unit: "1", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/container/cpu/limit_utilization", Description: "The fraction of the CPU limit that is currently in use on the instance", Unit: "1", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/container/cpu/request_cores", Description: "Number of CPU cores requested by the container", Unit: "1", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/container/cpu/request_utilization", Description: "The fraction of the requested CPU that is currently in use on the instance", Unit: "1", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/container/memory/limit_bytes", Description: "Memory limit of the container in bytes", Unit: "By", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/container/memory/limit_utilization", Description: "The fraction of the memory limit that is currently in use on the instance", Unit: "1", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/container/memory/page_fault_count", Description: "Number of page faults, broken down by type: major and minor", Unit: "1", InstrumentType: "counter", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/container/memory/request_bytes", Description: "Memory request of the container in bytes", Unit: "By", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/container/memory/request_utilization", Description: "The fraction of the requested memory that is currently in use on the instance", Unit: "1", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/container/memory/used_bytes", Description: "Memory usage in bytes", Unit: "By", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/container/ephemeral_storage/limit_bytes", Description: "Local ephemeral storage limit in bytes", Unit: "By", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/container/ephemeral_storage/request_bytes", Description: "Local ephemeral storage request in bytes", Unit: "By", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/container/ephemeral_storage/used_bytes", Description: "Local ephemeral storage usage in bytes", Unit: "By", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/container/restart_count", Description: "Number of times the container has restarted", Unit: "1", InstrumentType: "counter", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/container/uptime", Description: "Time in seconds that the container has been running", Unit: "s", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
	}
}

func nodeMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "kubernetes.io/node/cpu/core_usage_time", Description: "Cumulative CPU usage on all cores in seconds", Unit: "s", InstrumentType: "counter", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/node/cpu/total_cores", Description: "Total number of CPU cores on the node", Unit: "1", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/node/cpu/allocatable_cores", Description: "Number of allocatable CPU cores on the node", Unit: "1", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/node/cpu/allocatable_utilization", Description: "The fraction of the allocatable CPU that is currently in use on the instance", Unit: "1", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/node/memory/total_bytes", Description: "Number of bytes of memory allocatable on the node", Unit: "By", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/node/memory/allocatable_bytes", Description: "Cumulative number of bytes of memory available for scheduling", Unit: "By", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/node/memory/allocatable_utilization", Description: "The fraction of the allocatable memory that is currently in use on the instance", Unit: "1", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/node/memory/used_bytes", Description: "Cumulative number of bytes of memory used on the node", Unit: "By", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/node/memory/page_fault_count", Description: "Cumulative number of page faults on the node", Unit: "1", InstrumentType: "counter", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/node/network/received_bytes_count", Description: "Cumulative number of bytes received by the node over the network", Unit: "By", InstrumentType: "counter", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/node/network/sent_bytes_count", Description: "Cumulative number of bytes transmitted by the node over the network", Unit: "By", InstrumentType: "counter", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/node/ephemeral_storage/total_bytes", Description: "Total ephemeral storage bytes on the node", Unit: "By", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/node/ephemeral_storage/allocatable_bytes", Description: "Local ephemeral storage bytes available for allocation on the node", Unit: "By", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/node/ephemeral_storage/used_bytes", Description: "Local ephemeral storage bytes used by the node", Unit: "By", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/node/pid_limit", Description: "The max PID of the node", Unit: "1", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/node/pid_used", Description: "The number of running process IDs on the node", Unit: "1", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
	}
}

func podMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "kubernetes.io/pod/network/received_bytes_count", Description: "Cumulative number of bytes received by the pod over the network", Unit: "By", InstrumentType: "counter", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/pod/network/sent_bytes_count", Description: "Cumulative number of bytes transmitted by the pod over the network", Unit: "By", InstrumentType: "counter", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/pod/volume/total_bytes", Description: "Total number of disk bytes available to the pod", Unit: "By", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/pod/volume/used_bytes", Description: "Number of disk bytes used by the pod", Unit: "By", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/pod/volume/utilization", Description: "The fraction of the volume that is currently being used by the instance", Unit: "1", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
	}
}

func clusterMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "kubernetes.io/node_daemon/cpu/core_usage_time", Description: "Cumulative CPU usage of the node-level system daemon in seconds", Unit: "s", InstrumentType: "counter", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/node_daemon/memory/used_bytes", Description: "Memory usage by the system daemon in bytes", Unit: "By", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
		{Name: "kubernetes.io/autoscaler/container/cpu/per_replica_recommended_request_cores", Description: "Recommended CPU request per replica for the container from Vertical Pod Autoscaler", Unit: "1", InstrumentType: "gauge", ComponentName: "GKE", ComponentType: "platform", SourceLocation: "kubernetes.io"},
	}
}
