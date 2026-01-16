package golang

import (
	"testing"
)

func TestParseContent_Int64ObservableCounter(t *testing.T) {
	content := `
uptime, err := r.meter.Int64ObservableCounter(
	"runtime.uptime",
	metric.WithUnit("ms"),
	metric.WithDescription("Milliseconds since application was initialized"),
)
`
	metrics, err := parseContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	m := metrics[0]
	if m.Name != "runtime.uptime" {
		t.Errorf("expected name 'runtime.uptime', got '%s'", m.Name)
	}
	if m.InstrumentType != "counter" {
		t.Errorf("expected type 'counter', got '%s'", m.InstrumentType)
	}
	if m.Unit != "ms" {
		t.Errorf("expected unit 'ms', got '%s'", m.Unit)
	}
	if m.Description != "Milliseconds since application was initialized" {
		t.Errorf("expected description 'Milliseconds since application was initialized', got '%s'", m.Description)
	}
}

func TestParseContent_Int64ObservableUpDownCounter(t *testing.T) {
	content := `
goroutines, err := r.meter.Int64ObservableUpDownCounter(
	"process.runtime.go.goroutines",
	metric.WithDescription("Number of goroutines that currently exist"),
)
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
	if m.Name != "process.runtime.go.goroutines" {
		t.Errorf("expected name 'process.runtime.go.goroutines', got '%s'", m.Name)
	}
}

func TestParseContent_Int64Histogram(t *testing.T) {
	content := `
gcPauseNs, err := r.meter.Int64Histogram(
	"process.runtime.go.gc.pause_ns",
	metric.WithUnit("ns"),
	metric.WithDescription("GC pause time in nanoseconds"),
)
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
}

func TestParseContent_Float64Histogram(t *testing.T) {
	content := `
duration, err := meter.Float64Histogram(
	"http.server.request.duration",
	metric.WithUnit("s"),
	metric.WithDescription("Duration of HTTP server requests"),
)
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

func TestParseContent_MultipleMetrics(t *testing.T) {
	content := `
func (r *runtime) register() error {
	uptime, err := r.meter.Int64ObservableCounter(
		"runtime.uptime",
		metric.WithUnit("ms"),
		metric.WithDescription("Milliseconds since application was initialized"),
	)
	if err != nil {
		return err
	}

	goroutines, err := r.meter.Int64ObservableUpDownCounter(
		"process.runtime.go.goroutines",
		metric.WithDescription("Number of goroutines that currently exist"),
	)
	if err != nil {
		return err
	}

	cgoCalls, err := r.meter.Int64ObservableUpDownCounter(
		"process.runtime.go.cgo.calls",
		metric.WithDescription("Number of cgo calls made by the current process"),
	)
	if err != nil {
		return err
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

	expected := []string{"runtime.uptime", "process.runtime.go.goroutines", "process.runtime.go.cgo.calls"}
	for _, name := range expected {
		if !names[name] {
			t.Errorf("expected metric '%s' not found", name)
		}
	}
}

func TestParseContent_NoMetrics(t *testing.T) {
	content := `
package main

func main() {
    fmt.Println("Hello, World!")
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

func TestParseContent_Float64Counter(t *testing.T) {
	content := `
bytesCounter, err := meter.Float64Counter(
	"http.server.request.body.size",
	metric.WithUnit("By"),
	metric.WithDescription("Size of HTTP request body"),
)
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
}
