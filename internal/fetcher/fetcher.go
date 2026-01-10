package fetcher

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type FetchOptions struct {
	RepoURL string
	Commit  string
	Shallow bool
	Depth   int
	Force   bool
}

type FetchResult struct {
	RepoPath  string
	Commit    string
	Timestamp time.Time
}

type GitFetcher struct {
	cacheDir string
}

func NewGitFetcher(cacheDir string) *GitFetcher {
	return &GitFetcher{cacheDir: cacheDir}
}

func (f *GitFetcher) Fetch(ctx context.Context, opts FetchOptions) (*FetchResult, error) {
	repoDir := f.repoDir(opts.RepoURL)

	// Check if repo already exists
	if _, err := os.Stat(filepath.Join(repoDir, ".git")); err == nil && !opts.Force {
		return f.openExisting(ctx, repoDir, opts)
	}

	// Clone the repo
	return f.cloneRepo(ctx, repoDir, opts)
}

func (f *GitFetcher) cloneRepo(ctx context.Context, repoDir string, opts FetchOptions) (*FetchResult, error) {
	if err := os.MkdirAll(filepath.Dir(repoDir), 0750); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Remove existing directory if force
	if opts.Force {
		_ = os.RemoveAll(repoDir)
	}

	cloneOpts := &git.CloneOptions{
		URL:      opts.RepoURL,
		Progress: nil,
	}

	if opts.Shallow {
		cloneOpts.Depth = opts.Depth
		if cloneOpts.Depth == 0 {
			cloneOpts.Depth = 1
		}
	}

	repo, err := git.PlainCloneContext(ctx, repoDir, false, cloneOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to clone repository: %w", err)
	}

	return f.getResult(repo, repoDir, opts)
}

func (f *GitFetcher) openExisting(ctx context.Context, repoDir string, opts FetchOptions) (*FetchResult, error) {
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Pull latest changes if not pinned to specific commit
	if opts.Commit == "" {
		w, err := repo.Worktree()
		if err != nil {
			return nil, fmt.Errorf("failed to get worktree: %w", err)
		}

		pullOpts := &git.PullOptions{
			RemoteName: "origin",
		}

		_ = w.PullContext(ctx, pullOpts)
	}

	return f.getResult(repo, repoDir, opts)
}

func (f *GitFetcher) getResult(repo *git.Repository, repoDir string, opts FetchOptions) (*FetchResult, error) {
	// Checkout specific commit if provided
	if opts.Commit != "" {
		w, err := repo.Worktree()
		if err != nil {
			return nil, fmt.Errorf("failed to get worktree: %w", err)
		}

		err = w.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(opts.Commit),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to checkout commit %s: %w", opts.Commit, err)
		}
	}

	// Get HEAD commit
	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get commit object: %w", err)
	}

	return &FetchResult{
		RepoPath:  repoDir,
		Commit:    head.Hash().String(),
		Timestamp: commit.Author.When,
	}, nil
}

func (f *GitFetcher) repoDir(repoURL string) string {
	parsed, err := url.Parse(repoURL)
	if err != nil {
		// Fallback to simple hash if URL parsing fails
		return filepath.Join(f.cacheDir, sanitizePath(repoURL))
	}

	// Remove .git suffix if present
	path := strings.TrimSuffix(parsed.Path, ".git")

	return filepath.Join(f.cacheDir, parsed.Host, path[1:]) // Remove leading slash
}

func sanitizePath(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '/' || r == ':' || r == '?' || r == '&' || r == '=' {
			return '_'
		}
		return r
	}, s)
}
