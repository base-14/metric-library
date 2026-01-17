package rds

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
	return "cloudwatch-rds"
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
	return "https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-metrics.html"
}

func (a *Adapter) Fetch(_ context.Context, _ adapter.FetchOptions) (*adapter.FetchResult, error) {
	return &adapter.FetchResult{
		Commit:    time.Now().Format("2006-01-02"),
		Timestamp: time.Now(),
	}, nil
}

func (a *Adapter) Extract(_ context.Context, _ *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	metrics := []*adapter.RawMetric{}
	metrics = append(metrics, instanceMetrics()...)
	metrics = append(metrics, usageMetrics()...)
	return metrics, nil
}

func instanceMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "BinLogDiskUsage", Description: "The amount of disk space occupied by binary logs", Unit: "Bytes", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "BurstBalance", Description: "The percent of General Purpose SSD (gp2) burst-bucket I/O credits available", Unit: "Percent", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "CheckpointLag", Description: "The amount of time since the most recent checkpoint", Unit: "Seconds", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "ConnectionAttempts", Description: "The number of attempts to connect to an instance, whether successful or not", Unit: "Count", InstrumentType: "counter", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "CPUUtilization", Description: "The percentage of CPU utilization", Unit: "Percent", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "CPUCreditUsage", Description: "The number of CPU credits spent by the instance for CPU utilization", Unit: "Credits", InstrumentType: "counter", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "CPUCreditBalance", Description: "The number of earned CPU credits that an instance has accrued", Unit: "Credits", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "CPUSurplusCreditBalance", Description: "The number of surplus credits spent by an unlimited instance", Unit: "Credits", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "CPUSurplusCreditsCharged", Description: "The number of spent surplus credits that incur an additional charge", Unit: "Credits", InstrumentType: "counter", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "DatabaseConnections", Description: "The number of client network connections to the database instance", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "DiskQueueDepth", Description: "The number of outstanding I/Os waiting to access the disk", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "DiskQueueDepthLogVolume", Description: "The number of outstanding I/Os waiting to access the log volume disk", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "EBSByteBalance%", Description: "The percentage of throughput credits remaining in the burst bucket", Unit: "Percent", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "EBSIOBalance%", Description: "The percentage of I/O credits remaining in the burst bucket", Unit: "Percent", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "FailedSQLServerAgentJobsCount", Description: "The number of failed Microsoft SQL Server Agent jobs during the last minute", Unit: "Count", InstrumentType: "counter", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "FreeableMemory", Description: "The amount of available random access memory", Unit: "Bytes", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "FreeLocalStorage", Description: "The amount of available local storage space", Unit: "Bytes", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "FreeLocalStoragePercent", Description: "The percentage of available local storage space", Unit: "Percent", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "FreeStorageSpace", Description: "The amount of available storage space", Unit: "Bytes", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "FreeStorageSpaceLogVolume", Description: "The amount of available storage space on the log volume", Unit: "Bytes", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "IamDbAuthConnectionRequests", Description: "The number of connection requests using IAM authentication", Unit: "Count", InstrumentType: "counter", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "MaximumUsedTransactionIDs", Description: "The maximum transaction IDs that have been used", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "NetworkReceiveThroughput", Description: "The incoming network traffic on the DB instance", Unit: "Bytes/Second", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "NetworkTransmitThroughput", Description: "The outgoing network traffic on the DB instance", Unit: "Bytes/Second", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "OldestLogicalReplicationSlotLag", Description: "The lagging size of Amazon RDS commits on source vs replica", Unit: "Bytes", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "OldestReplicationSlotLag", Description: "The lagging size of the replica with most WAL data lag", Unit: "Bytes", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "ReadIOPS", Description: "The average number of disk read I/O operations per second", Unit: "Count/Second", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "ReadIOPSLocalStorage", Description: "The average number of disk read I/O operations to local storage per second", Unit: "Count/Second", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "ReadIOPSLogVolume", Description: "The average number of disk read I/O operations per second for the log volume", Unit: "Count/Second", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "ReadLatency", Description: "The average amount of time taken per disk I/O operation", Unit: "Seconds", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "ReadLatencyLocalStorage", Description: "The average amount of time taken per disk I/O operation for local storage", Unit: "Seconds", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "ReadLatencyLogVolume", Description: "The average amount of time taken per disk I/O operation for the log volume", Unit: "Seconds", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "ReadThroughput", Description: "The average number of bytes read from disk per second", Unit: "Bytes/Second", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "ReadThroughputLocalStorage", Description: "The average number of bytes read from disk per second for local storage", Unit: "Bytes/Second", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "ReadThroughputLogVolume", Description: "The average number of bytes read from disk per second for the log volume", Unit: "Bytes/Second", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "ReplicaLag", Description: "The amount of time a read replica lags behind the source DB instance", Unit: "Seconds", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "ReplicationChannelLag", Description: "The amount of time a multi-source replica channel lags behind the source", Unit: "Seconds", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "ReplicationSlotDiskUsage", Description: "The disk space used by replication slot files", Unit: "Bytes", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "SwapUsage", Description: "The amount of swap space used on the DB instance", Unit: "Bytes", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "TempDbAvailableDataSpace", Description: "The amount of available data space on tempdb", Unit: "Bytes", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "TempDbAvailableLogSpace", Description: "The amount of available log space on tempdb", Unit: "Bytes", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "TempDbDataFileUsage", Description: "The percentage of data files used on tempdb", Unit: "Percent", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "TempDbLogFileUsage", Description: "The percentage of log files used on tempdb", Unit: "Percent", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "TransactionLogsDiskUsage", Description: "The disk space used by transaction logs", Unit: "Bytes", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "TransactionLogsGeneration", Description: "The size of transaction logs generated per second", Unit: "Bytes/Second", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "WriteIOPS", Description: "The average number of disk write I/O operations per second", Unit: "Count/Second", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "WriteIOPSLocalStorage", Description: "The average number of disk write I/O operations per second on local storage", Unit: "Count/Second", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "WriteIOPSLogVolume", Description: "The average number of disk write I/O operations per second for the log volume", Unit: "Count/Second", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "WriteLatency", Description: "The average amount of time taken per disk I/O write operation", Unit: "Seconds", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "WriteLatencyLocalStorage", Description: "The average amount of time taken per disk I/O write operation on local storage", Unit: "Seconds", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "WriteLatencyLogVolume", Description: "The average amount of time taken per disk I/O write operation for the log volume", Unit: "Seconds", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "WriteThroughput", Description: "The average number of bytes written to disk per second", Unit: "Bytes/Second", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "WriteThroughputLocalStorage", Description: "The average number of bytes written to disk per second for local storage", Unit: "Bytes/Second", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
		{Name: "WriteThroughputLogVolume", Description: "The average number of bytes written to disk per second for the log volume", Unit: "Bytes/Second", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/RDS"},
	}
}

func usageMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "AllocatedStorage", Description: "The total storage for all DB instances", Unit: "Gigabytes", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/Usage"},
		{Name: "AuthorizationsPerDBSecurityGroup", Description: "The number of ingress rules per DB security group", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/Usage"},
		{Name: "CustomEndpointsPerDBCluster", Description: "The number of custom endpoints per DB cluster", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/Usage"},
		{Name: "CustomEngineVersions", Description: "The number of custom engine versions (CEVs) in your account", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/Usage"},
		{Name: "DBClusterParameterGroups", Description: "The number of DB cluster parameter groups in your account", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/Usage"},
		{Name: "DBClusterRoles", Description: "The number of associated IAM roles per DB cluster", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/Usage"},
		{Name: "DBClusters", Description: "The number of Amazon Aurora DB clusters in your account", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/Usage"},
		{Name: "DBInstanceRoles", Description: "The number of associated IAM roles per DB instance", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/Usage"},
		{Name: "DBInstances", Description: "The number of DB instances in your account", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/Usage"},
		{Name: "DBParameterGroups", Description: "The number of DB parameter groups in your account", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/Usage"},
		{Name: "DBSecurityGroups", Description: "The number of security groups in your account", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/Usage"},
		{Name: "DBSubnetGroups", Description: "The number of DB subnet groups in your account", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/Usage"},
		{Name: "EventSubscriptions", Description: "The number of event notification subscriptions in your account", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/Usage"},
		{Name: "Integrations", Description: "The number of zero-ETL integrations with Amazon Redshift", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/Usage"},
		{Name: "ManualClusterSnapshots", Description: "The number of manually created DB cluster snapshots", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/Usage"},
		{Name: "ManualSnapshots", Description: "The number of manually created DB snapshots", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/Usage"},
		{Name: "OptionGroups", Description: "The number of option groups in your account", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/Usage"},
		{Name: "Proxies", Description: "The number of RDS proxies in your account", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/Usage"},
		{Name: "ReadReplicasPerMaster", Description: "The number of read replicas per DB instance", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/Usage"},
		{Name: "ReservedDBInstances", Description: "The number of reserved DB instances in your account", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/Usage"},
		{Name: "SubnetsPerDBSubnetGroup", Description: "The number of subnets per DB subnet group", Unit: "Count", InstrumentType: "gauge", ComponentName: "RDS", ComponentType: "platform", SourceLocation: "AWS/Usage"},
	}
}
