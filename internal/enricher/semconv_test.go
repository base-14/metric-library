package enricher

import (
	"testing"

	"github.com/base-14/metric-library/internal/domain"
)

func TestSemconvEnricher_ExactMatch(t *testing.T) {
	index := []SemconvMetric{
		{Name: "http.server.request.duration", Stability: "stable"},
		{Name: "http.client.request.duration", Stability: "stable"},
		{Name: "system.cpu.time", Stability: "development"},
	}

	enricher := NewSemconvEnricher(index)

	metric := &domain.CanonicalMetric{
		MetricName: "http.server.request.duration",
	}

	enricher.Enrich(metric)

	if metric.SemconvMatch != domain.SemconvMatchExact {
		t.Errorf("expected exact match, got %s", metric.SemconvMatch)
	}
	if metric.SemconvName != "http.server.request.duration" {
		t.Errorf("expected semconv name http.server.request.duration, got %s", metric.SemconvName)
	}
	if metric.SemconvStability != "stable" {
		t.Errorf("expected stability stable, got %s", metric.SemconvStability)
	}
}

func TestSemconvEnricher_PrefixMatch(t *testing.T) {
	index := []SemconvMetric{
		{Name: "http.server.request.duration", Stability: "stable"},
		{Name: "http.server.active_requests", Stability: "development"},
	}

	enricher := NewSemconvEnricher(index)

	metric := &domain.CanonicalMetric{
		MetricName: "http.server.request.duration.bucket",
	}

	enricher.Enrich(metric)

	if metric.SemconvMatch != domain.SemconvMatchPrefix {
		t.Errorf("expected prefix match, got %s", metric.SemconvMatch)
	}
	if metric.SemconvName != "http.server.request.duration" {
		t.Errorf("expected semconv name http.server.request.duration, got %s", metric.SemconvName)
	}
}

func TestSemconvEnricher_NoMatch(t *testing.T) {
	index := []SemconvMetric{
		{Name: "http.server.request.duration", Stability: "stable"},
	}

	enricher := NewSemconvEnricher(index)

	metric := &domain.CanonicalMetric{
		MetricName: "custom_metric_name",
	}

	enricher.Enrich(metric)

	if metric.SemconvMatch != domain.SemconvMatchNone {
		t.Errorf("expected no match, got %s", metric.SemconvMatch)
	}
	if metric.SemconvName != "" {
		t.Errorf("expected empty semconv name, got %s", metric.SemconvName)
	}
}

func TestSemconvEnricher_UnderscoreNormalization(t *testing.T) {
	index := []SemconvMetric{
		{Name: "http.server.request.duration", Stability: "stable"},
	}

	enricher := NewSemconvEnricher(index)

	metric := &domain.CanonicalMetric{
		MetricName: "http_server_request_duration",
	}

	enricher.Enrich(metric)

	if metric.SemconvMatch != domain.SemconvMatchExact {
		t.Errorf("expected exact match after normalization, got %s", metric.SemconvMatch)
	}
}

func TestSemconvEnricher_EnrichAll(t *testing.T) {
	index := []SemconvMetric{
		{Name: "http.server.request.duration", Stability: "stable"},
		{Name: "system.cpu.time", Stability: "development"},
	}

	enricher := NewSemconvEnricher(index)

	metrics := []*domain.CanonicalMetric{
		{MetricName: "http.server.request.duration"},
		{MetricName: "system.cpu.time"},
		{MetricName: "custom_metric"},
	}

	enricher.EnrichAll(metrics)

	if metrics[0].SemconvMatch != domain.SemconvMatchExact {
		t.Errorf("first metric should have exact match")
	}
	if metrics[1].SemconvMatch != domain.SemconvMatchExact {
		t.Errorf("second metric should have exact match")
	}
	if metrics[2].SemconvMatch != domain.SemconvMatchNone {
		t.Errorf("third metric should have no match")
	}
}
