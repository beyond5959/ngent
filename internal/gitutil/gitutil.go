package gitutil

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

var (
	// ErrGitUnavailable indicates the host does not have the git binary.
	ErrGitUnavailable = errors.New("git is unavailable")
	// ErrNotRepository indicates the target cwd is not inside a git worktree.
	ErrNotRepository = errors.New("cwd is not inside a git repository")
	// ErrBranchNotFound indicates the requested branch does not exist locally.
	ErrBranchNotFound = errors.New("git branch not found")
)

// Branch describes one local branch.
type Branch struct {
	Name    string
	Current bool
}

// Status captures the visible git state for one working tree.
type Status struct {
	RepoRoot      string
	CurrentRef    string
	CurrentBranch string
	Detached      bool
	Branches      []Branch
}

// Inspect loads local branch information for one working tree.
func Inspect(ctx context.Context, cwd string) (Status, error) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return Status{}, ErrGitUnavailable
	}

	repoRoot, err := repoRoot(ctx, gitPath, cwd)
	if err != nil {
		return Status{}, err
	}

	currentBranch, detached, currentRef, err := currentRef(ctx, gitPath, cwd)
	if err != nil {
		return Status{}, err
	}

	branches, err := localBranches(ctx, gitPath, cwd, currentBranch)
	if err != nil {
		return Status{}, err
	}

	return Status{
		RepoRoot:      repoRoot,
		CurrentRef:    currentRef,
		CurrentBranch: currentBranch,
		Detached:      detached,
		Branches:      branches,
	}, nil
}

// Checkout switches the worktree to one existing local branch and returns the
// refreshed repository state.
func Checkout(ctx context.Context, cwd, branch string) (Status, error) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return Status{}, ErrGitUnavailable
	}

	branch = strings.TrimSpace(branch)
	if branch == "" {
		return Status{}, ErrBranchNotFound
	}

	if _, err := repoRoot(ctx, gitPath, cwd); err != nil {
		return Status{}, err
	}

	branches, err := localBranches(ctx, gitPath, cwd, "")
	if err != nil {
		return Status{}, err
	}
	if !hasBranch(branches, branch) {
		return Status{}, ErrBranchNotFound
	}

	if _, _, err := runGit(ctx, gitPath, cwd, "checkout", "--quiet", branch); err != nil {
		return Status{}, fmt.Errorf("checkout %q: %w", branch, err)
	}

	return Inspect(ctx, cwd)
}

func repoRoot(ctx context.Context, gitPath, cwd string) (string, error) {
	stdout, stderr, err := runGit(ctx, gitPath, cwd, "rev-parse", "--show-toplevel")
	if err != nil {
		if looksLikeNotRepository(stderr) {
			return "", ErrNotRepository
		}
		return "", fmt.Errorf("resolve repo root: %w", err)
	}

	root := filepath.Clean(strings.TrimSpace(stdout))
	if root == "." || root == "" {
		return "", fmt.Errorf("resolve repo root: empty output")
	}
	return root, nil
}

func currentRef(ctx context.Context, gitPath, cwd string) (branch string, detached bool, label string, err error) {
	stdout, stderr, err := runGit(ctx, gitPath, cwd, "symbolic-ref", "--quiet", "--short", "HEAD")
	if err == nil {
		branch = strings.TrimSpace(stdout)
		if branch == "" {
			return "", false, "", fmt.Errorf("resolve current branch: empty output")
		}
		return branch, false, branch, nil
	}
	if looksLikeNotRepository(stderr) {
		return "", false, "", ErrNotRepository
	}

	sha, shaStderr, shaErr := runGit(ctx, gitPath, cwd, "rev-parse", "--short", "HEAD")
	if shaErr != nil {
		if looksLikeNotRepository(shaStderr) {
			return "", false, "", ErrNotRepository
		}
		return "", false, "", fmt.Errorf("resolve detached HEAD: %w", shaErr)
	}

	shortSHA := strings.TrimSpace(sha)
	if shortSHA == "" {
		return "", false, "", fmt.Errorf("resolve detached HEAD: empty output")
	}
	return "", true, "Detached HEAD @" + " " + shortSHA, nil
}

func localBranches(ctx context.Context, gitPath, cwd, currentBranch string) ([]Branch, error) {
	stdout, stderr, err := runGit(ctx, gitPath, cwd, "for-each-ref", "--format=%(refname:short)\t%(HEAD)", "refs/heads")
	if err != nil {
		if looksLikeNotRepository(stderr) {
			return nil, ErrNotRepository
		}
		return nil, fmt.Errorf("list local branches: %w", err)
	}

	branches := make([]Branch, 0)
	seen := make(map[string]struct{})
	for _, line := range strings.Split(stdout, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "\t", 2)
		name := strings.TrimSpace(parts[0])
		if name == "" {
			continue
		}
		_, duplicate := seen[name]
		if duplicate {
			continue
		}
		seen[name] = struct{}{}

		headMarker := ""
		if len(parts) == 2 {
			headMarker = strings.TrimSpace(parts[1])
		}
		branches = append(branches, Branch{
			Name:    name,
			Current: headMarker == "*" || name == currentBranch,
		})
	}

	if currentBranch != "" {
		if _, ok := seen[currentBranch]; !ok {
			branches = append(branches, Branch{
				Name:    currentBranch,
				Current: true,
			})
		}
	}

	sort.Slice(branches, func(i, j int) bool {
		if branches[i].Current != branches[j].Current {
			return branches[i].Current
		}
		return strings.ToLower(branches[i].Name) < strings.ToLower(branches[j].Name)
	})

	return branches, nil
}

func hasBranch(branches []Branch, branch string) bool {
	for _, candidate := range branches {
		if candidate.Name == branch {
			return true
		}
	}
	return false
}

func runGit(ctx context.Context, gitPath, cwd string, args ...string) (stdout string, stderr string, err error) {
	cmd := exec.CommandContext(ctx, gitPath, args...)
	cmd.Dir = cwd

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()
	return stdoutBuf.String(), stderrBuf.String(), err
}

func looksLikeNotRepository(stderr string) bool {
	normalized := strings.ToLower(strings.TrimSpace(stderr))
	return strings.Contains(normalized, "not a git repository")
}
