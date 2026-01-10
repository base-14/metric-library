package parser

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Metadata struct {
	Type       string                         `yaml:"type"`
	Status     StatusDefinition               `yaml:"status"`
	Attributes map[string]AttributeDefinition `yaml:"attributes"`
	Metrics    map[string]MetricDefinition    `yaml:"metrics"`
}

type StatusDefinition struct {
	Class      string               `yaml:"class"`
	Stability  StabilityDefinition  `yaml:"stability"`
	Codeowners CodeownersDefinition `yaml:"codeowners"`
}

type StabilityDefinition struct {
	Beta        []string `yaml:"beta"`
	Alpha       []string `yaml:"alpha"`
	Development []string `yaml:"development"`
	Stable      []string `yaml:"stable"`
}

type CodeownersDefinition struct {
	Active   []string `yaml:"active"`
	Emeritus []string `yaml:"emeritus"`
}

type AttributeDefinition struct {
	Description string   `yaml:"description"`
	Type        string   `yaml:"type"`
	Enum        []string `yaml:"enum"`
}

type MetricDefinition struct {
	Enabled     bool                 `yaml:"enabled"`
	Description string               `yaml:"description"`
	Unit        string               `yaml:"unit"`
	Sum         *SumDefinition       `yaml:"sum"`
	Gauge       *GaugeDefinition     `yaml:"gauge"`
	Histogram   *HistogramDefinition `yaml:"histogram"`
	Attributes  []string             `yaml:"attributes"`
	Warnings    WarningsDefinition   `yaml:"warnings"`
}

type SumDefinition struct {
	ValueType              string `yaml:"value_type"`
	Monotonic              bool   `yaml:"monotonic"`
	AggregationTemporality string `yaml:"aggregation_temporality"`
}

type GaugeDefinition struct {
	ValueType string `yaml:"value_type"`
}

type HistogramDefinition struct {
	ValueType string `yaml:"value_type"`
}

type WarningsDefinition struct {
	IfEnabled          string `yaml:"if_enabled"`
	IfEnabledNotSet    string `yaml:"if_enabled_not_set"`
	IfConfigured       string `yaml:"if_configured"`
	IfConfiguredNotSet string `yaml:"if_configured_not_set"`
}

func (m MetricDefinition) InstrumentType() string {
	if m.Sum != nil {
		if m.Sum.Monotonic {
			return "counter"
		}
		return "updowncounter"
	}
	if m.Gauge != nil {
		return "gauge"
	}
	if m.Histogram != nil {
		return "histogram"
	}
	return "gauge"
}

func (m MetricDefinition) ValueType() string {
	if m.Sum != nil {
		return m.Sum.ValueType
	}
	if m.Gauge != nil {
		return m.Gauge.ValueType
	}
	if m.Histogram != nil {
		return m.Histogram.ValueType
	}
	return ""
}

type MetadataParser struct{}

func NewMetadataParser() *MetadataParser {
	return &MetadataParser{}
}

func (p *MetadataParser) Parse(content []byte) (*Metadata, error) {
	var meta Metadata
	if len(content) == 0 {
		return &meta, nil
	}

	if err := yaml.Unmarshal(content, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &meta, nil
}

func (p *MetadataParser) ParseFile(path string) (*Metadata, error) {
	content, err := os.ReadFile(path) //nolint:gosec // path is trusted from discovery
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return p.Parse(content)
}
