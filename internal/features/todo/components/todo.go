// Package todocomponents contains gomponents-based UI for the todo feature.
package todocomponents

import (
	"fmt"
	"strings"

	commoncomponents "github.com/yacobolo/datastar-go-blueprint/internal/features/common/components"
	appds "github.com/yacobolo/datastar-go-blueprint/internal/platform/ds"
	"github.com/yacobolo/datastar-go-blueprint/internal/ui"

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
				h.Div(
					h.Class(classNames(ui.Btn, ui.BtnSm, ui.BtnPrimary)),
					g.Text(TodoViewModeStrings[i]),
				),
			)
			continue
		}

		modeButtons = append(modeButtons,
			h.Button(
				h.Class(classNames(ui.Btn, ui.BtnSm, ui.BtnGhost)),
				data.On("click", appds.Putf("/api/todos/mode/%d", i)),
				g.Text(TodoViewModeStrings[i]),
			),
		)
	}

	footerActions := g.Group{}
	if completed > 0 {
		footerActions = append(footerActions,
			h.Div(
				h.Title(fmt.Sprintf("clear %d completed todos", completed)),
				h.Button(
					h.Class(classNames(ui.Btn, ui.BtnSm, ui.BtnError)),
					data.On("click", appds.Delete("/api/todos/-1")),
					commoncomponents.Icon("material-symbols:delete"),
				),
			),
		)
	}
	footerActions = append(footerActions,
		h.Div(
			h.Title("Reset list"),
			h.Button(
				h.Class(classNames(ui.Btn, ui.BtnSm, ui.BtnSecondary)),
				data.On("click", appds.Put("/api/todos/reset")),
				commoncomponents.Icon("material-symbols:delete-sweep"),
			),
		),
	)

	return h.Div(
		h.ID("todos-container"),
		h.Class(ui.TodoContainer),
		h.Div(
			h.Class(ui.TodoContent),
			data.Signals(map[string]any{"input": input}),
			h.Section(
				h.Class(ui.TodoHeader),
				h.Header(
					h.Class(ui.TodoHeader),
					h.Div(
						h.Class(ui.TodoTitleSection),
						h.H1(h.Class(ui.TodoTitle), g.Text("todos")),
					),
					h.Div(
						h.Class(ui.TodoInputControls),
						g.If(hasTodos, h.Div(
							h.Title("toggle all todos"),
							h.Button(
								h.ID("toggleAll"),
								h.Class(classNames(ui.Btn, ui.BtnLg, ui.BtnPrimary)),
								data.On("click", appds.Post("/api/todos/-1/toggle")),
								data.Indicator("toggleAllFetching"),
								data.Attr("disabled", "$toggleAllFetching"),
								commoncomponents.Icon("material-symbols:checklist"),
							),
						)),
						g.If(mvc.EditingIdx < 0, TodoInput(-1)),
						commoncomponents.SseIndicator("toggleAllFetching"),
					),
				),
				g.If(hasTodos, h.Section(
					h.Class(ui.TodoListContainer),
					h.Ul(
						h.Class(ui.TodoList),
						rows,
					),
				)),
				g.If(hasTodos, h.Footer(
					h.Class(ui.TodoFooter),
					h.Span(
						h.Class(ui.TodoFooterCount),
						h.Strong(
							g.Text(fmt.Sprint(left)),
							g.Text(itemLabel(len(mvc.Todos))),
						),
						g.Text(" left"),
					),
					h.Div(
						h.Class(ui.TodoFooterActions),
						modeButtons,
					),
					h.Div(
						h.Class(ui.TodoFooterActions),
						footerActions,
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
		h.Class(classNames(ui.TodoInput, ui.Input)),
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
		return TodoInput(i)
	}

	if mode != TodoViewModeAll &&
		(mode != TodoViewModeActive || todo.Completed) &&
		(mode != TodoViewModeCompleted || !todo.Completed) {
		return nil
	}

	iconName := "material-symbols:check-box-outline-blank"
	if todo.Completed {
		iconName = "material-symbols:check-box-outline"
	}

	return h.Li(
		h.Class(ui.TodoItem),
		h.ID(fmt.Sprintf("todo%d", i)),
		h.Label(
			h.ID(fmt.Sprintf("toggle%d", i)),
			h.Class(ui.TodoCheckboxLabel),
			h.TabIndex("0"),
			h.Role("checkbox"),
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
			commoncomponents.Icon(iconName),
		),
		h.Label(
			h.ID(indicatorID),
			h.Class(ui.TodoTextLabel),
			h.TabIndex("0"),
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
			h.Class(classNames(ui.Btn, ui.BtnSm, ui.BtnError)),
			data.On("click", appds.Deletef("/api/todos/%d", i)),
			data.Indicator(fetchingSignalName),
			data.Attr("disabled", fmt.Sprintf("$%s", fetchingSignalName)),
			g.Attr("data-testid", fmt.Sprintf("delete_todo%d", i)),
			commoncomponents.Icon("material-symbols:close"),
		),
	)
}

func itemLabel(total int) string {
	if total > 1 {
		return " items"
	}
	return " item"
}
