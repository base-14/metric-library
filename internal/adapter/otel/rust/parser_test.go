package rust

import (
	"testing"
)

func TestParseContent_F64Histogram(t *testing.T) {
	content := `
let http_server_duration = meter
    .f64_histogram("http.server.duration")
    .with_description("Measures the duration of inbound HTTP requests.")
    .with_unit("s")
    .build();
`
	metrics, err := parseContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	m := metrics[0]
	if m.Name != "http.server.duration" {
		t.Errorf("expected name 'http.server.duration', got '%s'", m.Name)
	}
	if m.InstrumentType != "histogram" {
		t.Errorf("expected type 'histogram', got '%s'", m.InstrumentType)
	}
	if m.Unit != "s" {
		t.Errorf("expected unit 's', got '%s'", m.Unit)
	}
	if m.Description != "Measures the duration of inbound HTTP requests." {
		t.Errorf("unexpected description: '%s'", m.Description)
	}
}

func TestParseContent_I64UpDownCounter(t *testing.T) {
	content := `
let active_requests = meter
    .i64_up_down_counter("http.server.active_requests")
    .with_description("Measures the number of concurrent HTTP requests.")
    .build();
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

func TestParseContent_U64Counter(t *testing.T) {
	content := `
let counter = meter
    .u64_counter("requests.total")
    .with_unit("{request}")
    .build();
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
	if m.Name != "requests.total" {
		t.Errorf("expected name 'requests.total', got '%s'", m.Name)
	}
}

func TestParseContent_WithConstant(t *testing.T) {
	content := `
const HTTP_SERVER_DURATION: &str = "http.server.duration";

let duration = meter
    .f64_histogram(HTTP_SERVER_DURATION)
    .with_description("Duration of HTTP requests.")
    .with_unit("s")
    .build();
`
	metrics, err := parseContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	m := metrics[0]
	if m.Name != "http.server.duration" {
		t.Errorf("expected name 'http.server.duration', got '%s'", m.Name)
	}
}

func TestParseContent_ObservableGauge(t *testing.T) {
	content := `
let gauge = meter
    .f64_observable_gauge("cpu.usage")
    .with_description("CPU usage percentage")
    .with_unit("%")
    .build();
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
}

func TestParseContent_MultipleMetrics(t *testing.T) {
	content := `
const HTTP_SERVER_DURATION: &str = "http.server.duration";
const HTTP_SERVER_ACTIVE_REQUESTS: &str = "http.server.active_requests";

let duration = meter
    .f64_histogram(HTTP_SERVER_DURATION)
    .with_description("Duration of HTTP requests.")
    .with_unit("s")
    .build();

let active = meter
    .i64_up_down_counter(HTTP_SERVER_ACTIVE_REQUESTS)
    .with_description("Active requests.")
    .build();
`
	metrics, err := parseContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 2 {
		t.Fatalf("expected 2 metrics, got %d", len(metrics))
	}

	names := map[string]bool{}
	for _, m := range metrics {
		names[m.Name] = true
	}

	if !names["http.server.duration"] {
		t.Error("expected metric 'http.server.duration' not found")
	}
	if !names["http.server.active_requests"] {
		t.Error("expected metric 'http.server.active_requests' not found")
	}
}

func TestParseContent_NoMetrics(t *testing.T) {
	content := `
fn main() {
    println!("Hello, World!");
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

func TestExtractConstants(t *testing.T) {
	content := `
const HTTP_SERVER_DURATION: &str = "http.server.duration";
const HTTP_SERVER_ACTIVE_REQUESTS: &str = "http.server.active_requests";
`
	constants := extractConstants(content)

	if constants["HTTP_SERVER_DURATION"] != "http.server.duration" {
		t.Errorf("expected 'http.server.duration', got '%s'", constants["HTTP_SERVER_DURATION"])
	}
	if constants["HTTP_SERVER_ACTIVE_REQUESTS"] != "http.server.active_requests" {
		t.Errorf("expected 'http.server.active_requests', got '%s'", constants["HTTP_SERVER_ACTIVE_REQUESTS"])
	}
}
