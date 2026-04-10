// Package todocomponents contains gomponents-based UI for the todo feature.
package todocomponents

import (
	"fmt"
	"strings"

	commoncomponents "github.com/yacobolo/datastar-go-blueprint/internal/features/common/components"
	appds "github.com/yacobolo/datastar-go-blueprint/internal/platform/ds"

	g "maragu.dev/gomponents"
	data "maragu.dev/gomponents-datastar"
	h "maragu.dev/gomponents/html"
)

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

// Todo is the UI representation of a single todo item.
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

func classNames(classes ...string) string {
	filtered := make([]string, 0, len(classes))
	for _, class := range classes {
		if class != "" {
			filtered = append(filtered, class)
		}
	}
	return strings.Join(filtered, " ")
}

func when(condition bool, value string) string {
	if condition {
		return value
	}
	return ""
}

// TodosMVCView renders the main todo UI.
func TodosMVCView(mvc *TodoMVC) g.Node {
	hasTodos := len(mvc.Todos) > 0
	left, completed := 0, 0
	for _, todo := range mvc.Todos {
		if !todo.Completed {
			left++
		} else {
			completed++
		}
	}

	input := ""
	if mvc.EditingIdx >= 0 {
		input = mvc.Todos[mvc.EditingIdx].Text
	}

	rows := make(g.Group, 0, len(mvc.Todos))
	for i, todo := range mvc.Todos {
		rows = append(rows, TodoRow(mvc.Mode, todo, i, i == mvc.EditingIdx))
	}

	modeButtons := make(g.Group, 0, int(TodoViewModeLast))
	for i := TodoViewModeAll; i < TodoViewModeLast; i++ {
		if i == mvc.Mode {
			modeButtons = append(modeButtons,
				h.Button(
					h.Type("button"),
					h.Class("btn btn-sm join-item btn-active"),
					g.Text(TodoViewModeStrings[i]),
				),
			)
			continue
		}

		modeButtons = append(modeButtons,
			h.Button(
				h.Type("button"),
				h.Class("btn btn-sm join-item"),
				data.On("click", appds.Putf("/api/todos/mode/%d", i)),
				g.Text(TodoViewModeStrings[i]),
			),
		)
	}

	footerActions := g.Group{}
	if completed > 0 {
		footerActions = append(footerActions,
			h.Button(
				h.Type("button"),
				h.Class("btn btn-ghost btn-sm btn-square text-error"),
				h.Title(fmt.Sprintf("Clear %d completed todos", completed)),
				g.Attr("aria-label", fmt.Sprintf("Clear %d completed todos", completed)),
				data.On("click", appds.Delete("/api/todos/-1")),
				commoncomponents.IconTrash(),
			),
		)
	}
	footerActions = append(footerActions,
		h.Button(
			h.Type("button"),
			h.Class("btn btn-outline btn-sm"),
			h.Title("Reset list"),
			data.On("click", appds.Put("/api/todos/reset")),
			commoncomponents.IconReset(),
			h.Span(g.Text("Reset")),
		),
	)

	return h.Section(
		h.ID("todos-container"),
		h.Class("mx-auto w-full max-w-4xl"),
		h.Div(
			h.Class("card bg-base-100 shadow-xl"),
			data.Signals(map[string]any{"input": input}),
			h.Div(
				h.Class("card-body gap-6"),
				h.Div(
					h.Class("flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"),
					h.Div(
						h.Class("space-y-1"),
						h.P(h.Class("text-sm uppercase tracking-[0.2em] text-base-content/60"), g.Text("Template")),
						h.H1(h.Class("text-4xl font-bold tracking-tight text-primary"), g.Text("todos")),
					),
					h.Div(
						h.Class("flex w-full flex-col gap-3 sm:w-auto sm:flex-row sm:items-center"),
						g.If(hasTodos, h.Button(
							h.ID("toggleAll"),
							h.Type("button"),
							h.Class("btn btn-primary btn-square shrink-0"),
							h.Title("Toggle all todos"),
							g.Attr("aria-label", "Toggle all todos"),
							data.On("click", appds.Post("/api/todos/-1/toggle")),
							data.Indicator("toggleAllFetching"),
							data.Attr("disabled", "$toggleAllFetching"),
							commoncomponents.IconToggleAll(),
						)),
						g.If(mvc.EditingIdx < 0, TodoInput(-1)),
						commoncomponents.SseIndicator("toggleAllFetching"),
					),
				),
				g.If(hasTodos, h.Section(
					h.Class("space-y-3"),
					h.Ul(
						h.Class("space-y-3"),
						rows,
					),
				)),
				g.If(hasTodos, h.Div(
					h.Class("flex flex-col gap-4 border-t border-base-300 pt-4 sm:flex-row sm:items-center sm:justify-between"),
					h.Div(
						h.Class("text-sm text-base-content/70"),
						h.Strong(g.Text(fmt.Sprint(left)+itemLabel(left))),
						g.Text(" left"),
					),
					h.Div(
						h.Class("flex flex-col gap-3 sm:flex-row sm:items-center"),
						h.Div(h.Class("join"), modeButtons),
						h.Div(h.Class("flex items-center gap-2"), footerActions),
					),
				)),
			),
		),
	)
}

// TodoInput renders the add/edit input.
func TodoInput(i int) g.Node {
	expression := fmt.Sprintf(`
			if (evt.key !== 'Enter' || !$input.trim().length) return;
			%s;
			$input = '';
		`, appds.Putf("/api/todos/%d/edit", i))

	return h.Input(
		h.ID("todoInput"),
		g.Attr("data-testid", "todos_input"),
		h.Class("input input-bordered w-full sm:min-w-80"),
		h.Placeholder("What needs to be done?"),
		data.Bind("input"),
		data.On("keydown", expression),
		g.If(i >= 0, data.On("click", appds.Put("/api/todos/cancel"), data.ModifierOutside)),
	)
}

// TodoRow renders a single todo row.
func TodoRow(mode TodoViewMode, todo *Todo, i int, isEditing bool) g.Node {
	indicatorID := fmt.Sprintf("indicator%d", i)
	fetchingSignalName := fmt.Sprintf("fetching%d", i)

	if isEditing {
		return h.Li(
			h.ID(fmt.Sprintf("todo%d", i)),
			h.Class("rounded-box border border-base-300 bg-base-100 p-3 shadow-sm"),
			TodoInput(i),
		)
	}

	if mode != TodoViewModeAll &&
		(mode != TodoViewModeActive || todo.Completed) &&
		(mode != TodoViewModeCompleted || !todo.Completed) {
		return nil
	}

	return h.Li(
		h.Class(classNames(
			"flex items-center gap-3 rounded-box border border-base-300 bg-base-100 px-4 py-3 shadow-sm",
			when(todo.Completed, "opacity-70"),
		)),
		h.ID(fmt.Sprintf("todo%d", i)),
		h.Button(
			h.ID(fmt.Sprintf("toggle%d", i)),
			h.Type("button"),
			h.Class("btn btn-ghost btn-sm btn-circle"),
			g.Attr("aria-label", fmt.Sprintf("Toggle todo %d", i+1)),
			g.Attr("aria-checked", fmt.Sprintf("%t", todo.Completed)),
			data.On("click", appds.Postf("/api/todos/%d/toggle", i)),
			data.On(
				"keydown",
				fmt.Sprintf(
					"if (evt.key === 'Enter' || evt.key === ' ') { evt.preventDefault(); %s }",
					appds.Postf("/api/todos/%d/toggle", i),
				),
			),
			data.Indicator(fetchingSignalName),
			commoncomponents.IconCheckbox(todo.Completed),
		),
		h.Button(
			h.ID(indicatorID),
			h.Type("button"),
			h.Class(classNames(
				"btn btn-ghost h-auto min-h-0 flex-1 justify-start px-3 py-2 text-left text-sm font-normal normal-case sm:text-base",
				when(todo.Completed, "line-through text-base-content/50"),
			)),
			h.Title("Edit todo"),
			data.On("click", appds.Getf("/api/todos/%d/edit", i)),
			data.On(
				"keydown",
				fmt.Sprintf(
					"if (evt.key === 'Enter') { evt.preventDefault(); %s }",
					appds.Getf("/api/todos/%d/edit", i),
				),
			),
			data.Indicator(fetchingSignalName),
			g.Text(todo.Text),
		),
		commoncomponents.SseIndicator(fetchingSignalName),
		h.Button(
			h.ID(fmt.Sprintf("delete%d", i)),
			h.Type("button"),
			h.Class("btn btn-ghost btn-sm btn-square text-error"),
			h.Title("Delete todo"),
			g.Attr("aria-label", fmt.Sprintf("Delete todo %d", i+1)),
			data.On("click", appds.Deletef("/api/todos/%d", i)),
			data.Indicator(fetchingSignalName),
			data.Attr("disabled", fmt.Sprintf("$%s", fetchingSignalName)),
			g.Attr("data-testid", fmt.Sprintf("delete_todo%d", i)),
			commoncomponents.IconClose(),
		),
	)
}

func itemLabel(total int) string {
	if total > 1 {
		return " items"
	}
	return " item"
}
