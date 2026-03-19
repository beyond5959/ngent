import { api } from '../api.ts'
import { escHtml } from '../utils.ts'
import type { PermissionRequiredPayload } from '../sse.ts'

// Server-side default permissionTimeout is 2 hours
export const PERMISSION_TIMEOUT_MS = 2 * 60 * 60 * 1000
const TICK_MS = 1_000
type ResolveOutcome = 'approved' | 'declined' | 'timeout'

interface MountOptions {
  deadlineMs?: number
  onResolved?: (outcome: ResolveOutcome) => void
}

// ── Public entry point ────────────────────────────────────────────────────

/**
 * Appends a permission card to `listEl` and starts its countdown timer.
 * The card manages its own lifecycle; no cleanup from the caller is needed.
 */
export function mountPermissionCard(
  listEl: HTMLElement,
  event: PermissionRequiredPayload,
  options?: MountOptions,
): void {
  const timeoutMs = options?.deadlineMs
    ? Math.max(0, options.deadlineMs - Date.now())
    : PERMISSION_TIMEOUT_MS

  const wrapper = document.createElement('div')
  wrapper.className = 'message message--agent'
  wrapper.innerHTML = buildHtml(event, timeoutMs)
  listEl.appendChild(wrapper)
  listEl.scrollTop = listEl.scrollHeight

  // Elements are now in the DOM — bind interactivity
  bindCard(event.permissionId, timeoutMs, options?.onResolved)
}

// ── HTML template ─────────────────────────────────────────────────────────

function buildHtml(event: PermissionRequiredPayload, timeoutMs: number): string {
  const { permissionId: pid, approval, command } = event

  return `
    <div class="message-avatar perm-avatar">!</div>
    <div class="message-group" style="max-width:min(480px,90%)">
      <div class="permission-card" id="perm-card-${pid}">

        <div class="permission-card-header">
          <span class="permission-badge permission-badge--${escHtml(approval)}">${escHtml(approval)}</span>
          <span class="permission-card-title">Permission Required</span>
        </div>

        <div class="permission-card-body">
          <code class="permission-command">${escHtml(command)}</code>
        </div>

        <div class="permission-card-footer">
          <span class="permission-countdown" id="perm-cd-${pid}">${formatRemaining(timeoutMs)}</span>
          <div class="permission-actions" id="perm-actions-${pid}">
            <button class="btn btn-sm btn-success" id="perm-allow-${pid}">Allow</button>
            <button class="btn btn-sm btn-danger"  id="perm-deny-${pid}">Deny</button>
          </div>
        </div>

      </div>
    </div>`
}

// ── Countdown + button binding ────────────────────────────────────────────

function bindCard(
  pid: string,
  timeoutMs: number,
  onResolved?: (outcome: ResolveOutcome) => void,
): void {
  let resolved = false
  let elapsed  = 0
  if (timeoutMs <= 0) {
    showResolved(pid, 'timeout', onResolved)
    return
  }

  const tick = setInterval(() => {
    const cardEl = document.getElementById(`perm-card-${pid}`)
    if (!cardEl) {
      clearInterval(tick)
      return
    }

    elapsed += TICK_MS
    const remaining = Math.max(0, timeoutMs - elapsed)

    const cdEl = document.getElementById(`perm-cd-${pid}`)
    if (cdEl) cdEl.textContent = formatRemaining(remaining)

    if (remaining === 0 && !resolved) {
      resolved = true
      clearInterval(tick)
      showResolved(pid, 'timeout', onResolved)
    }
  }, TICK_MS)

  document.getElementById(`perm-allow-${pid}`)?.addEventListener('click', () => {
    if (resolved) return
    resolved = true
    clearInterval(tick)
    void doResolve(pid, 'approved', onResolved)
  })

  document.getElementById(`perm-deny-${pid}`)?.addEventListener('click', () => {
    if (resolved) return
    resolved = true
    clearInterval(tick)
    void doResolve(pid, 'declined', onResolved)
  })
}

function formatRemaining(remainingMs: number): string {
  const totalSeconds = Math.ceil(Math.max(0, remainingMs) / 1000)
  const hours = Math.floor(totalSeconds / 3600)
  const minutes = Math.floor((totalSeconds % 3600) / 60)
  const seconds = totalSeconds % 60
  return `${hours}:${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
}

// ── API call ──────────────────────────────────────────────────────────────

async function doResolve(
  pid: string,
  outcome: 'approved' | 'declined',
  onResolved?: (outcome: ResolveOutcome) => void,
): Promise<void> {
  // Disable buttons immediately so the user can't double-click
  document
    .getElementById(`perm-actions-${pid}`)
    ?.querySelectorAll<HTMLButtonElement>('button')
    .forEach(b => { b.disabled = true })

  try {
    await api.resolvePermission(pid, outcome)
  } catch {
    // 409 = already resolved by server — show the intended outcome anyway
  }

  showResolved(pid, outcome, onResolved)
}

// ── Resolved state ────────────────────────────────────────────────────────

function showResolved(
  pid: string,
  outcome: ResolveOutcome,
  onResolved?: (outcome: ResolveOutcome) => void,
): void {
  const cardEl    = document.getElementById(`perm-card-${pid}`)
  const actionsEl = document.getElementById(`perm-actions-${pid}`)
  const cdEl      = document.getElementById(`perm-cd-${pid}`)
  if (!cardEl) {
    onResolved?.(outcome)
    return
  }

  // Hide countdown text
  if (cdEl) cdEl.textContent = ''

  // Replace Allow/Deny buttons with a resolved label
  if (actionsEl) {
    const label =
      outcome === 'approved' ? '✓ Approved'             :
      outcome === 'declined' ? '✗ Denied'               :
                               '⏱ Timed out (auto-denied)'
    actionsEl.innerHTML = `
      <span class="permission-resolved permission-resolved--${outcome}">${label}</span>`
  }

  // Recolor card header and border to reflect outcome
  cardEl.classList.add(`permission-card--${outcome}`)
  onResolved?.(outcome)
}
