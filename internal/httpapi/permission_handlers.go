package httpapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/beyond5959/ngent/internal/agents"
)

var (
	errPermissionNotFound        = errors.New("permission not found")
	errPermissionAlreadyResolved = errors.New("permission already resolved")
	errPermissionInvalidOption   = errors.New("permission option is invalid")
	errPermissionOutcomeRequired = errors.New("permission outcome is required")
)

type pendingPermission struct {
	turnID  string
	options map[string]agents.PermissionOption

	ch   chan agents.PermissionResponse
	once sync.Once
}

func (s *Server) handlePermissionDecision(w http.ResponseWriter, r *http.Request, clientID, permissionID string) {
	if err := requireMethod(r, http.MethodPost); err != nil {
		writeMethodNotAllowed(w, r)
		return
	}

	var req struct {
		Outcome  string `json:"outcome"`
		OptionID string `json:"optionId"`
	}
	if err := decodeJSONBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, codeInvalidArgument, "invalid JSON body", map[string]any{"reason": err.Error()})
		return
	}

	response := agents.PermissionResponse{
		SelectedOptionID: strings.TrimSpace(req.OptionID),
	}
	if rawOutcome := strings.TrimSpace(req.Outcome); rawOutcome != "" {
		outcome, ok := normalizePermissionOutcome(rawOutcome)
		if !ok {
			writeError(w, http.StatusBadRequest, codeInvalidArgument, "outcome must be approved, declined, or cancelled", map[string]any{
				"field": "outcome",
			})
			return
		}
		response.Outcome = outcome
	}
	if response.Outcome == "" && response.SelectedOptionID == "" {
		writeError(w, http.StatusBadRequest, codeInvalidArgument, "outcome or optionId is required", map[string]any{
			"field": "outcome",
		})
		return
	}

	resolvedResponse, err := s.resolvePermission(permissionID, response)
	if err != nil {
		if errors.Is(err, errPermissionNotFound) {
			writeError(w, http.StatusNotFound, codeNotFound, "permission not found", map[string]any{})
			return
		}
		if errors.Is(err, errPermissionInvalidOption) {
			writeError(w, http.StatusBadRequest, codeInvalidArgument, "optionId must match one of the advertised permission options", map[string]any{
				"field":        "optionId",
				"permissionId": permissionID,
			})
			return
		}
		if errors.Is(err, errPermissionOutcomeRequired) {
			writeError(w, http.StatusBadRequest, codeInvalidArgument, "outcome is required for this permission selection", map[string]any{
				"field":        "outcome",
				"permissionId": permissionID,
			})
			return
		}
		if errors.Is(err, errPermissionAlreadyResolved) {
			writeError(w, http.StatusConflict, codeConflict, "permission already resolved", map[string]any{
				"permissionId": permissionID,
			})
			return
		}
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to resolve permission", map[string]any{
			"reason": err.Error(),
		})
		return
	}

	payload := map[string]any{
		"permissionId": permissionID,
		"status":       "recorded",
		"outcome":      string(resolvedResponse.Outcome),
	}
	if resolvedResponse.SelectedOptionID != "" {
		payload["optionId"] = resolvedResponse.SelectedOptionID
	}
	writeJSON(w, http.StatusOK, payload)
}

func newPendingPermission(turnID string, options []agents.PermissionOption) *pendingPermission {
	optionMap := make(map[string]agents.PermissionOption, len(options))
	for _, option := range options {
		optionID := strings.TrimSpace(option.OptionID)
		if optionID == "" {
			continue
		}
		optionMap[optionID] = agents.PermissionOption{
			OptionID: optionID,
			Name:     strings.TrimSpace(option.Name),
			Kind:     strings.TrimSpace(option.Kind),
		}
	}
	return &pendingPermission{
		turnID:  strings.TrimSpace(turnID),
		options: optionMap,
		ch:      make(chan agents.PermissionResponse, 1),
	}
}

func (p *pendingPermission) Resolve(response agents.PermissionResponse) bool {
	resolved := false
	p.once.Do(func() {
		p.ch <- normalizePermissionResponse(response)
		close(p.ch)
		resolved = true
	})
	return resolved
}

func (p *pendingPermission) normalizeDecision(response agents.PermissionResponse) (agents.PermissionResponse, error) {
	response = normalizePermissionResponse(response)
	if response.SelectedOptionID != "" {
		if len(p.options) > 0 {
			option, ok := p.options[response.SelectedOptionID]
			if !ok {
				return agents.PermissionResponse{}, errPermissionInvalidOption
			}
			if response.Outcome == "" {
				if inferred, ok := permissionOutcomeForOptionKind(option.Kind); ok {
					response.Outcome = inferred
				}
			}
		}
	}
	if response.Outcome == "" {
		return agents.PermissionResponse{}, errPermissionOutcomeRequired
	}
	return response, nil
}

func (p *pendingPermission) resolvedEventPayload(permissionID string, response agents.PermissionResponse) map[string]any {
	payload := map[string]any{
		"turnId":       p.turnID,
		"permissionId": strings.TrimSpace(permissionID),
		"outcome":      string(response.Outcome),
	}
	if optionID := strings.TrimSpace(response.SelectedOptionID); optionID != "" {
		payload["optionId"] = optionID
		if option, ok := p.options[optionID]; ok {
			if name := strings.TrimSpace(option.Name); name != "" {
				payload["optionName"] = name
			}
			if kind := strings.TrimSpace(option.Kind); kind != "" {
				payload["optionKind"] = kind
			}
		}
	}
	return payload
}

func normalizePermissionOutcome(raw string) (agents.PermissionOutcome, bool) {
	switch strings.TrimSpace(strings.ToLower(raw)) {
	case string(agents.PermissionOutcomeApproved):
		return agents.PermissionOutcomeApproved, true
	case string(agents.PermissionOutcomeDeclined):
		return agents.PermissionOutcomeDeclined, true
	case string(agents.PermissionOutcomeCancelled):
		return agents.PermissionOutcomeCancelled, true
	default:
		return "", false
	}
}

func permissionOutcomeForOptionKind(raw string) (agents.PermissionOutcome, bool) {
	switch strings.TrimSpace(strings.ToLower(raw)) {
	case "allow", "allow_once", "allow_always", "approve", "approved":
		return agents.PermissionOutcomeApproved, true
	case "deny", "denied", "decline", "declined", "reject", "reject_once", "reject_always":
		return agents.PermissionOutcomeDeclined, true
	case "cancel", "cancelled", "canceled":
		return agents.PermissionOutcomeCancelled, true
	default:
		return "", false
	}
}

func (s *Server) nextPermissionID(requestID string) string {
	seq := atomic.AddUint64(&s.permissionSeq, 1)
	safeRequestID := sanitizePermissionIDComponent(requestID)
	if safeRequestID == "" {
		return fmt.Sprintf("perm_%d", seq)
	}
	return fmt.Sprintf("perm_%s_%d", safeRequestID, seq)
}

func sanitizePermissionIDComponent(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	var builder strings.Builder
	builder.Grow(len(raw))
	for _, r := range raw {
		switch {
		case r >= 'a' && r <= 'z':
			builder.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			builder.WriteRune(r)
		case r >= '0' && r <= '9':
			builder.WriteRune(r)
		default:
			builder.WriteByte('_')
		}
	}
	return strings.Trim(builder.String(), "_")
}

func (s *Server) registerPermission(permissionID string, pending *pendingPermission) {
	s.permissionsMu.Lock()
	s.permissions[permissionID] = pending
	s.permissionsMu.Unlock()
}

func (s *Server) unregisterPermission(permissionID string, pending *pendingPermission) {
	s.permissionsMu.Lock()
	current, ok := s.permissions[permissionID]
	if ok && current == pending {
		delete(s.permissions, permissionID)
	}
	s.permissionsMu.Unlock()
}

func permissionFailClosedResponse() agents.PermissionResponse {
	return agents.PermissionResponse{Outcome: agents.PermissionOutcomeDeclined}
}

func normalizePermissionResponse(response agents.PermissionResponse) agents.PermissionResponse {
	response.SelectedOptionID = strings.TrimSpace(response.SelectedOptionID)
	if response.Outcome == "" {
		return response
	}
	switch response.Outcome {
	case agents.PermissionOutcomeApproved, agents.PermissionOutcomeDeclined, agents.PermissionOutcomeCancelled:
		return response
	default:
		return agents.PermissionResponse{
			Outcome:          agents.PermissionOutcomeDeclined,
			SelectedOptionID: response.SelectedOptionID,
		}
	}
}

func (s *Server) waitPermissionResponse(ctx context.Context, pending *pendingPermission) agents.PermissionResponse {
	timeout := s.permissionTimeout
	if timeout <= 0 {
		timeout = defaultPermissionTimeout
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case response := <-pending.ch:
		if response.Outcome == "" {
			return permissionFailClosedResponse()
		}
		return response
	case <-timer.C:
		pending.Resolve(permissionFailClosedResponse())
	case <-ctx.Done():
		pending.Resolve(permissionFailClosedResponse())
	}

	response, ok := <-pending.ch
	if !ok || response.Outcome == "" {
		return permissionFailClosedResponse()
	}
	return response
}

func (s *Server) resolvePermission(permissionID string, response agents.PermissionResponse) (agents.PermissionResponse, error) {
	s.permissionsMu.Lock()
	pending, ok := s.permissions[permissionID]
	s.permissionsMu.Unlock()
	if !ok {
		return agents.PermissionResponse{}, errPermissionNotFound
	}
	normalized, err := pending.normalizeDecision(response)
	if err != nil {
		return agents.PermissionResponse{}, err
	}
	if !pending.Resolve(normalized) {
		return agents.PermissionResponse{}, errPermissionAlreadyResolved
	}
	return normalized, nil
}
