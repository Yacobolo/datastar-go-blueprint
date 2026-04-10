import { LitElement, html } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import type { Todo } from './types';

/**
 * A Lit web component for rendering a TODO table.
 * Demonstrates integration with Datastar for reactive updates.
 *
 * Usage:
 * <todo-table
 *   data-bind-todos="$todos"
 *   data-bind-mode="$mode">
 * </todo-table>
 */
@customElement('todo-table')
export class TodoTable extends LitElement {
  @property({ type: Array }) todos: Todo[] = [];
  @property({ type: String }) mode: 'all' | 'active' | 'completed' = 'all';

  override createRenderRoot() {
    return this;
  }

  get filteredTodos() {
    switch (this.mode) {
      case 'active':
        return this.todos.filter(t => !t.completed);
      case 'completed':
        return this.todos.filter(t => t.completed);
      default:
        return this.todos;
    }
  }

  render() {
    const filtered = this.filteredTodos;

    if (filtered.length === 0) {
      return html`
        <div class="alert alert-info shadow-sm">
          <span>
            ${this.mode === 'all'
              ? 'No todos yet. Add one to get started!'
              : `No ${this.mode} todos.`}
          </span>
        </div>
      `;
    }

    return html`
      <div class="space-y-3">
        ${filtered.map((todo, index) => this.renderTodoItem(todo, index))}
      </div>
    `;
  }

  renderTodoItem(todo: Todo, index: number) {
    return html`
      <div class="flex items-center gap-3 rounded-box border border-base-300 bg-base-100 px-4 py-3 shadow-sm ${todo.completed ? 'opacity-70' : ''}">
        <button
          class="btn btn-ghost btn-sm btn-circle"
          type="button"
          aria-label="Toggle todo"
          data-on-click="@post('/api/todos/${index}/toggle')"
        >
          <span class="${todo.completed ? 'text-primary' : 'text-base-content/50'}">
            ${todo.completed ? '☑' : '☐'}
          </span>
        </button>
        <span class="flex-1 text-sm sm:text-base ${todo.completed ? 'line-through text-base-content/50' : ''}">
          ${todo.task}
        </span>
        <div class="flex items-center gap-2">
          <button
            class="btn btn-ghost btn-xs"
            type="button"
            data-on-click="@get('/api/todos/${index}/edit')"
          >
            Edit
          </button>
          <button
            class="btn btn-ghost btn-xs text-error"
            type="button"
            data-on-click="@delete('/api/todos/${index}')"
          >
            Delete
          </button>
        </div>
      </div>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'todo-table': TodoTable;
  }
}
