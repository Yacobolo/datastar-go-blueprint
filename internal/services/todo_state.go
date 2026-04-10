package services

// TodoViewMode controls which todos are visible in the list.
type TodoViewMode int

const (
	// TodoViewModeAll shows every todo item.
	TodoViewModeAll TodoViewMode = iota
	// TodoViewModeActive shows only incomplete todos.
	TodoViewModeActive
	// TodoViewModeCompleted shows only completed todos.
	TodoViewModeCompleted
	// TodoViewModeLast marks the exclusive upper bound for view modes.
	TodoViewModeLast
)

// TodoViewModeStrings maps view modes to their UI labels.
var TodoViewModeStrings = []string{"All", "Active", "Completed"}

// Todo is the state representation of a single todo item.
type Todo struct {
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

// TodoMVC represents the full todo page state stored per session.
type TodoMVC struct {
	Todos      []*Todo      `json:"todos"`
	EditingIdx int          `json:"editingIdx"`
	Mode       TodoViewMode `json:"mode"`
}
