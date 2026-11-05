package fetcher

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
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

func newSourceRepo(t *testing.T) (string, func(msg string) string) {
	t.Helper()
	dir := t.TempDir()
	repo, err := git.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("init source repo: %v", err)
	}
	w, err := repo.Worktree()
	if err != nil {
		t.Fatalf("worktree: %v", err)
	}
	commit := func(msg string) string {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, "file.txt"), []byte(msg), 0600); err != nil {
			t.Fatalf("write file: %v", err)
		}
		if _, err := w.Add("file.txt"); err != nil {
			t.Fatalf("add: %v", err)
		}
		hash, err := w.Commit(msg, &git.CommitOptions{
			Author: &object.Signature{Name: "test", Email: "test@example.com", When: time.Now()},
		})
		if err != nil {
			t.Fatalf("commit: %v", err)
		}
		return hash.String()
	}
	return dir, commit
}

func TestGitFetcher_Fetch_PullsNewCommits(t *testing.T) {
	src, commit := newSourceRepo(t)
	c1 := commit("one")
	fetcher := NewGitFetcher(t.TempDir())
	opts := FetchOptions{RepoURL: src}

	result, err := fetcher.Fetch(context.Background(), opts)
	if err != nil {
		t.Fatalf("first fetch: %v", err)
	}
	if result.Commit != c1 {
		t.Fatalf("first fetch commit = %s, want %s", result.Commit, c1)
	}

	c2 := commit("two")
	result, err = fetcher.Fetch(context.Background(), opts)
	if err != nil {
		t.Fatalf("second fetch: %v", err)
	}
	if result.Commit != c2 {
		t.Errorf("second fetch commit = %s, want %s", result.Commit, c2)
	}
}

func TestGitFetcher_Fetch_RecloneOnPullFailure(t *testing.T) {
	src, commit := newSourceRepo(t)
	commit("one")
	fetcher := NewGitFetcher(t.TempDir())
	opts := FetchOptions{RepoURL: src}

	result, err := fetcher.Fetch(context.Background(), opts)
	if err != nil {
		t.Fatalf("first fetch: %v", err)
	}

	c2 := commit("two")
	cached, err := git.PlainOpen(result.RepoPath)
	if err != nil {
		t.Fatalf("open cached repo: %v", err)
	}
	cfg, err := cached.Config()
	if err != nil {
		t.Fatalf("cached config: %v", err)
	}
	cfg.Remotes["origin"].URLs = []string{filepath.Join(t.TempDir(), "gone")}
	if err := cached.SetConfig(cfg); err != nil {
		t.Fatalf("set config: %v", err)
	}

	result, err = fetcher.Fetch(context.Background(), opts)
	if err != nil {
		t.Fatalf("second fetch: %v", err)
	}
	if result.Commit != c2 {
		t.Errorf("second fetch commit = %s, want %s", result.Commit, c2)
	}
}

func TestGitFetcher_Fetch_UpToDateKeepsCache(t *testing.T) {
	src, commit := newSourceRepo(t)
	c1 := commit("one")
	fetcher := NewGitFetcher(t.TempDir())
	opts := FetchOptions{RepoURL: src}

	result, err := fetcher.Fetch(context.Background(), opts)
	if err != nil {
		t.Fatalf("first fetch: %v", err)
	}
	marker := filepath.Join(result.RepoPath, ".git", "marker")
	if err := os.WriteFile(marker, []byte("x"), 0600); err != nil {
		t.Fatalf("write marker: %v", err)
	}

	result, err = fetcher.Fetch(context.Background(), opts)
	if err != nil {
		t.Fatalf("second fetch: %v", err)
	}
	if result.Commit != c1 {
		t.Errorf("second fetch commit = %s, want %s", result.Commit, c1)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Errorf("cached repo was replaced: %v", err)
	}
}
