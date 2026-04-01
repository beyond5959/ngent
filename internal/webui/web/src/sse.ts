import type { PermissionOption, PlanEntry, SessionUsage } from './types.ts'

// ── SSE event payloads (mirror server API contract) ───────────────────────

interface TurnEventMeta {
  seq?: number
}

export interface TurnStartedPayload extends TurnEventMeta {
  turnId: string
}

export interface MessageDeltaPayload extends TurnEventMeta {
  turnId: string
  delta: string
}

export interface MessageContentPayload extends TurnEventMeta {
  turnId: string
  content?: unknown
}

export interface ReasoningDeltaPayload extends TurnEventMeta {
  turnId: string
  delta: string
}

export interface TurnCompletedPayload extends TurnEventMeta {
  turnId: string
  stopReason: string
}

export interface PlanUpdatePayload extends TurnEventMeta {
  turnId: string
  entries: PlanEntry[]
}

export interface ToolCallPayload extends TurnEventMeta {
  turnId: string
  toolCallId: string
  title?: string
  kind?: string
  status?: string
  content?: unknown
  locations?: unknown
  rawInput?: unknown
  rawOutput?: unknown
}

export interface SessionBoundPayload extends TurnEventMeta {
  threadId: string
  sessionId: string
}

export interface SessionInfoUpdatePayload extends TurnEventMeta {
  turnId: string
  sessionId: string
  title: string
}

export interface SessionUsageUpdatePayload extends SessionUsage, TurnEventMeta {
  turnId: string
}

export interface TurnErrorPayload extends TurnEventMeta {
  turnId: string
  code: string
  message: string
}

export interface PermissionRequiredPayload extends TurnEventMeta {
  turnId: string
  permissionId: string
  approval: string
  command: string
  requestId: string
  options?: PermissionOption[]
}

export interface PermissionResolvedPayload extends TurnEventMeta {
  turnId: string
  permissionId: string
  outcome: 'approved' | 'declined' | 'cancelled'
  optionId?: string
  optionName?: string
  optionKind?: string
}

// ── Callbacks ─────────────────────────────────────────────────────────────

export interface TurnStreamCallbacks {
  onTurnStarted?:        (e: TurnStartedPayload) => void
  onDelta?:              (e: MessageDeltaPayload) => void
  onMessageContent?:     (e: MessageContentPayload) => void
  onReasoningDelta?:     (e: ReasoningDeltaPayload) => void
  onPlanUpdate?:         (e: PlanUpdatePayload) => void
  onToolCall?:           (e: ToolCallPayload) => void
  onToolCallUpdate?:     (e: ToolCallPayload) => void
  onSessionBound?:       (e: SessionBoundPayload) => void
  onSessionInfoUpdate?:  (e: SessionInfoUpdatePayload) => void
  onSessionUsageUpdate?: (e: SessionUsageUpdatePayload) => void
  onCompleted?:          (e: TurnCompletedPayload) => void
  onError?:              (e: TurnErrorPayload) => void
  onPermissionRequired?: (e: PermissionRequiredPayload) => void
  onPermissionResolved?: (e: PermissionResolvedPayload) => void
  /** Called when the fetch connection closes unexpectedly (not via abort). */
  onDisconnect?:         () => void
}

// ── TurnStream ────────────────────────────────────────────────────────────

/**
 * Reads a POST SSE stream using the Fetch API.
 * EventSource is not used because it only supports GET without custom headers.
 */
export class TurnStream {
  private readonly aborter = new AbortController()
  /** Set to true once a terminal event (completed/error) or abort() is received. */
  private terminated = false

  constructor(
    private readonly fetchUrl: string,
    private readonly request: {
      method: 'GET' | 'POST'
      headers: Record<string, string>
      body?: FormData | Record<string, unknown> | null
    },
    private readonly callbacks: TurnStreamCallbacks,
  ) {}

  /** Starts the stream. Resolves when the stream ends (for any reason). */
  async start(): Promise<void> {
    let res: Response
    try {
      res = await fetch(this.fetchUrl, {
        method: this.request.method,
        headers: this.request.headers,
        body: this.request.body instanceof FormData
          ? this.request.body
          : (this.request.body ? JSON.stringify(this.request.body) : undefined),
        signal: this.aborter.signal,
      })
    } catch {
      if (!this.terminated) this.callbacks.onDisconnect?.()
      return
    }

    if (!res.ok || !res.body) {
      let code = 'INTERNAL'
      let message = `HTTP ${res.status}`
      try {
        const payload = (await res.json()) as {
          error?: { code?: string; message?: string }
        }
        if (payload.error) {
          code    = payload.error.code    ?? code
          message = payload.error.message ?? message
        }
      } catch { /* ignore */ }
      this.callbacks.onError?.({ turnId: '', code, message })
      return
    }

    const reader = res.body.getReader()
    const decoder = new TextDecoder()
    let buf = ''

    try {
      for (;;) {
        const { done, value } = await reader.read()
        if (done) break
        buf += decoder.decode(value, { stream: true })
        // SSE messages are separated by blank lines (\n\n)
        const parts = buf.split('\n\n')
        buf = parts.pop() ?? ''
        for (const part of parts) {
          this.dispatch(part.trim())
        }
        if (this.terminated) break
      }
    } catch {
      if (!this.terminated) this.callbacks.onDisconnect?.()
    } finally {
      reader.releaseLock()
    }
  }

  /** Aborts the underlying fetch request. */
  abort(): void {
    this.terminated = true
    this.aborter.abort()
  }

  // ── SSE parsing ──────────────────────────────────────────────────────────

  private dispatch(raw: string): void {
    if (!raw) return

    let eventType = 'message'
    let dataStr = ''

    for (const line of raw.split('\n')) {
      if (line.startsWith('event:')) {
        eventType = line.slice(6).trim()
      } else if (line.startsWith('data:')) {
        dataStr = line.slice(5).trim()
      }
    }

    if (!dataStr) return

    let payload: Record<string, unknown>
    try {
      payload = JSON.parse(dataStr) as Record<string, unknown>
    } catch {
      return
    }

    switch (eventType) {
      case 'turn_started':
        this.callbacks.onTurnStarted?.(payload as unknown as TurnStartedPayload)
        break
      case 'message_delta':
        this.callbacks.onDelta?.(payload as unknown as MessageDeltaPayload)
        break
      case 'message_content':
        this.callbacks.onMessageContent?.(payload as unknown as MessageContentPayload)
        break
      case 'reasoning_delta':
        this.callbacks.onReasoningDelta?.(payload as unknown as ReasoningDeltaPayload)
        break
      case 'plan_update':
        this.callbacks.onPlanUpdate?.(payload as unknown as PlanUpdatePayload)
        break
      case 'tool_call':
        this.callbacks.onToolCall?.(payload as unknown as ToolCallPayload)
        break
      case 'tool_call_update':
        this.callbacks.onToolCallUpdate?.(payload as unknown as ToolCallPayload)
        break
      case 'session_bound':
        this.callbacks.onSessionBound?.(payload as unknown as SessionBoundPayload)
        break
      case 'session_info_update':
        this.callbacks.onSessionInfoUpdate?.(payload as unknown as SessionInfoUpdatePayload)
        break
      case 'session_usage_update':
        this.callbacks.onSessionUsageUpdate?.(payload as unknown as SessionUsageUpdatePayload)
        break
      case 'turn_completed':
        this.terminated = true
        this.callbacks.onCompleted?.(payload as unknown as TurnCompletedPayload)
        break
      case 'error':
        this.terminated = true
        this.callbacks.onError?.(payload as unknown as TurnErrorPayload)
        break
      case 'permission_required':
        this.callbacks.onPermissionRequired?.(payload as unknown as PermissionRequiredPayload)
        break
      case 'permission_resolved':
        this.callbacks.onPermissionResolved?.(payload as unknown as PermissionResolvedPayload)
        break
    }
  }
}
