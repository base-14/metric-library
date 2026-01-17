package ec2

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
	return "cloudwatch-ec2"
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
	return "https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/viewing_metrics_with_cloudwatch.html"
}

func (a *Adapter) Fetch(_ context.Context, _ adapter.FetchOptions) (*adapter.FetchResult, error) {
	return &adapter.FetchResult{
		Commit:    time.Now().Format("2006-01-02"),
		Timestamp: time.Now(),
	}, nil
}

func (a *Adapter) Extract(_ context.Context, _ *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	metrics := []*adapter.RawMetric{}

	// Instance metrics
	metrics = append(metrics, ec2Metrics()...)

	// CPU Credit metrics (burstable instances)
	metrics = append(metrics, cpuCreditMetrics()...)

	// EBS metrics for Nitro instances
	metrics = append(metrics, ebsMetrics()...)

	// Status check metrics
	metrics = append(metrics, statusCheckMetrics()...)

	return metrics, nil
}

func ec2Metrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{
			Name:           "CPUUtilization",
			Description:    "The percentage of physical CPU time that Amazon EC2 uses to run the EC2 instance, including time spent to run both the user code and the Amazon EC2 code",
			Unit:           "Percent",
			InstrumentType: "gauge",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "DiskReadOps",
			Description:    "Completed read operations from all instance store volumes available to the instance",
			Unit:           "Count",
			InstrumentType: "counter",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "DiskWriteOps",
			Description:    "Completed write operations to all instance store volumes available to the instance",
			Unit:           "Count",
			InstrumentType: "counter",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "DiskReadBytes",
			Description:    "Bytes read from all instance store volumes available to the instance",
			Unit:           "Bytes",
			InstrumentType: "counter",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "DiskWriteBytes",
			Description:    "Bytes written to all instance store volumes available to the instance",
			Unit:           "Bytes",
			InstrumentType: "counter",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "MetadataNoToken",
			Description:    "The number of times the instance metadata service was successfully accessed using a method that does not use a token (IMDSv1)",
			Unit:           "Count",
			InstrumentType: "counter",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "MetadataNoTokenRejected",
			Description:    "The number of times an IMDSv1 call was attempted and rejected after IMDSv1 was disabled",
			Unit:           "Count",
			InstrumentType: "counter",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "NetworkIn",
			Description:    "The number of bytes received on all network interfaces by the instance",
			Unit:           "Bytes",
			InstrumentType: "counter",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "NetworkOut",
			Description:    "The number of bytes sent out on all network interfaces by the instance",
			Unit:           "Bytes",
			InstrumentType: "counter",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "NetworkPacketsIn",
			Description:    "The number of packets received on all network interfaces by the instance",
			Unit:           "Count",
			InstrumentType: "counter",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "NetworkPacketsOut",
			Description:    "The number of packets sent out on all network interfaces by the instance",
			Unit:           "Count",
			InstrumentType: "counter",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "GPUPowerUtilization",
			Description:    "Active power usage as percentage of maximum for accelerated computing instances",
			Unit:           "Percent",
			InstrumentType: "gauge",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "DedicatedHostCPUUtilization",
			Description:    "The percentage of allocated compute capacity in use on a Dedicated Host",
			Unit:           "Percent",
			InstrumentType: "gauge",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
	}
}

func cpuCreditMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{
			Name:           "CPUCreditUsage",
			Description:    "The number of CPU credits spent by the instance for CPU utilization (burstable instances)",
			Unit:           "Credits",
			InstrumentType: "counter",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "CPUCreditBalance",
			Description:    "The number of earned CPU credits that an instance has accrued since it was launched or started (burstable instances)",
			Unit:           "Credits",
			InstrumentType: "gauge",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "CPUSurplusCreditBalance",
			Description:    "The number of surplus credits that have been spent by an unlimited instance when its CPUCreditBalance value is zero",
			Unit:           "Credits",
			InstrumentType: "gauge",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "CPUSurplusCreditsCharged",
			Description:    "The number of spent surplus credits that are not paid down by earned CPU credits, and which thus incur an additional charge",
			Unit:           "Credits",
			InstrumentType: "counter",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
	}
}

func ebsMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{
			Name:           "EBSReadOps",
			Description:    "Completed read operations from all EBS volumes attached to the instance (Nitro instances)",
			Unit:           "Count",
			InstrumentType: "counter",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "EBSWriteOps",
			Description:    "Completed write operations to all EBS volumes attached to the instance (Nitro instances)",
			Unit:           "Count",
			InstrumentType: "counter",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "EBSReadBytes",
			Description:    "Bytes read from all EBS volumes attached to the instance (Nitro instances)",
			Unit:           "Bytes",
			InstrumentType: "counter",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "EBSWriteBytes",
			Description:    "Bytes written to all EBS volumes attached to the instance (Nitro instances)",
			Unit:           "Bytes",
			InstrumentType: "counter",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "EBSIOBalance%",
			Description:    "Percentage of I/O credits remaining in the burst bucket (Nitro instances)",
			Unit:           "Percent",
			InstrumentType: "gauge",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "EBSByteBalance%",
			Description:    "Percentage of throughput credits remaining in the burst bucket (Nitro instances)",
			Unit:           "Percent",
			InstrumentType: "gauge",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "InstanceEBSIOPSExceededCheck",
			Description:    "Returns 1 if the instance has exceeded the IOPS limit, otherwise returns 0",
			Unit:           "None",
			InstrumentType: "gauge",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "InstanceEBSThroughputExceededCheck",
			Description:    "Returns 1 if the instance has exceeded the throughput limit, otherwise returns 0",
			Unit:           "None",
			InstrumentType: "gauge",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
	}
}

func statusCheckMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{
			Name:           "StatusCheckFailed",
			Description:    "Reports whether the instance has passed both the instance status check and the system status check in the last minute",
			Unit:           "Count",
			InstrumentType: "gauge",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "StatusCheckFailed_Instance",
			Description:    "Reports whether the instance has passed the instance status check in the last minute",
			Unit:           "Count",
			InstrumentType: "gauge",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "StatusCheckFailed_System",
			Description:    "Reports whether the instance has passed the system status check in the last minute",
			Unit:           "Count",
			InstrumentType: "gauge",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
		{
			Name:           "StatusCheckFailed_AttachedEBS",
			Description:    "Reports whether the instance has passed the attached EBS status check in the last minute",
			Unit:           "Count",
			InstrumentType: "gauge",
			ComponentName:  "EC2",
			ComponentType:  "platform",
			SourceLocation: "AWS/EC2",
		},
	}
}
