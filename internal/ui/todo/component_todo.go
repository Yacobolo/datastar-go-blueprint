// Package todo contains gomponents-based UI for the todo feature.
package todo

import (
	"fmt"

	appds "github.com/yacobolo/datastar-go-blueprint/internal/platform/ds"
	"github.com/yacobolo/datastar-go-blueprint/internal/services"
	"github.com/yacobolo/datastar-go-blueprint/internal/ui/core"

	g "maragu.dev/gomponents"
	data "maragu.dev/gomponents-datastar"
	h "maragu.dev/gomponents/html"
)

// ViewMode controls which todos are visible in the list.
type ViewMode = services.TodoViewMode

const (
	// ViewModeAll shows every todo item.
	ViewModeAll = services.TodoViewModeAll
	// ViewModeActive shows only incomplete todos.
	ViewModeActive = services.TodoViewModeActive
	// ViewModeCompleted shows only completed todos.
	ViewModeCompleted = services.TodoViewModeCompleted
	// ViewModeLast marks the exclusive upper bound for view modes.
	ViewModeLast = services.TodoViewModeLast
)

// ViewModeStrings maps view modes to their UI labels.
var ViewModeStrings = services.TodoViewModeStrings

// Todo is the UI representation of a single todo item.
type Todo = services.Todo

// MVC represents the full todo page state stored per session.
type MVC = services.TodoMVC

// TodosMVCView renders the main todo UI.
func TodosMVCView(mvc *MVC) g.Node {
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
		rows = append(rows, Row(mvc.Mode, todo, i, i == mvc.EditingIdx))
	}

	modeButtons := make(g.Group, 0, int(ViewModeLast))
	for i := ViewModeAll; i < ViewModeLast; i++ {
		if i == mvc.Mode {
			modeButtons = append(modeButtons,
				h.Button(
					h.Type("button"),
					h.Class("btn btn-sm join-item btn-primary"),
					g.Text(ViewModeStrings[i]),
				),
			)
			continue
		}

		modeButtons = append(modeButtons,
			h.Button(
				h.Type("button"),
				h.Class("btn btn-sm join-item"),
				data.On("click", appds.Putf("/api/todos/mode/%d", i)),
				g.Text(ViewModeStrings[i]),
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
				core.IconTrash(),
			),
		)
	}
	footerActions = append(footerActions,
		h.Button(
			h.Type("button"),
			h.Class("btn btn-outline btn-sm"),
			h.Title("Reset list"),
			data.On("click", appds.Put("/api/todos/reset")),
			core.IconListChecks(),
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
						h.P(
							h.Class("max-w-xl text-sm text-base-content/70"),
							g.Text("A small server-rendered TodoMVC starter built with Go, Datastar, Tailwind, and DaisyUI."),
						),
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
							core.IconListChecks(),
						)),
						g.If(mvc.EditingIdx < 0, Input(-1)),
						core.SseIndicator("toggleAllFetching"),
					),
				),
				g.If(hasTodos, h.Section(
					h.Class("space-y-3"),
					h.Ul(
						h.Class("space-y-3"),
						rows,
					),
				)),
				g.If(!hasTodos, h.Div(
					h.Class("hero rounded-box border border-dashed border-base-300 bg-base-200/60 py-8"),
					h.Div(
						h.Class("hero-content text-center"),
						h.Div(
							h.Class("max-w-md space-y-3"),
							h.Div(
								h.Class("flex justify-center"),
								h.Span(h.Class("badge badge-primary badge-outline"), g.Text("Ready")),
							),
							h.H2(h.Class("text-2xl font-semibold"), g.Text("Start with your first todo")),
							h.P(
								h.Class("text-sm text-base-content/70"),
								g.Text("Type something above and press Enter to seed the list. The rest of the UI updates live from the server."),
							),
						),
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

// Input renders the add/edit input.
func Input(i int) g.Node {
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

// Row renders a single todo row.
func Row(mode ViewMode, todo *Todo, i int, isEditing bool) g.Node {
	indicatorID := fmt.Sprintf("indicator%d", i)
	fetchingSignalName := fmt.Sprintf("fetching%d", i)
	rowClass := "flex items-center gap-3 rounded-box border border-base-300 bg-base-100 px-4 py-3 shadow-sm"
	labelClass := "btn btn-ghost h-auto min-h-0 flex-1 justify-start px-3 py-2 text-left text-sm font-normal normal-case sm:text-base"

	if isEditing {
		return h.Li(
			h.ID(fmt.Sprintf("todo%d", i)),
			h.Class("rounded-box border border-primary/30 bg-base-100 p-3 shadow-sm"),
			h.Div(
				h.Class("space-y-3"),
				h.Div(
					h.Class("flex items-center justify-between gap-3"),
					h.Span(h.Class("badge badge-primary badge-outline"), g.Text("Editing")),
					h.Span(h.Class("text-xs text-base-content/60"), g.Text("Press Enter to save")),
				),
				Input(i),
			),
		)
	}

	if mode != ViewModeAll &&
		(mode != ViewModeActive || todo.Completed) &&
		(mode != ViewModeCompleted || !todo.Completed) {
		return nil
	}

	if todo.Completed {
		rowClass += " opacity-70"
		labelClass += " line-through text-base-content/50"
	}

	return h.Li(
		h.Class(rowClass),
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
			core.IconCheckbox(todo.Completed),
		),
		h.Button(
			h.ID(indicatorID),
			h.Type("button"),
			h.Class(labelClass),
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
		core.SseIndicator(fetchingSignalName),
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
			core.IconClose(),
		),
	)
}

func itemLabel(total int) string {
	if total > 1 {
		return " items"
	}
	return " item"
}
