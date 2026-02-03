import { LitElement, html, css } from 'lit';
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
  
  static styles = css`
    :host {
      display: block;
    }

    .todo-list {
      display: flex;
      flex-direction: column;
      gap: var(--ui-space-sm);
    }

    .todo-item {
      display: flex;
      align-items: center;
      gap: var(--ui-space-md);
      padding: var(--ui-space-md);
      background: var(--ui-color-surface);
      border: var(--ui-border-md) solid var(--ui-color-outline);
      border-radius: var(--ui-radius-md);
      transition: background var(--ui-duration-fast) var(--ui-ease-default);
    }

    .todo-item:hover {
      background: var(--ui-color-surface-variant);
    }

    .todo-item.completed {
      opacity: 0.6;
    }

    .todo-checkbox {
      width: 20px;
      height: 20px;
      cursor: pointer;
    }

    .todo-task {
      flex: 1;
      font-size: var(--ui-type-size-base);
      color: var(--ui-color-on-surface);
    }

    .todo-task.completed {
      text-decoration: line-through;
    }

    .todo-actions {
      display: flex;
      gap: var(--ui-space-sm);
    }

    .btn {
      padding: var(--ui-space-xs) var(--ui-space-sm);
      font-size: var(--ui-type-size-xs);
      border: none;
      border-radius: var(--ui-radius-md);
      cursor: pointer;
      transition: opacity var(--ui-duration-fast);
    }

    .btn:hover {
      opacity: 0.8;
    }

    .btn-error {
      background: var(--ui-color-error);
      color: var(--ui-color-on-error);
    }

    .btn-ghost {
      background: transparent;
      color: var(--ui-color-on-background);
    }

    .empty-state {
      padding: var(--ui-space-xl);
      text-align: center;
      color: var(--ui-color-on-surface-variant);
    }
  `;

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
        <div class="empty-state">
          ${this.mode === 'all' 
            ? 'No todos yet. Add one to get started!' 
            : `No ${this.mode} todos.`}
        </div>
      `;
    }

    return html`
      <div class="todo-list">
        ${filtered.map((todo, index) => this.renderTodoItem(todo, index))}
      </div>
    `;
  }

  renderTodoItem(todo: Todo, index: number) {
    return html`
      <div class="todo-item ${todo.completed ? 'completed' : ''}">
        <input
          type="checkbox"
          class="todo-checkbox"
          .checked=${todo.completed}
          data-on-change="$$post('/api/todos/${index}/toggle')"
        />
        <span class="todo-task ${todo.completed ? 'completed' : ''}">
          ${todo.task}
        </span>
        <div class="todo-actions">
          <button
            class="btn btn-ghost"
            data-on-click="$$post('/api/todos/${index}/start-edit')"
          >
            Edit
          </button>
          <button
            class="btn btn-error"
            data-on-click="$$delete('/api/todos/${index}')"
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
