# Datastar Go Starter Kit

A modern, opinionated starter kit for building reactive web applications with Go, Datastar, SQLite, and Lit web components.

## Stack

**Backend:**

- [Go](https://go.dev/doc/) - Type-safe backend with excellent tooling
- [SQLite](https://www.sqlite.org/) + [sqlc](https://sqlc.dev/) - Embedded database with type-safe queries
- [Datastar](https://github.com/starfederation/datastar) - Hypermedia-driven reactivity
- [datastar-templ](https://github.com/Yacobolo/datastar-templ) - Type-safe Datastar attribute helpers for templ
- [Templ](https://templ.guide/) - Type-safe HTML templating
- [Chi](https://github.com/go-chi/chi) - Lightweight, composable HTTP router

**Frontend:**

- [Lit](https://lit.dev/) - Fast, modern web components
- Native CSS with CSS layers, nesting, and custom properties
- [esbuild](https://esbuild.github.io/) - Lightning-fast JavaScript bundling
- [Datastar Client](https://www.jsdelivr.com/package/gh/starfederation/datastar) - Reactive signals and SSE

## Features

- **Feature-based architecture** - Organize code by feature, not layer
- **Type-safe everything** - Go types, SQL queries (sqlc), and HTML templates (templ)
- **Hot reload** - Fast development with Air + esbuild + templ watchers
- **Modern CSS** - No framework needed, uses CSS layers, nesting, and variables
- **Reactive components** - Lit web components integrated with Datastar signals
- **Single binary** - Static assets embedded, no external dependencies
- **Production ready** - Docker support, optimized builds with UPX compression

## Setup

1. Clone this repository

```shell
git clone https://github.com/yacobolo/datastar-go-starter-kit.git
cd datastar-go-starter-kit
```

2. Install Dependencies

```shell
go mod tidy
```

3. Install development tools

```shell
task tools:install
# Installs: templ, air, sqlc, golangci-lint, cssgen
```

Or install individually:
```shell
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/yacobolo/cssgen/cmd/cssgen@latest
```

4. Generate code

```shell
task generate:all  # Generates database code, CSS constants, and templ files
```

5. Start developing!

```shell
go tool task live
```

## Development

Live reload is set up out of the box, powered by:

- [Air](https://github.com/air-verse/air) - Go hot reload
- [esbuild](cmd/web/build/main.go) - TypeScript/JavaScript bundling
- [templ](https://templ.guide/) - Template hot reload

Use the [live task](./Taskfile.yml) from the [Taskfile](https://taskfile.dev/) to start with live reload:

```shell
go tool task live
```

Navigate to [`http://localhost:8080`](http://localhost:8080) in your browser.

### Available Tasks

```shell
# Development
go tool task live           # Start with hot reload
go tool task debug          # Start with debugger attached

# Building
go tool task build          # Production build
go tool task build:templ    # Generate Go from .templ files
go tool task build:web:bundle # Bundle TypeScript components

# Code Generation
sqlc generate               # Generate type-safe database code
```

## Project Structure

```
.
├── cmd/web/              # Application entry point
│   ├── main.go          # Server initialization
│   └── build/           # esbuild bundler for web components
├── config/              # Environment configuration
├── db/                  # Database layer
│   ├── schema.sql       # SQLite schema
│   ├── queries.sql      # SQL queries for sqlc
│   └── db.go           # Database initialization
├── features/            # Feature-based modules
│   ├── common/         # Shared layouts & components
│   └── index/          # TODO app feature
│       ├── components/  # Templ components
│       ├── pages/      # Page templates
│       ├── services/   # Business logic
│       ├── handlers.go # HTTP handlers
│       └── routes.go   # Route registration
├── router/             # HTTP router setup
├── web/
│   ├── libs/lit/       # Lit web components
│   └── resources/      # Static assets & CSS
└── Taskfile.yml        # Build automation
```

## Database

This starter kit uses **SQLite** with **sqlc** for type-safe database queries.

### Why SQLite + sqlc?

- **Zero configuration** - No database server to set up
- **Type-safe queries** - sqlc generates Go code from SQL
- **Perfect for starter projects** - Easy to scale to PostgreSQL/MySQL later
- **Single file** - `data/todos.db` contains all your data

### Working with the database

1. Edit `db/schema.sql` to modify the database schema
2. Edit `db/queries.sql` to add/modify queries
3. Run `sqlc generate` to generate Go code
4. Use the generated code in your services: `queries.GetTodosByUser(ctx, userID)`

### Viewing data

```shell
sqlite3 data/todos.db "SELECT * FROM todos;"
```

## Styling with Native CSS

This starter kit uses **modern native CSS** - no Tailwind, no frameworks.

### CSS Architecture

- `styles/reset.css` - Modern CSS reset
- `styles/tokens.css` - CSS custom properties (design tokens)
- `styles/base.css` - Base element styles using @layer
- `styles/components.css` - Reusable component patterns with nesting
- `styles/utilities.css` - Minimal utility classes
- `styles/main.css` - Entry point that imports all styles

### Key Features

- **CSS Layers** - Proper cascade control with `@layer`
- **CSS Nesting** - Write cleaner, more maintainable styles
- **CSS Variables** - Theme tokens like `--color-primary`, `--space-md`
- **Component-scoped styles** - Lit components use Shadow DOM

### Example

```css
/* styles/components.css */
@layer components {
  .btn {
    padding: var(--space-sm) var(--space-md);
    background: var(--color-primary);
    border-radius: var(--radius-md);

    &:hover {
      background: var(--color-primary-dark);
    }
  }
}
```

## Type-Safe CSS with cssgen

This starter kit includes **cssgen**, a tool that generates Go constants from CSS classes and provides build-time validation.

### Why Type-Safe CSS?

- **Catch typos at build time** - `class="btn btn--primray"` fails compilation
- **IDE autocomplete** - Get suggestions for all available CSS classes
- **Refactor with confidence** - Rename `.btn-primary` and find all usages
- **Zero runtime overhead** - Pure compile-time tool

### Usage

**Before (hardcoded strings):**
```go
<button class="btn btn-primary btn-lg">Click</button>
```

**After (type-safe constants):**
```go
import "github.com/yacobolo/datastar-go-starter-kit/internal/ui"

<button class={ ui.Btn, ui.BtnPrimary, ui.BtnLg }>Click</button>
```

### Available Commands

```shell
# Generate CSS constants from your CSS files
task css:gen

# Lint CSS usage in templ files (golangci-lint style)
task css:lint

# Full report with statistics and Quick Wins
task css:lint:full

# Weekly adoption report
task css:report

# Export Markdown report for documentation
task css:report:md
```

### Integration

CSS constants are automatically generated when you run:

```shell
task generate:all  # Runs sqlc, css:gen, and templ generation
task dev           # Auto-regenerates on CSS file changes
```

### Linting in CI

The `check` task includes CSS linting:

```shell
task check  # Runs tests, Go linter, and CSS linter
```

For more details, see the [cssgen documentation](https://github.com/Yacobolo/cssgen).

### Example: Todo Component Refactoring

The todo component demonstrates a complete migration from Tailwind-style utility combinations to semantic, native CSS classes.

**Before (Tailwind-style utilities):**
```go
<div class="flex flex-col gap-sm">
    <div class="flex items-center gap-sm justify-center">
        <h1 class="text-4xl font-bold text-primary">TODO</h1>
    </div>
    <div class="todo-info-box p-md rounded-md" style="background: var(--ui-color-surface);">
        <span class="italic font-bold uppercase text-primary">single get request!</span>
    </div>
</div>
```

**After (Semantic native CSS with cssgen):**
```go
import "github.com/yacobolo/datastar-go-starter-kit/internal/ui"

<div class={ ui.TodoHeader }>
    <div class={ ui.TodoTitleSection }>
        <h1 class={ ui.TodoTitle }>TODO</h1>
    </div>
    <div class={ ui.Callout, ui.CalloutInfo }>
        <div class={ ui.CalloutContent }>
            <span class={ ui.Italic, ui.FontBold, ui.Uppercase, ui.TextPrimary }>single get request!</span>
        </div>
    </div>
</div>
```

**Benefits:**
- ✅ Self-documenting class names (`.todo-title` vs `.text-4xl .font-bold .text-primary`)
- ✅ Consistent with design system (all values use CSS variables)
- ✅ Type-safe with IDE autocomplete
- ✅ Easier to maintain and refactor
- ✅ No inline styles needed
- ✅ Reusable components (`.callout` works across features)

**Results:**
- Adoption rate increased from 5.6% to 42.0%
- All inline styles eliminated from todo.templ
- Created 15+ semantic component classes
- Added generic `.callout` component for reuse

## Type-Safe Datastar Attributes with datastar-templ

This starter kit uses **datastar-templ** for compile-time type safety when working with Datastar attributes in templ templates.

### Why datastar-templ?

- **Compile-time type checking** - Catch errors before runtime
- **IDE autocomplete** - Get suggestions for all Datastar attributes and modifiers
- **Consistent API** - Clean, predictable functions for all Datastar features
- **Zero runtime overhead** - Pure compile-time helpers

### Usage

Import the library (commonly aliased as `ds`):

```go
import ds "github.com/Yacobolo/datastar-templ"
```

**Before (inline strings):**
```go
<div data-signals={ fmt.Sprintf("{count: %d}", count) }>
  <button data-on:click={ datastar.PostSSE("/increment") }>+</button>
  <span data-text="$count"></span>
</div>
```

**After (type-safe):**
```go
<div { ds.Signals(ds.Int("count", count))... }>
  <button { ds.OnClick(ds.Post("/increment"))... }>+</button>
  <span { ds.Text("$count")... }></span>
</div>
```

### Common Patterns

**Signals:**
```go
{ ds.Signals(
    ds.String("name", ""),
    ds.Int("count", 0),
    ds.Bool("isOpen", false),
    ds.JSON("user", userData),
)... }
```

**Events:**
```go
{ ds.OnClick(ds.Post("/submit"))... }
{ ds.OnInput(ds.Get("/search?q=$query"), ds.ModDebounce, ds.Ms(300))... }
{ ds.OnEvent("custom-event", "$handler()")... }
```

**Bindings:**
```go
{ ds.Bind("email")... }
{ ds.Text("$message")... }
{ ds.Show("$isVisible")... }
```

**Multiple attributes:**
```go
{ ds.Merge(
    ds.OnClick(ds.Post("/submit")),
    ds.Indicator("loading"),
    ds.Attr(ds.Pair("disabled", "$loading")),
)... }
```

For more details, see the [datastar-templ documentation](https://github.com/Yacobolo/datastar-templ).

## Lit Web Components

The TODO table is built as a Lit web component that integrates seamlessly with Datastar.

### Creating Components

```typescript
// web/libs/lit/src/my-component.ts
import { LitElement, html, css } from "lit";
import { customElement, property } from "lit/decorators.js";

@customElement("my-component")
export class MyComponent extends LitElement {
  @property({ type: String }) message = "";

  static styles = css`
    :host {
      display: block;
      padding: var(--space-md);
    }
  `;

  render() {
    return html`<div>${this.message}</div>`;
  }
}
```

### Integrating with Datastar

```html
<div data-store='{"message": "Hello"}'>
  <my-component data-bind-message="message"></my-component>
</div>
```

## Debugging

The [debug task](./Taskfile.yml) will launch [delve](https://github.com/go-delve/delve) for debugging:

```shell
go tool task debug
```

### Visual Studio Code Integration

A `Debug Main` configuration is included in [.vscode/launch.json](./.vscode/launch.json).

## Deployment

### Building an Executable

```shell
go tool task build
```

This creates a single binary with all assets embedded.

### Docker

```shell
# Build image
docker build -t datastar-go-starter-kit:latest .

# Run container
docker run --name datastar-app -p 8080:9001 datastar-go-starter-kit:latest
```

The [Dockerfile](./Dockerfile) uses multi-stage builds and UPX compression for minimal image size.

## Contributing

Pull requests and feature requests are welcome!

## References

### Backend

- [Go](https://go.dev/)
- [SQLite](https://www.sqlite.org/)
- [sqlc](https://sqlc.dev/)
- [Datastar SDK](https://github.com/starfederation/datastar/tree/develop/sdk)
- [datastar-templ](https://github.com/Yacobolo/datastar-templ)
- [Templ](https://templ.guide/)
- [Chi Router](https://github.com/go-chi/chi)

### Frontend

- [Datastar Client](https://www.jsdelivr.com/package/gh/starfederation/datastar)
- [Lit](https://lit.dev/)
- [esbuild](https://esbuild.github.io/)
- [Modern CSS](https://developer.mozilla.org/en-US/docs/Web/CSS)

## License

MIT
