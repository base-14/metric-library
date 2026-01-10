package enricher

import (
	"strings"

	"github.com/base-14/metric-library/internal/domain"
)

type SemconvMetric struct {
	Name      string
	Stability string
}

type SemconvEnricher struct {
	exactIndex  map[string]SemconvMetric
	prefixIndex []SemconvMetric
}

func NewSemconvEnricher(metrics []SemconvMetric) *SemconvEnricher {
	exactIndex := make(map[string]SemconvMetric)
	prefixIndex := make([]SemconvMetric, 0, len(metrics))

	for _, m := range metrics {
		normalized := normalize(m.Name)
		exactIndex[normalized] = m
		prefixIndex = append(prefixIndex, m)
	}

	return &SemconvEnricher{
		exactIndex:  exactIndex,
		prefixIndex: prefixIndex,
	}
}

func (e *SemconvEnricher) Enrich(metric *domain.CanonicalMetric) {
	normalized := normalize(metric.MetricName)

	if m, ok := e.exactIndex[normalized]; ok {
		metric.SemconvMatch = domain.SemconvMatchExact
		metric.SemconvName = m.Name
		metric.SemconvStability = m.Stability
		return
	}

	for _, m := range e.prefixIndex {
		semconvNorm := normalize(m.Name)
		if strings.HasPrefix(normalized, semconvNorm+".") || strings.HasPrefix(normalized, semconvNorm+"_") {
			metric.SemconvMatch = domain.SemconvMatchPrefix
			metric.SemconvName = m.Name
			metric.SemconvStability = m.Stability
			return
		}
	}

	metric.SemconvMatch = domain.SemconvMatchNone
}

func (e *SemconvEnricher) EnrichAll(metrics []*domain.CanonicalMetric) {
	for _, m := range metrics {
		e.Enrich(m)
	}
}

func normalize(name string) string {
	return strings.ReplaceAll(name, "_", ".")
}
