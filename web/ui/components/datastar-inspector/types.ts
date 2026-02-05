/**
 * Type definitions for Datastar Inspector
 */

/**
 * View mode for signal display
 */
export type ViewMode = 'json' | 'table'

/**
 * Persisted inspector state (stored in sessionStorage)
 */
export interface InspectorState {
  /** Whether the panel is expanded */
  expanded: boolean
  /** Current filter text */
  filter: string
  /** Current view mode */
  viewMode: ViewMode
  /** Panel width in pixels */
  panelWidth?: number
  /** Panel height in pixels */
  panelHeight?: number
}

/**
 * Signal data structure (recursive key-value)
 */
export type SignalValue = string | number | boolean | null | SignalValue[] | SignalObject

export interface SignalObject {
  [key: string]: SignalValue
}
