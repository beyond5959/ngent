import { getLanguageLabel, setLanguage, t } from '../i18n.ts'
import { store } from '../store.ts'
import type { Language, Theme } from '../types.ts'
import { debounce, escHtml } from '../utils.ts'

// ── Icons ──────────────────────────────────────────────────────────────────

const iconClose = `<svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden="true">
  <path d="M3 3l10 10M13 3L3 13" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
</svg>`

// ── Render ─────────────────────────────────────────────────────────────────

interface SettingsDraftValues {
  authToken: string
  serverUrl: string
}

function readDraftValues(): SettingsDraftValues {
  const state = store.get()
  return {
    authToken: container?.querySelector<HTMLInputElement>('#auth-token-input')?.value ?? state.authToken,
    serverUrl: container?.querySelector<HTMLInputElement>('#server-url-input')?.value ?? state.serverUrl,
  }
}

function renderPanel(drafts: SettingsDraftValues = readDraftValues()): string {
  const { theme, language } = store.get()
  const { authToken, serverUrl } = drafts

  const themeBtn = (value: Theme, label: string) => `
    <button
      class="theme-btn ${theme === value ? 'theme-btn--active' : ''}"
      data-theme-value="${value}"
      type="button"
    >${label}</button>`

  const languageBtn = (value: Language, label: string) => `
    <button
      class="theme-btn ${language === value ? 'theme-btn--active' : ''}"
      data-language-value="${value}"
      type="button"
    >${label}</button>`

  return `
    <div class="settings-overlay" id="settings-overlay" role="dialog" aria-modal="true" aria-label="${escHtml(t('settingsDialogLabel'))}">
      <div class="settings-panel" id="settings-panel">

        <div class="settings-header">
          <div class="settings-header-copy">
            <div class="settings-kicker">${escHtml(t('browserPreferences'))}</div>
            <h2 class="settings-title">${escHtml(t('settingsTitle'))}</h2>
            <p class="settings-header-desc">${escHtml(t('settingsHeaderDesc'))}</p>
          </div>
          <button class="btn btn-icon" id="settings-close-btn" aria-label="${escHtml(t('closeSettings'))}">
            ${iconClose}
          </button>
        </div>

        <div class="settings-body">
          <div class="settings-intro">
            <div class="settings-intro-badge">${escHtml(t('storedLocally'))}</div>
            <p class="settings-intro-copy">${escHtml(t('settingsIntroCopy'))}</p>
          </div>

          <section class="settings-section">
            <h3 class="settings-section-title">${escHtml(t('connection'))}</h3>
            <label class="settings-label" for="server-url-input">${escHtml(t('serverUrl'))}</label>
            <p class="settings-description">
              ${escHtml(t('serverUrlDesc'))}
            </p>
            <input
              id="server-url-input"
              class="settings-input"
              type="url"
              placeholder="http://127.0.0.1:8686"
              value="${escHtml(serverUrl)}"
              autocomplete="off"
              spellcheck="false"
            />
          </section>

          <section class="settings-section">
            <h3 class="settings-section-title">${escHtml(t('security'))}</h3>
            <label class="settings-label" for="auth-token-input">${escHtml(t('bearerToken'))}</label>
            <p class="settings-description">
              ${escHtml(t('bearerTokenDesc'))}
            </p>
            <input
              id="auth-token-input"
              class="settings-input"
              type="password"
              placeholder="${escHtml(t('bearerTokenPlaceholder'))}"
              value="${escHtml(authToken)}"
              autocomplete="off"
              spellcheck="false"
            />
          </section>

          <section class="settings-section">
            <h3 class="settings-section-title">${escHtml(t('appearance'))}</h3>
            <label class="settings-label">${escHtml(t('theme'))}</label>
            <div class="theme-btn-group">
              ${themeBtn('light', t('light'))}
              ${themeBtn('system', t('system'))}
              ${themeBtn('dark', t('dark'))}
            </div>
          </section>

          <section class="settings-section">
            <h3 class="settings-section-title">${escHtml(t('language'))}</h3>
            <label class="settings-label">${escHtml(t('language'))}</label>
            <p class="settings-description">${escHtml(t('languageDesc'))}</p>
            <div class="theme-btn-group">
              ${languageBtn('en', getLanguageLabel('en'))}
              ${languageBtn('zh-CN', getLanguageLabel('zh-CN'))}
              ${languageBtn('es', getLanguageLabel('es'))}
              ${languageBtn('fr', getLanguageLabel('fr'))}
            </div>
          </section>

        </div>
      </div>
    </div>`
}

// ── Mount / Unmount ────────────────────────────────────────────────────────

let container: HTMLDivElement | null = null

function rerender(): void {
  if (!container) return
  const drafts = readDraftValues()
  container.innerHTML = renderPanel(drafts)
  bindEvents()
}

function unmount(): void {
  if (container) {
    container.remove()
    container = null
  }
  store.set({ settingsOpen: false })
}

function mount(): void {
  if (container) return

  container = document.createElement('div')
  container.innerHTML = renderPanel(readDraftValues())
  document.body.appendChild(container)

  bindEvents()

  // Focus the panel for keyboard navigation
  ;(container.querySelector('#settings-panel') as HTMLElement | null)?.focus()
}

// ── Event binding ──────────────────────────────────────────────────────────

function bindEvents(): void {
  if (!container) return

  // Close on backdrop click
  container.querySelector<HTMLElement>('#settings-overlay')?.addEventListener('click', e => {
    if ((e.target as HTMLElement).id === 'settings-overlay') unmount()
  })

  // Close button
  container.querySelector('#settings-close-btn')?.addEventListener('click', unmount)

  // Escape key
  const onKey = (e: KeyboardEvent) => {
    if (e.key === 'Escape') { unmount(); document.removeEventListener('keydown', onKey) }
  }
  document.addEventListener('keydown', onKey)

  // Auth token — save on change (debounced)
  const saveToken = debounce((v: string) => store.set({ authToken: v }), 400)
  container.querySelector<HTMLInputElement>('#auth-token-input')?.addEventListener('input', e => {
    saveToken((e.target as HTMLInputElement).value)
  })

  // Theme buttons
  container.querySelector('.theme-btn-group')?.addEventListener('click', e => {
    const btn = (e.target as HTMLElement).closest<HTMLButtonElement>('[data-theme-value]')
    if (!btn) return
    const value = btn.dataset.themeValue as Theme
    store.set({ theme: value })
    applyTheme(value)
    rerender()
  })

  // Language buttons
  container.querySelectorAll('.theme-btn-group').forEach(group => {
    group.addEventListener('click', e => {
      const btn = (e.target as HTMLElement).closest<HTMLButtonElement>('[data-language-value]')
      if (!btn) return
      const value = btn.dataset.languageValue as Language
      store.set({ language: value })
      setLanguage(value)
      rerender()
    })
  })

  // Server URL — save on change (debounced)
  const saveUrl = debounce((v: string) => store.set({ serverUrl: v }), 400)
  container.querySelector<HTMLInputElement>('#server-url-input')?.addEventListener('input', e => {
    saveUrl((e.target as HTMLInputElement).value)
  })
}

// ── Theme application ──────────────────────────────────────────────────────

function getSystemTheme(): 'light' | 'dark' {
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

export function applyTheme(theme: Theme): void {
  document.documentElement.dataset.theme = theme === 'system' ? getSystemTheme() : theme
}

// ── Public API ─────────────────────────────────────────────────────────────

export const settingsPanel = {
  open(): void {
    store.set({ settingsOpen: true })
    mount()
  },
  close(): void {
    unmount()
  },
}
