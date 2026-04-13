package httpapi

import (
	"bufio"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/beyond5959/ngent/internal/agents"
	"github.com/beyond5959/ngent/internal/observability"
)

func parseThreadPath(path string) (threadID, subresource string, ok bool) {
	const prefix = "/v1/threads/"
	if !strings.HasPrefix(path, prefix) {
		return "", "", false
	}
	parts := strings.Split(strings.TrimPrefix(path, prefix), "/")
	if len(parts) == 0 || parts[0] == "" {
		return "", "", false
	}
	threadID = parts[0]
	if len(parts) == 1 {
		return threadID, "", true
	}
	if len(parts) == 2 && parts[1] != "" {
		return threadID, parts[1], true
	}
	return "", "", false
}

func parseAttachmentPath(path string) (attachmentID string, ok bool) {
	const prefix = "/attachments/"
	if !strings.HasPrefix(path, prefix) {
		return "", false
	}
	raw := strings.Trim(strings.TrimPrefix(path, prefix), "/")
	if raw == "" || strings.Contains(raw, "/") {
		return "", false
	}
	return raw, true
}

func parseAgentModelsPath(path string) (agentID string, ok bool) {
	const prefix = "/v1/agents/"
	const suffix = "/models"
	if !strings.HasPrefix(path, prefix) || !strings.HasSuffix(path, suffix) {
		return "", false
	}
	raw := strings.TrimSuffix(strings.TrimPrefix(path, prefix), suffix)
	raw = strings.Trim(raw, "/")
	if raw == "" || strings.Contains(raw, "/") {
		return "", false
	}
	return raw, true
}

func parsePermissionPath(path string) (permissionID string, ok bool) {
	const prefix = "/v1/permissions/"
	if !strings.HasPrefix(path, prefix) {
		return "", false
	}
	raw := strings.Trim(strings.TrimPrefix(path, prefix), "/")
	if raw == "" || strings.Contains(raw, "/") {
		return "", false
	}
	return raw, true
}

func parseTurnCancelPath(path string) (turnID string, ok bool) {
	const prefix = "/v1/turns/"
	const suffix = "/cancel"
	if !strings.HasPrefix(path, prefix) || !strings.HasSuffix(path, suffix) {
		return "", false
	}
	raw := strings.TrimSuffix(strings.TrimPrefix(path, prefix), suffix)
	raw = strings.Trim(raw, "/")
	if raw == "" || strings.Contains(raw, "/") {
		return "", false
	}
	return raw, true
}

func newThreadID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("th_%d", time.Now().UTC().UnixMicro())
	}
	return fmt.Sprintf("th_%d_%s", time.Now().UTC().UnixMicro(), hex.EncodeToString(buf))
}

func newTurnID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("tu_%d", time.Now().UTC().UnixMicro())
	}
	return fmt.Sprintf("tu_%d_%s", time.Now().UTC().UnixMicro(), hex.EncodeToString(buf))
}

func newAttachmentID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("att_%d", time.Now().UTC().UnixMicro())
	}
	return fmt.Sprintf("att_%d_%s", time.Now().UTC().UnixMicro(), hex.EncodeToString(buf))
}

func parseBoolQuery(r *http.Request, key string) bool {
	return parseBoolString(r.URL.Query().Get(key))
}

func parseFormBoolValue(value string) bool {
	return parseBoolString(value)
}

func parseBoolString(value string) bool {
	value = strings.TrimSpace(strings.ToLower(value))
	return value == "1" || value == "true" || value == "yes"
}

func requireMethod(r *http.Request, method string) error {
	if r.Method != method {
		return errors.New("method not allowed")
	}
	return nil
}

func decodeJSONBody(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return err
	}
	if decoder.More() {
		return errors.New("extra JSON values are not allowed")
	}
	return nil
}

func (s *Server) decodeTurnCreateRequest(r *http.Request) (turnCreateRequest, error) {
	contentType := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type")))
	if strings.HasPrefix(contentType, "multipart/form-data") {
		return decodeMultipartTurnCreateRequest(r, s.dataDir)
	}

	var req struct {
		Input      string `json:"input"`
		Stream     bool   `json:"stream"`
		FullAccess bool   `json:"fullAccess"`
	}
	if err := decodeJSONBody(r, &req); err != nil {
		return turnCreateRequest{}, err
	}

	return turnCreateRequest{
		Stream:     req.Stream,
		FullAccess: req.FullAccess,
		Prompt:     agents.TextPrompt(req.Input),
	}, nil
}

func decodeMultipartTurnCreateRequest(r *http.Request, dataDir string) (turnCreateRequest, error) {
	if err := r.ParseMultipartForm(maxTurnMultipartMemory); err != nil {
		return turnCreateRequest{}, err
	}
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}

	text := strings.TrimSpace(r.FormValue("input"))
	stream := parseFormBoolValue(r.FormValue("stream"))
	attachments, err := persistTurnAttachments(dataDir, r.MultipartForm.File["attachments"])
	if err != nil {
		return turnCreateRequest{}, err
	}

	content := make([]agents.PromptContent, 0, len(attachments)+1)
	if text != "" {
		content = append(content, agents.PromptContent{
			Type: agents.PromptContentTypeText,
			Text: text,
		})
	}
	for _, attachment := range attachments {
		content = append(content, attachment.PromptContent)
	}

	return turnCreateRequest{
		Stream:     stream,
		FullAccess: parseFormBoolValue(r.FormValue("fullAccess")),
		Prompt:     agents.NormalizePrompt(agents.Prompt{Content: content}),
		Uploads:    attachments,
	}, nil
}

func persistTurnAttachments(dataDir string, files []*multipart.FileHeader) ([]storedTurnAttachment, error) {
	if len(files) == 0 {
		return nil, nil
	}

	attachments := make([]storedTurnAttachment, 0, len(files))
	for _, fileHeader := range files {
		attachment, err := persistTurnAttachment(dataDir, fileHeader)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, attachment)
	}
	return attachments, nil
}

func persistTurnAttachment(dataDir string, fileHeader *multipart.FileHeader) (storedTurnAttachment, error) {
	if fileHeader == nil {
		return storedTurnAttachment{}, errors.New("attachment is required")
	}

	src, err := fileHeader.Open()
	if err != nil {
		return storedTurnAttachment{}, fmt.Errorf("open attachment %q: %w", fileHeader.Filename, err)
	}
	defer src.Close()

	attachmentID := newAttachmentID()
	displayName := normalizeUploadFilename(fileHeader.Filename)
	dstFile, dstPath, err := createUploadTempFile(dataDir, attachmentID, displayName)
	if err != nil {
		return storedTurnAttachment{}, err
	}

	size, mimeType, copyErr := copyUploadToTempFile(dstFile, src, displayName, fileHeader.Header.Get("Content-Type"))
	closeErr := dstFile.Close()
	if copyErr != nil {
		_ = os.Remove(dstPath)
		return storedTurnAttachment{}, copyErr
	}
	if closeErr != nil {
		_ = os.Remove(dstPath)
		return storedTurnAttachment{}, fmt.Errorf("close temp upload %q: %w", displayName, closeErr)
	}

	finalPath, err := finalizeUploadPath(dataDir, attachmentID, displayName, mimeType)
	if err != nil {
		_ = os.Remove(dstPath)
		return storedTurnAttachment{}, err
	}
	if err := os.Rename(dstPath, finalPath); err != nil {
		_ = os.Remove(dstPath)
		return storedTurnAttachment{}, fmt.Errorf("move stored upload %q: %w", displayName, err)
	}

	return storedTurnAttachment{
		PromptContent: agents.PromptContent{
			Type:         agents.PromptContentTypeResourceLink,
			URI:          fileURIForPath(finalPath),
			Name:         displayName,
			MimeType:     mimeType,
			Size:         size,
			AttachmentID: attachmentID,
		},
		FilePath: finalPath,
	}, nil
}

func createUploadTempFile(dataDir, attachmentID, displayName string) (*os.File, string, error) {
	displayName = normalizeUploadFilename(displayName)
	ext := filepath.Ext(displayName)
	stem := sanitizeUploadTempStem(strings.TrimSuffix(displayName, ext))
	pattern := fmt.Sprintf("%s-%s-*%s", attachmentID, stem, ext)
	tempDir := filepath.Join(filepath.Clean(dataDir), "attachments", ".incoming")
	if err := os.MkdirAll(tempDir, 0o755); err != nil {
		return nil, "", fmt.Errorf("create upload staging dir %q: %w", tempDir, err)
	}
	file, err := os.CreateTemp(tempDir, pattern)
	if err != nil {
		return nil, "", fmt.Errorf("create temp upload for %q: %w", displayName, err)
	}
	return file, file.Name(), nil
}

func copyUploadToTempFile(dst *os.File, src multipart.File, displayName, headerMime string) (int64, string, error) {
	if dst == nil {
		return 0, "", errors.New("temp upload file is required")
	}
	if src == nil {
		return 0, "", errors.New("upload source is required")
	}

	sniffBuf := make([]byte, 512)
	n, readErr := io.ReadFull(src, sniffBuf)
	switch {
	case readErr == nil:
	case errors.Is(readErr, io.EOF), errors.Is(readErr, io.ErrUnexpectedEOF):
	default:
		return 0, "", fmt.Errorf("read upload %q: %w", displayName, readErr)
	}

	total := int64(0)
	if n > 0 {
		written, err := dst.Write(sniffBuf[:n])
		total += int64(written)
		if err != nil {
			return 0, "", fmt.Errorf("write upload %q: %w", displayName, err)
		}
		if written != n {
			return 0, "", io.ErrShortWrite
		}
	}

	written, err := io.Copy(dst, src)
	total += written
	if err != nil {
		return 0, "", fmt.Errorf("copy upload %q: %w", displayName, err)
	}

	return total, detectUploadMimeType(displayName, headerMime, sniffBuf[:n]), nil
}

func detectUploadMimeType(displayName, headerMime string, sniff []byte) string {
	if mediaType, _, err := mime.ParseMediaType(strings.TrimSpace(headerMime)); err == nil && mediaType != "" && mediaType != "application/octet-stream" {
		return mediaType
	}
	if len(sniff) > 0 {
		if detected := http.DetectContentType(sniff); detected != "" {
			return detected
		}
	}
	if detected := mime.TypeByExtension(strings.ToLower(filepath.Ext(displayName))); detected != "" {
		return detected
	}
	return "application/octet-stream"
}

func normalizeUploadFilename(name string) string {
	name = strings.TrimSpace(filepath.Base(name))
	name = strings.Map(func(r rune) rune {
		switch r {
		case 0, '/', '\\':
			return -1
		default:
			return r
		}
	}, name)
	if name == "" || name == "." {
		return "attachment"
	}
	return name
}

func sanitizeUploadTempStem(stem string) string {
	stem = strings.ToLower(strings.TrimSpace(stem))
	if stem == "" {
		return "attachment"
	}
	var builder strings.Builder
	for _, r := range stem {
		switch {
		case r >= 'a' && r <= 'z':
			builder.WriteRune(r)
		case r >= '0' && r <= '9':
			builder.WriteRune(r)
		case r == '-' || r == '_':
			builder.WriteRune(r)
		default:
			builder.WriteByte('-')
		}
		if builder.Len() >= 32 {
			break
		}
	}
	result := strings.Trim(builder.String(), "-_")
	if result == "" {
		return "attachment"
	}
	return result
}

func buildStoredUploadFilename(attachmentID, displayName string) string {
	displayName = normalizeUploadFilename(displayName)
	ext := strings.ToLower(filepath.Ext(displayName))
	stem := sanitizeUploadTempStem(strings.TrimSuffix(displayName, ext))
	if ext == "" {
		return fmt.Sprintf("%s-%s", attachmentID, stem)
	}
	return fmt.Sprintf("%s-%s%s", attachmentID, stem, ext)
}

func uploadDirectoryCategory(displayName, mimeType string) string {
	mimeType = strings.ToLower(strings.TrimSpace(mimeType))
	ext := strings.ToLower(strings.TrimSpace(filepath.Ext(displayName)))
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return "images"
	case strings.HasPrefix(mimeType, "audio/"):
		return "audio"
	case strings.HasPrefix(mimeType, "video/"):
		return "video"
	case strings.HasPrefix(mimeType, "text/"):
		return "text"
	case mimeType == "application/pdf",
		mimeType == "application/msword",
		mimeType == "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		mimeType == "application/vnd.ms-excel",
		mimeType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		mimeType == "application/vnd.ms-powerpoint",
		mimeType == "application/vnd.openxmlformats-officedocument.presentationml.presentation":
		return "documents"
	case mimeType == "application/zip",
		mimeType == "application/x-gzip",
		mimeType == "application/gzip",
		mimeType == "application/x-tar",
		mimeType == "application/x-7z-compressed",
		mimeType == "application/x-rar-compressed":
		return "archives"
	}

	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".bmp", ".svg":
		return "images"
	case ".mp3", ".wav", ".m4a", ".aac", ".flac", ".ogg":
		return "audio"
	case ".mp4", ".mov", ".avi", ".mkv", ".webm":
		return "video"
	case ".txt", ".md", ".json", ".yaml", ".yml", ".csv", ".log":
		return "text"
	case ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx":
		return "documents"
	case ".zip", ".tar", ".gz", ".tgz", ".bz2", ".7z", ".rar":
		return "archives"
	default:
		return "files"
	}
}

func finalizeUploadPath(dataDir, attachmentID, displayName, mimeType string) (string, error) {
	category := uploadDirectoryCategory(displayName, mimeType)
	dir := filepath.Join(filepath.Clean(dataDir), "attachments", category)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create upload dir %q: %w", dir, err)
	}
	return filepath.Join(dir, buildStoredUploadFilename(attachmentID, displayName)), nil
}

func removeStoredAttachments(attachments []storedTurnAttachment) {
	for _, attachment := range attachments {
		path := strings.TrimSpace(attachment.FilePath)
		if path == "" {
			continue
		}
		_ = os.Remove(path)
	}
}

func uploadTempDir() string {
	if info, err := os.Stat("/tmp"); err == nil && info.IsDir() {
		return "/tmp"
	}
	return os.TempDir()
}

func fileURIForPath(path string) string {
	path = filepath.Clean(strings.TrimSpace(path))
	slashPath := filepath.ToSlash(path)
	if volume := filepath.VolumeName(path); volume != "" && !strings.HasPrefix(slashPath, "/") {
		slashPath = "/" + slashPath
	}
	return (&url.URL{Scheme: "file", Path: slashPath}).String()
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		ResponseWriter: w,
		statusCode:     0,
		bytesWritten:   0,
	}
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	if w.statusCode == 0 {
		w.statusCode = statusCode
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *loggingResponseWriter) Write(body []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(body)
	w.bytesWritten += n
	return n, err
}

func (w *loggingResponseWriter) StatusCode() int {
	if w.statusCode == 0 {
		return http.StatusOK
	}
	return w.statusCode
}

func (w *loggingResponseWriter) BytesWritten() int {
	return w.bytesWritten
}

func (w *loggingResponseWriter) Flush() {
	flusher, ok := w.ResponseWriter.(http.Flusher)
	if !ok {
		return
	}
	flusher.Flush()
}

func (w *loggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("response writer does not support hijacking")
	}
	return hijacker.Hijack()
}

func (w *loggingResponseWriter) Push(target string, opts *http.PushOptions) error {
	pusher, ok := w.ResponseWriter.(http.Pusher)
	if !ok {
		return http.ErrNotSupported
	}
	return pusher.Push(target, opts)
}

func (w *loggingResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func requestClientAddr(r *http.Request) string {
	if r == nil {
		return "unknown"
	}

	if xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			if ip := strings.TrimSpace(parts[0]); ip != "" {
				return ip
			}
		}
	}

	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}

	remoteAddr := strings.TrimSpace(r.RemoteAddr)
	if remoteAddr == "" {
		return "unknown"
	}

	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil && strings.TrimSpace(host) != "" {
		return strings.TrimSpace(host)
	}

	return remoteAddr
}

func requestLogPath(r *http.Request) string {
	if r == nil || r.URL == nil {
		return "/"
	}

	path := strings.TrimSpace(r.URL.RequestURI())
	if path == "" {
		path = strings.TrimSpace(r.URL.Path)
	}
	if path == "" {
		path = "/"
	}
	return observability.RedactString(path)
}

func (s *Server) isAuthorized(r *http.Request) bool {
	if s.authToken == "" {
		return true
	}

	const prefix = "Bearer "
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, prefix) {
		return false
	}

	provided := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
	return s.matchesAuthToken(provided)
}

func (s *Server) isAttachmentAuthorized(r *http.Request) bool {
	if s.isAuthorized(r) {
		return true
	}
	if s.authToken == "" || r == nil || r.URL == nil {
		return s.authToken == ""
	}
	return s.matchesAuthToken(strings.TrimSpace(r.URL.Query().Get("access_token")))
}

func (s *Server) matchesAuthToken(provided string) bool {
	if provided == "" {
		return false
	}

	if len(provided) != len(s.authToken) {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(provided), []byte(s.authToken)) == 1
}

func writeMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusMethodNotAllowed, codeInvalidArgument, "method is not allowed for this endpoint", map[string]any{
		"method": r.Method,
		"path":   r.URL.Path,
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	encoder := json.NewEncoder(w)
	_ = encoder.Encode(payload)
}

func writeError(w http.ResponseWriter, statusCode int, code, message string, details map[string]any) {
	if details == nil {
		details = map[string]any{}
	}

	writeJSON(w, statusCode, map[string]any{
		"error": map[string]any{
			"code":    code,
			"message": message,
			"details": details,
		},
	})
}

func isPathAllowed(path string, roots []string) bool {
	path = filepath.Clean(path)
	for _, root := range roots {
		root = filepath.Clean(root)
		rel, err := filepath.Rel(root, path)
		if err != nil {
			continue
		}
		if rel == "." {
			return true
		}
		if rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return true
		}
	}
	return false
}

func sortedAgentIDs(allowed map[string]struct{}) []string {
	if len(allowed) == 0 {
		return []string{}
	}
	ids := make([]string, 0, len(allowed))
	for id := range allowed {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func expandPath(path string) (string, error) {
	if path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return home, nil
	}
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}
