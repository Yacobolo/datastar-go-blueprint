export interface Todo {
  id: string;
  task: string;
  completed: boolean;
}

export interface TodoStore {
  todos: Todo[];
  mode: 'all' | 'active' | 'completed';
  editingIdx: number;
  input: string;
}
