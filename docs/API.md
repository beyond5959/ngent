# API

This document defines the current HTTP API contract.

## Common Conventions

- JSON response content type: `application/json; charset=utf-8`.
- Except `/healthz`, every `/v1/*` endpoint requires `X-Client-ID` header (non-empty).
- `X-Client-ID` is retained as a required compatibility header, but it is not persisted in SQLite and it is not a thread/session access boundary.
- threads, sessions, permissions, persisted attachments, and recent-directory suggestions are shared across callers connected to the same ngent instance.
- Optional auth switch:
  - if server starts with `--auth-token=<token>`, `/v1/*` also requires `Authorization: Bearer <token>`.

## Runtime Logging Conventions

- Startup prints a human-readable multi-line summary on `stderr` with `Time`, `HTTP`, `Web`, `DB`, `Agents`, and `Help`.
- Every HTTP request emits one human-readable access-log line on `stderr`, for example:
  - `INFO: 2026-03-23 15:30:45 127.0.0.1 - "GET /v1/threads HTTP/1.1" 200 OK 12.4ms`
- When `stderr` is attached to a TTY, access logs and level labels may use ANSI colors; redirected output stays plain text.
- When server starts with `--debug=true`, `stderr` also emits readable `acp.message` debug lines for ACP JSON-RPC traffic with:
  - `component`
  - `direction` (`inbound|outbound`)
  - `rpcType` (`request|response|notification`)
  - `method` when present
  - sanitized `rpc` payload with sensitive fields redacted

## Unified Error Envelope

All errors use:

```json
{
  "error": {
    "code": "INVALID_ARGUMENT",
    "message": "human-readable message",
    "details": {
      "field": "cwd"
    }
  }
}
```

## Implemented Endpoints

### Frontend (Web UI)

11. `GET /`
- No authentication required.
- Returns the embedded web UI (`index.html`).
- Response `200`: `text/html; charset=utf-8`.

12. `GET /assets/*`
- No authentication required.
- Returns embedded static assets (JS, CSS, fonts) produced by the frontend build.
- SPA fallback: any non-API, non-asset path also returns `index.html` so the client-side router can handle it.

### Health

1. `GET /healthz`
- Response `200`:

```json
{
  "ok": true
}
```

2. `GET /v1/agents`
- Headers: `X-Client-ID` (required), optional bearer auth if enabled.
- agent status contract:
  - each agent entry reports readiness as `available|unavailable`.
  - current built-in ids are `codex`, `claude`, `cursor`, `gemini`, `kimi`, `qwen`, `opencode`, and `blackbox`.
- Response `200`:

```json
{
  "agents": [
    {
      "id": "codex",
      "name": "Codex",
      "status": "available"
    },
    {
      "id": "claude",
      "name": "Claude Code",
      "status": "unavailable"
    }
  ]
}
```

2.1 `GET /v1/agents/{agentId}/models`
- Headers: `X-Client-ID` (required), optional bearer auth if enabled.
- Behavior:
  - queries the target agent via ACP (`initialize` + `session/new`) and returns runtime-reported model options.
  - returns `503 UPSTREAM_UNAVAILABLE` when the agent runtime is unavailable or model discovery handshake fails.
- Response `200`:

```json
{
  "agentId": "codex",
  "models": [
    {"id": "gpt-5", "name": "GPT-5"},
    {"id": "gpt-5-mini", "name": "GPT-5 Mini"}
  ]
}
```

3. `POST /v1/threads`
- Headers: `X-Client-ID` (required), optional bearer auth if enabled.
- Request:

```json
{
  "agent": "<agent-id>",
  "cwd": "/abs/path",
  "title": "optional",
  "agentOptions": {
    "mode": "safe",
    "modelId": "gpt-5"
  }
}
```

- Validation:
  - `agent` must be in the current runtime allowlist (derived from agents whose startup preflight succeeds in the running environment).
  - `cwd` must be absolute.
  - server default policy accepts any absolute `cwd`.
  - create thread only persists row; no agent process is started.

- Response `200`:

```json
{
  "threadId": "th_..."
}
```

4. `GET /v1/threads`
- Headers: `X-Client-ID` (required), optional bearer auth if enabled.
- Behavior:
  - returns every persisted thread on the current ngent instance, not just threads created by the current `X-Client-ID`.
  - each thread payload includes `hasActiveSession`, which is true when ngent currently has at least one active session-scoped turn running on that thread.
- Response `200`:

```json
{
  "threads": [
    {
      "threadId": "th_...",
      "agent": "<agent-id>",
      "cwd": "/abs/path",
      "title": "optional",
      "agentOptions": {},
      "summary": "",
      "hasActiveSession": false,
      "createdAt": "2026-02-28T00:00:00Z",
      "updatedAt": "2026-02-28T00:00:00Z"
    }
  ]
}
```

5. `GET /v1/threads/{threadId}`
- Headers: `X-Client-ID` (required), optional bearer auth if enabled.
- Visibility rule:
  - if thread does not exist, return `404`.
- Response `200`:

```json
{
  "thread": {
    "threadId": "th_...",
    "agent": "<agent-id>",
    "cwd": "/abs/path",
    "title": "optional",
    "agentOptions": {},
    "summary": "",
    "hasActiveSession": false,
    "createdAt": "2026-02-28T00:00:00Z",
    "updatedAt": "2026-02-28T00:00:00Z"
  }
}
```

5.1 `PATCH /v1/threads/{threadId}`
- Headers: `X-Client-ID` (required), optional bearer auth if enabled.
- Visibility rule:
  - if thread does not exist, return `404`.
- Request:

```json
{
  "title": "optional new title",
  "agentOptions": {
    "modelId": "gpt-5"
  }
}
```

- Behavior:
  - when `title` is present, trims surrounding whitespace, persists `thread.title`, and updates `updatedAt`.
  - when `agentOptions` is present, updates persisted `thread.agentOptions` and `updatedAt`.
  - if the update changes shared thread state (`title`, `modelId`, `configOverrides`, or other non-session fields) while any session on the thread is active, returns `409 CONFLICT`.
  - session-only `agentOptions.sessionId` updates are allowed while a different session on the same thread is active.
 - closes cached thread-scoped agent providers only when the update changes non-session agent options, so the next turn uses updated shared options.
- Response `200`:

```json
{
  "thread": {
    "threadId": "th_...",
    "agent": "<agent-id>",
    "cwd": "/abs/path",
    "title": "optional",
    "agentOptions": {
      "modelId": "gpt-5"
    },
    "summary": "",
    "createdAt": "2026-02-28T00:00:00Z",
    "updatedAt": "2026-02-28T00:05:00Z"
  }
}
```

5.2 `GET /v1/threads/{threadId}/git`
- Headers: `X-Client-ID` (required), optional bearer auth if enabled.
- Behavior:
  - inspects git state for the thread `cwd`.
  - if the host does not have `git`, or the `cwd` is not inside a git repository, returns `200` with `available=false`.
  - otherwise returns repository root, current ref metadata, and local branches.
- Response `200` (repository-backed thread):

```json
{
  "threadId": "th_...",
  "available": true,
  "repoRoot": "/abs/path/to/repo",
  "currentRef": "main",
  "currentBranch": "main",
  "detached": false,
  "branches": [
    {"name": "main", "current": true},
    {"name": "feature/demo", "current": false}
  ]
}
```

- Response `200` (non-git or gitless host):

```json
{
  "threadId": "th_...",
  "available": false
}
```

5.3 `POST /v1/threads/{threadId}/git`
- Headers: `X-Client-ID` (required), optional bearer auth if enabled.
- Request:

```json
{
  "branch": "feature/demo"
}
```

- Behavior:
  - only existing local branches are accepted.
  - branch checkout uses the thread-wide exclusive guard.
  - if any turn on the thread is active, returns `409 CONFLICT`.
  - success response reuses the same payload shape as `GET /v1/threads/{threadId}/git`.

5.2 `DELETE /v1/threads/{threadId}`
- Headers: `X-Client-ID` (required), optional bearer auth if enabled.
- Visibility rule:
  - if thread does not exist, return `404`.
- Behavior:
  - hard-deletes thread history (thread row + turns + events).
  - if any session on the thread has an active turn, returns `409 CONFLICT`.
- Response `200`:

```json
{
  "threadId": "th_...",
  "status": "deleted"
}
```

6. `POST /v1/threads/{threadId}/turns`
- Headers: `X-Client-ID` (required), optional bearer auth if enabled.
- Request:

```json
{
  "input": "hello",
  "stream": true
}
```

- Behavior:
  - response is SSE (`text/event-stream`).
  - same `(thread, sessionId)` scope allows only one active turn at a time.
  - if another turn is active on that same scope, return `409 CONFLICT`.
  - different sessions on the same thread may run concurrently after switching `agentOptions.sessionId`.
  - disconnecting this specific viewer does not cancel a healthy active turn by itself.
  - if provider requests runtime permission, server emits `permission_required` and pauses turn until decision/timeout.

- SSE event types:
  - all replayable event payloads include a monotonic per-turn `seq`.
  - `turn_started`: `{"turnId":"...","seq":1}`
  - `message_delta`: `{"turnId":"...","seq":2,"delta":"..."}`
  - `plan_update`: `{"turnId":"...","seq":3,"entries":[{"content":"...","status":"pending|in_progress|completed","priority":"low|medium|high"}]}`
  - `permission_required`: `{"turnId":"...","seq":4,"permissionId":"...","approval":"command|file|network|mcp","command":"...","requestId":"...","options":[{"optionId":"...","name":"...","kind":"allow_once|allow_always|reject_once|reject_always|..."}]}`
  - `permission_resolved`: `{"turnId":"...","seq":5,"permissionId":"...","outcome":"approved|declined|cancelled","optionId":"...","optionName":"..."}` (`optionId` / `optionName` are present when the decision selected one advertised option)
  - `turn_completed`: `{"turnId":"...","seq":6,"stopReason":"end_turn|cancelled|error"}`
  - `error`: `{"turnId":"...","seq":6,"code":"...","message":"..."}`
  - for ACP `sessionUpdate == "plan"`, the server emits `plan_update` and treats each payload as a full replacement of the current plan list.

- Permission fail-closed contract:
  - permission request timeout or missing/invalid decision defaults to `declined`.
  - disconnecting one viewer does not auto-decline by itself; another viewer may still reconnect and resolve the pending permission before timeout.
  - fake ACP flow uses terminal `stopReason="cancelled"` for `declined`/`cancelled`.

6.1 `GET /v1/turns/{turnId}/events`
- Headers: `X-Client-ID` (required), optional bearer auth if enabled.
- Query:
  - `after=<seq>` (optional integer, default `0`)
- Behavior:
  - response is SSE (`text/event-stream`).
  - replays persisted events for the turn with `seq > after`, then tails newly published live events for the same turn.
  - multiple clients may attach to the same active turn concurrently.
  - disconnecting this stream does not cancel the underlying turn.
  - stream closes naturally after a terminal event (`turn_completed` or `error`) or when the client disconnects.
- Event payloads:
  - same shapes as the `POST /v1/threads/{threadId}/turns` SSE event types above, including per-turn `seq`.

7. `POST /v1/turns/{turnId}/cancel`
- Headers: `X-Client-ID` (required), optional bearer auth if enabled.
- Behavior:
  - requests cancellation for active turn.
  - terminal stream event should end with `stopReason=cancelled` if cancellation wins race.
- Response `200`:

```json
{
  "turnId": "tu_...",
  "threadId": "th_...",
  "status": "cancelling"
}
```

8. `GET /v1/threads/{threadId}/history`
- Headers: `X-Client-ID` (required), optional bearer auth if enabled.
- Query:
  - `includeEvents=true|1` (optional, default false)
  - `includeInternal=true|1` (optional, default false)
  - `sessionId=<id>` (optional)
- Response `200`:

```json
{
  "sessionTranscript": {
    "supported": true,
    "messages": [
      {
        "role": "user",
        "content": "hello",
        "timestamp": "2026-02-28T00:00:00Z"
      },
      {
        "role": "assistant",
        "content": "world",
        "timestamp": "2026-02-28T00:00:01Z"
      }
    ]
  },
  "turns": [
    {
      "turnId": "tu_...",
      "requestText": "hello",
      "responseText": "hello",
      "status": "completed",
      "stopReason": "end_turn",
      "errorMessage": "",
      "createdAt": "2026-02-28T00:00:00Z",
      "completedAt": "2026-02-28T00:00:01Z",
      "events": [
        {
          "eventId": 1,
          "seq": 1,
          "type": "turn_started",
          "data": {
            "turnId": "tu_..."
          },
          "createdAt": "2026-02-28T00:00:00Z"
        }
      ]
    }
  ]
}
```

- Behavior:
  - without `sessionId`, returns the thread's persisted ngent turns and optional events as before.
  - with `sessionId`, filters persisted turns server-side by `session_bound`.
  - `/history` may still attach cached provider transcript replay under `sessionTranscript` when sqlite `session_transcript_cache` already has a snapshot for that selected session.
  - only when that selected session currently has no filtered turns in this response does `/history` fall back to provider/ACP `session/load` on a cache miss.
  - if provider transcript replay fails, the request still returns the persisted `turns`; `sessionTranscript` may be omitted.

9. `POST /v1/permissions/{permissionId}`
- Headers: `X-Client-ID` (required), optional bearer auth if enabled.
- Request:

```json
{
  "outcome": "approved"
}
```

Or submit the provider's exact option id when the permission request advertised options:

```json
{
  "optionId": "allow_always_opt"
}
```

- Behavior:
  - `outcome` remains supported for generic approve / decline / cancel flows.
  - `optionId` lets clients return the provider's exact permission choice when multiple options are available.
  - clients may send both `outcome` and `optionId`; when `optionId` is present, the server forwards that exact selection back to option-aware providers.

10. `POST /v1/threads/{threadId}/compact`
- Headers: `X-Client-ID` (required), optional bearer auth if enabled.
- Request (optional body):

```json
{
  "maxSummaryChars": 1200
}
```

- Behavior:
  - triggers one internal summarization turn (`is_internal=1`).
  - updates `threads.summary` on success.
  - internal compact turn is hidden from default history.

- Response `200`:

```json
{
  "threadId": "th_...",
  "turnId": "tu_...",
  "status": "completed",
  "stopReason": "end_turn",
  "summary": "updated summary text",
  "summaryChars": 324
}
```

- Validation:
  - `outcome` must be one of `approved|declined|cancelled`.
  - `permissionId` must exist.
  - already-resolved permission returns `409 CONFLICT`.

- Response `200`:

```json
{
  "permissionId": "perm_...",
  "status": "recorded",
  "outcome": "approved"
}
```

## Baseline Error Codes

- `INVALID_ARGUMENT`: validation failed.
- `UNAUTHORIZED`: bearer token missing or invalid.
- `FORBIDDEN`: path/policy denied.
- `NOT_FOUND`: endpoint/resource missing.
- `CONFLICT`: active-turn conflict or invalid cancel state.
- `TIMEOUT`: upstream/model operation exceeded allowed time budget.
- `UPSTREAM_UNAVAILABLE`: configured agent/provider is unavailable or failed to start/respond.
- `INTERNAL`: unexpected server/storage failure.
