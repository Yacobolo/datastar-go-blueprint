<h1 align="center">Datastar Go Blueprint</h1>

<p align="center">
  <strong>A modern, opinionated blueprint for building reactive web applications.</strong>
  <br>
  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square" alt="License: MIT">
  </a>
</p>

<p align="center">
  <img src="mascot.png" alt="Datastar Go Blueprint mascot" width="600">
</p>

---

## The Stack

- **Backend:** [Go](https://go.dev/), [SQLite](https://www.sqlite.org/) + [sqlc](https://sqlc.dev/), [Datastar](https://github.com/starfederation/datastar), [datastar-templ](https://github.com/Yacobolo/datastar-templ), [Templ](https://templ.guide/), [Chi](https://github.com/go-chi/chi)
- **Frontend:** [Lit](https://lit.dev/), Native CSS (Layers/Nesting), [esbuild](https://esbuild.github.io/)
- **Tools:** [Air](https://github.com/air-verse/air), [Hivemind](https://github.com/DarthSim/hivemind), [UPX](https://upx.github.io/)

---

## Quick Start

1. **Clone & Install**

   Shell

   ```
   git clone https://github.com/yacobolo/datastar-go-blueprint.git
   cd datastar-go-blueprint
   go mod tidy
   ```

2. **Environment Setup**

   Shell

   ```
   task tools:install  # Installs templ, air, sqlc, golangci-lint, hivemind
   task generate:all   # Generates SQL, CSS constants, and Templ files
   ```

3. **Run Development Server**

   Shell

   ```
   task dev            # Starts Air, Templ watcher, and esbuild via Hivemind
   ```

   Visit [`http://localhost:8080`](https://www.google.com/search?q=http://localhost:8080).

---

## Project Structure

```
.
├── cmd/server/          # App entry point & build scripts
├── data/                # Local SQLite database files
├── internal/
│   ├── app/             # Application lifecycle & initialization
│   ├── domain/          # Core business logic and entities
│   ├── features/        # Feature-based modules (Templ, Handlers, Routes)
│   ├── platform/        # Shared infra (Router, PubSub)
│   ├── store/           # Database layer (Migrations, SQLC, Repositories)
│   └── ui/              # Generated type-safe CSS constants (cssgen)
├── web/
│   ├── resources/       # Static assets & embedded Go files
│   └── ui/              # Frontend source (Lit components, CSS Layers, TS)
├── Procfile.dev         # Dev process management (Hivemind)
├── Taskfile.yml         # Project automation tasks
└── sqlc.yaml            # SQL compiler configuration
```

---

## Key Features

- **Feature-Based Architecture:** Logic grouped by domain, not layer.
- **Type-Safe Everything:** \* [sqlc](https://sqlc.dev/) for database queries.
  - [cssgen](https://github.com/Yacobolo/cssgen) for type-safe CSS classes.
  - [datastar-templ](https://github.com/Yacobolo/datastar-templ) for Datastar attributes.
- **Native CSS:** No framework. Uses standard CSS `@layer` and nesting.
- **Single Binary:** Static assets embedded using `go:embed`.

---

## Essential Commands

| **Command**         | **Description**                             |
| ------------------- | ------------------------------------------- |
| `task dev`          | Start full hot-reload dev environment       |
| `task generate:all` | Run all code generators (SQL, Templ, CSS)   |
| `task build`        | Create a production-ready compressed binary |
| `task check`        | Run tests and linters                       |
| `task docker:build` | Build optimized Docker image                |

---

## Deployment

**Docker Quick Start:**

Shell

```
docker build -t datastar-app .
docker run -p 8080:9001 datastar-app
```

**Binary Build:**

Executables are optimized with UPX compression for minimal footprint. Use `task build`.

---

**License:** MIT
