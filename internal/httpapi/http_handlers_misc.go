package httpapi

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/beyond5959/ngent/internal/storage"
)

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, r)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleAttachment(w http.ResponseWriter, r *http.Request, attachmentID string) {
	if err := requireMethod(r, http.MethodGet); err != nil {
		writeMethodNotAllowed(w, r)
		return
	}
	if !s.isAttachmentAuthorized(r) {
		writeError(w, http.StatusUnauthorized, codeUnauthorized, "missing or invalid attachment token", map[string]any{})
		return
	}

	if s.store == nil {
		writeError(w, http.StatusInternalServerError, codeInternal, "storage is not configured", map[string]any{})
		return
	}

	attachment, err := s.store.GetTurnAttachment(r.Context(), attachmentID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			writeError(w, http.StatusNotFound, codeNotFound, "attachment not found", map[string]any{})
			return
		}
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to load attachment", map[string]any{
			"reason": err.Error(),
		})
		return
	}

	if _, err := s.store.GetTurn(r.Context(), attachment.TurnID); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			writeError(w, http.StatusNotFound, codeNotFound, "attachment not found", map[string]any{})
			return
		}
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to load attachment turn", map[string]any{
			"reason": err.Error(),
		})
		return
	}
	attachmentPath := filepath.Clean(strings.TrimSpace(attachment.FilePath))
	if attachmentPath == "" || !isPathAllowed(attachmentPath, []string{s.dataDir}) {
		writeError(w, http.StatusInternalServerError, codeInternal, "attachment path is invalid", map[string]any{})
		return
	}

	file, err := os.Open(attachmentPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			writeError(w, http.StatusNotFound, codeNotFound, "attachment file not found", map[string]any{})
			return
		}
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to open attachment file", map[string]any{
			"reason": err.Error(),
		})
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to stat attachment file", map[string]any{
			"reason": err.Error(),
		})
		return
	}
	if info.IsDir() {
		writeError(w, http.StatusNotFound, codeNotFound, "attachment file not found", map[string]any{})
		return
	}

	contentType := strings.TrimSpace(attachment.MimeType)
	if contentType == "" {
		contentType = mime.TypeByExtension(strings.ToLower(filepath.Ext(attachment.Name)))
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "private, max-age=31536000, immutable")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", normalizeUploadFilename(attachment.Name)))
	http.ServeContent(w, r, attachment.Name, info.ModTime(), file)
}

func (s *Server) handleAgents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, r)
		return
	}

	writeJSON(w, http.StatusOK, struct {
		Agents []AgentInfo `json:"agents"`
	}{Agents: s.agents})
}

func (s *Server) handleThreadsCollection(w http.ResponseWriter, r *http.Request, clientID string) {
	switch r.Method {
	case http.MethodPost:
		s.handleCreateThread(w, r, clientID)
	case http.MethodGet:
		s.handleListThreads(w, r, clientID)
	default:
		writeMethodNotAllowed(w, r)
	}
}

func (s *Server) handleThreadResource(w http.ResponseWriter, r *http.Request, clientID, threadID, subresource string) {
	switch subresource {
	case "":
		switch r.Method {
		case http.MethodGet:
			s.handleGetThread(w, r, clientID, threadID)
		case http.MethodPatch:
			s.handleUpdateThread(w, r, clientID, threadID)
		case http.MethodDelete:
			s.handleDeleteThread(w, r, clientID, threadID)
		default:
			writeMethodNotAllowed(w, r)
		}
	case "turns":
		s.handleCreateTurnStream(w, r, clientID, threadID)
	case "compact":
		s.handleCompactThread(w, r, clientID, threadID)
	case "history":
		s.handleThreadHistory(w, r, clientID, threadID)
	case "sessions":
		s.handleThreadSessions(w, r, clientID, threadID)
	case "session-usage":
		s.handleThreadSessionUsage(w, r, clientID, threadID)
	case "config-options":
		s.handleThreadConfigOptions(w, r, clientID, threadID)
	case "slash-commands":
		s.handleThreadSlashCommands(w, r, clientID, threadID)
	case "git":
		s.handleThreadGit(w, r, clientID, threadID)
	case "git-diff":
		s.handleThreadGitDiff(w, r, threadID)
	case "git-diff-file":
		s.handleThreadGitDiffFile(w, r, threadID)
	default:
		writeError(w, http.StatusNotFound, codeNotFound, "endpoint not found", map[string]any{"path": r.URL.Path})
	}
}

// handlePathSearch handles path search requests for the working directory input.
// It searches for directories under $HOME matching the query.
// Search is triggered only when query has 3 or more characters.
// Priority: first level, then second level, then third level.
func (s *Server) handlePathSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, r)
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if len(query) < 3 {
		writeJSON(w, http.StatusOK, map[string]any{
			"query":   query,
			"results": []string{},
		})
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		s.logger.Warn("path_search.home_dir_failed", "error", err.Error())
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to get home directory", nil)
		return
	}

	results := s.searchDirectories(homeDir, query)

	writeJSON(w, http.StatusOK, map[string]any{
		"query":   query,
		"results": results,
	})
}

// searchDirectories searches for directories matching the query.
// Searches all levels and returns up to maxPathSearchResults matches.
// Skips hidden directories (those starting with a dot).
const maxPathSearchResults = 5

func (s *Server) searchDirectories(homeDir, query string) []string {
	queryLower := strings.ToLower(query)
	var results []string

	entries, err := os.ReadDir(homeDir)
	if err != nil {
		s.logger.Warn("path_search.read_dir_failed", "dir", homeDir, "error", err.Error())
		return results
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		if strings.Contains(strings.ToLower(name), queryLower) {
			results = append(results, filepath.Join(homeDir, name))
			if len(results) >= maxPathSearchResults {
				return results
			}
		}
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		firstLevelPath := filepath.Join(homeDir, name)

		subEntries, err := os.ReadDir(firstLevelPath)
		if err != nil {
			continue
		}

		for _, subEntry := range subEntries {
			if !subEntry.IsDir() {
				continue
			}
			subName := subEntry.Name()
			if strings.HasPrefix(subName, ".") {
				continue
			}
			if strings.Contains(strings.ToLower(subName), queryLower) {
				results = append(results, filepath.Join(firstLevelPath, subName))
				if len(results) >= maxPathSearchResults {
					return results
				}
			}
		}
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		firstLevelPath := filepath.Join(homeDir, name)

		subEntries, err := os.ReadDir(firstLevelPath)
		if err != nil {
			continue
		}

		for _, subEntry := range subEntries {
			if !subEntry.IsDir() {
				continue
			}
			subName := subEntry.Name()
			if strings.HasPrefix(subName, ".") {
				continue
			}
			secondLevelPath := filepath.Join(firstLevelPath, subName)

			thirdEntries, err := os.ReadDir(secondLevelPath)
			if err != nil {
				continue
			}

			for _, thirdEntry := range thirdEntries {
				if !thirdEntry.IsDir() {
					continue
				}
				thirdName := thirdEntry.Name()
				if strings.HasPrefix(thirdName, ".") {
					continue
				}
				if strings.Contains(strings.ToLower(thirdName), queryLower) {
					results = append(results, filepath.Join(secondLevelPath, thirdName))
					if len(results) >= maxPathSearchResults {
						return results
					}
				}
			}
		}
	}

	return results
}

func (s *Server) handleRecentDirectories(w http.ResponseWriter, r *http.Request, clientID string) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, r)
		return
	}

	dirs, err := s.store.ListRecentDirectories(r.Context(), clientID, 5)
	if err != nil {
		s.logger.Warn("recent_directories.query_failed", "error", err.Error())
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to get recent directories", nil)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"directories": dirs,
	})
}
