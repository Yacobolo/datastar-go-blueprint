// Import theme switcher (auto-initializes)
import '../src/plugins/theme-switcher';

// Export all Lit components
export { TodoTable } from './todo-table/todo-table';
export { DatastarInspector } from './datastar-inspector/datastar-inspector';

// Re-export Lit for convenience
export { LitElement, html, css } from 'lit';
export type { TemplateResult } from 'lit';