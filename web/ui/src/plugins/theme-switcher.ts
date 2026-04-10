/**
 * Theme Switcher Module
 * Handles light/dark/system theme switching with localStorage persistence
 */

export type ThemeMode = 'system' | 'light' | 'dark'
export type ResolvedTheme = 'light' | 'dark'

const STORAGE_KEY = 'theme'
const THEME_ATTR = 'data-theme'

/** Cycle order for theme toggle */
const CYCLE_ORDER: ThemeMode[] = ['system', 'light', 'dark']

/**
 * Get the resolved theme based on system preference
 */
function getSystemTheme(): ResolvedTheme {
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

function resolveTheme(mode: ThemeMode): ResolvedTheme {
  return mode === 'system' ? getSystemTheme() : mode
}

function emitThemeChange(mode: ThemeMode): void {
  window.dispatchEvent(
    new CustomEvent('themechange', {
      detail: {
        mode,
        resolved: resolveTheme(mode),
      },
    })
  )
}

/**
 * Apply the resolved theme to the document
 */
function applyTheme(mode: ThemeMode): void {
  const resolved = resolveTheme(mode)
  document.documentElement.setAttribute(THEME_ATTR, resolved)
  emitThemeChange(mode)
}

/**
 * Get the current theme preference from localStorage
 */
export function getTheme(): ThemeMode {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored === 'light' || stored === 'dark' || stored === 'system') {
    return stored
  }
  return 'system'
}

export function getResolvedTheme(): ResolvedTheme {
  return resolveTheme(getTheme())
}

/**
 * Set the theme and persist to localStorage
 */
export function setTheme(mode: ThemeMode): void {
  localStorage.setItem(STORAGE_KEY, mode)
  applyTheme(mode)
}

/**
 * Cycle to the next theme in order: system → light → dark → system
 * Returns the new theme mode
 */
export function cycleTheme(): ThemeMode {
  const current = getTheme()
  const currentIndex = CYCLE_ORDER.indexOf(current)
  const nextIndex = (currentIndex + 1) % CYCLE_ORDER.length
  const next = CYCLE_ORDER[nextIndex]
  setTheme(next)
  return next
}

/**
 * Initialize theme on page load
 * Should be called as early as possible (ideally in <head>) to prevent flash
 */
export function initTheme(): void {
  const mode = getTheme()
  applyTheme(mode)

  // Listen for system preference changes when in 'system' mode
  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  mediaQuery.addEventListener('change', () => {
    if (getTheme() === 'system') {
      applyTheme('system')
    }
  })
}

/**
 * Theme toggle API for use with Datastar
 */
export const themeToggle = {
  get: getTheme,
  resolved: getResolvedTheme,
  set: setTheme,
  cycle: cycleTheme,
  init: initTheme,
}

// Attach to window for Datastar access
declare global {
  interface Window {
    themeToggle: typeof themeToggle
  }
}

window.themeToggle = themeToggle

// Auto-initialize when module loads
initTheme()
