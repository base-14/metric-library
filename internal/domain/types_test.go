package domain

import "testing"

func TestInstrumentType_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		t     InstrumentType
		valid bool
	}{
		{"counter", InstrumentCounter, true},
		{"updowncounter", InstrumentUpDownCounter, true},
		{"gauge", InstrumentGauge, true},
		{"histogram", InstrumentHistogram, true},
		{"summary", InstrumentSummary, true},
		{"invalid", InstrumentType("invalid"), false},
		{"empty", InstrumentType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.IsValid(); got != tt.valid {
				t.Errorf("InstrumentType(%q).IsValid() = %v, want %v", tt.t, got, tt.valid)
			}
		})
	}
}

func TestComponentType_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		t     ComponentType
		valid bool
	}{
		{"receiver", ComponentReceiver, true},
		{"exporter", ComponentExporter, true},
		{"processor", ComponentProcessor, true},
		{"extension", ComponentExtension, true},
		{"connector", ComponentConnector, true},
		{"instrumentation", ComponentInstrumentation, true},
		{"platform", ComponentPlatform, true},
		{"invalid", ComponentType("invalid"), false},
		{"empty", ComponentType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.IsValid(); got != tt.valid {
				t.Errorf("ComponentType(%q).IsValid() = %v, want %v", tt.t, got, tt.valid)
			}
		})
	}
}

func TestSourceCategory_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		c     SourceCategory
		valid bool
	}{
		{"otel", SourceOTEL, true},
		{"prometheus", SourcePrometheus, true},
		{"kubernetes", SourceKubernetes, true},
		{"cloud", SourceCloud, true},
		{"vendor", SourceVendor, true},
		{"invalid", SourceCategory("invalid"), false},
		{"empty", SourceCategory(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValid(); got != tt.valid {
				t.Errorf("SourceCategory(%q).IsValid() = %v, want %v", tt.c, got, tt.valid)
			}
		})
	}
}

func TestExtractionMethod_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		m     ExtractionMethod
		valid bool
	}{
		{"metadata", ExtractionMetadata, true},
		{"ast", ExtractionAST, true},
		{"scrape", ExtractionScrape, true},
		{"hybrid", ExtractionHybrid, true},
		{"invalid", ExtractionMethod("invalid"), false},
		{"empty", ExtractionMethod(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsValid(); got != tt.valid {
				t.Errorf("ExtractionMethod(%q).IsValid() = %v, want %v", tt.m, got, tt.valid)
			}
		})
	}
}

func TestConfidenceLevel_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		c     ConfidenceLevel
		valid bool
	}{
		{"authoritative", ConfidenceAuthoritative, true},
		{"derived", ConfidenceDerived, true},
		{"documented", ConfidenceDocumented, true},
		{"vendor_claimed", ConfidenceVendorClaimed, true},
		{"invalid", ConfidenceLevel("invalid"), false},
		{"empty", ConfidenceLevel(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValid(); got != tt.valid {
				t.Errorf("ConfidenceLevel(%q).IsValid() = %v, want %v", tt.c, got, tt.valid)
			}
		})
	}
}
