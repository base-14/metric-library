package parser

import (
	"testing"
)

func TestMetadataParser_Parse(t *testing.T) {
	content := []byte(`
type: mysqlreceiver

status:
  class: receiver
  stability:
    beta: [metrics]
  codeowners:
    active: [antonblock]

attributes:
  buffer_pool_data:
    description: The status of buffer pool data.
    type: string
    enum: [dirty, clean]

metrics:
  mysql.buffer_pool.pages:
    enabled: true
    description: The number of pages in the InnoDB buffer pool.
    unit: "{pages}"
    sum:
      value_type: int
      monotonic: false
      aggregation_temporality: cumulative
    attributes: [buffer_pool_data]

  mysql.commands:
    enabled: false
    description: The number of times each type of command has been executed.
    unit: "{commands}"
    sum:
      value_type: int
      monotonic: true
      aggregation_temporality: cumulative
`)

	parser := NewMetadataParser()
	meta, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if meta.Type != "mysqlreceiver" {
		t.Errorf("expected type 'mysqlreceiver', got %q", meta.Type)
	}

	if meta.Status.Class != "receiver" {
		t.Errorf("expected class 'receiver', got %q", meta.Status.Class)
	}

	if len(meta.Metrics) != 2 {
		t.Errorf("expected 2 metrics, got %d", len(meta.Metrics))
	}

	if m, ok := meta.Metrics["mysql.buffer_pool.pages"]; ok {
		if !m.Enabled {
			t.Error("expected mysql.buffer_pool.pages to be enabled")
		}
		if m.Description != "The number of pages in the InnoDB buffer pool." {
			t.Errorf("unexpected description: %q", m.Description)
		}
		if m.Unit != "{pages}" {
			t.Errorf("unexpected unit: %q", m.Unit)
		}
		if m.Sum == nil {
			t.Error("expected sum to be set")
		} else {
			if m.Sum.ValueType != "int" {
				t.Errorf("unexpected value_type: %q", m.Sum.ValueType)
			}
			if m.Sum.Monotonic {
				t.Error("expected monotonic to be false")
			}
		}
		if len(m.Attributes) != 1 || m.Attributes[0] != "buffer_pool_data" {
			t.Errorf("unexpected attributes: %v", m.Attributes)
		}
	} else {
		t.Error("expected metric mysql.buffer_pool.pages not found")
	}
}

func TestMetadataParser_Parse_Gauge(t *testing.T) {
	content := []byte(`
type: hostmetricsreceiver

status:
  class: receiver

metrics:
  system.cpu.utilization:
    enabled: true
    description: The CPU utilization percentage.
    unit: "1"
    gauge:
      value_type: double
`)

	parser := NewMetadataParser()
	meta, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	m, ok := meta.Metrics["system.cpu.utilization"]
	if !ok {
		t.Fatal("expected metric system.cpu.utilization not found")
	}

	if m.Gauge == nil {
		t.Fatal("expected gauge to be set")
	}

	if m.Gauge.ValueType != "double" {
		t.Errorf("expected value_type 'double', got %q", m.Gauge.ValueType)
	}
}

func TestMetadataParser_Parse_Attributes(t *testing.T) {
	content := []byte(`
type: test

status:
  class: receiver

attributes:
  state:
    description: The state of something.
    type: string
    enum: [idle, busy, waiting]
  direction:
    description: The direction.
    type: string
    enum: [read, write]
`)

	parser := NewMetadataParser()
	meta, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(meta.Attributes) != 2 {
		t.Errorf("expected 2 attributes, got %d", len(meta.Attributes))
	}

	if attr, ok := meta.Attributes["state"]; ok {
		if attr.Description != "The state of something." {
			t.Errorf("unexpected description: %q", attr.Description)
		}
		if attr.Type != "string" {
			t.Errorf("unexpected type: %q", attr.Type)
		}
		if len(attr.Enum) != 3 {
			t.Errorf("expected 3 enum values, got %d", len(attr.Enum))
		}
	} else {
		t.Error("expected attribute 'state' not found")
	}
}

func TestMetadataParser_Parse_InvalidYAML(t *testing.T) {
	content := []byte(`
type: invalid
  indentation: wrong
`)

	parser := NewMetadataParser()
	_, err := parser.Parse(content)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestMetadataParser_Parse_EmptyContent(t *testing.T) {
	parser := NewMetadataParser()
	meta, err := parser.Parse([]byte(""))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if meta.Type != "" {
		t.Errorf("expected empty type, got %q", meta.Type)
	}
}

func TestMetadataParser_ParseFile(t *testing.T) {
	parser := NewMetadataParser()
	_, err := parser.ParseFile("/nonexistent/path/metadata.yaml")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestMetadataParser_GetInstrumentType(t *testing.T) {
	tests := []struct {
		name     string
		metric   MetricDefinition
		expected string
	}{
		{
			name:     "sum_monotonic",
			metric:   MetricDefinition{Sum: &SumDefinition{Monotonic: true}},
			expected: "counter",
		},
		{
			name:     "sum_non_monotonic",
			metric:   MetricDefinition{Sum: &SumDefinition{Monotonic: false}},
			expected: "updowncounter",
		},
		{
			name:     "gauge",
			metric:   MetricDefinition{Gauge: &GaugeDefinition{ValueType: "double"}},
			expected: "gauge",
		},
		{
			name:     "histogram",
			metric:   MetricDefinition{Histogram: &HistogramDefinition{}},
			expected: "histogram",
		},
		{
			name:     "unknown",
			metric:   MetricDefinition{},
			expected: "gauge",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.metric.InstrumentType()
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}
