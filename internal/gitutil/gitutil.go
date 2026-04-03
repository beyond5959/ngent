package gitutil

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	// ErrGitUnavailable indicates the host does not have the git binary.
	ErrGitUnavailable = errors.New("git is unavailable")
	// ErrNotRepository indicates the target cwd is not inside a git worktree.
	ErrNotRepository = errors.New("cwd is not inside a git repository")
	// ErrBranchNotFound indicates the requested branch does not exist locally.
	ErrBranchNotFound = errors.New("git branch not found")
	// ErrDiffFileNotFound indicates the requested path is not part of the visible diff.
	ErrDiffFileNotFound = errors.New("git diff file not found")
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

// DiffSummary captures aggregate diff counts for a worktree.
type DiffSummary struct {
	FilesChanged int
	Insertions   int
	Deletions    int
}

// DiffFile captures per-file diff counts for one worktree.
type DiffFile struct {
	Path      string
	Added     int
	Deleted   int
	Binary    bool
	Untracked bool
	Viewable  bool
}

// DiffStatus captures the current visible diff state for one worktree.
type DiffStatus struct {
	RepoRoot string
	Summary  DiffSummary
	Files    []DiffFile
}

const (
	// DiffFileDetailKindDiff indicates the content came from `git diff -- <path>`.
	DiffFileDetailKindDiff = "diff"
	// DiffFileDetailKindFile indicates the content came from directly reading the file.
	DiffFileDetailKindFile = "file"

	// DiffFileDetailReasonBinary indicates git marked the file as binary.
	DiffFileDetailReasonBinary = "binary"
	// DiffFileDetailReasonNonText indicates the file is not a text file we can preview.
	DiffFileDetailReasonNonText = "non_text"
)

// DiffFileDetail captures the preview payload for one changed file.
type DiffFileDetail struct {
	RepoRoot  string
	Path      string
	Kind      string
	Content   string
	Supported bool
	Reason    string
}

var (
	shortstatFilesPattern      = regexp.MustCompile(`(\d+)\s+files?\s+changed`)
	shortstatInsertionsPattern = regexp.MustCompile(`(\d+)\s+insertions?\(\+\)`)
	shortstatDeletionsPattern  = regexp.MustCompile(`(\d+)\s+deletions?\(-\)`)
)

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

// Diff loads the current visible working-tree diff summary for one worktree,
// including untracked files.
func Diff(ctx context.Context, cwd string) (DiffStatus, error) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return DiffStatus{}, ErrGitUnavailable
	}

	repoRoot, err := repoRoot(ctx, gitPath, cwd)
	if err != nil {
		return DiffStatus{}, err
	}

	shortstatOut, stderr, err := runGit(ctx, gitPath, cwd, "--no-pager", "diff", "--shortstat")
	if err != nil {
		if looksLikeNotRepository(stderr) {
			return DiffStatus{}, ErrNotRepository
		}
		return DiffStatus{}, fmt.Errorf("load diff shortstat: %w", err)
	}

	numstatOut, stderr, err := runGit(ctx, gitPath, cwd, "--no-pager", "diff", "--numstat")
	if err != nil {
		if looksLikeNotRepository(stderr) {
			return DiffStatus{}, ErrNotRepository
		}
		return DiffStatus{}, fmt.Errorf("load diff numstat: %w", err)
	}

	untrackedOut, stderr, err := runGit(ctx, gitPath, cwd, "ls-files", "--others", "--exclude-standard", "-z")
	if err != nil {
		if looksLikeNotRepository(stderr) {
			return DiffStatus{}, ErrNotRepository
		}
		return DiffStatus{}, fmt.Errorf("load untracked files: %w", err)
	}

	summary, err := parseShortstat(shortstatOut)
	if err != nil {
		return DiffStatus{}, fmt.Errorf("parse diff shortstat: %w", err)
	}
	files, err := parseNumstat(numstatOut)
	if err != nil {
		return DiffStatus{}, fmt.Errorf("parse diff numstat: %w", err)
	}
	for i := range files {
		files[i].Viewable = !files[i].Binary
	}
	untrackedFiles := parseUntrackedFiles(untrackedOut)
	markUntrackedDiffFiles(repoRoot, untrackedFiles)
	if len(untrackedFiles) > 0 {
		summary.FilesChanged += len(untrackedFiles)
		files = append(files, untrackedFiles...)
	}
	if summary.FilesChanged == 0 && len(files) > 0 {
		summary.FilesChanged = len(files)
	}

	return DiffStatus{
		RepoRoot: repoRoot,
		Summary:  summary,
		Files:    files,
	}, nil
}

// FileDetail loads previewable diff content for one visible diff file.
func FileDetail(ctx context.Context, cwd, rawPath string) (DiffFileDetail, error) {
	normalizedPath, err := normalizeRepoRelativePath(rawPath)
	if err != nil {
		return DiffFileDetail{}, err
	}

	status, err := Diff(ctx, cwd)
	if err != nil {
		return DiffFileDetail{}, err
	}

	file, ok := diffFileByPath(status.Files, normalizedPath)
	if !ok {
		return DiffFileDetail{}, ErrDiffFileNotFound
	}
	if !file.Viewable {
		reason := DiffFileDetailReasonBinary
		if file.Untracked {
			reason = DiffFileDetailReasonNonText
		}
		return DiffFileDetail{
			RepoRoot:  status.RepoRoot,
			Path:      normalizedPath,
			Supported: false,
			Reason:    reason,
		}, nil
	}

	if file.Untracked {
		fullPath := filepath.Join(status.RepoRoot, filepath.FromSlash(normalizedPath))
		content, ok, readErr := readTextFile(fullPath)
		if readErr != nil {
			return DiffFileDetail{}, readErr
		}
		if !ok {
			return DiffFileDetail{
				RepoRoot:  status.RepoRoot,
				Path:      normalizedPath,
				Supported: false,
				Reason:    DiffFileDetailReasonNonText,
			}, nil
		}
		return DiffFileDetail{
			RepoRoot:  status.RepoRoot,
			Path:      normalizedPath,
			Kind:      DiffFileDetailKindFile,
			Content:   content,
			Supported: true,
		}, nil
	}

	gitPath, err := exec.LookPath("git")
	if err != nil {
		return DiffFileDetail{}, ErrGitUnavailable
	}

	diffOut, stderr, err := runGit(ctx, gitPath, cwd, "--no-pager", "diff", "--", normalizedPath)
	if err != nil {
		if looksLikeNotRepository(stderr) {
			return DiffFileDetail{}, ErrNotRepository
		}
		return DiffFileDetail{}, fmt.Errorf("load diff for %q: %w", normalizedPath, err)
	}

	return DiffFileDetail{
		RepoRoot:  status.RepoRoot,
		Path:      normalizedPath,
		Kind:      DiffFileDetailKindDiff,
		Content:   diffOut,
		Supported: true,
	}, nil
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

func parseShortstat(raw string) (DiffSummary, error) {
	line := strings.TrimSpace(raw)
	if line == "" {
		return DiffSummary{}, nil
	}

	filesChanged, err := parseShortstatCount(shortstatFilesPattern, line)
	if err != nil {
		return DiffSummary{}, fmt.Errorf("files changed: %w", err)
	}
	insertions, err := parseShortstatCount(shortstatInsertionsPattern, line)
	if err != nil {
		return DiffSummary{}, fmt.Errorf("insertions: %w", err)
	}
	deletions, err := parseShortstatCount(shortstatDeletionsPattern, line)
	if err != nil {
		return DiffSummary{}, fmt.Errorf("deletions: %w", err)
	}

	return DiffSummary{
		FilesChanged: filesChanged,
		Insertions:   insertions,
		Deletions:    deletions,
	}, nil
}

func parseShortstatCount(pattern *regexp.Regexp, line string) (int, error) {
	match := pattern.FindStringSubmatch(line)
	if len(match) != 2 {
		return 0, nil
	}

	value, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, fmt.Errorf("parse %q: %w", match[1], err)
	}
	return value, nil
}

func parseNumstat(raw string) ([]DiffFile, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	files := make([]DiffFile, 0)
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "\t", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("unexpected line %q", line)
		}

		path := strings.TrimSpace(parts[2])
		if path == "" {
			return nil, fmt.Errorf("missing path in line %q", line)
		}

		file := DiffFile{Path: path}
		if parts[0] == "-" || parts[1] == "-" {
			file.Binary = true
			files = append(files, file)
			continue
		}

		added, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("parse added count in line %q: %w", line, err)
		}
		deleted, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("parse deleted count in line %q: %w", line, err)
		}
		file.Added = added
		file.Deleted = deleted
		files = append(files, file)
	}

	return files, nil
}

func parseUntrackedFiles(raw string) []DiffFile {
	if raw == "" {
		return nil
	}

	files := make([]DiffFile, 0)
	seen := make(map[string]struct{})
	for _, path := range strings.Split(raw, "\x00") {
		if path == "" {
			continue
		}
		if _, duplicate := seen[path]; duplicate {
			continue
		}
		seen[path] = struct{}{}
		files = append(files, DiffFile{
			Path:      path,
			Untracked: true,
		})
	}
	return files
}

func diffFileByPath(files []DiffFile, normalizedPath string) (DiffFile, bool) {
	for _, file := range files {
		if filepath.ToSlash(strings.TrimSpace(file.Path)) == normalizedPath {
			return file, true
		}
	}
	return DiffFile{}, false
}

func markUntrackedDiffFiles(repoRoot string, files []DiffFile) {
	for i := range files {
		fullPath := filepath.Join(repoRoot, filepath.FromSlash(files[i].Path))
		ok, err := isTextFile(fullPath)
		files[i].Viewable = err == nil && ok
		files[i].Binary = !files[i].Viewable
	}
}

func normalizeRepoRelativePath(rawPath string) (string, error) {
	trimmed := strings.TrimSpace(rawPath)
	if trimmed == "" {
		return "", fmt.Errorf("path is required")
	}

	cleaned := filepath.Clean(filepath.FromSlash(trimmed))
	if cleaned == "." || cleaned == string(filepath.Separator) || filepath.IsAbs(cleaned) {
		return "", fmt.Errorf("path must be a repository-relative file path")
	}
	if cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path must stay within the repository root")
	}
	return filepath.ToSlash(cleaned), nil
}

func readTextFile(path string) (content string, ok bool, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false, err
	}
	if !isLikelyText(data) {
		return "", false, nil
	}
	return string(data), true, nil
}

func isTextFile(path string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	return isLikelyText(data), nil
}

func isLikelyText(data []byte) bool {
	if len(data) == 0 {
		return true
	}
	if bytes.IndexByte(data, 0) >= 0 {
		return false
	}
	return utf8.Valid(data)
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
