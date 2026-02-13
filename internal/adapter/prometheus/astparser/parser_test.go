package astparser

import (
	"testing"
)

func TestParseFile_SimpleNewDesc(t *testing.T) {
	src := `
package collector

import "github.com/prometheus/client_golang/prometheus"

var metricDesc = prometheus.NewDesc(
	"pg_up",
	"Whether the PostgreSQL server is up",
	nil,
	nil,
)
`
	metrics, err := ParseSource("test.go", []byte(src))
	if err != nil {
		t.Fatalf("ParseSource failed: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	m := metrics[0]
	if m.Name != "pg_up" {
		t.Errorf("expected name 'pg_up', got '%s'", m.Name)
	}
	if m.Help != "Whether the PostgreSQL server is up" {
		t.Errorf("expected help 'Whether the PostgreSQL server is up', got '%s'", m.Help)
	}
	if len(m.Labels) != 0 {
		t.Errorf("expected 0 labels, got %d", len(m.Labels))
	}
}

func TestParseFile_WithLabels(t *testing.T) {
	src := `
package collector

import "github.com/prometheus/client_golang/prometheus"

var metricDesc = prometheus.NewDesc(
	"pg_stat_database_numbackends",
	"Number of backends currently connected",
	[]string{"datid", "datname"},
	nil,
)
`
	metrics, err := ParseSource("test.go", []byte(src))
	if err != nil {
		t.Fatalf("ParseSource failed: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	m := metrics[0]
	if m.Name != "pg_stat_database_numbackends" {
		t.Errorf("expected name 'pg_stat_database_numbackends', got '%s'", m.Name)
	}
	if len(m.Labels) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(m.Labels))
	}
	if m.Labels[0] != "datid" || m.Labels[1] != "datname" {
		t.Errorf("expected labels [datid, datname], got %v", m.Labels)
	}
}

func TestParseFile_BuildFQName(t *testing.T) {
	src := `
package collector

import "github.com/prometheus/client_golang/prometheus"

const namespace = "pg"
const subsystem = "stat"

var metricDesc = prometheus.NewDesc(
	prometheus.BuildFQName(namespace, subsystem, "connections"),
	"Number of connections",
	[]string{"state"},
	nil,
)
`
	metrics, err := ParseSource("test.go", []byte(src))
	if err != nil {
		t.Fatalf("ParseSource failed: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	m := metrics[0]
	if m.Name != "pg_stat_connections" {
		t.Errorf("expected name 'pg_stat_connections', got '%s'", m.Name)
	}
}

func TestParseFile_BuildFQNameWithLiterals(t *testing.T) {
	src := `
package collector

import "github.com/prometheus/client_golang/prometheus"

var metricDesc = prometheus.NewDesc(
	prometheus.BuildFQName("pg", "replication", "lag_bytes"),
	"Replication lag in bytes",
	[]string{"slot_name"},
	nil,
)
`
	metrics, err := ParseSource("test.go", []byte(src))
	if err != nil {
		t.Fatalf("ParseSource failed: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	m := metrics[0]
	if m.Name != "pg_replication_lag_bytes" {
		t.Errorf("expected name 'pg_replication_lag_bytes', got '%s'", m.Name)
	}
}

func TestParseFile_MultipleMetrics(t *testing.T) {
	src := `
package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	metric1 = prometheus.NewDesc(
		"pg_up",
		"Whether PostgreSQL is up",
		nil,
		nil,
	)
	metric2 = prometheus.NewDesc(
		"pg_info",
		"PostgreSQL version info",
		[]string{"version"},
		nil,
	)
)
`
	metrics, err := ParseSource("test.go", []byte(src))
	if err != nil {
		t.Fatalf("ParseSource failed: %v", err)
	}

	if len(metrics) != 2 {
		t.Fatalf("expected 2 metrics, got %d", len(metrics))
	}

	names := make(map[string]bool)
	for _, m := range metrics {
		names[m.Name] = true
	}
	if !names["pg_up"] || !names["pg_info"] {
		t.Errorf("expected metrics pg_up and pg_info, got %v", names)
	}
}

func TestParseFile_EmptySubsystem(t *testing.T) {
	src := `
package collector

import "github.com/prometheus/client_golang/prometheus"

var metricDesc = prometheus.NewDesc(
	prometheus.BuildFQName("pg", "", "up"),
	"Whether PostgreSQL is up",
	nil,
	nil,
)
`
	metrics, err := ParseSource("test.go", []byte(src))
	if err != nil {
		t.Fatalf("ParseSource failed: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	m := metrics[0]
	if m.Name != "pg_up" {
		t.Errorf("expected name 'pg_up', got '%s'", m.Name)
	}
}

func TestParseFile_LabelsFromLocalVariable(t *testing.T) {
	src := `
package collector

import "github.com/prometheus/client_golang/prometheus"

func createCollector(system string) {
	summaryLabels := []string{"server_id"}
	_ = prometheus.NewDesc(
		prometheus.BuildFQName("ns", "sub", "metric"),
		"A metric with labels from a local variable",
		summaryLabels,
		nil,
	)
}
`
	metrics, err := ParseSource("test.go", []byte(src))
	if err != nil {
		t.Fatalf("ParseSource failed: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	m := metrics[0]
	if len(m.Labels) != 1 || m.Labels[0] != "server_id" {
		t.Errorf("expected labels [server_id], got %v", m.Labels)
	}
}

func TestParseFile_NoMetrics(t *testing.T) {
	src := `
package collector

import "fmt"

func doSomething() {
	fmt.Println("hello")
}
`
	metrics, err := ParseSource("test.go", []byte(src))
	if err != nil {
		t.Fatalf("ParseSource failed: %v", err)
	}

	if len(metrics) != 0 {
		t.Errorf("expected 0 metrics, got %d", len(metrics))
	}
}
