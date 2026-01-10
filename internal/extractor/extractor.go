package extractor

import (
	"time"

	"github.com/base-14/metric-library/internal/domain"
	"github.com/base-14/metric-library/internal/parser"
)

type MetricExtractor struct {
	sourceName    string
	componentName string
	componentType string
}

func NewMetricExtractor(sourceName, componentName, componentType string) *MetricExtractor {
	return &MetricExtractor{
		sourceName:    sourceName,
		componentName: componentName,
		componentType: componentType,
	}
}

func (e *MetricExtractor) Extract(meta *parser.Metadata) ([]*domain.CanonicalMetric, error) {
	var metrics []*domain.CanonicalMetric

	attrDefs := meta.Attributes

	for name, def := range meta.Metrics {
		m := &domain.CanonicalMetric{
			MetricName:       name,
			Description:      def.Description,
			Unit:             def.Unit,
			InstrumentType:   e.mapInstrumentType(def),
			EnabledByDefault: def.Enabled,
			ComponentType:    domain.ComponentType(e.componentType),
			ComponentName:    e.componentName,
			SourceCategory:   domain.SourceOTEL,
			SourceName:       e.sourceName,
			ExtractionMethod: domain.ExtractionMetadata,
			SourceConfidence: domain.ConfidenceAuthoritative,
			ExtractedAt:      time.Now(),
		}

		m.Attributes = e.extractAttributes(def.Attributes, attrDefs)
		m.EnsureID()

		metrics = append(metrics, m)
	}

	return metrics, nil
}

func (e *MetricExtractor) mapInstrumentType(def parser.MetricDefinition) domain.InstrumentType {
	if def.Sum != nil {
		if def.Sum.Monotonic {
			return domain.InstrumentCounter
		}
		return domain.InstrumentUpDownCounter
	}
	if def.Gauge != nil {
		return domain.InstrumentGauge
	}
	if def.Histogram != nil {
		return domain.InstrumentHistogram
	}
	return domain.InstrumentGauge
}

func (e *MetricExtractor) extractAttributes(
	attrNames []string,
	attrDefs map[string]parser.AttributeDefinition,
) []domain.Attribute {
	var attrs []domain.Attribute

	for _, name := range attrNames {
		attr := domain.Attribute{
			Name: name,
		}

		if def, ok := attrDefs[name]; ok {
			attr.Description = def.Description
			attr.Type = def.Type
			attr.Enum = def.Enum
		}

		attrs = append(attrs, attr)
	}

	return attrs
}
