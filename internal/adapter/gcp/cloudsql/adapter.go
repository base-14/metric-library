package cloudsql

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
	return "gcp-cloudsql"
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
	return "https://cloud.google.com/monitoring/api/metrics_gcp#gcp-cloudsql"
}

func (a *Adapter) Fetch(_ context.Context, _ adapter.FetchOptions) (*adapter.FetchResult, error) {
	return &adapter.FetchResult{
		Commit:    time.Now().Format("2006-01-02"),
		Timestamp: time.Now(),
	}, nil
}

func (a *Adapter) Extract(_ context.Context, _ *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	var metrics []*adapter.RawMetric
	metrics = append(metrics, generalMetrics()...)
	metrics = append(metrics, cpuMetrics()...)
	metrics = append(metrics, memoryMetrics()...)
	metrics = append(metrics, diskMetrics()...)
	metrics = append(metrics, networkMetrics()...)
	metrics = append(metrics, replicationMetrics()...)
	metrics = append(metrics, mysqlMetrics()...)
	metrics = append(metrics, postgresMetrics()...)
	return metrics, nil
}

func generalMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "cloudsql.googleapis.com/database/up", Description: "Indicates if the server is up or not. On-demand instances are spun down if no connections are made for a sufficient amount of time", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/uptime", Description: "Delta count of the time in seconds the instance has been running", Unit: "s", InstrumentType: "counter", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/state", Description: "The current serving state of the Cloud SQL instance", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/instance_state", Description: "The current state of the instance", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/available_for_failover", Description: "Whether failover operation is available on the instance", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
	}
}

func cpuMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "cloudsql.googleapis.com/database/cpu/utilization", Description: "Current CPU utilization represented as a percentage of the reserved CPU that is currently in use", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/cpu/reserved_cores", Description: "Number of cores reserved for the database", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
	}
}

func memoryMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "cloudsql.googleapis.com/database/memory/utilization", Description: "Memory utilization represented as a percentage", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/memory/total_usage", Description: "Total RAM usage in bytes. This includes OS buffer and cache memory", Unit: "By", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/memory/usage", Description: "RAM usage in bytes excluding OS buffer and cache memory", Unit: "By", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/memory/quota", Description: "Maximum RAM size in bytes", Unit: "By", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
	}
}

func diskMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "cloudsql.googleapis.com/database/disk/bytes_used", Description: "Data utilization in bytes", Unit: "By", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/disk/quota", Description: "Maximum data disk size in bytes", Unit: "By", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/disk/utilization", Description: "The fraction of the disk quota that is currently in use", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/disk/read_ops_count", Description: "Delta count of data disk read IO operations", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/disk/write_ops_count", Description: "Delta count of data disk write IO operations", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/disk/bytes_used_by_data_type", Description: "Data utilization in bytes broken down by data type", Unit: "By", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
	}
}

func networkMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "cloudsql.googleapis.com/database/network/received_bytes_count", Description: "Delta count of bytes received through the network", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/network/sent_bytes_count", Description: "Delta count of bytes sent through the network", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/network/connections", Description: "Number of connections to databases on the Cloud SQL instance", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
	}
}

func replicationMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "cloudsql.googleapis.com/database/replication/replica_lag", Description: "Number of seconds the read replica is behind its primary", Unit: "s", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/replication/network_lag", Description: "Indicates time taken from primary binary log to IO thread on replica", Unit: "s", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/replication/replica_byte_count", Description: "Number of bytes that the replica has received from the primary", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/replication/state", Description: "The current state of replication", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
	}
}

func mysqlMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "cloudsql.googleapis.com/database/mysql/queries", Description: "Delta count of statements executed by the server", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/mysql/questions", Description: "Delta count of statements executed by the server sent by the client", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/mysql/received_bytes_count", Description: "Delta count of bytes received by the MySQL process", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/mysql/sent_bytes_count", Description: "Delta count of bytes sent by the MySQL process", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/mysql/innodb_buffer_pool_pages_dirty", Description: "Number of dirty pages in the InnoDB buffer pool", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/mysql/innodb_buffer_pool_pages_free", Description: "Number of free pages in the InnoDB buffer pool", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/mysql/innodb_buffer_pool_pages_total", Description: "Total number of pages in the InnoDB buffer pool", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/mysql/innodb_data_fsyncs", Description: "Delta count of InnoDB fsync() calls", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/mysql/innodb_os_log_fsyncs", Description: "Delta count of InnoDB log fsync() calls", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/mysql/innodb_pages_read", Description: "Delta count of InnoDB pages read", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/mysql/innodb_pages_written", Description: "Delta count of InnoDB pages written", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/mysql/replication_seconds_behind_master", Description: "Number of seconds the read replica is behind its primary (approximation)", Unit: "s", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
	}
}

func postgresMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "cloudsql.googleapis.com/database/postgresql/num_backends", Description: "Number of connections to the Cloud SQL PostgreSQL instance", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/postgresql/transaction_count", Description: "Delta count of number of transactions", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/postgresql/insights/aggregate/execution_time", Description: "Accumulated query execution time per user per database. This is the sum of cpu time, IO wait time, lock wait time, process context switch, and scheduling for all the processes involved in the query execution", Unit: "us", InstrumentType: "counter", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/postgresql/insights/aggregate/io_time", Description: "Accumulated IO time per user per database", Unit: "us", InstrumentType: "counter", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/postgresql/insights/aggregate/latencies", Description: "Query latency distribution per user per database", Unit: "us", InstrumentType: "histogram", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/postgresql/insights/aggregate/lock_time", Description: "Accumulated lock wait time per user per database", Unit: "us", InstrumentType: "counter", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/postgresql/insights/aggregate/row_count", Description: "Total number of rows affected during query execution", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/postgresql/insights/aggregate/shared_blk_access_count", Description: "Shared blocks (regular tables and indexed) accessed by statement execution", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/postgresql/replication/replica_byte_count", Description: "Number of bytes that the replica has received from the primary", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
		{Name: "cloudsql.googleapis.com/database/postgresql/vacuum/oldest_transaction_age", Description: "Age of the oldest transaction yet to be vacuumed in the Cloud SQL PostgreSQL instance", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud SQL", ComponentType: "platform", SourceLocation: "cloudsql.googleapis.com"},
	}
}
