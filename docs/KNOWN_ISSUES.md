# KNOWN ISSUES

## Issue Template

```text
- ID: KI-XXX
- Title:
- Status: Open | Mitigated | Closed
- Severity: Low | Medium | High
- Affects:
- Symptom:
- Workaround:
- Follow-up plan:
```

## Active Issues

Status is authoritative for each entry. This section keeps only open and
mitigated items so active work does not share space with already closed issues.

- ID: KI-051
- Title: Grouped thread session-collapse state resets after a full page reload
- Status: Open
- Severity: Low
- Affects: browser users who manually collapse one or more thread groups in the embedded Web UI left rail
- Symptom:
  - each thread group now supports collapsing its inline session list from the leading agent glyph, but that state currently lives only in browser-local runtime memory.
  - refreshing the page or reopening the Web UI expands every thread group again, even if the user had intentionally collapsed several of them before.
- Workaround:
  - re-collapse the desired thread groups after reload; the behavior remains stable during ordinary in-page re-renders while the browser tab stays open.
- Follow-up plan:
  - evaluate whether the collapse map should be persisted in browser local storage without turning it into shared server-side thread metadata.

- ID: KI-050
- Title: Cross-client session-row activity spinner is refreshed on fetch, not pushed live to already-open browsers
- Status: Open
- Severity: Low
- Affects: secondary browsers that already have the grouped thread/session rail open before another browser starts a new turn
- Symptom:
  - `GET /v1/threads/{threadId}/sessions` now marks concrete session rows with `isActive`, so a newly opened browser or a manual session-list refresh shows the correct spinner for the active session.
  - however, the grouped rail still does not subscribe to a separate push channel for background thread/session activity, so an already-open browser will not see another browser's new session-row spinner until it reloads or re-fetches that thread's sessions.
- Workaround:
  - reload the page, reopen the relevant thread, or use the thread menu's session refresh action.
- Follow-up plan:
  - evaluate whether the grouped rail needs lightweight polling or a dedicated push/resume path for cross-browser background session-activity updates.

- ID: KI-049
- Title: Live plan overlay is only pinned for the currently open session
- Status: Open
- Severity: Low
- Affects: operators watching multiple active sessions/threads at once in the embedded Web UI
- Symptom:
  - the new live plan card stays pinned at the bottom only for the session currently open in the chat pane.
  - if another session/thread is also streaming in the background, its live plan is not shown globally until the user switches into that conversation.
- Workaround:
  - switch into the running session whose plan you want to inspect; the card rehydrates from persisted `plan_update` events and resumes live updates.
- Follow-up plan:
  - evaluate whether the grouped left rail eventually needs a compact cross-thread live-plan indicator for background sessions.

- ID: KI-059
- Title: Git-diff drawer preview is intentionally limited to tracked patches and new text files
- Status: Open
- Severity: Low
- Affects: Web UI users opening the per-file git-diff drawer for staged-only additions, binary assets, or other non-text untracked files
- Symptom:
  - tracked working-tree changes preview correctly through `git diff`, and new untracked text files preview through their current file contents.
  - binary/non-text untracked rows stay disabled, and staged-only additions that do not appear in the existing working-tree diff summary still have no drawer entry.
- Workaround:
  - inspect binary assets with external tooling, or stage/commit/re-diff through a workflow that makes the needed change visible in the current working-tree summary first.
- Follow-up plan:
  - evaluate whether the git-diff surface should eventually cover staged-only previews or a separate asset viewer without expanding the current local-only safety boundary.

- ID: KI-048
- Title: Git-diff file icons currently cover a curated subset of common file types
- Status: Open
- Severity: Low
- Affects: expanded Web UI git-diff rows whose file names/extensions are not yet in the frontend icon map
- Symptom:
  - common code/content/config files such as `*.py`, `*.ts`, `*.tsx`, `*.go`, `*.json`, `*.md`, `*.yml`, `Dockerfile`, and `go.mod` now show specific type icons.
  - rarer file types that are not yet mapped still fall back to the generic file icon in the expanded diff panel.
- Workaround:
  - rely on the visible file path text; unknown types still render correctly, just without a specialized icon.
- Follow-up plan:
  - extend the curated basename/extension map when real usage shows additional file types that need first-class icons.

- ID: KI-047
- Title: Web UI locale switch does not translate server-originated error text
- Status: Open
- Severity: Low
- Affects: browser users viewing HTTP/API/provider error messages inside the localized Web UI
- Symptom:
  - Settings, navigation, empty states, composer controls, and other client-owned Web UI strings switch between English, Simplified Chinese, Spanish, and French.
  - errors that come straight from the backend or upstream agent/provider can still appear in English because the frontend currently renders those payload strings verbatim.
- Workaround:
  - rely on the surrounding localized UI and exact error code/message, or switch to English if a support/debugging workflow expects backend wording unchanged.
- Follow-up plan:
  - evaluate a stable frontend error-code translation layer, or backend-side localization, if mixed-language error surfaces become a product problem.

- ID: KI-046
- Title: ACP session-usage visibility still depends on provider support
- Status: Open
- Severity: Low
- Affects: threads backed by providers that do not emit ACP prompt `usage` and/or `usage_update`
- Symptom:
  - sqlite `session_usage_cache` and the Web UI footer indicator update only when the upstream provider actually returns session-usage data.
  - on providers with no usage support, the composer footer stays unchanged and there is no per-session token/context snapshot to query later.
- Workaround:
  - none inside ngent today; use provider-native billing/usage surfaces when available.
- Follow-up plan:
  - keep watching provider ACP compatibility and add support automatically when a provider starts emitting usage.
  - evaluate optional provider-specific approximations only if product requirements later demand a universal badge.

- ID: KI-062
- Title: Embedded Pi provider currently inherits adapter-side MCP and plan gaps
- Status: Open
- Severity: Low
- Affects: threads backed by `pi`
- Symptom:
  - the new embedded Pi integration supports session startup, prompt streaming, session list/load, config options, slash commands, and permission bridging through `pkg/piacp`.
  - however, current Pi support in `acp-adapter` still does not expose ACP MCP routing or plan snapshots, so ngent cannot surface those features on Pi threads yet.
- Workaround:
  - use Codex or another provider when a workflow depends on MCP routing or live plan updates.
- Follow-up plan:
  - track future `acp-adapter` Pi capability expansion and wire those surfaces through ngent when the upstream runtime exposes them.

- ID: KI-045
- Title: Web UI git branch badge can lag behind out-of-band repository changes
- Status: Open
- Severity: Low
- Affects: repository-backed threads when the worktree branch changes outside ngent while the same thread remains open and idle
- Symptom:
  - the composer footer branch pill refreshes on thread open, after ngent turn completion, and after in-UI branch switches, but it does not continuously poll the filesystem.
  - if another terminal or tool checks out a different branch while the thread stays open, the visible badge can remain stale until the next refresh point.
- Workaround:
  - reload the page, switch away and back to the thread, or let the next turn complete to refresh the displayed branch state.
- Follow-up plan:
  - evaluate an explicit refresh affordance or low-cost polling strategy if out-of-band branch switching becomes common in real usage.

- ID: KI-043
- Title: `X-Client-ID` no longer isolates data between callers on the same ngent instance
- Status: Open
- Severity: Medium
- Affects: deployments that expected separate browsers or separate users behind one ngent instance to have isolated thread/session state
- Symptom:
  - thread list, thread metadata, session views, permission resolution, and persisted attachment fetches are shared across all callers that can reach the same ngent instance.
  - changing browser profiles changes the local `clientId`, but no longer hides or fences off previously created threads.
- Workaround:
  - run separate ngent instances or separate `--data-path` roots for each isolation boundary you need.
  - if the service is exposed beyond localhost, put it behind stronger caller authentication/authorization than `X-Client-ID`.
- Follow-up plan:
  - evaluate whether ngent should eventually add an explicit optional namespace/project isolation flag instead of relying on browser-local identifiers.

- ID: KI-044
- Title: Grouped left rail can still become dense on threads with very large provider session catalogs
- Status: Open
- Severity: Low
- Affects: Web UI threads whose providers return long session lists for the current working directory
- Symptom:
  - the redesign now groups sessions directly under each thread header, but the left rail is still a straight scrollable list with truncation and no in-rail filtering/search.
  - on threads with many historical sessions, scanning the grouped rail can still feel dense even though the separate session drawer is gone.
- Workaround:
  - rely on the provider's recent-session ordering near the top of each thread group and use `Show more` only when older sessions are needed.
- Follow-up plan:
  - evaluate adding optional session search/grouping if large-history threads remain common in real usage.

- ID: KI-035
- Title: Desktop-workbench Web UI visuals vary slightly by host browser/font stack
- Status: Open
- Severity: Low
- Affects: the embedded Web UI on different desktop/mobile browsers and operating systems
- Symptom:
  - assistant reply prose now ships with a bundled `Inter Variable` body font for closer Kimi-style transcript rendering, but CJK fallback glyphs, code rendering details, and the rest of the shell chrome still depend on host/browser font rasterization.
  - typography, subtle contrast steps, and spacing can therefore still look slightly different across browsers, operating systems, and displays even though behavior stays the same.
- Workaround: use a modern browser with good font rendering; functional behavior remains unchanged even when the exact presentation shifts slightly
- Follow-up plan: consider lightweight visual regression snapshots and, only if exact parity becomes a product requirement, broader bundled font coverage beyond the assistant-body surface

- ID: KI-034
- Title: Human-readable stderr logs are less machine-friendly than JSON logs
- Status: Open
- Severity: Low
- Affects: deployments that ingest ngent stderr into pipelines expecting one JSON object per line
- Symptom: request/access logs and ACP debug traces are emitted as readable text lines, so strict JSON log collectors cannot parse them directly
- Workaround: use a text log parser in the collector pipeline, or pin an earlier revision if JSON log envelopes are a hard requirement
- Follow-up plan: consider an opt-in `--log-format=json|pretty` switch if operator demand for machine-readable logs returns

- ID: KI-001
- Title: SSE viewer disconnect can still cause a temporary live-view gap
- Status: Mitigated
- Severity: Medium
- Affects: streaming clients on unstable links
- Symptom:
  - a dropped viewer connection no longer cancels the underlying turn, but that specific viewer still stops receiving live tokens until it reconnects.
  - if the network silently stalls rather than hard-closing, the browser can still experience a short gap before reconnect logic notices and resumes tailing the active turn.
- Workaround:
  - refresh/reopen the same thread or reconnect with `GET /v1/turns/{turnId}/events?after=<lastSeq>`.
  - the underlying turn keeps running unless the user explicitly cancels it or it times out while waiting on permission.
- Follow-up plan:
  - add heartbeat/idle detection so dead viewer connections are noticed faster and reconnect latency is more predictable.

- ID: KI-002
- Title: Permission decision timeout
- Status: Open
- Severity: Medium
- Affects: slow/offline client decision path
- Symptom: pending permission expires and turn is fail-closed (`outcome=declined`), typically ending with `stopReason=cancelled`
- Workaround: respond quickly to `permission_required`; custom permission integrations should keep any adapter-side timeout aligned with the hub timeout.
- Follow-up plan: expose timeout metadata in SSE payload and add client-side countdown UX

- ID: KI-003
- Title: SQLite lock contention under burst writes
- Status: Open
- Severity: Medium
- Affects: high-concurrency turn/event persistence
- Symptom: transient `database is locked` errors
- Workaround: enable WAL, busy timeout, and retry with jitter
- Follow-up plan: benchmark and tune connection settings in M2 and M8

- ID: KI-005
- Title: External agent process crash
- Status: Open
- Severity: High
- Affects: ACP/Codex provider turns
- Symptom: turn aborts unexpectedly; stream ends with provider error
- Workaround: detect process exit quickly, persist failure event, allow user retry
- Follow-up plan: supervised restart and backoff policy in M6 and M8

- ID: KI-006
- Title: Permission request races with SSE disconnect
- Status: Open
- Severity: Medium
- Affects: clients that close stream while permission is pending
- Symptom: decision endpoint may return `404/409` after auto fail-closed resolution
- Workaround: reconnect and inspect turn history terminal state; treat stale `permissionId` as non-retriable
- Follow-up plan: extend the existing `permission_resolved` event with an explicit machine-readable reason such as `timeout|disconnect|client_decision`

- ID: KI-007
- Title: Embedded codex runtime prerequisite mismatch
- Status: Open
- Severity: Medium
- Affects: deployments enabling embedded codex provider
- Symptom: codex turns fail when `codex app-server` prerequisites/auth/environment are not ready even though server binary is correctly configured
- Workaround: verify codex CLI/app-server availability and auth state before issuing codex turns; inspect startup preflight and turn error logs
- Follow-up plan: add richer preflight diagnostics and compatibility matrix checks for codex CLI vs linked `acp-adapter` module versions

- ID: KI-008
- Title: Character-based context budgeting can diverge from token budgets
- Status: Open
- Severity: Medium
- Affects: long multilingual threads with high token/char variance
- Symptom: prompt fits `context-max-chars` but may still be too large for model token limits
- Workaround: reduce `--context-max-chars` conservatively and run compact more frequently
- Follow-up plan: replace char-based policy with model-aware token estimation in M8

- ID: KI-009
- Title: Embedded codex local state/schema drift warnings
- Status: Open
- Severity: Medium
- Affects: real local embedded codex runs that depend on user `~/.codex` state and app-server version capabilities
- Symptom: stderr may show warnings like `state_5.sqlite migration ... missing` and endpoint compatibility errors such as `mcpServer/call unknown variant`; turn usually still completes but tool output can be empty
- Workaround: align local codex CLI/app-server version with linked `acp-adapter` schema expectations, and repair/reset local codex state DB when migration drift appears
- Follow-up plan: add explicit diagnostics/preflight endpoint to surface local state/schema compatibility before turn execution

- ID: KI-010
- Title: Qwen ACP environment/auth dependency
- Status: Open
- Severity: Medium
- Affects: implemented `qwen --acp` provider turns in constrained environments
- Symptom:
  - in sandboxed or permission-restricted environments, Qwen can fail before ACP initialize completes when local runtime files under `~/.qwen` are not writable.
  - hub-side symptom in those environments typically converges as `qwen: initialize: qwen: connection closed`.
  - prompt execution can still fail with upstream/internal errors when auth or network is not ready, even after handshake succeeds.
- Workaround:
  - ensure writable home/config directory for qwen runtime (`HOME`, `~/.qwen`).
  - ensure qwen authentication is completed and network path to model backend is available before turn execution.
  - run real qwen validation outside restrictive sandboxes when verifying `session/list` / `session/load` behavior against the locally installed CLI.
- Follow-up plan:
  - add clearer preflight diagnostics for qwen runtime prerequisites (filesystem writable check + auth hints).
  - map common qwen upstream errors to stable hub error details for easier operator debugging.

- ID: KI-011
- Title: Thread deletion is irreversible
- Status: Open
- Severity: Low
- Affects: users deleting historical threads via API/Web UI
- Symptom: deleting a thread permanently removes its thread/turn/event history and cannot be restored through server APIs.
- Workaround: export needed history before delete.
- Follow-up plan: evaluate optional soft-delete retention window and admin-only restore endpoint if product requirements demand recoverability.

- ID: KI-012
- Title: Model override id is accepted as free text
- Status: Open
- Severity: Low
- Affects: thread create/update flows using `agentOptions.modelId`
- Symptom: direct API clients can still submit any `modelId`; unsupported values fail later during provider/runtime execution instead of being rejected at create/update time.
- Workaround: query `GET /v1/agents/{agentId}/models` first, then submit a returned model id (or omit `modelId` to use provider default model).
- Follow-up plan: add optional server-side validation in create/update path against runtime-discovered model catalogs.

- ID: KI-013
- Title: Stdio providers apply config in transient ACP sessions
- Status: Open
- Severity: Low
- Affects: `opencode` / `qwen` / `gemini` / `kimi` thread config behavior
- Symptom: these providers are process-per-turn; config changes are mirrored into persisted thread metadata (`agentOptions.modelId` + `agentOptions.configOverrides`) and reapplied on the next ACP session, but there is still no long-lived runtime session to mutate between turns.
- Workaround: none required for normal usage; persisted config selections remain effective on future turns through thread metadata replay.
- Follow-up plan: evaluate persistent per-thread ACP runtime for stdio agents if future product requirements need truly in-session config mutations beyond thread-level replay.

- ID: KI-014
- Title: Web UI surfaces only model and reasoning config controls
- Status: Open
- Severity: Low
- Affects: advanced ACP config categories beyond `model` and `reasoning`
- Symptom: thread `configOptions` can contain additional categories/ids, but the composer footer currently exposes first-class controls only for model and reasoning.
- Workaround: use `GET/POST /v1/threads/{threadId}/config-options` directly to inspect and update other ACP config options.
- Follow-up plan: add a generic advanced-config surface in the Web UI if users need broader access to ACP session settings.

- ID: KI-015
- Title: Partial startup refresh keeps stale catalog rows for failed models
- Status: Open
- Severity: Low
- Affects: persisted `agent_config_catalogs` when startup background refresh succeeds for some models but fails for others
- Symptom: on partial refresh failure, the server intentionally keeps older sqlite catalog rows for models that could not be refreshed, so removed/changed upstream model metadata can remain temporarily stale until a later successful refresh or an explicit config change rewrites that model row.
- Workaround: restart again after upstream/provider health is restored, or trigger a config change on the affected model/thread so the latest snapshot is written through immediately.
- Follow-up plan: evaluate adding per-agent refresh status/age diagnostics in the API or Web UI so operators can see when catalog data is partially stale.

- ID: KI-016
- Title: Thread rename is blocked during an active turn
- Status: Open
- Severity: Low
- Affects: API/Web UI rename requests using `PATCH /v1/threads/{threadId}` with `title`
- Symptom: rename requests return `409 CONFLICT` while the thread is actively streaming because title updates share the same active-turn lock as other thread mutations.
- Workaround: wait for the current turn to finish or cancel it, then retry rename.
- Follow-up plan: evaluate whether title-only updates should continue using the shared mutation lock or move to a narrower metadata-only guard in a future revision.

- ID: KI-019
- Title: `item/tool/requestUserInput` currently uses fallback auto-selection
- Status: Open
- Severity: Medium
- Affects: codex app-server flows that require real user-entered answers or multi-choice semantics beyond first-option selection
- Symptom: adapter now avoids hard `-32000` errors, but for `requestUserInput` it auto-selects first option labels and does not expose full interactive question UI in hub frontend.
- Workaround: prefer MCP/tool flows that do not require complex interactive follow-up prompts; if needed, run the same operation in an environment with native codex app UI support.
- Follow-up plan: add first-class user-input request bridge and frontend interaction model for arbitrary question/option responses.

- ID: KI-020
- Title: Kimi ACP runtime/auth prerequisites and command-form drift
- Status: Open
- Severity: Medium
- Affects: implemented `kimi` provider turns in environments with uninitialized Kimi CLI state
- Symptom:
  - Kimi upstream docs currently show both `kimi acp` and `kimi --acp`; the hub retries both, but a local CLI that exits before ACP initialize still surfaces as a provider startup failure.
  - prompt execution can fail after handshake when local Kimi authentication or upstream network access is not ready.
- Workaround:
  - ensure the local `kimi` CLI is installed and already logged in before issuing turns.
  - inspect startup preflight and turn error logs when ACP mode closes immediately.
- Follow-up plan:
  - add richer preflight diagnostics for Kimi auth/runtime readiness beyond PATH existence.
  - keep validating future Kimi CLI releases and narrow the fallback path once upstream command syntax stabilizes.

- ID: KI-021
- Title: Resumed ACP sessions do not backfill prior transcript into hub history
- Status: Open
- Severity: Medium
- Affects: threads that select an existing ACP `sessionId` from the Web UI/API
- Symptom: ngent now caches prior provider transcript snapshots in SQLite for session-scoped `GET /v1/threads/{threadId}/history?sessionId=...` replies, but that replay is still not imported into SQLite `turns/events`; history APIs remain source-of-truth only for hub-created turns.
- Workaround: use the grouped session-list replay for provider-owned historical context, but rely on persisted hub `/history` for turns created through ngent itself.
- Follow-up plan: evaluate importing selected provider transcript into local persisted history, or exposing an explicit merged-history view, without duplicating future hub-originated turns.

- ID: KI-022
- Title: Codex grouped session-list titles can still show provider wrapper text
- Status: Open
- Severity: Low
- Affects: Codex `session/list` entries rendered in the Web UI grouped left rail
- Symptom: Codex provider metadata and replayed transcript can expose long wrapper-generated text such as `[Conversation Summary] ... [Current User Input] ...` or IDE context blocks because ngent now shows the raw provider-owned ACP replay.
- Workaround: use the thread title or the most recent visible turn content when the grouped-rail session label or replayed prompt body is noisy.
- Follow-up plan: normalize Codex `session/list` display titles in the backend, likely by preferring the first replayable user prompt over raw provider preview text when available.

- ID: KI-023
- Title: Fresh Kimi ACP sessions may resume before they appear in Kimi session browsing surfaces
- Status: Open
- Severity: Medium
- Affects: newly created Kimi sessions bound through ngent ACP turns
- Symptom:
  - a just-created Kimi `sessionId` can be resumed successfully through ACP `session/load`, but may still be absent from Kimi's own `session/list`, `kimi export`, and local `~/.kimi/sessions/*/<sessionId>` files for a while.
  - ngent can continue the bound session on the same or another thread if the `sessionId` is already known, but the grouped left rail may not show the new session immediately after creation.
- Workaround:
  - continue using the bound `sessionId` directly in ngent even if Kimi's own session browser has not caught up yet.
  - retry session browsing later after Kimi finishes persisting its own session index/files.
- Follow-up plan:
  - keep validating Kimi CLI session persistence timing across upstream releases.
  - decide whether ngent should annotate newly bound sessions as "not yet discoverable upstream" when local evidence supports that distinction.

- ID: KI-024
- Title: Kimi CLI 1.20.0 does not replay transcript messages during historical session/load
- Status: Open
- Severity: Medium
- Affects: Kimi historical session replay through session-scoped `GET /v1/threads/{threadId}/history?sessionId=...`
- Symptom:
  - Kimi `session/list` returns historical sessions and ACP `session/load` succeeds for those session ids.
  - Kimi CLI 1.20.0 currently emits no replay `session/update` notifications for those historical loads, so ngent returns `supported=true` with an empty transcript under the ACP-only implementation.
- Workaround:
  - continue the selected Kimi session normally; `session/load` still restores provider context for subsequent turns.
  - use ngent-local `/history` for turns created through ngent itself.
- Follow-up plan:
  - keep validating newer Kimi CLI releases and switch to transcript replay immediately if Kimi starts emitting standard `session/update` history during `session/load`.

- ID: KI-025
- Title: Session transcript replay cache does not auto-refresh from provider metadata
- Status: Open
- Severity: Medium
- Affects: repeated session-scoped `GET /v1/threads/{threadId}/history?sessionId=...` requests for the same `(agent, cwd, sessionId)`
- Symptom:
  - after the first successful replay, ngent serves the cached SQLite snapshot on later requests and server restarts.
  - if the provider session later gains more messages outside that cached snapshot, ngent does not yet compare provider `updatedAt` metadata before returning the cached transcript.
- Workaround:
  - continue the conversation through the current ngent thread so new hub-local turns remain visible in `/history`.
  - if a full provider replay refresh is required immediately, clear the cached row from sqlite and request session-scoped `/history` again.
- Follow-up plan:
  - persist `session/list.updatedAt` metadata and invalidate or refresh `session_transcript_cache` when that metadata advances.

- ID: KI-026
- Title: Thread-wide config changes still serialize across concurrent sessions
- Status: Open
- Severity: Low
- Affects: `PATCH /v1/threads/{threadId}` non-session updates and `POST /v1/threads/{threadId}/config-options`
- Symptom:
  - the server now allows concurrent turns across different sessions on the same thread.
  - shared thread metadata such as title, model, and config overrides still returns `409 CONFLICT` while any session on that thread is active.
- Workaround:
  - wait for active sessions to finish or cancel them before renaming the thread or changing shared model/config options.
- Follow-up plan:
  - evaluate whether some metadata-only updates can move to a narrower guard without letting shared provider state drift across sessions.

- ID: KI-027
- Title: Slash-command cache currently assumes one stable command set per agent
- Status: Open
- Severity: Low
- Affects: agents whose ACP `available_commands_update` payload may vary by workspace, session, or model
- Symptom:
  - ngent currently stores the latest slash-command snapshot in SQLite keyed only by `agent_id`.
  - if a provider later emits different slash-command sets for different contexts, the most recently observed snapshot for that agent will replace the earlier one and the Web UI composer may show commands from the wrong context.
- Workaround:
  - for Codex, simply opening the thread is now enough because `config-options` backfills the latest provider snapshot into sqlite before the first turn.
  - for Kimi, Qwen, OpenCode, and Gemini, the first `/slash-commands` request now backfills sqlite directly from the live provider snapshot cache when the provider has already emitted one, so typing `/` on a fresh thread is enough as long as the underlying CLI actually publishes `available_commands_update`.
  - for other agents, run a fresh turn in the target context before relying on slash-command suggestions, so the latest provider snapshot overwrites the cache.
- Follow-up plan:
  - keep provider-specific delivery timing fixes in place; codex caches the initial `session/new` / `session/load` snapshot before the first prompt, and Kimi/Qwen/OpenCode/Gemini now share both the same early ACP notification handling and the same provider-local slash-command cache across turn/config-session flows, so the remaining risk here is cache-key scope or provider emission behavior, not notification loss inside ngent.
  - revisit the cache key and move to `(agent_id, cwd)` or another provider-specific scope if a real agent starts varying slash commands by context.

- ID: KI-028
- Title: Real OpenCode ACP `session/new` can stall until hub timeout
- Status: Open
- Severity: Medium
- Affects: real `opencode` ACP turns and the host smoke test path
- Symptom:
  - on the current host validation run dated 2026-03-13, `E2E_OPENCODE=1 go test ./internal/agents/opencode -run TestOpenCodeE2ESmoke -count=1 -v -timeout 180s` failed with `opencode: session/new: context deadline exceeded`.
  - the process starts and initializes, but the upstream CLI does not complete `session/new` before the 45-second test context expires.
- Workaround:
  - verify local OpenCode auth/session readiness and backend reachability outside ngent, then retry.
  - if turns are business-critical, prefer a provider whose local CLI is already known-good on the host until OpenCode readiness is restored.
- Follow-up plan:
  - add richer diagnostics around stalled OpenCode `session/new` calls so auth/backend/readiness failures are distinguishable from protocol regressions.
  - keep rerunning the real smoke after local OpenCode environment fixes to confirm the shared ACP driver is not the blocking factor.

- ID: KI-029
- Title: Denied ACP permission turns currently collapse into an empty agent bubble
- Status: Open
- Severity: Low
- Affects: Web UI turns where the user denies a provider permission request
- Symptom:
  - the Web UI now correctly renders the pending `Permission Required` card for direct ACP providers such as Kimi and OpenCode, but once the user clicks `Deny`, the subsequent completed turn still has empty `responseText`.
  - after the final re-render, the ephemeral permission card disappears and the chat shows the existing empty-agent fallback (`…`) instead of a clearer provider rejection message.
- Workaround:
  - use the permission card itself as the source of truth for what was denied; the underlying tool action remains fail-closed and is not executed.
  - inspect `/history?includeEvents=1` if you need to confirm that the turn did emit `permission_required` before completion.
- Follow-up plan:
  - decide whether denied-permission turns should persist a lightweight terminal message, or whether the Web UI should keep the resolved permission card visible after turn completion.

- ID: KI-030
- Title: Provider-owned historical session replay still omits hidden reasoning, tool timeline, and rich content blocks
- Status: Open
- Severity: Low
- Affects: session-scoped `GET /v1/threads/{threadId}/history?sessionId=...` and Web UI session-sidebar replay for pre-existing provider sessions
- Symptom:
  - ngent now surfaces hidden reasoning and ordered tool/content/thought segments for hub-created turns by persisting turn events in normal history.
  - provider-owned historical replay returned in `/history.sessionTranscript` still exposes only visible `user` / `assistant` transcript messages, so switching to an older external session in the Web UI does not reconstruct past hidden reasoning blocks, the tool-call timeline, or non-text assistant content such as images/embedded resources.
- Workaround:
  - use regular ngent turn history for turns created through ngent itself; those now preserve reasoning and the ordered assistant segment timeline after reload.
  - treat provider-owned session replay as visible transcript-only until the replay contract is extended.
- Follow-up plan:
  - evaluate whether session-transcript schema should grow optional reasoning/tool metadata, or whether provider-owned replay should intentionally remain visible-transcript-only for privacy/product reasons.

- ID: KI-031
- Title: Web UI thinking expand state is not persisted across page reload
- Status: Open
- Severity: Low
- Affects: finalized agent messages that include reasoning/thinking content
- Symptom: users can expand a collapsed `Thinking` panel during the current page session, but a full page reload or fresh history load resets it to the default collapsed presentation.
- Workaround: expand the needed `Thinking` panel again after reload.
- Follow-up plan: evaluate persisting per-message UI presentation preferences in browser-local state if users need sticky behavior across reloads.

- ID: KI-032
- Title: Missing target-model catalog can temporarily leave stale reasoning choices after a picker change
- Status: Open
- Severity: Low
- Affects: threads that switch to a model whose sqlite config catalog has not been refreshed yet
- Symptom:
  - `POST /v1/threads/{threadId}/config-options` now persists the selected model immediately without mutating the live provider.
  - if sqlite does not yet have a catalog snapshot for the newly selected model, the immediate response can only fall back to the current in-memory option set, so the Web UI may temporarily show the previous reasoning list until a later real session lifecycle teaches ngent the new model snapshot.
- Workaround:
  - send the next turn on that model, or switch back after a real `session/new` / `session/load` has populated sqlite for the target model.
- Follow-up plan:
  - add an explicit on-demand fetch path for missing target-model catalogs so the picker can self-heal without waiting for the next turn.

- ID: KI-060
- Title: Some ACP tool-call payload shapes render as generic JSON
- Status: Open
- Severity: Low
- Affects: Web UI display of `tool_call` / `tool_call_update` events carrying non-text or provider-specific content blocks
- Symptom:
  - common text, diff, command, and path/location payloads render as structured cards.
  - richer ACP payloads such as media/resource-specific blocks currently fall back to raw JSON sections in the same tool-call card.
- Workaround:
  - inspect the JSON block shown in the tool-call card; the full structured payload is still preserved in SSE/history.
- Follow-up plan:
  - add richer renderers for additional ACP content block variants once real provider payloads stabilize.

- ID: KI-061
- Title: Full-suite `go test ./...` can flake when ACP-heavy agent packages run in parallel
- Status: Open
- Severity: Low
- Affects: local validation runs that execute the whole Go test suite with default package parallelism
- Symptom:
  - on some hosts, `go test ./...` can intermittently hit startup timeouts in ACP-heavy fake-process packages such as `internal/agents/codex`, `internal/agents/opencode`, or `internal/agents/qwen`.
  - the same packages pass when re-run individually, which points to host resource contention rather than a deterministic functional regression in those packages.
- Workaround:
  - rerun the full suite with serialized package scheduling, for example `GOFLAGS=-p=1 go test ./...`.
  - if one package still fails, rerun that package directly to distinguish a transient host-timeout from a real regression.
- Follow-up plan:
  - reduce test startup latency and remove package-parallel sensitivity in the ACP-heavy fake-process suites so default `go test ./...` stays stable without extra flags.

- ID: KI-036
- Title: Fresh agents do not expose model/reasoning metadata until one real session reports it
- Status: Open
- Severity: Low
- Affects: brand-new threads or fresh installs before any real turn / resumed-session turn has returned `configOptions`
- Symptom:
  - `GET /v1/agents/{agentId}/models` can return an empty list for an agent that has never reported config metadata into sqlite.
  - `GET /v1/threads/{threadId}/config-options` returns an empty `configOptions` list for a fresh thread until a real `session/new` or `session/load` runs during a turn.
  - the Web UI therefore hides model/reasoning controls until after that first real session snapshot arrives.
  - once a specific session has already been learned, switching back to that same session now restores its cached model/reasoning snapshot immediately.
  - switching directly onto an unseen existing session can also reveal controls immediately if that user-triggered `session/load` returns config metadata.
  - only sessions whose real `session/new` / `session/load` have never yielded config metadata remain empty.
- Workaround:
  - send one real turn on the thread, or switch to an existing session and let ngent load it once so the session's config snapshot can be learned.
- Follow-up plan:
  - if upstream CLIs eventually expose model/config catalogs without creating sessions, consider adopting those non-session surfaces so ngent can prefill metadata without reintroducing empty probe sessions.

- ID: KI-037
- Title: BLACKBOX ACP currently lacks session resume and catalog discovery surfaces
- Status: Open
- Severity: Medium
- Affects: `blackbox` threads, session-scoped `/history`, grouped-rail session browsing, and model picker/catalog endpoints for BLACKBOX
- Symptom:
  - local probing on 2026-03-22 against `blackbox 1.2.47` showed `initialize` advertising `agentCapabilities.loadSession=false`.
  - real `session/load` currently returns `-32601 method not found`.
  - `session/new` currently returns only `sessionId`, without `models.availableModels` or `configOptions`.
  - BLACKBOX can also emit stdout noise unrelated to ACP frames; ngent now tolerates that noise during transport, but the missing ACP resume/catalog surfaces still limit UX.
- Workaround:
  - use BLACKBOX through normal ngent thread turns; ngent-local `/history` still preserves turns created through ngent itself.
  - if a specific model must be forced, set `agentOptions.modelId` directly via API/advanced thread options even though no picker/catalog is currently available.
- Follow-up plan:
  - keep validating newer BLACKBOX CLI releases for `session/list` / `session/load` / model-catalog support.
  - wire those surfaces into ngent immediately once upstream ACP exposes them consistently.

- ID: KI-038
- Title: Cursor ACP runtime depends on pre-authenticated local CLI state
- Status: Open
- Severity: Medium
- Affects: implemented `cursor` provider turns in environments without a ready local Cursor login
- Symptom:
  - Cursor ACP `initialize` advertises `authMethods=[cursor_login]`.
  - if local Cursor login state is missing or unusable, the provider can fail during the explicit ACP `authenticate` step before `session/new`.
  - even after ACP authentication succeeds, prompt execution still depends on the local account's current entitlement/quota state.
- Workaround:
  - run `agent login` before issuing Cursor turns, or provide supported Cursor CLI auth flags/environment when launching the CLI outside ngent.
  - if turns authenticate but return product-gating text rather than the requested answer, verify the local Cursor account plan/quota state separately from ngent.
- Follow-up plan:
  - add richer preflight diagnostics for Cursor auth readiness beyond PATH existence if this becomes a recurring support issue.
  - keep validating whether future Cursor CLI releases expose a more explicit non-interactive auth-health probe.

- ID: KI-039
- Title: Persisted Web UI uploads currently have no automatic janitor
- Status: Open
- Severity: Low
- Affects: long-running ngent instances that handle many uploaded Web UI attachments
- Symptom:
  - `POST /v1/threads/{threadId}/turns` now persists uploaded files under the configured `data-path/attachments/*` tree so ACP providers can read them through stable `file://` resource links and the Web UI can keep rendering them after reload.
  - ngent currently keeps those persisted files indefinitely and does not yet run an age-based or reference-count-based cleanup sweep.
- Workaround:
  - periodically clear stale files from the configured `data-path/attachments/` tree if disk usage matters in a long-lived environment.
- Follow-up plan:
  - add an age-based or thread-reference-aware attachment janitor once real usage clarifies safe retention expectations for provider retries and history-driven debugging.

## Closed Issues

Closed issues are kept in one place to avoid split bookkeeping across the file.
Newer closures should appear first when practical.

- ID: KI-058
- Title: `/git-diff-file` duplicated raw text and per-line JSON, inflating preview payload size
- Status: Closed
- Severity: Low
- Affects: thread-scoped git-diff file preview responses consumed by the embedded Web UI and any future mobile client
- Symptom:
  - the endpoint previously returned both the full raw `content` string and a fully expanded `lines[]` array for the same preview.
  - large patches therefore paid for the same payload twice and also repeated `showLineNumbers` on every rendered row even though tone/line-number presence already implied that presentation.
- Workaround:
  - none required after the 2026-04-04 response compaction.
- Follow-up plan:
  - keep future clients on the compact grouped `blocks[]` response shape unless a separate raw-patch export use case appears that genuinely needs the unrendered text.

- ID: KI-057
- Title: Git-diff drawer used a synthetic single line counter instead of real old/new diff line numbers
- Status: Closed
- Severity: Medium
- Affects: embedded Web UI users reading tracked git patches in the right-side drawer
- Symptom:
  - the drawer previously rendered one monotonically increasing counter for every visible row, regardless of whether the row came from the old file, the new file, or a hunk header.
  - raw diff header metadata such as `diff --git`, `index`, `---`, and `+++` also remained visible even though they were not part of the requested code-reading surface.
- Workaround:
  - none required after the 2026-04-04 unified-diff line-number rendering fix.
- Follow-up plan:
  - none.

- ID: KI-056
- Title: Git-diff drawer rows used excessive vertical spacing
- Status: Closed
- Severity: Low
- Affects: embedded Web UI users reading multi-line diff/file content in the right-side drawer
- Symptom:
  - each rendered drawer row used relatively tall line-height and vertical padding, so long diffs consumed too much vertical space.
  - the requested fix was to make rows denser without reducing the visible text size.
- Workaround:
  - none required after the 2026-04-04 row-density polish.
- Follow-up plan:
  - none.

- ID: KI-055
- Title: Git-diff drawer code nodes still fell back to the browser default code font
- Status: Closed
- Severity: Low
- Affects: embedded Web UI users comparing the drawer content typography to the composer footer's model-label typography
- Symptom:
  - even after the drawer row container switched to the model-label typography, the actual rendered content still sat inside `<code>` nodes.
  - browser default `code` styling could therefore keep showing a different font family until the drawer content explicitly inherited the parent font settings.
- Workaround:
  - none required after the 2026-04-04 explicit font-inheritance fix.
- Follow-up plan:
  - none.

- ID: KI-054
- Title: Git-diff drawer content used heavy row separators and mismatched typography
- Status: Closed
- Severity: Low
- Affects: embedded Web UI users reading the right-side git-diff file preview drawer
- Symptom:
  - the drawer content used horizontal separators between every rendered line, which made the preview feel visually noisy.
  - the drawer text styling also differed from the composer footer's model-label typography requested for this surface.
- Workaround:
  - none required after the 2026-04-04 drawer typography polish.
- Follow-up plan:
  - none.

- ID: KI-053
- Title: Git-diff drawer reused the previous file-detail payload when reopening the same file
- Status: Closed
- Severity: Low
- Affects: embedded Web UI users who open file A, switch to another file, and then return to file A in the git-diff drawer
- Symptom:
  - reopening the same file could reuse the last fetched drawer detail from browser memory instead of issuing a fresh `/git-diff-file` request.
  - this made the drawer look stale when the user expected a new backend fetch on every file selection.
- Workaround:
  - none required after the 2026-04-04 follow-up fix.
- Follow-up plan:
  - none.

- ID: KI-052
- Title: Git-diff drawer used to auto-close on outside clicks and flicker under summary polling
- Status: Closed
- Severity: Medium
- Affects: embedded Web UI users reading a per-file git-diff preview in the right-side drawer
- Symptom:
  - after opening a changed file from the git-diff chip, clicking elsewhere in the workspace could close the drawer even though the user had not clicked its close affordance.
  - collapsing the git-diff chip itself could also dismiss the right-side drawer even though the user only intended to hide the changed-file list.
  - while the drawer remained open, periodic `/git-diff` summary refreshes could rebuild the drawer and force a detail refresh, causing visible loading flashes and scroll/content jumping.
- Workaround:
  - none required after the 2026-04-04 Web UI stability fix.
- Follow-up plan:
  - none.

- ID: KI-042
- Title: Web UI session browsing raised 409 conflicts while another session was still streaming
- Status: Closed
- Severity: Medium
- Affects: switching between sessions in the Web UI while the thread already has an active turn
- Symptom:
  - choosing another session from the session browser previously always issued `PATCH /v1/threads/{threadId}` immediately.
  - if any turn in that thread was still active, the server returned `409 thread has an active turn`, so even read-only session browsing produced an error dialog.
- Workaround:
  - none; fixed on 2026-03-26 by separating local viewed-session state from backend thread session binding and deferring backend sync until the thread is idle.
- Follow-up plan:
  - monitor whether users also need an explicit visible indicator when they are browsing a session that has not yet been synced back into backend thread state.

- ID: KI-041
- Title: Web UI session switch fetched whole-thread history and stalled on large multi-session threads
- Status: Closed
- Severity: Medium
- Affects: Web UI session switching on threads with many persisted turns/events across multiple sessions
- Symptom:
  - selecting one historical session previously still triggered `GET /v1/threads/{threadId}/history?includeEvents=1` for the entire thread, then filtered by `session_bound` in the browser.
  - on real Codex threads this could force the UI to parse roughly 19 MB / 42k events just to render one session that only needed one persisted turn.
- Workaround:
  - none; fixed on 2026-03-26 by adding `sessionId` filtering to `/history`, compacting historical delta runs on read, and yielding during heavy message-list replay in the Web UI.
- Follow-up plan:
  - monitor whether any remaining session-switch lag is dominated by `tool_call_update` volume or provider transcript replay load time rather than persisted thread history size.

- ID: KI-040
- Title: Inline base64 image placeholders in user messages rendered as raw text
- Status: Closed
- Severity: Low
- Affects: Web UI threads whose user messages include bracketed placeholders such as `[Image: data:image/png;base64,...]`
- Symptom:
  - the user bubble previously sent the entire message through plain markdown rendering, so bracketed base64 image placeholders showed up as long unreadable text blobs instead of visible image previews.
- Workaround:
  - none; fixed on 2026-03-26.
- Follow-up plan:
  - monitor whether any upstream source emits materially different placeholder shapes that should also be normalized into the same inline-image render path.

- ID: KI-033
- Title: Repeated `New session` reused stale anonymous Web UI scope after fast cancel
- Status: Closed
- Severity: Low
- Affects: Web UI fresh-session flows cancelled before ACP emits `session_bound`
- Symptom:
  - a fast `send -> cancel -> New session` sequence could leave the chat pane showing the just-cancelled content instead of returning to an empty composer because the UI kept reusing the same empty-session scope.
  - reopening the thread after reload could also replay that empty cancelled placeholder from local history hydration.
- Workaround:
  - none; fixed on 2026-03-16.
- Follow-up plan:
  - monitor whether explicit fresh-session scope state ever needs backend persistence beyond the current Web UI behavior.

- ID: KI-004
- Title: `cwd` validation false positives
- Status: Closed
- Severity: Low
- Affects: legacy deployments that used restrictive allow-root policies
- Symptom: historical issue where valid paths could be rejected as outside allowed roots
- Workaround: N/A after ADR-016 default absolute-cwd policy
- Follow-up plan: none
