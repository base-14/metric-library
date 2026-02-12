package aks

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
	return "azure-aks"
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
	return "https://learn.microsoft.com/en-us/azure/azure-monitor/reference/supported-metrics/microsoft-containerservice-managedclusters-metrics"
}

func (a *Adapter) Fetch(_ context.Context, _ adapter.FetchOptions) (*adapter.FetchResult, error) {
	return &adapter.FetchResult{
		Commit:    time.Now().Format("2006-01-02"),
		Timestamp: time.Now(),
	}, nil
}

func (a *Adapter) Extract(_ context.Context, _ *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	var metrics []*adapter.RawMetric
	metrics = append(metrics, nodeMetrics()...)
	metrics = append(metrics, podMetrics()...)
	metrics = append(metrics, clusterMetrics()...)
	metrics = append(metrics, apiServerMetrics()...)
	return metrics, nil
}

func nodeMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.aks.node_cpu_usage_millicores", Description: "Aggregated measurement of CPU utilization in millicores across the cluster", Unit: "m", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.node_cpu_usage_percentage", Description: "Aggregated average CPU utilization measured in percentage across the cluster", Unit: "%", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.node_memory_rss_bytes", Description: "Container RSS memory used in bytes", Unit: "By", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.node_memory_rss_percentage", Description: "Container RSS memory used in percent", Unit: "%", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.node_memory_working_set_bytes", Description: "Container working set memory used in bytes", Unit: "By", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.node_memory_working_set_percentage", Description: "Container working set memory used in percent", Unit: "%", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.node_disk_usage_bytes", Description: "Disk used in bytes by device", Unit: "By", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.node_disk_usage_percentage", Description: "Disk used in percent by device", Unit: "%", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.node_network_in_bytes", Description: "Network received bytes", Unit: "By", InstrumentType: "counter", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.node_network_out_bytes", Description: "Network transmitted bytes", Unit: "By", InstrumentType: "counter", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
	}
}

func podMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.aks.kube_pod_status_ready", Description: "Number of pods in Ready state", Unit: "1", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.kube_pod_status_phase", Description: "Number of pods by phase", Unit: "1", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.kube_pod_containers_ready", Description: "Number of containers in Ready state per pod", Unit: "1", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.kube_pod_containers_restarts", Description: "Number of container restarts per pod", Unit: "1", InstrumentType: "counter", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.kube_pod_containers_last_state_terminated", Description: "Number of containers in last terminated state per pod", Unit: "1", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.kube_deployment_status_replicas_ready", Description: "Number of ready replicas per deployment", Unit: "1", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.kube_deployment_spec_replicas", Description: "Number of desired replicas per deployment", Unit: "1", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.kube_deployment_status_replicas_available", Description: "Number of available replicas per deployment", Unit: "1", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.kube_deployment_status_replicas_unavailable", Description: "Number of unavailable replicas per deployment", Unit: "1", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
	}
}

func clusterMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.aks.kube_node_status_condition", Description: "Statuses for various node conditions", Unit: "1", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.kube_node_status_allocatable_cpu_cores", Description: "Total number of available CPU cores in a managed cluster", Unit: "1", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.kube_node_status_allocatable_memory_bytes", Description: "Total amount of available memory in a managed cluster", Unit: "By", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.cluster_autoscaler_cluster_safe_to_autoscale", Description: "Determines whether the cluster autoscaler will take action on the cluster", Unit: "1", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.cluster_autoscaler_scale_down_in_cooldown", Description: "Determines if the scale down is in cooldown", Unit: "1", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.cluster_autoscaler_unneeded_nodes_count", Description: "Cluster autoscaler marks those nodes as candidates for deletion and are eventually deleted", Unit: "1", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.cluster_autoscaler_unschedulable_pods_count", Description: "Number of pods that are currently unschedulable in the cluster", Unit: "1", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
	}
}

func apiServerMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.aks.apiserver_current_inflight_requests", Description: "Maximum number of currently used inflight request limit on API Server per request kind", Unit: "1", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.node_count", Description: "Number of nodes in the managed cluster", Unit: "1", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.total_number_of_cpu_cores_in_managed_cluster", Description: "Total number of CPU cores in a managed cluster", Unit: "1", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
		{Name: "azure.aks.total_amount_of_memory_in_managed_cluster", Description: "Total amount of memory in a managed cluster", Unit: "By", InstrumentType: "gauge", ComponentName: "AKS", ComponentType: "platform", SourceLocation: "Microsoft.ContainerService/managedClusters"},
	}
}
