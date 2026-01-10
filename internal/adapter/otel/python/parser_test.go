package python

import (
	"testing"
)

func TestParseContent_CreateHistogram(t *testing.T) {
	content := `
meter = get_meter(__name__, __version__)
duration_histogram = meter.create_histogram(
    name="http.client.request.duration",
    unit="s",
    description="Duration of HTTP client requests.",
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
	if m.Name != "http.client.request.duration" {
		t.Errorf("expected name 'http.client.request.duration', got '%s'", m.Name)
	}
	if m.InstrumentType != "histogram" {
		t.Errorf("expected type 'histogram', got '%s'", m.InstrumentType)
	}
	if m.Unit != "s" {
		t.Errorf("expected unit 's', got '%s'", m.Unit)
	}
	if m.Description != "Duration of HTTP client requests." {
		t.Errorf("expected description 'Duration of HTTP client requests.', got '%s'", m.Description)
	}
}

func TestParseContent_CreateCounter(t *testing.T) {
	content := `
self._meter.create_counter(
    name="http.server.active_requests",
    unit="requests",
    description="Number of active HTTP requests",
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
	if m.Name != "http.server.active_requests" {
		t.Errorf("expected name 'http.server.active_requests', got '%s'", m.Name)
	}
	if m.InstrumentType != "counter" {
		t.Errorf("expected type 'counter', got '%s'", m.InstrumentType)
	}
}

func TestParseContent_ObservableGauge(t *testing.T) {
	content := `
self._meter.create_observable_gauge(
    name="system.memory.usage",
    callbacks=[self._get_memory],
    description="System memory usage",
    unit="By",
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
	if m.Name != "system.memory.usage" {
		t.Errorf("expected name 'system.memory.usage', got '%s'", m.Name)
	}
	if m.InstrumentType != "gauge" {
		t.Errorf("expected type 'gauge', got '%s'", m.InstrumentType)
	}
}

func TestParseContent_ObservableCounter(t *testing.T) {
	content := `
self._meter.create_observable_counter(
    name="system.cpu.time",
    callbacks=[self._get_cpu_time],
    description="System CPU time",
    unit="s",
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
	if m.Name != "system.cpu.time" {
		t.Errorf("expected name 'system.cpu.time', got '%s'", m.Name)
	}
	if m.InstrumentType != "counter" {
		t.Errorf("expected type 'counter', got '%s'", m.InstrumentType)
	}
}

func TestParseContent_MultipleMetrics(t *testing.T) {
	content := `
meter = get_meter(__name__)

request_histogram = meter.create_histogram(
    name="http.request.duration",
    unit="ms",
    description="Request duration",
)

response_histogram = meter.create_histogram(
    name="http.response.size",
    unit="By",
    description="Response size",
)

error_counter = meter.create_counter(
    name="http.errors",
    unit="1",
    description="Error count",
)
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

	expected := []string{"http.request.duration", "http.response.size", "http.errors"}
	for _, name := range expected {
		if !names[name] {
			t.Errorf("expected metric '%s' not found", name)
		}
	}
}

func TestParseContent_UpDownCounter(t *testing.T) {
	content := `
meter.create_up_down_counter(
    name="queue.length",
    description="Queue length",
    unit="items",
)
`
	metrics, err := parseContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	if metrics[0].InstrumentType != "updowncounter" {
		t.Errorf("expected type 'updowncounter', got '%s'", metrics[0].InstrumentType)
	}
}

func TestParseContent_ObservableUpDownCounter(t *testing.T) {
	content := `
self._meter.create_observable_up_down_counter(
    name="system.network.connections",
    callbacks=[self._get_connections],
    description="Network connections",
    unit="connections",
)
`
	metrics, err := parseContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	if metrics[0].InstrumentType != "updowncounter" {
		t.Errorf("expected type 'updowncounter', got '%s'", metrics[0].InstrumentType)
	}
}

func TestParseContent_NoMetrics(t *testing.T) {
	content := `
def some_function():
    pass
`
	metrics, err := parseContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 0 {
		t.Errorf("expected 0 metrics, got %d", len(metrics))
	}
}

func TestParseContent_StringWithParens(t *testing.T) {
	content := `
meter.create_histogram(
    name="db.query.duration",
    description="Duration of database queries (including retries)",
    unit="s",
)
`
	metrics, err := parseContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	if metrics[0].Description != "Duration of database queries (including retries)" {
		t.Errorf("unexpected description: %s", metrics[0].Description)
	}
}

func TestFindMatchingParen(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		start    int
		expected int
	}{
		{
			name:     "simple",
			content:  "(abc)",
			start:    0,
			expected: 5,
		},
		{
			name:     "nested",
			content:  "(a(b)c)",
			start:    0,
			expected: 7,
		},
		{
			name:     "string with paren",
			content:  `("hello)")`,
			start:    0,
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findMatchingParen(tt.content, tt.start)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
