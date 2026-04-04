package httpapi

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/beyond5959/ngent/internal/gitutil"
	"github.com/beyond5959/ngent/internal/runtime"
	"github.com/beyond5959/ngent/internal/storage"
)

type threadGitResponse struct {
	ThreadID      string            `json:"threadId"`
	Available     bool              `json:"available"`
	RepoRoot      string            `json:"repoRoot,omitempty"`
	CurrentRef    string            `json:"currentRef,omitempty"`
	CurrentBranch string            `json:"currentBranch,omitempty"`
	Detached      bool              `json:"detached,omitempty"`
	Branches      []threadGitBranch `json:"branches,omitempty"`
}

type threadGitBranch struct {
	Name    string `json:"name"`
	Current bool   `json:"current"`
}

type threadGitDiffResponse struct {
	ThreadID  string                 `json:"threadId"`
	Available bool                   `json:"available"`
	RepoRoot  string                 `json:"repoRoot,omitempty"`
	Summary   threadGitDiffSummary   `json:"summary"`
	Files     []threadGitDiffFileRow `json:"files,omitempty"`
}

type threadGitDiffSummary struct {
	FilesChanged int `json:"filesChanged"`
	Insertions   int `json:"insertions"`
	Deletions    int `json:"deletions"`
}

type threadGitDiffFileRow struct {
	Path      string `json:"path"`
	Added     int    `json:"added"`
	Deleted   int    `json:"deleted"`
	Binary    bool   `json:"binary,omitempty"`
	Untracked bool   `json:"untracked,omitempty"`
	Viewable  bool   `json:"viewable"`
}

type threadGitDiffFileResponse struct {
	ThreadID  string                     `json:"threadId"`
	Available bool                       `json:"available"`
	RepoRoot  string                     `json:"repoRoot,omitempty"`
	Path      string                     `json:"path,omitempty"`
	Supported bool                       `json:"supported"`
	Kind      string                     `json:"kind,omitempty"`
	Blocks    []threadGitDiffRenderBlock `json:"blocks,omitempty"`
	Reason    string                     `json:"reason,omitempty"`
}

type threadGitDiffRenderBlock struct {
	Tone           string   `json:"tone"`
	Text           []string `json:"text"`
	OldLineNumbers []int    `json:"oldLineNumbers,omitempty"`
	NewLineNumbers []int    `json:"newLineNumbers,omitempty"`
}

func (s *Server) handleThreadGit(w http.ResponseWriter, r *http.Request, clientID, threadID string) {
	thread, ok := s.loadAccessibleThreadOrWriteNotFound(w, r.Context(), threadID)
	if !ok {
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.handleGetThreadGit(w, r, thread)
	case http.MethodPost:
		s.handleSwitchThreadGitBranch(w, r, thread)
	default:
		writeMethodNotAllowed(w, r)
	}
}

func (s *Server) handleThreadGitDiff(w http.ResponseWriter, r *http.Request, threadID string) {
	thread, ok := s.loadAccessibleThreadOrWriteNotFound(w, r.Context(), threadID)
	if !ok {
		return
	}

	if err := requireMethod(r, http.MethodGet); err != nil {
		writeMethodNotAllowed(w, r)
		return
	}

	status, err := gitutil.Diff(r.Context(), thread.CWD)
	if err != nil {
		if errors.Is(err, gitutil.ErrGitUnavailable) || errors.Is(err, gitutil.ErrNotRepository) {
			writeJSON(w, http.StatusOK, unavailableThreadGitDiffResponse(thread.ThreadID))
			return
		}
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to inspect git diff", map[string]any{
			"threadId": thread.ThreadID,
			"cwd":      thread.CWD,
			"reason":   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, threadGitDiffResponseForStatus(thread.ThreadID, status))
}

func (s *Server) handleThreadGitDiffFile(w http.ResponseWriter, r *http.Request, threadID string) {
	thread, ok := s.loadAccessibleThreadOrWriteNotFound(w, r.Context(), threadID)
	if !ok {
		return
	}

	if err := requireMethod(r, http.MethodGet); err != nil {
		writeMethodNotAllowed(w, r)
		return
	}

	rawPath := strings.TrimSpace(r.URL.Query().Get("path"))
	path, err := sanitizeThreadGitDiffFilePath(rawPath)
	if err != nil {
		writeError(w, http.StatusBadRequest, codeInvalidArgument, "path must be a repository-relative file path", map[string]any{
			"field": "path",
		})
		return
	}

	detail, err := gitutil.FileDetail(r.Context(), thread.CWD, path)
	if err != nil {
		switch {
		case errors.Is(err, gitutil.ErrGitUnavailable), errors.Is(err, gitutil.ErrNotRepository):
			writeJSON(w, http.StatusOK, unavailableThreadGitDiffFileResponse(thread.ThreadID))
			return
		case errors.Is(err, gitutil.ErrDiffFileNotFound), errors.Is(err, os.ErrNotExist):
			writeError(w, http.StatusNotFound, codeNotFound, "git diff file not found", map[string]any{
				"threadId": thread.ThreadID,
				"path":     path,
			})
			return
		default:
			writeError(w, http.StatusInternalServerError, codeInternal, "failed to inspect git diff file", map[string]any{
				"threadId": thread.ThreadID,
				"cwd":      thread.CWD,
				"path":     path,
				"reason":   err.Error(),
			})
			return
		}
	}

	writeJSON(w, http.StatusOK, threadGitDiffFileResponseForDetail(thread.ThreadID, detail))
}

func (s *Server) handleGetThreadGit(w http.ResponseWriter, r *http.Request, thread storage.Thread) {
	status, err := gitutil.Inspect(r.Context(), thread.CWD)
	if err != nil {
		if errors.Is(err, gitutil.ErrGitUnavailable) || errors.Is(err, gitutil.ErrNotRepository) {
			writeJSON(w, http.StatusOK, unavailableThreadGitResponse(thread.ThreadID))
			return
		}
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to inspect git state", map[string]any{
			"threadId": thread.ThreadID,
			"cwd":      thread.CWD,
			"reason":   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, threadGitResponseForStatus(thread.ThreadID, status))
}

func (s *Server) handleSwitchThreadGitBranch(w http.ResponseWriter, r *http.Request, thread storage.Thread) {
	var req struct {
		Branch string `json:"branch"`
	}

	if err := decodeJSONBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, codeInvalidArgument, "invalid JSON body", map[string]any{"reason": err.Error()})
		return
	}

	req.Branch = strings.TrimSpace(req.Branch)
	if req.Branch == "" {
		writeError(w, http.StatusBadRequest, codeInvalidArgument, "branch is required", map[string]any{"field": "branch"})
		return
	}

	guardTurnID := "git-checkout-" + newTurnID()
	switchCtx, cancelSwitch := context.WithCancel(r.Context())
	if err := s.turns.ActivateThreadExclusive(thread.ThreadID, guardTurnID, cancelSwitch); err != nil {
		cancelSwitch()
		if errors.Is(err, runtime.ErrActiveTurnExists) {
			writeError(w, http.StatusConflict, codeConflict, "thread has an active turn", map[string]any{"threadId": thread.ThreadID})
			return
		}
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to lock thread for git checkout", map[string]any{
			"threadId": thread.ThreadID,
			"reason":   err.Error(),
		})
		return
	}
	defer func() {
		cancelSwitch()
		s.turns.ReleaseThreadExclusive(thread.ThreadID, guardTurnID)
	}()

	status, err := gitutil.Checkout(switchCtx, thread.CWD, req.Branch)
	if err != nil {
		switch {
		case errors.Is(err, gitutil.ErrGitUnavailable), errors.Is(err, gitutil.ErrNotRepository):
			writeJSON(w, http.StatusOK, unavailableThreadGitResponse(thread.ThreadID))
			return
		case errors.Is(err, gitutil.ErrBranchNotFound):
			writeError(w, http.StatusBadRequest, codeInvalidArgument, "git branch does not exist locally", map[string]any{
				"field":    "branch",
				"branch":   req.Branch,
				"threadId": thread.ThreadID,
			})
			return
		default:
			writeError(w, http.StatusConflict, codeConflict, "failed to switch git branch", map[string]any{
				"threadId": thread.ThreadID,
				"branch":   req.Branch,
				"reason":   err.Error(),
			})
			return
		}
	}

	writeJSON(w, http.StatusOK, threadGitResponseForStatus(thread.ThreadID, status))
}

func unavailableThreadGitResponse(threadID string) threadGitResponse {
	return threadGitResponse{
		ThreadID:  threadID,
		Available: false,
	}
}

func unavailableThreadGitDiffResponse(threadID string) threadGitDiffResponse {
	return threadGitDiffResponse{
		ThreadID:  threadID,
		Available: false,
	}
}

func unavailableThreadGitDiffFileResponse(threadID string) threadGitDiffFileResponse {
	return threadGitDiffFileResponse{
		ThreadID:  threadID,
		Available: false,
	}
}

func threadGitResponseForStatus(threadID string, status gitutil.Status) threadGitResponse {
	branches := make([]threadGitBranch, 0, len(status.Branches))
	for _, branch := range status.Branches {
		name := strings.TrimSpace(branch.Name)
		if name == "" {
			continue
		}
		branches = append(branches, threadGitBranch{
			Name:    name,
			Current: branch.Current,
		})
	}

	return threadGitResponse{
		ThreadID:      threadID,
		Available:     true,
		RepoRoot:      status.RepoRoot,
		CurrentRef:    status.CurrentRef,
		CurrentBranch: status.CurrentBranch,
		Detached:      status.Detached,
		Branches:      branches,
	}
}

func threadGitDiffResponseForStatus(threadID string, status gitutil.DiffStatus) threadGitDiffResponse {
	files := make([]threadGitDiffFileRow, 0, len(status.Files))
	for _, file := range status.Files {
		path := strings.TrimSpace(file.Path)
		if path == "" {
			continue
		}
		files = append(files, threadGitDiffFileRow{
			Path:      path,
			Added:     file.Added,
			Deleted:   file.Deleted,
			Binary:    file.Binary,
			Untracked: file.Untracked,
			Viewable:  file.Viewable,
		})
	}

	return threadGitDiffResponse{
		ThreadID:  threadID,
		Available: true,
		RepoRoot:  status.RepoRoot,
		Summary: threadGitDiffSummary{
			FilesChanged: status.Summary.FilesChanged,
			Insertions:   status.Summary.Insertions,
			Deletions:    status.Summary.Deletions,
		},
		Files: files,
	}
}

func threadGitDiffFileResponseForDetail(threadID string, detail gitutil.DiffFileDetail) threadGitDiffFileResponse {
	blocks := make([]threadGitDiffRenderBlock, 0, len(detail.Blocks))
	for _, block := range detail.Blocks {
		blocks = append(blocks, threadGitDiffRenderBlock{
			Tone:           block.Tone,
			Text:           append([]string(nil), block.Text...),
			OldLineNumbers: append([]int(nil), block.OldLineNumbers...),
			NewLineNumbers: append([]int(nil), block.NewLineNumbers...),
		})
	}
	return threadGitDiffFileResponse{
		ThreadID:  threadID,
		Available: true,
		RepoRoot:  detail.RepoRoot,
		Path:      detail.Path,
		Supported: detail.Supported,
		Kind:      detail.Kind,
		Blocks:    blocks,
		Reason:    detail.Reason,
	}
}

func sanitizeThreadGitDiffFilePath(rawPath string) (string, error) {
	trimmed := strings.TrimSpace(rawPath)
	if trimmed == "" {
		return "", errors.New("path is required")
	}

	cleaned := filepath.Clean(filepath.FromSlash(trimmed))
	if cleaned == "." || cleaned == string(filepath.Separator) || filepath.IsAbs(cleaned) {
		return "", errors.New("path must be relative")
	}
	if cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return "", errors.New("path escapes repository root")
	}
	return filepath.ToSlash(cleaned), nil
}
