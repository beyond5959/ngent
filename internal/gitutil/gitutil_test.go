package gitutil

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
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
