package sqldatabase

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
	return "azure-sqldatabase"
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
	return "https://learn.microsoft.com/en-us/azure/azure-monitor/reference/supported-metrics/microsoft-sql-servers-databases-metrics"
}

func (a *Adapter) Fetch(_ context.Context, _ adapter.FetchOptions) (*adapter.FetchResult, error) {
	return &adapter.FetchResult{
		Commit:    time.Now().Format("2006-01-02"),
		Timestamp: time.Now(),
	}, nil
}

func (a *Adapter) Extract(_ context.Context, _ *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	var metrics []*adapter.RawMetric
	metrics = append(metrics, computeMetrics()...)
	metrics = append(metrics, storageMetrics()...)
	metrics = append(metrics, connectionMetrics()...)
	metrics = append(metrics, dtuMetrics()...)
	metrics = append(metrics, queryMetrics()...)
	metrics = append(metrics, replicationMetrics()...)
	metrics = append(metrics, tempdbMetrics()...)
	return metrics, nil
}

func computeMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.sqldatabase.cpu_percent", Description: "CPU percentage", Unit: "%", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.cpu_limit", Description: "CPU limit for vCore-based databases", Unit: "1", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.cpu_used", Description: "CPU used for vCore-based databases", Unit: "1", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.physical_data_read_percent", Description: "Data IO percentage", Unit: "%", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.log_write_percent", Description: "Log IO percentage", Unit: "%", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.workers_percent", Description: "Workers percentage", Unit: "%", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.sessions_percent", Description: "Sessions percentage", Unit: "%", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.sessions_count", Description: "Number of active sessions", Unit: "1", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
	}
}

func storageMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.sqldatabase.storage", Description: "Data space used", Unit: "By", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.storage_percent", Description: "Data space used percent", Unit: "%", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.allocated_data_storage", Description: "Data space allocated", Unit: "By", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.xtp_storage_percent", Description: "In-Memory OLTP storage percent", Unit: "%", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
	}
}

func connectionMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.sqldatabase.connection_successful", Description: "Successful connections", Unit: "1", InstrumentType: "counter", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.connection_failed", Description: "Failed connections", Unit: "1", InstrumentType: "counter", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.blocked_by_firewall", Description: "Connections blocked by firewall", Unit: "1", InstrumentType: "counter", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.deadlock", Description: "Deadlocks", Unit: "1", InstrumentType: "counter", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
	}
}

func dtuMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.sqldatabase.dtu_consumption_percent", Description: "DTU percentage", Unit: "%", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.dtu_limit", Description: "DTU limit", Unit: "1", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.dtu_used", Description: "DTU used", Unit: "1", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.edtu_limit", Description: "eDTU limit for elastic pool databases", Unit: "1", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.edtu_used", Description: "eDTU used for elastic pool databases", Unit: "1", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
	}
}

func queryMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.sqldatabase.dwu_limit", Description: "DWU limit for data warehouse databases", Unit: "1", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.dwu_consumption_percent", Description: "DWU percentage for data warehouse databases", Unit: "%", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.dwu_used", Description: "DWU used for data warehouse databases", Unit: "1", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.full_backup_size_bytes", Description: "Cumulative full backup storage size", Unit: "By", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.diff_backup_size_bytes", Description: "Cumulative differential backup storage size", Unit: "By", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.log_backup_size_bytes", Description: "Cumulative log backup storage size", Unit: "By", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.sqlserver_process_core_percent", Description: "CPU usage as a percentage of the SQL DB process", Unit: "%", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.sqlserver_process_memory_percent", Description: "Memory usage as a percentage of the SQL DB process", Unit: "%", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
	}
}

func replicationMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.sqldatabase.geo_replication_lag_seconds", Description: "Geo-replication lag in seconds", Unit: "s", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.active_geo_replication_health", Description: "Health status of active geo-replication", Unit: "1", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.ledger_digest_upload_success", Description: "Successful ledger digest uploads", Unit: "1", InstrumentType: "counter", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.ledger_digest_upload_failed", Description: "Failed ledger digest uploads", Unit: "1", InstrumentType: "counter", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
	}
}

func tempdbMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "azure.sqldatabase.tempdb_data_size", Description: "Space used in tempdb data files in kilobytes", Unit: "KBy", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.tempdb_log_size", Description: "Space used in tempdb transaction log file in kilobytes", Unit: "KBy", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.tempdb_log_used_percent", Description: "Space used percentage in tempdb transaction log file", Unit: "%", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.app_cpu_billed", Description: "App CPU billed for serverless databases", Unit: "1", InstrumentType: "counter", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.app_memory_percent", Description: "App memory used percentage for serverless databases", Unit: "%", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
		{Name: "azure.sqldatabase.app_cpu_percent", Description: "App CPU percentage for serverless databases", Unit: "%", InstrumentType: "gauge", ComponentName: "SQL Database", ComponentType: "platform", SourceLocation: "Microsoft.Sql/servers/databases"},
	}
}
