package fetcher

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestGitFetcher_Fetch_ClonesRepo(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("skipping integration test; set INTEGRATION_TEST=true to run")
	}

	tmpDir, err := os.MkdirTemp("", "fetcher-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	fetcher := NewGitFetcher(tmpDir)

	// Use a small public repo for testing
	result, err := fetcher.Fetch(context.Background(), FetchOptions{
		RepoURL: "https://github.com/open-telemetry/opentelemetry-collector-contrib",
		Shallow: true,
		Depth:   1,
	})
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	if result.RepoPath == "" {
		t.Error("RepoPath should not be empty")
	}

	if result.Commit == "" {
		t.Error("Commit should not be empty")
	}

	// Verify repo was cloned
	if _, err := os.Stat(filepath.Join(result.RepoPath, ".git")); os.IsNotExist(err) {
		t.Error(".git directory should exist")
	}
}

func TestGitFetcher_Fetch_UsesCachedRepo(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("skipping integration test; set INTEGRATION_TEST=true to run")
	}

	tmpDir, err := os.MkdirTemp("", "fetcher-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	fetcher := NewGitFetcher(tmpDir)
	opts := FetchOptions{
		RepoURL: "https://github.com/open-telemetry/opentelemetry-collector-contrib",
		Shallow: true,
		Depth:   1,
	}

	// First fetch
	result1, err := fetcher.Fetch(context.Background(), opts)
	if err != nil {
		t.Fatalf("First fetch failed: %v", err)
	}

	// Second fetch should use cache
	result2, err := fetcher.Fetch(context.Background(), opts)
	if err != nil {
		t.Fatalf("Second fetch failed: %v", err)
	}

	if result1.RepoPath != result2.RepoPath {
		t.Error("Second fetch should use same repo path")
	}
}

func TestGitFetcher_repoDir(t *testing.T) {
	fetcher := NewGitFetcher("/tmp/cache")

	dir := fetcher.repoDir("https://github.com/open-telemetry/opentelemetry-collector-contrib")
	expected := "/tmp/cache/github.com/open-telemetry/opentelemetry-collector-contrib"

	if dir != expected {
		t.Errorf("repoDir() = %q, want %q", dir, expected)
	}
}

func TestGitFetcher_repoDir_WithGitSuffix(t *testing.T) {
	fetcher := NewGitFetcher("/tmp/cache")

	dir := fetcher.repoDir("https://github.com/open-telemetry/opentelemetry-collector-contrib.git")
	expected := "/tmp/cache/github.com/open-telemetry/opentelemetry-collector-contrib"

	if dir != expected {
		t.Errorf("repoDir() = %q, want %q", dir, expected)
	}
}
