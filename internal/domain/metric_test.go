package domain

import (
	"errors"
	"testing"
	"time"
)

func validMetric() *CanonicalMetric {
	return &CanonicalMetric{
		MetricName:       "system.cpu.utilization",
		InstrumentType:   InstrumentGauge,
		Description:      "CPU utilization",
		Unit:             "1",
		ComponentType:    ComponentReceiver,
		ComponentName:    "hostmetrics",
		SourceCategory:   SourceOTEL,
		SourceName:       "opentelemetry-collector-contrib",
		ExtractionMethod: ExtractionMetadata,
		SourceConfidence: ConfidenceAuthoritative,
		ExtractedAt:      time.Now(),
	}
}

func TestCanonicalMetric_Validate(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(*CanonicalMetric)
		wantErr error
	}{
		{
			name:    "valid metric",
			modify:  func(m *CanonicalMetric) {},
			wantErr: nil,
		},
		{
			name:    "empty metric name",
			modify:  func(m *CanonicalMetric) { m.MetricName = "" },
			wantErr: ErrEmptyMetricName,
		},
		{
			name:    "empty component name",
			modify:  func(m *CanonicalMetric) { m.ComponentName = "" },
			wantErr: ErrEmptyComponentName,
		},
		{
			name:    "empty source name",
			modify:  func(m *CanonicalMetric) { m.SourceName = "" },
			wantErr: ErrEmptySourceName,
		},
		{
			name:    "invalid instrument type",
			modify:  func(m *CanonicalMetric) { m.InstrumentType = "invalid" },
			wantErr: ErrInvalidInstrument,
		},
		{
			name:    "invalid component type",
			modify:  func(m *CanonicalMetric) { m.ComponentType = "invalid" },
			wantErr: ErrInvalidComponent,
		},
		{
			name:    "invalid source category",
			modify:  func(m *CanonicalMetric) { m.SourceCategory = "invalid" },
			wantErr: ErrInvalidSource,
		},
		{
			name:    "invalid extraction method",
			modify:  func(m *CanonicalMetric) { m.ExtractionMethod = "invalid" },
			wantErr: ErrInvalidExtraction,
		},
		{
			name:    "invalid confidence level",
			modify:  func(m *CanonicalMetric) { m.SourceConfidence = "invalid" },
			wantErr: ErrInvalidConfidence,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := validMetric()
			tt.modify(m)
			err := m.Validate()

			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
				return
			}

			if err == nil {
				t.Errorf("Validate() expected error %v, got nil", tt.wantErr)
				return
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestCanonicalMetric_GenerateID(t *testing.T) {
	m := validMetric()
	id := m.GenerateID()

	if id == "" {
		t.Error("GenerateID() returned empty string")
	}

	if len(id) != 32 {
		t.Errorf("GenerateID() returned ID of length %d, want 32", len(id))
	}

	id2 := m.GenerateID()
	if id != id2 {
		t.Error("GenerateID() not deterministic")
	}
}

func TestCanonicalMetric_GenerateID_Deterministic(t *testing.T) {
	m1 := validMetric()
	m2 := validMetric()

	if m1.GenerateID() != m2.GenerateID() {
		t.Error("GenerateID() should produce same ID for identical metrics")
	}

	m2.MetricName = "different.metric"
	if m1.GenerateID() == m2.GenerateID() {
		t.Error("GenerateID() should produce different IDs for different metric names")
	}
}

func TestCanonicalMetric_GenerateID_Components(t *testing.T) {
	base := validMetric()
	baseID := base.GenerateID()

	tests := []struct {
		name   string
		modify func(*CanonicalMetric)
	}{
		{"different source category", func(m *CanonicalMetric) { m.SourceCategory = SourcePrometheus }},
		{"different source name", func(m *CanonicalMetric) { m.SourceName = "different-source" }},
		{"different component name", func(m *CanonicalMetric) { m.ComponentName = "different-component" }},
		{"different metric name", func(m *CanonicalMetric) { m.MetricName = "different.metric" }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := validMetric()
			tt.modify(m)
			if m.GenerateID() == baseID {
				t.Errorf("GenerateID() should produce different ID for %s", tt.name)
			}
		})
	}
}

func TestCanonicalMetric_EnsureID(t *testing.T) {
	m := validMetric()
	m.ID = ""
	m.EnsureID()

	if m.ID == "" {
		t.Error("EnsureID() did not set ID")
	}

	originalID := m.ID
	m.EnsureID()
	if m.ID != originalID {
		t.Error("EnsureID() should not change existing ID")
	}
}
