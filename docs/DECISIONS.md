# DECISIONS

## ADR Index

- ADR-095: Preview absolute local markdown file links through the shared right-side drawer. (Accepted)
- ADR-094: Skip live provider transcript replay once a selected session already has filtered turns. (Accepted)
- ADR-093: Integrate Pi Agent through the embedded `pkg/piacp` runtime. (Accepted)
- ADR-092: Fold provider session replay into session-scoped `/history` responses. (Accepted)
- ADR-091: Persist learned config snapshots per provider session. (Accepted)
- ADR-090: Learn model/reasoning metadata only from real session lifecycle events. (Accepted)
- ADR-089: Share repeated ACP discovery and session-param helpers across built-in providers. (Accepted)
- ADR-088: Derive the runtime agent list from startup preflight results. (Accepted)
- ADR-087: Render assistant turns as ordered UI segments instead of one aggregated bubble. (Accepted)
- ADR-086: Preserve ACP tool-call updates as first-class turn events. (Accepted)
- ADR-085: Normalize rich ACP permission requests before bridging them into ngent. (Accepted)
- ADR-084: Persist ACP slash commands as agent-level SQLite snapshots. (Accepted)
- ADR-083: Allow concurrent turns across different sessions on the same thread. (Accepted)
- ADR-082: Scope Web UI chat playback to the selected ACP session. (Accepted)
- ADR-081: Persist thread-level ACP session selection and resume through provider sessions. (Accepted)
- ADR-080: Persist agent config catalogs in SQLite and refresh them asynchronously on startup. (Accepted)
- ADR-079: Thread-level model switching via ACP session config options. (Accepted)
- ADR-078: Keep the Web UI git-diff drawer explicitly dismissed and stable across summary polling. (Accepted)
- ADR-077: Preview git-diff file details in a right-side drawer, with text-only support for new files. (Accepted)
- ADR-076: Make grouped thread session lists collapsible from the leading agent glyph. (Accepted)
- ADR-075: Decorate thread session-list responses with cross-client active-session state. (Accepted)
- ADR-074: Merge thread and session browsing into one grouped left rail. (Accepted)
- ADR-073: Render live ACP plan updates as an ephemeral bottom overlay instead of transcript history. (Accepted)
- ADR-072: Render git-diff expanded file rows with locally vendored suffix-based file icons. (Accepted)
- ADR-071: Surface session-selected git diff summaries above the Web UI composer, including untracked files. (Accepted)
- ADR-070: Expand embedded Web UI localization and repository READMEs to Spanish and French. (Accepted)
- ADR-069: Keep active turns independent of individual SSE viewers and resume them through per-turn event streams. (Accepted)
- ADR-068: Add browser-default English/Simplified Chinese localization to the embedded Web UI. (Accepted)
- ADR-067: Persist ACP session usage snapshots and surface context-window pressure in the Web UI. (Accepted)
- ADR-066: Surface thread-scoped git branch state in the Web UI composer. (Accepted)
- ADR-065: Recast the embedded Web UI as a restrained desktop workbench. (Accepted)
- ADR-064: Share threads and sessions across browser-scoped client IDs on the same ngent instance. (Accepted)
- ADR-062: Decouple viewed Web UI session from backend thread session during active turns. (Accepted)
- ADR-061: Compact historical delta runs on read and render large chats incrementally. (Accepted)
- ADR-060: Make thread history session-scoped for Web UI session switches. (Accepted)
- ADR-059: Store uploaded attachments under the configurable data directory and serve them back through a stable attachment route. (Accepted)
- ADR-058: Render bracketed inline base64 user-image placeholders as safe Web UI previews. (Accepted)
- ADR-057: Persist Web UI uploads as local temp files and forward them as ACP resource links. (Accepted)
- ADR-056: Preserve exact provider permission options through the hub permission flow. (Accepted)
- ADR-055: Preserve non-text assistant ACP content as first-class turn events. (Accepted)
- ADR-054: Refresh the embedded Web UI as a premium workbench without changing behavior. (Superseded)
- ADR-053: Replace `slog` JSON output with a human-readable stderr logger and colored access logs. (Accepted)
- ADR-001: HTTP/JSON API with SSE streaming transport. (Accepted)
- ADR-002: Client identity via `X-Client-ID` header. (Superseded)
- ADR-003: SQLite append-only events table as interaction source of truth. (Accepted)
- ADR-004: Permission handling defaults to fail-closed. (Accepted)
- ADR-005: Default bind is localhost only. (Accepted)
- ADR-006: M1 API baseline for health/auth/agents. (Accepted)
- ADR-007: M3 thread API tenancy and path policy. (Accepted)
- ADR-008: M4 turn streaming over SSE with persisted event log. (Accepted)
- ADR-009: M5 ACP stdio provider and permission bridge. (Accepted)
- ADR-010: M6 codex-acp-go runtime wiring. (Superseded)
- ADR-011: M7 context window injection and compact policy. (Accepted)
- ADR-012: M8 reliability alignment (TTL, shutdown, error codes). (Accepted)
- ADR-013: Codex provider migration from sidecar binary to embedded library. (Accepted)
- ADR-015: First-turn prompt passthrough for slash-command compatibility in embedded codex mode. (Accepted)
- ADR-016: Remove `--allowed-root` runtime parameter and default to absolute-cwd policy. (Accepted)
- ADR-017: Human-readable startup summary and request completion access logs. (Accepted)
- ADR-018: Embedded Web UI via Go embed. (Accepted)
- ADR-019: OpenCode ACP stdio provider. (Accepted)
- ADR-020: Gemini CLI ACP stdio provider. (Accepted)
- ADR-022: Qwen Code ACP stdio provider integration. (Accepted)
- ADR-023: Shared ACP stdio transport for OpenCode and Qwen providers. (Accepted)
- ADR-024: Claude Code embedded provider via claudeacp runtime. (Accepted)
- ADR-025: Hard-delete thread endpoint with active-turn lock. (Accepted)
- ADR-026: Thread-level model override update API and provider reset. (Accepted)
- ADR-048: Web UI fresh-session scopes for repeated `New session`. (Accepted)
- ADR-027: ACP-backed agent model catalog endpoint and UI dropdown wiring. (Accepted)
- ADR-028: Persist thread config overrides and surface reasoning control in Web UI. (Accepted)
- ADR-029: Consolidate sidebar thread actions into a drawer and reuse thread patch for rename. (Accepted)
- ADR-030: Pin local acp-adapter hotfix for codex app-server server-request compatibility. (Accepted)
- ADR-031: Kimi CLI ACP stdio provider with dual startup syntax fallback. (Accepted)
- ADR-032: Shared common agent config/state helper without protocol unification. (Accepted)
- ADR-033: Surface ACP plan updates as first-class SSE and Web UI state. (Accepted)
- ADR-035: Add opt-in ACP debug tracing behind `--debug`. (Accepted)
- ADR-036: Persist stable Codex session ids and normalize Codex transcript replay. (Superseded)
- ADR-037: Replay Kimi session transcript from local Kimi session files. (Superseded)
- ADR-038: Replay OpenCode session transcript from local OpenCode SQLite storage. (Superseded)
- ADR-039: Standardize session transcript replay on ACP `session/load`. (Accepted)
- ADR-040: Cache session transcript replay snapshots in SQLite. (Accepted)
- ADR-041: Treat Web UI "New session" as provider-cache reset for the empty session scope. (Accepted)
- ADR-042: Treat explicit Web UI "New session" as a fresh turn with no injected thread context. (Accepted)
- ADR-043: Share one ACP CLI driver across Kimi/Qwen/OpenCode/Gemini. (Accepted)
- ADR-044: Normalize path-like ACP permission previews across direct ACP providers. (Accepted)
- ADR-045: Surface hidden agent reasoning as first-class SSE/history events in the Web UI. (Accepted)
- ADR-046: Collapse finalized Web UI thinking panels by default. (Accepted)
- ADR-047: Defer thread config-option apply until the next turn boundary. (Accepted)
- ADR-049: Align Web UI navigation with a left agent rail and left session panel. (Accepted)
- ADR-050: Keep the left agent rail permanently expanded. (Accepted)
- ADR-051: BLACKBOX AI ACP provider integration via shared ACP CLI driver. (Accepted)
- ADR-052: Cursor CLI ACP provider integration with explicit ACP authentication. (Accepted)

## ADR-095: Preview Absolute Local Markdown File Links Through The Shared Right-Side Drawer

- Status: Accepted
- Date: 2026-04-10
- Context:
  - the embedded SPA previously left those links as ordinary browser navigations, which pushed the browser to a meaningless local-app URL and showed no useful file content.
  - product requires that these message-linked files behave like the existing git-diff file surface:
    - open inside the same right-side inspection panel instead of navigating away.
    - support text preview and image preview.
    - stay fail-closed for unsupported types.
- Decision:
  - add thread-scoped preview endpoints:
    - `GET /v1/threads/{threadId}/file-preview?path=...`
    - `GET /v1/threads/{threadId}/file-preview-content?path=...`
  - authorize preview targets only when the requested absolute path resolves inside configured allowed roots; resolve symlinks before the allowed-root check so symlink escapes are rejected.
  - support only text files and image files from this surface.
    - text files return grouped rendered `blocks[]` for at most the first 10000 lines of the file.
    - image files return preview metadata from `file-preview` and raw bytes from `file-preview-content`.
    - other content types stay unsupported and do not get inline viewers or downloads here.
  - parse optional `#L<number>` markdown fragments on the client and pass the line number to the preview endpoint; the drawer highlights that focused line after render when it falls inside the returned 10000-line window.
  - render absolute local markdown links in finalized message markdown as ordinary inline links without pill chrome or extra file-type icons; unsupported extensions render as visibly disabled inline links.
  - reuse the existing right-side drawer shell for both git-diff previews and message-linked file previews, with message-linked opens explicitly dismissing any current git-diff file drawer first.
- Consequences:
  - clicking an absolute local file link in assistant output no longer navigates the browser away from the conversation.
  - large text files stay readable because the drawer loads only a bounded 10000-line prefix instead of trying to page through the entire file in the browser.
  - authenticated image previews remain local-first without leaking bearer tokens into query-string image URLs.
  - preview authorization stays within ngent's existing absolute-path and allowed-root safety boundary.
- Alternatives considered:
  - leave message links as normal anchors and rely on browser navigation (rejected: the destination view is empty/useless in the embedded app shell).
  - embed full file contents directly into the message stream (rejected: inflates turn payloads and makes large files impractical).
  - support arbitrary binary downloads/viewers from the same surface (rejected: broader security/UI surface than required for the current workflow).

## ADR-094: Skip Live Provider Transcript Replay Once A Selected Session Already Has Filtered Turns

- Status: Accepted
- Date: 2026-04-10
- Context:
  - ADR-092 let one session-scoped `/history` response carry both filtered persisted `turns` and optional provider-owned `sessionTranscript` replay.
  - the embedded Web UI merged those two sources when replay was available so provider-owned prehistory could appear ahead of ngent-owned turns.
  - real usage exposed a duplication bug: once a selected session already had visible filtered turns in `/history`, the provider replay could still reintroduce the same visible `user` / `assistant` transcript content, leaving the chat pane with duplicated history.
- Decision:
  - move the duplication guard to the backend history path instead of keeping it as a frontend-only merge rule.
  - when `GET /v1/threads/{threadId}/history?sessionId=...` finds one or more filtered turns for that selected session in the current response, skip live provider transcript loading through `session/load`.
  - in that filtered-turn case, `/history` may still attach `sessionTranscript` when sqlite `session_transcript_cache` already holds a snapshot for the same `(agent, cwd, sessionId)`, but it must not call the provider again on a cache miss.
  - only when the selected session currently has no filtered turns does `/history` fall back to provider `session/load` after checking cache.
  - keep the existing fresh-session promotion path in the Web UI; it still reuses browser-local messages after the first bind, and the UI continues rendering replay first and local turns after it whenever `sessionTranscript` is present.
- Consequences:
  - once filtered turns exist, session switches no longer trigger a fresh provider replay for chat reconstruction, removing one major source of duplicated visible transcript content.
  - provider-owned historical context remains available without a new provider call when that transcript snapshot was already warmed earlier into sqlite cache.
  - if a selected session later returns one or more filtered turns before any transcript snapshot was ever cached, later `/history` requests will not live-load the missing provider prehistory.
- Alternatives considered:
  - keep unconditional frontend merging and rely on overlap heuristics alone (rejected: brittle and already producing user-visible duplicates).
  - import provider replay directly into persisted `turns/events` (rejected for now: still conflates provider-owned transcript with ngent-owned durable history and needs a broader storage migration).

## ADR-093: Integrate Pi Agent Through The Embedded `pkg/piacp` Runtime

- Status: Accepted
- Date: 2026-04-07
- Context:
  - `acp-adapter` now exposes Pi through `pkg/piacp`, which matches ngent's existing embedded-provider pattern much more closely than the direct ACP CLI providers.
  - spawning `acp-adapter --adapter pi` as a second subprocess inside ngent would duplicate the adapter boundary even though ngent already embeds Codex and Claude directly.
  - the embedded Web UI also needs first-class Pi branding so users can recognize Pi-backed threads and distinguish the pure-white `pi.svg` in light theme.
- Decision:
  - add `internal/agents/pi` as a new embedded provider built on `github.com/beyond5959/acp-adapter/pkg/piacp`.
  - wire Pi through startup preflight, `/v1/agents`, per-thread lazy runtime resolution, model/config discovery, session list/load replay, slash-command caching, and the existing fail-closed permission bridge.
  - derive Pi runtime defaults from the same host environment knobs used by the standalone adapter (`PI_BIN`, `PI_PROVIDER`, `PI_MODEL`, `PI_SESSION_DIR`, `PI_DISABLE_GATE`) instead of inventing separate ngent-only flags.
  - render the provided `pi.svg` in the embedded Web UI and add a dark icon backing only in light theme so the white mark stays legible without changing the source asset.
- Consequences:
  - Pi now participates in the same in-process provider lifecycle as Codex and Claude, including idle-scope caching and lazy ACP session initialization.
  - thread/session replay, config pickers, slash commands, and permission cards work through the existing shared HTTP/UI surfaces with no Pi-specific API additions.
  - Pi threads still inherit the current upstream Pi limitations from `acp-adapter` such as missing MCP routing and plan snapshots; ngent documents those as a known limitation instead of simulating incomplete support.
- Alternatives considered:
  - wrap `acp-adapter --adapter pi` as a child process from ngent (rejected: unnecessary extra boundary and more process management).
  - treat Pi as a generic direct ACP CLI provider under `internal/agents/acpcli` (rejected: ngent does not own a raw Pi ACP stdio implementation; `pkg/piacp` is the supported embedding surface).

## ADR-092: Remove the standalone provider transcript replay endpoint and fold replay into session-scoped `/history`

- Status: Accepted
- Date: 2026-04-06
- Context:
  - the Web UI session picker reconstructed one selected session by issuing two backend requests: `GET /v1/threads/{threadId}/history?includeEvents=1&sessionId=...` for ngent-owned turns plus a second dedicated provider transcript replay request.
  - that split forced the frontend to own duplicate cache/provider fallback behavior, extra error handling, and a second round-trip even though both payloads belong to the same visible selected-session replay view.
  - the existing transcript loader already had the desired backend behavior: sqlite `session_transcript_cache` first, provider/ACP `session/load` only on cache miss, and one extra live load when config metadata still needed to be learned.
- Decision:
  - keep `GET /v1/threads/{threadId}/history` as the canonical selected-session replay endpoint.
  - when `sessionId` is present, `/history` still returns filtered persisted `turns`, and may also attach `sessionTranscript { supported, messages }` populated through the shared sqlite-first/provider-fallback transcript loader.
  - factor the replay loader into one shared backend helper owned by `/history`.
  - if transcript replay fails for reasons such as upstream unavailability or a stale session id, keep returning the persisted `turns` and omit `sessionTranscript` instead of failing the whole `/history` request.
  - remove the standalone provider transcript replay endpoint entirely now that the embedded Web UI and tests no longer depend on it.
- Consequences:
  - the Web UI session-switch flow now needs only one history request.
  - persisted ngent turn/event history remains the source of truth for rich artifacts such as reasoning, tool calls, and structured content; provider replay is still kept separate and is not imported into `turns/events`.
  - there is now only one selected-session replay entrypoint to maintain, document, and test.
- Alternatives considered:
  - keep a long-lived dedicated transcript replay alias (rejected: preserves avoidable surface area and duplicated documentation/testing burden).
  - fully merge provider replay into persisted `turns/events` (rejected: still blurs provider-owned transcript data with hub-owned history and risks duplication).
  - keep two frontend requests and only document the split more clearly (rejected: complexity remains in the least reusable layer).

## ADR-078: Keep The Web UI Git-Diff Drawer Explicitly Dismissed And Stable Across Summary Polling

- Status: Accepted
- Date: 2026-04-04
- Context:
  - ADR-077 introduced the right-side git-diff file drawer, but follow-up usage showed two UX regressions in the embedded SPA:
    - after opening a file preview, clicking elsewhere in the workspace could close the drawer because both widget-level `focusout` and document-level click handlers treated it like a transient popover.
    - the existing `/git-diff` polling path re-rendered the open drawer and force-refetched its selected file detail, causing visible flicker, content jumping, and scroll resets even when the user was still reading the same file.
  - the summary chip and expanded file list still need fresh repository state, but the opened drawer should behave like a stable inspection surface rather than a poll-driven tooltip.
- Decision:
  - stop using outside-click/focus-leave behavior to dismiss the git-diff drawer.
  - collapsing or re-expanding the git-diff summary chip must not dismiss the already open right-side drawer.
  - the right-side drawer itself is dismissed only from its own close button, not from summary-chip toggles or keyboard escape handling.
  - keep `/git-diff` polling responsible for the summary chip and expanded changed-file list only.
  - once a file drawer is opened, render it from a browser-local preview snapshot keyed by thread-session scope, and update that snapshot only when the user selects another file or that file-detail request itself changes state.
  - every explicit file-row selection must reissue the backend `git-diff-file` request for that path; reopening the same file later must not reuse a stale cached detail payload from a previous open.
  - keep the drawer's line-content presentation visually lightweight: no horizontal row separators, and use the same mono typography treatment as the composer footer's model label.
  - because the drawer content is rendered inside `<code>` nodes, explicitly force those nodes to inherit the parent font settings so the browser's default code font cannot override the intended model-label typography.
  - keep drawer rows compact by reducing vertical padding/line-height rather than shrinking the text itself.
  - for tracked unified diffs, derive visible line numbers from each hunk header and render both old/new line-number columns instead of a synthetic monotonically increasing row counter.
  - keep hunk headers visible as structural markers, but do not assign them line numbers; hide raw diff header metadata lines (`diff --git`, `index`, `---`, `+++`, and similar pre-hunk patch headers) from the drawer content.
  - serialize `/git-diff-file` as grouped rendered `blocks[]` rather than duplicating both raw `content` and per-line `lines[]`; adjacent rows with the same tone are compacted into one block containing `text[]` plus only the relevant old/new line-number arrays.
  - drop the old `showLineNumbers` response field and let clients infer visible columns from the presence of `oldLineNumbers[]` / `newLineNumbers[]` for each rendered block.
  - keep server-side parsing cheap under parallel requests by scanning the already produced diff/file text once per request in memory, and perform block compaction within that same pass without adding any extra git subprocesses beyond the existing file-detail fetch.
  - keep the preview snapshot browser-local and session-scoped; do not persist it to backend thread/session metadata.
- Consequences:
  - users can inspect a file diff while clicking elsewhere in the workspace without losing the drawer.
  - users can collapse the changed-file list to recover vertical space without losing the file they are reading in the right-side drawer.
  - periodic summary refreshes no longer rebuild the open drawer DOM or force a loading-state flash for the currently viewed file.
  - reselecting a file may briefly show the loading state again, but the drawer content now always reflects a fresh backend read instead of a reused earlier response.
  - switching sessions within the same thread no longer leaks one session's open drawer into another session's chat pane; each session restores only its own last local preview snapshot.
  - `/git-diff-file` responses are materially smaller because they no longer carry duplicate raw text plus one JSON object per rendered row.
- Alternatives considered:
  - keep outside-click dismissal and try to special-case a few targets (rejected: brittle and still wrong for a reading surface meant to stay open).
  - continue rendering the drawer directly from the live polled diff state (rejected: keeps poll cadence coupled to preview stability and scroll position).
  - persist drawer-open state to backend thread/session metadata (rejected: this is local presentation state, not shared runtime/domain state).

## ADR-077: Preview Git-Diff File Details In A Right-Side Drawer, With Text-Only Support For New Files

- Status: Accepted
- Date: 2026-04-03
- Context:
  - ADR-071 already established the session-scoped git diff summary chip and expandable file list above the Web UI composer.
  - product now requires the next interaction step after expanding that chip:
    - clicking a changed file should open its concrete diff/details.
    - switching between files should reuse one panel instead of stacking multiple overlays.
    - newly created text files should still be previewable even though `git diff` has no tracked patch for an untracked path.
  - the embedded SPA remains browser-only, so all git/file access still has to happen server-side against the validated thread `cwd`.
- Decision:
  - add `GET /v1/threads/{threadId}/git-diff-file?path=...` as a thread-scoped read-only detail endpoint alongside the existing summary endpoint.
  - keep summary parsing/backend authority in `internal/gitutil` and extend per-file rows with `viewable`, so the frontend can disable unsupported rows before the user clicks them.
  - for tracked text rows, return raw patch content from `git --no-pager diff -- <path>`.
  - for untracked text rows, return direct file contents instead of synthesizing a fake patch.
  - keep binary/non-text rows non-previewable and render them as disabled UI state instead of attempting inline binary viewers or downloads from the git-diff surface.
  - render the detail payload in one browser-local right-side drawer whose open/selected-file state is scoped to the active thread session and never persisted to backend thread metadata.
- Consequences:
  - users can inspect one changed file without leaving the conversation or opening separate modal layers.
  - new text files are previewable even before they are staged or committed.
  - the feature remains local-first and safe because path resolution stays repo-relative on the backend and non-text files are fail-closed.
- Alternatives considered:
  - inline-expand raw diff text under each row in the chip panel (rejected: too tall/noisy above the composer and awkward for file switching).
  - open a full-screen modal (rejected: heavier than necessary and visually disconnected from the existing git chip surface).
  - attempt to preview arbitrary binary files (rejected: higher complexity, broader security/UI surface, and not required for the requested workflow).

## ADR-076: Make Grouped Thread Session Lists Collapsible From The Leading Agent Glyph

- Status: Accepted
- Date: 2026-04-03
- Context:
  - ADR-074 already merged thread and session browsing into one grouped left rail, but every thread group always rendered its inline session list expanded.
  - once several threads each expose multiple recent sessions, the rail becomes visually dense and harder to scan, especially when the user only needs the thread header as a waypoint into the main chat pane.
  - product now requires a Codex-style interaction where the thread's leading agent glyph also acts as the disclosure control:
    - expanded at rest still looks like the provider avatar.
    - hover/focus reveals a down-chevron affordance for collapse.
    - collapsed groups keep a right-chevron visible so the hidden sessions remain discoverable.
- Decision:
  - keep the backend thread/session APIs, paging, and runtime model unchanged; this is a Web UI-only behavior layered on top of the existing grouped rail.
  - treat collapse state as browser-local per-thread UI state keyed by `threadId`, independent from thread selection, session binding, or provider session data.
  - collapse only the inline session-list region for that thread group; keep the thread header itself, `New session`, and overflow actions available even while the group is collapsed.
  - automatically reopen a collapsed group when the user explicitly starts `New session` from that thread header so the fresh session can surface immediately.
- Consequences:
  - dense thread/session rails become easier to scan without adding a second navigation drawer back into the layout.
  - the leading glyph now carries both branding and disclosure affordance, so the UI stays compact without adding another dedicated expand/collapse column.
  - because the state is browser-local UI state and not persisted server-side, full-page reloads currently reset all groups back to expanded.
- Alternatives considered:
  - add a dedicated disclosure column before every thread (rejected: adds visual weight and wastes horizontal space in an already compact rail).
  - persist collapsed state in backend thread metadata (rejected: this is presentation state, not shared thread/domain state).
  - collapse groups automatically based on active selection only (rejected: removes user control and makes the rail jump unexpectedly during navigation).

## ADR-075: Decorate Thread Session-List Responses With Cross-Client Active-Session State

- Status: Accepted
- Date: 2026-04-03
- Context:
  - ADR-069 already exposed thread-level `hasActiveSession`, which lets a fresh browser see that some session on a thread is live.
  - the grouped left rail introduced by ADR-074 renders concrete session rows and already has a session-row spinner, but before this change that spinner only reflected browser-local `streamStates`.
  - product now requires another browser opening the same ngent instance to immediately tell which specific session row is streaming, without first selecting that session and reconstructing its running turn state in the chat pane.
- Decision:
  - keep runtime ownership of active-turn truth in `internal/runtime.TurnController`; do not add a second persisted/session-active cache.
  - extend `GET /v1/threads/{threadId}/sessions` so each returned session row may include `isActive=true` when that concrete `(thread, session)` scope currently has a live turn.
  - continue prepending the thread's currently bound `sessionId` ahead of provider-listed sessions before computing `isActive`, so stale upstream `session/list` catalogs do not hide the active session row.
  - in the embedded Web UI, merge server-reported `session.isActive` with the existing browser-local `streamStates` set so the same row spinner works both for locally started turns and for active sessions discovered from another browser on first load/refresh.
- Consequences:
  - a newly opened or manually refreshed secondary browser can spot the exact active session row directly in the grouped rail.
  - session-row loading UI now has one shared rendering path regardless of whether the active state came from local streaming state or the server's session-list response.
  - already-open background browsers still do not receive push-based rail updates for another browser's newly started session until they reload or re-fetch that thread's session list.
- Alternatives considered:
  - add a dedicated push channel for session-list activity state (rejected: heavier protocol/UI coordination than needed for the requested "open another browser and see it" behavior).
  - infer the active row purely from thread-level `hasActiveSession` (rejected: does not identify which concrete session is active).
  - persist session-active flags in SQLite (rejected: runtime already owns the authoritative live-turn state and persistence would introduce stale-state cleanup risk).

## ADR-074: Merge Thread And Session Browsing Into One Grouped Left Rail

- Status: Accepted
- Date: 2026-04-03
- Context:
  - ADR-049 and ADR-050 moved session browsing onto the left side of the workspace, but the resulting two-column navigation still kept thread selection and session selection visually separate.
  - product now requires the left side to behave closer to Codex App's grouped project list:
    - one left rail only.
    - each thread/project header followed immediately by its recent sessions.
    - no dedicated session drawer and no separate collapse affordance.
  - the existing backend session APIs, paging contract, and thread/session runtime model remain valid and do not need schema changes.
- Decision:
  - keep `GET /v1/threads/{threadId}/sessions`, `nextCursor`, `PATCH /v1/threads/{threadId}` session binding, and fresh-session scope semantics exactly as they are.
  - in the embedded Web UI, render one left navigation rail where each thread row becomes a grouped header and its ACP session rows render inline beneath that header.
  - remove the dedicated session sidebar column and the chat-edge hover toggle used to collapse/expand it.
  - keep per-thread `New session`, refresh, and live title updates inside each grouped thread block, but expose refresh from the thread overflow menu instead of as a dedicated inline control.
  - remove persistent selected styling from thread headers/groups; only session rows carry the active selection state.
  - use the provider/agent icon as the leading grouped-thread glyph instead of a generic folder icon.
  - cap the initially visible rows to 5 sessions per thread in the Web UI and use `Show more` to reveal the next chunk while still honoring backend `nextCursor` pagination when additional provider pages exist.
- Consequences:
  - thread and session context are scanned together in one place, so switching to another thread's historical session is a single visual/interaction step.
  - chat width increases because the middle session drawer no longer reserves a second navigation column.
  - thread headers now read as neutral group labels rather than a second competing selection state next to the actual selected session.
  - the backend remains unchanged; the grouping and 5-row initial cap are frontend concerns layered on top of the existing session-list API.
- Alternatives considered:
  - keep the two-column left layout and only restyle it (rejected: does not satisfy the requested Codex-style grouped browsing model).
  - keep the separate session drawer but auto-expand it for every thread (rejected: still preserves the wrong interaction model and wastes horizontal space).
  - move grouping/paging into a new backend endpoint (rejected: unnecessary because current thread/session APIs already provide the required data).

## ADR-073: Render Live ACP Plan Updates As An Ephemeral Bottom Overlay Instead Of Transcript History

- Status: Accepted
- Date: 2026-04-03
- Context:
  - ADR-033 already promoted ACP `plan_update` to a first-class SSE/history signal, and the Web UI rendered the latest plan inside the assistant message itself.
  - in long-running turns, that placement forces users to scroll back toward the top of the reply to see whether the agent has advanced to the next step.
  - product now requires the plan to behave more like a live session-status surface:
    - always visible at the bottom of the active chat while the turn is running.
    - gone once that turn finishes.
    - still visible to refreshed or secondary browsers that attach to the same active session.
- Decision:
  - keep `plan_update` persistence and replay exactly as-is on the backend; do not add a new session-plan API or separate runtime cache.
  - in the embedded Web UI, treat the latest running-turn plan as streaming-only session UI state sourced from `streamPlanByScope`, hydrated from persisted running-turn events during history replay and then kept fresh via the resumed per-turn SSE stream.
  - render that live plan in a dedicated bottom-floating card outside the transcript list, and reserve dynamic bottom inset space in the message list so the overlay does not cover the latest messages or the scroll-to-bottom affordance.
  - stop rendering plan sections inside finalized assistant transcript bubbles, so the card disappears after completion/error/cancel just like other streaming-only session affordances.
- Consequences:
  - users can watch plan progress continuously during long replies without leaving the bottom of the conversation.
  - multi-browser/session-refresh recovery keeps working without any protocol change because the live card reuses persisted `plan_update` history plus `GET /v1/turns/{turnId}/events?after=<seq>`.
  - completed transcript history becomes cleaner, but the final plan snapshot is no longer visible after the turn ends unless a future product requirement explicitly brings back historical plan browsing.
- Alternatives considered:
  - keep plan inside the assistant message and rely on auto-scroll/top anchors (rejected: still hides progress in long replies).
  - persist/render the last plan snapshot inside finalized history as well as the live overlay (rejected: conflicts with the requested streaming-only lifecycle).
  - add a separate backend-maintained session-plan endpoint or SSE channel (rejected: duplicates existing turn-event persistence and replay mechanisms).

## ADR-072: Render Git-Diff Expanded File Rows With Locally Vendored Suffix-Based File Icons

- Status: Accepted
- Date: 2026-04-02
- Context:
  - ADR-071 added a session-scoped git diff summary chip and expandable per-file list above the Web UI composer.
  - product now requires those expanded file rows to show recognizable file-type icons, while the embedded SPA must stay local-first and avoid runtime icon fetches from GitHub or external CDNs.
  - the requested upstream source, `file-icons/vscode`, ships its icon theme as font glyphs plus palette metadata rather than a small SVG bundle.
- Decision:
  - locally vendor only the `woff2` font assets needed by the current git-diff surface and keep the upstream MIT license text in-repo with those files.
  - resolve icons in the frontend from a curated basename/extension map:
    - basename-first for special files such as `Dockerfile`, `.dockerignore`, `go.mod`, `go.sum`, `.bashrc`, and `.zshrc`.
    - otherwise by the final lowercase file extension, so names such as `test.py` still resolve to the Python icon.
  - render those glyphs inside subtle tinted tiles whose border/background derive from the icon color plus the active theme surface, and fall back to the existing generic file icon when no mapping exists.
- Consequences:
  - the expanded git-diff panel now reads more like a code workbench while remaining fully local and theme-aware.
  - unknown or niche file types still render safely with the generic file icon instead of a missing/broken asset state.
  - icon coverage is intentionally curated, so expanding support later is a frontend-only mapping change.
- Alternatives considered:
  - load icons directly from GitHub/raw URLs at runtime (rejected: weaker offline behavior, more latency, and an unnecessary external dependency).
  - vendor the entire upstream icon theme and all mapping tables (rejected: much heavier than this narrow UI surface needs).
  - draw a bespoke inline SVG icon set from scratch (rejected: slower to build and less aligned with the explicitly requested upstream visual source).

## ADR-071: Surface Session-Selected Git Diff Summaries Above The Web UI Composer, Including Untracked Files

- Status: Accepted
- Date: 2026-04-02
- Context:
  - ADR-066 already established thread-level git repository inspection plus an in-composer branch pill/switcher.
  - product now also requires a Kimi-style diff summary above the composer that reflects current working-tree changes for the selected session.
  - the UI must not show anything for non-git directories or when the host does not have `git`, and the frontend must own polling instead of receiving pushed diff updates.
- Decision:
  - add `GET /v1/threads/{threadId}/git-diff` as a lightweight thread-scoped git capability endpoint.
  - back the endpoint with direct host git commands:
    - `git --no-pager diff --shortstat`
    - `git --no-pager diff --numstat`
    - `git ls-files --others --exclude-standard -z`
  - parse those commands server-side into structured JSON (`summary`, per-file rows, `repoRoot`) so the frontend stays presentation-focused and does not need to parse raw git output.
  - treat missing `git` binaries and non-repository `cwd` values as optional capability absence by returning `available=false` instead of surfacing a hard failure.
  - keep the Web UI-side gating so polling still happens only when the user has chosen a concrete historical/live session, but do not require or echo `sessionId` at the API layer because the repository identity already comes from the thread `cwd`.
  - in the embedded Web UI, poll every 15 seconds only for the active thread's selected session, force an immediate fetch on session switch, and hide the surface entirely when `available=false` or there are no tracked/untracked rows.
- Consequences:
  - the composer now exposes repository change pressure in the active session without adding backend push complexity or extra SSE events.
  - git diff rendering remains consistent across browsers because parsing happens once on the backend.
  - repositories with no tracked modifications still surface newly created untracked files, so the chip matches the user's visible working tree more closely.
  - repositories with no visible modifications stay visually clean because the chip is omitted.
- Alternatives considered:
  - push git diff changes over SSE after every tool/file mutation (rejected: much higher coupling to agent behavior and no reliable signal for out-of-band filesystem edits).
  - have the frontend shell out or parse raw git output directly (rejected: impossible in-browser and would duplicate parsing logic).
  - use `git status --short` instead of `git diff --shortstat/--numstat` (rejected: product explicitly wants diff-based counts aligned with Kimi's presentation).

## ADR-070: Expand Embedded Web UI Localization And Repository READMEs To Spanish And French

- Status: Accepted
- Date: 2026-04-01
- Context:
  - ADR-068 already established browser-local UI localization for English and Simplified Chinese.
  - product now requires the embedded Web UI to also support Spanish and French without introducing any backend locale negotiation or server-side translation layer.
  - repository onboarding docs also need equivalent Spanish and French entry points alongside the existing English and Simplified Chinese READMEs.
- Decision:
  - expand the frontend `language` preference to four supported values:
    - `en`
    - `zh-CN`
    - `es`
    - `fr`
  - extend browser-locale detection so first load chooses the closest supported locale:
    - any `zh-*` locale => `zh-CN`
    - any `es-*` locale => `es`
    - any `fr-*` locale => `fr`
    - all other locales => `en`
  - keep localization entirely client-side in the embedded Web UI store/settings flow.
  - add `README.es.md` and `README.fr.md`, and cross-link all root README variants from each translated entry point.
- Consequences:
  - the embedded Web UI now covers four languages without changing HTTP/SSE/API contracts.
  - users whose browsers prefer Spanish or French now get a better first-load default instead of falling straight back to English.
  - repository landing docs are now available in English, Simplified Chinese, Spanish, and French.
  - backend/provider-originated error payloads still remain untranslated and can still produce mixed-language surfaces.
- Alternatives considered:
  - introduce region-specific stored variants such as `es-ES`, `es-MX`, `fr-FR`, and `fr-CA` (rejected: unnecessary complexity for the current UI copy scope).
  - add automatic machine-translated README generation in CI (rejected: poorer reviewability and weaker quality control over user-facing docs).
  - keep README translations limited to English and Simplified Chinese while only translating the UI (rejected: incomplete onboarding for the requested languages).

## ADR-069: Keep Active Turns Independent Of Individual SSE Viewers And Resume Them Through Per-Turn Event Streams

- Status: Accepted
- Date: 2026-04-01
- Context:
  - the original `POST /v1/threads/{threadId}/turns` SSE request previously owned the turn lifetime, so refreshing the browser or losing that viewer connection cancelled the active turn even when the user did not ask to stop it.
  - product now requires the opposite behavior: only explicit cancel, provider terminal state, or fail-closed permission timeout should end a healthy in-flight turn.
  - refreshed browsers and other concurrent browsers also need to attach to the same active turn, receive the remaining live deltas in order, and keep the same active session visible.
  - ngent already persisted turn events, but write-time delta merging weakened a strict `afterSeq` resume contract for live replay.
- Decision:
  - decouple turn execution lifetime from the original SSE viewer connection by running the turn under a context that survives request disconnect; explicit cancel and existing fail-closed permission timeout behavior remain intact.
  - add `GET /v1/turns/{turnId}/events?after=<seq>` as a per-turn SSE endpoint that:
    - subscribes to a live in-memory broker for that turn.
    - replays persisted events with `seq > after`.
    - continues tailing newly published events until terminal completion.
  - emit a first-class `permission_resolved` turn event whenever a pending permission is decided or times out, so every viewer of that turn converges on the same permission state.
  - expose `hasActiveSession` on thread list/get responses so a newly opened browser can tell which thread currently has a live session before it opens history for that thread.
  - keep turn-event storage append-only on write, including consecutive delta events, and expose monotonic per-turn `seq` values in replayable SSE payloads.
  - update the embedded Web UI to hydrate already-running turns from persisted history and then reattach via the new per-turn stream, so refresh and multi-browser viewing preserve the same active conversation.
- Consequences:
  - browser refreshes and other non-user transport disconnects no longer cancel healthy turns.
  - multiple browsers can observe one running turn concurrently and receive the same ordered live event stream.
  - SQLite now stores more raw delta rows than the earlier write-time merge approach, but resume correctness improves; any delta coalescing stays a read-path concern only.
  - turn cancellation semantics are clearer: only explicit cancel, timeout, or provider terminal state stop the turn.
- Alternatives considered:
  - keep request-bound execution and require the browser to resubmit on refresh (rejected: duplicates work and violates the requested UX).
  - use one shared global SSE channel for all turns (rejected: worse isolation and unnecessary cross-turn fanout noise).
  - keep write-time delta merging and resume from timestamps instead of sequence numbers (rejected: weaker ordering guarantees and more ambiguous recovery behavior).

## ADR-068: Add Browser-Default English/Simplified Chinese Localization To The Embedded Web UI

- Status: Accepted
- Date: 2026-04-01
- Context:
  - the embedded Web UI was English-only even though it is now used as the primary local workstation for a broader set of users.
  - product required one new language immediately: Simplified Chinese.
  - users also needed an explicit per-browser language override in Settings instead of being locked to the first detected locale forever.
  - the existing frontend is a no-framework SPA, so the change needed to stay client-local and avoid backend/session contract changes.
- Decision:
  - add a persisted browser-local `language` preference to the frontend store with two supported values:
    - `en`
    - `zh-CN`
  - on first load, when no explicit preference is stored yet:
    - map any browser locale matching `zh-*` to `zh-CN`
    - fall back to `en` for every other locale
  - move language selection into the Settings drawer and keep it entirely client-side, alongside theme/auth/server URL preferences.
  - re-render the Web UI shell when language changes so visible chrome updates immediately without a page reload.
  - localize client-owned UI strings and relative-time formatting only; pass through backend/provider text unchanged for now.
- Consequences:
  - the first visit now feels locale-aware without introducing a third `system/follow browser` mode into long-term settings state.
  - different browsers/profiles on the same machine can intentionally keep different UI languages because the setting lives in each browser's local storage.
  - unsupported browser locales degrade safely to English.
  - mixed-language surfaces are still possible when the backend returns English error text; that remains an accepted limitation for now.
- Alternatives considered:
  - add a persistent `system` language mode that tracks browser locale on every visit (rejected: unnecessary complexity for the current two-language scope).
  - localize on the backend and send translated text over HTTP/SSE (rejected: outside the current UI-only scope and would complicate protocol contracts).
  - auto-map only exact `zh-CN` and leave other `zh-*` locales in English (rejected: worse default experience for Chinese-language browsers).

## ADR-067: Persist ACP Session Usage Snapshots And Surface Context-Window Pressure In The Web UI

- Status: Accepted
- Date: 2026-03-30
- Context:
  - ACP's session-usage RFD adds two related but different signals:
    - cumulative prompt-token usage on `session/prompt` responses.
    - `session/update` `usage_update` snapshots with current context-window `used` and `size`.
  - product needed both:
    - sqlite persistence so session usage can be queried again later.
    - a Codex-style context-pressure indicator in the embedded Web UI.
  - provider support is partial today, so the feature must fail soft when usage is absent instead of showing placeholders or misleading zeros.
- Decision:
  - treat ACP session usage as a cumulative per-session snapshot keyed by `(agent_id, cwd, session_id)`.
  - persist those snapshots in sqlite `session_usage_cache`, storing:
    - token totals (`total/input/output/thought/cached_*`)
    - current context-window `used/size`
    - optional cost amount/currency
  - merge partial updates with `COALESCE` semantics on conflict so providers can send prompt-response totals and later `usage_update` context data independently without clearing earlier fields.
  - emit a first-class `session_usage_update` event into SSE/history whenever usage arrives during a turn.
  - expose `GET /v1/threads/{threadId}/session-usage?sessionId=...` so the Web UI and future tooling can query the latest stored snapshot without replaying all history.
  - render the Web UI badge only when `contextUsed` and `contextSize` are both present and `contextSize > 0`; otherwise hide the badge entirely.
- Consequences:
  - session usage survives page reloads and server restarts.
  - the browser can hydrate usage cheaply from sqlite and then refine it with live SSE updates during the current turn.
  - providers that only expose token totals still contribute useful stored data, but the UI remains silent because a context-window percentage would be incomplete.
  - provider gaps remain visible only by absence, not by noisy error chrome.
- Alternatives considered:
  - reconstruct usage only from persisted turn events in the browser (rejected: repeated replay work and no simple single-snapshot query surface).
  - show a placeholder badge for every session even when providers do not emit usage (rejected: adds UI noise and implies unsupported precision).
  - persist only total token counts and ignore context-window size (rejected: does not satisfy the requested percentage indicator).

## ADR-066: Surface Thread-Scoped Git Branch State In The Web UI Composer

- Status: Accepted
- Date: 2026-03-30
- Context:
  - users wanted the embedded Web UI to show the active repository branch near the composer, similar to Codex App, and to allow fast local branch switching from that same surface.
  - ngent threads already carry one validated `cwd`, so git state needs to be resolved per thread rather than globally for the whole server.
  - the feature must stay optional:
    - if the host does not have `git`, do not show anything.
    - if the thread `cwd` is not inside a repository, do not show anything.
  - branch checkout mutates the shared worktree for the thread, so it must respect ngent's existing whole-thread shared-state locking model instead of racing active turns.
- Decision:
  - add `internal/gitutil` as the small host-integration layer that:
    - detects `git` availability via `exec.LookPath("git")`.
    - inspects repository root, current ref, detached-head state, and local branches for one `cwd`.
    - checks out only existing local branches.
  - expose thread-scoped HTTP APIs:
    - `GET /v1/threads/{threadId}/git`
    - `POST /v1/threads/{threadId}/git` with `{branch}`
  - make the git capability fail-soft for inspection:
    - missing git binary or non-repository cwd returns `200` with `available=false`.
    - the Web UI uses that as a signal to hide the branch control entirely.
  - make branch checkout a thread-exclusive operation:
    - `POST /git` acquires the runtime's thread-wide exclusive guard.
    - if any turn is active on that thread, return `409 CONFLICT`.
  - render the control in the composer footer only when `available=true`, list only local branches, and refresh git state when the thread is opened and after turns settle.
- Consequences:
  - repository-aware threads now expose one compact branch affordance without adding noise to non-git or gitless environments.
  - branch switching honors the same concurrency guarantees as delete/compact and other shared-state thread operations.
  - the server API surface now owns git worktree introspection, so the browser does not need direct filesystem or shell access.
- Alternatives considered:
  - resolve git state entirely in the browser or via shell-side hacks (rejected: the browser cannot safely inspect host repositories directly).
  - make missing git / non-repository states hard API errors (rejected: the UI should quietly omit the optional affordance).
  - allow free-form checkout targets including remotes or revision expressions (rejected: the requested UX is specifically local-branch switching and the narrower contract is safer).

## ADR-065: Recast The Embedded Web UI As A Restrained Desktop Workbench

- Status: Accepted
- Date: 2026-03-30
- Context:
  - the earlier premium refresh in ADR-054 improved quality, but the resulting shell still leaned too far toward decorative glass-panel SaaS styling rather than a serious local desktop tool.
  - product feedback for this iteration explicitly asked for a colder, denser, calmer workbench feel:
    - no ASCII brand lockup.
    - no background grids/orbs/glow-heavy blur.
    - less pill/card-kit styling.
    - a more precise hierarchy around working directories, threads, sessions, and streamed output.
  - the runtime and protocol behavior were already correct, so the redesign needed to stay UI-only and preserve the existing store/SSE/DOM interaction model.
- Decision:
  - keep all backend, protocol, and streaming behavior unchanged; this redesign is limited to frontend structure, typography, tokens, and interaction treatment.
  - replace the prior glass-heavy direction with a restrained desktop-workbench system built on:
    - neutral graphite/light-steel surfaces.
    - a single restrained blue accent reserved for actions, focus, and active markers.
    - tighter radii, flatter shadows, and stronger border hierarchy.
    - an offline-friendly local/system font stack with a monospace secondary voice for tool metadata.
  - make the chat workspace the clear visual center:
    - left thread rail stays present but quieter.
    - session rail becomes an auxiliary context column rather than a third heavy primary container.
    - empty states anchor to bounded panels instead of floating in large whitespace.
  - make working directory selection the primary field in the thread-creation flow, and visually demote advanced JSON agent options.
  - unify reasoning, plan, tool-call, markdown, and permission blocks under one document-like section language instead of distinct card families.
- Consequences:
  - the Web UI reads more like a mature local development tool and less like a decorative AI SaaS demo.
  - interaction semantics remain stable because ids, bindings, and stream/history contracts were preserved.
  - the visual system is now less dependent on blur/backdrop support and therefore degrades more predictably across browsers.
  - ADR-054's glass-heavy styling direction and the earlier shared ASCII browser-branding experiment are no longer current guidance for the embedded Web UI.
- Alternatives considered:
  - keep the previous shell and only tweak color/radius values (rejected: too shallow for the requested direction change).
  - push further into animated/editor-chrome spectacle (rejected: would conflict with the requested calm, professional tool posture).
  - add bundled webfonts or a component framework to chase richer presentation (rejected: unnecessary for the requested result and inconsistent with the local-first no-framework Web UI constraints).

## ADR-054: Refresh The Embedded Web UI As A Premium Workbench Without Changing Behavior

- Status: Superseded by ADR-065
- Date: 2026-03-26
- Context:
  - the existing embedded Web UI was functionally complete, but the visual quality still read like an internal tool prototype.
  - users explicitly asked for a more premium, higher-quality desktop feel while keeping the product local-first and operationally unchanged.
  - the frontend milestone is already on a no-framework Vite + TypeScript SPA, so the refresh needed to stay within the existing DOM/store/SSE architecture.
- Decision:
  - keep all runtime logic, API usage, SSE semantics, and store behavior unchanged; the refresh is limited to markup hierarchy, styling tokens, and interaction polish.
  - adopt a glass-panel workbench direction with:
    - layered backdrop treatment.
    - elevated sidebars and chat surface.
    - stronger typographic hierarchy in headers and empty states.
    - richer composer, modal, drawer, and permission-card presentation.
  - use a restrained teal accent and warmer neutrals instead of the earlier default blue-on-flat-gray presentation.
  - rely on a curated system-font stack instead of bundling a remote webfont, so the embedded UI remains self-contained and offline-friendly.
- Consequences:
  - the UI feels more intentional and mature without adding framework/runtime complexity.
  - visual output remains stable at the interaction/API level because business logic did not move.
  - exact typography can vary slightly by host OS because the stack prefers locally available premium system fonts.
- Alternatives considered:
  - keep the current structure and only tweak colors (rejected: insufficient improvement).
  - introduce a JS component framework for richer visuals (rejected: conflicts with the repo's no-framework Web UI direction).
  - bundle external webfonts/assets from a CDN (rejected: weakens the local-first/offline posture).

## ADR-053: Replace `slog` JSON Output With A Human-Readable Stderr Logger And Colored Access Logs

- Status: Accepted
- Date: 2026-03-23
- Context:
  - the previous stderr logger emitted JSON `slog` envelopes for every runtime event and HTTP request.
  - operators reviewing ngent locally primarily read logs in a terminal, and the JSON shape made common request traffic harder to scan quickly.
  - the product still needs the existing safety properties:
    - stderr-only logging.
    - leveled output with debug gating.
    - secret redaction in errors and ACP traces.
- Decision:
  - replace the shared `slog` logger with a repo-local human-readable logger in `internal/observability`.
  - keep the current logging call sites on a simple leveled API (`Debug/Info/Warn/Error`) to minimize churn outside the observability package.
  - emit HTTP completion logs in a compact access-log shape:
    - `INFO: <local-time> <client-ip> - "<method> <path> <proto>" <status> <statusText> <duration>`
  - keep ACP debug tracing behind `--debug=true`, but print it as readable text fields with the sanitized RPC payload embedded as JSON.
  - enable ANSI colors only when stderr is attached to a TTY, so redirected output remains plain text.
- Consequences:
  - local terminal logs are easier to scan, especially for high-volume request traffic.
  - stderr output is no longer one-JSON-object-per-line, so external log collectors may need a text parser if they ingest ngent logs directly.
  - secret redaction remains centralized in `internal/observability`, so changing the presentation layer does not weaken the fail-closed logging policy.
- Alternatives considered:
  - keep JSON `slog` and only add a prettier startup banner (rejected: request traffic would remain noisy).
  - keep `slog` but swap in a custom pretty handler (rejected: the desired output is simpler as a small repo-local logger than as a `slog` compatibility layer).
  - add a new pretty logger only for access logs and leave everything else as JSON (rejected: mixed formats on stderr would be inconsistent and harder to reason about).

## ADR-052: Cursor CLI ACP Provider Integration With Explicit ACP Authentication

- Status: Accepted
- Date: 2026-03-23
- Context:
  - Cursor CLI now documents ACP mode via `agent acp`, and the user already has a working local Cursor CLI install.
  - official Cursor ACP docs describe the expected session flow as `initialize -> authenticate(methodId="cursor_login") -> session/new|session/load -> session/prompt`.
  - local probing against the installed CLI confirmed that:
    - `initialize` advertises `authMethods=[{"id":"cursor_login",...}]`.
    - skipping `authenticate` causes real `session/new` to stall without returning a response.
    - `session/new.model` and `session/new.modelId` do not select the requested model.
    - `session/set_config_option("model", value)` updates Cursor's selected model and returns updated `configOptions`.
- Decision:
  - add `internal/agents/cursor` as a first-class provider on top of the shared `internal/agents/acpcli` driver.
  - start Cursor ACP with a binary fallback order of `agent acp` then `cursor-agent acp`, so ngent tolerates both common local install shapes.
  - perform provider-local ACP authentication immediately after `initialize` whenever `cursor_login` is advertised in `authMethods`.
  - treat Cursor model selection as a config-option concern rather than a startup/prompt hint:
    - do not rely on `session/new.model` / `modelId`.
    - apply selected `modelId` through `session/set_config_option("model", ...)`.
  - keep Cursor's optional ACP extension methods out of scope for now; ngent continues to rely on standard ACP `session/update`, `session/request_permission`, and `session/cancel`.
- Consequences:
  - ngent can list, configure, and run Cursor threads through the same shared ACP CLI lifecycle as the other direct ACP providers.
  - Cursor support remains robust even when only one of `agent` or `cursor-agent` is present in PATH.
  - thread-selected Cursor models now behave consistently with real Cursor ACP semantics instead of silently being ignored at `session/new`.
- Alternatives considered:
  - treat Cursor like the other direct ACP providers and skip `authenticate` (rejected: local probing showed `session/new` hangs).
  - key model selection off `session/new.model` or `session/prompt.model` only (rejected: real Cursor ignores those hints).
  - build a Cursor-specific lifecycle stack instead of reusing `acpcli` (rejected: the delta is limited to auth and model-selection hooks).

## ADR-051: BLACKBOX AI ACP Provider Integration Via Shared ACP CLI Driver

- Status: Accepted
- Date: 2026-03-22
- Context:
  - BLACKBOX AI CLI now exposes ACP mode through `blackbox --experimental-acp`, and the user already has the local CLI installed/configured.
  - local probing against `blackbox 1.2.47` showed that the CLI can emit non-JSON stdout noise (process-info and telemetry-related lines) before or between ACP frames.
  - the same probing showed current ACP capability limits:
    - `initialize` advertises `agentCapabilities.loadSession=false`.
    - real `session/load` returns `-32601 method not found`.
    - `session/new` currently does not expose `models.availableModels` or `configOptions`.
  - the CLI can start ACP without explicitly passing `BLACKBOX_API_KEY` when another local auth method is already configured, but prompt execution still depends on valid upstream auth/network readiness.
- Decision:
  - add `internal/agents/blackbox` as a first-class provider on top of the shared `internal/agents/acpcli` driver.
  - start BLACKBOX with `blackbox --experimental-acp`, plus hub-side compatibility flags `--skip-update` and `--telemetry=false`.
  - enable `acpstdio.ConnOptions.AllowStdoutNoise` for BLACKBOX so provider stdout noise does not corrupt the JSON-RPC stream.
  - forward thread-selected `modelId` through both process startup (`--model`) and ACP request hints (`session/new.model` / `modelId`, `session/prompt.model`).
  - keep session browsing/replay unavailable until BLACKBOX exposes standard ACP resume surfaces.
- Consequences:
  - ngent can now run BLACKBOX turns through the same direct ACP path as Qwen/OpenCode/Gemini/Kimi without introducing another custom lifecycle stack.
  - BLACKBOX threads remain usable for normal multi-turn conversation through ngent's own persisted history/context injection even though provider-owned `session/load` is not available.
  - the Web UI/API will not show resumable BLACKBOX session replay or model catalogs until upstream ACP surfaces them.
- Alternatives considered:
  - delay BLACKBOX support until `session/load` and model discovery are fully available upstream.
  - build a provider-local non-ACP integration path just for BLACKBOX.
  - require `BLACKBOX_API_KEY` unconditionally at ngent startup instead of reusing whichever auth method the local CLI already has configured.

## ADR-050: Keep The Left Agent Rail Permanently Expanded

- Status: Accepted
- Date: 2026-03-19
- Current-status note: the permanently expanded left-rail part of this ADR still stands, but ADR-074 later removed the separate session drawer entirely and merged session browsing back into the same grouped left rail.
- Context:
  - ADR-049 introduced a collapsible compact agent rail to mimic OpenCode's left-most project strip more closely.
  - in follow-up product review, that compact state was judged less useful than expected because the ngent left column represents full agent/thread items rather than tiny project icons, and collapsing it hid search plus thread metadata too aggressively.
  - the session panel still benefits from collapsibility because it is secondary context and can be hidden to reclaim chat width.
- Decision:
  - keep the left agent rail permanently expanded on desktop and mobile overlay states.
  - remove the agent-rail collapse/expand trigger and compact monogram-only rendering path.
  - retain the session panel as the only navigation-width toggle, but expose its desktop collapse/expand affordance from a hover-revealed button on the chat panel's left edge instead of from the panel header itself.
  - keep the session panel contextual to the current selection: if no thread is active yet, do not render a placeholder session column.
- Consequences:
  - the main navigation always exposes thread metadata and thread actions without an extra click.
  - layout stays simpler because only one left-side panel now owns collapse state.
  - first-load navigation density improves because chat sits directly beside the agent rail until a thread is chosen.
  - collapsing the session panel now fully retracts it instead of leaving behind a narrow visible strip.
  - the left rail no longer mirrors OpenCode's narrow icon strip exactly, but preserves the more useful ngent-specific thread browsing surface.
- Alternatives considered:
  - keep both rails collapsible (rejected: too much state and weaker scanning ergonomics for ngent's denser thread rows).
  - remove all collapse controls entirely (rejected: the session panel still benefits from being dismissible when chat width matters).

## ADR-049: Align Web UI Navigation With A Left Agent Rail And Left Session Panel

- Status: Accepted
- Date: 2026-03-19
- Current-status note: this ADR established the first left-side session-navigation move, but its two-column left layout was later superseded by ADR-074, which merged thread and session browsing back into one grouped left rail.
- Context:
  - the Web UI had been using a wide left thread list plus a separate right session sidebar.
  - users wanted the navigation model to feel closer to OpenCode's web UI, where project/session browsing sits on the left side of the workspace and the first column can collapse into a compact rail.
  - the old layout also forced session context away from the currently selected agent/thread, making it harder to see which project path and session set were active together.
- Decision:
  - render navigation as two left-side columns: a collapsible agent/thread rail and a collapsible session panel between that rail and the chat area.
  - default the agent rail to its collapsed state; collapsed items use the displayed agent/thread title's first character, uppercasing ASCII letters and preserving non-Latin first characters directly.
  - keep the full agent list, search box, and thread action menu available in the expanded rail.
  - move `New agent` below the agent list instead of keeping it in the rail header.
  - when the session panel is expanded, show the active agent/thread title, provider badge, project path, and a full-width `New session` action before the session list.
- Consequences:
  - desktop navigation matches the requested OpenCode-style mental model more closely and keeps both agent selection and session browsing on the same side of the workspace.
  - the compact default rail saves horizontal space while still allowing quick switching between agents/threads.
  - the session panel can still be collapsed independently when users want more room for chat content.
  - mobile behavior stays conservative: the existing small-screen sidebar overlay path remains the fallback, and the dedicated session panel still hides below the narrower desktop breakpoint.
- Alternatives considered:
  - keep the old right-side session sidebar and only restyle it (rejected: does not satisfy the requested navigation model).
  - move sessions left but keep the agent list permanently wide (rejected: wastes space and misses the requested compact rail).
  - use provider logos for collapsed rail markers (rejected: product request explicitly prefers first-character monograms).

## ADR-047: Defer Thread Config-Option Apply Until The Next Turn Boundary

- Status: Accepted
- Date: 2026-03-16
- Context:
  - the Web UI model/reasoning pickers had been calling `POST /v1/threads/{threadId}/config-options`, and that endpoint immediately pushed the selection into the cached provider via ACP `session/set_config_option`.
  - this coupled a pure UI selection change to live provider mutation, even when the user had not sent another message yet.
  - cached provider instances were also keyed by the full `agentOptions` blob, so changing `modelId` or `configOverrides` could unnecessarily rotate the provider cache instead of reusing the same thread/session runtime.
- Decision:
  - keep `POST /v1/threads/{threadId}/config-options` as a persistence-only API: validate against available config options, update sqlite thread state, and return the selected config view without mutating the live provider.
  - narrow cached provider scope to thread + session/fresh-session identity, so model/reasoning edits do not evict the existing session provider on their own.
  - when a new turn starts, compare the persisted thread selections against the cached provider's current selections and apply only the changed options immediately before `session/prompt`.
- Consequences:
  - Web UI config changes feel immediate in persisted thread state, but agent-side mutation is deferred until the user actually sends the next message.
  - cached sessions survive picker edits, so session continuity is preserved and redundant provider churn is reduced.
  - if a provider has no cached runtime yet, the next turn still starts from the persisted thread selections, so restart recovery behavior stays intact.
- Alternatives considered:
  - keep immediate apply on every picker change (rejected: mutates provider state without a user turn boundary).
  - add a separate explicit Apply button (rejected: product requirement keeps no-button picker UX).
  - persist both desired and applied config state in sqlite (rejected: the cached provider already holds current session state; the extra durable copy would add complexity without near-term product value).

## ADR-046: Collapse Finalized Web UI Thinking Panels by Default

- Status: Accepted
- Date: 2026-03-14
- Context:
  - ADR-045 made hidden reasoning visible in the Web UI, but a fully expanded reasoning block on every completed agent message made longer threads harder to scan.
  - the product still needs live reasoning to stay visible while the turn is actively streaming.
  - the message list is re-rendered from store state, so manual expand/collapse choice needs explicit local UI tracking if it should survive later store updates.
- Decision:
  - use a sparkles icon + italic label + rotating chevron trigger, with expanded content shown as indented text behind a left border.
  - use tense-sensitive labels: live reasoning stays `Thinking`, and finalized reasoning switches to `Thought`.
  - render finalized reasoning with the same sanitized markdown pipeline used for finalized assistant messages, while keeping in-flight reasoning as plain text during streaming.
  - keep the streaming panel expanded while `reasoning_delta` is still arriving.
  - once the turn is finalized into message history, render the panel collapsed by default.
  - preserve manual expand/collapse state for finalized messages in page-local UI state across later list re-renders.
- Consequences:
  - active reasoning remains visible during execution, but completed threads are denser and easier to review.
  - markdown affordances such as headings, lists, links, and code blocks now render consistently inside finalized `Thinking` content.
  - collapse state is a local Web UI concern and is not persisted into server-side turn history.
  - page reload still returns to the default product behavior: finalized thinking starts collapsed.
- Alternatives considered:
  - keep finalized reasoning always expanded (rejected: too noisy for longer chats).
  - collapse reasoning immediately even during streaming (rejected: hides the live signal users asked to watch).
  - persist expand/collapse state in backend history (rejected: presentation state does not belong in the turn/event model).

## ADR-045: Surface Hidden Agent Reasoning as First-Class SSE/History Events in the Web UI

- Status: Accepted
- Date: 2026-03-14
- Context:
  - shared ACP update parsing already recognized provider thought chunks such as `thought_message_chunk` and `agent_thought_chunk`, but the turn pipeline only forwarded visible assistant text through `message_delta`.
  - as a result, hidden reasoning emitted by supporting agents was discarded before it reached SSE clients or persisted turn history, so the Web UI could not show it during streaming or after reload.
  - merging reasoning into `responseText` would blur the boundary between visible assistant output and hidden provider reasoning.
- Decision:
  - add a context-bound reasoning callback in `internal/agents`, parallel to the existing plan/session/permission callback pattern.
  - route ACP thought chunks into that callback from the shared ACP notification handler.
  - persist and stream reasoning as a separate `reasoning_delta` event in the HTTP API layer, without changing `responseText`.
  - reconstruct reasoning in the Web UI from live `reasoning_delta` SSE events and from persisted turn `events[]`, rendering it in a dedicated `Thinking` section above the final assistant answer.
- Consequences:
  - users can see reasoning for ngent-created turns both live and from history reloads.
  - `responseText` remains the visible assistant answer only, so existing prompt-compaction/history semantics stay stable.
  - provider-owned transcript replay returned in session-scoped history still contains only visible user/assistant messages; hidden reasoning for historical external sessions is not backfilled there.
- Alternatives considered:
  - append reasoning directly into `responseText` (rejected: mixes two different product surfaces and would break existing history assumptions).
  - keep reasoning UI-only without persisting it in turn events (rejected: reload/history would lose the data).
  - add a separate top-level reasoning column to `turns` (rejected: event log already models streamed deltas cleanly and preserves ordering).

## ADR-044: Normalize Path-Like ACP Permission Previews Across Direct ACP Providers

- Status: Accepted
- Date: 2026-03-14
- Context:
  - after the shared ACP CLI refactor, direct ACP providers only bridge `session/request_permission` when they explicitly install a `HandlePermissionRequest` hook.
  - OpenCode did not install that hook, so real permission-gated turns still returned JSON-RPC `method not found` even though Kimi/Qwen had already been fixed.
  - real OpenCode permission payloads can represent file-like access without `content[]` diffs, using `toolCall.locations[]` or `toolCall.rawInput.filepath` and generic titles such as `external_directory`.
- Decision:
  - require every direct ACP stdio provider that advertises permissions to install a request-permission hook in the shared ACP CLI driver.
  - extend the shared permission parser to extract a first path-like preview from `content[]`, `locations[]`, or `rawInput` keys such as `filepath`, `path`, and `parentDir`.
  - treat directory/path-oriented permission requests as `file` approvals in the Web UI, and prefer the resolved path over generic provider titles when building the user-facing command label.
  - reuse the same shared normalization path for OpenCode instead of adding another provider-local decoder.
- Consequences:
  - real OpenCode file-creation requests now surface through the normal ngent permission card flow instead of failing invisibly.
  - future ACP providers can attach path previews in multiple shapes without forcing another adapter-specific parser.
- Alternatives considered:
  - leave OpenCode on provider-default RPC handling and only fix Kimi/Qwen (rejected: direct ACP providers would keep drifting on identical permission semantics).
  - special-case OpenCode path extraction inside `internal/agents/opencode` (rejected: the shape difference belongs in shared ACP normalization, not a one-off adapter fork).

## ADR-043: Share One ACP CLI Driver Across Kimi/Qwen/OpenCode/Gemini

- Status: Accepted
- Date: 2026-03-13
- Context:
  - `qwen`, `opencode`, `kimi`, and `gemini` all run as ACP-over-stdio providers inside ngent.
  - before this refactor, each provider duplicated the same lifecycle code for process startup, `initialize`, `session/new`, `session/load`, `session/list`, `session/prompt`, config-option probing, model discovery, and transcript replay.
  - the duplication made protocol fixes expensive and increased the cost of adding another ACP-capable CLI in the future.
  - the real differences between these providers are comparatively small:
    - startup command and environment.
    - `session/new/load/prompt` parameter shapes.
    - permission request/response encoding.
    - cancel strategy.
    - provider quirks such as Kimi local config fallback and Gemini stdout noise before JSON-RPC frames.
- Decision:
  - introduce shared package `internal/agents/acpcli` as the common driver for ACP CLI providers.
  - keep provider-specific behavior as hooks:
    - process launcher/open-connection logic.
    - ACP request parameter builders.
    - permission-response mapping.
    - cancel behavior.
    - config-session planning for provider quirks such as Kimi model selection at process startup.
  - standardize all four providers on `internal/agents/acpstdio.Conn`.
  - extend `acpstdio` with opt-in stdout-noise tolerance so Gemini can reuse the same shared transport instead of keeping its own JSON-RPC connection implementation.
  - reuse the same shared driver for model discovery and `session/load` transcript replay, not only turn streaming.
- Consequences:
  - future ACP CLI providers can usually be added by supplying a small provider spec/hook layer instead of copying full lifecycle code.
  - transport/protocol fixes now land once and benefit all ACP CLI providers together.
  - provider-specific quirks remain isolated and explicit rather than being hidden in near-identical forked implementations.
  - real-provider regressions can still come from local CLI readiness/auth/network state; sharing the driver does not hide upstream availability issues.
- Alternatives considered:
  - keep separate provider implementations and continue copy-pasting fixes.
  - introduce one fully generic provider configured only by command string/args, without explicit provider hook points.
  - keep Gemini on its own transport while only partially sharing logic for the other providers.

## ADR-042: Treat Explicit Web UI "New session" as a Fresh Turn with No Injected Thread Context

- Status: Accepted
- Date: 2026-03-13
- Context:
  - the root cause was `buildInjectedPrompt()`: ngent treats any empty `sessionId` as a local-context continuation and wraps the next prompt with prior thread turns.
  - disabling context injection for all empty-session turns would be too broad because brand-new no-session threads still benefit from local prompt continuation after earlier turns.
- Decision:
  - persist one internal thread option `_ngentFreshSession=true` only when an existing non-empty `sessionId` is explicitly cleared through the Web UI `New session` flow.
  - while `_ngentFreshSession=true` and `sessionId` is empty, bypass local context injection and send the raw user input into ACP `session/new` / `session/prompt`.
  - automatically clear `_ngentFreshSession` when a new session binds or when the user selects an explicit existing `sessionId`.
  - strip `_ngentFreshSession` from public thread responses and from session-only diff checks so the flag remains server-internal and does not change client-visible API semantics.
- Consequences:
  - explicit `New session` now means a blank fresh ACP session from the user's perspective, not merely a new durable `sessionId`.
  - ordinary threads that have no bound session for other reasons still keep existing context-window behavior.
  - stale plain-empty and fresh-empty cached providers are both evicted on explicit reset, so neither scope can silently reuse the wrong runtime.
- Alternatives considered:
  - disable context injection for every empty-session turn.
  - expose a dedicated public `session/new` endpoint and model fresh-session state outside `PATCH /v1/threads/{threadId}`.
  - leave the internal fresh-session marker in API responses.

## ADR-041: Treat Web UI "New session" as Provider-Cache Reset for the Empty Session Scope

- Status: Accepted
- Date: 2026-03-13
- Context:
  - the Web UI uses `PATCH /v1/threads/{threadId}` with `agentOptions.sessionId` cleared to represent "New session", and the next turn is expected to force ACP `session/new`.
  - ngent caches managed providers by `(threadId, normalized agentOptions)` scope.
  - if a stale provider is still cached under the empty-session scope, simply clearing `sessionId` can reuse that provider and route the next prompt into an older ACP session instead of a fresh one.
- Decision:
  - when a session-only thread update changes `sessionId` from a non-empty value to empty, evict any idle cached provider for the target empty-session scope before the next turn starts.
  - keep the current selected-session provider cache intact; only the provisional empty-session scope is force-reset.
  - keep the Web UI composer disabled while a session switch request is in flight so the user cannot submit a turn against stale selection state.
- Consequences:
  - `New session` semantics become deterministic: the next turn must resolve a fresh provider/session for the empty scope.
  - existing cached providers for explicit session ids remain reusable, so historical-session switching and multi-session concurrency are preserved.
  - the empty-session-scope eviction is intentionally narrow and skips active turns.
- Alternatives considered:
  - optimistically trust that no stale empty-scope provider can exist.
  - close every cached provider on any session selection update.
  - add a dedicated `session/new` API endpoint instead of continuing to encode "new session" as `sessionId=""`.

## ADR-040: Cache Session Transcript Replay Snapshots in SQLite

- Status: Accepted
- Date: 2026-03-13
- Context:
  - the Web UI requests provider transcript replay whenever the user selects a historical ACP session in the right sidebar.
  - after ADR-039, every selection hit the provider again through ACP `session/load`, even if ngent had already replayed that exact provider session earlier.
  - repeated `session/load` calls add avoidable latency, reopen provider processes/runtimes, and make session browsing depend on live provider availability even for transcript content already observed locally.
  - the product still does not want to import provider-owned replay into hub `turns/events`, because that would blur the boundary between provider-owned history and ngent-created turns.
- Decision:
  - add SQLite table `session_transcript_cache(agent_id, cwd, session_id, messages_json, updated_at)`.
  - key cached transcript snapshots by provider session identity `(agent_id, cwd, session_id)` so the same session replay can be reused across threads and across server restarts.
  - make provider transcript replay read sqlite first; only call provider `LoadSessionTranscript` on cache miss.
  - write successful replay results, including empty transcript snapshots, back into sqlite after the provider call completes.
  - keep cache write failures non-fatal to the API response; the replay request itself should still succeed when the provider load succeeded.
- Consequences:
  - repeated historical session selection becomes local-first after the first successful replay.
  - session replay browsing survives server restart without requiring a new provider `session/load`.
  - provider-owned replay remains separate from durable hub `turns/events`; `/history` semantics do not change.
  - cache freshness is currently write-on-miss/write-on-success only; ngent does not yet compare provider `updatedAt` metadata before serving the cached snapshot.
- Alternatives considered:
  - keep always hitting provider `session/load` for every selection.
  - import replayed transcript into `turns/events` instead of caching a separate snapshot table.
  - key cache rows by `thread_id` only (rejected because the same provider session should be reusable across threads).
- Follow-up actions:
  - persist `session/list.updatedAt` metadata and refresh cached snapshots when the provider session advances.
  - evaluate whether a merged history view should combine cached provider replay with hub-local turns without duplicating context.

## ADR-039: Standardize Session Transcript Replay on ACP `session/load`

- Status: Accepted
- Date: 2026-03-12
- Context:
  - the Web UI already uses one generic provider transcript replay flow after session sidebar selection.
  - ACP `session/load` standard behavior replays prior conversation through `session/update` notifications before the RPC returns.
  - earlier ngent implementations reconstructed replay from provider-local files or databases, which diverged from ACP and created provider-specific behavior that the product no longer wants.
  - real-provider validation showed important runtime nuances:
    - Codex raw ACP `sessionId` values returned by `session/list` are scoped to the same embedded runtime and cannot be safely reused across a second runtime.
    - OpenCode replays transcript correctly over ACP `session/load`.
    - Qwen replays transcript correctly over ACP `session/load` for locally created sessions.
    - Qwen historical replay also includes non-message notifications such as `tool_call_update`, whose `content` payload does not follow the text-chunk schema.
    - Kimi CLI 1.20.0 resumes historical sessions over ACP `session/load`, but currently emits no replay transcript updates for those historical loads.
- Decision:
  - implement `SessionTranscriptLoader` for `codex`, `kimi`, `opencode`, and `qwen` by:
    - resolving the requested session through ACP `session/list`
    - calling ACP `session/load`
    - collecting replayed `user_message_chunk` and `agent_message_chunk` updates into transcript messages
  - add shared ACP update parsing and a shared replay collector for `session/load` transcript reconstruction.
  - keep the shared replay parser tolerant of non-message `session/update` variants so provider-specific tool/metadata updates do not abort transcript replay.
  - for Codex, resolve and load the session within the same embedded runtime so runtime-scoped raw session ids remain valid.
  - do not add provider-local transcript fallbacks behind the standard ACP path.
- Consequences:
  - provider transcript replay now follows ACP semantics instead of provider-local storage formats.
  - OpenCode, Codex, and Qwen replay their provider-owned transcript through standard ACP `session/load`.
  - providers may interleave replayable text chunks with tool or metadata updates; ngent now ignores those non-message updates during transcript reconstruction instead of treating them as transport errors.
  - Kimi currently remains limited by upstream behavior: historical `session/load` succeeds but yields no replay transcript messages on CLI 1.20.0.
  - Codex replay now reflects raw provider-owned prompt content, including wrapper text that was previously normalized away.
- Alternatives considered:
  - keep provider-local transcript parsing for Kimi/OpenCode/Codex/Qwen.
  - mix ACP `session/load` with local fallback parsing when replay updates are absent.
  - add a non-standard transcript import step into SQLite.
- Follow-up actions:
  - keep validating newer Kimi CLI releases for proper historical replay over standard ACP `session/load`.
  - decide later whether Codex replay text should be normalized again despite the ACP-first policy.

## ADR-038: Replay OpenCode Session Transcript from Local OpenCode SQLite Storage

- Status: Accepted
- Date: 2026-03-12
- Context:
  - the Web UI uses the same generic provider transcript replay flow after session sidebar selection for all ACP agents.
  - real OpenCode validation showed `session/list` and `session/load` both worked, but `opencode` exposed no `SessionTranscriptLoader`, so transcript replay always returned `supported=false`.
  - OpenCode stores replayable session content locally in `XDG_DATA_HOME/opencode/opencode.db`, with session metadata in `session` and visible text split across `message` + `part` rows.
- Decision:
  - implement `agents.SessionTranscriptLoader` for `internal/agents/opencode`.
  - read the local OpenCode SQLite database directly in read-only mode instead of shelling out to `opencode export`.
  - reconstruct transcript messages by:
    - selecting `user` / `assistant` rows from `message`.
    - appending only `part.type == "text"` payloads in insertion order.
    - dropping `reasoning`, `tool`, `step-start`, and `step-finish` parts.
  - reject transcript loads whose stored OpenCode `session.directory` does not match the active thread cwd.
- Consequences:
  - selecting an existing OpenCode session now replays provider-owned history in the center chat pane.
  - replay no longer depends on OpenCode CLI subcommands that may mutate state or fail due local DB write behavior.
  - the implementation now depends on OpenCode's local DB schema remaining compatible enough for read-only transcript reconstruction.
- Alternatives considered:
  - keep OpenCode session replay unsupported and rely only on ngent-local history.
  - invoke `opencode export <sessionID>` for transcript reconstruction.
  - parse only `session_diff` JSON files, even though they do not contain full chat transcript.
- Follow-up actions:
  - monitor future OpenCode schema changes and add a CLI-export fallback if the local DB layout becomes incompatible.

## ADR-037: Replay Kimi Session Transcript from Local Kimi Session Files

- Status: Accepted
- Date: 2026-03-12
- Context:
  - the Web UI already requests provider transcript replay after the user selects a provider-owned session from the right sidebar.
  - real Kimi debugging showed the backend successfully persisted and resumed `sessionId` through ACP `session/load`, but `kimi` exposed no `SessionTranscriptLoader`, so transcript replay always returned `supported=false`.
  - Kimi stores replayable chat history locally in `KIMI_HOME/sessions/*/<sessionId>/context.jsonl` for listable historical sessions, with assistant `think` blocks interleaved alongside visible text.
- Decision:
  - implement `agents.SessionTranscriptLoader` for `internal/agents/kimi`.
  - resolve transcript files from the local Kimi home directory and parse `context.jsonl` directly instead of trying to reconstruct transcript from hub-local turns.
  - keep only visible `user` / `assistant` text in replay payloads and drop `_checkpoint`, `_usage`, tool messages, and assistant `think` blocks before returning the transcript to the Web UI.
- Consequences:
  - selecting an existing Kimi session now replays provider-owned history in the center chat pane, matching the existing Codex session browsing behavior.
  - the replay remains on-demand only; ngent still does not import provider transcript into persisted SQLite turns/events.
  - fresh ACP-created Kimi sessions may still resume before they become visible in Kimi's own list/export surfaces; that remains an upstream/runtime limitation.
- Alternatives considered:
  - leave Kimi session replay unsupported and rely only on ngent-local history.
  - shell out to `kimi export` for every replay request instead of reading local session files.
  - attempt transcript reconstruction from `session/load` side effects during the next prompt.
- Follow-up actions:
  - monitor future Kimi CLI layout changes and add an export-based fallback if the local session directory format stops being stable enough for direct parsing.

## ADR-036: Persist Stable Codex Session IDs and Normalize Codex Transcript Replay

- Status: Accepted
- Date: 2026-03-11
- Context:
  - embedded Codex `session/new` can initially return only a provisional runtime-scoped id like `session-1`, while the durable session identity arrives later via `session/list` metadata (`_meta.threadId`).
  - persisting the provisional id caused fresh `New session` turns to collapse back onto the same thread session binding.
  - Codex transcript files also include wrapper-generated user messages for bootstrap context (`AGENTS.md`, `environment_context`) and prompt wrappers (`[Current User Input]`, IDE setup metadata), which polluted Web UI session replay.
- Decision:
  - treat the durable Codex session id as the canonical thread binding and defer fresh-session `session_bound` persistence/emission when the only known id still matches the provisional raw runtime id.
  - after the first prompt completes, retry `session/list` briefly to resolve the stable id, then update in-memory client state and persisted thread `agentOptions.sessionId` with that durable id.
  - normalize Codex session transcript replay on the backend before returning it to the Web UI:
    - drop bootstrap wrapper messages injected by the desktop environment.
    - extract the actual user prompt from known wrapper formats such as `[Current User Input]` and `## My request for Codex:`.
- Consequences:
  - fresh `New session` flows now bind to durable Codex session identities and no longer merge unrelated turns under one provisional raw id.
  - Web UI session replay shows user-visible prompts instead of provider/bootstrap scaffolding.
  - session list titles from provider metadata remain raw for now; only replayed message bodies are normalized.
- Alternatives considered:
  - persist the first raw runtime id immediately and rely on later correction.
  - leave transcript normalization to the frontend merge logic.
- Follow-up actions:
  - evaluate normalizing Codex `session/list` titles/previews in the backend so the session sidebar also hides wrapper-generated summary text.

## ADR-018: Embedded Web UI via Go embed

- Status: Accepted
- Date: 2026-02-28
- Context: users need a visual client to interact with the Ngent without writing curl commands or building a separate frontend project.
- Decision:
  - add a Vite + TypeScript (no framework) frontend under `internal/webui/web/src/`.
  - build output lands in `internal/webui/web/dist/`, embedded via `//go:embed web/dist` in `internal/webui/webui.go`.
  - register `GET /` and `GET /assets/*` in `httpapi` (lower priority than all `/v1/*` and `/healthz` routes).
  - SPA fallback: any non-API path returns `index.html`.
  - `make build-web` produces the dist; `web/dist` is not committed and must be rebuilt locally or in CI before packaging.
- Consequences: single-binary distribution with no external file dependencies; Go binary size increases by the size of the minified JS/CSS bundle (~200–400 KB estimated). Build pipeline requires Node.js for frontend changes and for any build that packages fresh web assets.
- Alternatives considered: separate static file directory (requires deployment of two artifacts); WebSocket-only SPA (rejected: SSE already implemented); React/Vue framework (rejected: adds runtime bundle weight and build complexity).
- Follow-up actions: add `npm run build` to CI pipeline; version-pin Node.js in project tooling docs.

## ADR-025: Hard-delete Thread Endpoint with Active-turn Lock

- Status: Accepted
- Date: 2026-03-03
- Context: users need to clean historical threads from both API and Web UI, while preserving the one-active-turn-per-thread guarantee and avoiding partial deletes.
- Decision:
  - add `DELETE /v1/threads/{threadId}` with ownership enforcement based on `X-Client-ID`.
  - return `409 CONFLICT` when the target thread currently has an active turn.
  - reserve a temporary turn-controller slot during deletion so no new turn can start on that thread while delete is in progress.
  - perform storage deletion in one transaction with explicit dependency order: `events` -> `turns` -> `threads`.
  - close and evict cached per-thread agent provider after successful delete.
- Consequences: deletion is deterministic and race-safe with active turn startup, but remains irreversible (no soft-delete/recover endpoint).
- Alternatives considered: soft-delete tombstone model, relying only on foreign-key cascades, and best-effort delete without turn-controller lock.
- Follow-up actions: add optional audit trail for delete operations if compliance requirements increase.

## ADR-028: Persist Thread Config Overrides and Surface Reasoning Control in Web UI

- Status: Accepted
- Date: 2026-03-06
- Context:
  - thread-scoped `config-options` support already enabled immediate model switching, but non-model settings like reasoning only lived inside the current ACP session.
  - for stdio providers (`opencode`, `qwen`, `gemini`), a reasoning change would otherwise disappear on the next turn because each turn starts a fresh ACP process.
  - Web UI only surfaced model in the composer footer, even though ACP can report model-specific reasoning choices in the same `configOptions` payload.
- Decision:
  - persist non-model current values returned by `POST /v1/threads/{threadId}/config-options` into `agentOptions.configOverrides`.
  - keep `agentOptions.modelId` as the durable model mirror, and store other current config values by config id in `configOverrides`.
  - update all built-in providers to reapply persisted non-model overrides on new ACP sessions:
    - embedded (`codex`, `claude`) reapply overrides after `session/new` on cached runtime initialization.
    - stdio (`opencode`, `qwen`, `gemini`, `kimi`) reapply overrides after `session/new` before `session/prompt`.
  - extend the Web UI composer footer to show both `Model` and `Reasoning` controls sourced from thread `configOptions`, with reasoning options refreshed after model changes.
  - cache config catalogs in the Web UI by agent id so same-agent threads reuse one shared model/reasoning option list without re-querying on every thread switch.
- Consequences:
  - reasoning-style settings now survive across future turns and restart/provider reinitialization boundaries.
  - reasoning remains model-specific because the UI always redraws from the latest ACP-reported thread config options after model changes.
  - switching between threads on the same agent is cheaper because the UI reuses the cached agent catalog and only applies thread-specific current values locally.
  - other non-model config categories are persisted server-side but are not yet surfaced as first-class UI controls beyond reasoning.
- Alternatives considered:
  - UI-only reasoning selector without persistence.
  - provider-specific persistence logic for reasoning rather than generic thread config override storage.
- Follow-up actions:
  - evaluate surfacing additional ACP config categories in the UI when product requirements justify more advanced controls.

## ADR-026: Thread-level Model Override Update API and Provider Reset

- Status: Accepted
- Date: 2026-03-05
- Context: users need to choose model at thread creation time and switch model on existing threads from Web UI/API without breaking one-active-turn-per-thread guarantees.
- Decision:
  - standardize thread-level model override as `agentOptions.modelId`.
  - add `PATCH /v1/threads/{threadId}` to update `agentOptions` for owned threads.
  - reject updates with `409 CONFLICT` while the thread has an active turn.
  - close cached per-thread provider after successful update so next turn re-initializes provider/session with new model config.
  - wire `modelId` into all provider factories; embedded `codex`/`claude` pass it to ACP `session/new` as `model`; `gemini` passes it to `session/new` and `session/prompt`; `opencode` passes `modelId` in `session/prompt`; `qwen` passes `model` in `session/prompt`.
- Consequences: model switching becomes explicit and deterministic at thread boundary; switching takes effect on next turn (not mid-turn) and may incur one provider re-init.
- Alternatives considered: in-session mutable model switching without provider reset, and provider-specific update endpoints.
- Follow-up actions: expose optional model catalogs per agent for richer dropdown UX and validate model ids against runtime-reported config options when available.

## ADR-027: ACP-backed Agent Model Catalog Endpoint and UI Dropdown Wiring

- Status: Accepted
- Date: 2026-03-05
- Context: Web UI hardcoded model lists drift from runtime reality; users need model options sourced directly from each agent's ACP runtime so create/switch flows stay accurate across codex/claude/opencode/gemini/kimi/qwen.
- Decision:
  - add `GET /v1/agents/{agentId}/models` and wire it to a backend `AgentModelsFactory`.
  - implement per-agent ACP discovery handshake (`initialize` + `session/new`) and parse model options from:
    - `session/new.configOptions` (`id=model`) for embedded codex/claude (acp-adapter latest).
    - `session/new.models.availableModels` for opencode/gemini/kimi/qwen.
  - normalize response shape to `[{id,name}]`, de-duplicate by `id`, and return `503 UPSTREAM_UNAVAILABLE` on discovery failure.
  - replace active-thread free-text model control with dropdown powered by the new endpoint:
    - active-thread header switches via dropdown + `PATCH /v1/threads/{threadId}`.
  - keep new-thread modal focused on agent/cwd/title creation and advanced JSON (no dedicated model selector).
- Consequences: model selection UX is runtime-accurate and provider-specific without frontend hardcoding; failure mode is explicit (upstream unavailable) and localized to model discovery.
- Alternatives considered: keep hardcoded frontend catalogs, or validate only at prompt-time without exposing a catalog endpoint.
- Follow-up actions: optionally add short-lived server-side model catalog cache and server-side create/update validation against discovered options.

## ADR Template

Use this template for new decisions.

```text
# ADR-XXX: <title>
- Status: Proposed | Accepted | Superseded
- Date: YYYY-MM-DD
- Context:
- Decision:
- Consequences:
- Alternatives considered:
- Follow-up actions:
```

## ADR-001: HTTP/JSON API with SSE Streaming

- Status: Accepted
- Date: 2026-02-28
- Context: turn output is incremental and long-running; clients need low-latency updates.
- Decision: use HTTP/JSON for request-response operations and SSE for server-to-client event streaming.
- Consequences: simpler client/network compatibility than WebSocket for one-way stream; requires reconnect/resume handling.
- Alternatives considered: WebSocket-only transport, polling.
- Follow-up actions: define event replay semantics and heartbeat policy.

## ADR-002: Client Identity via `X-Client-ID`

- Status: Superseded by ADR-064
- Date: 2026-02-28
- Context: server must isolate resources across multiple clients.
- Decision: require `X-Client-ID` on authenticated endpoints and scope data by that identity.
- Consequences: easy stateless routing and testing; header contract must be documented and validated strictly.
- Alternatives considered: query parameter id, session cookie.
- Follow-up actions: add optional auth token binding for production mode.

## ADR-003: SQLite Events as Source of Truth

- Status: Accepted
- Date: 2026-02-28
- Context: stream continuity and restart recovery require durable event history.
- Decision: persist turn events in append-only `events` table, indexed by thread/turn and sequence.
- Consequences: enables replay and audits; requires careful handling of SQLite contention.
- Alternatives considered: in-memory stream only, external queue.
- Follow-up actions: implement WAL mode, busy timeout, and compaction policy.

## ADR-004: Permission Workflow Fail-Closed

- Status: Accepted
- Date: 2026-02-28
- Context: runtime permissions are security-sensitive.
- Decision: when client decision is missing/invalid/late, default to deny.
- Consequences: safer default posture; may interrupt slow clients.
- Alternatives considered: fail-open with audit warning.
- Follow-up actions: add configurable timeout and clear UX hints.

## ADR-005: Localhost-by-Default Network Policy

- Status: Accepted
- Date: 2026-02-28
- Context: server may expose local filesystem and command capabilities.
- Decision: default bind `127.0.0.1:8686`; require explicit `--allow-public=true` for public interfaces.
- Consequences: secure local default; remote access requires intentional operator action.
- Alternatives considered: public by default.
- Follow-up actions: add warning log when public bind is enabled.

## ADR-006: M1 API Baseline for Health/Auth/Agents

- Status: Accepted
- Date: 2026-02-28
- Context: M1 requires a minimal but stable API contract before thread/turn APIs are implemented.
- Decision:
  - define `GET /healthz` response as `{ "ok": true }`.
  - define `GET /v1/agents` response key as `agents` with `id/name/status` fields.
  - gate only `/v1/*` endpoints behind optional bearer token (`--auth-token`).
  - standardize error envelope as `{ "error": { "code", "message", "details" } }`.
- Consequences: contract is simpler for clients and tests; request-id/hint fields are deferred to later milestones.
- Alternatives considered: keep earlier draft response schemas with extra fields.
- Follow-up actions: extend same error envelope to all future endpoints and document additional error codes as APIs expand.

## ADR-007: M3 Thread API Tenancy and Path Policy

- Status: Superseded by ADR-016
- Date: 2026-02-28
- Context: thread APIs introduce per-client resource ownership and filesystem-scoped execution context.
- Decision:
  - require `X-Client-ID` on all `/v1/*` endpoints and upsert client heartbeat on each request.
  - enforce `cwd` as absolute path under configured allowed roots.
  - return `404` for cross-client thread access to avoid existence leakage.
  - thread creation only persists metadata and does not start any agent process.
- Consequences: stronger tenancy boundaries and safer path policy at API edge; clients must always include identity header.
- Alternatives considered: permissive cross-client errors (`403`) and late validation in turn execution stage.
- Follow-up actions: wire same tenancy/path checks through turns and permission endpoints in M4+.

## ADR-008: M4 Turn Streaming over SSE with Persisted Event Log

- Status: Accepted
- Date: 2026-02-28
- Context: turns must stream incremental output while preserving durable history, cancellation, and per-thread single-active constraints.
- Decision:
  - use `POST /v1/threads/{threadId}/turns` as SSE response endpoint.
  - persist each emitted SSE event (`turn_started`, `message_delta`, `turn_completed`, `error`) into `events` table.
  - enforce one active turn per thread with in-memory controller; concurrent start on same thread returns `409 CONFLICT`.
  - implement `POST /v1/turns/{turnId}/cancel` to cancel active turn promptly.
  - expose `GET /v1/threads/{threadId}/history` with optional `includeEvents` query.
- Consequences: simple, testable streaming pipeline before provider integration; active-turn state is process-local and will require restart-recovery work in later milestones.
- Alternatives considered: separate stream endpoint per turn, websocket transport for M4.
- Follow-up actions: add restart-safe active-turn recovery and provider-backed execution in M5+.

## ADR-009: M5 ACP Stdio Provider and Permission Bridge

- Status: Accepted
- Date: 2026-02-28
- Context: M5 requires talking to ACP agents over stdio JSON-RPC and forwarding permission requests to HTTP clients with fail-closed semantics.
- Decision:
  - add `internal/agents/acp` provider that launches one external ACP agent process per streamed turn and handles newline-delimited JSON-RPC on stdio.
  - support inbound ACP message classes: response (pending id match), notification (`session/update`), request (`session/request_permission`; unknown methods return JSON-RPC method-not-found).
  - add `POST /v1/permissions/{permissionId}` and SSE event `permission_required`; bridge decisions back to ACP request responses.
  - timeout/disconnect default to `declined` (fail-closed); fake ACP flow converges with `stopReason="cancelled"`.
  - persist turn/event writes with `context.WithoutCancel(r.Context())` so terminal state is still durable after stream disconnect.
- Consequences: permission path is secure by default and testable without real codex dependency; pending permission lifecycle is process-local and late decisions can race with auto-close.
- Alternatives considered: fail-open timeout policy, websocket permission callbacks, delaying persistence until stream close.
- Follow-up actions: expose permission timeout metadata and add explicit permission-resolution terminal events.

## ADR-010: M6 Codex-ACP-Go Runtime Wiring

- Status: Superseded by ADR-013
- Date: 2026-02-28
- Context: M6 needs real codex provider enablement while keeping default tests stable in environments without codex binaries.
- Decision:
  - add runtime flags `--codex-acp-go-bin` and `--codex-acp-go-args` for codex ACP process configuration.
  - `GET /v1/agents` reports codex status as `unconfigured` when binary is absent and `available` when configured.
  - resolve turn provider lazily via a per-turn factory; codex turns create `internal/agents/acp` clients on demand.
  - use persisted `thread.cwd` as ACP process working directory for each turn.
  - keep default automated tests codex-independent; add env-gated optional smoke test (`E2E_CODEX=1` + `CODEX_ACP_GO_BIN`).
- Consequences: production path can run real codex providers without starting background processes at server boot; optional integration remains explicit and non-blocking for CI.
- Alternatives considered: eager provider startup on server boot, replacing existing tests with codex-only integration.
- Follow-up actions: add richer codex health diagnostics and startup validation in M8.

## ADR-011: M7 Context Window Injection and Compact Policy

- Status: Accepted
- Date: 2026-02-28
- Context: turns must preserve continuity across long threads and server restarts without relying on provider in-memory session state.
- Decision:
  - build per-turn injected prompt from `threads.summary` + recent non-internal turns + current input.
  - add runtime controls `context-recent-turns`, `context-max-chars`, and `compact-max-chars`.
  - enforce deterministic trimming order: drop oldest recent turns, then shrink summary, then shrink current input only as last resort.
  - add manual compact endpoint `POST /v1/threads/{threadId}/compact` that runs an internal summarization turn and persists updated `threads.summary`.
  - add `turns.is_internal` to mark compact/system turns and hide them from default history (`includeInternal=true` opt-in).
  - rebuild context solely from durable SQLite data after restart.
- Consequences: predictable context budget behavior and restart-safe continuity with auditable compact turns.
- Alternatives considered: provider-only session memory, in-memory context cache without durable summary, token-based approximate truncation only.
- Follow-up actions: add automatic compact trigger heuristics and token-aware budgeting in M8.

## ADR-012: M8 Reliability Alignment (TTL, Shutdown, Error Codes)

- Status: Accepted
- Date: 2026-02-28
- Context: final milestone requires explicit guarantees for concurrency conflicts, idle resource cleanup, shutdown behavior, and consistent API error semantics.
- Decision:
  - keep one-active-turn-per-thread invariant with `409 CONFLICT` and add additional concurrent multi-thread coverage.
  - add thread-agent cache with idle janitor (`agent-idle-ttl`) and JSON logs for idle reclaim/close actions.
  - add graceful shutdown flow: stop accepting requests, wait active turns, then force-cancel on timeout with structured logs.
  - unify API/SSE error code set to: `INVALID_ARGUMENT`, `UNAUTHORIZED`, `FORBIDDEN`, `NOT_FOUND`, `CONFLICT`, `TIMEOUT`, `INTERNAL`, `UPSTREAM_UNAVAILABLE`.
  - keep SSE disconnect fail-fast/fail-closed behavior to avoid hanging turns.
  - align acceptance checklist to executable `go test` plus `curl` verification commands.
- Consequences: operational behavior is predictable under contention, disconnects, and process lifecycle transitions.
- Alternatives considered: no idle janitor (manual cleanup only), immediate hard shutdown without grace period, preserving non-unified legacy error codes.
- Follow-up actions: optional enhancements after M8 include WebSocket transport, paginated history, RBAC, and audit expansion.

## ADR-013: Codex Provider Migration to Embedded Library

- Status: Accepted
- Date: 2026-02-28
- Context: sidecar mode required user-facing binary path configuration (`--codex-acp-go-bin`) and made deployment ergonomics/error modes depend on path wiring.
- Decision:
  - replace codex turn execution from external `codex-acp-go` process spawning to in-process `github.com/beyond5959/acp-adapter/pkg/codexacp` embedded runtime.
  - remove user-facing codex binary path flags; server now links acp-adapter library directly.
  - keep lazy startup and per-thread isolation by creating one embedded runtime per thread provider on first turn.
  - keep existing HTTP/SSE/permission/history contracts unchanged; permission round-trip remains fail-closed.
  - set `/v1/agents` codex status by embedded runtime preflight (`available`/`unavailable`) instead of path-config presence.
- Consequences: simpler operator UX and fewer path misconfiguration failures; server binary is now more tightly coupled to acp-adapter module/runtime behavior.
- Alternatives considered: keep sidecar-only mode; dual mode (embedded + sidecar fallback).
- Follow-up actions: define acp-adapter version pin/upgrade policy and add compatibility smoke checks across codex CLI/app-server versions.

## ADR-015: First-Turn Prompt Passthrough for Embedded Slash Commands

- Status: Accepted
- Date: 2026-02-28
- Context: context-window injection always wrapped prompts with `[Conversation Summary]` / `[Recent Turns]` / `[Current User Input]`, which masked first-turn slash commands (for example `/mcp call`) in embedded acp-adapter flows.
- Decision:
  - keep context wrapper for normal multi-turn continuity.
  - when `summary == ""` and there are no visible recent turns, pass through raw `currentInput` (still bounded by `context-max-chars`) instead of wrapping.
- Consequences:
  - first-turn slash commands remain functional in embedded mode, enabling deterministic permission round-trip validation (`approved` / `declined`).
  - first-turn request text persisted in history no longer includes synthetic wrapper headings.
- Alternatives considered:
  - parse wrapped `[Current User Input]` inside acp-adapter slash-command parser.
  - keep always-wrapped behavior and accept slash-command incompatibility.
- Follow-up actions:
  - evaluate an explicit API-level raw-input toggle if future providers need slash-command compatibility beyond first turn.

## ADR-016: Remove `--allowed-root` Runtime Parameter

- Status: Accepted
- Date: 2026-02-28
- Context: operators requested simpler startup without path allowlist configuration and required that `cwd` can be any user-specified absolute directory.
- Decision:
  - remove CLI flag `--allowed-root`.
  - server startup now configures allowed-roots internally as filesystem root.
  - keep `cwd` validation for absolute path only and retain tenancy/ownership rules.
- Consequences:
  - simpler startup and fewer configuration errors.
  - path-boundary restriction is effectively disabled in default runtime behavior.
- Alternatives considered:
  - keep `--allowed-root` and add a separate opt-out flag.
  - preserve strict allowlist-only behavior.
- Follow-up actions:
  - evaluate policy controls (for example opt-in restrictive mode) if deployments need stronger path boundaries.

## ADR-017: Human-Readable Startup Summary and Request Access Logs

- Status: Accepted
- Date: 2026-02-28
- Context: local operators found single-line JSON startup output hard to scan quickly; runtime troubleshooting also needed stable request completion telemetry.
- Decision:
  - print a concise multi-line startup summary to stderr with a QR code; print the service port and a concrete URL under the QR code.
  - keep structured request completion logs via `slog` for all HTTP traffic.
  - include `req_time`, `method`, `path`, `ip`, `status`, `duration_ms`, and `resp_bytes` in completion logs.
  - standardize log `time` and request `req_time` to UTC `time.DateTime` second precision for easier human scanning and stable parsing.
- Consequences:
  - local startup UX is easier to read without parsing JSON.
  - request observability is consistent across normal JSON responses and long-lived SSE requests.
- Alternatives considered:
  - keep startup as JSON only.
  - add ad-hoc per-endpoint logging instead of one centralized completion logger.
- Follow-up actions:
  - add optional request id correlation in completion logs and outbound SSE error events.

## ADR-019: OpenCode ACP stdio provider

- Status: Accepted
- Date: 2026-03-01
- Context: OpenCode supports ACP and is an actively developed coding agent; adding it as a provider gives users an alternative to the embedded Codex runtime.
- Decision: implement `internal/agents/opencode` as a standalone ACP stdio provider. One `opencode acp --cwd <dir>` process is spawned per turn. The package is self-contained with its own JSON-RPC 2.0 transport layer to avoid coupling with the internal `acp` package.
- Protocol differences from codex ACP that drove a separate implementation:
  - `protocolVersion` field is an integer (`1`) not a string.
  - `session/new` does not accept a client-supplied sessionId; the server assigns one and also returns a model list.
  - `session/prompt` uses a `prompt` array of content items instead of a flat `input` string.
  - `session/update` notifications carry delta text under `update.content.text` for `agent_message_chunk` events, not a flat `delta` field.
  - No `session/request_permission` requests from server to client (OpenCode handles tool permissions internally via MCP).
- Consequences:
  - `opencode` binary must be in PATH for the provider to be available; Preflight() is called at startup.
  - Model selection is optional via `agentOptions.modelId` in thread creation; defaults to OpenCode's configured default.
  - Turn cancel sends `session/cancel` and kills the process within 2s if it doesn't exit cleanly.

## ADR-020: Gemini CLI ACP stdio provider

- Status: Accepted
- Date: 2026-03-01
- Context: Gemini CLI (v0.31+) supports ACP via `--experimental-acp` flag; it uses the `@agentclientprotocol/sdk` npm package which speaks standard newline-delimited JSON-RPC 2.0 over stdio.
- Decision: implement `internal/agents/gemini` as a standalone ACP stdio provider. One `gemini --experimental-acp` process is spawned per turn. Protocol flow: `initialize` → `authenticate` → `session/new` → `session/prompt` with streaming `session/update` notifications.
- Key protocol details:
  - `PROTOCOL_VERSION = 1` (integer).
  - An explicit `authenticate({methodId: "gemini-api-key"})` call is required between `initialize` and `session/new` so Gemini reads `GEMINI_API_KEY` from the environment.
  - `GEMINI_CLI_HOME` is set to a fresh temp directory per turn, containing a minimal `settings.json` that selects API key auth; this prevents Gemini CLI from writing OAuth browser prompts to stdout, which would corrupt the JSON-RPC stream.
  - `session/update` notifications carry delta text under `update.content.text` for `agent_message_chunk` events (same structure as OpenCode).
  - Gemini can send `session/request_permission` requests; the provider bridges these through the hub server's `PermissionHandler` context mechanism. Approved maps to `{outcome: {outcome: "selected", optionId: "allow_once"}}`, declined to `reject_once`, cancelled to `{outcome: {outcome: "cancelled"}}`.
  - Turn cancel sends a `session/cancel` notification (no id, no response expected) and kills the process within 2s.
- Consequences:
  - `gemini` binary must be in PATH and `GEMINI_API_KEY` must be set for the provider to be available.
  - No model selection option at thread creation time; model is controlled by Gemini CLI's own configuration.

## ADR-022: Qwen Code ACP stdio provider integration

- Status: Accepted
- Date: 2026-03-03
- Context:
  - Qwen Code is available locally and supports ACP via `qwen --acp`.
  - hub requirements remain strict: one-active-turn-per-thread, fast cancel, fail-closed permissions, and no regressions for existing providers.
  - protocol inspection shows required ACP fields (`clientCapabilities.fs`, `mcpServers`) and provider-specific response variants.
- Decision:
  - implemented `internal/agents/qwen` as a standalone ACP stdio provider (one process per turn).
  - process command is fixed as `qwen --acp` (no user-supplied binary path in server config).
  - protocol flow is `initialize -> session/new -> session/prompt`, with required params:
    - `initialize.protocolVersion = 1`
    - `initialize.clientCapabilities.fs.readTextFile = false`
    - `initialize.clientCapabilities.fs.writeTextFile = false`
    - `session/new` includes `cwd` and `mcpServers: []`
    - `session/prompt` uses ACP prompt blocks (`[{type:"text", text:...}]`)
  - stream output is parsed from `session/update` when `update.sessionUpdate == "agent_message_chunk"` and delta comes from `update.content.text`.
  - handle `session/request_permission` by mapping hub decisions into ACP outcome format:
    - approve/decline: `outcome=selected` with matching `optionId`
    - cancel: `outcome=cancelled`
    - default deny on timeout/errors/no handler (fail-closed)
  - cancellation path sends `session/cancel` with `sessionId` on context cancellation and converges to `stopReason=cancelled`.
  - stderr is drained/discarded to avoid protocol stream corruption; existing providers (`codex`, `opencode`, `gemini`) behavior remains unchanged.
- Consequences:
  - qwen availability is startup-preflight dependent (`qwen` in PATH).
  - real qwen turns depend on local runtime prerequisites (writable qwen home/config + auth/network readiness), so environment misconfiguration can fail before ACP turn execution.
  - provider must tolerate schema drift across qwen versions (new optional fields in `session/new` response).
  - test surface expanded:
    - fake-process ACP tests for initialize/session/new/session/prompt/session/update
    - permission mapping tests (`approved`, `declined`, `cancelled`)
    - optional real smoke test (`E2E_QWEN=1`)
- Alternatives considered:
  - reuse a single generic ACP provider configured by command/args at runtime.
  - force qwen through existing `opencode`/`gemini` adapters.
  - postpone qwen until ACP schema is frozen upstream.
- Follow-up actions:
  - improve qwen preflight diagnostics for filesystem/auth prerequisites (beyond PATH existence).
  - keep validating against newer qwen releases for ACP schema compatibility.

## ADR-023: Shared ACP stdio transport for OpenCode and Qwen providers

- Status: Accepted
- Date: 2026-03-03
- Context:
  - `internal/agents/opencode` and `internal/agents/qwen` had duplicated JSON-RPC stdio transport code (request id/pending map, read loop, inbound request handling, notify/call framing, process termination helpers).
  - duplicated protocol plumbing increased maintenance risk and made bug fixes easy to diverge across providers.
- Decision:
  - extracted shared package `internal/agents/acpstdio` with:
    - newline-delimited JSON-RPC connection (`Conn`) supporting `Call`, `Notify`, notifications, and inbound request handling.
    - shared JSON-RPC message/error types.
    - shared helpers: `ParseSessionID`, `ParseStopReason`, `TerminateProcess`.
  - refactored both providers to use the shared transport while keeping provider-specific ACP behavior unchanged:
    - OpenCode flow and modelId handling unchanged.
    - Qwen permission mapping and fail-closed behavior unchanged.
- Consequences:
  - transport-layer fixes are now centralized and consistent across providers.
  - provider files are shorter and focused on protocol semantics instead of wire plumbing.
  - regression risk from this refactor is controlled by fake-process tests + full `go test ./...` + real E2E smoke tests for both providers.
- Alternatives considered:
  - keep duplicated transport code and only copy fixes manually.
  - extract only tiny helper funcs without shared connection type.
- Follow-up actions:
  - if Gemini migration value is clear, evaluate moving Gemini transport to `acpstdio` in a separate change (to keep current refactor blast radius limited).

## ADR-024: Claude Code embedded provider via claudeacp runtime

- Status: Accepted
- Date: 2026-03-03
- Context:
  - Claude Code is the primary Anthropic coding agent; it was listed as a planned provider (`🔜`) since project inception.
  - `github.com/beyond5959/acp-adapter` already contained a complete parallel `pkg/claudeacp` package with identical API surface to `pkg/codexacp`; no new library dependency was needed.
  - Preflight for Claude does not require a binary path check — availability is determined entirely by the presence of `ANTHROPIC_AUTH_TOKEN` in the environment.
- Decision:
  - implement `internal/agents/claude` as an embedded provider package mirroring `internal/agents/codex`.
  - replace `codexacp` references with `claudeacp`; `Preflight()` checks `ANTHROPIC_AUTH_TOKEN != ""` (no binary lookup).
  - `DefaultRuntimeConfig()` delegates to `claudeacp.DefaultRuntimeConfig()`, which reads `ANTHROPIC_AUTH_TOKEN` and `ANTHROPIC_BASE_URL` from environment.
  - wire into server startup: preflight, `/v1/agents` status, `AllowedAgentIDs`, and `TurnAgentFactory`.
- Consequences:
  - claude availability is purely environment-variable dependent; no binary installation required beyond valid API credentials.
  - `ANTHROPIC_BASE_URL` allows pointing at a compatible proxy or local endpoint (e.g., for testing or corporate gateways).
- Alternatives considered:
  - implement as an ACP stdio provider wrapping the `claude` CLI binary (rejected: CLI spawns its own runtime per invocation with higher latency and no direct permission bridge).
  - share implementation with codex via generics/interface (rejected: would couple two independently-versioned runtimes).
- Follow-up actions:
  - add permission round-trip E2E test for Claude (approved/declined/cancelled paths).

## ADR-079: Thread-level model switching via ACP session config options

- Status: Accepted
- Date: 2026-03-05
- Context:
  - model switching previously used thread metadata patch (`PATCH /v1/threads/{threadId}` with `agentOptions.modelId`) and recreated provider state, while ACP now standardizes runtime config through `session/new.configOptions` + `session/set_config_option`.
  - Web UI requirement is immediate switch on model select (no extra apply action), and model list/selected value must come from ACP session config data (including per-model descriptions).
- Decision:
  - add thread-scoped config option endpoints:
    - `GET /v1/threads/{threadId}/config-options`
    - `POST /v1/threads/{threadId}/config-options` (`configId`, `value`)
  - `POST` applies changes through provider `SetConfigOption`, backed by ACP `session/set_config_option`.
  - support `ConfigOptionManager` on all built-in providers:
    - embedded (`codex`, `claude`): mutate cached session in-place.
    - stdio (`opencode`, `qwen`, `gemini`, `kimi`): perform ACP handshake/apply flow and persist resulting model id for subsequent turns.
  - keep `agentOptions.modelId` as durable thread metadata mirror when `configId=model` succeeds.
  - Web UI model selector is bound to thread-level `configOptions` model option and applies immediately on `change`.
- Consequences:
  - model UX is consistent with ACP protocol semantics and no longer depends on separate apply state.
  - per-option descriptions are available to UI and rendered beneath selector.
  - active-turn safety remains strict (`409` on config mutation while turn is active).
- Alternatives considered:
  - continue using thread patch only and skip ACP `session/set_config_option`.
  - expose only agent-level model catalog and infer current model from local metadata.
- Follow-up actions:
  - optionally add richer Web UI rendering for non-model config categories (e.g. reasoning level) using the same API.

## ADR-080: Persist agent config catalogs in SQLite and refresh them asynchronously on startup

- Status: Accepted
- Date: 2026-03-06
- Context:
  - thread config UX now depends on ACP `configOptions`, especially model-specific reasoning choices.
  - querying live providers on every thread switch is unnecessary and becomes user-visible after server restart because model/reasoning catalogs are metadata, not per-turn output.
  - different threads on the same agent can keep different selected model/reasoning values, while still reusing the same underlying catalog data.
- Decision:
  - add sqlite table `agent_config_catalogs(agent_id, model_id, config_options_json, updated_at)` and persist normalized ACP config-options snapshots there.
  - keep per-thread selected values in `threads.agent_options_json`:
    - `modelId`
    - `configOverrides`
  - read path:
    - `GET /v1/threads/{threadId}/config-options` first loads the stored catalog row for the thread's selected model (or a reserved default snapshot when no model is selected yet), then overlays the thread's own selected current values.
    - `GET /v1/agents/{agentId}/models` derives model list from stored catalogs before falling back to live discovery.
  - write path:
    - `POST /v1/threads/{threadId}/config-options` still mutates the live provider/session, then persists both thread selection state and the current model's returned config-options snapshot.
  - startup behavior:
    - launch a background refresher goroutine after server initialization.
    - refresher queries default + per-model config catalogs for built-in agents and updates sqlite without delaying HTTP startup or blocking frontend requests.
    - on partial refresh failure, keep previously stored rows for models that could not be refreshed instead of deleting them.
- Consequences:
  - restart no longer forces frontend-visible catalog discovery in the common case; stored model/reasoning metadata remains immediately available.
  - reasoning lists remain model-specific while thread-selected current values stay isolated per thread.
  - partial refresh strategy trades perfect freshness for availability and catalog continuity.
- Alternatives considered:
  - keep catalogs only in frontend memory and re-query on restart.
  - store only an agent-level flat reasoning list (rejected because reasoning is model-dependent).
  - block startup until all agents refresh live catalogs (rejected because it would directly impact frontend responsiveness).

## ADR-029: Consolidate sidebar thread actions into a drawer and reuse thread patch for rename

- Status: Accepted
- Date: 2026-03-06
- Context:
  - sidebar thread rows already show dense metadata and activity state; a direct delete affordance makes the list visually noisy and increases the chance of accidental destructive clicks.
  - thread rename is metadata-only and should share the existing thread ownership/conflict semantics instead of adding a dedicated endpoint.
- Decision:
  - replace the direct sidebar delete button with a per-thread drawer trigger.
  - render drawer actions as text-only controls in this order:
    - `Rename`
    - `Delete` (danger styling)
  - extend `PATCH /v1/threads/{threadId}` to accept optional `title` in addition to `agentOptions`.
  - keep rename under the same active-turn conflict guard as other thread mutations.
- Consequences:
  - sidebar actions are less visually noisy and destructive actions are one click deeper.
  - rename and thread metadata/config mutations now share one API contract.
  - users must wait for a running turn to finish, or cancel it, before renaming the thread.

## ADR-030: Handle codex `item/tool/requestUserInput` and `item/tool/call` without hard RPC failure

- Status: Accepted
- Date: 2026-03-06
- Context:
  - codex app-server may emit server requests `item/tool/requestUserInput` and `item/tool/call` during MCP-related operations.
  - previous bridge behavior returned `-32000 ... is not supported`, which caused app-server-side hard errors and interrupted command flow.
- Decision:
  - implement a compatibility fallback in embedded adapter path:
    - for `item/tool/requestUserInput`: return schema-compatible `answers` payload by auto-selecting the first option label for each question.
    - for `item/tool/call`: return schema-compatible `DynamicToolCallResponse` with `success=false` and one text content item, instead of throwing RPC method error.
- Consequences:
  - removes immediate `-32000` hard-fail class for these methods and allows app-server to continue handling tool flow.
  - behavior is still a fallback (not full interactive user-input / dynamic-tool execution support).
- Alternatives considered:
  - keep fail-closed `-32000` hard error (rejected: user-facing breakage for MCP flows).
  - fully implement interactive user-input/dynamic tool execution end-to-end in hub UI + protocol bridge in one step (rejected for hotfix scope).

## ADR-031: Kimi CLI ACP stdio provider with dual startup syntax fallback

- Status: Accepted
- Date: 2026-03-09
- Context:
  - Kimi CLI now exposes ACP mode in upstream docs, but current official pages show both `kimi acp` and `kimi --acp` as startup forms.
  - hub requirements remain unchanged: one active turn per thread, fast cancel, fail-closed permissions, and no user-supplied binary-path flags.
- Decision:
  - implement `internal/agents/kimi` as a standalone ACP stdio provider (one process per turn).
  - startup tries `kimi acp` first, then retries with `kimi --acp` if the first form fails before ACP initialize completes.
  - keep protocol flow aligned with existing stdio providers:
    - `initialize -> session/new -> session/prompt`
    - stream deltas from `session/update` / `agent_message_chunk`
    - handle `session/request_permission` with fail-closed mapping
    - send `session/cancel` on context cancellation
  - wire `kimi` into startup preflight, `/v1/agents`, thread allowlist, turn factory, model discovery, and startup catalog refresh.
- Consequences:
  - Kimi becomes a built-in provider without adding new runtime flags or config surface.
  - startup is resilient to the current upstream command-form drift across Kimi CLI docs/releases.
  - real Kimi turns still depend on local Kimi authentication and network readiness.
- Alternatives considered:
  - hardcode only `kimi acp` (rejected: conflicts with current upstream `--acp` docs).
  - hardcode only `kimi --acp` (rejected: conflicts with IDE integration docs and user expectation).
  - add a user-facing command override flag (rejected: unnecessary config surface and contrary to current built-in provider policy).

## ADR-032: Shared common agent config/state helper without protocol unification

- Status: Accepted
- Date: 2026-03-09
- Context:
  - built-in agents duplicated the same common fields and methods repeatedly:
    - `Dir`
    - `ModelID`
    - `ConfigOverrides`
    - model/config override normalization, cloning, and state updates after `SetConfigOption`
  - the duplicated code existed across both stdio providers (`gemini`, `opencode`, `qwen`, `kimi`) and embedded providers (`codex`, `claude`).
  - previous experience in this repo shows that protocol flows diverge materially between providers, so a single generic provider abstraction would be riskier than the duplication it removes.
- Decision:
  - extract shared common config/state into `internal/agents/agentutil`:
    - `agentutil.Config`
    - `agentutil.State`
  - move only common constructor validation and mutable model/config-override state management into the shared helper.
  - keep each provider's transport/runtime/session logic independent.
- Consequences:
  - repeated per-provider bookkeeping is reduced without coupling unrelated ACP/runtime flows.
  - future providers can reuse the same common state helper while still implementing provider-specific protocol behavior.
  - embedded and stdio providers remain free to evolve independently where their protocol/runtime requirements differ.
- Alternatives considered:
  - create one generic ACP provider with pluggable commands and hooks (rejected: protocol differences are too large and already visible across current providers).
  - leave duplicated `Config`/`Client` state in place (rejected: ongoing maintenance cost and drift risk).

## ADR-033: Surface ACP plan updates as first-class SSE and Web UI state

- Status: Accepted
- Date: 2026-03-09
- Context:
  - ACP agents can emit `session/update` notifications with `sessionUpdate == "plan"` and a full `entries[]` list describing the current execution plan.
  - the hub previously only mapped `agent_message_chunk` into `message_delta`, so users could not see plan progress in the Web UI and history reloads lost that context entirely.
- Decision:
  - normalize ACP `session/update` payloads in shared agent code, recognizing both `agent_message_chunk` and `plan`.
  - route plan replacements through a new per-turn `PlanHandler` context callback, parallel to the existing permission callback pattern.
  - emit and persist a dedicated SSE/history event `plan_update` with payload `{"turnId","entries":[]}`.
  - render `plan_update` in the Web UI as a live plan card above the streaming agent bubble, and rebuild the final plan state from persisted turn events when loading history.
- Consequences:
  - ACP plan state is now visible during live execution without overloading `message_delta`.
  - history replay preserves the last known plan instead of dropping it on refresh.
  - empty `entries[]` remains meaningful as "clear the current plan", so the hub must preserve replacement semantics instead of merging incrementally.
- Alternatives considered:
  - fold plan text into `message_delta` (rejected: mixes distinct ACP concepts and loses replacement semantics).
  - keep plan rendering purely transient in the browser (rejected: history reload would still discard plan state).


## ADR-035: Add opt-in ACP debug tracing behind `--debug`

- Status: Accepted
- Date: 2026-03-11
- Context:
  - operators need a direct way to inspect the exact ACP methods and payloads exchanged with built-in agents when debugging protocol/runtime issues.
  - normal production logs must remain concise, protocol-only, and continue redacting sensitive information.
  - ACP traffic currently flows through three different boundaries:
    - stdio JSON-RPC transport (`acpstdio`)
    - legacy stdio ACP client (`internal/agents/acp`)
    - embedded runtimes (`codex`, `claude`)
- Decision:
  - add a startup flag `--debug` that raises the shared `slog` logger to debug level.
  - when enabled, emit structured stderr log entries `acp.message` for inbound/outbound ACP JSON-RPC messages across all supported transport paths.
  - include `component`, `direction`, `rpcType`, `method`, `id`, and the sanitized `rpc` payload in each trace log.
  - keep a redaction pass in front of debug logging so sensitive keys and common token formats are masked before serialization.
- Consequences:
  - ACP handshakes, prompts, updates, permission requests, and permission responses are inspectable without changing HTTP/stdout protocol behavior.
  - debug mode can produce high log volume and should be used only when investigating runtime issues.
  - the tracing mechanism remains transport-local and does not require changing provider/public API contracts.
- Alternatives considered:
  - always log ACP payloads at info level (rejected: too noisy and unsafe for normal operation).
  - add per-provider bespoke debug flags (rejected: fragmented UX and duplicated plumbing).
  - expose raw unredacted ACP dumps (rejected: conflicts with repository logging/redaction requirements).

## ADR-081: Persist thread-level ACP session selection and resume through provider sessions

- Status: Accepted
- Date: 2026-03-11
- Context:
  - users need to browse an agent's historical ACP sessions and continue a conversation from an existing provider-owned session.
  - hub threads already persist local turn history, but blindly injecting that history into prompts duplicates context once ACP `session/load` has restored the provider's own transcript.
  - the frontend needs paginated session discovery and a lightweight way to switch between "new session" and "existing session" without changing the SQLite schema.
- Decision:
  - persist the selected ACP session id in `threads.agent_options_json` as `sessionId`.
  - expose `GET /v1/threads/{threadId}/sessions` backed by ACP `session/list`, using a fresh provider instance so sidebar discovery does not disturb cached turn runtimes.
  - extend built-in providers to:
    - load `sessionId` through ACP `session/load` when present.
    - create a fresh session through `session/new` otherwise.
    - report the effective session id back during turn setup so the server can persist it and emit SSE `session_bound`.
  - once `sessionId` is present on a thread, skip local recent-turn prompt injection and rely on ACP session state for continuation.
  - keep Web UI session selection as a thread metadata mutation (`PATCH /v1/threads/{threadId}`) and model the "New session" action as clearing `sessionId`.
- Consequences:
  - session continuation survives provider restart/server restart when the agent supports ACP `session/load`.
  - right-sidebar session browsing stays paginated (`nextCursor`/`Show more`) and does not require schema changes.
  - local SQLite history remains a hub-local view and is no longer the source of truth for resumed ACP context on bound threads.
  - historical ACP transcript import is deferred and tracked separately as a known limitation.
- Alternatives considered:
  - import the full ACP transcript from `session/load` into SQLite immediately (rejected: larger behavioral change, requires reliable transcript reconstruction).
  - keep relying on hub prompt injection even after binding to an ACP session (rejected: duplicates already-restored conversation context).
  - add a dedicated sessions table instead of reusing `agentOptions` JSON (rejected: unnecessary schema churn for a single thread-scoped selection value).

## ADR-082: Scope Web UI chat playback to the selected ACP session

- Status: Accepted
- Date: 2026-03-11
- Context:
  - the Web UI session sidebar switches `agentOptions.sessionId` on the active thread without changing `threadId`.
  - local history remains stored per thread, and each turn's effective ACP session is persisted through the `session_bound` event stream.
  - refreshing the chat area only on thread changes leaves stale messages visible after choosing a different session from the sidebar.
- Decision:
  - treat `(threadId, sessionId)` as the client-side chat render scope.
  - when the active thread's selected session changes outside an active turn, rebuild the chat area and reload history for that scope.
  - when the selected scope changes to one that is still streaming in the background, rebuild the chat area unless that same scope's streaming bubble is already mounted in the current DOM.
  - filter locally persisted turns by their `session_bound` event so the center chat panel renders only turns recorded for the selected session; an empty `sessionId` shows only unbound turns.
- Consequences:
  - clicking a session in the sidebar replays that session's ngent-recorded turns instead of keeping the previously rendered session on screen.
  - revisiting a background-streaming session restores its live typing/loading state from scope-local client buffers instead of dropping back to the persisted history snapshot.
  - session changes reported during a live turn do not wipe the streaming bubble; the full refresh is deferred until the turn finishes.
  - transcript content that predates ngent participation is still not imported from the provider and remains covered by KI-021.
- Alternatives considered:
  - add a session-scoped history endpoint immediately (rejected: larger server contract change while turn events already contain the session discriminator).
  - keep all thread turns visible regardless of selected session (rejected: does not meet the expected session playback behavior).

## ADR-083: Allow concurrent turns across different sessions on the same thread

- Status: Accepted
- Date: 2026-03-13
- Context:
  - once ACP session switching shipped, users needed to leave one session streaming while switching the same thread to another session and continuing work there.
  - the existing runtime/controller and provider cache were keyed only by `threadId`, so a running turn blocked `PATCH /v1/threads/{threadId}` session changes and forced all sessions on that thread into a single execution lane.
  - provider instances keep mutable session/model/config state, so simply removing the conflict check would have let different sessions reuse the wrong cached provider.
- Decision:
  - change turn concurrency from thread-wide to `(threadId, sessionId)` scope, with empty `sessionId` representing the provisional "new session" scope until `session_bound` arrives.
  - keep delete/compact and shared thread-state mutations guarded at whole-thread scope.
  - rebind active turn scope when `session_bound` reports the effective session id so same-session re-entry remains blocked after a new session is created.
  - key cached providers by thread/session/shared-option scope instead of raw thread id, and evict all cached providers for a thread after shared config changes.
  - update the Web UI to store messages and live stream state per chat scope `(threadId, sessionId)` so background turns in one session do not overwrite another session's visible transcript.
- Consequences:
  - users can switch the right-side session picker and start work in session B while session A is still streaming on the same thread.
  - same-session re-entry protection, cancellation, and permission handling remain unchanged.
  - thread-wide config changes remain serialized and are documented as KI-021.
- Alternatives considered:
  - keep thread-wide turn serialization and force users to create separate threads per ACP session (rejected: poor UX and redundant thread duplication).
  - remove the conflict check without changing provider cache scope (rejected: would mix session-bound provider state and route turns to the wrong ACP session).

## ADR-084: Persist ACP slash commands as agent-level SQLite snapshots

- Status: Accepted
- Date: 2026-03-13
- Context:
  - ACP agents can emit `available_commands_update` during `session/update`, and the Web UI needs a durable source for slash-command suggestions instead of depending on the current in-memory stream only.
  - the same agent can be opened from multiple threads, and slash-command suggestions should survive server restart once any thread has observed them.
  - the requested UI interaction is lightweight composer assistance, not a separate command palette service.
- Decision:
  - extend shared ACP update parsing with `available_commands_update` and normalize the payload into a common `SlashCommand` model.
  - add a shared per-turn `SlashCommandsHandler` callback and wire all built-in ACP providers to forward the latest slash-command snapshot through that callback.
  - persist the snapshot in SQLite table `agent_slash_commands` keyed by `agent_id`, replacing the previous row each time a new update arrives.
  - expose `GET /v1/threads/{threadId}/slash-commands` so the Web UI can read the cached commands for the active thread's agent while still enforcing normal thread ownership checks.
  - in the Web UI composer, only open the slash-command picker when the current input starts with `/`, so `abc /plan` remains ordinary message text.
  - fetch slash commands lazily when the user types `/`; the first bare `/` in each slash interaction forces a backend refresh for the active thread even if the browser already has an in-memory agent cache, and if the refreshed snapshot is empty, keep the composer as plain text and do not render an empty or loading picker.
  - provider adapters must start observing `session/update` before `session/new` / `session/load` completes if they want to capture capability snapshots, because real agents such as Kimi can emit `available_commands_update` before the first `session/prompt`.
- Consequences:
  - once any thread observes an agent's slash commands, later threads for that agent can reuse them immediately and server restart does not clear the cache.
  - slash-command updates do not become part of persisted turn history; they are stored as agent capability snapshots instead of user-visible transcript events.
  - agents that never emit `available_commands_update` still behave normally in the composer because `/` falls back to plain input instead of trapping the UI in a retry loop.
  - the first `/` after entering slash mode always re-checks sqlite through the thread endpoint, so the browser no longer silently masks stale or newly refreshed slash-command snapshots behind a hot in-memory cache.
  - adapters that also support transcript replay must still suppress pre-prompt message chunks, otherwise `session/load` history replays could leak into the visible answer stream.
  - embedded adapters whose runtime does not replay historical notifications must install a slash-command monitor before `session/new` / `session/load`, cache the initial snapshot on the provider instance, and replay that cached snapshot into the active turn before `session/prompt`; codex now follows this rule so config-option queries and first turns observe the same slash-command state.
  - stdio adapters must also install `session/update` handlers immediately after `initialize`, because `acpstdio.Conn` drops notifications that arrive before `SetNotificationHandler`; Qwen and OpenCode now follow the same pre-prompt capability-capture rule as Kimi.
  - Kimi, Qwen, OpenCode, and Gemini now share one internal ACP notification handler builder so the providers cannot silently drift on pre-prompt slash-command, message-chunk, or plan handling semantics even when they use different connection implementations.
  - when a provider can expose its current slash-command snapshot outside a turn, ngent may backfill a missing sqlite row from that live provider state during thread initialization flows; codex now does this on the `config-options` path so a fresh thread can show slash commands before the first prompt.
  - `GET /v1/threads/{threadId}/slash-commands` must also perform the same best-effort backfill on sqlite miss, because users can type `/` before any parallel thread-initialization request finishes; Qwen now relies on this path for deterministic fresh-thread behavior.
  - direct ACP stdio providers must apply the same slash-command cache logic to both their turn path and their config-session path; Kimi, Qwen, OpenCode, and Gemini now all use one provider-local `SlashCommandsCache` so fresh-thread slash probes and later turns observe the same snapshot source.
  - if a future provider varies slash commands by workspace, model, or session, the current `agent_id` cache key may be too coarse and will need refinement.
- Alternatives considered:
  - keep slash commands in memory only (rejected: loses state on restart and leaves fresh threads without suggestions until another turn streams).
  - persist slash commands per thread (rejected: duplicates identical agent data and prevents reuse across threads).
  - append slash-command updates into `turns/events` only (rejected: complicates retrieval for the composer and mixes capability cache with transcript history).

## ADR-085: Normalize rich ACP permission requests before bridging them into ngent

- Status: Accepted
- Date: 2026-03-14
- Context:
  - real ACP providers do not guarantee that `session/request_permission.toolCall` is a flat string map.
  - Kimi CLI 1.22.0 sends structured previews containing diff/content arrays, `toolCallId`, and a human-readable `title`.
  - ngent's direct ACP adapters for Kimi and Qwen previously decoded that payload as `map[string]string`, which caused JSON decode failure and an immediate fail-closed reject before HTTP/SSE could emit `permission_required`.
- Decision:
  - add one shared ACP permission-request parser that accepts structured `toolCall` payloads, preserves key metadata (`sessionId`, `toolCallId`, first `path`), and derives a normalized `PermissionRequest` for the rest of ngent.
  - classify permission badges into the existing Web UI families (`file`, `command`, `network`, `mcp`) from the tool title/content instead of trusting provider-specific `kind` strings such as `execute`.
  - use the tool title as the primary user-facing command label (for example `WriteFile: soul.md`) so the UI can display the same preview text the provider asked the user to approve.
  - reuse this shared normalization in both Kimi and Qwen so direct ACP providers do not drift on permission bridging semantics.
- Consequences:
  - real Kimi file-write requests now surface through the normal ngent permission workflow instead of being auto-rejected invisibly.
  - fail-closed behavior is preserved for malformed payloads or missing handlers, but only after attempting structured decode.
  - future ACP providers can attach richer tool-call previews without forcing another provider-specific permission decoder.
- Alternatives considered:
  - keep provider-local bespoke permission decoding (rejected: already diverged across adapters and failed on a real payload shape).
  - flatten structured tool-call previews into strings at the transport layer (rejected: hides provider metadata and makes badge classification/path extraction harder later).

## ADR-048: Web UI fresh-session scopes for repeated `New session`

- Status: Accepted
- Date: 2026-03-16
- Context:
  - the Web UI previously keyed every unbound thread view to the same empty-session scope `${threadId}::`.
  - when a user started a fresh session, cancelled the turn before ACP emitted `session_bound`, and clicked `New session` again, the UI stayed on that same anonymous scope and kept showing the cancelled placeholder/content.
  - fast-cancelled turns with no `session_bound` and no visible response text are transient UI artifacts, not durable conversation state the user expects to keep re-entering.
- Decision:
  - treat explicit `New session` in the Web UI as a client-side fresh-session scope with a temporary key `${threadId}::@fresh:<uuid>` until a real ACP session binds.
  - allow `New session` even when the active thread already has no persisted `sessionId`; in that case, rotate to a new fresh-session scope locally instead of treating the action as a no-op.
  - seed that fresh-session scope with an empty message cache and skip server history replay for it until the first turn binds to a real ACP session id.
  - filter cancelled turns with no `session_bound` and no visible response text out of empty-session replay reconstruction so page reload does not resurrect those transient placeholders.
- Consequences:
  - `send -> cancel -> New session` now returns to a blank composer even when the cancelled turn never acquired a session id.
  - late completion/cancel callbacks for the abandoned scope still land in their original scope and no longer leak into the newly opened fresh session.
  - reopening the thread after reload no longer repaints empty cancelled placeholders from those abandoned pre-bind attempts.
- Alternatives considered:
  - keep using the single empty-session scope `${threadId}::` for all fresh-session attempts (rejected: this is the bug).
  - persist a backend-generated fresh-session nonce in thread metadata (rejected for now: more invasive than needed for the Web UI reset bug).

## ADR-086: Preserve ACP tool-call updates as first-class turn events

- Status: Accepted
- Date: 2026-03-16
- Context:
  - ACP tool execution progress is reported through `session/update` variants `tool_call` and `tool_call_update`, not through normal assistant text deltas.
  - those payloads can carry structured fields such as `toolCallId`, `status`, `content[]`, `locations[]`, `rawInput`, and `rawOutput`.
  - ngent previously tolerated non-text tool-call payloads during transcript replay, but live turn streaming still discarded them because only message/plan/reasoning updates were bridged into HTTP/SSE and the Web UI.
- Decision:
  - extend shared ACP update parsing to normalize `tool_call` and `tool_call_update` into one structured ngent event model without flattening the payload into plain text.
  - preserve the raw structured JSON for `content`, `locations`, `rawInput`, and `rawOutput` so downstream clients can evolve their rendering without changing the transport contract again.
  - add a per-turn tool-call callback in the agent layer, persist those updates into SQLite turn events, and emit them over SSE using the same event names (`tool_call`, `tool_call_update`).
  - have the Web UI merge tool-call events by `toolCallId`, so live streaming and history reload reconstruct the same final tool state.
- Consequences:
  - clients can now observe structured tool activity during a turn and after reload/history fetch.
  - ngent keeps ACP semantics intact instead of inventing a hub-specific flattened tool transcript format.
  - the current Web UI renders common text/diff/command/path payloads directly and falls back to generic JSON blocks for unsupported tool-call content shapes.
- Alternatives considered:
  - continue ignoring tool-call updates outside transcript replay (rejected: loses important execution state in the UI).
  - flatten tool-call payloads into `message_delta` text (rejected: destroys structure and makes incremental updates ambiguous).
  - keep tool-call state only in browser memory (rejected: reload/history would still lose it).

## ADR-087: Render assistant turns as ordered UI segments instead of one aggregated bubble

- Status: Accepted
- Date: 2026-03-19
- Context:
  - ngent already persisted `message_delta`, `reasoning_delta`, `tool_call`, and `tool_call_update` turn events in order.
  - the Web UI still collapsed those into one final assistant bubble plus aggregated reasoning/tool sections, which hid the actual execution sequence and felt noticeably worse than Kimi/Codex-style timelines during tool-heavy turns.
  - tool-call updates can arrive after later reasoning or assistant text, so the UI needs stable per-tool identity without losing the first-seen ordering.
- Decision:
  - extend the frontend assistant-message model with ordered `segments` and rebuild those segments from persisted turn events for finalized history.
  - treat `message_delta` and `reasoning_delta` as append-only text segments that coalesce only while the stream stays in the same mode.
  - treat each `tool_call` / `tool_call_update` as a stable segment keyed by `toolCallId`; later updates mutate that segment in place instead of appending a second duplicate card.
  - render assistant content segments as plain answer blocks rather than chat bubbles so visible answer text sits in the same timeline grammar as thought/tool blocks.
  - drop the IM-style left/right transcript layout in the Web UI; both user prompts and agent output render in one left-aligned reading column so a user prompt can flow directly into the agent timeline below it.
  - during streaming, keep only the currently active answer segment in the plain-text typing state; once the stream moves on, completed answer segments should immediately render through the finalized markdown pipeline so tables and lists appear without waiting for turn completion.
  - during streaming, keep only the currently active thought segment in the expanded/plain-text state; once a later event arrives, the completed thought segment should immediately become a collapsed finalized panel even before the overall turn ends.
  - apply the same panel model to tool-call segments: keep the currently updating tool call expanded during streaming, but render finalized tool calls as collapsed panels that can be manually reopened.
  - keep permission-request cards out of the ordered segment collapse model because approval prompts are independent workflow interruptions and should remain immediately visible/actionable.
  - attach copy actions to finalized assistant content segments instead of the whole assistant message footer, so copy semantics match the visible `Thought / Tool / Answer` segmentation.
  - keep those per-segment copy actions attached to their own answer blocks, but render each block's timestamp and copy control together on one local meta row with time first, to avoid burning an extra line.
  - style rendered markdown tables explicitly in answer/thought content instead of relying on browser defaults, so GFM tables remain visually recognizable inside the single-column transcript.
  - wrap rendered markdown tables in a dedicated fit-content scroll container, so table borders track the natural table width while still allowing horizontal overflow for wide tables.
  - keep plan updates outside the ordered segment list for now because `plan_update` is a replace-style progress snapshot rather than append-only assistant content.
- Consequences:
  - live streaming and history reload now show `Thought -> Tool -> Answer` in the order the agent actually emitted them.
  - finalized assistant messages can contain multiple visible answer blocks when the model alternates between tool use and visible output.
  - completed answer blocks can render markdown tables/lists during long-running turns instead of exposing raw markdown syntax until the final completion event.
  - markdown tables now read as actual tabular UI with cell boundaries and scroll behavior, especially for wider technical summaries.
  - narrow tables no longer show an oversized empty right gutter inside the table border.
  - longer in-flight tool runs no longer keep every previous thought segment expanded; completed thoughts collapse as soon as the stream moves on.
  - tool-call-heavy turns no longer leave large verbose tool cards permanently expanded in the transcript; users see a compact per-tool timeline and can reopen details selectively.
  - permission prompts continue to render outside those collapsible tool panels, so approval state is not hidden behind a disclosure control.
  - copying a multi-answer turn no longer merges unrelated answer fragments together; users can copy only the specific answer block they intend.
  - each finalized answer block now owns its own compact `time + copy` meta row, so copy stays visually attached to that block without adding another stacked control line.
  - provider-owned transcript replay still falls back to transcript-only bubbles when no turn-event history exists.
- Alternatives considered:
  - keep the old single-bubble layout and only restyle the tool/reasoning sections (rejected: still loses chronological structure).
  - flatten tool/thought events into one markdown transcript string (rejected: harder to update incrementally and loses structured tool metadata).
  - move `plan_update` into the same segment list immediately (deferred: it is replace-style state and needs separate UX rules).

## ADR-088: Derive the runtime agent list from startup preflight results

- Status: Accepted
- Date: 2026-03-20
- Context:
  - ngent supported multiple optional local agent CLIs, but startup wiring still treated the full static provider set as active.
  - as a result, servers that only had a subset of binaries installed still attempted startup config-catalog refresh for missing agents and emitted repeated `config_catalog.refresh_failed` warnings.
  - the Web UI also received unavailable agent rows even though those agents could never be used successfully in that environment.
- Decision:
  - compute the active agent set from startup `Preflight()` success only.
  - use that same derived set for `/v1/agents`, request-time `AllowedAgentIDs` validation, and startup config-catalog refresh.
  - continue logging provider-specific startup preflight diagnostics, but do not attempt model/config refresh for agents that are not runnable in the current environment.
- Consequences:
  - missing local binaries no longer trigger startup config-catalog refresh failures.
  - the frontend only shows agents that are actually usable on the running server.
  - thread creation now rejects agent ids that are unsupported in the current runtime, even if ngent knows about that provider type in other environments.
- Alternatives considered:
  - keep returning unavailable agents and suppress only refresh warnings (rejected: frontend/runtime behavior would still disagree about what is usable).
  - make refresh failures silent while leaving the static allowlist intact (rejected: still permits users to create threads for agents that cannot start).

## ADR-089: Share repeated ACP discovery and session-param helpers across built-in providers

- Status: Accepted
- Date: 2026-03-21
- Context:
  - `gemini`, `qwen`, `kimi`, and `opencode` each had an identical package-level `DiscoverModels(ctx, cfg)` implementation that only constructed the provider client and delegated to `client.DiscoverModels(ctx)`.
  - the shared behavior already lived in `internal/agents/acpcli.Client.DiscoverModels`; the duplication was only in package entrypoints.
  - those same four providers also duplicated the ACP request parameter builders for `session/new`, `session/load`, `session/list`, and the `cwd` fallback helper used by `session/list`.
- Decision:
  - add `acpcli.DiscoverModelsWithClient`, a small shared helper that constructs a provider client and invokes its `DiscoverModels` method.
  - add shared `acpcli.SessionNewParams`, `acpcli.DiscoverModelsParams`, `acpcli.SessionLoadParams`, `acpcli.SessionListParams`, and `acpcli.SessionCWD` helpers for the common ACP request shapes used by these providers.
  - route `gemini`, `qwen`, `kimi`, and `opencode` package-level `DiscoverModels` through that helper.
  - keep Kimi's local-config fast path in `internal/agents/kimi/models.go`, but use the same helper for its ACP fallback path.
  - wire `gemini`, `qwen`, `opencode`, and `kimi` to those shared ACP param helpers, leaving only genuinely provider-specific hooks local.
- Consequences:
  - package-level model discovery entrypoints now share one implementation for the repeated constructor-plus-delegate pattern.
  - the four ACP-backed providers no longer carry copy-pasted param-builder functions for the common `cwd + mcpServers + optional model/cursor/sessionId` request shapes.
  - future ACP providers with the same shape can reuse the helper instead of adding another near-identical `models.go` body.
  - Kimi still preserves its provider-specific local configuration behavior without forking the common ACP fallback.
- Alternatives considered:
  - leave the wrappers duplicated because they are small (rejected: unnecessary repetition across multiple providers).
  - leave the param builders local because they are short (rejected: they were identical across four providers and changed together conceptually).
  - move the local-config branch into `acpcli` as well (rejected: that behavior is provider-specific and should stay in `kimi`).

## ADR-090: Learn model/reasoning metadata only from real session lifecycle events

- Status: Accepted
- Date: 2026-03-22
- Context:
  - ngent had accumulated several probe-only metadata paths for model/reasoning discovery:
    - startup config-catalog refresh
    - `GET /v1/threads/{threadId}/config-options` live fallback
    - `GET /v1/agents/{agentId}/models` live fallback
  - for ACP-backed providers, those paths ultimately used `session/new`, which created provider-side empty sessions even when the user had not actually started or resumed a conversation.
  - the old agent-scoped "default catalog" assumption was also too coarse once threads could point at arbitrary existing sessions; two threads without explicit `modelId` could legitimately load different actual session configs.
- Decision:
  - treat user-initiated `session/new` / `session/load` as the authoritative source for model/reasoning metadata.
  - add a shared config snapshot callback so providers can report the `configOptions` returned by those session lifecycle calls.
  - persist those real snapshots immediately into sqlite:
    - update the thread row with the session's actual current `modelId` / `configOverrides`
    - update the agent config catalog under the actual current model id
  - remove proactive startup refresh and stop `/config-options` / `/models` endpoints from opening probe sessions just to discover metadata.
  - when a thread switches to an existing session, clear stale thread-local model/reasoning selections first so the next user-triggered `session/load` can repopulate the real values.
- Consequences:
  - ngent no longer creates empty provider sessions just because the UI opened a thread or the server started up.
  - fresh threads now show no model/reasoning metadata until at least one real turn (or resumed session turn) reports it.
  - model/reasoning controls in the Web UI become session-driven instead of agent-default-driven, which avoids cross-thread leakage from a shared default snapshot.
  - `/v1/agents/{agentId}/models` may legitimately return an empty list until sqlite has learned at least one real config snapshot for that agent.
- Alternatives considered:
  - keep proactive refresh and only hide the frontend controls (rejected: still creates empty upstream sessions and keeps sqlite detached from real session state).
  - keep using one shared default catalog row for threads without explicit `modelId` (rejected: stale or unrelated session config can leak between threads).
  - update only the agent catalog and leave thread rows untouched (rejected: thread-level model/reasoning state would stay ambiguous when multiple sessions of the same agent differ).

## ADR-091: Persist learned config snapshots per provider session

- Status: Accepted
- Date: 2026-03-22
- Context:
  - ADR-090 removed probe sessions and made real `session/new` / `session/load` the only source of config metadata.
  - ngent persisted those learned snapshots into the thread row and `agent_config_catalogs`, but a later switch to another existing session intentionally clears stale thread-local `modelId` / `configOverrides`.
  - when the user switched back to the original session before sending another turn, ngent only had the `sessionId`; the learned model/reasoning snapshot was no longer addressable, so the Web UI hid the controls again.
- Decision:
  - add sqlite table `session_config_cache(agent_id, cwd, session_id, config_options_json, updated_at)`.
  - whenever a user-triggered turn/session load reports config for a bound session, persist the normalized snapshot under that `agent + cwd + sessionId` key in addition to the thread row and `agent_config_catalogs`.
  - if a fresh session reports config before `session_bound`, replay the already persisted thread/model snapshot into `session_config_cache` as soon as the session id becomes known.
  - when session-scoped history is served from cached transcript data but the destination session still has no cached config snapshot, bypass the transcript-only short circuit and perform one live `session/load` so that user-triggered session switching can still teach sqlite that session's config.
  - change `GET /v1/threads/{threadId}/config-options` to restore from session cache when the thread currently points at a known session but has no thread-local `modelId`.
- Consequences:
  - switching away from a learned session and then back to it restores model/reasoning controls immediately without requiring another turn.
  - switching directly onto an unseen existing session can now also reveal model/reasoning controls immediately when that session's user-triggered `session/load` returns config metadata.
  - session-specific config state no longer depends on thread-local mirrors surviving every session change.
  - unseen sessions still legitimately return no config metadata until one real turn (or other real session lifecycle event during a turn) teaches ngent that session's snapshot.
- Alternatives considered:
  - keep only thread-local state and accept that switching back requires another turn (rejected: poor UX and contradicts the intent of session-driven discovery).
  - repopulate by proactively calling `session/load` on session switch (rejected: reintroduces probe-style metadata fetches outside real turn flow).
  - key everything only by `agent + model` (rejected: multiple sessions can share a model id but differ in other config values such as reasoning/mode).

## ADR-055: Preserve non-text assistant ACP content as first-class turn events

- Status: Accepted
- Date: 2026-03-22
- Context:
  - ACP `agent_message_chunk` is not limited to plain text; providers can emit structured visible assistant content such as image blocks or embedded resources.
  - ngent previously treated ordinary assistant message chunks as text-only, so any non-text payload was silently ignored unless it happened to appear inside a tool-call payload.
  - flattening those blocks into markdown or ad-hoc strings would lose protocol structure and make future richer rendering harder.
- Decision:
  - extend shared ACP update parsing so text `agent_message_chunk` payloads continue to populate `message_delta`, while non-text assistant blocks are preserved as raw structured `content`.
  - add a per-turn assistant message-content callback in the agent layer and persist/stream those blocks as first-class `message_content` turn events.
  - keep `responseText` as the aggregate visible text only, and let rich assistant content be reconstructed from persisted turn events in the Web UI.
  - render common assistant image/resource blocks directly in the Web UI while leaving unknown block shapes on a JSON fallback path.
- Consequences:
  - hub-created turns no longer lose assistant images or embedded resources during live streaming or history reload.
  - downstream clients can evolve richer renderers without changing the transport contract again because the raw ACP `content` JSON is preserved.
  - provider-owned transcript replay remains text-only for now because its transcript schema has not yet grown structured content support.
- Alternatives considered:
  - ignore non-text assistant blocks outside tool calls (rejected: visibly loses model output).
  - stringify image/resource payloads into `message_delta` (rejected: destroys structure and produces poor UI).
  - change `responseText` into a heterogeneous rich-content blob (rejected: larger API break and unnecessary when turn events already provide the ordered detail).

## ADR-056: Preserve exact provider permission options through the hub permission flow

- Status: Accepted
- Date: 2026-03-22
- Context:
  - ngent's permission UI and HTTP endpoint originally collapsed every permission request into a binary `approved` / `declined` choice.
  - real ACP-backed providers such as Kimi, Qwen, and OpenCode can advertise richer permission option sets such as `allow_once`, `allow_always`, `reject_once`, and `reject_always`, each with a distinct `optionId`.
  - once the hub reduced those requests to a generic outcome, the Web UI could no longer present the provider's real choices and the backend could only approximate the response by picking a default allow/reject option.
- Decision:
  - extend the shared permission request model so provider-advertised `options[]` are preserved and forwarded in the `permission_required` SSE payload.
  - extend `POST /v1/permissions/{permissionId}` to accept `optionId` alongside the existing generic `outcome`.
  - when a client submits an exact `optionId`, prefer forwarding that exact selection back to option-aware providers; keep generic outcome fallback behavior for providers that still only understand `approved` / `declined` / `cancelled`.
  - keep fail-closed behavior unchanged: missing, invalid, timed-out, or disconnected permission decisions still default to deny.
- Consequences:
  - the Web UI can render all provider permission choices instead of hard-coding `Allow` / `Deny`.
  - Kimi/Qwen/OpenCode-style permission flows preserve provider-specific semantics such as `allow always` vs `allow once`.
  - existing outcome-only clients and generic providers remain compatible.
- Alternatives considered:
  - keep the UI binary and continue mapping approvals to the first allow/reject option (rejected: loses provider semantics and hides real choices from the user).
  - expose provider options in SSE but keep the HTTP endpoint outcome-only (rejected: UI would still be unable to return the exact selected option).

## ADR-057: Persist Web UI uploads as local temp files and forward them as ACP resource links

- Status: Accepted
- Date: 2026-03-23
- Context:
  - ACP `session/prompt` supports structured prompt content, including `resource_link`, so users can send text together with local files/images.
  - ngent's turn pipeline previously accepted only one plain `input` string from HTTP through to the agent layer, which meant Web UI uploads had no transport path.
  - base64-inlining uploaded files into prompt text would bloat requests, lose MIME/name metadata, and diverge from ACP's native content model.
- Decision:
  - introduce a shared structured prompt model in the agent layer with `text` and `resource_link` items plus a plain-text fallback renderer for non-ACP paths.
  - extend `POST /v1/threads/{threadId}/turns` to accept `multipart/form-data`; persist uploaded files into the local temp directory and convert them into `file://` ACP resource links with `name`, `mimeType`, and `size`.
  - keep `requestText` as a readable fallback summary that mentions attached resources, and persist the original structured user prompt separately as a `user_prompt` turn event so history/UI can reconstruct attachments without polluting the visible user text.
  - teach ACP CLI and embedded providers to send structured `session/prompt.prompt[]` arrays instead of forcing plain strings.
- Consequences:
  - Web UI users can now send attachment-only turns and text+attachment turns without leaving ACP semantics.
  - hub-created history preserves enough information to re-render user attachment cards after reload.
  - temp files remain local-only and provider-facing through `file://` URIs; ngent does not need to expose a new download endpoint.
- Alternatives considered:
  - upload files to a new ngent HTTP asset endpoint and send remote URLs (rejected: unnecessary extra surface and weaker local-first story).
  - embed file contents directly into prompt text (rejected: loses ACP structure, MIME metadata, and scales poorly for binary files).
  - add a separate pre-upload API that returns opaque ids (rejected: more round trips and more state than needed for the current Web UI flow).

## ADR-058: Render bracketed inline base64 user-image placeholders as safe Web UI previews

- Status: Accepted
- Date: 2026-03-26
- Context:
  - some user messages now arrive with inline image payloads serialized directly into the visible text as placeholders beginning with `[Image: data:image/...;base64,...]`.
  - the Web UI previously fed the entire user message through markdown rendering unchanged, so those placeholders appeared as long raw base64 blobs inside the chat bubble.
  - attachment cards already cover structured uploads, but these placeholder-style images live inside ordinary `msg.content` text and therefore need a separate render path.
- Decision:
  - add a user-message-only renderer that scans for bracketed inline placeholders matching `data:image/*;base64,...`.
  - accept only image data URLs, normalize incidental whitespace out of the base64 payload, and render each valid placeholder as an inline preview image inside the existing user bubble.
  - keep all surrounding text on the existing markdown path and leave malformed or unsupported placeholders untouched as literal text.
  - preserve the original raw message string for copy actions and stored history; the change is presentation-only.
- Consequences:
  - user messages containing inline base64 image placeholders become immediately readable in the Web UI without changing backend contracts.
  - the renderer remains fail-closed against non-image `data:` payloads because only `data:image/*;base64,...` is promoted to an `<img>`.
  - markdown behavior around ordinary text stays consistent, though placeholder parsing is intentionally targeted at the bracketed image form rather than arbitrary embedded binary text.
- Alternatives considered:
  - require upstream senders to convert these images into structured attachments first (rejected: existing message streams already contain the placeholder form).
  - render every `data:` URL found in user text as an image (rejected: too permissive and unsafe).
  - leave the placeholder as raw text and rely on copy/paste elsewhere (rejected: poor UX for image-bearing prompts).

## ADR-059: Store uploaded attachments under the configurable data directory and serve them back through a stable attachment route

- Status: Accepted
- Date: 2026-03-26
- Context:
  - ADR-057 deliberately used local `file://` uploads so ACP providers could read user attachments without introducing a remote object store.
  - that first implementation stored uploads in the OS temp directory and only persisted the raw `file://` resource link into history.
  - temp-directory storage was too fragile for long-running local-first usage because the OS can clean it at any time, and browsers cannot reliably keep rendering persisted `file://` image/resource links after the live in-memory `blob:` preview disappears.
- Decision:
  - replace the CLI `--db-path` flag with `--data-path`, with default root `$HOME/.ngent/`.
  - derive sqlite as `data-path/ngent.db` and store uploaded files under `data-path/attachments/<category>/`, where category is chosen from MIME/extension families such as `images`, `documents`, `text`, `audio`, `video`, `archives`, and `files`.
  - persist uploaded attachment metadata in sqlite as `turn_attachments(attachment_id, turn_id, name, mime_type, size, file_path, created_at)`.
  - include a stable `attachmentId` in persisted `user_prompt` events and serve persisted files back to the Web UI through `GET /attachments/{attachmentId}` with client ownership checks and optional query-token auth for image tags.
- Consequences:
  - uploaded attachments now survive service restarts and OS temp cleanup.
  - Web UI attachment cards continue to render after stream completion and after history reload because they no longer depend on ephemeral `blob:` URLs.
  - the CLI data model is clearer: one configurable data root now owns both sqlite state and uploaded-file state.
  - ngent still needs a later janitor policy for old attachment files inside the data directory.
- Alternatives considered:
  - keep `--db-path` and add a second unrelated attachment-root flag (rejected: splits one local state root across two flags and makes defaults harder to reason about).
  - keep storing uploads in temp and only persist extra preview metadata (rejected: does not solve OS cleanup and still leaves the attachment file itself non-durable).
  - expose raw absolute file paths directly to the browser (rejected: browsers cannot reliably use persisted `file://` paths from the HTTP Web UI, and a routed attachment fetch keeps ownership/auth checks server-side).

## ADR-060: Make thread history session-scoped for Web UI session switches

- Status: Accepted
- Date: 2026-03-26
- Context:
  - the Web UI session picker reconstructs the selected chat by combining provider-owned transcript replay with ngent-owned persisted turn history.
  - before this change, every session switch still fetched `GET /v1/threads/{threadId}/history?includeEvents=1` for the entire thread and only then filtered the turns client-side by `session_bound`.
  - on real Codex threads this became expensive enough to stall the UI: one thread with 21 turns produced roughly 19 MB of history JSON and about 42k persisted events, even though the target session needed only a single turn.
- Decision:
  - add an optional `sessionId` query parameter to `GET /v1/threads/{threadId}/history`.
  - apply session filtering on the server using the same rules the Web UI already depended on:
    - match turns by their latest `session_bound` event.
    - if the thread has no annotated turns at all, keep returning non-ephemeral history instead of hiding legacy data.
    - include unannotated legacy turns only when the thread has exactly one annotated session and it matches the requested `sessionId`.
    - continue dropping cancelled turns that have no visible response text.
  - update the Web UI to request `history?includeEvents=1&sessionId=...` when a concrete session is selected, while keeping the local filter as a compatibility fallback.
- Consequences:
  - session switches no longer require the browser to parse and walk the full thread history when only one historical session is being opened.
  - the frontend keeps its existing session-reconstruction behavior, but the largest JSON payload and event traversal now happen server-side where they do not block the UI thread.
  - the base `/history` endpoint remains unchanged for callers that still need whole-thread history.
- Alternatives considered:
  - keep whole-thread history and add more frontend caching only (rejected: first-open session switches from `New session` would still pay the full parse/render cost).
  - rely only on provider transcript replay for session views (rejected: loses ngent-owned rich turn artifacts such as persisted reasoning, tool-call, and other turn events).
  - add a brand-new endpoint just for session-filtered history (rejected: the existing `/history` contract already fit the need with one optional query parameter).

## ADR-061: Compact historical delta runs on read and render large chats incrementally

- Status: Accepted
- Date: 2026-03-26
- Context:
  - session-scoped `/history` removed unrelated sessions from the payload, but old databases still stored every historical `message_delta` / `reasoning_delta` row separately.
  - the provided Codex repro still showed residual UI jank on the selected session because the browser had to replay a delta-heavy turn and then rebuild the full chat DOM in one synchronous step.
  - write-side event merging alone was insufficient because it only helps turns created after the fix ships.
- Decision:
  - merge consecutive same-turn `message_delta`, `reasoning_delta`, and `thought_delta` runs when serializing `/v1/threads/{threadId}/history?includeEvents=1`.
  - keep preserving boundaries when event type changes so ordering relative to `tool_call*`, `plan_update`, `session_bound`, and other event kinds stays intact.
  - also merge those same delta runs at storage write time in `AppendEvent(...)` so new turns persist fewer rows.
  - render large Web UI message lists incrementally across animation frames instead of always doing one synchronous `innerHTML = msgs.map(renderMessage).join('')` rebuild.
- Consequences:
  - old databases immediately benefit from smaller history payloads and fewer replay events without a migration step.
  - new databases avoid accumulating the same redundant delta rows over time.
  - session switches remain responsive even when one persisted turn contains many delta updates, because the browser now has less replay work and can yield while rebuilding the chat.
- Alternatives considered:
  - migrate existing SQLite data in place (rejected: more invasive, riskier on user state, and unnecessary when read-time compaction solves the API cost directly).
  - compact every event type aggressively, including `tool_call_update` (rejected for now: those updates can carry meaningful intermediate state transitions that should stay ordered until a clearer merge contract is defined).
  - keep synchronous DOM rebuilds and rely only on smaller payloads (rejected: large rendered transcripts can still block the main thread even after history parsing becomes cheap).

## ADR-062: Decouple viewed Web UI session from backend thread session during active turns

- Status: Accepted
- Date: 2026-03-26
- Context:
  - the Web UI session picker previously reused `thread.agentOptions.sessionId` as both the backend agent scope and the frontend "currently viewed chat" selection.
  - when a turn was still active, switching to another session always tried `PATCH /v1/threads/{threadId}` and the server correctly rejected it with `409 thread has an active turn`.
  - users still need to inspect older sessions while another session is waiting on a long-running response.
- Decision:
  - add a frontend-only per-thread session-selection override that can temporarily differ from `thread.agentOptions.sessionId`.
  - while a thread has any active stream, session switching updates only the visible chat/history scope in the browser and skips the backend patch.
  - once the stream finishes, or immediately before the next send when no stream is active, sync the selected session back into `PATCH /v1/threads/{threadId}` so future turns still execute in the session currently shown in the UI.
  - represent unsaved "New session" views with a stable local `@fresh:<nonce>` selection so message/history caches keep a distinct fresh-session scope before ACP emits `session_bound`.
- Consequences:
  - users can browse away from a waiting session and later return without seeing a conflict alert.
  - backend concurrency semantics stay unchanged because thread/session binding is still only mutated when the thread is idle.
  - the frontend now owns a small amount of ephemeral session-selection state and must clear it when thread state catches up or the thread is deleted.
- Alternatives considered:
  - keep the old behavior and surface the 409 as a user error (rejected: browsing another session is a read-only UI action and should not feel blocked by an unrelated active turn).
  - relax the server to accept session-changing `PATCH` during active turns (rejected: would violate the existing whole-thread conflict model and risks mutating active agent scope mid-turn).

## ADR-064: Share threads and sessions across browser-scoped client IDs on the same ngent instance

- Status: Accepted
- Date: 2026-03-27
- Context:
  - browser-local `clientId` generation made the same ngent data look artificially partitioned by browser profile, so threads created in Safari were invisible from another browser even though both were talking to the same local service and sqlite database.
  - the product request is for one local shared workspace per ngent instance: thread list, session list, permissions, and persisted attachments should all be visible regardless of which browser-scoped `clientId` originated them.
  - ngent still benefits from a caller identifier for API compatibility, so dropping `X-Client-ID` entirely would be a larger breaking change than needed.
- Decision:
  - keep requiring `X-Client-ID` on `/v1/*` requests as a compatibility header, but stop persisting it in sqlite.
  - stop using `clientId` as an authorization/tenancy gate for thread-scoped APIs, permission resolution, or persisted attachment fetches.
  - list threads globally from sqlite and resolve thread access by `threadId` only.
  - remove `threads.client_id`, drop the obsolete `clients` table, and make `recent-directories` global instead of browser-scoped.
  - remove the Web UI's visible Client ID display/reset controls and send one fixed compatibility header internally from the browser.
- Consequences:
  - multiple browsers connected to the same ngent instance now see the same thread/session state and can continue the same conversation without copying browser-local IDs around.
  - `X-Client-ID` is now a compatibility field, not a persisted identity record and not a security boundary.
  - operators who need true user isolation must use separate ngent instances, separate `--data-path` values, or an external auth/proxy layer instead of relying on browser-local client IDs.
- Alternatives considered:
  - force users to copy one shared `clientId` across browsers (rejected: brittle manual workaround and still couples visibility to browser-local storage).
  - remove `X-Client-ID` from the API entirely (rejected for now: unnecessary breakage for existing clients when header-compatibility is enough).
  - keep thread ownership but special-case only `/sessions` (rejected: inconsistent UX because the main thread list would still disappear across browsers).
