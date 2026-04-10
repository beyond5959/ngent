package httpapi

import (
	"bufio"
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/beyond5959/ngent/internal/storage"
)

const threadFilePreviewMaxLines = 10000

const (
	threadFilePreviewKindText  = "text"
	threadFilePreviewKindImage = "image"
)

var errThreadFilePreviewForbidden = errors.New("file preview path is outside allowed roots")

type threadFilePreviewResponse struct {
	ThreadID  string                     `json:"threadId"`
	Path      string                     `json:"path,omitempty"`
	Supported bool                       `json:"supported"`
	Kind      string                     `json:"kind,omitempty"`
	MimeType  string                     `json:"mimeType,omitempty"`
	StartLine int                        `json:"startLine,omitempty"`
	EndLine   int                        `json:"endLine,omitempty"`
	FocusLine int                        `json:"focusLine,omitempty"`
	Blocks    []threadGitDiffRenderBlock `json:"blocks,omitempty"`
	Reason    string                     `json:"reason,omitempty"`
}

type threadPreviewFileInfo struct {
	ResolvedPath string
	MimeType     string
	Kind         string
	Supported    bool
	Reason       string
}

type threadFilePreviewTextContent struct {
	StartLine int
	EndLine   int
	FocusLine int
	Blocks    []threadGitDiffRenderBlock
}

func (s *Server) handleThreadFilePreview(w http.ResponseWriter, r *http.Request, threadID string) {
	thread, ok := s.loadAccessibleThreadOrWriteNotFound(w, r.Context(), threadID)
	if !ok {
		return
	}
	if err := requireMethod(r, http.MethodGet); err != nil {
		writeMethodNotAllowed(w, r)
		return
	}

	rawPath := strings.TrimSpace(r.URL.Query().Get("path"))
	resolvedPath, err := resolveThreadPreviewPath(rawPath, s.allowedRoots)
	if err != nil {
		writeThreadFilePreviewPathError(w, err)
		return
	}

	info, err := inspectThreadPreviewFile(resolvedPath)
	if err != nil {
		writeThreadFilePreviewInspectError(w, thread, rawPath, err)
		return
	}

	response := threadFilePreviewResponse{
		ThreadID:  thread.ThreadID,
		Path:      resolvedPath,
		Supported: info.Supported,
		Kind:      info.Kind,
		MimeType:  info.MimeType,
		Reason:    info.Reason,
	}
	if !info.Supported {
		writeJSON(w, http.StatusOK, response)
		return
	}
	if info.Kind == threadFilePreviewKindImage {
		writeJSON(w, http.StatusOK, response)
		return
	}

	focusLine, err := parseThreadFilePreviewLineQuery(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, codeInvalidArgument, "line must be a positive integer", map[string]any{
			"line": r.URL.Query().Get("line"),
		})
		return
	}

	textContent, err := loadThreadFilePreviewTextContent(resolvedPath, focusLine)
	if err != nil {
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to read preview file", map[string]any{
			"threadId": thread.ThreadID,
			"path":     rawPath,
			"reason":   err.Error(),
		})
		return
	}

	response.StartLine = textContent.StartLine
	response.EndLine = textContent.EndLine
	response.FocusLine = textContent.FocusLine
	response.Blocks = textContent.Blocks
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) handleThreadFilePreviewContent(w http.ResponseWriter, r *http.Request, threadID string) {
	thread, ok := s.loadAccessibleThreadOrWriteNotFound(w, r.Context(), threadID)
	if !ok {
		return
	}
	if err := requireMethod(r, http.MethodGet); err != nil {
		writeMethodNotAllowed(w, r)
		return
	}

	rawPath := strings.TrimSpace(r.URL.Query().Get("path"))
	resolvedPath, err := resolveThreadPreviewPath(rawPath, s.allowedRoots)
	if err != nil {
		writeThreadFilePreviewPathError(w, err)
		return
	}

	info, err := inspectThreadPreviewFile(resolvedPath)
	if err != nil {
		writeThreadFilePreviewInspectError(w, thread, rawPath, err)
		return
	}
	if !info.Supported || info.Kind != threadFilePreviewKindImage {
		writeError(w, http.StatusBadRequest, codeInvalidArgument, "path does not point to a previewable image file", map[string]any{
			"path": resolvedPath,
		})
		return
	}

	file, err := os.Open(resolvedPath)
	if err != nil {
		writeThreadFilePreviewInspectError(w, thread, rawPath, err)
		return
	}
	defer file.Close()

	infoStat, err := file.Stat()
	if err != nil {
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to stat preview file", map[string]any{
			"threadId": thread.ThreadID,
			"path":     rawPath,
			"reason":   err.Error(),
		})
		return
	}
	if infoStat.IsDir() {
		writeError(w, http.StatusNotFound, codeNotFound, "preview file not found", map[string]any{
			"path": rawPath,
		})
		return
	}

	w.Header().Set("Content-Type", info.MimeType)
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("Content-Disposition", "inline")
	http.ServeContent(w, r, filepath.Base(resolvedPath), infoStat.ModTime(), file)
}

func resolveThreadPreviewPath(rawPath string, roots []string) (string, error) {
	cleaned := filepath.Clean(strings.TrimSpace(rawPath))
	if cleaned == "" || !filepath.IsAbs(cleaned) {
		return "", errors.New("path must be an absolute path")
	}

	resolved, err := filepath.EvalSymlinks(cleaned)
	if err != nil {
		return "", err
	}
	resolved = filepath.Clean(resolved)
	allowedRoots := make([]string, 0, len(roots)*2)
	for _, root := range roots {
		trimmedRoot := filepath.Clean(strings.TrimSpace(root))
		if trimmedRoot == "" {
			continue
		}
		allowedRoots = append(allowedRoots, trimmedRoot)
		if evaluatedRoot, evalErr := filepath.EvalSymlinks(trimmedRoot); evalErr == nil {
			allowedRoots = append(allowedRoots, filepath.Clean(evaluatedRoot))
		}
	}
	if !isPathAllowed(resolved, allowedRoots) {
		return "", errThreadFilePreviewForbidden
	}
	return resolved, nil
}

func inspectThreadPreviewFile(path string) (threadPreviewFileInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return threadPreviewFileInfo{}, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return threadPreviewFileInfo{}, err
	}
	if info.IsDir() {
		return threadPreviewFileInfo{}, os.ErrNotExist
	}

	sniff := make([]byte, 8192)
	n, readErr := file.Read(sniff)
	if readErr != nil && !errors.Is(readErr, io.EOF) {
		return threadPreviewFileInfo{}, readErr
	}
	sniff = sniff[:n]

	extMime := strings.TrimSpace(mime.TypeByExtension(strings.ToLower(filepath.Ext(path))))
	detectedMime := strings.TrimSpace(http.DetectContentType(sniff))
	if previewImageMime(extMime) || previewImageMime(detectedMime) {
		return threadPreviewFileInfo{
			ResolvedPath: path,
			MimeType:     choosePreviewMime(extMime, detectedMime, "image/png"),
			Kind:         threadFilePreviewKindImage,
			Supported:    true,
		}, nil
	}
	if previewTextMime(extMime) || previewTextMime(detectedMime) {
		return threadPreviewFileInfo{
			ResolvedPath: path,
			MimeType:     choosePreviewMime(extMime, detectedMime, "text/plain; charset=utf-8"),
			Kind:         threadFilePreviewKindText,
			Supported:    true,
		}, nil
	}
	if previewExplicitNonTextMime(extMime) || previewExplicitNonTextMime(detectedMime) {
		return threadPreviewFileInfo{
			ResolvedPath: path,
			MimeType:     choosePreviewMime(extMime, detectedMime, "application/octet-stream"),
			Supported:    false,
			Reason:       "non_text",
		}, nil
	}
	if previewLikelyText(sniff) {
		return threadPreviewFileInfo{
			ResolvedPath: path,
			MimeType:     choosePreviewMime(extMime, detectedMime, "text/plain; charset=utf-8"),
			Kind:         threadFilePreviewKindText,
			Supported:    true,
		}, nil
	}

	return threadPreviewFileInfo{
		ResolvedPath: path,
		MimeType:     choosePreviewMime(extMime, detectedMime, "application/octet-stream"),
		Supported:    false,
		Reason:       "non_text",
	}, nil
}

func parseThreadFilePreviewLineQuery(r *http.Request) (int, error) {
	if rawLine := strings.TrimSpace(r.URL.Query().Get("line")); rawLine != "" {
		return parsePositiveInt(rawLine)
	}
	return 0, nil
}

func parsePositiveInt(value string) (int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, errors.New("value is required")
	}

	total := 0
	for _, ch := range value {
		if ch < '0' || ch > '9' {
			return 0, errors.New("value must be numeric")
		}
		total = total*10 + int(ch-'0')
	}
	if total <= 0 {
		return 0, errors.New("value must be positive")
	}
	return total, nil
}

func loadThreadFilePreviewTextContent(path string, focusLine int) (threadFilePreviewTextContent, error) {
	file, err := os.Open(path)
	if err != nil {
		return threadFilePreviewTextContent{}, err
	}
	defer file.Close()

	scanner := newThreadPreviewLineScanner(file)
	lines := make([]string, 0, min(threadFilePreviewMaxLines, 1024))
	for scanner.Scan() {
		if len(lines) >= threadFilePreviewMaxLines {
			break
		}
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return threadFilePreviewTextContent{}, err
	}

	if len(lines) == 0 {
		return threadFilePreviewTextContent{
			Blocks: []threadGitDiffRenderBlock{{
				Tone: "plain",
				Text: []string{""},
			}},
		}, nil
	}

	lineNumbers := make([]int, 0, len(lines))
	for offset := range lines {
		lineNumbers = append(lineNumbers, offset+1)
	}

	content := threadFilePreviewTextContent{
		StartLine: 1,
		EndLine:   len(lines),
		Blocks: []threadGitDiffRenderBlock{{
			Tone:           "plain",
			Text:           lines,
			NewLineNumbers: lineNumbers,
		}},
	}
	if focusLine >= 1 && focusLine <= len(lines) {
		content.FocusLine = focusLine
	}
	return content, nil
}

func newThreadPreviewLineScanner(src io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(src)
	scanner.Buffer(make([]byte, 0, 64*1024), 2*1024*1024)
	return scanner
}

func previewImageMime(value string) bool {
	base := previewMimeBase(value)
	return strings.HasPrefix(base, "image/")
}

func previewTextMime(value string) bool {
	base := previewMimeBase(value)
	if base == "" {
		return false
	}
	if strings.HasPrefix(base, "text/") {
		return true
	}
	if base == "application/json" || base == "application/javascript" || base == "application/x-javascript" {
		return true
	}
	return strings.HasSuffix(base, "+json") || strings.HasSuffix(base, "+xml")
}

func previewExplicitNonTextMime(value string) bool {
	base := previewMimeBase(value)
	if base == "" || base == "application/octet-stream" {
		return false
	}
	return !previewTextMime(base) && !previewImageMime(base)
}

func previewMimeBase(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return ""
	}
	if index := strings.Index(value, ";"); index >= 0 {
		value = strings.TrimSpace(value[:index])
	}
	return value
}

func choosePreviewMime(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return "application/octet-stream"
}

func previewLikelyText(data []byte) bool {
	if len(data) == 0 {
		return true
	}
	for _, b := range data {
		if b == 0 {
			return false
		}
	}
	return utf8.Valid(data)
}

func writeThreadFilePreviewPathError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errThreadFilePreviewForbidden):
		writeError(w, http.StatusForbidden, codeForbidden, "path is outside allowed roots", map[string]any{
			"field": "path",
		})
	case errors.Is(err, os.ErrNotExist):
		writeError(w, http.StatusNotFound, codeNotFound, "preview file not found", map[string]any{
			"field": "path",
		})
	default:
		writeError(w, http.StatusBadRequest, codeInvalidArgument, "path must be an absolute file path inside allowed roots", map[string]any{
			"field": "path",
		})
	}
}

func writeThreadFilePreviewInspectError(w http.ResponseWriter, thread storage.Thread, rawPath string, err error) {
	switch {
	case errors.Is(err, os.ErrNotExist):
		writeError(w, http.StatusNotFound, codeNotFound, "preview file not found", map[string]any{
			"threadId": thread.ThreadID,
			"path":     rawPath,
		})
	default:
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to inspect preview file", map[string]any{
			"threadId": thread.ThreadID,
			"path":     rawPath,
			"reason":   err.Error(),
		})
	}
}
