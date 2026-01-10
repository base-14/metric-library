package semconv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "metrics.yaml")
	content := `groups:
  - id: metric.http.server.request.duration
    type: metric
    metric_name: http.server.request.duration
    brief: "Duration of HTTP server requests."
    instrument: histogram
    unit: "s"
    stability: stable
    attributes:
      - ref: http.request.method
        requirement_level: required
      - ref: url.scheme
        requirement_level: required

  - id: metric.http.server.active_requests
    type: metric
    metric_name: http.server.active_requests
    stability: development
    brief: "Number of active HTTP server requests."
    instrument: updowncounter
    unit: "{request}"

  - id: metric_attributes.http.server
    type: attribute_group
    brief: 'HTTP server attributes'
`
	if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	defs, err := ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(defs) != 2 {
		t.Errorf("expected 2 metrics, got %d", len(defs))
		for _, d := range defs {
			t.Logf("  metric: %s", d.Name)
		}
	}

	expected := map[string]struct {
		brief      string
		instrument string
		unit       string
		stability  string
	}{
		"http.server.request.duration": {"Duration of HTTP server requests.", "histogram", "s", "stable"},
		"http.server.active_requests":  {"Number of active HTTP server requests.", "updowncounter", "{request}", "development"},
	}

	for _, def := range defs {
		exp, ok := expected[def.Name]
		if !ok {
			t.Errorf("unexpected metric: %s", def.Name)
			continue
		}

		if def.Brief != exp.brief {
			t.Errorf("metric %s: expected brief %q, got %q", def.Name, exp.brief, def.Brief)
		}

		if def.Instrument != exp.instrument {
			t.Errorf("metric %s: expected instrument %q, got %q", def.Name, exp.instrument, def.Instrument)
		}

		if def.Unit != exp.unit {
			t.Errorf("metric %s: expected unit %q, got %q", def.Name, exp.unit, def.Unit)
		}

		if def.Stability != exp.stability {
			t.Errorf("metric %s: expected stability %q, got %q", def.Name, exp.stability, def.Stability)
		}
	}
}

func TestParseFileWithAttributes(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "metrics.yaml")
	content := `groups:
  - id: metric.system.cpu.time
    type: metric
    metric_name: system.cpu.time
    brief: "Seconds each logical CPU spent on each mode."
    instrument: counter
    unit: "s"
    stability: stable
    attributes:
      - ref: cpu.mode
        requirement_level: required
      - ref: cpu.logical_number
        requirement_level: opt_in
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

	if len(defs[0].Attributes) != 2 {
		t.Errorf("expected 2 attributes, got %d", len(defs[0].Attributes))
	}
}

func TestParseFileNonExistent(t *testing.T) {
	_, err := ParseFile("/nonexistent/file.yaml")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}
