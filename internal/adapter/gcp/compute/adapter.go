package compute

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
	return "gcp-compute"
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
	return "https://cloud.google.com/monitoring/api/metrics_gcp#gcp-compute"
}

func (a *Adapter) Fetch(_ context.Context, _ adapter.FetchOptions) (*adapter.FetchResult, error) {
	return &adapter.FetchResult{
		Commit:    time.Now().Format("2006-01-02"),
		Timestamp: time.Now(),
	}, nil
}

func (a *Adapter) Extract(_ context.Context, _ *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	var metrics []*adapter.RawMetric
	metrics = append(metrics, cpuMetrics()...)
	metrics = append(metrics, diskMetrics()...)
	metrics = append(metrics, networkMetrics()...)
	metrics = append(metrics, instanceMetrics()...)
	metrics = append(metrics, firewallMetrics()...)
	return metrics, nil
}

func cpuMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "compute.googleapis.com/instance/cpu/utilization", Description: "Fractional utilization of allocated CPU on the instance", Unit: "1", InstrumentType: "gauge", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/cpu/usage_time", Description: "CPU usage in seconds", Unit: "s", InstrumentType: "counter", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/cpu/reserved_cores", Description: "Number of vCPUs reserved on the host of the instance", Unit: "1", InstrumentType: "gauge", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/cpu/scheduler_wait_time", Description: "Wait time is the time a vCPU is ready to run, but unexpectedly not scheduled to run", Unit: "s", InstrumentType: "counter", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/cpu/guest_visible_vcpus", Description: "Number of vCPUs visible inside the guest", Unit: "1", InstrumentType: "gauge", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
	}
}

func diskMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "compute.googleapis.com/instance/disk/read_bytes_count", Description: "Count of bytes read from disk", Unit: "By", InstrumentType: "counter", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/disk/read_ops_count", Description: "Count of disk read IO operations", Unit: "1", InstrumentType: "counter", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/disk/write_bytes_count", Description: "Count of bytes written to disk", Unit: "By", InstrumentType: "counter", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/disk/write_ops_count", Description: "Count of disk write IO operations", Unit: "1", InstrumentType: "counter", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/disk/throttled_read_bytes_count", Description: "Count of bytes in throttled read operations", Unit: "By", InstrumentType: "counter", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/disk/throttled_read_ops_count", Description: "Count of throttled read operations", Unit: "1", InstrumentType: "counter", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/disk/throttled_write_bytes_count", Description: "Count of bytes in throttled write operations", Unit: "By", InstrumentType: "counter", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/disk/throttled_write_ops_count", Description: "Count of throttled write operations", Unit: "1", InstrumentType: "counter", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
	}
}

func networkMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "compute.googleapis.com/instance/network/received_bytes_count", Description: "Count of bytes received from the network", Unit: "By", InstrumentType: "counter", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/network/received_packets_count", Description: "Count of packets received from the network", Unit: "1", InstrumentType: "counter", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/network/sent_bytes_count", Description: "Count of bytes sent over the network", Unit: "By", InstrumentType: "counter", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/network/sent_packets_count", Description: "Count of packets sent over the network", Unit: "1", InstrumentType: "counter", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/network/received_packets_dropped_count", Description: "Count of incoming packets dropped by the network", Unit: "1", InstrumentType: "counter", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/network/sent_packets_dropped_count", Description: "Count of outgoing packets dropped by the network", Unit: "1", InstrumentType: "counter", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
	}
}

func instanceMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "compute.googleapis.com/instance/uptime", Description: "How long the VM has been running in seconds", Unit: "s", InstrumentType: "counter", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/uptime_total", Description: "Elapsed time since the VM was started in seconds", Unit: "s", InstrumentType: "counter", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/memory/balloon/ram_used", Description: "Memory used by the VM as seen by the hypervisor", Unit: "By", InstrumentType: "gauge", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/memory/balloon/ram_size", Description: "Total memory of the VM as seen by the hypervisor", Unit: "By", InstrumentType: "gauge", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/memory/balloon/swap_in_bytes_count", Description: "Amount of memory read into the guest from its own swap space", Unit: "By", InstrumentType: "counter", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/memory/balloon/swap_out_bytes_count", Description: "Amount of memory written from the guest to its own swap space", Unit: "By", InstrumentType: "counter", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/integrity/early_boot_validation_status", Description: "Validation status of early boot integrity policy", Unit: "1", InstrumentType: "gauge", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/instance/integrity/late_boot_validation_status", Description: "Validation status of late boot integrity policy", Unit: "1", InstrumentType: "gauge", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
	}
}

func firewallMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "compute.googleapis.com/firewall/dropped_bytes_count", Description: "Count of incoming bytes dropped by the firewall", Unit: "By", InstrumentType: "counter", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
		{Name: "compute.googleapis.com/firewall/dropped_packets_count", Description: "Count of incoming packets dropped by the firewall", Unit: "1", InstrumentType: "counter", ComponentName: "Compute Engine", ComponentType: "platform", SourceLocation: "compute.googleapis.com"},
	}
}
