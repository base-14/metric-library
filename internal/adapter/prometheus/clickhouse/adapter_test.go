package clickhouse

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
)

func TestClickHouseAdapter_Name(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Name() != "prometheus-clickhouse" {
		t.Errorf("expected name 'prometheus-clickhouse', got %q", a.Name())
	}
}

func TestClickHouseAdapter_SourceCategory(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.SourceCategory() != domain.SourcePrometheus {
		t.Errorf("expected source category 'prometheus', got %q", a.SourceCategory())
	}
}

func TestClickHouseAdapter_Confidence(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.Confidence() != domain.ConfidenceAuthoritative {
		t.Errorf("expected confidence 'authoritative', got %q", a.Confidence())
	}
}

func TestClickHouseAdapter_ExtractionMethod(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.ExtractionMethod() != domain.ExtractionAST {
		t.Errorf("expected extraction method 'ast', got %q", a.ExtractionMethod())
	}
}

func TestClickHouseAdapter_RepoURL(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	if a.RepoURL() != "https://github.com/ClickHouse/ClickHouse" {
		t.Errorf("unexpected repo URL: %q", a.RepoURL())
	}
}

func TestClickHouseAdapter_ImplementsAdapter(t *testing.T) {
	a := NewAdapter("/tmp/cache")
	var _ adapter.Adapter = a
}

func TestClickHouseAdapter_Extract_CurrentMetrics(t *testing.T) {
	tmpDir := setupTestRepo(t)

	writeFile(t, filepath.Join(tmpDir, "src", "Common", "CurrentMetrics.cpp"), `#include <Common/CurrentMetrics.h>

#define APPLY_FOR_BUILTIN_METRICS(M) \
    M(Query, "Number of executing queries") \
    M(Merge, "Number of executing background merges") \
    M(MemoryTracking, "Total amount of memory (bytes) allocated by the server")
`)

	a := NewAdapter("/tmp/cache")
	metrics, err := a.Extract(context.Background(), &adapter.FetchResult{
		RepoPath: tmpDir,
		Commit:   "abc123",
	})
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 3 {
		t.Fatalf("expected 3 metrics, got %d", len(metrics))
	}

	names := metricsByName(metrics)

	m := names["ClickHouseMetrics_Query"]
	if m == nil {
		t.Fatal("expected metric 'ClickHouseMetrics_Query'")
	}
	if m.Description != "Number of executing queries" {
		t.Errorf("unexpected description: %q", m.Description)
	}
	if m.InstrumentType != "gauge" {
		t.Errorf("expected gauge, got %q", m.InstrumentType)
	}
	if m.ComponentName != "current_metrics" {
		t.Errorf("expected component 'current_metrics', got %q", m.ComponentName)
	}
}

func TestClickHouseAdapter_Extract_ProfileEvents(t *testing.T) {
	tmpDir := setupTestRepo(t)

	writeFile(t, filepath.Join(tmpDir, "src", "Common", "ProfileEvents.cpp"), `#include <Common/ProfileEvents.h>

#define APPLY_FOR_BUILTIN_EVENTS(M) \
    M(Query, "Number of queries to be interpreted.", ValueType::Number) \
    M(ReadBufferFromFileDescriptorReadBytes, "Number of bytes read from file descriptors.", ValueType::Bytes) \
    M(QueryTimeMicroseconds, "Total time of all queries.", ValueType::Microseconds)
`)

	a := NewAdapter("/tmp/cache")
	metrics, err := a.Extract(context.Background(), &adapter.FetchResult{
		RepoPath: tmpDir,
		Commit:   "abc123",
	})
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 3 {
		t.Fatalf("expected 3 metrics, got %d", len(metrics))
	}

	names := metricsByName(metrics)

	m := names["ClickHouseProfileEvents_Query"]
	if m == nil {
		t.Fatal("expected metric 'ClickHouseProfileEvents_Query'")
	}
	if m.InstrumentType != "counter" {
		t.Errorf("expected counter, got %q", m.InstrumentType)
	}
	if m.ComponentName != "profile_events" {
		t.Errorf("expected component 'profile_events', got %q", m.ComponentName)
	}
	if m.Unit != "" {
		t.Errorf("expected no unit for Number type, got %q", m.Unit)
	}

	bytesMetric := names["ClickHouseProfileEvents_ReadBufferFromFileDescriptorReadBytes"]
	if bytesMetric == nil {
		t.Fatal("expected bytes metric")
	}
	if bytesMetric.Unit != "bytes" {
		t.Errorf("expected unit 'bytes', got %q", bytesMetric.Unit)
	}

	timeMetric := names["ClickHouseProfileEvents_QueryTimeMicroseconds"]
	if timeMetric == nil {
		t.Fatal("expected time metric")
	}
	if timeMetric.Unit != "microseconds" {
		t.Errorf("expected unit 'microseconds', got %q", timeMetric.Unit)
	}
}

func TestClickHouseAdapter_Extract_AsyncMetrics(t *testing.T) {
	tmpDir := setupTestRepo(t)

	writeFile(t, filepath.Join(tmpDir, "src", "Interpreters", "ServerAsynchronousMetrics.cpp"), `
void ServerAsynchronousMetrics::updateImpl() {
    new_values["Uptime"] = { getContext()->getUptimeSeconds(),
        "Server uptime in seconds." };
    new_values["MaxPartCountForPartition"] = { max_part_count_for_partition,
        "Maximum number of parts per partition." };
}
`)

	writeFile(t, filepath.Join(tmpDir, "src", "Common", "AsynchronousMetrics.cpp"), `
void AsynchronousMetrics::update() {
    new_values["OSMemoryTotal"] = { mem_total,
        "Total memory on the host system." };
    new_values["CGroupUserTime"]
        = { value, "User time in microseconds." };
}
`)

	a := NewAdapter("/tmp/cache")
	metrics, err := a.Extract(context.Background(), &adapter.FetchResult{
		RepoPath: tmpDir,
		Commit:   "abc123",
	})
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 4 {
		t.Fatalf("expected 4 metrics, got %d", len(metrics))
	}

	names := metricsByName(metrics)

	m := names["ClickHouseAsyncMetrics_Uptime"]
	if m == nil {
		t.Fatal("expected metric 'ClickHouseAsyncMetrics_Uptime'")
	}
	if m.Description != "Server uptime in seconds." {
		t.Errorf("unexpected description: %q", m.Description)
	}
	if m.InstrumentType != "gauge" {
		t.Errorf("expected gauge, got %q", m.InstrumentType)
	}
	if m.ComponentName != "async_metrics" {
		t.Errorf("expected component 'async_metrics', got %q", m.ComponentName)
	}

	if names["ClickHouseAsyncMetrics_OSMemoryTotal"] == nil {
		t.Error("expected metric 'ClickHouseAsyncMetrics_OSMemoryTotal'")
	}
	if names["ClickHouseAsyncMetrics_CGroupUserTime"] == nil {
		t.Error("expected metric 'ClickHouseAsyncMetrics_CGroupUserTime'")
	}
}

func TestClickHouseAdapter_Extract_SkipsDynamicAsyncMetrics(t *testing.T) {
	tmpDir := setupTestRepo(t)

	writeFile(t, filepath.Join(tmpDir, "src", "Common", "AsynchronousMetrics.cpp"), `
void AsynchronousMetrics::update() {
    new_values[fmt::format("BlockReadOps_{}", device_name)] = { value, "Read ops per device." };
    new_values["OSMemoryTotal"] = { mem_total, "Total memory." };
}
`)

	a := NewAdapter("/tmp/cache")
	metrics, err := a.Extract(context.Background(), &adapter.FetchResult{
		RepoPath: tmpDir,
		Commit:   "abc123",
	})
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	// Only the static metric should be extracted, not the fmt::format one
	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric (dynamic skipped), got %d", len(metrics))
	}

	if metrics[0].Name != "ClickHouseAsyncMetrics_OSMemoryTotal" {
		t.Errorf("expected 'ClickHouseAsyncMetrics_OSMemoryTotal', got %q", metrics[0].Name)
	}
}

func TestClickHouseAdapter_Extract_AllSources(t *testing.T) {
	tmpDir := setupTestRepo(t)

	writeFile(t, filepath.Join(tmpDir, "src", "Common", "CurrentMetrics.cpp"), `
#define APPLY_FOR_BUILTIN_METRICS(M) \
    M(Query, "Number of executing queries") \
    M(Merge, "Number of executing background merges")
`)

	writeFile(t, filepath.Join(tmpDir, "src", "Common", "ProfileEvents.cpp"), `
#define APPLY_FOR_BUILTIN_EVENTS(M) \
    M(SelectQuery, "SELECT queries.", ValueType::Number) \
    M(InsertQuery, "INSERT queries.", ValueType::Number)
`)

	writeFile(t, filepath.Join(tmpDir, "src", "Interpreters", "ServerAsynchronousMetrics.cpp"), `
void update() {
    new_values["Uptime"] = { val, "Server uptime." };
}
`)

	writeFile(t, filepath.Join(tmpDir, "src", "Common", "AsynchronousMetrics.cpp"), `
void update() {
    new_values["OSMemoryTotal"] = { val, "Total memory." };
}
`)

	a := NewAdapter("/tmp/cache")
	metrics, err := a.Extract(context.Background(), &adapter.FetchResult{
		RepoPath: tmpDir,
		Commit:   "abc123",
	})
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 6 {
		t.Fatalf("expected 6 metrics total, got %d", len(metrics))
	}

	names := metricsByName(metrics)

	// current_metrics (gauges)
	if names["ClickHouseMetrics_Query"] == nil {
		t.Error("missing ClickHouseMetrics_Query")
	}
	// profile_events (counters)
	if names["ClickHouseProfileEvents_SelectQuery"] == nil {
		t.Error("missing ClickHouseProfileEvents_SelectQuery")
	}
	// async_metrics (gauges)
	if names["ClickHouseAsyncMetrics_Uptime"] == nil {
		t.Error("missing ClickHouseAsyncMetrics_Uptime")
	}
	if names["ClickHouseAsyncMetrics_OSMemoryTotal"] == nil {
		t.Error("missing ClickHouseAsyncMetrics_OSMemoryTotal")
	}
}

func TestClickHouseAdapter_Extract_EmptyRepo(t *testing.T) {
	tmpDir := setupTestRepo(t)

	a := NewAdapter("/tmp/cache")
	metrics, err := a.Extract(context.Background(), &adapter.FetchResult{
		RepoPath: tmpDir,
		Commit:   "abc123",
	})
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 0 {
		t.Errorf("expected 0 metrics for empty repo, got %d", len(metrics))
	}
}

func TestParseCurrentMetrics(t *testing.T) {
	src := []byte(`
#define APPLY_FOR_BUILTIN_METRICS(M) \
    M(Query, "Number of executing queries") \
    M(Merge, "Number of executing background merges") \
    M(MemoryTracking, "Total amount of memory (bytes) allocated by the server")
`)

	defs := parseCurrentMetrics(src)
	if len(defs) != 3 {
		t.Fatalf("expected 3 metrics, got %d", len(defs))
	}

	if defs[0].Name != "Query" || defs[0].Description != "Number of executing queries" {
		t.Errorf("unexpected first metric: %+v", defs[0])
	}
}

func TestParseProfileEvents(t *testing.T) {
	src := []byte(`
#define APPLY_FOR_BUILTIN_EVENTS(M) \
    M(Query, "Number of queries.", ValueType::Number) \
    M(ReadBytes, "Bytes read.", ValueType::Bytes) \
    M(ElapsedMicroseconds, "Time elapsed.", ValueType::Microseconds) \
    M(ElapsedMilliseconds, "Time elapsed ms.", ValueType::Milliseconds) \
    M(ElapsedNanoseconds, "Time elapsed ns.", ValueType::Nanoseconds)
`)

	defs := parseProfileEvents(src)
	if len(defs) != 5 {
		t.Fatalf("expected 5 events, got %d", len(defs))
	}

	if defs[0].Unit != "" {
		t.Errorf("expected no unit for Number, got %q", defs[0].Unit)
	}
	if defs[1].Unit != "bytes" {
		t.Errorf("expected 'bytes', got %q", defs[1].Unit)
	}
	if defs[2].Unit != "microseconds" {
		t.Errorf("expected 'microseconds', got %q", defs[2].Unit)
	}
	if defs[3].Unit != "milliseconds" {
		t.Errorf("expected 'milliseconds', got %q", defs[3].Unit)
	}
	if defs[4].Unit != "nanoseconds" {
		t.Errorf("expected 'nanoseconds', got %q", defs[4].Unit)
	}
}

func TestParseAsyncMetrics(t *testing.T) {
	src := []byte(`
void update() {
    new_values["Uptime"] = { val, "Server uptime in seconds." };
    new_values["MaxPartCount"] = { count,
        "Maximum number of parts." };
    new_values[fmt::format("Disk_{}", name)] = { val, "Dynamic metric." };
    new_values["OSMemoryTotal"] = { val, "Total memory." };
}
`)

	defs := parseAsyncMetrics(src)
	if len(defs) != 3 {
		t.Fatalf("expected 3 static metrics (dynamic skipped), got %d", len(defs))
	}

	names := make(map[string]bool)
	for _, d := range defs {
		names[d.Name] = true
	}
	if !names["Uptime"] {
		t.Error("missing Uptime")
	}
	if !names["MaxPartCount"] {
		t.Error("missing MaxPartCount")
	}
	if !names["OSMemoryTotal"] {
		t.Error("missing OSMemoryTotal")
	}
}

func TestValueTypeToUnit(t *testing.T) {
	tests := []struct {
		valueType string
		expected  string
	}{
		{"Number", ""},
		{"Bytes", "bytes"},
		{"Microseconds", "microseconds"},
		{"Milliseconds", "milliseconds"},
		{"Nanoseconds", "nanoseconds"},
		{"Unknown", ""},
	}

	for _, tt := range tests {
		t.Run(tt.valueType, func(t *testing.T) {
			got := valueTypeToUnit(tt.valueType)
			if got != tt.expected {
				t.Errorf("valueTypeToUnit(%q) = %q, want %q", tt.valueType, got, tt.expected)
			}
		})
	}
}

// Test helpers

func setupTestRepo(t *testing.T) string {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "clickhouse-adapter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tmpDir) })

	for _, dir := range []string{
		filepath.Join(tmpDir, "src", "Common"),
		filepath.Join(tmpDir, "src", "Interpreters"),
	} {
		if err := os.MkdirAll(dir, 0750); err != nil {
			t.Fatalf("failed to create dir %s: %v", dir, err)
		}
	}

	return tmpDir
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write %s: %v", path, err)
	}
}

func metricsByName(metrics []*adapter.RawMetric) map[string]*adapter.RawMetric {
	m := make(map[string]*adapter.RawMetric)
	for _, metric := range metrics {
		m[metric.Name] = metric
	}
	return m
}
