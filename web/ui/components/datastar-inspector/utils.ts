/**
 * Utility functions for Datastar Inspector
 */

import type { SignalObject } from './types.js'

// ============================================
// HTML/Regex Escaping
// ============================================

/**
 * Escape HTML special characters to prevent XSS
 */
export function escapeHtml(str: string): string {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

/**
 * Escape special regex characters in a string
 */
export function escapeRegex(str: string): string {
  return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}

// ============================================
// Signal Counting & Flattening
// ============================================

/**
 * Count the number of leaf signals in an object tree
 */
export function countSignals(obj: unknown, count = 0): number {
  if (typeof obj !== 'object' || obj === null) return count + 1
  for (const value of Object.values(obj as Record<string, unknown>)) {
    count = countSignals(value, count)
  }
  return count
}

/**
 * Flatten nested signals into an array of [path, value] tuples
 */
export function flattenSignals(
  obj: Record<string, unknown>,
  prefix = ''
): Array<[string, unknown]> {
  const result: Array<[string, unknown]> = []

  for (const [key, value] of Object.entries(obj)) {
    const path = prefix ? `${prefix}.${key}` : key

    if (typeof value === 'object' && value !== null && !Array.isArray(value)) {
      result.push(...flattenSignals(value as Record<string, unknown>, path))
    } else {
      result.push([path, value])
    }
  }

  return result
}

// ============================================
// Signal Filtering
// ============================================

/**
 * Parse a filter string into a RegExp
 * Supports: plain text, /regex/, and wildcards (*)
 */
export function parseFilterPattern(filterText: string): RegExp {
  // Check if filter is a regex (starts and ends with /)
  if (filterText.startsWith('/') && filterText.lastIndexOf('/') > 0) {
    const lastSlash = filterText.lastIndexOf('/')
    const pattern = filterText.slice(1, lastSlash)
    const flags = filterText.slice(lastSlash + 1)
    try {
      return new RegExp(pattern, flags || 'i')
    } catch {
      return new RegExp(escapeRegex(filterText), 'i')
    }
  }

  // Wildcard pattern
  if (filterText.includes('*')) {
    const pattern = escapeRegex(filterText).replace(/\\\*/g, '.*')
    return new RegExp(pattern, 'i')
  }

  // Plain text search
  return new RegExp(escapeRegex(filterText), 'i')
}

/**
 * Filter an object tree by a regex pattern
 * Matches against both paths and values
 */
export function filterObject(
  obj: Record<string, unknown>,
  regex: RegExp,
  path = ''
): Record<string, unknown> {
  const result: Record<string, unknown> = {}

  for (const [key, value] of Object.entries(obj)) {
    const fullPath = path ? `${path}.${key}` : key

    if (typeof value === 'object' && value !== null && !Array.isArray(value)) {
      const filtered = filterObject(value as Record<string, unknown>, regex, fullPath)
      if (Object.keys(filtered).length > 0) {
        result[key] = filtered
      }
    } else if (regex.test(fullPath) || regex.test(String(value))) {
      result[key] = value
    }
  }

  return result
}

// ============================================
// Change Detection
// ============================================

/**
 * Find paths that changed between two signal objects
 */
export function findChangedPaths(
  oldObj: SignalObject,
  newObj: SignalObject,
  prefix = ''
): Set<string> {
  const changed = new Set<string>()

  // Check all keys in new object
  for (const [key, newValue] of Object.entries(newObj)) {
    const path = prefix ? `${prefix}.${key}` : key
    const oldValue = oldObj[key]

    if (typeof newValue === 'object' && newValue !== null && !Array.isArray(newValue)) {
      if (typeof oldValue === 'object' && oldValue !== null && !Array.isArray(oldValue)) {
        // Recurse into nested objects
        const nestedChanged = findChangedPaths(
          oldValue as SignalObject,
          newValue as SignalObject,
          path
        )
        nestedChanged.forEach((p) => changed.add(p))
      } else {
        // Type changed from non-object to object
        changed.add(path)
      }
    } else if (JSON.stringify(oldValue) !== JSON.stringify(newValue)) {
      changed.add(path)
    }
  }

  // Check for removed keys
  for (const key of Object.keys(oldObj)) {
    const path = prefix ? `${prefix}.${key}` : key
    if (!(key in newObj)) {
      changed.add(path)
    }
  }

  return changed
}

// ============================================
// JSON Rendering
// ============================================

/**
 * Render a value as syntax-highlighted HTML for JSON display
 */
export function renderJsonValue(
  value: unknown,
  changedPaths: Set<string>,
  indent = 0,
  path = ''
): string {
  const pad = '  '.repeat(indent)

  if (value === null) {
    return `<span class="ds-inspector-null">null</span>`
  }
  if (typeof value === 'boolean') {
    return `<span class="ds-inspector-boolean">${value}</span>`
  }
  if (typeof value === 'number') {
    return `<span class="ds-inspector-number">${value}</span>`
  }
  if (typeof value === 'string') {
    return `<span class="ds-inspector-string">"${escapeHtml(value)}"</span>`
  }
  if (Array.isArray(value)) {
    if (value.length === 0) return '[]'
    const items = value
      .map((v, i) => {
        const itemPath = `${path}[${i}]`
        return `${pad}  ${renderJsonValue(v, changedPaths, indent + 1, itemPath)}`
      })
      .join(',\n')
    return `[\n${items}\n${pad}]`
  }
  if (typeof value === 'object') {
    const entries = Object.entries(value)
    if (entries.length === 0) return '{}'
    const items = entries
      .map(([k, v]) => {
        const keyPath = path ? `${path}.${k}` : k
        const isChanged = changedPaths.has(keyPath)
        const flashClass = isChanged ? ' ds-inspector-flash' : ''
        const lineContent = `<span class="ds-inspector-key">"${escapeHtml(k)}"</span>: ${renderJsonValue(v, changedPaths, indent + 1, keyPath)}`
        return `${pad}  <span class="ds-inspector-line${flashClass}">${lineContent}</span>`
      })
      .join(',\n')
    return `{\n${items}\n${pad}}`
  }
  return String(value)
}
