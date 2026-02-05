/**
 * Datastar Inspector - Dev-only debugging tool
 *
 * A self-contained web component for inspecting Datastar signals.
 * Works in any Datastar project with zero configuration.
 *
 * Usage:
 *   <script type="module" src="datastar-inspector.js"></script>
 *   <datastar-inspector></datastar-inspector>
 *
 * Features:
 * - Real-time signal display (JSON or Table view)
 * - Filter signals by text, regex, or wildcards
 * - Flash animation on signal changes
 * - Persisted state (expanded, filter, viewMode)
 */
import { LitElement, html, nothing } from 'lit'
import { customElement, state } from 'lit/decorators.js'

import type { InspectorState, ViewMode, SignalObject } from './types.js'
import { STYLES_ID, STORAGE_KEY, FLASH_DURATION, INSPECTOR_STYLES } from './styles.js'
import {
  countSignals,
  flattenSignals,
  filterObject,
  findChangedPaths,
  parseFilterPattern,
  renderJsonValue,
} from './utils.js'

/**
 * Datastar Inspector web component.
 *
 * Displays current Datastar signals in a floating panel with:
 * - Toggle between collapsed (DS button) and expanded (panel) states
 * - JSON view with syntax highlighting
 * - Table view with flattened signal paths
 * - Filter by text, /regex/, or wildcards (*)
 * - Flash animation when signals change
 *
 * @element datastar-inspector
 */
@customElement('datastar-inspector')
export class DatastarInspector extends LitElement {
  // ============================================
  // Light DOM (required for Datastar compatibility)
  // ============================================

  override createRenderRoot() {
    return this
  }

  // ============================================
  // Reactive State
  // ============================================

  @state() private expanded = false
  @state() private filter = ''
  @state() private viewMode: ViewMode = 'json'
  @state() private signals: SignalObject = {}
  @state() private signalCount = 0
  @state() private changedPaths: Set<string> = new Set()
  @state() private hasUnseenChanges = false
  @state() private panelWidth = 340
  @state() private panelHeight = 280

  // ============================================
  // Private Properties
  // ============================================

  private observer: MutationObserver | null = null
  private signalsElementId = `ds-inspector-signals-${Math.random().toString(36).slice(2, 9)}`
  private previousSignals: SignalObject = {}
  private flashTimeout: number | null = null
  private isResizing = false
  private resizeStartX = 0
  private resizeStartY = 0
  private resizeStartWidth = 0
  private resizeStartHeight = 0

  // Size constraints
  private readonly minWidth = 280
  private readonly maxWidth = 600
  private readonly minHeight = 200
  private readonly maxHeight = 800

  // ============================================
  // Lifecycle
  // ============================================

  override connectedCallback() {
    super.connectedCallback()
    this.loadState()
    this.injectStyles()
  }

  override disconnectedCallback() {
    super.disconnectedCallback()
    this.observer?.disconnect()
    if (this.flashTimeout) {
      clearTimeout(this.flashTimeout)
    }
    // Clean up resize listeners
    document.removeEventListener('mousemove', this.handleResizeMove)
    document.removeEventListener('mouseup', this.handleResizeEnd)
  }

  override firstUpdated() {
    this.setupSignalObserver()
  }

  // ============================================
  // State Persistence
  // ============================================

  private loadState() {
    try {
      const saved = sessionStorage.getItem(STORAGE_KEY)
      if (saved) {
        const state: InspectorState = JSON.parse(saved)
        this.expanded = state.expanded ?? false
        this.filter = state.filter ?? ''
        this.viewMode = state.viewMode ?? 'json'
        this.panelWidth = state.panelWidth ?? 340
        this.panelHeight = state.panelHeight ?? 280
      }
    } catch {
      /* ignore parse errors */
    }
  }

  private saveState() {
    const state: InspectorState = {
      expanded: this.expanded,
      filter: this.filter,
      viewMode: this.viewMode,
      panelWidth: this.panelWidth,
      panelHeight: this.panelHeight,
    }
    sessionStorage.setItem(STORAGE_KEY, JSON.stringify(state))
  }

  // ============================================
  // Style Injection
  // ============================================

  private injectStyles() {
    if (!document.getElementById(STYLES_ID)) {
      const style = document.createElement('style')
      style.id = STYLES_ID
      style.textContent = INSPECTOR_STYLES
      document.head.appendChild(style)
    }
  }

  // ============================================
  // Signal Observation
  // ============================================

  private setupSignalObserver() {
    const el = document.getElementById(this.signalsElementId)
    if (!el) return

    // Initial read (don't flash on initial load)
    this.parseSignals(el.textContent || '{}', true)

    // Watch for changes
    this.observer = new MutationObserver(() => {
      this.parseSignals(el.textContent || '{}', false)
    })
    this.observer.observe(el, { childList: true, characterData: true, subtree: true })
  }

  private parseSignals(json: string, isInitial: boolean) {
    try {
      const newSignals = JSON.parse(json) as SignalObject

      // Detect changed paths (skip on initial load)
      if (!isInitial && Object.keys(this.previousSignals).length > 0) {
        const changed = findChangedPaths(this.previousSignals, newSignals)
        if (changed.size > 0) {
          this.changedPaths = changed

          // Track unseen changes when collapsed
          if (!this.expanded) {
            this.hasUnseenChanges = true
          }

          // Clear flash after animation duration
          if (this.flashTimeout) {
            clearTimeout(this.flashTimeout)
          }
          this.flashTimeout = window.setTimeout(() => {
            this.changedPaths = new Set()
            this.hasUnseenChanges = false
          }, FLASH_DURATION)
        }
      }

      this.previousSignals = JSON.parse(json) as SignalObject // Deep copy
      this.signals = newSignals
      this.signalCount = countSignals(this.signals)
    } catch {
      this.signals = {}
      this.signalCount = 0
    }
  }

  // ============================================
  // Filtering
  // ============================================

  private getFilteredSignals(): SignalObject {
    if (!this.filter.trim()) return this.signals

    const regex = parseFilterPattern(this.filter.trim())
    return filterObject(this.signals, regex) as SignalObject
  }

  // ============================================
  // Event Handlers
  // ============================================

  private toggle() {
    this.expanded = !this.expanded
    this.saveState()
    // Re-setup observer after expansion (element may have been created)
    if (this.expanded) {
      this.hasUnseenChanges = false // Clear unseen changes indicator when opening
      requestAnimationFrame(() => this.setupSignalObserver())
    }
  }

  private close() {
    this.expanded = false
    this.saveState()
  }

  private handleFilterInput(e: Event) {
    this.filter = (e.target as HTMLInputElement).value
    this.saveState()
  }

  private clearFilter() {
    this.filter = ''
    this.saveState()
  }

  private setViewMode(mode: ViewMode) {
    this.viewMode = mode
    this.saveState()
  }

  // ============================================
  // Resize Handlers
  // ============================================

  private handleResizeStart = (e: MouseEvent) => {
    e.preventDefault()
    this.isResizing = true
    this.resizeStartX = e.clientX
    this.resizeStartY = e.clientY
    this.resizeStartWidth = this.panelWidth
    this.resizeStartHeight = this.panelHeight

    document.addEventListener('mousemove', this.handleResizeMove)
    document.addEventListener('mouseup', this.handleResizeEnd)
  }

  private handleResizeMove = (e: MouseEvent) => {
    if (!this.isResizing) return

    // Since panel is anchored bottom-right and handle is top-left:
    // - Moving mouse left (negative deltaX) increases width
    // - Moving mouse up (negative deltaY) increases height
    const deltaX = this.resizeStartX - e.clientX
    const deltaY = this.resizeStartY - e.clientY

    const newWidth = Math.min(this.maxWidth, Math.max(this.minWidth, this.resizeStartWidth + deltaX))
    const newHeight = Math.min(this.maxHeight, Math.max(this.minHeight, this.resizeStartHeight + deltaY))

    this.panelWidth = newWidth
    this.panelHeight = newHeight
  }

  private handleResizeEnd = () => {
    this.isResizing = false
    document.removeEventListener('mousemove', this.handleResizeMove)
    document.removeEventListener('mouseup', this.handleResizeEnd)
    this.saveState()
  }

  // ============================================
  // Render
  // ============================================

  override render() {
    const filteredSignals = this.getFilteredSignals()
    const filteredCount = countSignals(filteredSignals)
    const hasFilter = this.filter.trim().length > 0

    return html`
      <!-- Hidden element for Datastar signal access -->
      <pre id="${this.signalsElementId}" class="ds-inspector-hidden" data-json-signals></pre>

      ${this.expanded ? this.renderPanel(filteredSignals, filteredCount, hasFilter) : this.renderToggle()}
    `
  }

  private renderToggle() {
    const toggleClass = this.hasUnseenChanges
      ? 'ds-inspector-toggle ds-inspector-toggle--changed'
      : 'ds-inspector-toggle'
    return html`
      <button class="${toggleClass}" @click=${this.toggle} title="Open Datastar Inspector">
        DS
      </button>
    `
  }

  private renderPanel(filteredSignals: SignalObject, filteredCount: number, hasFilter: boolean) {
    const panelStyle = `width: ${this.panelWidth}px; height: ${this.panelHeight}px;`
    return html`
      <div class="ds-inspector-panel" style="${panelStyle}">
        <div class="ds-inspector-resize-handle" @mousedown=${this.handleResizeStart}></div>
        ${this.renderHeader(filteredCount, hasFilter)}
        ${this.renderContent(filteredSignals, hasFilter)}
      </div>
    `
  }

  private renderHeader(filteredCount: number, hasFilter: boolean) {
    const placeholder = hasFilter
      ? `${filteredCount}/${this.signalCount} match...`
      : `Filter ${this.signalCount} signals...`
    return html`
      <div class="ds-inspector-header">
        <span class="ds-inspector-logo" title="Datastar Inspector">DS</span>
        <div class="ds-inspector-filter-wrapper">
          <input
            type="text"
            class="ds-inspector-filter"
            placeholder="${placeholder}"
            .value=${this.filter}
            @input=${this.handleFilterInput}
          />
          ${hasFilter
            ? html`<button class="ds-inspector-filter-clear" @click=${this.clearFilter}>&times;</button>`
            : nothing}
        </div>
        <div class="ds-inspector-view-toggle">
          <button
            class="ds-inspector-view-btn ${this.viewMode === 'json' ? 'active' : ''}"
            @click=${() => this.setViewMode('json')}
            title="JSON view"
          >
            { }
          </button>
          <button
            class="ds-inspector-view-btn ${this.viewMode === 'table' ? 'active' : ''}"
            @click=${() => this.setViewMode('table')}
            title="Table view"
          >
            ≡
          </button>
        </div>
        <button class="ds-inspector-btn" @click=${this.close} title="Close">&times;</button>
      </div>
    `
  }

  private renderContent(filteredSignals: SignalObject, hasFilter: boolean) {
    const isEmpty = Object.keys(filteredSignals).length === 0

    return html`
      <div class="ds-inspector-content">
        ${isEmpty
          ? html`<div class="ds-inspector-empty">
              ${hasFilter ? 'No signals match filter' : 'No signals found'}
            </div>`
          : this.viewMode === 'json'
            ? this.renderJsonView(filteredSignals)
            : this.renderTableView(filteredSignals)}
      </div>
    `
  }

  private renderJsonView(signals: SignalObject) {
    return html`
      <pre class="ds-inspector-json" .innerHTML=${renderJsonValue(signals, this.changedPaths)}></pre>
    `
  }

  private renderTableView(signals: SignalObject) {
    return html`
      <table class="ds-inspector-table">
        <thead>
          <tr>
            <th>Signal</th>
            <th>Value</th>
          </tr>
        </thead>
        <tbody>
          ${flattenSignals(signals).map(([path, value]) => {
            const isChanged = this.changedPaths.has(path)
            return html`
              <tr class="${isChanged ? 'ds-inspector-flash' : ''}">
                <td>${path}</td>
                <td class="ds-inspector-table-value" title=${JSON.stringify(value)}>
                  ${JSON.stringify(value)}
                </td>
              </tr>
            `
          })}
        </tbody>
      </table>
    `
  }
}

// ============================================
// TypeScript HTMLElementTagNameMap Extension
// ============================================

declare global {
  interface HTMLElementTagNameMap {
    'datastar-inspector': DatastarInspector
  }
}
