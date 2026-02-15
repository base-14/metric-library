package domain

type InstrumentType string

const (
	InstrumentCounter       InstrumentType = "counter"
	InstrumentUpDownCounter InstrumentType = "updowncounter"
	InstrumentGauge         InstrumentType = "gauge"
	InstrumentHistogram     InstrumentType = "histogram"
	InstrumentSummary       InstrumentType = "summary"
)

func (t InstrumentType) IsValid() bool {
	switch t {
	case InstrumentCounter, InstrumentUpDownCounter, InstrumentGauge, InstrumentHistogram, InstrumentSummary:
		return true
	}
	return false
}

type ComponentType string

const (
	ComponentReceiver        ComponentType = "receiver"
	ComponentExporter        ComponentType = "exporter"
	ComponentProcessor       ComponentType = "processor"
	ComponentExtension       ComponentType = "extension"
	ComponentConnector       ComponentType = "connector"
	ComponentInstrumentation ComponentType = "instrumentation"
	ComponentPlatform        ComponentType = "platform"
)

func (t ComponentType) IsValid() bool {
	switch t {
	case ComponentReceiver, ComponentExporter, ComponentProcessor, ComponentExtension, ComponentConnector, ComponentInstrumentation, ComponentPlatform:
		return true
	}
	return false
}

type SourceCategory string

const (
	SourceOTEL        SourceCategory = "otel"
	SourcePrometheus  SourceCategory = "prometheus"
	SourceKubernetes  SourceCategory = "kubernetes"
	SourceCloud       SourceCategory = "cloud"
	SourceVendor      SourceCategory = "vendor"
	SourceCodingAgent SourceCategory = "codingagent"
)

func (c SourceCategory) IsValid() bool {
	switch c {
	case SourceOTEL, SourcePrometheus, SourceKubernetes, SourceCloud, SourceVendor, SourceCodingAgent:
		return true
	}
	return false
}

type ExtractionMethod string

const (
	ExtractionMetadata ExtractionMethod = "metadata"
	ExtractionAST      ExtractionMethod = "ast"
	ExtractionScrape   ExtractionMethod = "scrape"
	ExtractionHybrid   ExtractionMethod = "hybrid"
)

func (m ExtractionMethod) IsValid() bool {
	switch m {
	case ExtractionMetadata, ExtractionAST, ExtractionScrape, ExtractionHybrid:
		return true
	}
	return false
}

type ConfidenceLevel string

const (
	ConfidenceAuthoritative ConfidenceLevel = "authoritative"
	ConfidenceDerived       ConfidenceLevel = "derived"
	ConfidenceDocumented    ConfidenceLevel = "documented"
	ConfidenceVendorClaimed ConfidenceLevel = "vendor_claimed"
)

func (c ConfidenceLevel) IsValid() bool {
	switch c {
	case ConfidenceAuthoritative, ConfidenceDerived, ConfidenceDocumented, ConfidenceVendorClaimed:
		return true
	}
	return false
}
