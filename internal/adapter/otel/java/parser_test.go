package java

import (
	"testing"
)

func TestParseContent_CounterBuilder(t *testing.T) {
	content := `
meter.counterBuilder("http.server.requests")
    .setDescription("Number of HTTP requests")
    .setUnit("{request}")
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
	if m.Name != "http.server.requests" {
		t.Errorf("expected name 'http.server.requests', got '%s'", m.Name)
	}
	if m.InstrumentType != "counter" {
		t.Errorf("expected type 'counter', got '%s'", m.InstrumentType)
	}
	if m.Unit != "{request}" {
		t.Errorf("expected unit '{request}', got '%s'", m.Unit)
	}
	if m.Description != "Number of HTTP requests" {
		t.Errorf("expected description 'Number of HTTP requests', got '%s'", m.Description)
	}
}

func TestParseContent_HistogramBuilder(t *testing.T) {
	content := `
meter.histogramBuilder("http.server.request.duration")
    .setDescription("Duration of HTTP requests")
    .setUnit("ms")
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
	if m.InstrumentType != "histogram" {
		t.Errorf("expected type 'histogram', got '%s'", m.InstrumentType)
	}
}

func TestParseContent_GaugeBuilder(t *testing.T) {
	content := `
meter.gaugeBuilder("jvm.memory.used")
    .setDescription("JVM memory used")
    .setUnit("By")
    .buildWithCallback(callback);
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
	if m.Name != "jvm.memory.used" {
		t.Errorf("expected name 'jvm.memory.used', got '%s'", m.Name)
	}
}

func TestParseContent_UpDownCounterBuilder(t *testing.T) {
	content := `
meter.upDownCounterBuilder("db.client.connections.idle")
    .setDescription("Idle database connections")
    .setUnit("{connection}")
    .buildObserver();
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

func TestParseContent_MultipleMetrics(t *testing.T) {
	content := `
public class Metrics {
    private void register() {
        meter.counterBuilder("requests.total")
            .setDescription("Total requests")
            .setUnit("1")
            .build();

        meter.histogramBuilder("request.duration")
            .setDescription("Request duration")
            .setUnit("ms")
            .build();

        meter.gaugeBuilder("active.connections")
            .setDescription("Active connections")
            .setUnit("{connection}")
            .buildWithCallback(obs -> {});
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

	expected := []string{"requests.total", "request.duration", "active.connections"}
	for _, name := range expected {
		if !names[name] {
			t.Errorf("expected metric '%s' not found", name)
		}
	}
}

func TestParseContent_NoMetrics(t *testing.T) {
	content := `
public class SomeClass {
    public void doSomething() {
        System.out.println("Hello");
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

func TestParseContent_MultilineChain(t *testing.T) {
	content := `
meter
    .counterBuilder("multi.line.metric")
    .setDescription("A metric defined across multiple lines")
    .setUnit("1")
    .build();
`
	metrics, err := parseContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	if metrics[0].Description != "A metric defined across multiple lines" {
		t.Errorf("unexpected description: %s", metrics[0].Description)
	}
}

func TestParseContent_WithoutDescription(t *testing.T) {
	content := `
meter.counterBuilder("simple.counter")
    .setUnit("1")
    .build();
`
	metrics, err := parseContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	if metrics[0].Description != "" {
		t.Errorf("expected empty description, got '%s'", metrics[0].Description)
	}
	if metrics[0].Unit != "1" {
		t.Errorf("expected unit '1', got '%s'", metrics[0].Unit)
	}
}
