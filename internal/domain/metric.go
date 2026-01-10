package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

type Attribute struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Required    bool     `json:"required"`
	Enum        []string `json:"enum,omitempty"`
}

type SemconvMatch string

const (
	SemconvMatchExact  SemconvMatch = "exact"
	SemconvMatchPrefix SemconvMatch = "prefix"
	SemconvMatchNone   SemconvMatch = "none"
)

type CanonicalMetric struct {
	ID               string           `json:"id"`
	MetricName       string           `json:"metric_name"`
	InstrumentType   InstrumentType   `json:"instrument_type"`
	Description      string           `json:"description"`
	Unit             string           `json:"unit"`
	Attributes       []Attribute      `json:"attributes"`
	EnabledByDefault bool             `json:"enabled_by_default"`
	ComponentType    ComponentType    `json:"component_type"`
	ComponentName    string           `json:"component_name"`
	SourceCategory   SourceCategory   `json:"source_category"`
	SourceName       string           `json:"source_name"`
	SourceLocation   string           `json:"source_location"`
	ExtractionMethod ExtractionMethod `json:"extraction_method"`
	SourceConfidence ConfidenceLevel  `json:"source_confidence"`
	Repo             string           `json:"repo"`
	Path             string           `json:"path"`
	Commit           string           `json:"commit"`
	ExtractedAt      time.Time        `json:"extracted_at"`

	// Semantic conventions enrichment
	SemconvMatch     SemconvMatch `json:"semconv_match,omitempty"`
	SemconvName      string       `json:"semconv_name,omitempty"`
	SemconvStability string       `json:"semconv_stability,omitempty"`
}

var (
	ErrEmptyMetricName    = errors.New("metric name is required")
	ErrEmptyComponentName = errors.New("component name is required")
	ErrEmptySourceName    = errors.New("source name is required")
	ErrInvalidInstrument  = errors.New("invalid instrument type")
	ErrInvalidComponent   = errors.New("invalid component type")
	ErrInvalidSource      = errors.New("invalid source category")
	ErrInvalidExtraction  = errors.New("invalid extraction method")
	ErrInvalidConfidence  = errors.New("invalid confidence level")
)

func (m *CanonicalMetric) Validate() error {
	if m.MetricName == "" {
		return ErrEmptyMetricName
	}
	if m.ComponentName == "" {
		return ErrEmptyComponentName
	}
	if m.SourceName == "" {
		return ErrEmptySourceName
	}
	if !m.InstrumentType.IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidInstrument, m.InstrumentType)
	}
	if !m.ComponentType.IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidComponent, m.ComponentType)
	}
	if !m.SourceCategory.IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidSource, m.SourceCategory)
	}
	if !m.ExtractionMethod.IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidExtraction, m.ExtractionMethod)
	}
	if !m.SourceConfidence.IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidConfidence, m.SourceConfidence)
	}
	return nil
}

func (m *CanonicalMetric) GenerateID() string {
	data := fmt.Sprintf("%s:%s:%s:%s",
		m.SourceCategory,
		m.SourceName,
		m.ComponentName,
		m.MetricName,
	)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
}

func (m *CanonicalMetric) EnsureID() {
	if m.ID == "" {
		m.ID = m.GenerateID()
	}
}
