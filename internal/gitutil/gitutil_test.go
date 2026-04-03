package gitutil

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInspect(t *testing.T) {
	repo := newGitRepo(t)

	status, err := Inspect(context.Background(), filepath.Join(repo, "nested"))
	if err != nil {
		t.Fatalf("Inspect(): %v", err)
	}

	if got, want := evalPath(t, status.RepoRoot), evalPath(t, repo); got != want {
		t.Fatalf("RepoRoot = %q, want %q", got, want)
	}
	if got, want := status.CurrentBranch, "main"; got != want {
		t.Fatalf("CurrentBranch = %q, want %q", got, want)
	}
	if got, want := status.CurrentRef, "main"; got != want {
		t.Fatalf("CurrentRef = %q, want %q", got, want)
	}
	if status.Detached {
		t.Fatal("Detached = true, want false")
	}
	if len(status.Branches) != 2 {
		t.Fatalf("len(Branches) = %d, want 2", len(status.Branches))
	}
	if !status.Branches[0].Current || status.Branches[0].Name != "main" {
		t.Fatalf("Branches[0] = %#v, want current main", status.Branches[0])
	}
}

func TestCheckout(t *testing.T) {
	repo := newGitRepo(t)

	status, err := Checkout(context.Background(), repo, "feature/demo")
	if err != nil {
		t.Fatalf("Checkout(): %v", err)
	}

	if got, want := status.CurrentBranch, "feature/demo"; got != want {
		t.Fatalf("CurrentBranch = %q, want %q", got, want)
	}
	if !status.Branches[0].Current || status.Branches[0].Name != "feature/demo" {
		t.Fatalf("Branches[0] = %#v, want current feature/demo", status.Branches[0])
	}
}

func TestInspectNotRepository(t *testing.T) {
	ensureGitAvailable(t)

	_, err := Inspect(context.Background(), t.TempDir())
	if !errors.Is(err, ErrNotRepository) {
		t.Fatalf("Inspect() error = %v, want %v", err, ErrNotRepository)
	}
}

func TestCheckoutMissingBranch(t *testing.T) {
	repo := newGitRepo(t)

	_, err := Checkout(context.Background(), repo, "missing")
	if !errors.Is(err, ErrBranchNotFound) {
		t.Fatalf("Checkout() error = %v, want %v", err, ErrBranchNotFound)
	}
}

func TestDiff(t *testing.T) {
	repo := newGitRepo(t)

	appPath := filepath.Join(repo, "pkg", "app.txt")
	if err := os.MkdirAll(filepath.Dir(appPath), 0o755); err != nil {
		t.Fatalf("os.MkdirAll(): %v", err)
	}
	if err := os.WriteFile(appPath, []byte("alpha\nbeta\n"), 0o644); err != nil {
		t.Fatalf("os.WriteFile(appPath): %v", err)
	}
	runGitCommand(t, repo, "add", "pkg/app.txt")
	runGitCommand(t, repo, "commit", "--quiet", "-m", "add app")

	if err := os.WriteFile(filepath.Join(repo, "README.md"), []byte("hello\nworld\n"), 0o644); err != nil {
		t.Fatalf("os.WriteFile(README.md): %v", err)
	}
	if err := os.WriteFile(appPath, []byte("alpha\n"), 0o644); err != nil {
		t.Fatalf("os.WriteFile(appPath): %v", err)
	}
	untrackedPath := filepath.Join(repo, "docs", "todo.txt")
	if err := os.MkdirAll(filepath.Dir(untrackedPath), 0o755); err != nil {
		t.Fatalf("os.MkdirAll(untrackedPath): %v", err)
	}
	if err := os.WriteFile(untrackedPath, []byte("draft\n"), 0o644); err != nil {
		t.Fatalf("os.WriteFile(untrackedPath): %v", err)
	}

	status, err := Diff(context.Background(), repo)
	if err != nil {
		t.Fatalf("Diff(): %v", err)
	}

	if got, want := evalPath(t, status.RepoRoot), evalPath(t, repo); got != want {
		t.Fatalf("RepoRoot = %q, want %q", got, want)
	}
	if got, want := status.Summary.FilesChanged, 3; got != want {
		t.Fatalf("Summary.FilesChanged = %d, want %d", got, want)
	}
	if got, want := status.Summary.Insertions, 1; got != want {
		t.Fatalf("Summary.Insertions = %d, want %d", got, want)
	}
	if got, want := status.Summary.Deletions, 1; got != want {
		t.Fatalf("Summary.Deletions = %d, want %d", got, want)
	}
	if got, want := len(status.Files), 3; got != want {
		t.Fatalf("len(Files) = %d, want %d", got, want)
	}

	byPath := map[string]DiffFile{}
	for _, file := range status.Files {
		byPath[file.Path] = file
	}

	readme, ok := byPath["README.md"]
	if !ok {
		t.Fatalf("README.md not found in Files: %#v", status.Files)
	}
	if readme.Added != 1 || readme.Deleted != 0 || readme.Binary {
		t.Fatalf("README.md diff = %#v, want Added=1 Deleted=0 Binary=false", readme)
	}
	if !readme.Viewable {
		t.Fatalf("README.md Viewable = false, want true")
	}

	app, ok := byPath["pkg/app.txt"]
	if !ok {
		t.Fatalf("pkg/app.txt not found in Files: %#v", status.Files)
	}
	if app.Added != 0 || app.Deleted != 1 || app.Binary {
		t.Fatalf("pkg/app.txt diff = %#v, want Added=0 Deleted=1 Binary=false", app)
	}
	if !app.Viewable {
		t.Fatalf("pkg/app.txt Viewable = false, want true")
	}

	untracked, ok := byPath["docs/todo.txt"]
	if !ok {
		t.Fatalf("docs/todo.txt not found in Files: %#v", status.Files)
	}
	if untracked.Added != 0 || untracked.Deleted != 0 || untracked.Binary || !untracked.Untracked {
		t.Fatalf("docs/todo.txt diff = %#v, want Added=0 Deleted=0 Binary=false Untracked=true", untracked)
	}
	if !untracked.Viewable {
		t.Fatalf("docs/todo.txt Viewable = false, want true")
	}
}

func TestDiffCleanRepository(t *testing.T) {
	repo := newGitRepo(t)

	status, err := Diff(context.Background(), repo)
	if err != nil {
		t.Fatalf("Diff(): %v", err)
	}

	if status.Summary.FilesChanged != 0 || status.Summary.Insertions != 0 || status.Summary.Deletions != 0 {
		t.Fatalf("Summary = %#v, want zero counts", status.Summary)
	}
	if len(status.Files) != 0 {
		t.Fatalf("len(Files) = %d, want 0", len(status.Files))
	}
}

func TestDiffNotRepository(t *testing.T) {
	ensureGitAvailable(t)

	_, err := Diff(context.Background(), t.TempDir())
	if !errors.Is(err, ErrNotRepository) {
		t.Fatalf("Diff() error = %v, want %v", err, ErrNotRepository)
	}
}

func TestDiffMarksUntrackedBinaryFilesAsNotViewable(t *testing.T) {
	repo := newGitRepo(t)

	binaryPath := filepath.Join(repo, "assets", "logo.png")
	if err := os.MkdirAll(filepath.Dir(binaryPath), 0o755); err != nil {
		t.Fatalf("os.MkdirAll(): %v", err)
	}
	if err := os.WriteFile(binaryPath, []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a, 0x00}, 0o644); err != nil {
		t.Fatalf("os.WriteFile(binaryPath): %v", err)
	}

	status, err := Diff(context.Background(), repo)
	if err != nil {
		t.Fatalf("Diff(): %v", err)
	}

	asset, ok := mapDiffFilesByPath(status.Files)["assets/logo.png"]
	if !ok {
		t.Fatalf("assets/logo.png not found in Files: %#v", status.Files)
	}
	if !asset.Untracked || !asset.Binary || asset.Viewable {
		t.Fatalf("assets/logo.png diff = %#v, want Untracked=true Binary=true Viewable=false", asset)
	}
}

func TestFileDetailReturnsPatchForTrackedFile(t *testing.T) {
	repo := newGitRepo(t)

	readmePath := filepath.Join(repo, "README.md")
	if err := os.WriteFile(readmePath, []byte("hello\nworld\n"), 0o644); err != nil {
		t.Fatalf("os.WriteFile(README.md): %v", err)
	}

	detail, err := FileDetail(context.Background(), repo, "README.md")
	if err != nil {
		t.Fatalf("FileDetail(): %v", err)
	}

	if !detail.Supported {
		t.Fatalf("Supported = false, want true: %#v", detail)
	}
	if got, want := detail.Kind, DiffFileDetailKindDiff; got != want {
		t.Fatalf("Kind = %q, want %q", got, want)
	}
	if !strings.Contains(detail.Content, "diff --git a/README.md b/README.md") {
		t.Fatalf("Content = %q, want git diff header", detail.Content)
	}
	if !strings.Contains(detail.Content, "+world") {
		t.Fatalf("Content = %q, want added line", detail.Content)
	}
}

func TestFileDetailReturnsFileContentsForUntrackedTextFile(t *testing.T) {
	repo := newGitRepo(t)

	draftPath := filepath.Join(repo, "docs", "todo.txt")
	if err := os.MkdirAll(filepath.Dir(draftPath), 0o755); err != nil {
		t.Fatalf("os.MkdirAll(): %v", err)
	}
	if err := os.WriteFile(draftPath, []byte("draft\nnext\n"), 0o644); err != nil {
		t.Fatalf("os.WriteFile(draftPath): %v", err)
	}

	detail, err := FileDetail(context.Background(), repo, "docs/todo.txt")
	if err != nil {
		t.Fatalf("FileDetail(): %v", err)
	}

	if !detail.Supported {
		t.Fatalf("Supported = false, want true: %#v", detail)
	}
	if got, want := detail.Kind, DiffFileDetailKindFile; got != want {
		t.Fatalf("Kind = %q, want %q", got, want)
	}
	if got, want := detail.Content, "draft\nnext\n"; got != want {
		t.Fatalf("Content = %q, want %q", got, want)
	}
}

func TestFileDetailRejectsBinaryUntrackedFile(t *testing.T) {
	repo := newGitRepo(t)

	binaryPath := filepath.Join(repo, "assets", "logo.png")
	if err := os.MkdirAll(filepath.Dir(binaryPath), 0o755); err != nil {
		t.Fatalf("os.MkdirAll(): %v", err)
	}
	if err := os.WriteFile(binaryPath, []byte{0x89, 'P', 'N', 'G', 0x00}, 0o644); err != nil {
		t.Fatalf("os.WriteFile(binaryPath): %v", err)
	}

	detail, err := FileDetail(context.Background(), repo, "assets/logo.png")
	if err != nil {
		t.Fatalf("FileDetail(): %v", err)
	}

	if detail.Supported {
		t.Fatalf("Supported = true, want false: %#v", detail)
	}
	if got, want := detail.Reason, DiffFileDetailReasonNonText; got != want {
		t.Fatalf("Reason = %q, want %q", got, want)
	}
}

func TestFileDetailRejectsUnsafePath(t *testing.T) {
	repo := newGitRepo(t)

	if _, err := FileDetail(context.Background(), repo, "../README.md"); err == nil {
		t.Fatal("FileDetail() error = nil, want non-nil")
	}
}

func TestFileDetailReturnsNotFoundWhenPathIsNotInDiff(t *testing.T) {
	repo := newGitRepo(t)

	_, err := FileDetail(context.Background(), repo, "README.md")
	if !errors.Is(err, ErrDiffFileNotFound) {
		t.Fatalf("FileDetail() error = %v, want %v", err, ErrDiffFileNotFound)
	}
}

func newGitRepo(t *testing.T) string {
	t.Helper()

	ensureGitAvailable(t)

	repo := t.TempDir()
	runGitCommand(t, repo, "init", "--quiet")
	runGitCommand(t, repo, "config", "user.name", "Ngent Test")
	runGitCommand(t, repo, "config", "user.email", "ngent-test@example.com")
	runGitCommand(t, repo, "checkout", "--quiet", "-b", "main")

	if err := os.WriteFile(filepath.Join(repo, "README.md"), []byte("hello\n"), 0o644); err != nil {
		t.Fatalf("os.WriteFile(): %v", err)
	}
	runGitCommand(t, repo, "add", "README.md")
	runGitCommand(t, repo, "commit", "--quiet", "-m", "initial")
	runGitCommand(t, repo, "branch", "feature/demo")

	nested := filepath.Join(repo, "nested")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("os.MkdirAll(): %v", err)
	}
	return repo
}

func ensureGitAvailable(t *testing.T) {
	t.Helper()

	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git binary is not available")
	}
}

func runGitCommand(t *testing.T, repo string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = repo
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, output)
	}
}

func evalPath(t *testing.T, path string) string {
	t.Helper()

	resolved, err := filepath.EvalSymlinks(path)
	if err == nil {
		return resolved
	}
	return filepath.Clean(path)
}

func mapDiffFilesByPath(files []DiffFile) map[string]DiffFile {
	result := make(map[string]DiffFile, len(files))
	for _, file := range files {
		result[file.Path] = file
	}
	return result
}
