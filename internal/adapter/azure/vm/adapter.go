package vm

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
	return "azure-vm"
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
	return "https://learn.microsoft.com/en-us/azure/azure-monitor/reference/supported-metrics/microsoft-compute-virtualmachines-metrics"
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
	metrics = append(metrics, memoryMetrics()...)
	metrics = append(metrics, diskMetrics()...)
	metrics = append(metrics, networkMetrics()...)
	metrics = append(metrics, dataDiskMetrics()...)
	metrics = append(metrics, osDiskMetrics()...)
	return metrics, nil
}

func cpuMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.vm.percentage_cpu", Description: "The percentage of allocated compute units that are currently in use by the Virtual Machine(s)", Unit: "%", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.cpu_credits_remaining", Description: "Total number of credits available to burst", Unit: "1", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.cpu_credits_consumed", Description: "Total number of credits consumed by the Virtual Machine", Unit: "1", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
	}
}

func memoryMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.vm.available_memory_bytes", Description: "Amount of physical memory in bytes immediately available for allocation to a process or for system use", Unit: "By", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
	}
}

func diskMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.vm.disk_read_bytes", Description: "Bytes read from disk during monitoring period", Unit: "By", InstrumentType: "counter", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.disk_write_bytes", Description: "Bytes written to disk during monitoring period", Unit: "By", InstrumentType: "counter", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.disk_read_operations_per_sec", Description: "Disk Read IOPS", Unit: "{operations}/s", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.disk_write_operations_per_sec", Description: "Disk Write IOPS", Unit: "{operations}/s", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
	}
}

func networkMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.vm.network_in_total", Description: "The number of bytes received on all network interfaces by the Virtual Machine(s)", Unit: "By", InstrumentType: "counter", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.network_out_total", Description: "The number of bytes out on all network interfaces by the Virtual Machine(s)", Unit: "By", InstrumentType: "counter", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.network_in_billable", Description: "The number of billable bytes received on all network interfaces by the Virtual Machine(s)", Unit: "By", InstrumentType: "counter", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.network_out_billable", Description: "The number of billable bytes out on all network interfaces by the Virtual Machine(s)", Unit: "By", InstrumentType: "counter", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
	}
}

func dataDiskMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.vm.data_disk_read_bytes_per_sec", Description: "Bytes per second read from a single disk during monitoring period", Unit: "By/s", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.data_disk_write_bytes_per_sec", Description: "Bytes per second written to a single disk during monitoring period", Unit: "By/s", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.data_disk_read_operations_per_sec", Description: "Read IOPS from a single disk during monitoring period", Unit: "{operations}/s", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.data_disk_write_operations_per_sec", Description: "Write IOPS from a single disk during monitoring period", Unit: "{operations}/s", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.data_disk_queue_depth", Description: "Data Disk Queue Depth (or Queue Length)", Unit: "1", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.data_disk_bandwidth_consumed_percentage", Description: "Percentage of data disk bandwidth consumed per minute", Unit: "%", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.data_disk_iops_consumed_percentage", Description: "Percentage of data disk I/Os consumed per minute", Unit: "%", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.data_disk_target_bandwidth", Description: "Baseline bytes per second throughput data disk can achieve without bursting", Unit: "By/s", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.data_disk_target_iops", Description: "Baseline IOPS data disk can achieve without bursting", Unit: "{operations}/s", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.data_disk_max_burst_bandwidth", Description: "Maximum bytes per second throughput data disk can achieve with bursting", Unit: "By/s", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.data_disk_max_burst_iops", Description: "Maximum IOPS data disk can achieve with bursting", Unit: "{operations}/s", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.data_disk_used_burst_io_credits_percentage", Description: "Percentage of data disk burst I/O credits used so far", Unit: "%", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.data_disk_used_burst_bps_credits_percentage", Description: "Percentage of data disk burst bandwidth credits used so far", Unit: "%", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
	}
}

func osDiskMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.vm.os_disk_read_bytes_per_sec", Description: "Bytes per second read from a single disk during monitoring period for OS disk", Unit: "By/s", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.os_disk_write_bytes_per_sec", Description: "Bytes per second written to a single disk during monitoring period for OS disk", Unit: "By/s", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.os_disk_read_operations_per_sec", Description: "Read IOPS from a single disk during monitoring period for OS disk", Unit: "{operations}/s", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.os_disk_write_operations_per_sec", Description: "Write IOPS from a single disk during monitoring period for OS disk", Unit: "{operations}/s", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.os_disk_queue_depth", Description: "OS Disk Queue Depth (or Queue Length)", Unit: "1", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.os_disk_bandwidth_consumed_percentage", Description: "Percentage of operating system disk bandwidth consumed per minute", Unit: "%", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.os_disk_iops_consumed_percentage", Description: "Percentage of operating system disk I/Os consumed per minute", Unit: "%", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.os_disk_target_bandwidth", Description: "Baseline bytes per second throughput OS disk can achieve without bursting", Unit: "By/s", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.os_disk_target_iops", Description: "Baseline IOPS OS disk can achieve without bursting", Unit: "{operations}/s", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.os_disk_max_burst_bandwidth", Description: "Maximum bytes per second throughput OS disk can achieve with bursting", Unit: "By/s", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
		{Name: "azure.vm.os_disk_max_burst_iops", Description: "Maximum IOPS OS disk can achieve with bursting", Unit: "{operations}/s", InstrumentType: "gauge", ComponentName: "Virtual Machines", ComponentType: "platform", SourceLocation: "Microsoft.Compute/virtualMachines"},
	}
}
