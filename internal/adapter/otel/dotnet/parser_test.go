package dotnet

import (
	"testing"
)

func TestParseContent_CreateCounter(t *testing.T) {
	content := `
_meter.CreateCounter<long>(
    name: "http.server.requests",
    unit: "{request}",
    description: "Number of HTTP server requests");
`
	metrics, err := parseContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	m := metrics[0]
	if m.Name != "http.server.requests" {
		t.Errorf("expected name 'http.server.requests', got '%s'", m.Name)
	}
	if m.InstrumentType != "counter" {
		t.Errorf("expected type 'counter', got '%s'", m.InstrumentType)
	}
	if m.Unit != "{request}" {
		t.Errorf("expected unit '{request}', got '%s'", m.Unit)
	}
	if m.Description != "Number of HTTP server requests" {
		t.Errorf("expected description 'Number of HTTP server requests', got '%s'", m.Description)
	}
}

func TestParseContent_CreateHistogram(t *testing.T) {
	content := `
meter.CreateHistogram<double>(
    "http.server.request.duration",
    unit: "s",
    description: "Duration of HTTP requests");
`
	metrics, err := parseContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	m := metrics[0]
	if m.InstrumentType != "histogram" {
		t.Errorf("expected type 'histogram', got '%s'", m.InstrumentType)
	}
	if m.Name != "http.server.request.duration" {
		t.Errorf("expected name 'http.server.request.duration', got '%s'", m.Name)
	}
}

func TestParseContent_CreateObservableGauge(t *testing.T) {
	content := `
this.CreateObservableGauge<double>(
    "process.cpu.utilization",
    unit: "%",
    description: "CPU utilization");
`
	metrics, err := parseContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	m := metrics[0]
	if m.InstrumentType != "gauge" {
		t.Errorf("expected type 'gauge', got '%s'", m.InstrumentType)
	}
	if m.Name != "process.cpu.utilization" {
		t.Errorf("expected name 'process.cpu.utilization', got '%s'", m.Name)
	}
}

func TestParseContent_CreateUpDownCounter(t *testing.T) {
	content := `
_meter.CreateUpDownCounter<long>("db.client.connections.usage",
    unit: "{connection}",
    description: "Active database connections");
`
	metrics, err := parseContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	m := metrics[0]
	if m.InstrumentType != "updowncounter" {
		t.Errorf("expected type 'updowncounter', got '%s'", m.InstrumentType)
	}
}

func TestParseContent_CreateObservableCounter(t *testing.T) {
	content := `
Meter.CreateObservableCounter<long>(
    "process.runtime.dotnet.gc.collections.count",
    () => GetGCCollectionCounts(),
    unit: "{collections}",
    description: "Number of garbage collections");
`
	metrics, err := parseContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	m := metrics[0]
	if m.InstrumentType != "counter" {
		t.Errorf("expected type 'counter', got '%s'", m.InstrumentType)
	}
	if m.Name != "process.runtime.dotnet.gc.collections.count" {
		t.Errorf("expected name 'process.runtime.dotnet.gc.collections.count', got '%s'", m.Name)
	}
}

func TestParseContent_MultipleMetrics(t *testing.T) {
	content := `
public class RuntimeMetrics
{
    private void RegisterMetrics()
    {
        _meter.CreateCounter<long>("gc.collections", description: "GC collections");
        _meter.CreateHistogram<double>("gc.pause.duration", unit: "ms", description: "GC pause duration");
        _meter.CreateObservableGauge<long>("memory.usage", description: "Memory usage");
    }
}
`
	metrics, err := parseContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 3 {
		t.Fatalf("expected 3 metrics, got %d", len(metrics))
	}

	names := map[string]bool{}
	for _, m := range metrics {
		names[m.Name] = true
	}

	expected := []string{"gc.collections", "gc.pause.duration", "memory.usage"}
	for _, name := range expected {
		if !names[name] {
			t.Errorf("expected metric '%s' not found", name)
		}
	}
}

func TestParseContent_NoMetrics(t *testing.T) {
	content := `
public class SomeClass
{
    public void DoSomething()
    {
        Console.WriteLine("Hello");
    }
}
`
	metrics, err := parseContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 0 {
		t.Errorf("expected 0 metrics, got %d", len(metrics))
	}
}

func TestParseContent_WithoutGenericType(t *testing.T) {
	content := `
meter.CreateCounter("simple.counter", unit: "1", description: "A simple counter");
`
	metrics, err := parseContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	if metrics[0].Name != "simple.counter" {
		t.Errorf("expected name 'simple.counter', got '%s'", metrics[0].Name)
	}
}

func TestParseContent_ObservableUpDownCounter(t *testing.T) {
	content := `
_meter.CreateObservableUpDownCounter<long>(
    "process.runtime.dotnet.thread_pool.threads.count",
    () => GetThreadPoolCount(),
    description: "Thread pool thread count");
`
	metrics, err := parseContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	m := metrics[0]
	if m.InstrumentType != "updowncounter" {
		t.Errorf("expected type 'updowncounter', got '%s'", m.InstrumentType)
	}
}

func TestExtractConstants(t *testing.T) {
	content := `
public const string MetricName = "my.metric.name";
private static readonly string OtherMetric = "other.metric";
`
	constants := extractConstants(content)

	if constants["MetricName"] != "my.metric.name" {
		t.Errorf("expected 'my.metric.name', got '%s'", constants["MetricName"])
	}
	if constants["OtherMetric"] != "other.metric" {
		t.Errorf("expected 'other.metric', got '%s'", constants["OtherMetric"])
	}
}
