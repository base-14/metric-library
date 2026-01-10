package extractor

import (
	"testing"

	"github.com/base-14/metric-library/internal/domain"
	"github.com/base-14/metric-library/internal/parser"
)

func TestMetricExtractor_Extract(t *testing.T) {
	meta := &parser.Metadata{
		Type: "mysqlreceiver",
		Status: parser.StatusDefinition{
			Class: "receiver",
		},
		Attributes: map[string]parser.AttributeDefinition{
			"buffer_pool_data": {
				Description: "The status of buffer pool data.",
				Type:        "string",
				Enum:        []string{"dirty", "clean"},
			},
		},
		Metrics: map[string]parser.MetricDefinition{
			"mysql.buffer_pool.pages": {
				Enabled:     true,
				Description: "The number of pages in the InnoDB buffer pool.",
				Unit:        "{pages}",
				Sum: &parser.SumDefinition{
					ValueType:              "int",
					Monotonic:              false,
					AggregationTemporality: "cumulative",
				},
				Attributes: []string{"buffer_pool_data"},
			},
		},
	}

	extractor := NewMetricExtractor(
		"opentelemetry-collector-contrib",
		"mysql",
		"receiver",
	)

	metrics, err := extractor.Extract(meta)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	m := metrics[0]
	if m.MetricName != "mysql.buffer_pool.pages" {
		t.Errorf("expected metric name 'mysql.buffer_pool.pages', got %q", m.MetricName)
	}

	if m.Description != "The number of pages in the InnoDB buffer pool." {
		t.Errorf("unexpected description: %q", m.Description)
	}

	if m.Unit != "{pages}" {
		t.Errorf("unexpected unit: %q", m.Unit)
	}

	if m.InstrumentType != domain.InstrumentUpDownCounter {
		t.Errorf("expected instrument type 'updowncounter', got %q", m.InstrumentType)
	}

	if m.ComponentType != domain.ComponentReceiver {
		t.Errorf("expected component type 'receiver', got %q", m.ComponentType)
	}

	if m.ComponentName != "mysql" {
		t.Errorf("expected component name 'mysql', got %q", m.ComponentName)
	}

	if m.SourceName != "opentelemetry-collector-contrib" {
		t.Errorf("expected source name 'opentelemetry-collector-contrib', got %q", m.SourceName)
	}

	if m.SourceCategory != domain.SourceOTEL {
		t.Errorf("expected source category 'otel', got %q", m.SourceCategory)
	}

	if m.ExtractionMethod != domain.ExtractionMetadata {
		t.Errorf("expected extraction method 'metadata', got %q", m.ExtractionMethod)
	}

	if m.SourceConfidence != domain.ConfidenceAuthoritative {
		t.Errorf("expected confidence 'authoritative', got %q", m.SourceConfidence)
	}

	if len(m.Attributes) != 1 {
		t.Errorf("expected 1 attribute, got %d", len(m.Attributes))
	}
}

func TestMetricExtractor_Extract_Counter(t *testing.T) {
	meta := &parser.Metadata{
		Type: "test",
		Metrics: map[string]parser.MetricDefinition{
			"test.counter": {
				Enabled:     true,
				Description: "A counter metric.",
				Unit:        "1",
				Sum: &parser.SumDefinition{
					ValueType: "int",
					Monotonic: true,
				},
			},
		},
	}

	extractor := NewMetricExtractor("test-source", "test", "receiver")
	metrics, err := extractor.Extract(meta)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	if metrics[0].InstrumentType != domain.InstrumentCounter {
		t.Errorf("expected 'counter', got %q", metrics[0].InstrumentType)
	}
}

func TestMetricExtractor_Extract_Gauge(t *testing.T) {
	meta := &parser.Metadata{
		Type: "test",
		Metrics: map[string]parser.MetricDefinition{
			"test.gauge": {
				Enabled:     true,
				Description: "A gauge metric.",
				Unit:        "1",
				Gauge: &parser.GaugeDefinition{
					ValueType: "double",
				},
			},
		},
	}

	extractor := NewMetricExtractor("test-source", "test", "receiver")
	metrics, err := extractor.Extract(meta)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if metrics[0].InstrumentType != domain.InstrumentGauge {
		t.Errorf("expected 'gauge', got %q", metrics[0].InstrumentType)
	}
}

func TestMetricExtractor_Extract_Histogram(t *testing.T) {
	meta := &parser.Metadata{
		Type: "test",
		Metrics: map[string]parser.MetricDefinition{
			"test.histogram": {
				Enabled:     true,
				Description: "A histogram metric.",
				Unit:        "ms",
				Histogram: &parser.HistogramDefinition{
					ValueType: "double",
				},
			},
		},
	}

	extractor := NewMetricExtractor("test-source", "test", "receiver")
	metrics, err := extractor.Extract(meta)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if metrics[0].InstrumentType != domain.InstrumentHistogram {
		t.Errorf("expected 'histogram', got %q", metrics[0].InstrumentType)
	}
}

func TestMetricExtractor_Extract_MultipleMetrics(t *testing.T) {
	meta := &parser.Metadata{
		Type: "test",
		Metrics: map[string]parser.MetricDefinition{
			"test.metric1": {
				Enabled:     true,
				Description: "First metric",
				Unit:        "1",
				Gauge:       &parser.GaugeDefinition{ValueType: "int"},
			},
			"test.metric2": {
				Enabled:     true,
				Description: "Second metric",
				Unit:        "1",
				Gauge:       &parser.GaugeDefinition{ValueType: "int"},
			},
			"test.disabled": {
				Enabled:     false,
				Description: "Disabled metric",
				Unit:        "1",
				Gauge:       &parser.GaugeDefinition{ValueType: "int"},
			},
		},
	}

	extractor := NewMetricExtractor("test-source", "test", "receiver")
	metrics, err := extractor.Extract(meta)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 3 {
		t.Errorf("expected 3 metrics (including disabled), got %d", len(metrics))
	}
}

func TestMetricExtractor_Extract_WithAttributes(t *testing.T) {
	meta := &parser.Metadata{
		Type: "test",
		Attributes: map[string]parser.AttributeDefinition{
			"state": {
				Description: "The state",
				Type:        "string",
				Enum:        []string{"active", "idle"},
			},
			"direction": {
				Description: "The direction",
				Type:        "string",
				Enum:        []string{"read", "write"},
			},
		},
		Metrics: map[string]parser.MetricDefinition{
			"test.metric": {
				Enabled:     true,
				Description: "A metric with attributes",
				Unit:        "1",
				Gauge:       &parser.GaugeDefinition{ValueType: "int"},
				Attributes:  []string{"state", "direction"},
			},
		},
	}

	extractor := NewMetricExtractor("test-source", "test", "receiver")
	metrics, err := extractor.Extract(meta)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	m := metrics[0]
	if len(m.Attributes) != 2 {
		t.Fatalf("expected 2 attributes, got %d", len(m.Attributes))
	}

	attrMap := make(map[string]domain.Attribute)
	for _, a := range m.Attributes {
		attrMap[a.Name] = a
	}

	if attr, ok := attrMap["state"]; ok {
		if attr.Description != "The state" {
			t.Errorf("unexpected state description: %q", attr.Description)
		}
		if attr.Type != "string" {
			t.Errorf("unexpected state type: %q", attr.Type)
		}
		if len(attr.Enum) != 2 {
			t.Errorf("expected 2 possible values, got %d", len(attr.Enum))
		}
	} else {
		t.Error("missing 'state' attribute")
	}
}

func TestMetricExtractor_Extract_EmptyMetadata(t *testing.T) {
	meta := &parser.Metadata{}

	extractor := NewMetricExtractor("test-source", "test", "receiver")
	metrics, err := extractor.Extract(meta)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(metrics) != 0 {
		t.Errorf("expected 0 metrics, got %d", len(metrics))
	}
}

func TestMetricExtractor_Extract_GeneratesID(t *testing.T) {
	meta := &parser.Metadata{
		Type: "test",
		Metrics: map[string]parser.MetricDefinition{
			"test.metric": {
				Enabled:     true,
				Description: "Test metric",
				Unit:        "1",
				Gauge:       &parser.GaugeDefinition{ValueType: "int"},
			},
		},
	}

	extractor := NewMetricExtractor("test-source", "test", "receiver")
	metrics, err := extractor.Extract(meta)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if metrics[0].ID == "" {
		t.Error("expected ID to be generated")
	}
}
