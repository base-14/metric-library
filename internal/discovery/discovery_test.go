package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMetadataDiscovery_FindMetadataFiles(t *testing.T) {
	tmpDir := setupTestRepo(t)

	discovery := NewMetadataDiscovery()
	files, err := discovery.FindMetadataFiles(tmpDir)
	if err != nil {
		t.Fatalf("FindMetadataFiles failed: %v", err)
	}

	if len(files) != 3 {
		t.Errorf("expected 3 metadata files, got %d", len(files))
	}

	expectedPaths := map[string]bool{
		filepath.Join(tmpDir, "receiver", "hostmetrics", "metadata.yaml"): false,
		filepath.Join(tmpDir, "processor", "transform", "metadata.yaml"):  false,
		filepath.Join(tmpDir, "exporter", "prometheus", "metadata.yaml"):  false,
	}

	for _, f := range files {
		if _, ok := expectedPaths[f.Path]; ok {
			expectedPaths[f.Path] = true
		}
	}

	for path, found := range expectedPaths {
		if !found {
			t.Errorf("expected to find %s", path)
		}
	}
}

func TestMetadataDiscovery_FindMetadataFiles_ComponentType(t *testing.T) {
	tmpDir := setupTestRepo(t)

	discovery := NewMetadataDiscovery()
	files, err := discovery.FindMetadataFiles(tmpDir)
	if err != nil {
		t.Fatalf("FindMetadataFiles failed: %v", err)
	}

	componentTypes := make(map[string]string)
	for _, f := range files {
		componentTypes[f.ComponentName] = f.ComponentType
	}

	tests := []struct {
		component    string
		expectedType string
	}{
		{"hostmetrics", "receiver"},
		{"transform", "processor"},
		{"prometheus", "exporter"},
	}

	for _, tc := range tests {
		if componentTypes[tc.component] != tc.expectedType {
			t.Errorf("expected %s to be %s, got %s", tc.component, tc.expectedType, componentTypes[tc.component])
		}
	}
}

func TestMetadataDiscovery_FindMetadataFiles_IgnoresNonMetadata(t *testing.T) {
	tmpDir := setupTestRepo(t)

	otherYAML := filepath.Join(tmpDir, "receiver", "hostmetrics", "config.yaml")
	if err := os.WriteFile(otherYAML, []byte("config: test"), 0600); err != nil {
		t.Fatalf("failed to create config.yaml: %v", err)
	}

	discovery := NewMetadataDiscovery()
	files, err := discovery.FindMetadataFiles(tmpDir)
	if err != nil {
		t.Fatalf("FindMetadataFiles failed: %v", err)
	}

	for _, f := range files {
		if filepath.Base(f.Path) != "metadata.yaml" {
			t.Errorf("found non-metadata file: %s", f.Path)
		}
	}
}

func TestMetadataDiscovery_FindMetadataFiles_EmptyRepo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "discovery-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	discovery := NewMetadataDiscovery()
	files, err := discovery.FindMetadataFiles(tmpDir)
	if err != nil {
		t.Fatalf("FindMetadataFiles failed: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("expected 0 files in empty repo, got %d", len(files))
	}
}

func TestMetadataDiscovery_FindMetadataFiles_AllComponentTypes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "discovery-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	componentTypes := []string{"receiver", "processor", "exporter", "extension", "connector"}
	for _, ct := range componentTypes {
		dir := filepath.Join(tmpDir, ct, "test"+ct)
		if err := os.MkdirAll(dir, 0750); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
		if err := os.WriteFile(filepath.Join(dir, "metadata.yaml"), []byte("type: "+ct), 0600); err != nil {
			t.Fatalf("failed to create metadata.yaml: %v", err)
		}
	}

	discovery := NewMetadataDiscovery()
	files, err := discovery.FindMetadataFiles(tmpDir)
	if err != nil {
		t.Fatalf("FindMetadataFiles failed: %v", err)
	}

	if len(files) != 5 {
		t.Errorf("expected 5 metadata files, got %d", len(files))
	}

	foundTypes := make(map[string]bool)
	for _, f := range files {
		foundTypes[f.ComponentType] = true
	}

	for _, ct := range componentTypes {
		if !foundTypes[ct] {
			t.Errorf("missing component type: %s", ct)
		}
	}
}

func setupTestRepo(t *testing.T) string {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "discovery-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tmpDir) })

	structure := map[string]string{
		"receiver/hostmetrics/metadata.yaml": `type: hostmetricsreceiver
status:
  class: receiver
`,
		"processor/transform/metadata.yaml": `type: transformprocessor
status:
  class: processor
`,
		"exporter/prometheus/metadata.yaml": `type: prometheusexporter
status:
  class: exporter
`,
	}

	for path, content := range structure {
		fullPath := filepath.Join(tmpDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0750); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0600); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
	}

	return tmpDir
}
