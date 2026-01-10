package cadvisor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "test.go")
	content := `package metrics

import "github.com/prometheus/client_golang/prometheus"

type containerMetric struct {
	name        string
	help        string
	valueType   prometheus.ValueType
	extraLabels []string
}

func NewCollector() {
	c := &Collector{
		containerMetrics: []containerMetric{
			{
				name:      "container_cpu_user_seconds_total",
				help:      "Cumulative user cpu time consumed in seconds.",
				valueType: prometheus.CounterValue,
			},
			{
				name:        "container_cpu_usage_seconds_total",
				help:        "Cumulative cpu time consumed in seconds.",
				valueType:   prometheus.CounterValue,
				extraLabels: []string{"cpu"},
			},
			{
				name:      "container_memory_usage_bytes",
				help:      "Current memory usage in bytes, including all memory regardless of when it was accessed",
				valueType: prometheus.GaugeValue,
			},
		},
	}
}
`
	if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	defs, err := ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(defs) != 3 {
		t.Errorf("expected 3 metrics, got %d", len(defs))
		for _, d := range defs {
			t.Logf("  metric: %s", d.Name)
		}
	}

	expected := map[string]struct {
		help       string
		metricType string
		labels     []string
	}{
		"container_cpu_user_seconds_total":  {"Cumulative user cpu time consumed in seconds.", "counter", nil},
		"container_cpu_usage_seconds_total": {"Cumulative cpu time consumed in seconds.", "counter", []string{"cpu"}},
		"container_memory_usage_bytes":      {"Current memory usage in bytes, including all memory regardless of when it was accessed", "gauge", nil},
	}

	for _, def := range defs {
		exp, ok := expected[def.Name]
		if !ok {
			t.Errorf("unexpected metric: %s", def.Name)
			continue
		}

		if def.Help != exp.help {
			t.Errorf("metric %s: expected help %q, got %q", def.Name, exp.help, def.Help)
		}

		if def.MetricType != exp.metricType {
			t.Errorf("metric %s: expected type %q, got %q", def.Name, exp.metricType, def.MetricType)
		}
	}
}

func TestParseFileMachineMetrics(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "test.go")
	content := `package metrics

import "github.com/prometheus/client_golang/prometheus"

type machineMetric struct {
	name        string
	help        string
	valueType   prometheus.ValueType
	extraLabels []string
}

func NewMachineCollector() {
	c := &Collector{
		machineMetrics: []machineMetric{
			{
				name:      "machine_cpu_physical_cores",
				help:      "Number of physical CPU cores.",
				valueType: prometheus.GaugeValue,
			},
			{
				name:      "machine_memory_bytes",
				help:      "Amount of memory installed on the machine.",
				valueType: prometheus.GaugeValue,
			},
		},
	}
}
`
	if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	defs, err := ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(defs) != 2 {
		t.Fatalf("expected 2 metrics, got %d", len(defs))
	}
}

func TestParseFileNonExistent(t *testing.T) {
	_, err := ParseFile("/nonexistent/file.go")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}
