/**
 * Styles for Datastar Inspector
 *
 * Self-contained CSS that gets injected into document.head.
 * Uses scoped custom properties to ensure dark mode regardless of page theme.
 */

/** ID for the injected style element */
export const STYLES_ID = 'ds-inspector-styles'

/** Flash animation duration in milliseconds */
export const FLASH_DURATION = 400

/** Storage key for persisted state */
export const STORAGE_KEY = 'ds-inspector'

/**
 * Complete CSS styles for the inspector.
 * Injected into document.head on component mount.
 */
export const INSPECTOR_STYLES = `
/* ============================================
   Datastar Inspector - Self-contained styles
   Always dark mode - properties scoped to component
   ============================================ */

/* ----------------------------------------
   Animations
   ---------------------------------------- */

@keyframes ds-inspector-flash {
  0% { background-color: rgba(250, 204, 21, 0.5); }
  100% { background-color: transparent; }
}

@keyframes ds-inspector-toggle-flash {
  0% { 
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4), 0 0 0 0 rgba(250, 204, 21, 0.4);
    border-color: rgba(250, 204, 21, 0.6);
  }
  100% { 
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4), 0 0 0 6px rgba(250, 204, 21, 0);
    border-color: rgba(250, 204, 21, 0);
  }
}

.ds-inspector-flash {
  animation: ds-inspector-flash ${FLASH_DURATION}ms ease-out;
  border-radius: 2px;
}

/* ----------------------------------------
   Base Styles & Custom Properties
   ---------------------------------------- */

.ds-inspector-toggle,
.ds-inspector-panel {
  /* Force dark color scheme for native elements (scrollbars, inputs) */
  color-scheme: dark;

  /* Scoped custom properties - always dark */
  --ds-inspector-bg: #1e1e2e;
  --ds-inspector-surface: #313244;
  --ds-inspector-text: #cdd6f4;
  --ds-inspector-text-dim: #a6adc8;
  --ds-inspector-border: #45475a;
  --ds-inspector-accent: #89b4fa;
  --ds-inspector-key: #f38ba8;
  --ds-inspector-string: #a6e3a1;
  --ds-inspector-number: #fab387;
  --ds-inspector-boolean: #cba6f7;
  --ds-inspector-null: #6c7086;
  --ds-inspector-flash: #facc15;
  --ds-inspector-font: ui-monospace, 'SF Mono', 'Cascadia Code', 'Fira Code', Consolas, monospace;
  --ds-inspector-z-index: 99999;
}

/* ----------------------------------------
   Toggle Button (collapsed state)
   ---------------------------------------- */

.ds-inspector-toggle {
  position: fixed;
  bottom: 16px;
  right: 16px;
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: var(--ds-inspector-bg);
  border: 1px solid var(--ds-inspector-border);
  color: var(--ds-inspector-accent);
  font-family: var(--ds-inspector-font);
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: var(--ds-inspector-z-index);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.ds-inspector-toggle:hover {
  transform: scale(1.05);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.5);
}

.ds-inspector-toggle--changed {
  animation: ds-inspector-toggle-flash ${FLASH_DURATION}ms ease-out;
}

/* ----------------------------------------
   Panel (expanded state)
   ---------------------------------------- */

.ds-inspector-panel {
  position: fixed;
  bottom: 16px;
  right: 16px;
  background: var(--ds-inspector-bg);
  border: 1px solid var(--ds-inspector-border);
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  z-index: var(--ds-inspector-z-index);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
  font-family: var(--ds-inspector-font);
  font-size: 12px;
  color: var(--ds-inspector-text);
}

/* ----------------------------------------
   Resize Handle (top-left corner)
   ---------------------------------------- */

.ds-inspector-resize-handle {
  position: absolute;
  top: 0;
  left: 0;
  width: 16px;
  height: 16px;
  cursor: nwse-resize;
  border-radius: 12px 0 0 0;
  z-index: 2;
}

.ds-inspector-resize-handle::before {
  content: '';
  position: absolute;
  top: 4px;
  left: 4px;
  width: 6px;
  height: 6px;
  border-left: 2px solid #6c7086;
  border-top: 2px solid #6c7086;
  transition: border-color 0.15s ease;
}

.ds-inspector-resize-handle:hover::before {
  border-color: var(--ds-inspector-accent);
}

/* ----------------------------------------
   Header (condensed - contains all controls)
   ---------------------------------------- */

.ds-inspector-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  border-bottom: 1px solid var(--ds-inspector-border);
  background: var(--ds-inspector-surface);
  border-radius: 12px 12px 0 0;
}

.ds-inspector-logo {
  color: var(--ds-inspector-accent);
  font-weight: 700;
  font-size: 11px;
  flex-shrink: 0;
}

.ds-inspector-filter-wrapper {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 4px;
  background: var(--ds-inspector-bg);
  border: 1px solid var(--ds-inspector-border);
  border-radius: 4px;
  padding: 0 6px;
  min-width: 0;
}

.ds-inspector-filter-wrapper:focus-within {
  border-color: var(--ds-inspector-accent);
}

.ds-inspector-filter {
  flex: 1;
  background: transparent;
  border: none;
  color: var(--ds-inspector-text);
  font-family: var(--ds-inspector-font);
  font-size: 11px;
  padding: 4px 0;
  outline: none;
  min-width: 0;
}

.ds-inspector-filter::placeholder {
  color: var(--ds-inspector-text-dim);
}

.ds-inspector-filter-clear {
  background: transparent;
  border: none;
  color: var(--ds-inspector-text-dim);
  cursor: pointer;
  padding: 0;
  font-size: 12px;
  line-height: 1;
  font-family: var(--ds-inspector-font);
  flex-shrink: 0;
}

.ds-inspector-filter-clear:hover {
  color: var(--ds-inspector-text);
}

.ds-inspector-view-toggle {
  display: flex;
  background: var(--ds-inspector-bg);
  border: 1px solid var(--ds-inspector-border);
  border-radius: 4px;
  overflow: hidden;
  flex-shrink: 0;
}

.ds-inspector-view-btn {
  padding: 4px 6px;
  border: none;
  background: transparent;
  color: var(--ds-inspector-text-dim);
  font-family: var(--ds-inspector-font);
  font-size: 10px;
  cursor: pointer;
}

.ds-inspector-view-btn:hover {
  color: var(--ds-inspector-text);
}

.ds-inspector-view-btn.active {
  background: var(--ds-inspector-accent);
  color: var(--ds-inspector-bg);
}

.ds-inspector-btn {
  width: 20px;
  height: 20px;
  border: none;
  background: transparent;
  color: var(--ds-inspector-text-dim);
  cursor: pointer;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-family: var(--ds-inspector-font);
  flex-shrink: 0;
}

.ds-inspector-btn:hover {
  background: var(--ds-inspector-border);
  color: var(--ds-inspector-text);
}

/* ----------------------------------------
   Content Area
   ---------------------------------------- */

.ds-inspector-content {
  flex: 1;
  overflow: auto;
  padding: 12px;
  min-height: 0;
}

/* ----------------------------------------
   JSON View
   ---------------------------------------- */

.ds-inspector-json {
  white-space: pre;
  font-size: 11px;
  line-height: 1.5;
  margin: 0;
}

.ds-inspector-key { color: var(--ds-inspector-key); }
.ds-inspector-string { color: var(--ds-inspector-string); }
.ds-inspector-number { color: var(--ds-inspector-number); }
.ds-inspector-boolean { color: var(--ds-inspector-boolean); }
.ds-inspector-null { color: var(--ds-inspector-null); font-style: italic; }

/* JSON line wrapper for flash effect - inline to preserve pre whitespace */
.ds-inspector-line {
  border-radius: 2px;
}

/* ----------------------------------------
   Table View
   ---------------------------------------- */

.ds-inspector-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 11px;
}

.ds-inspector-table th,
.ds-inspector-table td {
  text-align: left;
  padding: 6px 8px;
  border-bottom: 1px solid var(--ds-inspector-border);
}

.ds-inspector-table th {
  color: var(--ds-inspector-text-dim);
  font-weight: 500;
  text-transform: uppercase;
  font-size: 10px;
  position: sticky;
  top: -12px;
  background: var(--ds-inspector-bg);
  z-index: 1;
}

.ds-inspector-table td:first-child {
  color: var(--ds-inspector-key);
  font-weight: 500;
}

.ds-inspector-table-value {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Table row flash */
.ds-inspector-table tr.ds-inspector-flash td {
  animation: ds-inspector-flash ${FLASH_DURATION}ms ease-out;
}

/* ----------------------------------------
   Empty State
   ---------------------------------------- */

.ds-inspector-empty {
  color: var(--ds-inspector-text-dim);
  text-align: center;
  padding: 24px;
  font-style: italic;
}

/* ----------------------------------------
   Utilities
   ---------------------------------------- */

.ds-inspector-hidden {
  display: none !important;
}

/* ----------------------------------------
   Scrollbar
   ---------------------------------------- */

.ds-inspector-content::-webkit-scrollbar {
  width: 8px;
}

.ds-inspector-content::-webkit-scrollbar-track {
  background: transparent;
}

.ds-inspector-content::-webkit-scrollbar-thumb {
  background: var(--ds-inspector-border);
  border-radius: 4px;
}

.ds-inspector-content::-webkit-scrollbar-thumb:hover {
  background: var(--ds-inspector-text-dim);
}
`
