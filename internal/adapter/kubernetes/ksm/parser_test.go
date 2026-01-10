package ksm

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "test.go")
	content := `package store

import (
	basemetrics "k8s.io/component-base/metrics"
	"k8s.io/kube-state-metrics/v2/pkg/metric"
	generator "k8s.io/kube-state-metrics/v2/pkg/metric_generator"
)

func podMetricFamilies() []generator.FamilyGenerator {
	return []generator.FamilyGenerator{
		*generator.NewFamilyGeneratorWithStability(
			"kube_pod_info",
			"Information about pod.",
			metric.Gauge,
			basemetrics.STABLE,
			"",
			wrapPodFunc(func(p *v1.Pod) *metric.Family {
				return nil
			}),
		),
		*generator.NewFamilyGeneratorWithStability(
			"kube_pod_start_time",
			"Start time in unix timestamp for a pod.",
			metric.Gauge,
			basemetrics.STABLE,
			"",
			wrapPodFunc(func(p *v1.Pod) *metric.Family {
				return nil
			}),
		),
	}
}

func createPodCompletionTimeFamilyGenerator() generator.FamilyGenerator {
	return *generator.NewFamilyGeneratorWithStability(
		"kube_pod_completion_time",
		"Completion time in unix timestamp for a pod.",
		metric.Gauge,
		basemetrics.STABLE,
		"",
		wrapPodFunc(func(p *v1.Pod) *metric.Family {
			return nil
		}),
	)
}
`
	if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	defs, err := ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(defs) != 3 {
		t.Errorf("expected 3 metrics, got %d", len(defs))
		for _, d := range defs {
			t.Logf("  metric: %s", d.Name)
		}
	}

	expected := map[string]struct {
		help       string
		metricType string
	}{
		"kube_pod_info":            {"Information about pod.", "gauge"},
		"kube_pod_start_time":      {"Start time in unix timestamp for a pod.", "gauge"},
		"kube_pod_completion_time": {"Completion time in unix timestamp for a pod.", "gauge"},
	}

	for _, def := range defs {
		exp, ok := expected[def.Name]
		if !ok {
			t.Errorf("unexpected metric: %s", def.Name)
			continue
		}

		if def.Help != exp.help {
			t.Errorf("metric %s: expected help %q, got %q", def.Name, exp.help, def.Help)
		}

		if def.MetricType != exp.metricType {
			t.Errorf("metric %s: expected type %q, got %q", def.Name, exp.metricType, def.MetricType)
		}
	}
}

func TestParseFileWithCounter(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "test.go")
	content := `package store

import (
	basemetrics "k8s.io/component-base/metrics"
	"k8s.io/kube-state-metrics/v2/pkg/metric"
	generator "k8s.io/kube-state-metrics/v2/pkg/metric_generator"
)

func containerMetricFamilies() []generator.FamilyGenerator {
	return []generator.FamilyGenerator{
		*generator.NewFamilyGeneratorWithStability(
			"kube_pod_container_status_restarts_total",
			"Number of container restarts.",
			metric.Counter,
			basemetrics.STABLE,
			"",
			nil,
		),
	}
}
`
	if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	defs, err := ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(defs) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(defs))
	}

	if defs[0].MetricType != "counter" {
		t.Errorf("expected type counter, got %s", defs[0].MetricType)
	}
}

func TestParseFileNonExistent(t *testing.T) {
	_, err := ParseFile("/nonexistent/file.go")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}
